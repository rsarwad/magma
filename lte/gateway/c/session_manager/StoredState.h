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

#include <lte/protos/session_manager.grpc.pb.h>
#include <lte/protos/pipelined.grpc.pb.h>
#include <lte/protos/session_manager.grpc.pb.h>

#include "CreditKey.h"

namespace magma {
struct StoredQoSInfo {
  bool enabled;
  uint32_t qci;
};

struct StoredSessionConfig {
  std::string ue_ipv4;
  std::string spgw_ipv4;
  std::string msisdn;
  std::string apn;
  std::string imei;
  std::string plmn_id;
  std::string imsi_plmn_id;
  std::string user_location;
  RATType rat_type;
  std::string mac_addr; // MAC Address for WLAN
  std::string hardware_addr; // MAC Address for WLAN (binary)
  std::string radius_session_id;
  uint32_t bearer_id;
  StoredQoSInfo qos_info;
};

// Session Credit

struct StoredFinalActionInfo {
  ChargingCredit_FinalAction final_action;
  RedirectServer redirect_server;
};

/**
 * A bucket is a counter used for tracking credit volume across sessiond.
 * These are independently incremented and reset
 * Each value is in terms of a volume unit - either bytes or seconds
 */
enum Bucket {
  USED_TX = 0,
  USED_RX = 1,
  ALLOWED_TOTAL = 2,
  ALLOWED_TX = 3,
  ALLOWED_RX = 4,
  REPORTING_TX = 5,
  REPORTING_RX = 6,
  REPORTED_TX = 7,
  REPORTED_RX = 8,
  MAX_VALUES = 9,
};

enum ReAuthState {
  REAUTH_NOT_NEEDED = 0,
  REAUTH_REQUIRED = 1,
  REAUTH_PROCESSING = 2,
};

enum ServiceState {
  SERVICE_ENABLED = 0,
  SERVICE_NEEDS_DEACTIVATION = 1,
  SERVICE_DISABLED = 2,
  SERVICE_NEEDS_ACTIVATION = 3,
};

struct StoredSessionCredit {
  bool reporting;
  bool is_final;
  bool unlimited_quota;
  StoredFinalActionInfo final_action_info;
  ReAuthState reauth_state;
  ServiceState service_state;
  std::time_t  expiry_time;
  std::unordered_map<Bucket, uint64_t> buckets;
  uint64_t usage_reporting_limit;
};

struct StoredMonitor {
  StoredSessionCredit credit;
  MonitoringLevel level;
};

struct StoredChargingCreditPool {
  std::string imsi;
  std::unordered_map<
    CreditKey, StoredSessionCredit,
    decltype(&ccHash), decltype(&ccEqual)> credit_map;
};

struct StoredUsageMonitoringCreditPool {
  std::string imsi;
  std::string session_level_key; // "" maps to nullptr
  std::unordered_map<std::string, StoredMonitor> monitor_map;
};

struct StoredSessionState {
  StoredSessionConfig config;
  StoredChargingCreditPool charging_pool;
  StoredUsageMonitoringCreditPool monitor_pool;
  std::string imsi;
  std::string session_id;
  std::string core_session_id;
  magma::lte::SubscriberQuotaUpdate_Type subscriber_quota_state;
  magma::lte::TgppContext tgpp_context;
  std::vector<std::string> static_rule_ids;
  std::vector<PolicyRule> dynamic_rules;
  uint32_t request_number;
};

// Update Criteria

struct SessionCreditUpdateCriteria {
  bool is_final;
  bool reporting;
  ReAuthState reauth_state;
  ServiceState service_state;
  std::time_t  expiry_time;
  // Do not mark REPORTING buckets, but do mark REPORTED
  std::unordered_map<Bucket, uint64_t> bucket_deltas;
  uint64_t usage_reporting_limit;
};

struct SessionStateUpdateCriteria {
  std::vector<std::string> static_rules_to_install;
  std::vector<std::string> static_rules_to_uninstall;
  std::vector<PolicyRule> dynamic_rules_to_install;
  std::vector<std::string> dynamic_rules_to_uninstall;
  std::unordered_map<
    CreditKey, StoredSessionCredit,
    decltype(&ccHash), decltype(&ccEqual)> charging_credit_to_install;
  std::unordered_map<
    CreditKey, SessionCreditUpdateCriteria,
    decltype(&ccHash), decltype(&ccEqual)> charging_credit_map;
  std::unordered_map<std::string, StoredMonitor> monitor_credit_to_install;
  std::unordered_map<std::string, SessionCreditUpdateCriteria> monitor_credit_map;
  TgppContext updated_tgpp_context;
  magma::lte::SubscriberQuotaUpdate_Type updated_subscriber_quota_state;
};

SessionStateUpdateCriteria get_default_update_criteria();

} // namespace magma
