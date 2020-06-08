/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers_test

import (
	"fmt"
	"testing"
	"time"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/multiplex"
	"magma/feg/gateway/policydb"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/feg/gateway/services/session_proxy/credit_control/gx"
	"magma/feg/gateway/services/session_proxy/credit_control/gy"
	"magma/feg/gateway/services/session_proxy/servicers"
	"magma/gateway/mconfig"
	"magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thoas/go-funk"
	"golang.org/x/net/context"
)

const (
	IMSI1          = "IMSI00101"
	IMSI2          = "IMSI00102"
	IMSI1_NOPREFIX = "00101"
	IMSI2_NOPREFIX = "00102"
	IMSI1_uint64   = uint64(101)
	IMSI2_uint64   = uint64(102)
	NUMBER_SERVERS = 5 // Must be always bigger or equal to num of imsis
)

var (
	imsis          = []string{IMSI1, IMSI2}
	imsis_noprefix = []string{IMSI1_NOPREFIX, IMSI2_NOPREFIX}
	ismis_uint64   = []uint64{IMSI1_uint64, IMSI2_uint64}
	// as many ports as servers
	ocs_server_ports  = []string{"3869", "3870", "3871", "3872", "3873"}
	pcrf_server_ports = []string{"3879", "3880", "3881", "3882", "3883"}
)

// ---- MockPolicyClient ----
type MockPolicyClient struct {
	mock.Mock
}

func (p *MockPolicyClient) SendCreditControlRequest(
	server *diameter.DiameterServerConfig,
	done chan interface{},
	request *gx.CreditControlRequest,
) error {
	args := p.Called(server, done, request)
	return args.Error(0)
}

func (p *MockPolicyClient) IgnoreAnswer(request *gx.CreditControlRequest) {
	return
}

func (p *MockPolicyClient) EnableConnections() error {
	p.Called()
	return nil
}

func (p *MockPolicyClient) DisableConnections(period time.Duration) {
	p.Called(period)
	return
}

// ---- MockCreditClient ----
type MockCreditClient struct {
	mock.Mock
}

func (cc *MockCreditClient) SendCreditControlRequest(
	server *diameter.DiameterServerConfig,
	done chan interface{},
	request *gy.CreditControlRequest,
) error {
	args := cc.Called(server, done, request)
	return args.Error(0)
}

func (cc *MockCreditClient) IgnoreAnswer(request *gy.CreditControlRequest) {
	return
}

func (cc *MockCreditClient) EnableConnections() error {
	cc.Called()
	return nil
}

func (cc *MockCreditClient) DisableConnections(period time.Duration) {
	cc.Called(period)
	return
}

// ---- MockPolicyDBClient ----
type MockPolicyDBClient struct {
	mock.Mock
}

func (client *MockPolicyDBClient) GetChargingKeysForRules(ruleIDs []string, ruleDefs []*protos.PolicyRule) []policydb.ChargingKey {

	args := client.Called(ruleIDs)
	return args.Get(0).([]policydb.ChargingKey)
}

func (client *MockPolicyDBClient) GetRuleIDsForBaseNames(baseNames []string) []string {
	args := client.Called(baseNames)
	return args.Get(0).([]string)
}

func (client *MockPolicyDBClient) GetPolicyRuleByID(id string) (*protos.PolicyRule, error) {
	return nil, nil
}

func (client *MockPolicyDBClient) GetOmnipresentRules() ([]string, []string) {
	args := client.Called()
	return args.Get(0).([]string), args.Get(1).([]string)
}

//  ---- MockMultiplexor ----
type MockMultiplexor struct {
	mock.Mock
	imsiToIndex map[uint64]int
}

func (mp *MockMultiplexor) GetIndex(muxCtx *multiplex.Context) (int, error) {
	imsi, err := muxCtx.GetIMSI()
	if err != nil {
		return -1, err
	}
	return mp.imsiToIndex[imsi], nil
}

// ---- TESTS ----

func TestSessionControllerInit(t *testing.T) {
	mockConfig := getTestConfig()
	mockControlParams := getMockControllerParams(mockConfig)
	mockPolicyDb := &MockPolicyDBClient{}
	mockMux := getMockMultiplexor(NUMBER_SERVERS)
	srv := servicers.NewCentralSessionControllers(mockControlParams, mockPolicyDb, mockMux)
	standardUsageTest(t, srv, mockControlParams, mockPolicyDb, mockMux, gy.PerKeyInit)
}

func TestStartSessionGxFail(t *testing.T) {
	// Set up mocks
	mockConfig := getTestConfig()
	mockControlParams := getMockControllerParams(mockConfig)
	mockPolicyDb := &MockPolicyDBClient{}
	mockMux := getMockMultiplexor(NUMBER_SERVERS)

	// Get the controller for this imsi
	idx, err := mockMux.GetIndex(multiplex.NewContext().WithIMSI(IMSI1))
	assert.NoError(t, err)

	mocksGx := mockControlParams[idx].PolicyClient.(*MockPolicyClient)
	// Send back DIAMETER_RATING_FAILED (5031) from gx
	mocksGx.On("SendCreditControlRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)
		done <- &gx.CreditControlAnswer{
			ResultCode:    uint32(diameter.DiameterRatingFailed),
			SessionID:     request.SessionID,
			RequestNumber: request.RequestNumber,
		}
	}).Once()
	// If gx fails gy should not be used at all

	srv := servicers.NewCentralSessionControllers(mockControlParams, mockPolicyDb, mockMux)
	ctx := context.Background()
	_, err = srv.CreateSession(ctx, &protos.CreateSessionRequest{
		Subscriber: &protos.SubscriberID{
			Id: IMSI1,
		},
		SessionId: genSessionID(IMSI1),
	})
	mocksGx.AssertExpectations(t)
	assert.Error(t, err)
}

func TestStartSessionGyFail(t *testing.T) {
	// Set up mocks
	mockConfig := getTestConfig()
	mockControlParams := getMockControllerParams(mockConfig)
	mockPolicyDb := &MockPolicyDBClient{}
	mockMux := getMockMultiplexor(NUMBER_SERVERS)

	// Get the controller for this imsi
	idx, err := mockMux.GetIndex(multiplex.NewContext().WithIMSI(IMSI1))
	assert.NoError(t, err)
	mocksGx := mockControlParams[idx].PolicyClient.(*MockPolicyClient)
	mocksGy := mockControlParams[idx].CreditClient.(*MockCreditClient)

	// Send back DIAMETER_SUCCESS (2001) from gx
	mocksGx.On("SendCreditControlRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)

		activationTime := time.Unix(1, 0)
		deactivationTime := time.Unix(2, 0)
		ruleInstalls := []*gx.RuleInstallAVP{
			&gx.RuleInstallAVP{
				RuleNames:            []string{"static_rule_1"},
				RuleActivationTime:   &activationTime,
				RuleDeactivationTime: &deactivationTime,
			},
		}

		done <- &gx.CreditControlAnswer{
			ResultCode:     uint32(diameter.SuccessCode),
			SessionID:      request.SessionID,
			RequestNumber:  request.RequestNumber,
			RuleInstallAVP: ruleInstalls,
		}
	}).Once()

	mockPolicyDb.On("GetChargingKeysForRules", mock.Anything, mock.Anything).Return(
		[]policydb.ChargingKey{{RatingGroup: 1}}, nil).Once()
	// no omnipresent rules
	mockPolicyDb.On("GetOmnipresentRules").Return([]string{}, []string{}).Once()
	mockPolicyDb.On("GetRuleIDsForBaseNames", mock.Anything).Return([]string{}).Once()

	// Send back DIAMETER_RATING_FAILED (5031) from gy
	mocksGy.On("SendCreditControlRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gy.CreditControlRequest)
		done <- &gy.CreditControlAnswer{
			ResultCode:    uint32(diameter.DiameterRatingFailed),
			SessionID:     request.SessionID,
			RequestNumber: request.RequestNumber,
		}
	}).Once()

	srv := servicers.NewCentralSessionControllers(mockControlParams, mockPolicyDb, mockMux)
	ctx := context.Background()
	_, err = srv.CreateSession(ctx, &protos.CreateSessionRequest{
		Subscriber: &protos.SubscriberID{
			Id: IMSI1,
		},
		SessionId: genSessionID(IMSI1),
	})
	mocksGx.AssertExpectations(t)
	assert.Error(t, err)
}

func standardUsageTest(
	t *testing.T,
	srv servicers.CentralSessionControllerServerWithHealth,
	controllerParams []*servicers.ControllerParam,
	policyDb policydb.PolicyDBClient,
	mux multiplex.Multiplexor,
	initMethod gy.InitMethod,
) error {
	ctx := context.Background()
	mockPolicyDb := policyDb.(*MockPolicyDBClient)

	// Create a structure to store the pointers to the type assertions. his is needed later to
	// be used on Enable/Disable. If it were not saved here the reference of the type to be
	// asserted will be different than the reference of the type inside the srv
	mocksGxs := make([]*MockPolicyClient, 0, len(controllerParams))
	mocksGys := make([]*MockCreditClient, 0, len(controllerParams))
	for _, cp := range controllerParams {
		mocksGxs = append(mocksGxs, cp.PolicyClient.(*MockPolicyClient))
		mocksGys = append(mocksGys, cp.CreditClient.(*MockCreditClient))
	}

	// Get the controller for this imsi
	idx, err := mux.GetIndex(multiplex.NewContext().WithIMSI(IMSI1))

	assert.NoError(t, err)

	mocksGx := mocksGxs[idx]
	mocksGy := mocksGys[idx]

	maxReqBWUL := uint32(128000)
	maxReqBWDL := uint32(128000)
	key1 := []byte("key1")

	// send static rules back
	mocksGx.On("SendCreditControlRequest", mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTInit)),
	).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)
		qos := gx.QosInformation{MaxReqBwUL: &maxReqBWUL, MaxReqBwDL: &maxReqBWDL}
		redirect := gx.RedirectInformation{
			RedirectSupport:       1,
			RedirectAddressType:   2,
			RedirectServerAddress: "http://www.example.com/",
		}

		var (
			rg20 uint32 = 20
			si20 uint32 = 201
			rg21 uint32 = 21
		)
		activationTime := time.Unix(1, 0)
		deactivationTime := time.Unix(2, 0)
		ruleInstalls := []*gx.RuleInstallAVP{
			&gx.RuleInstallAVP{
				RuleNames:     []string{"static_rule_1", "static_rule_2"},
				RuleBaseNames: []string{"base_10"},
				RuleDefinitions: []*gx.RuleDefinition{
					&gx.RuleDefinition{
						RuleName:            "dyn_rule_20",
						RatingGroup:         &rg20,
						ServiceIdentifier:   &si20,
						Precedence:          100,
						MonitoringKey:       key1,
						RedirectInformation: &redirect,
						Qos:                 &qos,
						FlowDescriptions: []string{
							"permit out ip from any to any",
							"permit in ip from any to 0.0.0.1",
						},
					},
					&gx.RuleDefinition{
						RuleName:    "dyn_rule_21",
						RatingGroup: &rg21,
						Precedence:  200,
					},
				},
				RuleActivationTime:   &activationTime,
				RuleDeactivationTime: &deactivationTime,
			},
		}

		done <- &gx.CreditControlAnswer{
			ResultCode:     uint32(diameter.SuccessCode),
			SessionID:      request.SessionID,
			RequestNumber:  request.RequestNumber,
			RuleInstallAVP: ruleInstalls,
		}
	}).Once()

	// send rating groups back
	mockPolicyDb.On("GetRuleIDsForBaseNames", []string{"base_10"}).Return([]string{"base_rule_1", "base_rule_2"})
	mockPolicyDb.On("GetChargingKeysForRules", mock.Anything, mock.Anything).Return(
		[]policydb.ChargingKey{
			policydb.ChargingKey{RatingGroup: 1},
			policydb.ChargingKey{RatingGroup: 2},
			policydb.ChargingKey{RatingGroup: 10},
			policydb.ChargingKey{RatingGroup: 11},
			policydb.ChargingKey{RatingGroup: 11},
			policydb.ChargingKey{RatingGroup: 20, ServiceIdTracking: true, ServiceIdentifier: 201},
			policydb.ChargingKey{RatingGroup: 21}}, nil).Once()
	// no omnipresent rules
	mockPolicyDb.On("GetOmnipresentRules").Return([]string{}, []string{}).Once()
	mockPolicyDb.On("GetRuleIDsForBaseNames", mock.Anything).Return([]string{}).Once()
	multiReqType := credit_control.CRTInit // type of CCR sent to get credits
	if initMethod == gy.PerSessionInit {
		mocksGy.On(
			"SendCreditControlRequest",
			mock.Anything,
			mock.Anything,
			mock.MatchedBy(getGyCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTInit)),
		).Return(nil).Run(returnDefaultGyResponse).Once()
		multiReqType = credit_control.CRTUpdate // on per session init, credits are received through CCR-Updates
	}
	// return default responses for gy CCR's, depending on init method
	mocksGy.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(IMSI1_NOPREFIX, multiReqType)),
	).Return(nil).Run(returnDefaultGyResponse).Once()
	createResponse, err := srv.CreateSession(ctx, &protos.CreateSessionRequest{
		Subscriber: &protos.SubscriberID{
			Id: IMSI1,
		},
		SessionId: genSessionID(IMSI1),
	})
	mocksGx.AssertExpectations(t)
	mocksGy.AssertExpectations(t)
	mockPolicyDb.AssertExpectations(t)
	assert.Equal(t, 6, len(createResponse.Credits)) // 2 static, 2 dynamic, 2 base
	assert.Equal(t, 2, len(createResponse.DynamicRules))

	allRuleIDs := []string{}
	for _, staticRule := range createResponse.StaticRules {
		allRuleIDs = append(allRuleIDs, staticRule.RuleId)
		assert.Equal(t, &timestamp.Timestamp{Seconds: 1}, staticRule.ActivationTime)
		assert.Equal(t, &timestamp.Timestamp{Seconds: 2}, staticRule.DeactivationTime)
	}
	assert.ElementsMatch(t, allRuleIDs, []string{"static_rule_1", "static_rule_2", "base_rule_1", "base_rule_2"})

	for _, rule := range createResponse.DynamicRules {
		if rule.PolicyRule.Id == "dyn_rule_20" {
			assert.Equal(t, protos.RedirectInformation_ENABLED, rule.PolicyRule.Redirect.Support)
			assert.Equal(t, protos.RedirectInformation_URL, rule.PolicyRule.Redirect.AddressType)
			assert.Equal(t, "http://www.example.com/", rule.PolicyRule.Redirect.ServerAddress)
			assert.Equal(t, maxReqBWUL, rule.PolicyRule.Qos.MaxReqBwUl)
			assert.Equal(t, maxReqBWDL, rule.PolicyRule.Qos.MaxReqBwDl)
			assert.Equal(t, &timestamp.Timestamp{Seconds: 1}, rule.ActivationTime)
			assert.Equal(t, &timestamp.Timestamp{Seconds: 2}, rule.DeactivationTime)
		} else if rule.PolicyRule.Id == "dyn_rule_21" {
			assert.Empty(t, rule.PolicyRule.Redirect)
			assert.Nil(t, rule.PolicyRule.Qos)
			assert.Equal(t, &timestamp.Timestamp{Seconds: 1}, rule.ActivationTime)
			assert.Equal(t, &timestamp.Timestamp{Seconds: 2}, rule.DeactivationTime)
		} else {
			assert.Fail(t, "Unknown rule id returned")
		}
	}
	ratingGroups := []uint32{}
	for _, update := range createResponse.Credits {
		assert.True(t, update.Success)
		assert.Equal(t, IMSI1, update.Sid)
		ratingGroups = append(ratingGroups, update.ChargingKey)
		if update.ChargingKey == 20 {
			assert.NotNil(t, update.ServiceIdentifier)
			assert.Equal(t, uint32(201), update.GetServiceIdentifier().GetValue())
		} else {
			assert.Nil(t, update.ServiceIdentifier)
		}
		assert.Equal(t, uint64(2048), update.Credit.GrantedUnits.Total.Volume)
		assert.True(t, update.Credit.GrantedUnits.Total.IsValid)
		assert.False(t, update.Credit.GrantedUnits.Rx.IsValid)
		assert.False(t, update.Credit.GrantedUnits.Tx.IsValid)
		assert.Equal(t, uint32(3600), update.Credit.ValidityTime)
	}
	assert.ElementsMatch(t, ratingGroups, []uint32{1, 2, 10, 11, 20, 21})

	// updates
	mocksGy.On("SendCreditControlRequest", mock.Anything, mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(returnDefaultGyResponse).Times(2)

	updateResponse, _ := srv.UpdateSession(ctx,
		&protos.UpdateSessionRequest{
			Updates: []*protos.CreditUsageUpdate{
				createUsageUpdate(IMSI1, 1, 1, protos.CreditUsage_QUOTA_EXHAUSTED),
				createUsageUpdate(IMSI1, 2, 2, protos.CreditUsage_TERMINATED),
			},
		},
	)
	mocksGy.AssertExpectations(t)
	assert.Equal(t, 2, len(updateResponse.Responses))
	for _, update := range updateResponse.Responses {
		assert.True(t, update.Success)
		assert.Equal(t, IMSI1, update.Sid)
		assert.True(t, update.ChargingKey == 1 || update.ChargingKey == 2)
	}

	// Connection Manager tests - Disable Connections for all configured servers
	confNumOfServers := len(controllerParams)
	for i := 0; i < confNumOfServers; i++ {
		mocksGxs[i].On("DisableConnections", mock.Anything).Return()
		mocksGys[i].On("DisableConnections", mock.Anything).Return()
	}
	void, err := srv.Disable(ctx, &fegprotos.DisableMessage{DisablePeriodSecs: 10})
	for i := 0; i < confNumOfServers; i++ {
		mocksGxs[i].AssertExpectations(t)
		mocksGys[i].AssertExpectations(t)
	}
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, void)

	// Connection Manager tests - Enable Connections
	for i := 0; i < confNumOfServers; i++ {
		mocksGxs[i].On("EnableConnections").Return()
		mocksGys[i].On("EnableConnections").Return()
	}
	void, err = srv.Enable(ctx, &orcprotos.Void{})

	for i := 0; i < confNumOfServers; i++ {
		mocksGxs[i].AssertExpectations(t)
		mocksGys[i].AssertExpectations(t)
	}
	assert.NoError(t, err)
	assert.Equal(t, &orcprotos.Void{}, void)

	return nil
}

func TestSessionCreateWithOmnipresentRules(t *testing.T) {
	// Set up mocks
	mockConfig := getTestConfig()
	mockControlParams := getMockControllerParams(mockConfig)
	mockPolicyDb := &MockPolicyDBClient{}
	mockMux := getMockMultiplexor(NUMBER_SERVERS)

	// Get the controller for this imsi
	idx, err := mockMux.GetIndex(multiplex.NewContext().WithIMSI(IMSI1))
	assert.NoError(t, err)
	mocksGx := mockControlParams[idx].PolicyClient.(*MockPolicyClient)

	// send static rules back
	mocksGx.On("SendCreditControlRequest", mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTInit)),
	).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)
		ruleInstalls := []*gx.RuleInstallAVP{
			&gx.RuleInstallAVP{
				RuleNames:     []string{"static_rule_1", "static_rule_2"},
				RuleBaseNames: []string{"base_10"},
			},
		}

		done <- &gx.CreditControlAnswer{
			ResultCode:     uint32(diameter.SuccessCode),
			SessionID:      request.SessionID,
			RequestNumber:  request.RequestNumber,
			RuleInstallAVP: ruleInstalls,
		}
	}).Once()
	mockPolicyDb.On("GetRuleIDsForBaseNames", []string{"base_10"}).Return([]string{"base_rule_1", "base_rule_2"})
	mockPolicyDb.On("GetRuleIDsForBaseNames", []string{"omnipresent_base_1"}).Return([]string{"omnipresent_rule_2"})
	mockPolicyDb.On("GetChargingKeysForRules", mock.Anything, mock.Anything).Return([]policydb.ChargingKey{}, nil).Once()
	mockPolicyDb.On("GetOmnipresentRules").Return([]string{"omnipresent_rule_1"}, []string{"omnipresent_base_1"})
	ctx := context.Background()
	srv := servicers.NewCentralSessionControllers(mockControlParams, mockPolicyDb, mockMux)
	response, err := srv.CreateSession(ctx, &protos.CreateSessionRequest{
		Subscriber: &protos.SubscriberID{
			Id: IMSI1,
		},
		SessionId: genSessionID(IMSI1),
	})
	assert.NoError(t, err)

	mocksGx.AssertExpectations(t)
	mockPolicyDb.AssertExpectations(t)

	assert.Equal(t, 6, len(response.StaticRules))
	expectedRuleIDs := []string{"static_rule_1", "static_rule_2", "base_rule_1", "base_rule_2", "omnipresent_rule_1", "omnipresent_rule_2"}
	actualRuleIDs := funk.Map(response.StaticRules, func(ruleInstall *protos.StaticRuleInstall) string { return ruleInstall.RuleId }).([]string)
	assert.ElementsMatch(t, expectedRuleIDs, actualRuleIDs)
}

func TestSessionControllerTimeouts(t *testing.T) {
	// Set up mocks
	mockConfig := getTestConfig()
	mockControlParams := getMockControllerParams(mockConfig)
	mockPolicyDb := &MockPolicyDBClient{}
	mockMux := getMockMultiplexor(NUMBER_SERVERS)

	srv := servicers.NewCentralSessionControllers(mockControlParams, mockPolicyDb, mockMux)

	idx1, err := mockMux.GetIndex(multiplex.NewContext().WithIMSI(IMSI1))
	assert.NoError(t, err)
	mocksGy_1 := mockControlParams[idx1].CreditClient.(*MockCreditClient)

	idx2, err := mockMux.GetIndex(multiplex.NewContext().WithIMSI(IMSI2))
	assert.NoError(t, err)
	mocksGy_2 := mockControlParams[idx2].CreditClient.(*MockCreditClient)

	ctx := context.Background()

	// depending on request number, "lose" request
	var units uint64 = 2048
	mocksGy_1.On("SendCreditControlRequest", mock.Anything, mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gy.CreditControlRequest)
		if request.RequestNumber%2 == 0 {
			return
		} else {
			done <- &gy.CreditControlAnswer{
				ResultCode:    uint32(diameter.SuccessCode),
				SessionID:     request.SessionID,
				RequestNumber: request.RequestNumber,
				Credits: []*gy.ReceivedCredits{&gy.ReceivedCredits{
					RatingGroup:  request.Credits[0].RatingGroup,
					GrantedUnits: &credit_control.GrantedServiceUnit{TotalOctets: &units},
					ValidityTime: 3600,
				}},
			}
		}
	}).Return(nil).Times(2)

	// This is the answer comming from the second server. NOTE THIS MAY NEED TO BE CHANGED IF idx1 and idx2 are the same
	mocksGy_2.On("SendCreditControlRequest", mock.Anything, mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(IMSI2_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gy.CreditControlRequest)
		if request.RequestNumber%2 == 0 {
			return
		} else {
			done <- &gy.CreditControlAnswer{
				ResultCode:    uint32(diameter.SuccessCode),
				SessionID:     request.SessionID,
				RequestNumber: request.RequestNumber,
				Credits: []*gy.ReceivedCredits{&gy.ReceivedCredits{
					RatingGroup:  request.Credits[0].RatingGroup,
					GrantedUnits: &credit_control.GrantedServiceUnit{TotalOctets: &units},
					ValidityTime: 3600,
				}},
			}
		}
	}).Return(nil).Times(1)

	updateResponse, _ := srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		Updates: []*protos.CreditUsageUpdate{
			createUsageUpdate(IMSI1, 1, 1, protos.CreditUsage_QUOTA_EXHAUSTED),
			createUsageUpdate(IMSI2, 2, 2, protos.CreditUsage_TERMINATED),
			createUsageUpdate(IMSI1, 1, 2, protos.CreditUsage_TERMINATED),
		},
	})
	mocksGy_1.AssertExpectations(t)
	mocksGy_2.AssertExpectations(t)
	assert.Equal(t, 3, len(updateResponse.Responses))
	// Every other request will fail
	countFailed := 0
	for _, update := range updateResponse.Responses {
		if !update.Success {
			countFailed++
		}
	}
	assert.Equal(t, 2, countFailed)
}

func TestSessionTermination(t *testing.T) {
	// Set up mocks
	mockConfig := getTestConfig()
	mockControlParams := getMockControllerParams(mockConfig)
	mockPolicyDb := &MockPolicyDBClient{}
	mockMux := getMockMultiplexor(NUMBER_SERVERS)

	srv := servicers.NewCentralSessionControllers(mockControlParams, mockPolicyDb, mockMux)

	// Get the controller for this imsi
	idx, err := mockMux.GetIndex(multiplex.NewContext().WithIMSI(IMSI2))
	assert.NoError(t, err)
	mocksGx := mockControlParams[idx].PolicyClient.(*MockPolicyClient)
	mocksGy := mockControlParams[idx].CreditClient.(*MockCreditClient)

	ctx := context.Background()

	// Return success for Gx termination
	mocksGx.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI2_NOPREFIX, credit_control.CRTTerminate)),
	).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)
		done <- &gx.CreditControlAnswer{
			ResultCode:    uint32(diameter.SuccessCode),
			SessionID:     request.SessionID,
			RequestNumber: request.RequestNumber,
		}
	}).Once()
	// Return success for Gy terminations
	mocksGy.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(IMSI2_NOPREFIX, credit_control.CRTTerminate)),
	).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gy.CreditControlRequest)
		done <- &gy.CreditControlAnswer{
			ResultCode:    uint32(diameter.SuccessCode),
			SessionID:     request.SessionID,
			RequestNumber: request.RequestNumber,
		}
	}).Once()

	termResponse, err := srv.TerminateSession(ctx, &protos.SessionTerminateRequest{
		Sid:       IMSI2,
		SessionId: genSessionID(IMSI2),
		CreditUsages: []*protos.CreditUsage{
			createUsage(2, protos.CreditUsage_TERMINATED),
			createUsage(1, protos.CreditUsage_TERMINATED),
		},
	})
	mocksGy.AssertExpectations(t)
	mocksGx.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, IMSI2, termResponse.Sid)
	assert.Equal(t, genSessionID(IMSI2), termResponse.SessionId)
}

func TestEventTriggerInUpdate(t *testing.T) {
	// Set up mocks
	mockConfig := getTestConfig()
	mockControlParams := getMockControllerParams(mockConfig)
	mockPolicyDb := &MockPolicyDBClient{}
	mockMux := getMockMultiplexor(NUMBER_SERVERS)

	// Get the controller for this imsi
	idx, err := mockMux.GetIndex(multiplex.NewContext().WithIMSI(IMSI1))
	assert.NoError(t, err)
	mocksGx := mockControlParams[idx].PolicyClient.(*MockPolicyClient)

	// send static rules back
	mocksGx.On("SendCreditControlRequest", mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRUMatcher(IMSI1_NOPREFIX, gx.RevalidationTimeout)),
	).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)
		done <- &gx.CreditControlAnswer{
			ResultCode:    uint32(diameter.SuccessCode),
			SessionID:     request.SessionID,
			RequestNumber: request.RequestNumber,
		}
	}).Once()
	ctx := context.Background()
	srv := servicers.NewCentralSessionControllers(mockControlParams, mockPolicyDb, mockMux)
	_, err = srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{{
			SessionId:     genSessionID(IMSI1),
			Sid:           IMSI1_NOPREFIX,
			RequestNumber: 4,
			EventTrigger:  protos.EventTrigger_REVALIDATION_TIMEOUT,
		}},
	})
	assert.NoError(t, err)

	mocksGx.AssertExpectations(t)
	mockPolicyDb.AssertExpectations(t)
}

func testGxUsageMonitoring(t *testing.T) {
	// Set up mocks
	mockConfig := getTestConfig()
	mockControlParams := getMockControllerParams(mockConfig)
	mockPolicyDb := &MockPolicyDBClient{}
	mockMux := getMockMultiplexor(NUMBER_SERVERS)

	srv := servicers.NewCentralSessionControllers(mockControlParams, mockPolicyDb, mockMux)
	ctx := context.Background()

	// Get the controller for this imsi
	idx_1, err := mockMux.GetIndex(multiplex.NewContext().WithIMSI(IMSI1))
	assert.NoError(t, err)
	mocksGx_1 := mockControlParams[idx_1].PolicyClient.(*MockPolicyClient)
	mocksGy_1 := mockControlParams[idx_1].CreditClient.(*MockCreditClient)

	idx_2, err := mockMux.GetIndex(multiplex.NewContext().WithIMSI(IMSI2))
	assert.NoError(t, err)
	mocksGx_2 := mockControlParams[idx_2].PolicyClient.(*MockPolicyClient)
	mocksGy_2 := mockControlParams[idx_2].CreditClient.(*MockCreditClient)

	// Return success for Gx Update
	mocksGy_1.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(returnDefaultGyResponse).Times(2)
	mocksGx_1.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(returnDefaultGxUpdateResponse).Times(2)

	mocksGy_2.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(IMSI2_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(returnDefaultGyResponse).Times(2)
	mocksGx_2.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI2_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(returnDefaultGxUpdateResponse).Times(2)

	updateResponse, _ := srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		Updates: []*protos.CreditUsageUpdate{
			createUsageUpdate(IMSI1, 1, 1, protos.CreditUsage_QUOTA_EXHAUSTED),
			createUsageUpdate(IMSI1, 2, 2, protos.CreditUsage_TERMINATED),
			createUsageUpdate(IMSI2, 3, 1, protos.CreditUsage_QUOTA_EXHAUSTED),
			createUsageUpdate(IMSI2, 4, 2, protos.CreditUsage_TERMINATED),
		},
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{
			createUsageMonitoringRequest(IMSI1, "mkey", 1, protos.MonitoringLevel_SESSION_LEVEL),
			createUsageMonitoringRequest(IMSI1, "mkey2", 2, protos.MonitoringLevel_PCC_RULE_LEVEL),
			createUsageMonitoringRequest(IMSI2, "mkey3", 1, protos.MonitoringLevel_SESSION_LEVEL),
			createUsageMonitoringRequest(IMSI2, "mkey4", 2, protos.MonitoringLevel_PCC_RULE_LEVEL),
		},
	})

	mocksGy_1.AssertExpectations(t)
	mocksGx_1.AssertExpectations(t)
	mocksGy_2.AssertExpectations(t)
	mocksGx_2.AssertExpectations(t)

	assert.Equal(t, 4, len(updateResponse.Responses))
	assert.Equal(t, 4, len(updateResponse.UsageMonitorResponses))

	for _, update := range updateResponse.Responses {
		assert.True(t, update.Success)
		assert.True(t,
			(IMSI1 == update.Sid && (update.ChargingKey == 1 || update.ChargingKey == 2)) ||
				(IMSI2 == update.Sid && (update.ChargingKey == 3 || update.ChargingKey == 4)),
		)
	}
	for _, update := range updateResponse.UsageMonitorResponses {
		assert.True(t, update.Success)
		assert.True(t, IMSI1 == update.Sid || IMSI2 == update.Sid)
		assert.Equal(t, protos.UsageMonitoringCredit_CONTINUE, update.Credit.Action)
		assert.Equal(t, uint64(2048), update.Credit.GrantedUnits.Total.Volume)
		if string(update.Credit.MonitoringKey) == "mkey" || string(update.Credit.MonitoringKey) == "mkey3" {
			assert.Equal(t, protos.MonitoringLevel_SESSION_LEVEL, update.Credit.Level)
		} else if string(update.Credit.MonitoringKey) == "mkey2" || string(update.Credit.MonitoringKey) == "mkey4" {
			assert.Equal(t, protos.MonitoringLevel_PCC_RULE_LEVEL, update.Credit.Level)
		} else {
			assert.True(t, false)
		}
	}

	// test usage monitoring disabling
	mocksGx_1.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(returnEmptyGxUpdateResponse).Times(1)

	mocksGx_2.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI2_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(returnEmptyGxUpdateResponse).Times(1)

	emptyUpdateResponse, _ := srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{
			createUsageMonitoringRequest(IMSI1, "mkey", 1, protos.MonitoringLevel_SESSION_LEVEL),
			createUsageMonitoringRequest(IMSI2, "mkey3", 1, protos.MonitoringLevel_SESSION_LEVEL),
		},
	})
	mocksGx_1.AssertExpectations(t)
	mocksGx_2.AssertExpectations(t)
	assert.Equal(t, 2, len(emptyUpdateResponse.UsageMonitorResponses))
	for _, update := range emptyUpdateResponse.UsageMonitorResponses {
		assert.True(t, update.Success)
		assert.True(t, IMSI1 == update.Sid || IMSI2 == update.Sid)
		assert.Equal(t, protos.UsageMonitoringCredit_DISABLE, update.Credit.Action)
		assert.Nil(t, update.Credit.GrantedUnits)
		assert.Equal(t, protos.MonitoringLevel_SESSION_LEVEL, update.Credit.Level)
	}

	// Test that static rule install avp in CCA-Update by rule names gets propagated properly
	mocksGx_1.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(getRuleInstallGxUpdateResponse([]string{"static1", "static2"}, []string{})).Times(1)

	mocksGx_2.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI2_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(getRuleInstallGxUpdateResponse([]string{"static3", "static4"}, []string{})).Times(1)

	ruleInstallUpdateResponse, _ := srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{
			createUsageMonitoringRequest(IMSI1, "mkey", 1, protos.MonitoringLevel_SESSION_LEVEL),
			createUsageMonitoringRequest(IMSI2, "mkey3", 1, protos.MonitoringLevel_SESSION_LEVEL),
		},
	})
	mocksGx_1.AssertExpectations(t)
	mocksGx_2.AssertExpectations(t)
	assert.Equal(t, 2, len(ruleInstallUpdateResponse.UsageMonitorResponses))
	for _, update := range ruleInstallUpdateResponse.UsageMonitorResponses {
		assert.True(t, update.Success)
		assert.Nil(t, update.Credit.GrantedUnits)
		if IMSI1 == update.Sid {
			assert.Equal(t, "static1", update.StaticRulesToInstall[0].RuleId)
			assert.Equal(t, "static2", update.StaticRulesToInstall[1].RuleId)
		} else if IMSI2 == update.Sid {
			assert.Equal(t, "static3", update.StaticRulesToInstall[0].RuleId)
			assert.Equal(t, "static4", update.StaticRulesToInstall[1].RuleId)
		} else {
			assert.True(t, false)
		}
	}
	// Test that static rule install avp in CCA-Update by rule base names gets propagated properly
	mocksGx_1.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(getRuleInstallGxUpdateResponse([]string{}, []string{"base_10"})).Times(1)
	mocksGx_2.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI2_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(getRuleInstallGxUpdateResponse([]string{}, []string{"base_30"})).Times(1)

	mockPolicyDb.On("GetRuleIDsForBaseNames", []string{"base_10"}).Return([]string{"base_rule_1", "base_rule_2"})
	mockPolicyDb.On("GetRuleIDsForBaseNames", []string{"base_30"}).Return([]string{"base_rule_2", "base_rule_3"})

	ruleInstallUpdateResponse, _ = srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{
			createUsageMonitoringRequest(IMSI1, "mkey", 1, protos.MonitoringLevel_SESSION_LEVEL),
			createUsageMonitoringRequest(IMSI2, "mkey3", 1, protos.MonitoringLevel_SESSION_LEVEL),
		},
	})
	mocksGx_1.AssertExpectations(t)
	mocksGx_2.AssertExpectations(t)
	assert.Equal(t, 2, len(ruleInstallUpdateResponse.UsageMonitorResponses))
	for _, update := range ruleInstallUpdateResponse.UsageMonitorResponses {
		assert.True(t, update.Success)
		assert.Nil(t, update.Credit.GrantedUnits)
		if IMSI1 == update.Sid {
			assert.Equal(t, "base_rule_1", update.StaticRulesToInstall[0].RuleId)
			assert.Equal(t, "base_rule_2", update.StaticRulesToInstall[1].RuleId)
		} else if IMSI2 == update.Sid {
			assert.Equal(t, "base_rule_3", update.StaticRulesToInstall[0].RuleId)
			assert.Equal(t, "base_rule_4", update.StaticRulesToInstall[1].RuleId)
		} else {
			assert.True(t, false)
		}
	}
	// Test that dynamic rule install avp in CCA-Update gets propagated properly
	mocksGx_1.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(returnDynamicRuleInstallGxUpdateResponse("dyn_rule_10")).Times(1)

	mocksGx_2.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI2_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(returnDynamicRuleInstallGxUpdateResponse("dyn_rule_30")).Times(1)

	ruleInstallUpdateResponse, _ = srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{
			createUsageMonitoringRequest(IMSI1, "mkey", 1, protos.MonitoringLevel_SESSION_LEVEL),
			createUsageMonitoringRequest(IMSI2, "mkey3", 1, protos.MonitoringLevel_SESSION_LEVEL),
		},
	})
	mocksGx_1.AssertExpectations(t)
	mocksGx_2.AssertExpectations(t)

	assert.Equal(t, 2, len(ruleInstallUpdateResponse.UsageMonitorResponses))
	for _, update := range ruleInstallUpdateResponse.UsageMonitorResponses {
		assert.True(t, update.Success)
		assert.Nil(t, update.Credit.GrantedUnits)
		assert.True(t, (IMSI1 == update.Sid || "dyn_rule_10" == update.DynamicRulesToInstall[0].PolicyRule.Id) ||
			(IMSI2 == update.Sid || "dyn_rule_30" == update.DynamicRulesToInstall[0].PolicyRule.Id),
		)
	}

	// Test that rule remove avp in CCA-Update by rule names gets propagated properly
	mocksGx_1.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(getRuleDisableGxUpdateResponse([]string{"rule1", "rule2"}, []string{})).Times(1)

	mocksGx_1.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI2_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(getRuleDisableGxUpdateResponse([]string{"rule3", "rule4"}, []string{})).Times(1)

	ruleDisableUpdateResponse, _ := srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{
			createUsageMonitoringRequest(IMSI1, "mkey", 1, protos.MonitoringLevel_SESSION_LEVEL),
			createUsageMonitoringRequest(IMSI2, "mkey3", 1, protos.MonitoringLevel_SESSION_LEVEL),
		},
	})
	mocksGx_1.AssertExpectations(t)
	mocksGx_2.AssertExpectations(t)
	assert.Equal(t, 2, len(ruleDisableUpdateResponse.UsageMonitorResponses))
	for _, update := range ruleDisableUpdateResponse.UsageMonitorResponses {
		assert.True(t, update.Success)
		assert.Nil(t, update.Credit.GrantedUnits)
		if IMSI1 == update.Sid {
			assert.Equal(t, []string{"rule1", "rule2"}, update.RulesToRemove)
		} else if IMSI2 == update.Sid {
			assert.Equal(t, []string{"rule3", "rule4"}, update.RulesToRemove)
		} else {
			assert.True(t, false)
		}
	}
	// Test that rule remove avp in CCA-Update by base names gets propagated properly
	mocksGx_1.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(getRuleDisableGxUpdateResponse([]string{}, []string{"base_10"})).Times(1)

	mocksGx_2.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI2_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(getRuleDisableGxUpdateResponse([]string{}, []string{"base_30"})).Times(1)

	mockPolicyDb.On("GetRuleIDsForBaseNames", []string{"base_10"}).Return([]string{"base_rule_1", "base_rule_2"})
	mockPolicyDb.On("GetRuleIDsForBaseNames", []string{"base_30"}).Return([]string{"base_rule_3", "base_rule_4"})

	ruleDisableUpdateResponse, _ = srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{
			createUsageMonitoringRequest(IMSI1, "mkey", 1, protos.MonitoringLevel_SESSION_LEVEL),
			createUsageMonitoringRequest(IMSI2, "mkey", 1, protos.MonitoringLevel_SESSION_LEVEL),
		},
	})
	mocksGx_1.AssertExpectations(t)
	mocksGx_2.AssertExpectations(t)
	assert.Equal(t, 2, len(ruleDisableUpdateResponse.UsageMonitorResponses))
	for _, update := range ruleDisableUpdateResponse.UsageMonitorResponses {
		assert.True(t, update.Success)
		assert.Nil(t, update.Credit.GrantedUnits)
		assert.Equal(t, []string{"base_rule_1", "base_rule_2"}, update.RulesToRemove)
		if IMSI1 == update.Sid {
			assert.Equal(t, []string{"base_rule_1", "base_rule_2"}, update.RulesToRemove)
		} else if IMSI2 == update.Sid {
			assert.Equal(t, []string{"base_rule_3", "base_rule_4"}, update.RulesToRemove)
		} else {
			assert.True(t, false)
		}
	}
}

func TestGetHealthStatus(t *testing.T) {
	err := initMconfig()
	assert.NoError(t, err)

	// Set up mocks
	mockConfig := getTestConfig()
	mockControlParams := getMockControllerParams(mockConfig)
	mockPolicyDb := &MockPolicyDBClient{}
	mockMux := getMockMultiplexor(NUMBER_SERVERS)

	srv := servicers.NewCentralSessionControllers(mockControlParams, mockPolicyDb, mockMux)
	ctx := context.Background()

	// Get the controller for two imsis
	idx_1, err := mockMux.GetIndex(multiplex.NewContext().WithIMSI(IMSI1))
	assert.NoError(t, err)
	mocksGx_1 := mockControlParams[idx_1].PolicyClient.(*MockPolicyClient)
	mocksGy_1 := mockControlParams[idx_1].CreditClient.(*MockCreditClient)

	idx_2, err := mockMux.GetIndex(multiplex.NewContext().WithIMSI(IMSI2))
	assert.NoError(t, err)
	mocksGx_2 := mockControlParams[idx_2].PolicyClient.(*MockPolicyClient)
	mocksGy_2 := mockControlParams[idx_2].CreditClient.(*MockCreditClient)

	// Return success for Gx/Gy CCR-Update in two different servers (2 PCRFs, 2 OCSs)
	mocksGy_1.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(returnDefaultGyResponse).Times(2)
	mocksGx_1.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(returnDefaultGxUpdateResponse).Times(2)

	mocksGy_2.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(IMSI2_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(returnDefaultGyResponse).Times(2)
	mocksGx_2.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI2_NOPREFIX, credit_control.CRTUpdate)),
	).Return(nil).Run(returnDefaultGxUpdateResponse).Times(2)

	updateResponse, err := srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		Updates: []*protos.CreditUsageUpdate{
			createUsageUpdate(IMSI1, 1, 1, protos.CreditUsage_QUOTA_EXHAUSTED),
			createUsageUpdate(IMSI1, 2, 2, protos.CreditUsage_TERMINATED),
			createUsageUpdate(IMSI2, 3, 1, protos.CreditUsage_QUOTA_EXHAUSTED),
			createUsageUpdate(IMSI2, 4, 2, protos.CreditUsage_TERMINATED),
		},
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{
			createUsageMonitoringRequest(IMSI1, "mkey", 1, protos.MonitoringLevel_SESSION_LEVEL),
			createUsageMonitoringRequest(IMSI1, "mkey2", 2, protos.MonitoringLevel_PCC_RULE_LEVEL),
			createUsageMonitoringRequest(IMSI2, "mkey3", 1, protos.MonitoringLevel_SESSION_LEVEL),
			createUsageMonitoringRequest(IMSI2, "mkey4", 2, protos.MonitoringLevel_PCC_RULE_LEVEL),
		},
	})
	mocksGy_1.AssertExpectations(t)
	mocksGx_1.AssertExpectations(t)
	mocksGy_2.AssertExpectations(t)
	mocksGx_2.AssertExpectations(t)

	assert.NoError(t, err)
	assert.Equal(t, 4, len(updateResponse.Responses))

	status, err := srv.GetHealthStatus(ctx, &orcprotos.Void{})
	assert.NoError(t, err)
	assert.Equal(t, fegprotos.HealthStatus_HEALTHY, status.Health)

	// Return error for Gx/Gy CCR-Updatee for 2 servers (2 OCSs, 2 PCRFs)
	mocksGy_1.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTUpdate)),
	).Return(fmt.Errorf("Failed to establish new diameter connection; will retry upon first request.")).Run(returnDefaultGyResponse).Times(2)
	mocksGx_1.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTUpdate)),
	).Return(fmt.Errorf("Failed to establish new diameter connection; will retry upon first request.")).Run(returnDefaultGxUpdateResponse).Times(2)

	mocksGy_2.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(IMSI2_NOPREFIX, credit_control.CRTUpdate)),
	).Return(fmt.Errorf("Failed to establish new diameter connection; will retry upon first request.")).Run(returnDefaultGyResponse).Times(2)
	mocksGx_2.On("SendCreditControlRequest",
		mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI2_NOPREFIX, credit_control.CRTUpdate)),
	).Return(fmt.Errorf("Failed to establish new diameter connection; will retry upon first request.")).Run(returnDefaultGxUpdateResponse).Times(2)

	updateResponse, err = srv.UpdateSession(ctx, &protos.UpdateSessionRequest{
		Updates: []*protos.CreditUsageUpdate{
			createUsageUpdate(IMSI1, 1, 1, protos.CreditUsage_QUOTA_EXHAUSTED),
			createUsageUpdate(IMSI1, 2, 2, protos.CreditUsage_TERMINATED),
			createUsageUpdate(IMSI2, 3, 1, protos.CreditUsage_QUOTA_EXHAUSTED),
			createUsageUpdate(IMSI2, 4, 2, protos.CreditUsage_TERMINATED),
		},
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{
			createUsageMonitoringRequest(IMSI1, "mkey", 1, protos.MonitoringLevel_SESSION_LEVEL),
			createUsageMonitoringRequest(IMSI1, "mkey2", 2, protos.MonitoringLevel_PCC_RULE_LEVEL),
			createUsageMonitoringRequest(IMSI2, "mkey3", 1, protos.MonitoringLevel_SESSION_LEVEL),
			createUsageMonitoringRequest(IMSI2, "mkey4", 2, protos.MonitoringLevel_PCC_RULE_LEVEL),
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 4, len(updateResponse.Responses))
	mocksGy_1.AssertExpectations(t)
	mocksGx_1.AssertExpectations(t)
	mocksGy_2.AssertExpectations(t)
	mocksGx_2.AssertExpectations(t)

	status, err = srv.GetHealthStatus(ctx, &orcprotos.Void{})
	assert.NoError(t, err)
	assert.Equal(t, fegprotos.HealthStatus_UNHEALTHY, status.Health)
}

func genSessionID(imsi string) string {
	return fmt.Sprintf("%s-1234", imsi)
}

// getMockMultiplexor loads mockMux with random controlers per each imsi and Imsi without prefix and
// session id (this way we don't need to parse IMSIs at all)
func getMockMultiplexor(numServers int) multiplex.Multiplexor {
	mockMux := &MockMultiplexor{
		imsiToIndex: make(map[uint64]int),
	}
	for i, imsi_uint64 := range ismis_uint64 {
		mockMux.imsiToIndex[imsi_uint64] = i % numServers
	}
	return mockMux
}

// getMockControllerParams generates total of NUMBER_SERVERS . Multiplexor will decide how many will be used
func getMockControllerParams(mockConfig []*servicers.SessionControllerConfig) []*servicers.ControllerParam {
	controlParams := make([]*servicers.ControllerParam, 0, NUMBER_SERVERS)
	for i := 0; i < NUMBER_SERVERS; i++ {
		cp := &servicers.ControllerParam{
			&MockCreditClient{},
			&MockPolicyClient{},
			mockConfig[i],
		}
		controlParams = append(controlParams, cp)
	}
	return controlParams
}

func getTestConfig() []*servicers.SessionControllerConfig {
	serverCfg := make([]*servicers.SessionControllerConfig, len(ocs_server_ports))
	for i := 0; i < NUMBER_SERVERS; i++ {
		ocs_port := ocs_server_ports[i]
		pcrf_port := pcrf_server_ports[i]
		srv := &servicers.SessionControllerConfig{
			OCSConfig: &diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
				Addr:     fmt.Sprintf("127.0.0.1:%s", ocs_port),
				Protocol: "tcp"},
			},
			PCRFConfig: &diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
				Addr:     fmt.Sprintf("127.0.0.1:%s", pcrf_port),
				Protocol: "tcp"},
			},
			RequestTimeout: time.Millisecond,
		}
		serverCfg[i] = srv
	}
	return serverCfg
}

func createUsageUpdate(
	sid string,
	chargingKey uint32,
	requestNumber uint32,
	requestType protos.CreditUsage_UpdateType,
) *protos.CreditUsageUpdate {
	return &protos.CreditUsageUpdate{
		Usage:         createUsage(chargingKey, requestType),
		SessionId:     genSessionID(sid),
		RequestNumber: requestNumber,
		Sid:           sid,
	}
}

func createUsageMonitoringRequest(
	sid string,
	monitoringKey string,
	requestNumber uint32,
	monitoringLevel protos.MonitoringLevel,
) *protos.UsageMonitoringUpdateRequest {
	return &protos.UsageMonitoringUpdateRequest{
		Update: &protos.UsageMonitorUpdate{
			BytesTx:       1024,
			BytesRx:       2048,
			MonitoringKey: []byte(monitoringKey),
			Level:         monitoringLevel,
		},
		SessionId:     genSessionID(sid),
		RequestNumber: requestNumber,
		Sid:           sid,
	}
}

func createUsage(
	chargingKey uint32,
	requestType protos.CreditUsage_UpdateType,
) *protos.CreditUsage {
	return &protos.CreditUsage{
		BytesTx:     1024,
		BytesRx:     2048,
		ChargingKey: chargingKey,
		Type:        requestType,
	}
}

func returnDefaultGyResponse(args mock.Arguments) {
	var units uint64 = 2048
	done := args.Get(1).(chan interface{})
	request := args.Get(2).(*gy.CreditControlRequest)
	credits := make([]*gy.ReceivedCredits, 0, len(request.Credits))

	for _, credit := range request.Credits {
		credits = append(credits, &gy.ReceivedCredits{
			RatingGroup:       credit.RatingGroup,
			ServiceIdentifier: credit.ServiceIdentifier,
			GrantedUnits:      &credit_control.GrantedServiceUnit{TotalOctets: &units},
			ValidityTime:      3600,
			ResultCode:        uint32(diameter.SuccessCode),
		})
	}

	done <- &gy.CreditControlAnswer{
		ResultCode:    uint32(diameter.SuccessCode),
		SessionID:     request.SessionID,
		RequestNumber: request.RequestNumber,
		Credits:       credits,
	}
}

func returnDefaultGxUpdateResponse(args mock.Arguments) {
	done := args.Get(1).(chan interface{})
	request := args.Get(2).(*gx.CreditControlRequest)
	monitors := make([]*gx.UsageMonitoringInfo, 0, len(request.UsageReports))
	for _, report := range request.UsageReports {
		totalOctets := uint64(2048)
		monitors = append(monitors, &gx.UsageMonitoringInfo{
			MonitoringKey: report.MonitoringKey,
			GrantedServiceUnit: &credit_control.GrantedServiceUnit{
				TotalOctets: &totalOctets,
			},
			Level: report.Level,
		})
	}
	done <- &gx.CreditControlAnswer{
		ResultCode:    uint32(diameter.SuccessCode),
		SessionID:     request.SessionID,
		RequestNumber: request.RequestNumber,
		UsageMonitors: monitors,
	}
}

func initMconfig() error {
	fegConfig := `{
      "configsByKey": {
         "session_proxy": {
            "@type": "type.googleapis.com/magma.mconfig.SessionProxyConfig",
            "logLevel": "INFO",
            "gx": {
               "server": {
                   "protocol": "tcp",
                   "address": "",
                   "retransmits": 3,
                   "watchdogInterval": 1,
                   "retryCount": 5,
                   "productName": "magma",
                  "realm": "magma.com",
                  "host": "magma-fedgw.magma.com"
               }
            },
            "gy": {
               "server": {
                   "protocol": "tcp",
                   "address": "",
                   "retransmits": 3,
                   "watchdogInterval": 1,
                   "retryCount": 5,
                   "productName": "magma",
                   "realm": "magma.com",
                   "host": "magma-fedgw.magma.com"
               },
               "initMethod": "PER_KEY"
            },
            "requestFailureThreshold": 0.5,
                "minimumRequestThreshold": 1
         }
      }
   }`

	err := mconfig.CreateLoadTempConfig(fegConfig)
	if err != nil {
		return err
	}
	return nil
}

func returnEmptyGxUpdateResponse(args mock.Arguments) {
	done := args.Get(1).(chan interface{})
	request := args.Get(2).(*gx.CreditControlRequest)
	monitors := make([]*gx.UsageMonitoringInfo, 0, len(request.UsageReports))
	for _, report := range request.UsageReports {
		monitors = append(monitors, &gx.UsageMonitoringInfo{
			MonitoringKey:      report.MonitoringKey,
			GrantedServiceUnit: &credit_control.GrantedServiceUnit{},
			Level:              report.Level,
		})
	}
	done <- &gx.CreditControlAnswer{
		ResultCode:    uint32(diameter.SuccessCode),
		SessionID:     request.SessionID,
		RequestNumber: request.RequestNumber,
		UsageMonitors: monitors,
	}
}

func getRuleInstallGxUpdateResponse(ruleNames, baseNames []string) func(mock.Arguments) {
	return func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)
		monitors := make([]*gx.UsageMonitoringInfo, 0, len(request.UsageReports))
		for _, report := range request.UsageReports {
			monitors = append(monitors, &gx.UsageMonitoringInfo{
				MonitoringKey:      report.MonitoringKey,
				GrantedServiceUnit: &credit_control.GrantedServiceUnit{},
				Level:              report.Level,
			})
		}
		done <- &gx.CreditControlAnswer{
			ResultCode:    uint32(diameter.SuccessCode),
			SessionID:     request.SessionID,
			RequestNumber: request.RequestNumber,
			UsageMonitors: monitors,
			RuleInstallAVP: []*gx.RuleInstallAVP{
				{
					RuleNames:     ruleNames,
					RuleBaseNames: baseNames,
				},
			},
		}
	}
}

func returnDynamicRuleInstallGxUpdateResponse(ruleName string) func(args mock.Arguments) {
	return func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)
		monitors := make([]*gx.UsageMonitoringInfo, 0, len(request.UsageReports))
		for _, report := range request.UsageReports {
			monitors = append(monitors, &gx.UsageMonitoringInfo{
				MonitoringKey:      report.MonitoringKey,
				GrantedServiceUnit: &credit_control.GrantedServiceUnit{},
				Level:              report.Level,
			})
		}
		activationTime := time.Unix(1, 0)
		deactivationTime := time.Unix(2, 0)
		done <- &gx.CreditControlAnswer{
			ResultCode:    uint32(diameter.SuccessCode),
			SessionID:     request.SessionID,
			RequestNumber: request.RequestNumber,
			UsageMonitors: monitors,
			RuleInstallAVP: []*gx.RuleInstallAVP{
				&gx.RuleInstallAVP{
					RuleDefinitions: []*gx.RuleDefinition{
						&gx.RuleDefinition{
							RuleName: ruleName,
							//RatingGroup: swag.Uint32(20),
						},
					},
					RuleActivationTime:   &activationTime,
					RuleDeactivationTime: &deactivationTime,
				},
			},
		}
	}
}

func getRuleDisableGxUpdateResponse(ruleNames []string, ruleBaseNames []string) func(mock.Arguments) {
	return func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)
		monitors := make([]*gx.UsageMonitoringInfo, 0, len(request.UsageReports))
		for _, report := range request.UsageReports {
			monitors = append(monitors, &gx.UsageMonitoringInfo{
				MonitoringKey:      report.MonitoringKey,
				GrantedServiceUnit: &credit_control.GrantedServiceUnit{},
				Level:              report.Level,
			})
		}
		done <- &gx.CreditControlAnswer{
			ResultCode:    uint32(diameter.SuccessCode),
			SessionID:     request.SessionID,
			RequestNumber: request.RequestNumber,
			UsageMonitors: monitors,
			RuleRemoveAVP: []*gx.RuleRemoveAVP{
				{
					RuleNames:     ruleNames,
					RuleBaseNames: ruleBaseNames,
				},
			},
		}
	}
}

func getGyCCRMatcher(imsi string, ccrType credit_control.CreditRequestType) interface{} {
	return func(request *gy.CreditControlRequest) bool {
		return request.Type == ccrType && request.IMSI == imsi
	}
}

func getGxCCRMatcher(imsi string, ccrType credit_control.CreditRequestType) interface{} {
	return func(request *gx.CreditControlRequest) bool {
		return request.Type == ccrType && request.IMSI == imsi
	}
}

func getGxCCRUMatcher(imsi string, eventTrigger gx.EventTrigger) interface{} {
	return func(request *gx.CreditControlRequest) bool {
		return request.Type == credit_control.CRTUpdate &&
			request.IMSI == imsi && request.EventTrigger == eventTrigger

	}
}

/***** UseGyForAuthOnlySuccess Test Cases *****/
func TestSessionControllerUseGyForAuthOnlySuccess(t *testing.T) {
	// Set up mocks
	mockConfig := getTestConfig()
	mockControlParams := getMockControllerParams(mockConfig)
	mockPolicyDb := &MockPolicyDBClient{}
	mockMux := getMockMultiplexor(NUMBER_SERVERS)

	// Get the controller for this imsi
	idx, err := mockMux.GetIndex(multiplex.NewContext().WithIMSI(IMSI1))
	assert.NoError(t, err)
	mocksGx := mockControlParams[idx].PolicyClient.(*MockPolicyClient)
	mocksGy := mockControlParams[idx].CreditClient.(*MockCreditClient)
	mockConfig[idx].UseGyForAuthOnly = true

	activationTime := time.Unix(1, 0)
	deactivationTime := time.Unix(2, 0)
	// send static rules back
	mocksGx.On("SendCreditControlRequest", mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTInit)),
	).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)
		ruleInstalls := []*gx.RuleInstallAVP{
			&gx.RuleInstallAVP{
				RuleNames:            []string{"static_rule_1"},
				RuleActivationTime:   &activationTime,
				RuleDeactivationTime: &deactivationTime,
			},
		}

		done <- &gx.CreditControlAnswer{
			ResultCode:     uint32(diameter.SuccessCode),
			SessionID:      request.SessionID,
			RequestNumber:  request.RequestNumber,
			RuleInstallAVP: ruleInstalls,
		}
	}).Once()

	mockPolicyDb.On("GetChargingKeysForRules", mock.Anything, mock.Anything).Return(
		[]policydb.ChargingKey{{RatingGroup: 3}}, nil).Once()
	mockPolicyDb.On("GetOmnipresentRules").Return([]string{"omnipresent_1"}, []string{}).Once()
	mockPolicyDb.On("GetRuleIDsForBaseNames", []string{}).Return([]string{}).Once()

	mocksGy.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTInit)),
	).Return(nil).Run(returnGySuccessNoRatingGroup).Once()

	srv := servicers.NewCentralSessionControllers(mockControlParams, mockPolicyDb, mockMux)
	ctx := context.Background()

	res, err := srv.CreateSession(ctx, &protos.CreateSessionRequest{
		Subscriber: &protos.SubscriberID{
			Id: IMSI1,
		},
		SessionId: genSessionID(IMSI1),
	})
	mocksGx.AssertExpectations(t)
	assert.NoError(t, err)
	expectedStaticRule1 := &protos.StaticRuleInstall{
		RuleId:           "static_rule_1",
		ActivationTime:   gx.ConvertToProtoTimestamp(&activationTime),
		DeactivationTime: gx.ConvertToProtoTimestamp(&deactivationTime),
	}
	assert.ElementsMatch(t, []*protos.StaticRuleInstall{{RuleId: "omnipresent_1"}, expectedStaticRule1}, res.StaticRules)
}

func TestSessionControllerUseGyForAuthOnlyNoRatingGroup(t *testing.T) {
	// Set up mocks
	mockConfig := getTestConfig()
	mockControlParams := getMockControllerParams(mockConfig)
	mockPolicyDb := &MockPolicyDBClient{}
	mockMux := getMockMultiplexor(NUMBER_SERVERS)

	// Get the controller for this imsi
	idx, err := mockMux.GetIndex(multiplex.NewContext().WithIMSI(IMSI1))
	assert.NoError(t, err)
	mocksGx := mockControlParams[idx].PolicyClient.(*MockPolicyClient)
	mocksGy := mockControlParams[idx].CreditClient.(*MockCreditClient)
	mockConfig[idx].UseGyForAuthOnly = true

	// Send back DIAMETER_SUCCESS (2001) from gx
	mocksGx.On("SendCreditControlRequest", mock.Anything, mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTInit)),
	).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)

		ruleInstalls := []*gx.RuleInstallAVP{
			&gx.RuleInstallAVP{
				RuleNames:       []string{"static_rule_1"},
				RuleDefinitions: []*gx.RuleDefinition{},
			},
		}

		done <- &gx.CreditControlAnswer{
			ResultCode:     uint32(diameter.SuccessCode),
			SessionID:      request.SessionID,
			RequestNumber:  request.RequestNumber,
			RuleInstallAVP: ruleInstalls,
		}
	}).Once()
	mockPolicyDb.On("GetChargingKeysForRules", mock.Anything, mock.Anything).Return(
		[]policydb.ChargingKey{}, nil).Once()
	// no omnipresent rule
	mockPolicyDb.On("GetOmnipresentRules").Return([]string{}, []string{}).Once()
	mockPolicyDb.On("GetRuleIDsForBaseNames", mock.Anything).Return([]string{}).Once()

	// Even if there are no rating groups, gy CCR-I will be called.
	mocksGy.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTInit)),
	).Return(nil).Run(returnGySuccessNoRatingGroup).Once()

	srv := servicers.NewCentralSessionControllers(mockControlParams, mockPolicyDb, mockMux)
	ctx := context.Background()
	_, err = srv.CreateSession(ctx, &protos.CreateSessionRequest{
		Subscriber: &protos.SubscriberID{
			Id: IMSI1,
		},
		SessionId: genSessionID(IMSI1),
	})
	mocksGx.AssertExpectations(t)
	assert.NoError(t, err)
}

func returnGySuccessNoRatingGroup(args mock.Arguments) {
	done := args.Get(1).(chan interface{})
	request := args.Get(2).(*gy.CreditControlRequest)
	credits := make([]*gy.ReceivedCredits, 0, len(request.Credits))
	done <- &gy.CreditControlAnswer{
		ResultCode:    uint32(diameter.SuccessCode),
		SessionID:     request.SessionID,
		RequestNumber: request.RequestNumber,
		Credits:       credits,
	}
}

func TestSessionControllerUseGyForAuthOnlyCreditLimitReached(t *testing.T) {
	// Set up mocks
	mockConfig := getTestConfig()
	mockControlParams := getMockControllerParams(mockConfig)
	mockPolicyDb := &MockPolicyDBClient{}
	mockMux := getMockMultiplexor(NUMBER_SERVERS)

	// Get the controller for this imsi
	idx, err := mockMux.GetIndex(multiplex.NewContext().WithIMSI(IMSI1))
	assert.NoError(t, err)
	mocksGx := mockControlParams[idx].PolicyClient.(*MockPolicyClient)
	mocksGy := mockControlParams[idx].CreditClient.(*MockCreditClient)
	mockConfig[idx].UseGyForAuthOnly = true

	// Send back DIAMETER_SUCCESS (2001) from gx
	mocksGx.On("SendCreditControlRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)

		ruleInstalls := []*gx.RuleInstallAVP{
			&gx.RuleInstallAVP{
				RuleNames:       []string{"static_rule_1"},
				RuleDefinitions: []*gx.RuleDefinition{},
			},
		}

		done <- &gx.CreditControlAnswer{
			ResultCode:     uint32(diameter.SuccessCode),
			SessionID:      request.SessionID,
			RequestNumber:  request.RequestNumber,
			RuleInstallAVP: ruleInstalls,
		}
	}).Once()
	mockPolicyDb.On("GetChargingKeysForRules", mock.Anything, mock.Anything).Return(
		[]policydb.ChargingKey{}, nil).Once()
	// no omnipresent rule
	mockPolicyDb.On("GetOmnipresentRules").Return([]string{}, []string{}).Once()
	mockPolicyDb.On("GetRuleIDsForBaseNames", mock.Anything).Return([]string{}).Once()

	// Even if there are no rating groups, gy CCR-I will be called.
	mocksGy.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTInit)),
	).Return(nil).Run(returnGySuccessCreditLimitReached).Once()

	srv := servicers.NewCentralSessionControllers(mockControlParams, mockPolicyDb, mockMux)
	ctx := context.Background()
	_, err = srv.CreateSession(ctx, &protos.CreateSessionRequest{
		Subscriber: &protos.SubscriberID{
			Id: IMSI1,
		},
		SessionId: genSessionID(IMSI1),
	})
	mocksGx.AssertExpectations(t)
	assert.NoError(t, err)
}

func returnGySuccessCreditLimitReached(args mock.Arguments) {
	done := args.Get(1).(chan interface{})
	request := args.Get(2).(*gy.CreditControlRequest)
	credits := []*gy.ReceivedCredits{
		&gy.ReceivedCredits{
			ResultCode: diameter.DiameterCreditLimitReached,
		},
	}

	done <- &gy.CreditControlAnswer{
		ResultCode:    uint32(diameter.SuccessCode),
		SessionID:     request.SessionID,
		RequestNumber: request.RequestNumber,
		Credits:       credits,
	}
}

func TestSessionControllerUseGyForAuthOnlySubscriberBarred(t *testing.T) {
	// Set up mocks
	mockConfig := getTestConfig()
	mockControlParams := getMockControllerParams(mockConfig)
	mockPolicyDb := &MockPolicyDBClient{}
	mockMux := getMockMultiplexor(NUMBER_SERVERS)

	// Get the controller for this imsi
	idx, err := mockMux.GetIndex(multiplex.NewContext().WithIMSI(IMSI1))
	assert.NoError(t, err)
	mocksGx := mockControlParams[idx].PolicyClient.(*MockPolicyClient)
	mocksGy := mockControlParams[idx].CreditClient.(*MockCreditClient)
	mockConfig[idx].UseGyForAuthOnly = true

	// Send back DIAMETER_SUCCESS (2001) from gx
	mocksGx.On("SendCreditControlRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		done := args.Get(1).(chan interface{})
		request := args.Get(2).(*gx.CreditControlRequest)

		ruleInstalls := []*gx.RuleInstallAVP{
			&gx.RuleInstallAVP{
				RuleNames:       []string{"static_rule_1"},
				RuleDefinitions: []*gx.RuleDefinition{},
			},
		}

		done <- &gx.CreditControlAnswer{
			ResultCode:     uint32(diameter.SuccessCode),
			SessionID:      request.SessionID,
			RequestNumber:  request.RequestNumber,
			RuleInstallAVP: ruleInstalls,
		}
	}).Once()
	mockPolicyDb.On("GetChargingKeysForRules", mock.Anything, mock.Anything).Return(
		[]policydb.ChargingKey{}, nil).Once()
	// no omnipresent rule
	mockPolicyDb.On("GetOmnipresentRules").Return([]string{}, []string{}).Once()
	mockPolicyDb.On("GetRuleIDsForBaseNames", mock.Anything).Return([]string{}).Once()

	// Even if there are no rating groups, gy CCR-I will be called.
	mocksGy.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGyCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTInit)),
	).Return(nil).Run(returnGySuccessSubscriberBarred).Once()

	srv := servicers.NewCentralSessionControllers(mockControlParams, mockPolicyDb, mockMux)
	ctx := context.Background()
	_, err = srv.CreateSession(ctx, &protos.CreateSessionRequest{
		Subscriber: &protos.SubscriberID{
			Id: IMSI1,
		},
		SessionId: genSessionID(IMSI1),
	})
	mocksGx.AssertExpectations(t)
	assert.Error(t, err)
}

func returnGySuccessSubscriberBarred(args mock.Arguments) {
	done := args.Get(1).(chan interface{})
	request := args.Get(2).(*gy.CreditControlRequest)
	credits := []*gy.ReceivedCredits{
		&gy.ReceivedCredits{
			ResultCode: diameter.DiameterRatingFailed,
		},
	}

	done <- &gy.CreditControlAnswer{
		ResultCode:    uint32(diameter.SuccessCode),
		SessionID:     request.SessionID,
		RequestNumber: request.RequestNumber,
		Credits:       credits,
	}
}

func returnGxSuccessRevalidationTimer(args mock.Arguments) {
	done := args.Get(1).(chan interface{})
	request := args.Get(2).(*gx.CreditControlRequest)
	ruleInstalls := []*gx.RuleInstallAVP{
		&gx.RuleInstallAVP{
			RuleNames: []string{"static_rule_1"},
		},
	}
	mkey := []byte("key")
	totalOctets := uint64(2048)
	monitors := []*gx.UsageMonitoringInfo{
		{
			MonitoringKey:      mkey,
			GrantedServiceUnit: &credit_control.GrantedServiceUnit{TotalOctets: &totalOctets},
		},
	}
	revalidationTime := time.Unix(1, 0)
	eventTrigger := []gx.EventTrigger{gx.RevalidationTimeout}

	done <- &gx.CreditControlAnswer{
		ResultCode:       uint32(diameter.SuccessCode),
		SessionID:        request.SessionID,
		RequestNumber:    request.RequestNumber,
		RuleInstallAVP:   ruleInstalls,
		UsageMonitors:    monitors,
		EventTriggers:    eventTrigger,
		RevalidationTime: &revalidationTime,
	}
}

func revalidationTimerTest(
	t *testing.T,
	srv servicers.CentralSessionControllerServerWithHealth,
	controllerParams []*servicers.ControllerParam,
	policyDb policydb.PolicyDBClient,
	mux multiplex.Multiplexor,
	useGyForAuthOnly bool,
	numberServers int,
) {
	ctx := context.Background()
	mockPolicyDb := policyDb.(*MockPolicyDBClient)

	// Get the controller for this imsi
	idx, err := mux.GetIndex(multiplex.NewContext().WithIMSI(IMSI1))
	assert.NoError(t, err)
	mocksGx := controllerParams[idx].PolicyClient.(*MockPolicyClient)
	mocksGy := controllerParams[idx].CreditClient.(*MockCreditClient)

	mocksGx.On(
		"SendCreditControlRequest",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(getGxCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTInit)),
	).Return(nil).Run(returnGxSuccessRevalidationTimer).Once()

	mockPolicyDb.On("GetOmnipresentRules").Return([]string{"omnipresent_rule_1"}, []string{"omnipresent_base_1"})
	mockPolicyDb.On("GetRuleIDsForBaseNames", []string{"omnipresent_base_1"}).Return([]string{"omnipresent_rule_2"})
	mockPolicyDb.On("GetChargingKeysForRules", mock.Anything, mock.Anything).Return([]policydb.ChargingKey{}, nil).Once()

	if useGyForAuthOnly {
		mocksGy.On(
			"SendCreditControlRequest",
			mock.Anything,
			mock.Anything,
			mock.MatchedBy(getGyCCRMatcher(IMSI1_NOPREFIX, credit_control.CRTInit)),
		).Return(nil).Run(returnGySuccessNoRatingGroup).Once()
	}

	createResponse, err := srv.CreateSession(ctx, &protos.CreateSessionRequest{
		Subscriber: &protos.SubscriberID{
			Id: IMSI1,
		},
		SessionId: genSessionID(IMSI1),
	})

	mocksGx.AssertExpectations(t)
	mocksGy.AssertExpectations(t)
	mockPolicyDb.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(createResponse.UsageMonitors))

	for _, monitor := range createResponse.GetUsageMonitors() {
		assert.Equal(t, &timestamp.Timestamp{Seconds: 1}, monitor.GetRevalidationTime())
		assert.ElementsMatch(t, monitor.GetEventTriggers(), []protos.EventTrigger{protos.EventTrigger_REVALIDATION_TIMEOUT})
	}
}

func TestSessionControllerRevalidationTimerUsed(t *testing.T) {
	// Set up mocks
	mockConfig := getTestConfig()
	mockControlParams := getMockControllerParams(mockConfig)
	mockPolicyDb := &MockPolicyDBClient{}
	mockMux := getMockMultiplexor(NUMBER_SERVERS)

	srv := servicers.NewCentralSessionControllers(mockControlParams, mockPolicyDb, mockMux)

	revalidationTimerTest(t, srv, mockControlParams, mockPolicyDb, mockMux, false, NUMBER_SERVERS)
}

func TestSessionControllerUseGyForAuthOnlyRevalidationTimerUsed(t *testing.T) {

	numberServers := 1
	mockConfig := getTestConfig()
	mockConfig[0].UseGyForAuthOnly = true
	mockControlParams := getMockControllerParams(mockConfig)
	mockPolicyDb := &MockPolicyDBClient{}
	mockMux := getMockMultiplexor(numberServers)
	srv := servicers.NewCentralSessionControllers(mockControlParams, mockPolicyDb, mockMux)

	revalidationTimerTest(t, srv, mockControlParams, mockPolicyDb, mockMux, mockConfig[0].UseGyForAuthOnly, 1)
}
