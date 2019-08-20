/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package device

import (
	"context"

	merrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/device/protos"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

func getDeviceClient() (protos.DeviceClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewDeviceClient(conn), err
}

func RegisterDevice(networkID, deviceType, deviceKey string, info interface{}) error {
	client, err := getDeviceClient()
	if err != nil {
		return err
	}

	serializedInfo, err := serde.Serialize(SerdeDomain, deviceType, info)
	if err != nil {
		return err
	}
	entity := &protos.PhysicalEntity{
		DeviceID: deviceKey,
		Type:     deviceType,
		Info:     serializedInfo,
	}
	req := &protos.RegisterOrUpdateDevicesRequest{
		NetworkID: networkID,
		Entities:  []*protos.PhysicalEntity{entity},
	}
	_, err = client.RegisterDevices(context.Background(), req)
	return err
}

func UpdateDevice(networkID, deviceType, deviceKey string, info interface{}) error {
	client, err := getDeviceClient()
	if err != nil {
		return err
	}

	serializedInfo, err := serde.Serialize(SerdeDomain, deviceType, info)
	if err != nil {
		return err
	}
	entity := &protos.PhysicalEntity{
		DeviceID: deviceKey,
		Type:     deviceType,
		Info:     serializedInfo,
	}
	req := &protos.RegisterOrUpdateDevicesRequest{
		NetworkID: networkID,
		Entities:  []*protos.PhysicalEntity{entity},
	}
	_, err = client.UpdateDevices(context.Background(), req)
	return err
}

func DeleteDevices(networkID string, deviceIDs []*protos.DeviceID) error {
	client, err := getDeviceClient()
	if err != nil {
		return err
	}

	req := &protos.DeleteDevicesRequest{NetworkID: networkID, DeviceIDs: deviceIDs}
	_, err = client.DeleteDevices(context.Background(), req)
	return err
}

func DeleteDevice(networkID, deviceType, deviceKey string) error {
	return DeleteDevices(networkID, []*protos.DeviceID{{DeviceID: deviceKey, Type: deviceType}})
}

func GetDevice(networkID, deviceType, deviceKey string) (interface{}, error) {
	client, err := getDeviceClient()
	if err != nil {
		return nil, err
	}
	deviceID := &protos.DeviceID{Type: deviceType, DeviceID: deviceKey}
	req := &protos.GetDeviceInfoRequest{NetworkID: networkID, DeviceIDs: []*protos.DeviceID{deviceID}}
	res, err := client.GetDeviceInfo(context.Background(), req)
	if err != nil {
		return nil, err
	}
	device, ok := res.DeviceMap[deviceKey]
	if !ok {
		return nil, merrors.ErrNotFound
	}
	return serde.Deserialize(SerdeDomain, deviceType, device.Info)
}

func GetDevices(networkID string, deviceType string, deviceIDs []string) (map[string]interface{}, error) {
	client, err := getDeviceClient()
	if err != nil {
		return nil, err
	}

	requestIDs := funk.Map(
		deviceIDs,
		func(id string) *protos.DeviceID { return &protos.DeviceID{Type: deviceType, DeviceID: id} },
	).([]*protos.DeviceID)
	req := &protos.GetDeviceInfoRequest{NetworkID: networkID, DeviceIDs: requestIDs}
	res, err := client.GetDeviceInfo(context.Background(), req)
	if err != nil {
		return map[string]interface{}{}, err
	}

	ret := make(map[string]interface{}, len(res.DeviceMap))
	for k, val := range res.DeviceMap {
		iVal, err := serde.Deserialize(SerdeDomain, deviceType, val.Info)
		if err != nil {
			return map[string]interface{}{}, errors.Wrapf(err, "failed to deserialize device %s", k)
		}
		ret[k] = iVal
	}
	return ret, nil
}

func DoesDeviceExist(networkID, deviceType, deviceID string) bool {
	_, err := GetDevice(networkID, deviceType, deviceID)
	if err != nil {
		return false
	}
	return true
}
