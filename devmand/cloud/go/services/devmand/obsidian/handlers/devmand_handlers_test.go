/*
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
*/

package handlers_test

import (
	"fmt"
	"os"
	"testing"

	"magma/orc8r/cloud/go/obsidian"
	obsidian_test "magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	configurator_test_utils "magma/orc8r/cloud/go/services/configurator/test_utils"
	"orc8r/devmand/cloud/go/devmand"
	devmandp "orc8r/devmand/cloud/go/plugin"
	"orc8r/devmand/cloud/go/services/devmand/test_utils"

	"github.com/stretchr/testify/assert"
)

func TestGetDeviceConfigs(t *testing.T) {
	os.Setenv(orc8r.UseConfiguratorEnv, "1")
	plugin.RegisterPluginForTests(t, &devmandp.DevmandOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testURLRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, obsidian.RestRoot)

	networkID := "Test Network 1"
	deviceID := "test_device_1"
	configurator_test_utils.RegisterNetwork(t, networkID, "")
	registerDevice(t, networkID, deviceID)

	// Happy path
	expectedConfig := test_utils.NewDefaultManagedDevice()
	marshaledCfg, err := expectedConfig.MarshalBinary()
	assert.NoError(t, err)
	expected := string(marshaledCfg)

	createConfigTestCase := obsidian_test.Testcase{
		Name:     "Create Device Config",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/configs/devices?requested_id=%s", testURLRoot, networkID, deviceID),
		Payload:  expected,
		Expected: `"test_device_1"`,
	}
	obsidian_test.RunTest(t, createConfigTestCase)

	happyPathTestCase := obsidian_test.Testcase{
		Name:     "Get Device Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/configs/devices/%s", testURLRoot, networkID, deviceID),
		Payload:  "",
		Expected: expected,
	}
	obsidian_test.RunTest(t, happyPathTestCase)
}

func TestSetDeviceConfigs(t *testing.T) {
	os.Setenv(orc8r.UseConfiguratorEnv, "1")
	plugin.RegisterPluginForTests(t, &devmandp.DevmandOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testURLRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, obsidian.RestRoot)

	networkID := "Test Network 1"
	deviceID := "test_device_1"
	configurator_test_utils.RegisterNetwork(t, networkID, "")
	registerDevice(t, networkID, deviceID)

	// Happy path
	deviceConfig := test_utils.NewDefaultManagedDevice()
	marshaledCfg, err := deviceConfig.MarshalBinary()
	assert.NoError(t, err)
	expected := string(marshaledCfg)

	createConfigTestCase := obsidian_test.Testcase{
		Name:     "Create Device Config",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/configs/devices?requested_id=%s", testURLRoot, networkID, deviceID),
		Payload:  expected,
		Expected: `"test_device_1"`,
	}
	obsidian_test.RunTest(t, createConfigTestCase)

	deviceConfig.Host = "0.0.0.0"
	marshaledCfg, err = deviceConfig.MarshalBinary()
	assert.NoError(t, err)
	configString := string(marshaledCfg)

	setConfigTestCase := obsidian_test.Testcase{
		Name:     "Set Device Config",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s/configs/devices/%s", testURLRoot, networkID, deviceID),
		Payload:  configString,
		Expected: "",
	}
	obsidian_test.RunTest(t, setConfigTestCase)
	happyPathTestCase := obsidian_test.Testcase{
		Name:     "Get Device Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/configs/devices/%s", testURLRoot, networkID, deviceID),
		Payload:  "",
		Expected: configString,
	}
	obsidian_test.RunTest(t, happyPathTestCase)
}

func TestGetGatewayConfigs(t *testing.T) {
	os.Setenv(orc8r.UseConfiguratorEnv, "1")
	plugin.RegisterPluginForTests(t, &devmandp.DevmandOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testURLRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, obsidian.RestRoot)

	networkID := "Test Network 1"
	gatewayID := "g1"
	d1 := "test_device_1"
	d2 := "test_device_2"
	configurator_test_utils.RegisterNetwork(t, networkID, "")
	configurator_test_utils.RegisterGateway(t, networkID, gatewayID, nil)
	registerDevice(t, networkID, d1)
	registerDevice(t, networkID, d2)

	// Happy path
	expectedConfig := test_utils.NewDefaultGatewayConfig()
	marshaledCfg, err := expectedConfig.MarshalBinary()
	assert.NoError(t, err)
	expected := string(marshaledCfg)

	createConfigTestCase := obsidian_test.Testcase{
		Name:     "Create Devmand Gateway Config",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/devmand", testURLRoot, networkID, gatewayID),
		Payload:  expected,
		Expected: `"g1"`,
	}
	obsidian_test.RunTest(t, createConfigTestCase)

	happyPathTestCase := obsidian_test.Testcase{
		Name:     "Get Devmand Gateway Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/devmand", testURLRoot, networkID, gatewayID),
		Payload:  "",
		Expected: expected,
	}
	obsidian_test.RunTest(t, happyPathTestCase)
}

func TestSetGatewayConfigs(t *testing.T) {
	os.Setenv(orc8r.UseConfiguratorEnv, "1")
	plugin.RegisterPluginForTests(t, &devmandp.DevmandOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testURLRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, obsidian.RestRoot)

	networkID := "Test Network 1"
	gatewayID := "g1"
	d1 := "test_device_1"
	d2 := "test_device_2"
	configurator_test_utils.RegisterNetwork(t, networkID, "")
	configurator_test_utils.RegisterGateway(t, networkID, gatewayID, nil)
	registerDevice(t, networkID, d1)
	registerDevice(t, networkID, d2)

	gatewayConfig := test_utils.NewDefaultGatewayConfig()
	marshaledCfg, err := gatewayConfig.MarshalBinary()
	assert.NoError(t, err)
	configString := string(marshaledCfg)

	createConfigTestCase := obsidian_test.Testcase{
		Name:     "Create Devmand Gateway Config",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/devmand", testURLRoot, networkID, gatewayID),
		Payload:  configString,
		Expected: `"g1"`,
	}
	obsidian_test.RunTest(t, createConfigTestCase)

	// Should fail if device is not registered
	gatewayConfig.ManagedDevices = []string{"test_device_1", "test_device_2", "test_device_3"}
	marshaledCfg, err = gatewayConfig.MarshalBinary()
	assert.NoError(t, err)
	configString = string(marshaledCfg)

	setConfigTestCase := obsidian_test.Testcase{
		Name:                     "Set Devmand Gateway Config Without Device Registered",
		Method:                   "PUT",
		Url:                      fmt.Sprintf("%s/%s/gateways/%s/configs/devmand", testURLRoot, networkID, gatewayID),
		Payload:                  configString,
		Expected:                 `{"message":"could not find entities matching [type:\"device\" key:\"test_device_3\" ]"}`,
		Expect_http_error_status: true,
	}

	d3 := "test_device_3"
	registerDevice(t, networkID, d3)

	setConfigTestCase = obsidian_test.Testcase{
		Name:    "Set Devmand Gateway Config",
		Method:  "PUT",
		Url:     fmt.Sprintf("%s/%s/gateways/%s/configs/devmand", testURLRoot, networkID, gatewayID),
		Payload: configString,
	}
	obsidian_test.RunTest(t, setConfigTestCase)

	getConfigTestCase := obsidian_test.Testcase{
		Name:     "Get Updated Devmand Gateway Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/devmand", testURLRoot, networkID, gatewayID),
		Payload:  "",
		Expected: configString,
	}
	obsidian_test.RunTest(t, getConfigTestCase)

	// remove devices from config
	gatewayConfig.ManagedDevices = []string{"test_device_1"}
	marshaledCfg, err = gatewayConfig.MarshalBinary()
	assert.NoError(t, err)
	configString = string(marshaledCfg)
	setConfigTestCase = obsidian_test.Testcase{
		Name:    "Set Devmand Gateway Config To Delete Device Association",
		Method:  "PUT",
		Url:     fmt.Sprintf("%s/%s/gateways/%s/configs/devmand", testURLRoot, networkID, gatewayID),
		Payload: configString,
	}
	obsidian_test.RunTest(t, setConfigTestCase)

	getConfigTestCase = obsidian_test.Testcase{
		Name:     "Get Updated Devmand Gateway Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/devmand", testURLRoot, networkID, gatewayID),
		Payload:  "",
		Expected: configString,
	}
	obsidian_test.RunTest(t, getConfigTestCase)

	// remove device entity and see configs get updated
	err = configurator.DeleteEntity(networkID, devmand.DeviceType, d1)
	assert.NoError(t, err)
	gatewayConfig.ManagedDevices = []string{}
	marshaledCfg, err = gatewayConfig.MarshalBinary()
	assert.NoError(t, err)
	configString = string(marshaledCfg)
	getConfigTestCase = obsidian_test.Testcase{
		Name:     "Get Updated Devmand Gateway Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/devmand", testURLRoot, networkID, gatewayID),
		Payload:  "",
		Expected: configString,
	}
	obsidian_test.RunTest(t, getConfigTestCase)
}

func TestDeleteGatewayConfigs(t *testing.T) {
	os.Setenv(orc8r.UseConfiguratorEnv, "1")
	plugin.RegisterPluginForTests(t, &devmandp.DevmandOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testURLRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, obsidian.RestRoot)

	networkID := "Test Network 1"
	gatewayID := "g1"
	d1 := "test_device_1"
	d2 := "test_device_2"
	configurator_test_utils.RegisterNetwork(t, networkID, "")
	configurator_test_utils.RegisterGateway(t, networkID, gatewayID, nil)
	registerDevice(t, networkID, d1)
	registerDevice(t, networkID, d2)

	gatewayConfig := test_utils.NewDefaultGatewayConfig()
	gatewayConfig.ManagedDevices = []string{d1, d2}
	marshaledCfg, err := gatewayConfig.MarshalBinary()
	assert.NoError(t, err)
	configString := string(marshaledCfg)

	createConfigTestCase := obsidian_test.Testcase{
		Name:     "Create Devmand Gateway Config",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/devmand", testURLRoot, networkID, gatewayID),
		Payload:  configString,
		Expected: `"g1"`,
	}
	obsidian_test.RunTest(t, createConfigTestCase)

	deleteConfigTestCase := obsidian_test.Testcase{
		Name:     "Delete Devmand Gateway Config",
		Method:   "DELETE",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/devmand", testURLRoot, networkID, gatewayID),
		Payload:  "",
		Expected: "",
	}
	obsidian_test.RunTest(t, deleteConfigTestCase)

	// test device loadedEntities still exist
	deviceType := "device"
	loadedEntities, _, err := configurator.LoadEntities(networkID, &deviceType, nil, nil, nil, configurator.EntityLoadCriteria{})
	assert.Equal(t, 2, len(loadedEntities))
}

func registerDevice(t *testing.T, networkID string, deviceID string) {
	entity := configurator.NetworkEntity{
		Key:  deviceID,
		Type: devmand.DeviceType,
	}
	_, err := configurator.CreateEntity(networkID, entity)
	assert.NoError(t, err)
}
