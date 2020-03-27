/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <limits>

#include "DiameterCodes.h"
#include "SessionCredit.h"
#include "magma_logging.h"

namespace magma {

SessionCreditUpdateCriteria SessionCredit::UNUSED_UPDATE_CRITERIA{};
float SessionCredit::USAGE_REPORTING_THRESHOLD = 0.8;
uint64_t SessionCredit::EXTRA_QUOTA_MARGIN = 1024;
bool SessionCredit::TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED = true;

std::unique_ptr<SessionCredit> SessionCredit::unmarshal(
  const StoredSessionCredit &marshaled,
  CreditType credit_type)
{
  auto session_credit = std::make_unique<SessionCredit>(credit_type);

  session_credit->reporting_ = marshaled.reporting;
  session_credit->is_final_grant_ = marshaled.is_final;
  session_credit->unlimited_quota_ = marshaled.unlimited_quota;

  // FinalActionInfo
  FinalActionInfo final_action_info;
  final_action_info.final_action = marshaled.final_action_info.final_action;
  final_action_info.redirect_server = marshaled.final_action_info.redirect_server;
  session_credit->final_action_info_ = final_action_info;

  session_credit->reauth_state_ = marshaled.reauth_state;
  session_credit->service_state_ = marshaled.service_state;
  session_credit->expiry_time_ = marshaled.expiry_time;

  for ( int bucket_int = USED_TX; bucket_int != MAX_VALUES; bucket_int++ )
  {
    Bucket bucket = static_cast<Bucket>(bucket_int);
    if (marshaled.buckets.find(bucket) != marshaled.buckets.end()) {
      session_credit->buckets_[bucket] = marshaled.buckets.find(bucket)->second;
    }
  }

  session_credit->usage_reporting_limit_ = marshaled.usage_reporting_limit;

  return session_credit;
}

StoredSessionCredit SessionCredit::marshal()
{
  StoredSessionCredit marshaled {};
  marshaled.reporting = reporting_;
  marshaled.is_final = is_final_grant_;
  marshaled.unlimited_quota = unlimited_quota_;

  marshaled.final_action_info.final_action = final_action_info_.final_action;
  marshaled.final_action_info.redirect_server = final_action_info_.redirect_server;

  marshaled.reauth_state = reauth_state_;
  marshaled.service_state = service_state_;
  marshaled.expiry_time = expiry_time_;

  for ( int bucket_int = USED_TX; bucket_int != MAX_VALUES; bucket_int++ )
  {
    Bucket bucket = static_cast<Bucket>(bucket_int);
    marshaled.buckets[bucket] = buckets_[bucket];
  }

  marshaled.usage_reporting_limit = usage_reporting_limit_;

  return marshaled;
}

SessionCreditUpdateCriteria SessionCredit::get_update_criteria()
{
  SessionCreditUpdateCriteria uc {};
  uc.is_final = is_final_grant_;
  uc.reauth_state = reauth_state_;
  uc.service_state = service_state_;
  uc.expiry_time = expiry_time_;
  for ( int bucket_int = USED_TX; bucket_int != MAX_VALUES; bucket_int++ )
  {
    Bucket bucket = static_cast<Bucket>(bucket_int);
    uc.bucket_deltas[bucket] = 0;
  }
  return uc;
}

SessionCredit::SessionCredit(CreditType credit_type, ServiceState start_state):
  credit_type_(credit_type),
  reporting_(false),
  reauth_state_(REAUTH_NOT_NEEDED),
  service_state_(start_state),
  unlimited_quota_(false),
  buckets_ {}
{
}

SessionCredit::SessionCredit(
  CreditType credit_type,
  ServiceState start_state,
  bool unlimited_quota):
  credit_type_(credit_type),
  reporting_(false),
  reauth_state_(REAUTH_NOT_NEEDED),
  service_state_(start_state),
  unlimited_quota_(unlimited_quota),
  buckets_ {}
{
}

// by default, enable service
SessionCredit::SessionCredit(CreditType credit_type):
  SessionCredit(credit_type, SERVICE_ENABLED, false)
{
}

void SessionCredit::set_expiry_time(
  uint32_t validity_time,
  SessionCreditUpdateCriteria& uc)
{
  if (validity_time == 0) {
    // set as max possible time
    expiry_time_ = std::numeric_limits<std::time_t>::max();
    uc.expiry_time = expiry_time_;
    return;
  }
  expiry_time_ = std::time(nullptr) + validity_time;
  uc.expiry_time = expiry_time_;
}

void SessionCredit::add_used_credit(uint64_t used_tx, uint64_t used_rx, SessionCreditUpdateCriteria& uc)
{
  buckets_[USED_TX] += used_tx;
  buckets_[USED_RX] += used_rx;
  uc.bucket_deltas[USED_TX] += used_tx;
  uc.bucket_deltas[USED_RX] += used_rx;

  if (should_deactivate_service()) {
    MLOG(MDEBUG) << "Quota exhausted. Deactivating service";
    service_state_ = SERVICE_NEEDS_DEACTIVATION;
    uc.service_state = SERVICE_NEEDS_DEACTIVATION;
  }
}

void SessionCredit::reset_reporting_credit(SessionCreditUpdateCriteria& update_criteria)
{
  buckets_[REPORTING_RX] = 0;
  buckets_[REPORTING_TX] = 0;
  reporting_ = false;
  update_criteria.reporting = false;
}

void SessionCredit::mark_failure(
  uint32_t code,
  SessionCreditUpdateCriteria& update_criteria)
{
  if (DiameterCodeHandler::is_transient_failure(code)) {
    buckets_[REPORTED_RX] += buckets_[REPORTING_RX];
    buckets_[REPORTED_TX] += buckets_[REPORTING_TX];
    update_criteria.bucket_deltas[REPORTED_RX] += buckets_[REPORTING_RX];
    update_criteria.bucket_deltas[REPORTED_TX] += buckets_[REPORTING_TX];
  }
  reset_reporting_credit(update_criteria);
  if (should_deactivate_service()) {
    service_state_ = SERVICE_NEEDS_DEACTIVATION;
    update_criteria.service_state = SERVICE_NEEDS_DEACTIVATION;
  }
}

void SessionCredit::receive_credit(
  uint64_t total_volume,
  uint64_t tx_volume,
  uint64_t rx_volume,
  uint32_t validity_time,
  bool is_final_grant,
  FinalActionInfo final_action_info,
  SessionCreditUpdateCriteria& update_criteria)
{
  MLOG(MDEBUG) << "Received the following credit"
               << " total_volume=" << total_volume
               << " tx_volume=" << tx_volume
               << " rx_volume=" << rx_volume
               << " w/ validity time=" << validity_time;
  if (is_final_grant) {
    MLOG(MDEBUG) << "This credit received is the final grant, with final "
                 << "action=" << final_action_info.final_action;
  }

  buckets_[ALLOWED_TOTAL] += total_volume;
  buckets_[ALLOWED_TX] += tx_volume;
  buckets_[ALLOWED_RX] += rx_volume;
  update_criteria.bucket_deltas[ALLOWED_TOTAL] += total_volume;
  update_criteria.bucket_deltas[ALLOWED_TX] += tx_volume;
  update_criteria.bucket_deltas[ALLOWED_RX] += rx_volume;
  MLOG(MDEBUG) << "Total amount received since start of session is "
               << " total=" << buckets_[ALLOWED_TOTAL]
               << " tx=" << buckets_[ALLOWED_TX]
               << " rx=" << buckets_[ALLOWED_RX];
  // transfer reporting usage to reported
  buckets_[REPORTED_RX] += buckets_[REPORTING_RX];
  buckets_[REPORTED_TX] += buckets_[REPORTING_TX];
  MLOG(MDEBUG) << "Total amount reported since start of session is "
               << " total=" << buckets_[REPORTED_RX] + buckets_[REPORTED_TX]
               << " rx=" << buckets_[REPORTED_RX]
               << " tx=" << buckets_[REPORTED_TX];

  // Set the usage_reporting_limit so that we never report more than grant
  // we've received.
  update_criteria.bucket_deltas[REPORTED_RX] += buckets_[REPORTING_RX];
  update_criteria.bucket_deltas[REPORTED_TX] += buckets_[REPORTING_TX];
  auto reported_sum = buckets_[REPORTED_RX] + buckets_[REPORTED_TX];
  if (buckets_[ALLOWED_TOTAL] > reported_sum) {
    usage_reporting_limit_ = buckets_[ALLOWED_TOTAL] - reported_sum;
  } else if (usage_reporting_limit_ != 0) {
    MLOG(MINFO) << "We have reported data usage for all credit received, the "
                 << "upper limit for reporting is now 0.";
    usage_reporting_limit_ = 0;
    update_criteria.usage_reporting_limit = usage_reporting_limit_;
  }

  set_expiry_time(validity_time);
  reset_reporting_credit();

  is_final_grant_ = is_final_grant;
  final_action_info_ = final_action_info;

  if (reauth_state_ == REAUTH_PROCESSING) {
    reauth_state_ = REAUTH_NOT_NEEDED; // done
  }
  if (!is_quota_exhausted() && (service_state_ == SERVICE_DISABLED ||
                             service_state_ == SERVICE_NEEDS_DEACTIVATION)) {
    // if quota no longer exhausted, reenable services as needed
    MLOG(MDEBUG) << "Quota available. Activating service";
    service_state_ = SERVICE_NEEDS_ACTIVATION;
  }
}

bool SessionCredit::is_quota_exhausted(
  float usage_reporting_threshold,
  uint64_t extra_quota_margin)
{
  // used quota since last report
  uint64_t total_reported_usage = buckets_[REPORTED_TX] + buckets_[REPORTED_RX];
  uint64_t total_usage_since_report = std::max(
    uint64_t(0), buckets_[USED_TX] + buckets_[USED_RX] - total_reported_usage);
  uint64_t tx_usage_since_report =
    std::max(uint64_t(0), buckets_[USED_TX] - buckets_[REPORTED_TX]);
  uint64_t rx_usage_since_report =
    std::max(uint64_t(0), buckets_[USED_RX] - buckets_[REPORTED_RX]);

  // available quota since last report
  auto total_usage_reporting_threshold =
    extra_quota_margin + std::max(
                           0.0f,
                           (buckets_[ALLOWED_TOTAL] - total_reported_usage) *
                             usage_reporting_threshold);

  // reported tx/rx could be greater than allowed tx/rx
  // because some OCS/PCRF might not track tx/rx,
  // and 0 is added to the allowed credit when an credit update is received
  auto tx_usage_reporting_threshold =
    extra_quota_margin + std::max(
                           0.0f,
                           (buckets_[ALLOWED_TX] - buckets_[REPORTED_TX]) *
                             usage_reporting_threshold);
  auto rx_usage_reporting_threshold =
    extra_quota_margin + std::max(
                           0.0f,
                           (buckets_[ALLOWED_RX] - buckets_[REPORTED_RX]) *
                             usage_reporting_threshold);

   MLOG(MDEBUG) << " Is Quota exhausted?"
               << "\n Total used: " << buckets_[USED_TX] + buckets_[USED_RX]
               << "\n Allowed total: " << buckets_[ALLOWED_TOTAL]
               << "\n Reported total: " << total_reported_usage;

  bool is_exhausted = false;
  if (total_usage_since_report >= total_usage_reporting_threshold) {
    is_exhausted = true;
  } else if (
    (buckets_[ALLOWED_TX] > 0) &&
    (tx_usage_since_report >= tx_usage_reporting_threshold)) {
    is_exhausted = true;
  } else if (
    (buckets_[ALLOWED_RX] > 0) &&
    (rx_usage_since_report >= rx_usage_reporting_threshold)) {
    is_exhausted = true;
  }
  if (is_exhausted) {
    MLOG(MDEBUG) << " YES Quota exhausted ";
  }
  return is_exhausted;
}

bool SessionCredit::should_deactivate_service()
{
  if (credit_type_ != CreditType::CHARGING) {
    // we only terminate on charging quota exhaustion
    return false;
  }
  if (unlimited_quota_) {
    return false;
  }
  if (!SessionCredit::TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED) {
    // configured in sessiond.yml
    return false;
  }
  if (is_final_grant_ && is_quota_exhausted()) {
    // If we've exhausted the last grant, we should terminate
    return true;
  }
  if (is_quota_exhausted(1, SessionCredit::EXTRA_QUOTA_MARGIN)) {
    MLOG(MINFO) << "Terminating service because we have exhausted the "
                << "given quota AND the extra quota margin="
                << SessionCredit::EXTRA_QUOTA_MARGIN;
    // extra quota margin is configured in sessiond.yml
    // We will terminate if we've exceeded (given quota + extra quota margin).
    // If the gateway loses connection to the reporter, we should not allow the
    // UE to use internet for too long. This quota should be reasonably big so
    // that we don't terminate the session too easily.
    return true;
  }
  return false;
}

bool SessionCredit::validity_timer_expired()
{
  return time(NULL) >= expiry_time_;
}

CreditUpdateType SessionCredit::get_update_type()
{
  if (is_reporting()) {
    return CREDIT_NO_UPDATE;
  } else if (is_reauth_required()) {
    return CREDIT_REAUTH_REQUIRED;
  } else if (is_final_grant_ && is_quota_exhausted()) {
    // Don't request updates if there's no quota left
    return CREDIT_NO_UPDATE;
  } else if (is_quota_exhausted(SessionCredit::USAGE_REPORTING_THRESHOLD, 0)) {
    return CREDIT_QUOTA_EXHAUSTED;
  } else if (validity_timer_expired()) {
    return CREDIT_VALIDITY_TIMER_EXPIRED;
  } else {
    return CREDIT_NO_UPDATE;
  }
}

SessionCredit::Usage SessionCredit::get_usage_for_reporting(
  bool is_termination,
  SessionCreditUpdateCriteria& update_criteria)
{
  // Send delta. If bytes are reporting, don't resend them
  auto report = buckets_[REPORTED_TX] + buckets_[REPORTING_TX];
  uint64_t tx = buckets_[USED_TX] > report ? buckets_[USED_TX] - report : 0;
  report = buckets_[REPORTED_RX] - buckets_[REPORTING_RX];
  uint64_t rx = buckets_[USED_RX] > report ? buckets_[USED_RX] - report : 0;

  MLOG(MDEBUG) << "Data usage since last report is tx=" << tx
               << " rx=" << rx;
  if (!is_termination && !is_final_grant_) {
    // Apply reporting limits since the user is not getting terminated.
    // The limits are applied on total usage (ie. tx + rx)
    tx = std::min(tx, usage_reporting_limit_);
    rx = std::min(rx, usage_reporting_limit_ - tx);
    MLOG(MDEBUG) << "Since this is not the last report, we will only report "
                 << "min(usage, usage_reporting_limit="
                 << usage_reporting_limit_ << ")";
  }

  if (get_update_type() == CREDIT_REAUTH_REQUIRED) {
    reauth_state_ = REAUTH_PROCESSING;
    update_criteria.reauth_state = REAUTH_PROCESSING;
  }

  buckets_[REPORTING_TX] += tx;
  buckets_[REPORTING_RX] += rx;
  reporting_ = true;
  update_criteria.reporting = true;

  MLOG(MDEBUG) << "Amount reporting for this report:"
               << " tx=" << tx << " rx=" << rx;
  MLOG(MDEBUG) << "The total amount currently being reported:"
               << " tx=" << buckets_[REPORTING_TX]
               << " rx=" << buckets_[REPORTING_RX];

  return SessionCredit::Usage {.bytes_tx = tx, .bytes_rx = rx};
}

ServiceActionType SessionCredit::get_action(SessionCreditUpdateCriteria& update_criteria)
{
  if (service_state_ == SERVICE_NEEDS_DEACTIVATION) {
    MLOG(MDEBUG) << "Service State: " << service_state_;
    service_state_ = SERVICE_DISABLED;
    update_criteria.service_state = SERVICE_DISABLED;
    return get_action_for_deactivating_service();
  } else if (service_state_ == SERVICE_NEEDS_ACTIVATION) {
    MLOG(MDEBUG) << "Service State: " << service_state_;
    service_state_ = SERVICE_ENABLED;
    update_criteria.service_state = SERVICE_ENABLED;
    return ACTIVATE_SERVICE;
  }
  return CONTINUE_SERVICE;
}

ServiceActionType SessionCredit::get_action_for_deactivating_service()
{
  if (is_final_grant_ &&
    final_action_info_.final_action == ChargingCredit_FinalAction_REDIRECT) {
    return REDIRECT;
  } else if (is_final_grant_ &&
    final_action_info_.final_action == ChargingCredit_FinalAction_RESTRICT_ACCESS) {
    return RESTRICT_ACCESS;
  } else {
    return TERMINATE_SERVICE;
  }
}

bool SessionCredit::is_reporting()
{
  return reporting_;
}

uint64_t SessionCredit::get_credit(Bucket bucket) const
{
  return buckets_[bucket];
}

bool SessionCredit::is_reauth_required()
{
  return reauth_state_ == REAUTH_REQUIRED;
}

void SessionCredit::reauth(SessionCreditUpdateCriteria& update_criteria)
{
  reauth_state_ = REAUTH_REQUIRED;
  update_criteria.reauth_state = REAUTH_REQUIRED;
}

RedirectServer SessionCredit::get_redirect_server() {
  return final_action_info_.redirect_server;
}

void SessionCredit::set_is_final_grant(
  bool is_final_grant,
  SessionCreditUpdateCriteria& update_criteria) {
  is_final_grant_ = is_final_grant;
  update_criteria.is_final = is_final_grant;
}

void SessionCredit::set_reauth(
  ReAuthState reauth_state,
  SessionCreditUpdateCriteria& update_criteria) {
  reauth_state_ = reauth_state;
  update_criteria.reauth_state = reauth_state;
}

void SessionCredit::set_service_state(
  ServiceState service_state,
  SessionCreditUpdateCriteria& update_criteria) {
  service_state_ = service_state;
  update_criteria.service_state = service_state;
}

void SessionCredit::set_expiry_time(
  std::time_t expiry_time,
  SessionCreditUpdateCriteria& update_criteria) {
  expiry_time_ = expiry_time;
  update_criteria.expiry_time = expiry_time;
}

void SessionCredit::add_credit(
  uint64_t credit,
  Bucket bucket,
  SessionCreditUpdateCriteria& update_criteria) {
  buckets_[bucket] += credit;
  update_criteria.bucket_deltas[bucket] += credit;
}

} // namespace magma
