/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#pragma once

#include "SessionCredit.h"
#include "SessionRules.h"

namespace magma {

/**
 * CreditPool is an interface that defines a group of credits to track. It is
 * keyed by some type and requires some update response type to receive credit.
 */
template<
  typename KeyType,
  typename UpdateResponseType,
  typename UpdateRequestType>
class CreditPool {
 public:
  /**
   * add_used_credit adds usage to a specific credit
   */
  virtual bool
  add_used_credit(const KeyType &key, uint64_t used_tx, uint64_t used_rx) = 0;

  /**
   * reset_reporting_credit resets the credit state machine by clearing any
   * credit that was in the reporting state
   */
  virtual bool reset_reporting_credit(const KeyType &key) = 0;

  /**
   * get_updates gets any usage updates required by the credits in the pool
   */
  virtual void get_updates(
    std::string imsi,
    std::string ip_addr,
    SessionRules *session_rules,
    std::vector<UpdateRequestType> *updates_out,
    std::vector<std::unique_ptr<ServiceAction>> *actions_out) = 0;

  /**
   * get_termination_updates gets updates from all credits in the pool at the
   * time of termination
   */
  virtual bool get_termination_updates(
    SessionTerminateRequest *termination_out) = 0;

  /**
   * receive_credit adds allowed credit from the cloud
   */
  virtual bool receive_credit(const UpdateResponseType &update) = 0;

  /**
   * get_credit is a helper function to return the bytes in a credit bucket
   */
  virtual uint64_t get_credit(const KeyType &key, Bucket bucket) = 0;
};

/**
 * ChargingCreditPool manages a pool of credits for OCS-based charging. It is
 * keyed by rating groups (uint32) and receives CreditUpdateResponses to update
 * credit
 */
class ChargingCreditPool :
  public CreditPool<uint32_t, CreditUpdateResponse, CreditUsage> {
 public:
  ChargingCreditPool(const std::string &imsi);

  bool add_used_credit(const uint32_t &key, uint64_t used_tx, uint64_t used_rx)
    override;

  bool reset_reporting_credit(const uint32_t &key) override;

  void get_updates(
    std::string imsi,
    std::string ip_addr,
    SessionRules *session_rules,
    std::vector<CreditUsage> *updates_out,
    std::vector<std::unique_ptr<ServiceAction>> *actions_out) override;

  bool get_termination_updates(
    SessionTerminateRequest *termination_out) override;

  bool receive_credit(const CreditUpdateResponse &update) override;

  uint64_t get_credit(const uint32_t &key, Bucket bucket) override;

  ChargingReAuthAnswer::Result reauth_key(uint32_t charging_key);

  ChargingReAuthAnswer::Result reauth_all();

 private:
  std::unordered_map<uint32_t, std::unique_ptr<SessionCredit>> credit_map_;
  std::string imsi_;

 private:
  bool init_new_credit(const CreditUpdateResponse &update);
};

/**
 * UsageMonitoringCreditPool manages a pool of credits for PCRF-based usage
 * monitoring. It is keyed by monitoring keys (string) and receives
 * UsageMonitoringUpdateResponse to update credit
 */
class UsageMonitoringCreditPool :
  public CreditPool<
    std::string,
    UsageMonitoringUpdateResponse,
    UsageMonitorUpdate> {
 public:
  UsageMonitoringCreditPool(const std::string &imsi);

  bool add_used_credit(
    const std::string &key,
    uint64_t used_tx,
    uint64_t used_rx) override;

  bool reset_reporting_credit(const std::string &key) override;

  void get_updates(
    std::string imsi,
    std::string ip_addr,
    SessionRules *session_rules,
    std::vector<UsageMonitorUpdate> *updates_out,
    std::vector<std::unique_ptr<ServiceAction>> *actions_out) override;

  bool get_termination_updates(
    SessionTerminateRequest *termination_out) override;

  bool receive_credit(const UsageMonitoringUpdateResponse &update) override;

  uint64_t get_credit(const std::string &key, Bucket bucket) override;

  std::unique_ptr<std::string> get_session_level_key();

 private:
  struct Monitor {
    SessionCredit credit;
    MonitoringLevel level;
  };

  std::unordered_map<std::string, std::unique_ptr<Monitor>> monitor_map_;
  std::string imsi_;
  std::unique_ptr<std::string> session_level_key_;

 private:
  void update_session_level_key(const UsageMonitoringUpdateResponse &update);
  bool init_new_credit(const UsageMonitoringUpdateResponse &update);
};

} // namespace magma
