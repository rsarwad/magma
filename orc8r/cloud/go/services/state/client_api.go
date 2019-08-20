/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package state

import (
	"context"
	"encoding/json"
	"sync"

	"magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/serde"

	"github.com/golang/glog"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc"
)

// State includes reported operational state and additional info about the reporter
type State struct {
	// ID of the entity reporting the state (hwID, cert serial number, etc)
	ReporterID string
	// Checkin Time
	Time uint64
	// Cert expiration Time
	CertExpirationTime int64
	ReportedState      interface{}
}

// SerializedStateWithMeta includes reported operational states and additional info
type SerializedStateWithMeta struct {
	ReporterID              string
	Time                    uint64
	CertExpirationTime      int64
	SerializedReportedState []byte
}

// StateID contains the identifying information of a state
type StateID struct {
	Type     string
	DeviceID string
}

// Global clientconn that can be reused for this service
var connSingleton = (*grpc.ClientConn)(nil)
var connGuard = sync.Mutex{}

func getStateClient() (protos.StateServiceClient, error) {
	if connSingleton == nil {
		// Reading the conn optimistically to avoid unnecessary overhead
		connGuard.Lock()
		if connSingleton == nil {
			conn, err := registry.GetConnection(ServiceName)
			if err != nil {
				initErr := errors.NewInitError(err, ServiceName)
				glog.Error(initErr)
				connGuard.Unlock()
				return nil, initErr
			}
			connSingleton = conn
		}
		connGuard.Unlock()
	}
	return protos.NewStateServiceClient(connSingleton), nil
}

// GetState returns the state specified by the networkID, typeVal, and hwID
func GetState(networkID string, typeVal string, hwID string) (State, error) {
	client, err := getStateClient()
	if err != nil {
		return State{}, err
	}

	stateID := &protos.StateID{
		Type:     typeVal,
		DeviceID: hwID,
	}

	ret, err := client.GetStates(
		context.Background(),
		&protos.GetStatesRequest{
			NetworkID: networkID,
			Ids:       []*protos.StateID{stateID},
		},
	)
	if err != nil {
		return State{}, err
	}
	if len(ret.States) == 0 {
		return State{}, errors.ErrNotFound
	}
	return toState(ret.States[0])
}

// GetStates returns a map of states specified by the networkID and a list of type and key
func GetStates(networkID string, stateIDs []StateID) (map[StateID]State, error) {
	client, err := getStateClient()
	if err != nil {
		return nil, err
	}

	res, err := client.GetStates(
		context.Background(), &protos.GetStatesRequest{
			NetworkID: networkID,
			Ids:       toProtosStateIDs(stateIDs),
		},
	)
	if err != nil {
		return nil, err
	}

	idToValue := map[StateID]State{}
	for _, pState := range res.States {
		stateID := StateID{Type: pState.Type, DeviceID: pState.DeviceID}
		state, err := toState(pState)
		if err != nil {
			return nil, err
		}
		idToValue[stateID] = state
	}
	return idToValue, nil
}

// DeleteStates deletes states specified by the networkID and a list of type and key
func DeleteStates(networkID string, stateIDs []StateID) error {
	client, err := getStateClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteStates(
		context.Background(),
		&protos.DeleteStatesRequest{
			NetworkID: networkID,
			Ids:       toProtosStateIDs(stateIDs),
		},
	)
	return err
}

func GetGatewayStatus(networkID string, deviceID string) (*models.GatewayStatus, error) {
	state, err := GetState(networkID, orc8r.GatewayStateType, deviceID)
	if err != nil {
		return nil, err
	}
	if state.ReportedState == nil {
		return nil, errors.ErrNotFound
	}
	return fillInGatewayStatusState(state), nil
}

func GetGatewayStatuses(networkID string, deviceIDs []string) (map[string]*models.GatewayStatus, error) {
	stateIDs := funk.Map(deviceIDs, func(id string) StateID { return StateID{Type: orc8r.GatewayStateType, DeviceID: id} }).([]StateID)
	res, err := GetStates(networkID, stateIDs)
	if err != nil {
		return map[string]*models.GatewayStatus{}, err
	}

	ret := make(map[string]*models.GatewayStatus, len(res))
	for stateID, state := range res {
		ret[stateID.DeviceID] = fillInGatewayStatusState(state)
	}
	return ret, nil
}

func fillInGatewayStatusState(state State) *models.GatewayStatus {
	if state.ReportedState == nil {
		return nil
	}

	gwStatus := state.ReportedState.(*models.GatewayStatus)
	gwStatus.CheckinTime = state.Time
	gwStatus.CertExpirationTime = state.CertExpirationTime
	gwStatus.HardwareID = state.ReporterID
	return gwStatus
}

func toProtosStateIDs(stateIDs []StateID) []*protos.StateID {
	ids := []*protos.StateID{}
	for _, state := range stateIDs {
		ids = append(ids, &protos.StateID{Type: state.Type, DeviceID: state.DeviceID})
	}
	return ids
}

func toState(pState *protos.State) (State, error) {
	serialized := &SerializedStateWithMeta{}
	err := json.Unmarshal(pState.Value, serialized)
	if err != nil {
		return State{}, err
	}
	iReportedState, err := serde.Deserialize(SerdeDomain, pState.Type, serialized.SerializedReportedState)
	state := State{
		ReporterID:         serialized.ReporterID,
		Time:               serialized.Time,
		CertExpirationTime: serialized.CertExpirationTime,
		ReportedState:      iReportedState,
	}
	return state, err
}
