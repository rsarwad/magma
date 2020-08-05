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

#include <memory>

#include <glog/logging.h>
#include <gtest/gtest.h>

#include "StoredState.h"
#include "magma_logging.h"

using ::testing::Test;

namespace magma {

class StoredStateTest : public ::testing::Test {
 protected:
  QoSInfo get_stored_qos_info() {
    QoSInfo stored;
    stored.enabled = true;
    stored.qci     = 123;
    return stored;
  }

  SessionConfig get_stored_session_config() {
    SessionConfig stored;
    stored.ue_ipv4           = "192.168.0.1";
    stored.spgw_ipv4         = "192.168.0.2";
    stored.msisdn            = "a";
    stored.apn               = "b";
    stored.imei              = "c";
    stored.plmn_id           = "d";
    stored.imsi_plmn_id      = "e";
    stored.user_location     = "f";
    stored.rat_type          = RATType::TGPP_WLAN;
    stored.mac_addr          = "g";  // MAC Address for WLAN
    stored.hardware_addr     = "h";  // MAC Address for WLAN (binary)
    stored.radius_session_id = "i";
    stored.bearer_id         = 321;
    stored.qos_info          = get_stored_qos_info();
    return stored;
  }

  FinalActionInfo get_stored_final_action_info() {
    FinalActionInfo stored;
    stored.final_action =
        ChargingCredit_FinalAction::ChargingCredit_FinalAction_REDIRECT;
    stored.redirect_server.set_redirect_address_type(
        RedirectServer_RedirectAddressType::
            RedirectServer_RedirectAddressType_IPV6);
    stored.redirect_server.set_redirect_server_address(
        "redirect_server_address");
    return stored;
  }

  StoredSessionCredit get_stored_session_credit() {
    StoredSessionCredit stored;

    stored.reporting         = true;
    stored.credit_limit_type = INFINITE_METERED;

    stored.buckets[USED_TX]       = 12345;
    stored.buckets[ALLOWED_TOTAL] = 54321;

    stored.grant_tracking_type = TX_ONLY;
    return stored;
  };

  StoredChargingGrant get_stored_charging_grant() {
    StoredChargingGrant stored;

    stored.is_final = true;

    stored.final_action_info.final_action =
        ChargingCredit_FinalAction::ChargingCredit_FinalAction_REDIRECT;
    stored.final_action_info.redirect_server.set_redirect_address_type(
        RedirectServer_RedirectAddressType::
            RedirectServer_RedirectAddressType_IPV6);
    stored.final_action_info.redirect_server.set_redirect_server_address(
        "redirect_server_address");

    stored.service_state = SERVICE_NEEDS_ACTIVATION;
    stored.reauth_state  = REAUTH_REQUIRED;

    stored.expiry_time = 32;
    stored.credit      = get_stored_session_credit();
    return stored;
  };

  StoredMonitor get_stored_monitor() {
    StoredMonitor stored;
    stored.credit = get_stored_session_credit();
    stored.level  = MonitoringLevel::PCC_RULE_LEVEL;
    return stored;
  }

  StoredChargingCreditMap get_stored_charging_credit_map() {
    StoredChargingCreditMap stored(4, &ccHash, &ccEqual);
    stored[CreditKey(1, 2)] = get_stored_charging_grant();
    return stored;
  }

  StoredMonitorMap get_stored_monitor_map() {
    StoredMonitorMap stored;
    stored["mk1"] = get_stored_monitor();
    return stored;
  }

  StoredSessionState get_stored_session() {
    StoredSessionState stored;

    stored.config                 = get_stored_session_config();
    stored.credit_map             = get_stored_charging_credit_map();
    stored.monitor_map            = get_stored_monitor_map();
    stored.session_level_key      = "session_level_key";
    stored.imsi                   = "IMSI1";
    stored.session_id             = "session_id";
    stored.core_session_id        = "core_session_id";
    stored.subscriber_quota_state = SubscriberQuotaUpdate_Type_VALID_QUOTA;
    stored.fsm_state              = SESSION_TERMINATING_FLOW_DELETED;

    magma::lte::TgppContext tgpp_context;
    tgpp_context.set_gx_dest_host("gx");
    tgpp_context.set_gy_dest_host("gy");
    stored.tgpp_context = tgpp_context;

    stored.pending_event_triggers[REVALIDATION_TIMEOUT] = READY;
    stored.revalidation_time.set_seconds(32);

    stored.request_number = 1;

    return stored;
  }
};

TEST_F(StoredStateTest, test_stored_qos_info) {
  auto stored = get_stored_qos_info();

  auto serialized   = serialize_stored_qos_info(stored);
  auto deserialized = deserialize_stored_qos_info(serialized);

  EXPECT_EQ(deserialized.enabled, true);
  EXPECT_EQ(deserialized.qci, 123);
}

TEST_F(StoredStateTest, test_stored_session_config) {
  auto stored = get_stored_session_config();

  std::string serialized     = serialize_stored_session_config(stored);
  SessionConfig deserialized = deserialize_stored_session_config(serialized);

  EXPECT_EQ(deserialized.ue_ipv4, "192.168.0.1");
  EXPECT_EQ(deserialized.spgw_ipv4, "192.168.0.2");
  EXPECT_EQ(deserialized.msisdn, "a");
  EXPECT_EQ(deserialized.apn, "b");
  EXPECT_EQ(deserialized.imei, "c");
  EXPECT_EQ(deserialized.plmn_id, "d");
  EXPECT_EQ(deserialized.imsi_plmn_id, "e");
  EXPECT_EQ(deserialized.user_location, "f");
  EXPECT_EQ(deserialized.rat_type, RATType::TGPP_WLAN);
  EXPECT_EQ(deserialized.mac_addr, "g");
  EXPECT_EQ(deserialized.hardware_addr, "h");
  EXPECT_EQ(deserialized.radius_session_id, "i");
  EXPECT_EQ(deserialized.bearer_id, 321);
  EXPECT_EQ(deserialized.qos_info.enabled, true);
  EXPECT_EQ(deserialized.qos_info.qci, 123);
}

TEST_F(StoredStateTest, test_stored_final_action_info) {
  auto stored = get_stored_final_action_info();

  auto serialized   = serialize_stored_final_action_info(stored);
  auto deserialized = deserialize_stored_final_action_info(serialized);

  EXPECT_EQ(
      deserialized.final_action,
      ChargingCredit_FinalAction::ChargingCredit_FinalAction_REDIRECT);
  EXPECT_EQ(
      deserialized.redirect_server.redirect_address_type(),
      RedirectServer_RedirectAddressType::
          RedirectServer_RedirectAddressType_IPV6);
  EXPECT_EQ(
      deserialized.redirect_server.redirect_server_address(),
      "redirect_server_address");
}

TEST_F(StoredStateTest, test_stored_session_credit) {
  auto stored = get_stored_session_credit();

  auto serialized   = serialize_stored_session_credit(stored);
  auto deserialized = deserialize_stored_session_credit(serialized);

  EXPECT_EQ(deserialized.reporting, true);
  EXPECT_EQ(deserialized.credit_limit_type, INFINITE_METERED);

  EXPECT_EQ(deserialized.buckets[USED_TX], 12345);
  EXPECT_EQ(deserialized.buckets[ALLOWED_TOTAL], 54321);

  EXPECT_EQ(deserialized.grant_tracking_type, TX_ONLY);
}

TEST_F(StoredStateTest, test_stored_monitor) {
  auto stored = get_stored_monitor();

  auto serialized   = serialize_stored_monitor(stored);
  auto deserialized = deserialize_stored_monitor(serialized);

  EXPECT_EQ(deserialized.credit.reporting, true);
  EXPECT_EQ(deserialized.credit.credit_limit_type, INFINITE_METERED);
  EXPECT_EQ(deserialized.credit.buckets[USED_TX], 12345);
  EXPECT_EQ(deserialized.credit.buckets[ALLOWED_TOTAL], 54321);
  EXPECT_EQ(deserialized.level, MonitoringLevel::PCC_RULE_LEVEL);
}

TEST_F(StoredStateTest, test_stored_charging_credit_map) {
  auto stored = get_stored_charging_credit_map();

  auto serialized   = serialize_stored_charging_credit_map(stored);
  auto deserialized = deserialize_stored_charging_credit_map(serialized);

  auto stored_charging_credit = deserialized[CreditKey(1, 2)];
  // test charging grant fields
  EXPECT_EQ(stored_charging_credit.is_final, true);
  EXPECT_EQ(
      stored_charging_credit.final_action_info.final_action,
      ChargingCredit_FinalAction::ChargingCredit_FinalAction_REDIRECT);
  EXPECT_EQ(
      stored_charging_credit.final_action_info.redirect_server
          .redirect_address_type(),
      RedirectServer_RedirectAddressType::
          RedirectServer_RedirectAddressType_IPV6);
  EXPECT_EQ(
      stored_charging_credit.final_action_info.redirect_server
          .redirect_server_address(),
      "redirect_server_address");
  EXPECT_EQ(stored_charging_credit.reauth_state, REAUTH_REQUIRED);
  EXPECT_EQ(stored_charging_credit.service_state, SERVICE_NEEDS_ACTIVATION);
  EXPECT_EQ(stored_charging_credit.expiry_time, 32);

  // test session credit fields
  auto credit = stored_charging_credit.credit;
  EXPECT_EQ(credit.reporting, true);
  EXPECT_EQ(credit.credit_limit_type, INFINITE_METERED);
  EXPECT_EQ(credit.buckets[USED_TX], 12345);
  EXPECT_EQ(credit.buckets[ALLOWED_TOTAL], 54321);
}

TEST_F(StoredStateTest, test_stored_monitor_map) {
  auto stored = get_stored_monitor_map();

  auto serialized   = serialize_stored_usage_monitor_map(stored);
  auto deserialized = deserialize_stored_usage_monitor_map(serialized);

  auto stored_monitor = deserialized["mk1"];
  EXPECT_EQ(stored_monitor.credit.reporting, true);
  EXPECT_EQ(stored_monitor.credit.credit_limit_type, INFINITE_METERED);
  EXPECT_EQ(stored_monitor.credit.buckets[USED_TX], 12345);
  EXPECT_EQ(stored_monitor.credit.buckets[ALLOWED_TOTAL], 54321);
  EXPECT_EQ(stored_monitor.level, MonitoringLevel::PCC_RULE_LEVEL);
}

TEST_F(StoredStateTest, test_stored_session) {
  auto stored = get_stored_session();

  auto serialized   = serialize_stored_session(stored);
  auto deserialized = deserialize_stored_session(serialized);

  auto config = deserialized.config;
  EXPECT_EQ(config.ue_ipv4, "192.168.0.1");
  EXPECT_EQ(config.spgw_ipv4, "192.168.0.2");
  EXPECT_EQ(config.msisdn, "a");
  EXPECT_EQ(config.apn, "b");
  EXPECT_EQ(config.imei, "c");
  EXPECT_EQ(config.plmn_id, "d");
  EXPECT_EQ(config.imsi_plmn_id, "e");
  EXPECT_EQ(config.user_location, "f");
  EXPECT_EQ(config.rat_type, RATType::TGPP_WLAN);
  EXPECT_EQ(config.mac_addr, "g");
  EXPECT_EQ(config.hardware_addr, "h");
  EXPECT_EQ(config.radius_session_id, "i");
  EXPECT_EQ(config.bearer_id, 321);
  EXPECT_EQ(config.qos_info.enabled, true);
  EXPECT_EQ(config.qos_info.qci, 123);

  auto stored_charging_credit = deserialized.credit_map[CreditKey(1, 2)];
  // test charging grant fields
  EXPECT_EQ(stored_charging_credit.is_final, true);
  EXPECT_EQ(
      stored_charging_credit.final_action_info.final_action,
      ChargingCredit_FinalAction::ChargingCredit_FinalAction_REDIRECT);
  EXPECT_EQ(
      stored_charging_credit.final_action_info.redirect_server
          .redirect_address_type(),
      RedirectServer_RedirectAddressType::
          RedirectServer_RedirectAddressType_IPV6);
  EXPECT_EQ(
      stored_charging_credit.final_action_info.redirect_server
          .redirect_server_address(),
      "redirect_server_address");
  EXPECT_EQ(stored_charging_credit.reauth_state, REAUTH_REQUIRED);
  EXPECT_EQ(stored_charging_credit.service_state, SERVICE_NEEDS_ACTIVATION);
  EXPECT_EQ(stored_charging_credit.expiry_time, 32);

  // test session credit fields
  auto credit = stored_charging_credit.credit;
  EXPECT_EQ(credit.reporting, true);
  EXPECT_EQ(credit.buckets[USED_TX], 12345);
  EXPECT_EQ(credit.buckets[ALLOWED_TOTAL], 54321);
  EXPECT_EQ(credit.credit_limit_type, INFINITE_METERED);

  EXPECT_EQ(deserialized.session_level_key, "session_level_key");

  auto stored_monitor = deserialized.monitor_map["mk1"];
  EXPECT_EQ(stored_monitor.credit.reporting, true);
  EXPECT_EQ(stored_monitor.credit.credit_limit_type, INFINITE_METERED);
  EXPECT_EQ(stored_monitor.credit.buckets[USED_TX], 12345);
  EXPECT_EQ(stored_monitor.credit.buckets[ALLOWED_TOTAL], 54321);
  EXPECT_EQ(stored_monitor.level, MonitoringLevel::PCC_RULE_LEVEL);

  EXPECT_EQ(stored.imsi, "IMSI1");
  EXPECT_EQ(stored.session_id, "session_id");
  EXPECT_EQ(stored.core_session_id, "core_session_id");
  EXPECT_EQ(
      stored.subscriber_quota_state, SubscriberQuotaUpdate_Type_VALID_QUOTA);
  EXPECT_EQ(stored.fsm_state, SESSION_TERMINATING_FLOW_DELETED);

  EXPECT_EQ(stored.tgpp_context.gx_dest_host(), "gx");
  EXPECT_EQ(stored.tgpp_context.gy_dest_host(), "gy");

  EXPECT_EQ(stored.pending_event_triggers.size(), 1);
  EXPECT_EQ(stored.pending_event_triggers[REVALIDATION_TIMEOUT], READY);
  EXPECT_EQ(stored.revalidation_time.seconds(), 32);

  EXPECT_EQ(stored.request_number, 1);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma
