/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <string>
#include <vector>
#include <time.h>

#include <google/protobuf/repeated_field.h>
#include <google/protobuf/timestamp.pb.h>
#include <google/protobuf/util/time_util.h>
#include <grpcpp/channel.h>

#include "LocalEnforcer.h"
#include "ServiceRegistrySingleton.h"
#include "magma_logging.h"

namespace {

std::chrono::milliseconds time_difference_from_now(
  const google::protobuf::Timestamp &timestamp)
{
  auto rule_time_sec =
    google::protobuf::util::TimeUtil::TimestampToSeconds(timestamp);
  auto now = time(NULL);
  auto delta = std::max(rule_time_sec - now, 0L);
  std::chrono::seconds sec(delta);
  return std::chrono::duration_cast<std::chrono::milliseconds>(sec);
}
} // namespace

namespace magma {

uint32_t LocalEnforcer::REDIRECT_FLOW_PRIORITY = 2000;

using google::protobuf::RepeatedPtrField;
using google::protobuf::util::TimeUtil;

// We will treat rule install/uninstall failures as all-or-nothing - that is,
// if we get a bad response from the pipelined client, we'll mark all the rules
// as failed in the response
static void mark_rule_failures(
  const bool &activate_success,
  const bool &deactivate_success,
  const PolicyReAuthRequest &request,
  PolicyReAuthAnswer &answer_out);

LocalEnforcer::LocalEnforcer(
  std::shared_ptr<SessionCloudReporter> reporter,
  std::shared_ptr<StaticRuleStore> rule_store,
  std::shared_ptr<PipelinedClient> pipelined_client,
  std::shared_ptr<aaa::AAAClient> aaa_client,
  long session_force_termination_timeout_ms):
  reporter_(reporter),
  rule_store_(rule_store),
  pipelined_client_(pipelined_client),
  aaa_client_(aaa_client),
  session_force_termination_timeout_ms_(session_force_termination_timeout_ms)
{
}

void LocalEnforcer::new_report()
{
  for (auto &session_pair : session_map_) {
    session_pair.second->new_report();
  }
}

void LocalEnforcer::finish_report()
{
  // Iterate through sessions and notify that report has finished. Terminate any
  // sessions that can be terminated.
  std::vector<std::string> imsi_to_terminate;
  for (auto &session_pair : session_map_) {
    session_pair.second->finish_report();
    if (session_pair.second->can_complete_termination()) {
      imsi_to_terminate.push_back(session_pair.first);
    }
  }
  for (std::string &imsi : imsi_to_terminate) {
    complete_termination(imsi, session_map_[imsi]->get_session_id());
  }
}

void LocalEnforcer::start()
{
  evb_->loopForever();
}

void LocalEnforcer::attachEventBase(folly::EventBase *evb)
{
  evb_ = evb;
}

void LocalEnforcer::stop()
{
  evb_->terminateLoopSoon();
}

folly::EventBase &LocalEnforcer::get_event_base()
{
  return *evb_;
}

bool LocalEnforcer::setup(
  const std::uint64_t &epoch,
  std::function<void(Status status, SetupFlowsResult)> callback)
{
  std::vector<SessionState::SessionInfo> session_infos;
  for(auto it = session_map_.begin(); it != session_map_.end(); it++)
  {
    SessionState::SessionInfo session_info;
    it->second->get_session_info(session_info);
    session_infos.push_back(session_info);
  }
  return pipelined_client_->setup(session_infos, epoch, callback);
}

void LocalEnforcer::aggregate_records(const RuleRecordTable &records)
{
  new_report(); // unmark all credits
  for (const RuleRecord &record : records.records()) {
    auto it = session_map_.find(record.sid());
    if (it == session_map_.end()) {
      MLOG(MERROR) << "Could not find session for IMSI " << record.sid()
                   << " during record aggregation";
      continue;
    }
    if (record.bytes_tx() > 0 || record.bytes_rx() > 0) {
      MLOG(MDEBUG) << "Subscriber " << record.sid() << " used "
                   << record.bytes_tx() << " tx bytes and " << record.bytes_rx()
                   << " rx bytes for rule " << record.rule_id();
    }
    it->second->add_used_credit(
      record.rule_id(), record.bytes_tx(), record.bytes_rx());
  }
  finish_report();
}

void LocalEnforcer::execute_actions(
  const std::vector<std::unique_ptr<ServiceAction>> &actions)
{
  for (auto &action_p : actions) {
    if (action_p->get_type() == TERMINATE_SERVICE) {
      terminate_service(
        action_p->get_imsi(),
        action_p->get_rule_ids(),
        action_p->get_rule_definitions());
    } else if (action_p->get_type() == ACTIVATE_SERVICE) {
      pipelined_client_->activate_flows_for_rules(
        action_p->get_imsi(),
        action_p->get_ip_addr(),
        action_p->get_rule_ids(),
        action_p->get_rule_definitions());
    } else if (action_p->get_type() == REDIRECT) {
      install_redirect_flow(action_p);
    } else if (action_p->get_type() == RESTRICT_ACCESS) {
      MLOG(MDEBUG) << "RESTRICT_ACCESS mode is unsupported"
                   << ", will just terminate the service.";
      terminate_service(
        action_p->get_imsi(),
        action_p->get_rule_ids(),
        action_p->get_rule_definitions());
    }
  }
}

void LocalEnforcer::terminate_service(
  const std::string &imsi,
  const std::vector<std::string> &rule_ids,
  const std::vector<PolicyRule> &dynamic_rules)
{
  pipelined_client_->deactivate_flows_for_rules(imsi, rule_ids, dynamic_rules);

  // tell AAA service to terminate radius session if necessary
  auto it = session_map_.find(imsi);
  if (it == session_map_.end()) {
    MLOG(MWARNING) << "Could not find session with IMSI " << imsi;
  } else if (it->second->is_radius_cwf_session()) {
    MLOG(MDEBUG) << "Asking AAA service to terminate session with "
                 << "Radius ID: " << it->second->get_radius_session_id()
                 << ", IMSI: " << imsi;
    aaa_client_->terminate_session(it->second->get_radius_session_id(), imsi);
  }
}

// TODO: make session_manager.proto and policydb.proto to use common field
static RedirectInformation_AddressType address_type_converter(
  RedirectServer_RedirectAddressType address_type)
{
  switch(address_type) {
    case RedirectServer_RedirectAddressType_IPV4:
      return RedirectInformation_AddressType_IPv4;
    case RedirectServer_RedirectAddressType_IPV6:
      return RedirectInformation_AddressType_IPv6;
    case RedirectServer_RedirectAddressType_URL:
      return RedirectInformation_AddressType_URL;
    case RedirectServer_RedirectAddressType_SIP_URI:
      return RedirectInformation_AddressType_SIP_URI;
  }
}

static PolicyRule create_redirect_rule(
  const std::unique_ptr<ServiceAction> &action)
{
  PolicyRule redirect_rule;
  redirect_rule.set_id("redirect");
  redirect_rule.set_priority(LocalEnforcer::REDIRECT_FLOW_PRIORITY);
  redirect_rule.set_rating_group(action->get_rating_group());

  RedirectInformation* redirect_info = redirect_rule.mutable_redirect();
  redirect_info->set_support(RedirectInformation_Support_ENABLED);

  auto redirect_server = action->get_redirect_server();
  redirect_info->set_address_type(
    address_type_converter(redirect_server.redirect_address_type()));
  redirect_info->set_server_address(
    redirect_server.redirect_server_address());

  return redirect_rule;
}

void LocalEnforcer::install_redirect_flow(
  const std::unique_ptr<ServiceAction> &action)
{
  std::vector<std::string> static_rules;
  std::vector<PolicyRule> dynamic_rules {create_redirect_rule(action)};

  pipelined_client_->activate_flows_for_rules(
    action->get_imsi(),
    action->get_ip_addr(),
    static_rules,
    dynamic_rules);
}

UpdateSessionRequest LocalEnforcer::collect_updates()
{
  UpdateSessionRequest request;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  for (auto &session_pair : session_map_) {
    session_pair.second->get_updates(&request, &actions);
  }
  execute_actions(actions);
  return request;
}

void LocalEnforcer::reset_updates(const UpdateSessionRequest &failed_request)
{
  for (const auto &update : failed_request.updates()) {
    auto it = session_map_.find(update.sid());
    if (it == session_map_.end()) {
      MLOG(MERROR) << "Could not reset credit for IMSI " << update.sid()
                   << " because it couldn't be found";
      return;
    }
    it->second->get_charging_pool().reset_reporting_credit(
      update.usage().charging_key());
  }
  for (const auto &update : failed_request.usage_monitors()) {
    auto it = session_map_.find(update.sid());
    if (it == session_map_.end()) {
      MLOG(MERROR) << "Could not reset credit for IMSI " << update.sid()
                   << " because it couldn't be found";
      return;
    }
    it->second->get_monitor_pool().reset_reporting_credit(
      update.update().monitoring_key());
  }
}

/*
 * If a rule needs to be tracked by the OCS, then it needs credit in order to
 * be activated. If it does not receive credit, it should not be installed.
 * If a rule has a monitoring key, it is not required that a usage monitor is
 * installed with quota
 */
static bool should_activate(
  const PolicyRule &rule,
  const std::unordered_set<uint32_t> &successful_credits)
{
  if (
    rule.tracking_type() == PolicyRule::ONLY_OCS ||
    rule.tracking_type() == PolicyRule::OCS_AND_PCRF) {
    return successful_credits.count(rule.rating_group()) > 0;
  }
  MLOG(MDEBUG) << "NO OCS TRACKING for this rule";
  ;
  // no tracking or PCRF-only tracking, activate
  return true;
}

void LocalEnforcer::schedule_static_rule_activation(
  const std::string &imsi,
  const std::string &ip_addr,
  const StaticRuleInstall &static_rule)
{
  std::vector<std::string> static_rules {static_rule.rule_id()};
  std::vector<PolicyRule> dynamic_rules;

  auto it = session_map_.find(imsi);
  auto delta = time_difference_from_now(static_rule.activation_time());
  MLOG(MDEBUG) << "Scheduling subscriber " << imsi << " static rule "
               << static_rule.rule_id() << " activation in "
               << (delta.count() / 1000) << " secs";
  evb_->runInEventBaseThread([=] {
    evb_->timer().scheduleTimeoutFn(
      std::move([=] {
        pipelined_client_->activate_flows_for_rules(
          imsi, ip_addr, static_rules, dynamic_rules);
        if (it == session_map_.end()) {
          MLOG(MWARNING) << "Could not find session for IMSI " << imsi
                         << "during installation of static rule "
                         << static_rule.rule_id();
        } else {
          it->second->activate_static_rule(static_rule.rule_id());
        }
      }),
      delta);
  });
}

void LocalEnforcer::schedule_dynamic_rule_activation(
  const std::string &imsi,
  const std::string &ip_addr,
  const DynamicRuleInstall &dynamic_rule)
{
  std::vector<std::string> static_rules;
  std::vector<PolicyRule> dynamic_rules {dynamic_rule.policy_rule()};

  auto it = session_map_.find(imsi);
  auto delta = time_difference_from_now(dynamic_rule.activation_time());
  MLOG(MDEBUG) << "Scheduling subscriber " << imsi << " dynamic rule "
               << dynamic_rule.policy_rule().id() << " activation in "
               << (delta.count() / 1000) << " secs";
  evb_->runInEventBaseThread([=] {
    evb_->timer().scheduleTimeoutFn(
      std::move([=] {
        pipelined_client_->activate_flows_for_rules(
          imsi, ip_addr, static_rules, dynamic_rules);
        if (it == session_map_.end()) {
          MLOG(MWARNING) << "Could not find session for IMSI " << imsi
                         << "during installation of dynamic rule "
                         << dynamic_rule.policy_rule().id();
        } else {
          it->second->insert_dynamic_rule(dynamic_rule.policy_rule());
        }
      }),
      delta);
  });
}

void LocalEnforcer::schedule_static_rule_deactivation(
  const std::string &imsi,
  const StaticRuleInstall &static_rule)
{
  std::vector<std::string> static_rules {static_rule.rule_id()};
  std::vector<PolicyRule> dynamic_rules;

  auto it = session_map_.find(imsi);
  auto delta = time_difference_from_now(static_rule.deactivation_time());
  MLOG(MDEBUG) << "Scheduling subscriber " << imsi << " static rule "
               << static_rule.rule_id() << " deactivation in "
               << (delta.count() / 1000) << " secs";
  evb_->runInEventBaseThread([=] {
    evb_->timer().scheduleTimeoutFn(
      std::move([=] {
        pipelined_client_->deactivate_flows_for_rules(
          imsi, static_rules, dynamic_rules);
        if (it == session_map_.end()) {
          MLOG(MWARNING) << "Could not find session for IMSI " << imsi
                         << "during removal of static rule "
                         << static_rule.rule_id();
        } else {
          if (!it->second->deactivate_static_rule(static_rule.rule_id()))
            MLOG(MWARNING) << "Could not find rule " << static_rule.rule_id()
                           << "for IMSI " << imsi
                           << " during static rule removal";
        }
      }),
      delta);
  });
}

void LocalEnforcer::schedule_dynamic_rule_deactivation(
  const std::string &imsi,
  const DynamicRuleInstall &dynamic_rule)
{
  std::vector<std::string> static_rules;
  std::vector<PolicyRule> dynamic_rules {dynamic_rule.policy_rule()};

  auto it = session_map_.find(imsi);
  auto delta = time_difference_from_now(dynamic_rule.deactivation_time());
  MLOG(MDEBUG) << "Scheduling subscriber " << imsi << " dynamic rule "
               << dynamic_rule.policy_rule().id() << " deactivation in "
               << (delta.count() / 1000) << " secs";
  evb_->runInEventBaseThread([=] {
    evb_->timer().scheduleTimeoutFn(
      std::move([=] {
        pipelined_client_->deactivate_flows_for_rules(
          imsi, static_rules, dynamic_rules);
        if (it == session_map_.end()) {
          MLOG(MWARNING) << "Could not find session for IMSI " << imsi
                         << "during removal of dynamic rule "
                         << dynamic_rule.policy_rule().id();
        } else {
          PolicyRule rule_dont_care;
          it->second->remove_dynamic_rule(
            dynamic_rule.policy_rule().id(), &rule_dont_care);
        }
      }),
      delta);
  });
}

void LocalEnforcer::process_create_session_response(
  const CreateSessionResponse &response,
  const std::unordered_set<uint32_t> &successful_credits,
  const std::string &imsi,
  const std::string &ip_addr,
  RulesToProcess *rules_to_activate,
  RulesToProcess *rules_to_deactivate)
{
  std::time_t current_time = time(NULL);
  for (const auto &static_rule : response.static_rules()) {
    auto id = static_rule.rule_id();
    PolicyRule rule;
    if (!rule_store_->get_rule(id, &rule)) {
      LOG(ERROR) << "Not activating rule " << id
                 << " because it could not be found";
      continue;
    }
    if (should_activate(rule, successful_credits)) {
      auto activation_time =
        TimeUtil::TimestampToSeconds(static_rule.activation_time());
      if (activation_time > current_time) {
        schedule_static_rule_activation(imsi, ip_addr, static_rule);
      } else {
        // activation time is an optional field in the proto message
        // it will be set as 0 by default
        // when it is 0 or some past time, the rule should be activated instanly
        rules_to_activate->static_rules.push_back(id);
        MLOG(MDEBUG) << "Activate Static rule id " << id;
      }

      auto deactivation_time =
        TimeUtil::TimestampToSeconds(static_rule.deactivation_time());
      if (deactivation_time > current_time) {
        schedule_static_rule_deactivation(imsi, static_rule);
      } else if (deactivation_time > 0) {
        // deactivation time is an optional field in the proto message
        // it will be set as 0 by default
        // when it is some past time, the rule should be deactivated instantly
        rules_to_deactivate->static_rules.push_back(id);
      }
    }
  }

  for (const auto &dynamic_rule : response.dynamic_rules()) {
    if (should_activate(dynamic_rule.policy_rule(), successful_credits)) {
      auto activation_time =
        TimeUtil::TimestampToSeconds(dynamic_rule.activation_time());
      if (activation_time > current_time) {
        schedule_dynamic_rule_activation(imsi, ip_addr, dynamic_rule);
      } else {
        rules_to_activate->dynamic_rules.push_back(dynamic_rule.policy_rule());
      }
      auto deactivation_time =
        TimeUtil::TimestampToSeconds(dynamic_rule.deactivation_time());
      if (deactivation_time > current_time) {
        schedule_dynamic_rule_deactivation(imsi, dynamic_rule);
      } else if (deactivation_time > 0) {
        rules_to_deactivate->dynamic_rules.push_back(
          dynamic_rule.policy_rule());
      }
    }
  }
}

// return true if any credit unit is valid and has non-zero volume
static bool contains_credit(const GrantedUnits &gsu)
{
  return (gsu.total().is_valid() && gsu.total().volume() > 0) ||
         (gsu.tx().is_valid() && gsu.tx().volume() > 0) ||
         (gsu.rx().is_valid() && gsu.rx().volume() > 0);
}

bool LocalEnforcer::init_session_credit(
  const std::string &imsi,
  const std::string &session_id,
  const SessionState::Config &cfg,
  const CreateSessionResponse &response)
{
  std::unordered_set<uint32_t> successful_credits;
  auto session_state = new SessionState(imsi, session_id, cfg, *rule_store_);
  for (const auto &credit : response.credits()) {
    session_state->get_charging_pool().receive_credit(credit);
    if (credit.success() && contains_credit(credit.credit().granted_units())) {
      successful_credits.insert(credit.charging_key());
    }
  }
  for (const auto &monitor : response.usage_monitors()) {
    session_state->get_monitor_pool().receive_credit(monitor);
  }
  session_map_[imsi] = std::unique_ptr<SessionState>(session_state);

  if (session_state->is_radius_cwf_session()) {
    MLOG(MDEBUG) << "Adding UE MAC flow for subscriber " << imsi;
    SubscriberID sid;
    sid.set_id(imsi);
    bool add_ue_mac_flow_success = pipelined_client_->add_ue_mac_flow(
      sid, session_state->get_mac_addr());
    if (!add_ue_mac_flow_success) {
      MLOG(MERROR) << "Failed to add UE MAC flow for subscriber " << imsi;
    }
  }

  auto ip_addr = session_state->get_subscriber_ip_addr();

  RulesToProcess rules_to_activate;
  RulesToProcess rules_to_deactivate;

  process_create_session_response(
    response,
    successful_credits,
    imsi,
    ip_addr,
    &rules_to_activate,
    &rules_to_deactivate);

  // activate_flows_for_rules() should be called even if there is no rule to
  // activate, because pipelined activates a "drop all packet" rule
  // when no rule is provided as the parameter
  for (const auto &static_rule : rules_to_activate.static_rules) {
    session_state->activate_static_rule(static_rule);
  }
  for (const auto &policy_rule : rules_to_activate.dynamic_rules) {
    session_state->insert_dynamic_rule(policy_rule);
  }
  bool activate_success = pipelined_client_->activate_flows_for_rules(
    imsi,
    ip_addr,
    rules_to_activate.static_rules,
    rules_to_activate.dynamic_rules);

  // deactivate_flows_for_rules() should not be called when there is no rule
  // to deactivate, because pipelined deactivates all rules
  // when no rule is provided as the parameter
  bool deactivate_success = true;
  if (
    rules_to_deactivate.static_rules.size() > 0 ||
    rules_to_deactivate.dynamic_rules.size() > 0) {
    for (const auto &static_rule : rules_to_deactivate.static_rules) {
      if (!session_state->deactivate_static_rule(static_rule))
        MLOG(MWARNING) << "Could not find rule " << static_rule  << "for IMSI "
                       << imsi << " during static rule removal";

    }
    for (const auto &policy_rule : rules_to_deactivate.dynamic_rules) {
      PolicyRule rule_dont_care;
      session_state->remove_dynamic_rule(policy_rule.id(), &rule_dont_care);
    }
    deactivate_success = pipelined_client_->deactivate_flows_for_rules(
      imsi,
      rules_to_deactivate.static_rules,
      rules_to_deactivate.dynamic_rules);
  }

  return activate_success && deactivate_success;
}

void LocalEnforcer::complete_termination(
  const std::string &imsi,
  const std::string &session_id)
{
  // If the session cannot be found in session_map_, or a new session has
  // already begun, do nothing.
  auto it = session_map_.find(imsi);
  if (it == session_map_.end() || it->second->get_session_id() != session_id) {
    // Session is already deleted, or new session already began, ignore.
    MLOG(MDEBUG) << "Could not find session for IMSI " << imsi
                 << " and session ID " << session_id
                 << ". Skipping termination.";
    return;
  }
  // Complete session termination and remove session from session_map_.
  it->second->complete_termination();
  session_map_.erase(imsi);
  MLOG(MDEBUG) << "Successfully terminated session for IMSI " << imsi
               << "session ID " << session_id;
}

void LocalEnforcer::update_session_credit(const UpdateSessionResponse &response)
{
  for (const auto &response : response.responses()) {
    auto it = session_map_.find(response.sid());
    if (it == session_map_.end()) {
      MLOG(MERROR) << "Could not find session for IMSI " << response.sid()
                   << " during update";
      return;
    }
    it->second->get_charging_pool().receive_credit(response);
  }
  for (const auto &usage_monitor_resp : response.usage_monitor_responses()) {
    if (revalidation_required(usage_monitor_resp.event_triggers())) {
      schedule_revalidation(usage_monitor_resp.revalidation_time());
    }
    auto it = session_map_.find(usage_monitor_resp.sid());
    if (it == session_map_.end()) {
      MLOG(MERROR) << "Could not find session for IMSI "
                   << usage_monitor_resp.sid() << " during update";
      return;
    }
    it->second->get_monitor_pool().receive_credit(usage_monitor_resp);
  }
}

void LocalEnforcer::terminate_subscriber(
  const std::string &imsi,
  std::function<void(SessionTerminateRequest)> on_termination_callback)
{
  auto it = session_map_.find(imsi);
  if (it == session_map_.end()) {
    MLOG(MERROR) << "Could not find session for IMSI " << imsi
                 << " during termination";
    throw SessionNotFound();
  }
  it->second->start_termination(on_termination_callback);

  if (!pipelined_client_->deactivate_all_flows(imsi)) {
    MLOG(MERROR) << "Could not deactivate flows for IMSI " << imsi
                 << " during termination";
  }
  std::string session_id = it->second->get_session_id();
  // The termination should be completed when aggregated usage record no longer
  // includes the imsi. If this has not occurred after the timeout, force
  // terminate the session.
  evb_->runAfterDelay(
    [this, imsi, session_id] {
      MLOG(MDEBUG) << "Completing forced termination for IMSI " << imsi;
      complete_termination(imsi, session_id);
    },
    session_force_termination_timeout_ms_);
}

uint64_t LocalEnforcer::get_charging_credit(
  const std::string &imsi,
  uint32_t charging_key,
  Bucket bucket) const
{
  auto it = session_map_.find(imsi);
  if (it == session_map_.end()) {
    return 0;
  }
  return it->second->get_charging_pool().get_credit(charging_key, bucket);
}

uint64_t LocalEnforcer::get_monitor_credit(
  const std::string &imsi,
  const std::string &mkey,
  Bucket bucket) const
{
  auto it = session_map_.find(imsi);
  if (it == session_map_.end()) {
    return 0;
  }
  return it->second->get_monitor_pool().get_credit(mkey, bucket);
}

ChargingReAuthAnswer::Result LocalEnforcer::init_charging_reauth(
  ChargingReAuthRequest request)
{
  auto it = session_map_.find(request.sid());
  if (it == session_map_.end()) {
    MLOG(MERROR) << "Could not find session for subscriber " << request.sid()
                 << " during reauth";
    return ChargingReAuthAnswer::SESSION_NOT_FOUND;
  }
  if (request.type() == ChargingReAuthRequest::SINGLE_SERVICE) {
    MLOG(MDEBUG) << "Initiating reauth of key " << request.charging_key()
                 << " for subscriber " << request.sid();
    return it->second->get_charging_pool().reauth_key(request.charging_key());
  }
  MLOG(MDEBUG) << "Initiating reauth of all keys for subscriber "
               << request.sid();
  return it->second->get_charging_pool().reauth_all();
}

void LocalEnforcer::init_policy_reauth(
  PolicyReAuthRequest request,
  PolicyReAuthAnswer &answer_out)
{
  auto it = session_map_.find(request.imsi());
  if (it == session_map_.end()) {
    MLOG(MERROR) << "Could not find session for subscriber " << request.imsi()
                 << " during policy reauth";
    answer_out.set_result(ReAuthResult::SESSION_NOT_FOUND);
    return;
  }

  receive_monitoring_credit_from_rar(request, it->second);

  RulesToProcess rules_to_activate;
  RulesToProcess rules_to_deactivate;

  get_rules_from_policy_reauth_request(
    request, it->second, &rules_to_activate, &rules_to_deactivate);

  auto ip_addr = it->second->get_subscriber_ip_addr();
  bool deactivate_success = true;
  if (
    rules_to_deactivate.static_rules.size() > 0 ||
    rules_to_deactivate.dynamic_rules.size() > 0) {
    deactivate_success = pipelined_client_->deactivate_flows_for_rules(
      request.imsi(),
      rules_to_deactivate.static_rules,
      rules_to_deactivate.dynamic_rules);
  }
  bool activate_success = true;
  if (
    rules_to_activate.static_rules.size() > 0 ||
    rules_to_activate.dynamic_rules.size() > 0) {
    activate_success = pipelined_client_->activate_flows_for_rules(
      request.imsi(),
      ip_addr,
      rules_to_activate.static_rules,
      rules_to_activate.dynamic_rules);
  }

  // Treat activate/deactivate as all-or-nothing when reporting rule failures
  answer_out.set_result(ReAuthResult::UPDATE_INITIATED);
  mark_rule_failures(activate_success, deactivate_success, request, answer_out);
}

void LocalEnforcer::receive_monitoring_credit_from_rar(
  const PolicyReAuthRequest &request,
  const std::unique_ptr<SessionState> &session)
{
  UsageMonitoringUpdateResponse monitoring_credit;
  monitoring_credit.set_session_id(request.session_id());
  monitoring_credit.set_sid("IMSI" + request.session_id());
  monitoring_credit.set_success(true);
  UsageMonitoringCredit* credit = monitoring_credit.mutable_credit();

  for (const auto &usage_monitoring_credit : request.usage_monitoring_credits()) {
    credit->CopyFrom(usage_monitoring_credit);
    session->get_monitor_pool().receive_credit(monitoring_credit);
  }
}

void LocalEnforcer::get_rules_from_policy_reauth_request(
  const PolicyReAuthRequest &request,
  const std::unique_ptr<SessionState> &session,
  RulesToProcess *rules_to_activate,
  RulesToProcess *rules_to_deactivate)
{
  MLOG(MDEBUG) << "Processing policy reauth for subscriber " << request.imsi();
  if (revalidation_required(request.event_triggers())) {
    schedule_revalidation(request.revalidation_time());
  }

  for (const auto &rule_id : request.rules_to_remove()) {
    // Try to remove as dynamic rule first
    PolicyRule dy_rule;
    bool is_dynamic = session->remove_dynamic_rule(rule_id, &dy_rule);

    if (is_dynamic) {
      rules_to_deactivate->dynamic_rules.push_back(dy_rule);
    } else {
      if (!session->deactivate_static_rule(rule_id))
        MLOG(MWARNING) << "Could not find rule " << rule_id << "for IMSI "
                       << request.imsi() << " during static rule removal";
      rules_to_deactivate->static_rules.push_back(rule_id);
    }
  }

  std::time_t current_time = time(NULL);
  std::string imsi = request.imsi();
  auto ip_addr = session->get_subscriber_ip_addr();
  for (const auto &static_rule : request.rules_to_install()) {
    auto activation_time =
      TimeUtil::TimestampToSeconds(static_rule.activation_time());
    if (activation_time > current_time) {
      schedule_static_rule_activation(imsi, ip_addr, static_rule);
    } else {
      session->activate_static_rule(static_rule.rule_id());
      rules_to_activate->static_rules.push_back(static_rule.rule_id());
    }

    auto deactivation_time =
      TimeUtil::TimestampToSeconds(static_rule.deactivation_time());
    if (deactivation_time > current_time) {
      schedule_static_rule_deactivation(imsi, static_rule);
    } else if (deactivation_time > 0) {
      if (!session->deactivate_static_rule(static_rule.rule_id()))
        MLOG(MWARNING) << "Could not find rule " << static_rule.rule_id()
                       << "for IMSI " << imsi << " during static rule removal";
      rules_to_deactivate->static_rules.push_back(static_rule.rule_id());
    }
  }

  for (const auto &dynamic_rule : request.dynamic_rules_to_install()) {
    auto activation_time =
      TimeUtil::TimestampToSeconds(dynamic_rule.activation_time());
    if (activation_time > current_time) {
      schedule_dynamic_rule_activation(imsi, ip_addr, dynamic_rule);
    } else {
      session->insert_dynamic_rule(dynamic_rule.policy_rule());
      rules_to_activate->dynamic_rules.push_back(dynamic_rule.policy_rule());
    }

    auto deactivation_time =
      TimeUtil::TimestampToSeconds(dynamic_rule.deactivation_time());
    if (deactivation_time > current_time) {
      schedule_dynamic_rule_deactivation(imsi, dynamic_rule);
    } else if (deactivation_time > 0) {
      PolicyRule rule_dont_care;
      session->remove_dynamic_rule(
        dynamic_rule.policy_rule().id(), &rule_dont_care);
      rules_to_deactivate->dynamic_rules.push_back(dynamic_rule.policy_rule());
    }
  }
}

bool LocalEnforcer::revalidation_required(
  const google::protobuf::RepeatedField<int> &event_triggers)
{
  auto it = std::find(
    event_triggers.begin(), event_triggers.end(), REVALIDATION_TIMEOUT);
  return it != event_triggers.end();
}

void LocalEnforcer::schedule_revalidation(
  const google::protobuf::Timestamp &revalidation_time)
{
  auto delta = time_difference_from_now(revalidation_time);
  evb_->runInEventBaseThread([=] {
    evb_->timer().scheduleTimeoutFn(
      std::move([=] {
        MLOG(MDEBUG) << "Revalidation timeout!";
        check_usage_for_reporting();
      }),
      delta);
  });
}

void LocalEnforcer::check_usage_for_reporting()
{
  auto request = collect_updates();
  if (request.updates_size() == 0 && request.usage_monitors_size() == 0) {
    return; // nothing to report
  }
  MLOG(MDEBUG) << "Sending " << request.updates_size()
               << " charging updates and " << request.usage_monitors_size()
               << " monitor updates to OCS and PCRF";

  // report to cloud
  (*reporter_).report_updates(
    request, [this, request](Status status, UpdateSessionResponse response) {
      if (!status.ok()) {
        reset_updates(request);
        MLOG(MERROR) << "Update of size " << request.updates_size()
                     << " to OCS failed entirely: " << status.error_message();
      } else {
        MLOG(MDEBUG) << "Received updated responses from OCS and PCRF";
        update_session_credit(response);
        // Check if we need to report more updates
        check_usage_for_reporting();
      }
    });
}

bool LocalEnforcer::is_imsi_duplicate(const std::string &imsi)
{
  auto it = session_map_.find(imsi);
  if (it == session_map_.end()) {
    return false;
  }
  return true;
}

bool LocalEnforcer::is_session_duplicate(
  const std::string &imsi, const magma::SessionState::Config &config)
{
  auto it = session_map_.find(imsi);
  if (it == session_map_.end()) {
    return false;
  }
  return it->second->is_same_config(config);
}

static void mark_rule_failures(
  const bool &activate_success,
  const bool &deactivate_success,
  const PolicyReAuthRequest &request,
  PolicyReAuthAnswer &answer_out)
{
  auto failed_rules = *answer_out.mutable_failed_rules();
  if (!deactivate_success) {
    for (const std::string &rule_id : request.rules_to_remove()) {
      failed_rules[rule_id] = PolicyReAuthAnswer::GW_PCEF_MALFUNCTION;
    }
  }
  if (!activate_success) {
    for (const StaticRuleInstall rule : request.rules_to_install()) {
      failed_rules[rule.rule_id()] = PolicyReAuthAnswer::GW_PCEF_MALFUNCTION;
    }
    for (const DynamicRuleInstall &d_rule :
         request.dynamic_rules_to_install()) {
      failed_rules[d_rule.policy_rule().id()] =
        PolicyReAuthAnswer::GW_PCEF_MALFUNCTION;
    }
  }
}
} // namespace magma
