/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <functional>

#include <grpc++/grpc++.h>
#include <lte/protos/session_manager.grpc.pb.h>

#include "LocalEnforcer.h"
#include "CloudReporter.h"
#include "SessionID.h"

using grpc::Server;
using grpc::ServerContext;
using grpc::Status;

namespace magma {
using namespace orc8r;

class LocalSessionManagerHandler {
 public:
  virtual ~LocalSessionManagerHandler() {}

  /**
   * Report flow stats from pipelined and track the usage per rule
   */
  virtual void ReportRuleStats(
    ServerContext *context,
    const RuleRecordTable *request,
    std::function<void(Status, Void)> response_callback) = 0;

  /**
   * Create a new session, initializing credit monitoring and requesting credit
   * from the cloud
   */
  virtual void CreateSession(
    ServerContext *context,
    const LocalCreateSessionRequest *request,
    std::function<void(Status, LocalCreateSessionResponse)>
      response_callback) = 0;

  /**
   * Terminate a session, untracking credit and terminating in the cloud
   */
  virtual void EndSession(
    ServerContext *context,
    const SubscriberID *request,
    std::function<void(Status, LocalEndSessionResponse)> response_callback) = 0;
};

/**
 * LocalSessionManagerHandler processes proxied gRPC requests to the session
 * manager. The handler uses a monitor and reporter to keep track of state
 * and report to the cloud, respectively
 */
class LocalSessionManagerHandlerImpl : public LocalSessionManagerHandler {
 public:
  LocalSessionManagerHandlerImpl(
    LocalEnforcer *monitor,
    SessionCloudReporter *reporter);

  ~LocalSessionManagerHandlerImpl() {}
  /**
   * Report flow stats from pipelined and track the usage per rule
   */
  void ReportRuleStats(
    ServerContext *context,
    const RuleRecordTable *request,
    std::function<void(Status, Void)> response_callback);

  /**
   * Create a new session, initializing credit monitoring and requesting credit
   * from the cloud
   */
  void CreateSession(
    ServerContext *context,
    const LocalCreateSessionRequest *request,
    std::function<void(Status, LocalCreateSessionResponse)> response_callback);

  /**
   * Terminate a session, untracking credit and terminating in the cloud
   */
  void EndSession(
    ServerContext *context,
    const SubscriberID *request,
    std::function<void(Status, LocalEndSessionResponse)> response_callback);

 private:
  LocalEnforcer *enforcer_;
  SessionCloudReporter *reporter_;
  SessionIDGenerator id_gen_;
  uint64_t current_epoch_;
  uint64_t reported_epoch_;
  std::chrono::seconds retry_timeout_;

 private:
  void check_usage_for_reporting();
  bool is_pipelined_restarted();
  bool restart_pipelined(const std::uint64_t &epoch);

  void send_create_session(
    const CreateSessionRequest &request,
    const std::string &imsi,
    const std::string &sid,
    const SessionState::Config &cfg,
    std::function<void(grpc::Status, LocalCreateSessionResponse)> response_callback);

  void handle_setup_callback(
    const std::uint64_t &epoch,
    Status status,
    SetupFlowsResult resp);
};

} // namespace magma
