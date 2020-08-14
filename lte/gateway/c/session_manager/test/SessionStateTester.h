/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#include <future>
#include <memory>
#include <utility>

#include <glog/logging.h>
#include <gtest/gtest.h>

#include "ProtobufCreators.h"
#include "SessionState.h"
#include "SessiondMocks.h"

using ::testing::Test;

namespace magma {

class SessionStateTest : public ::testing::Test {
 protected:
  virtual void SetUp() {
    SessionConfig test_sstate_cfg;
    auto tgpp_ctx = TgppContext();
    create_tgpp_context("gx.dest.com", "gy.dest.com", &tgpp_ctx);
    rule_store    = std::make_shared<StaticRuleStore>();
    session_state = std::make_shared<SessionState>(
        "imsi", "session", test_sstate_cfg, *rule_store, tgpp_ctx);
    update_criteria = get_default_update_criteria();
  }
  enum RuleType {
    STATIC  = 0,
    DYNAMIC = 1,
  };

  void insert_rule(
      uint32_t rating_group, const std::string& m_key,
      const std::string& rule_id, RuleType rule_type,
      std::time_t activation_time, std::time_t deactivation_time) {
    PolicyRule rule;
    create_policy_rule(rule_id, m_key, rating_group, &rule);
    RuleLifetime lifetime{
        .activation_time   = activation_time,
        .deactivation_time = deactivation_time,
    };
    switch (rule_type) {
      case STATIC:
        // insert into list of existing rules
        rule_store->insert_rule(rule);
        // mark the rule as active in session
        session_state->activate_static_rule(rule_id, lifetime, update_criteria);
        break;
      case DYNAMIC:
        session_state->insert_dynamic_rule(rule, lifetime, update_criteria);
        break;
    }
  }

  void schedule_rule(
      uint32_t rating_group, const std::string& m_key,
      const std::string& rule_id, RuleType rule_type,
      std::time_t activation_time, std::time_t deactivation_time) {
    PolicyRule rule;
    create_policy_rule(rule_id, m_key, rating_group, &rule);
    RuleLifetime lifetime{
        .activation_time   = activation_time,
        .deactivation_time = deactivation_time,
    };
    switch (rule_type) {
      case STATIC:
        // insert into list of existing rules
        rule_store->insert_rule(rule);
        // mark the rule as scheduled in the session
        session_state->schedule_static_rule(rule_id, lifetime, update_criteria);
        break;
      case DYNAMIC:
        session_state->schedule_dynamic_rule(rule, lifetime, update_criteria);
        break;
    }
  }

  // TODO: make session_manager.proto and policydb.proto to use common field
  static RedirectInformation_AddressType address_type_converter(
      RedirectServer_RedirectAddressType address_type) {
    switch (address_type) {
      case RedirectServer_RedirectAddressType_IPV4:
        return RedirectInformation_AddressType_IPv4;
      case RedirectServer_RedirectAddressType_IPV6:
        return RedirectInformation_AddressType_IPv6;
      case RedirectServer_RedirectAddressType_URL:
        return RedirectInformation_AddressType_URL;
      case RedirectServer_RedirectAddressType_SIP_URI:
        return RedirectInformation_AddressType_SIP_URI;
      default:
        return RedirectInformation_AddressType_IPv4;
    }
  }

  void insert_gy_redirection_rule(const std::string& rule_id) {
    PolicyRule redirect_rule;
    redirect_rule.set_id(rule_id);
    redirect_rule.set_priority(999);

    RedirectInformation* redirect_info = redirect_rule.mutable_redirect();
    redirect_info->set_support(RedirectInformation_Support_ENABLED);

    RedirectServer redirect_server;
    redirect_server.set_redirect_address_type(RedirectServer::URL);
    redirect_server.set_redirect_server_address("http://www.example.com/");

    redirect_info->set_address_type(
        address_type_converter(redirect_server.redirect_address_type()));
    redirect_info->set_server_address(
        redirect_server.redirect_server_address());

    RuleLifetime lifetime{};
    session_state->insert_gy_dynamic_rule(
        redirect_rule, lifetime, update_criteria);
  }

  void receive_credit_from_ocs(uint32_t rating_group, uint64_t volume) {
    CreditUpdateResponse charge_resp;
    create_credit_update_response("IMSI1", rating_group, volume, &charge_resp);
    session_state->receive_charging_credit(charge_resp, update_criteria);
  }

  void receive_credit_from_ocs(uint32_t rating_group, uint64_t total_volume,
                               uint64_t tx_volume,uint64_t rx_volume, bool is_final) {
    CreditUpdateResponse charge_resp;
    create_credit_update_response("IMSI1", rating_group,total_volume, tx_volume,
                                  rx_volume, is_final, &charge_resp);
    session_state->receive_charging_credit(charge_resp, update_criteria);
  }

  void receive_credit_from_pcrf(
      const std::string& mkey, uint64_t volume, MonitoringLevel level) {
    UsageMonitoringUpdateResponse monitor_resp;
    create_monitor_update_response("IMSI1", mkey, level, volume, &monitor_resp);
    session_state->receive_monitor(monitor_resp, update_criteria);
  }

  void activate_rule(
      uint32_t rating_group, const std::string& m_key,
      const std::string& rule_id, RuleType rule_type,
      std::time_t activation_time, std::time_t deactivation_time) {
    PolicyRule rule;
    create_policy_rule(rule_id, m_key, rating_group, &rule);
    RuleLifetime lifetime{
        .activation_time   = activation_time,
        .deactivation_time = deactivation_time,
    };
    switch (rule_type) {
      case STATIC:
        rule_store->insert_rule(rule);
        session_state->activate_static_rule(rule_id, lifetime, update_criteria);
        break;
      case DYNAMIC:
        session_state->insert_dynamic_rule(rule, lifetime, update_criteria);
        break;
    }
  }

 protected:
  std::shared_ptr<StaticRuleStore> rule_store;
  std::shared_ptr<SessionState> session_state;
  SessionStateUpdateCriteria update_criteria;
};
};  // namespace magma