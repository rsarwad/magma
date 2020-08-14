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

#include <folly/io/async/EventBaseManager.h>
#include <glog/logging.h>
#include <gtest/gtest.h>

#include "LocalEnforcer.h"
#include "MagmaService.h"
#include "ProtobufCreators.h"
#include "RuleStore.h"
#include "ServiceRegistrySingleton.h"
#include "SessionState.h"
#include "SessionStore.h"
#include "SessiondMocks.h"
#include "StoredState.h"
#include "magma_logging.h"

using ::testing::Test;

namespace magma {

class SessionManagerHandlerTest : public ::testing::Test {
 protected:
  virtual void SetUp() {
    monitoring_key = "mk1";

    reporter               = std::make_shared<MockSessionReporter>();
    rule_store             = std::make_shared<StaticRuleStore>();
    session_store          = std::make_shared<SessionStore>(rule_store);
    pipelined_client       = std::make_shared<MockPipelinedClient>();
    auto directoryd_client = std::make_shared<MockDirectorydClient>();
    auto spgw_client       = std::make_shared<MockSpgwServiceClient>();
    auto aaa_client        = std::make_shared<MockAAAClient>();
    events_reporter        = std::make_shared<MockEventsReporter>();
    auto default_mconfig   = get_default_mconfig();
    local_enforcer         = std::make_shared<LocalEnforcer>(
        reporter, rule_store, *session_store, pipelined_client,
        directoryd_client, events_reporter, spgw_client, aaa_client, 0, 0,
        default_mconfig);
    evb = new folly::EventBase();
    std::thread([&]() {
      std::cout << "Started event loop thread\n";
      folly::EventBaseManager::get()->setEventBase(evb, 0);
    })
        .detach();

    local_enforcer->attachEventBase(evb);
    session_map_ = SessionMap{};

    session_manager = std::make_shared<LocalSessionManagerHandlerImpl>(
        local_enforcer, reporter.get(), directoryd_client,
        events_reporter, *session_store);
  }

  void insert_static_rule(
      std::shared_ptr<StaticRuleStore> rule_store, const std::string& m_key,
      uint32_t charging_key, const std::string& rule_id) {
    PolicyRule rule;
    rule.set_id(rule_id);
    rule.set_rating_group(charging_key);
    rule.set_monitoring_key(m_key);
    rule.set_tracking_type(PolicyRule::OCS_AND_PCRF);
    rule_store->insert_rule(rule);
  }

 protected:
  std::string monitoring_key;

  std::shared_ptr<SessionStore> session_store;
  std::shared_ptr<StaticRuleStore> rule_store;
  std::shared_ptr<LocalSessionManagerHandlerImpl> session_manager;
  std::shared_ptr<MockSessionReporter> reporter;
  std::shared_ptr<MockPipelinedClient> pipelined_client;
  std::shared_ptr<LocalEnforcer> local_enforcer;
  std::shared_ptr<MockEventsReporter> events_reporter;
  SessionIDGenerator id_gen_;
  folly::EventBase* evb;
  SessionMap session_map_;
};

MATCHER_P(CheckCreateSession, imsi, "") {
  auto req = static_cast<const CreateSessionRequest*>(arg);
  return req->common_context().sid().id() == imsi;
}

MATCHER_P(CheckUpdateSessionRequest, request_number, "") {
  auto request = static_cast<const UpdateSessionRequest&>(arg);
  for (const auto& credit_usage_update : request.updates()) {
    int req_number = credit_usage_update.request_number();
    return req_number == request_number;
  }
  return false;
}

TEST_F(SessionManagerHandlerTest, test_create_session_cfg) {
  // 1) Insert the entry for a rule
  insert_static_rule(rule_store, monitoring_key, 1, "rule1");
  std::vector<std::string> static_rules{"rule1"};

  LocalCreateSessionRequest request;
  CreateSessionResponse response;
  std::string hardware_addr_bytes = {0x0f, 0x10, 0x2e, 0x12, 0x3a, 0x55};
  std::string imsi                = "IMSI1";
  std::string msisdn              = "5100001234";
  std::string radius_session_id =
      "AA-AA-AA-AA-AA-AA:TESTAP__"
      "0F-10-2E-12-3A-55";
  std::string mac_addr = "0f:10:2e:12:3a:55";
  auto sid             = id_gen_.gen_session_id(imsi);
  SessionConfig cfg    = {
      .mac_addr          = mac_addr,
      .hardware_addr     = hardware_addr_bytes,
      .radius_session_id = radius_session_id};
  build_common_context(
      imsi, "", "apn1", msisdn, TGPP_WLAN, &cfg.common_context);
  build_wlan_context(
      mac_addr, radius_session_id,
      cfg.rat_specific_context.mutable_wlan_context());

  response.set_session_id(sid);
  // Only the active sessions are not recycled, to ensure that
  // this session is not automatically scheduled for termination
  // when RAT Type is WLAN, it needs monitoring keys...
  create_session_create_response(imsi, monitoring_key, static_rules, &response);
  response.mutable_static_rules()->Add()->mutable_rule_id()->assign("rule1");
  create_credit_update_response(
      imsi, 1, 1536, response.mutable_credits()->Add());

  SessionRead req  = {"IMSI1"};
  auto session_map = session_store->read_sessions(req);
  local_enforcer->init_session_credit(session_map, imsi, sid, cfg, response);
  bool write_success =
      session_store->create_sessions(imsi, std::move(session_map[imsi]));
  EXPECT_TRUE(write_success);
  session_map = session_store->read_sessions(req);
  auto it     = session_map.find("IMSI1");
  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map["IMSI1"].size(), 1);
  auto& session = session_map["IMSI1"][0];
  EXPECT_EQ(session->get_config().common_context.apn(), "apn1");

  grpc::ServerContext create_context;
  request.mutable_sid()->set_id("IMSI1");
  request.set_rat_type(RATType::TGPP_WLAN);
  request.set_hardware_addr(hardware_addr_bytes);
  request.set_msisdn(msisdn);
  request.set_radius_session_id(radius_session_id);
  request.set_apn("apn2");  // Update APN
  build_common_context(
      imsi, "", "apn2", msisdn, TGPP_WLAN, request.mutable_common_context());
  build_wlan_context(
      mac_addr, radius_session_id,
      request.mutable_rat_specific_context()->mutable_wlan_context());

  // Ensure session is not reported as its a duplicate
  EXPECT_CALL(*reporter, report_create_session(_, _)).Times(0);
  session_manager->CreateSession(
      &create_context, &request,
      [this](grpc::Status status, LocalCreateSessionResponse response_out) {});

  // Run session creation in the EventBase loop
  // It needs to loop once here.
  evb->loopOnce();

  // Assert the internal session config is updated to the new one
  session_map = session_store->read_sessions(req);
  it          = session_map.find("IMSI1");
  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map["IMSI1"].size(), 1);
  auto& session_apn2 = session_map["IMSI1"][0];
  EXPECT_EQ(session_apn2->get_config().common_context.apn(), "apn2");
}

TEST_F(SessionManagerHandlerTest, test_session_recycling_lte) {
  // 1) Insert the entry for a rule
  insert_static_rule(rule_store, monitoring_key, 1, "rule1");
  std::vector<std::string> static_rules{"rule1"};

  CreateSessionResponse response;
  std::string imsi   = "IMSI1";
  std::string msisdn = "5100001234";
  auto sid           = id_gen_.gen_session_id(imsi);
  SessionConfig cfg;
  build_common_context(imsi, "", "apn1", msisdn, TGPP_LTE, &cfg.common_context);
  build_lte_context(
      "spgw_ip", "imei", "plmn_id", "imsi_plmn_id", "user_loc", 1, nullptr,
      cfg.rat_specific_context.mutable_lte_context());

  response.set_session_id(sid);
  create_session_create_response(imsi, monitoring_key, static_rules, &response);
  response.mutable_static_rules()->Add()->mutable_rule_id()->assign("rule1");
  create_credit_update_response(
      imsi, 1, 1536, response.mutable_credits()->Add());

  SessionRead req  = {"IMSI1"};
  auto session_map = session_store->read_sessions(req);
  local_enforcer->init_session_credit(session_map, imsi, sid, cfg, response);
  bool write_success =
      session_store->create_sessions(imsi, std::move(session_map[imsi]));
  EXPECT_TRUE(write_success);
  session_map = session_store->read_sessions(req);
  auto it     = session_map.find("IMSI1");
  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map["IMSI1"].size(), 1);
  auto& session = session_map["IMSI1"][0];
  EXPECT_EQ(session->get_config().common_context.apn(), "apn1");

  // Only active, identical sessions can be recycled for LTE
  // The previously created session is active and this request has the same
  // context
  LocalCreateSessionRequest request;
  grpc::ServerContext create_context;
  request.mutable_sid()->set_id(imsi);
  request.set_rat_type(RATType::TGPP_LTE);
  request.set_msisdn(msisdn);
  request.set_apn("apn1");
  request.set_spgw_ipv4("spgw_ip");
  request.set_imei("imei");
  request.set_plmn_id("plmn_id");
  request.set_imsi_plmn_id("imsi_plmn_id");
  build_common_context(
      imsi, "", "apn1", msisdn, TGPP_LTE, request.mutable_common_context());
  build_lte_context(
      "spgw_ip", "imei", "plmn_id", "imsi_plmn_id", "user_loc", 1, nullptr,
      request.mutable_rat_specific_context()->mutable_lte_context());

  // Ensure session is not reported as its a duplicate
  EXPECT_CALL(*reporter, report_create_session(_, _)).Times(0);
  // Termination process for the previous session is started
  EXPECT_CALL(
      *pipelined_client,
      deactivate_flows_for_rules("IMSI1", testing::_, testing::_, testing::_))
      .Times(1)
      .WillOnce(testing::Return(true));
  session_manager->CreateSession(
      &create_context, &request,
      [this](grpc::Status status, LocalCreateSessionResponse response_out) {});

  // Run session creation in the EventBase loop
  // It needs to loop once here.
  evb->loopOnce();

  // Assert the internal session config is updated to the new one
  session_map = session_store->read_sessions(req);
  it          = session_map.find("IMSI1");
  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map["IMSI1"].size(), 1);
  auto& session_apn2 = session_map["IMSI1"][0];
  EXPECT_EQ(session_apn2->get_config().common_context.apn(), "apn1");

  // Now make the config not identical but with the same APN=apn1, this should
  // trigger a terminate for the existing and a creation for the new session
  LocalCreateSessionRequest request2;
  grpc::ServerContext create_context2;
  request2.mutable_sid()->set_id(imsi);
  request2.set_rat_type(RATType::TGPP_LTE);
  request2.set_msisdn(msisdn + "magma :)");  // different msisdn
  request2.set_apn("apn1");
  build_common_context(
      imsi, "", "apn1", msisdn + "magma :)", TGPP_LTE,
      request2.mutable_common_context());
  build_lte_context(
      "spgw_ip", "imei", "plmn_id", "imsi_plmn_id", "user_loc", 1, nullptr,
      request2.mutable_rat_specific_context()->mutable_lte_context());

  // Ensure a create session for the new session is sent, the old one is
  // terminated
  EXPECT_CALL(*reporter, report_create_session(_, _)).Times(1);

  session_manager->CreateSession(
      &create_context2, &request2,
      [this](grpc::Status status, LocalCreateSessionResponse response_out) {});

  // Run session creation in the EventBase loop
  // It needs to loop once here.
  evb->loopOnce();
}

TEST_F(SessionManagerHandlerTest, test_create_session) {
  // 1) Create the session
  LocalCreateSessionRequest request;
  std::string hardware_addr_bytes = {0x0f, 0x10, 0x2e, 0x12, 0x3a, 0x55};
  std::string imsi                = "IMSI1";
  std::string msisdn              = "5100001234";
  std::string radius_session_id =
      "AA-AA-AA-AA-AA-AA:TESTAP__"
      "0F-10-2E-12-3A-55";

  grpc::ServerContext server_context;
  request.mutable_sid()->set_id(imsi);
  request.set_rat_type(RATType::TGPP_LTE);
  request.set_hardware_addr(hardware_addr_bytes);
  request.set_msisdn(msisdn);
  request.set_radius_session_id(radius_session_id);

  CreateSessionResponse create_response;
  create_response.mutable_static_rules()->Add()->mutable_rule_id()->assign(
      "rule1");
  create_response.mutable_static_rules()->Add()->mutable_rule_id()->assign(
      "rule2");
  create_response.mutable_static_rules()->Add()->mutable_rule_id()->assign(
      "rule3");
  create_credit_update_response(
      "IMSI1", 1, 1536, create_response.mutable_credits()->Add());
  create_credit_update_response(
      "IMSI1", 2, 1024, create_response.mutable_credits()->Add());

  // Ensure session is reported as it is not a duplicate
  EXPECT_CALL(*reporter, report_create_session(_, _)).Times(1);
  session_manager->CreateSession(
      &server_context, &request,
      [this](grpc::Status status, LocalCreateSessionResponse response_out) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();
  evb->loopOnce();
  evb->loopOnce();
}

TEST_F(SessionManagerHandlerTest, test_report_rule_stats) {
  // 1) Insert the entry for a rule
  insert_static_rule(rule_store, monitoring_key, 1, "rule1");

  // 2) Create a session
  CreateSessionResponse response;
  response.mutable_static_rules()->Add()->mutable_rule_id()->assign("rule1");
  create_credit_update_response(
      "IMSI1", 1, 1025, response.mutable_credits()->Add());
  std::string hardware_addr_bytes = {0x0f, 0x10, 0x2e, 0x12, 0x3a, 0x55};
  std::string imsi                = "IMSI1";
  std::string msisdn              = "5100001234";
  std::string radius_session_id =
      "AA-AA-AA-AA-AA-AA:TESTAP__"
      "0F-10-2E-12-3A-55";
  auto sid          = id_gen_.gen_session_id(imsi);
  SessionConfig cfg = {
      .mac_addr          = "0f:10:2e:12:3a:55",
      .hardware_addr     = hardware_addr_bytes,
      .radius_session_id = radius_session_id};

  SessionRead req  = {"IMSI1"};
  auto session_map = session_store->read_sessions(req);
  EXPECT_CALL(*events_reporter, session_created(testing::_)).Times(1);
  local_enforcer->init_session_credit(session_map, imsi, sid, cfg, response);
  bool write_success =
      session_store->create_sessions(imsi, std::move(session_map[imsi]));
  EXPECT_TRUE(write_success);

  // Check the request number
  auto session_map_2 = session_store->read_sessions(SessionRead{imsi});
  EXPECT_EQ(session_map_2[imsi].front()->get_request_number(), 1);
  // 2) ReportRuleStats
  grpc::ServerContext server_context;
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record("IMSI1", "rule1", 512, 512, record_list->Add());

  EXPECT_CALL(
      *reporter, report_updates(CheckUpdateSessionRequest(1), testing::_))
      .Times(1);
  session_manager->ReportRuleStats(
      &server_context, &table,
      [this](grpc::Status status, orc8r::Void response_out) {});
  evb->loopOnce();

  session_map_2 = session_store->read_sessions(SessionRead{imsi});
  EXPECT_EQ(session_map_2[imsi].front()->get_request_number(), 2);
  evb->loopOnce();
}

TEST_F(SessionManagerHandlerTest, test_end_session) {
  // 1) Insert the entry for a rule
  insert_static_rule(rule_store, monitoring_key, 1, "rule1");

  // 2) Create a session
  CreateSessionResponse response;
  response.mutable_static_rules()->Add()->mutable_rule_id()->assign("rule1");
  create_credit_update_response(
      "IMSI1", 1, 1025, response.mutable_credits()->Add());
  std::string hardware_addr_bytes = {0x0f, 0x10, 0x2e, 0x12, 0x3a, 0x55};
  std::string imsi                = "IMSI1";
  std::string msisdn              = "5100001234";
  std::string radius_session_id =
      "AA-AA-AA-AA-AA-AA:TESTAP__"
      "0F-10-2E-12-3A-55";
  std::string apn   = "apn1";
  std::string mac_addr = "0f:10:2e:12:3a:55";
  auto sid          = id_gen_.gen_session_id(imsi);
  SessionConfig cfg = {
      .mac_addr          = mac_addr,
      .hardware_addr     = hardware_addr_bytes,
      .radius_session_id = radius_session_id};

  build_common_context(imsi, "", apn, msisdn, TGPP_WLAN, &cfg.common_context);
  build_wlan_context(
      mac_addr, radius_session_id,
      cfg.rat_specific_context.mutable_wlan_context());

  SessionRead req  = {"IMSI1"};
  auto session_map = session_store->read_sessions(req);
  local_enforcer->init_session_credit(session_map, imsi, sid, cfg, response);
  bool write_success =
      session_store->create_sessions(imsi, std::move(session_map[imsi]));
  EXPECT_TRUE(write_success);

  // 3) EndSession
  session_map = session_store->read_sessions(req);
  EXPECT_EQ(session_map["IMSI1"].size(), 1);
  LocalEndSessionRequest end_request;
  end_request.mutable_sid()->set_id("IMSI1");
  end_request.set_apn("apn1");
  grpc::ServerContext server_context;

  EXPECT_CALL(*reporter, report_terminate_session(_, _)).Times(1);
  session_manager->EndSession(
      &server_context, &end_request,
      [this](grpc::Status status, LocalEndSessionResponse response_out) {});
  evb->loopOnce();
  session_map = session_store->read_sessions(req);
  EXPECT_EQ(session_map["IMSI1"].size(), 1);

  evb->loopOnce();

  session_map = session_store->read_sessions(req);
  EXPECT_EQ(session_map["IMSI1"].size(), 0);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma
