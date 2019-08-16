/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 */

package tests

import (
	"context"
	"encoding/json"
	"testing"

	"magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/orc8r"
	models2 "magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/serde"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
	configuratorTestUtils "magma/orc8r/cloud/go/services/configurator/test_utils"
	"magma/orc8r/cloud/go/services/device"
	deviceTestInit "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/state"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"

	"github.com/golang/glog"
	"github.com/stretchr/testify/assert"
)

const (
	typeName   = "typeName"
	testAgHwId = "Test-AGW-Hw-Id"
)

type stateBundle struct {
	state *protos.State
	ID    state.StateID
}

func makeStateBundle(typeVal string, key string, value interface{}) stateBundle {
	marshaledValue, _ := json.Marshal(value)
	ID := state.StateID{Type: typeVal, DeviceID: key}
	state := protos.State{Type: typeVal, DeviceID: key, Value: marshaledValue}
	return stateBundle{state: &state, ID: ID}
}

func TestStateService(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	// Set up test networkID, hwID, and encode into context
	stateTestInit.StartTestService(t)
	err := serde.RegisterSerdes(
		&Serde{},
		serde.NewBinarySerde(device.SerdeDomain, orc8r.AccessGatewayRecordType, &models2.GatewayDevice{}))
	assert.NoError(t, err)

	networkID := "state_service_test_network"
	configuratorTestUtils.RegisterNetwork(t, networkID, "State Service Test")
	gatewayID := testAgHwId
	configuratorTestUtils.RegisterGateway(t, networkID, gatewayID, &models2.GatewayDevice{HardwareID: testAgHwId})
	ctx := test_utils.GetContextWithCertificate(t, testAgHwId)

	// Create States, IDs, values
	value0 := Name{Name: "name0"}
	value1 := Name{Name: "name1"}
	value2 := NameAndAge{Name: "name2", Age: 20}
	bundle0 := makeStateBundle(typeName, "key0", value0)
	bundle1 := makeStateBundle(typeName, "key1", value1)
	bundle2 := makeStateBundle(typeName, "key2", value2)

	// Check contract for empty network
	states, err := state.GetStates(networkID, []state.StateID{bundle0.ID})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(states))

	// Report and read back
	_, err = reportStates(ctx, bundle0, bundle1)
	assert.NoError(t, err)
	states, err = state.GetStates(networkID, []state.StateID{bundle0.ID, bundle1.ID})
	assert.NoError(t, err)
	testGetStatesResponse(t, states, bundle0, bundle1)

	// Report a state with fields the corresponding serde does not expect
	_, err = reportStates(ctx, bundle2)
	assert.NoError(t, err)
	states, err = state.GetStates(networkID, []state.StateID{bundle2.ID})
	assert.NoError(t, err)
	testGetStatesResponse(t, states, bundle2)

	// Delete and read back
	err = state.DeleteStates(networkID, []state.StateID{bundle0.ID, bundle2.ID})
	assert.NoError(t, err)
	states, err = state.GetStates(networkID, []state.StateID{bundle0.ID, bundle1.ID, bundle2.ID})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(states))
	testGetStatesResponse(t, states, bundle1)

	// Send a valid state and a state with no corresponding serde
	unserializableBundle := makeStateBundle("nonexistent-serde", "key3", value0)
	resp, err := reportStates(ctx, bundle0, unserializableBundle)
	assert.NoError(t, err)
	assert.Equal(t, "nonexistent-serde", resp.UnreportedStates[0].Type)
	assert.Equal(t, "No Serde found for type nonexistent-serde", resp.UnreportedStates[0].Error)
	// Valid state should still be reported
	states, err = state.GetStates(networkID, []state.StateID{bundle0.ID, bundle1.ID, bundle2.ID})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(states))
	testGetStatesResponse(t, states, bundle0, bundle1)
}

type NameAndAge struct {
	// name
	Name string `json:"name"`
	// age
	Age int `json:"age"`
}

type Name struct {
	// name
	Name string `json:"name"`
}

type Serde struct {
}

func (*Serde) GetDomain() string {
	return state.SerdeDomain
}

func (*Serde) GetType() string {
	return typeName
}

func (*Serde) Serialize(in interface{}) ([]byte, error) {
	return json.Marshal(in)

}

func (*Serde) Deserialize(message []byte) (interface{}, error) {
	res := Name{}
	err := json.Unmarshal(message, &res)
	return res, err
}

func getClient() (protos.StateServiceClient, error) {
	conn, err := registry.GetConnection(state.ServiceName)
	if err != nil {
		initErr := errors.NewInitError(err, state.ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewStateServiceClient(conn), err
}

func reportStates(ctx context.Context, bundles ...stateBundle) (*protos.ReportStatesResponse, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}
	response, err := client.ReportStates(ctx, makeReportStatesRequest(bundles))
	return response, err
}

func testGetStatesResponse(t *testing.T, states map[state.StateID]state.State, bundles ...stateBundle) {
	for _, bundle := range bundles {
		value := states[bundle.ID]
		iState, err := serde.Deserialize(state.SerdeDomain, bundle.ID.Type, bundle.state.Value)
		assert.NoError(t, err)
		assert.Equal(t, iState, value.ReportedState)
	}
}

func makeReportStatesRequest(bundles []stateBundle) *protos.ReportStatesRequest {
	res := protos.ReportStatesRequest{}
	res.States = makeStates(bundles)
	return &res
}

func makeStates(bundles []stateBundle) []*protos.State {
	states := []*protos.State{}
	for _, bundle := range bundles {
		states = append(states, bundle.state)
	}
	return states
}
