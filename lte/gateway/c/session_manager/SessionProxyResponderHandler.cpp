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
#include <chrono>
#include <thread>

#include "EnumToString.h"
#include "SessionProxyResponderHandler.h"
#include "magma_logging.h"

#include "magma_logging.h"
#include "GrpcMagmaUtils.h"

using grpc::Status;

namespace magma {
SessionProxyResponderHandlerImpl::SessionProxyResponderHandlerImpl(
    std::shared_ptr<LocalEnforcer> enforcer, SessionStore& session_store)
    : enforcer_(enforcer), session_store_(session_store) {}

void SessionProxyResponderHandlerImpl::ChargingReAuth(
    ServerContext* context, const ChargingReAuthRequest* request,
    std::function<void(Status, ChargingReAuthAnswer)> response_callback) {
  auto& request_cpy = *request;
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request_cpy));
  MLOG(MDEBUG) << "Received a Gy (Charging) ReAuthRequest for "
               << request->session_id() << " and charging_key "
               << request->charging_key();
  enforcer_->get_event_base().runInEventBaseThread([this, request_cpy,
                                                    response_callback]() {
    auto session_map = get_sessions_for_charging(request_cpy);
    SessionUpdate update =
        SessionStore::get_default_session_update(session_map);
    auto result =
        enforcer_->init_charging_reauth(session_map, request_cpy, update);
    MLOG(MDEBUG) << "Result of Gy (Charging) ReAuthRequest "
                 << raa_result_to_str(result);
    ChargingReAuthAnswer ans;
    ans.set_result(result);

    bool update_success = session_store_.update_sessions(update);
    if (update_success) {
      MLOG(MDEBUG) << "Sending RAA response for Gy ReAuth "
                   << request_cpy.session_id();
      PrintGrpcMessage(static_cast<const google::protobuf::Message&>(ans));
      response_callback(Status::OK, ans);
    } else {
      // Todo If update fails, we should rollback changes from the request
      MLOG(MERROR) << "Failed to update Gy (Charging) ReAuthRequest changes...";
      auto status = Status(
          grpc::ABORTED,
          "ChargingReAuth no longer valid due to another update that "
          "updated the session first.");
      PrintGrpcMessage(static_cast<const google::protobuf::Message&>(ans));
      response_callback(status, ans);
    }
    MLOG(MDEBUG) << "Sent RAA response for Gy ReAuth "
                 << request_cpy.session_id();
  });
}

void SessionProxyResponderHandlerImpl::PolicyReAuth(
    ServerContext* context, const PolicyReAuthRequest* request,
    std::function<void(Status, PolicyReAuthAnswer)> response_callback) {
  auto& request_cpy = *request;
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request_cpy));
  MLOG(MDEBUG) << "Received a Gx (Policy) ReAuthRequest for session_id "
               << request->session_id();
  enforcer_->get_event_base().runInEventBaseThread([this, request_cpy,
                                                    response_callback]() {
    PolicyReAuthAnswer ans;
    auto session_map = get_sessions_for_policy(request_cpy);
    SessionUpdate update =
        SessionStore::get_default_session_update(session_map);
    enforcer_->init_policy_reauth(session_map, request_cpy, ans, update);
    MLOG(MDEBUG) << "Result of Gx (Policy) ReAuthRequest "
                 << raa_result_to_str(ans.result());
    bool update_success = session_store_.update_sessions(update);
    if (update_success) {
      MLOG(MDEBUG) << "Sending RAA response for Gx ReAuth "
                   << request_cpy.session_id();
      PrintGrpcMessage(static_cast<const google::protobuf::Message&>(ans));
      response_callback(Status::OK, ans);
    } else {
      // Todo If update fails, we should rollback changes from the request
      MLOG(MERROR) << "Failed to update Gx (Policy) ReAuthRequest changes...";
      auto status = Status(
          grpc::ABORTED,
          "PolicyReAuth no longer valid due to another update that "
          "updated the session first.");
      PrintGrpcMessage(static_cast<const google::protobuf::Message&>(ans));
      response_callback(status, ans);
    }
    MLOG(MDEBUG) << "Sent RAA response for Gx ReAuth "
                 << request_cpy.session_id();
  });
}

SessionMap SessionProxyResponderHandlerImpl::get_sessions_for_charging(
    const ChargingReAuthRequest& request) {
  SessionRead req = {request.sid()};
  return session_store_.read_sessions(req);
}

SessionMap SessionProxyResponderHandlerImpl::get_sessions_for_policy(
    const PolicyReAuthRequest& request) {
  SessionRead req = {request.imsi()};
  return session_store_.read_sessions(req);
}
}  // namespace magma
