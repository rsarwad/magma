/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <mutex>
#include <unordered_map>

#include <lte/protos/policydb.pb.h>
#include <lte/protos/pipelined.grpc.pb.h>

#include "GRPCReceiver.h"

using grpc::Status;

namespace magma {
using namespace lte;
/**
 * Template class for keeping track of a map of one key to many policy rules
 */
template<typename KeyType>
class PoliciesByKeyMap {
 public:
  void insert(const KeyType &key, std::shared_ptr<PolicyRule> rule_p);

  void remove(const KeyType &key, std::shared_ptr<PolicyRule> rule_p);

  bool get_rule_ids_for_key(
    const KeyType &key,
    std::vector<std::string> &rules_out);

  bool get_rule_definitions_for_key(
    const KeyType &key,
    std::vector<PolicyRule> &rules_out);

 private:
  std::unordered_map<KeyType, std::vector<std::shared_ptr<PolicyRule>>>
    rules_by_key_;
};

/**
 * RuleChargingKeyMapper is a class for querying a bi-directional map of
 * rule_id <-> charging_key
 */
class PolicyRuleBiMap {
 public:
  /**
   * Clear the maps and add in the given rules
   */
  virtual void sync_rules(const std::vector<PolicyRule> &rules);

  virtual void insert_rule(const PolicyRule &rule);

  virtual bool get_rule(const std::string &rule_id, PolicyRule *rule);

  // Remove a rule from the store by ID. Returns true if the rule ID was found.
  // The removed rule will be copied into rule_out
  virtual bool remove_rule(const std::string &rule_id, PolicyRule *rule_out);

  /**
   * Get the charging key for a particular rule id. The charging key is set in
   * the out parameter charging_key
   * @returns false if it doesn't exist, true if so
   */
  virtual bool get_charging_key_for_rule_id(
    const std::string &rule_id,
    uint32_t *charging_key);

  virtual bool get_monitoring_key_for_rule_id(
    const std::string &rule_id,
    std::string *monitoring_key);

  /**
   * Get all the rules for a given key. Rule ids are copied into rules_out
   */
  virtual bool get_rule_ids_for_charging_key(
    uint32_t charging_key,
    std::vector<std::string> &rules_out);

  virtual bool get_rule_ids_for_monitoring_key(
    const std::string &monitoring_key,
    std::vector<std::string> &rules_out);

  /**
   * Get all the rules for a given key. Rule ids are copied into rules_out
   */
  virtual bool get_rule_definitions_for_charging_key(
    uint32_t charging_key,
    std::vector<PolicyRule> &rules_out);

  virtual bool get_rule_definitions_for_monitoring_key(
    const std::string &monitoring_key,
    std::vector<PolicyRule> &rules_out);

  virtual bool get_rule_ids(
    std::vector<std::string> &rules_ids_out
  );

  virtual bool get_rules(
    std::vector<PolicyRule> &rules_out
  );

 protected:
  std::mutex map_mutex_;
  // rule_id -> PolicyRule
  std::unordered_map<std::string, std::shared_ptr<PolicyRule>>
    rules_by_rule_id_;
  // charging key -> [PolicyRule]
  PoliciesByKeyMap<uint32_t> rules_by_charging_key_;
  // monitoring key -> [PolicyRule]
  PoliciesByKeyMap<std::string> rules_by_monitoring_key_;
};

/**
 * StaticRuleStore holds the rules that are defined in policydb
 */
class StaticRuleStore : public PolicyRuleBiMap {
};

/**
 * DynamicRuleStore manages dynamic rules for a subscriber
 */
class DynamicRuleStore : public PolicyRuleBiMap {
};

} // namespace magma
