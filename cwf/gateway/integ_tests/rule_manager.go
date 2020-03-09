/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package integ_tests

import (
	"fmt"

	fegProtos "magma/feg/cloud/go/protos"
	lteProtos "magma/lte/cloud/go/protos"
)

// RuleManager keeps track of rules and monitors for the integration
// test. It keeps track of all successfully installed static/dynamic rules
// along with usage monitors. The dynamic rules and usage monitors are added
// into the mock PCRF service. The static rules are usually streamed down into
// gateway policyDB from the cloud. Since this integration test does not cover
// the cloud component, we will directly insert the static policies into the
// redis database by using policyDBWrapper.
type RuleManager struct {
	// List of static rules successfully inserted into the policyDB store
	staticRuleDefs []*lteProtos.PolicyRule
	// List of dynamic rules successfully installed into PCRF
	accountRules []*fegProtos.AccountRules
	// List of usage monitors successfully installed into PCRF
	monitors []*fegProtos.UsageMonitorInfo
	// List of network wide static rules successfully inserted into the policyDB store
	omniPresentRules []*lteProtos.AssignedPolicies
	// Wrapper around redis operations for policyDB objects
	policyDBWrapper *policyDBWrapper
}

// NewRuleManager initialized the struct
func NewRuleManager() (*RuleManager, error) {
	policyDBWrapper, err := initializePolicyDBWrapper()
	if err != nil {
		return nil, err
	}
	return &RuleManager{
		policyDBWrapper: policyDBWrapper,
	}, nil
}

// AddStaticPassAllToDB adds a static rule that passes all traffic to policyDB
// storage
func (manager *RuleManager) AddStaticPassAllToDB(ruleID string, monitoringKey string, ratingGroup uint32, trackingType string, priority uint32) error {
	fmt.Printf("************************* Adding a Pass-All static rule: %s\n", ruleID)
	staticPassAll := getStaticPassAll(ruleID, monitoringKey, ratingGroup, trackingType, priority)
	return manager.insertStaticRuleIntoRedis(staticPassAll)
}

// AddStaticRuleToDB adds the static rule to policyDB storage
func (manager *RuleManager) AddStaticRuleToDB(rule *lteProtos.PolicyRule) error {
	fmt.Printf("************************* Adding a static rule: %s\n", rule.Id)
	return manager.insertStaticRuleIntoRedis(rule)
}

// AddDynamicPassAllToPCRF adds a dynamic rule that passes all traffic into PCRF
func (manager *RuleManager) AddDynamicPassAllToPCRF(imsi, ruleID, monitoringKey string) error {
	fmt.Printf("************************* Adding Pass-All Dynamic Rule for UE with IMSI: %s, ruleID: %s\n", imsi, ruleID)
	dynamicPassAll := getAccountRulesWithDynamicPassAll(imsi, ruleID, monitoringKey)
	return manager.addAccountRules(dynamicPassAll)
}

// AddRulesToPCRF adds the dynamic rule into PCRF
func (manager *RuleManager) AddRulesToPCRF(imsi string, ruleNames, baseNames []string) error {
	fmt.Printf("************************* Adding PCRF Rule for UE with IMSI: %s"+
		" with ruleNames=%v, baseNames=%v\n", imsi, ruleNames, baseNames)
	rules := makeAccountRules(imsi, ruleNames, baseNames)
	return manager.addAccountRules(rules)
}

// AddOmniPresentRulesToDB adds the network wide static rule to policyDB storage
func (manager *RuleManager) AddOmniPresentRulesToDB(keyId string, ruleNames, baseNames []string) error {
	fmt.Printf("************************* Adding a network wide rule\n")
	rule := makeAssignedRules(ruleNames, baseNames)
	return manager.insertOmniPresentRuleIntoRedis(keyId, rule)
}

// GetInstalledRulesByIMSI returns all dynamic rule ids and static rules
// referenced by dynamic rules keyed by the IMSI they are attached to.
func (manager *RuleManager) GetInstalledRulesByIMSI() map[string][]string {
	installedRulesByIMSI := map[string][]string{}
	for _, accountRules := range manager.accountRules {
		rules, exists := installedRulesByIMSI[accountRules.Imsi]
		if !exists {
			rules = []string{}
		}
		for _, ruleID := range accountRules.StaticRuleNames {
			rules = append(rules, ruleID)
		}
		for _, dynamicRule := range accountRules.DynamicRuleDefinitions {
			rules = append(rules, dynamicRule.RuleName)
		}
		installedRulesByIMSI[accountRules.Imsi] = rules
	}
	return installedRulesByIMSI
}

// RemoveInstalledRules removes previously installed rules from PCRF and policyDB
func (manager *RuleManager) RemoveInstalledRules() error {
	rulesIDsByIMSI := manager.GetInstalledRulesByIMSI()
	for imsi, ruleIDs := range rulesIDsByIMSI {
		err := deactivateSubscriberFlows(imsi, ruleIDs)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddUsageMonitor constructs a usage monitor according to the parameters and
// inserts it into PCRF
func (manager *RuleManager) AddUsageMonitor(imsi, monitoringKey string, volume, bytesPerGrant uint64) error {
	fmt.Printf("************************* Adding PCRF Usage Monitor for UE with IMSI: %s\n", imsi)
	usageMonitor := makeUsageMonitor(imsi, monitoringKey, volume, bytesPerGrant)
	manager.monitors = append(manager.monitors, usageMonitor)
	return addPCRFUsageMonitors(usageMonitor)
}

func (manager *RuleManager) insertStaticRuleIntoRedis(rule *lteProtos.PolicyRule) error {
	err := manager.policyDBWrapper.policyMap.Set(rule.Id, rule)
	if err != nil {
		manager.staticRuleDefs = append(manager.staticRuleDefs, rule)
	}
	return err
}

func (manager *RuleManager) insertOmniPresentRuleIntoRedis(keyID string, rule *lteProtos.AssignedPolicies) error {
	err := manager.policyDBWrapper.omniPresentRules.Set(keyID, rule)
	if err != nil {
		manager.omniPresentRules = append(manager.omniPresentRules, rule)
	}
	return err
}

func (manager *RuleManager) addAccountRules(rules *fegProtos.AccountRules) error {
	err := addPCRFRules(rules)
	if err != nil {
		manager.accountRules = append(manager.accountRules, rules)
	}
	return err
}

func getAccountRulesWithDynamicPassAll(imsi, ruleID, monitoringKey string) *fegProtos.AccountRules {
	return &fegProtos.AccountRules{
		Imsi:                imsi,
		StaticRuleNames:     []string{},
		StaticRuleBaseNames: []string{},
		DynamicRuleDefinitions: []*fegProtos.RuleDefinition{
			{
				RuleName:         ruleID,
				Precedence:       100,
				FlowDescriptions: []string{"permit out ip from any to any", "permit in ip from any to any"},
				MonitoringKey:    monitoringKey,
			},
		},
	}
}

func makeAccountRules(imsi string, ruleNames []string, baseNames []string) *fegProtos.AccountRules {
	return &fegProtos.AccountRules{
		Imsi:                imsi,
		StaticRuleNames:     ruleNames,
		StaticRuleBaseNames: baseNames,
	}
}

func makeAssignedRules(ruleNames []string, baseNames []string) *lteProtos.AssignedPolicies {
	return &lteProtos.AssignedPolicies{
		AssignedPolicies:  ruleNames,
		AssignedBaseNames: baseNames,
	}
}

func makeUsageMonitor(imsi, monitoringKey string, volume, bytesPerGrant uint64) *fegProtos.UsageMonitorInfo {
	return &fegProtos.UsageMonitorInfo{
		Imsi: imsi,
		UsageMonitorCredits: []*fegProtos.UsageMonitorCredit{
			{
				MonitoringKey:        monitoringKey,
				TotalQuota:           &fegProtos.Octets{TotalOctets: volume},
				QuotaGrantPerRequest: &fegProtos.Octets{TotalOctets: bytesPerGrant},
				MonitoringLevel:      fegProtos.UsageMonitorCredit_RuleLevel,
			},
		},
	}
}
