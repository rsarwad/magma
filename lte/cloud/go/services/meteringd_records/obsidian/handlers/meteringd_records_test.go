/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Tests for Meteringd REST Endpoints
package handlers_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	lteplugin "magma/lte/cloud/go/plugin"
	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/meteringd_records"
	"magma/lte/cloud/go/services/meteringd_records/obsidian/models"
	meteringd_records_test_init "magma/lte/cloud/go/services/meteringd_records/test_init"
	sdb_test_init "magma/lte/cloud/go/services/subscriberdb/test_init"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	orcprotos "magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/service/middleware/unary/test_utils"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

// Update flows in tables as a gateway would
// NOTE: This endpoint exists for testing ONLY
// Real clients will use gRPC directly
func UpdateFlowsTest(csn string, tbl *protos.FlowTable) error {
	client, err := meteringd_records.GetMeteringdRecordsClient()
	if err != nil {
		return err
	}

	// Hack in the identity context
	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("x-magma-client-cert-serial", csn))
	_, err = client.UpdateFlows(ctx, tbl)
	return err
}

// TestMeteringdRecords is Obsidian Metering Records Integration Test intended to be run
// on cloud VM
func TestMeteringdRecords(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	_ = plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	sdb_test_init.StartTestService(t)
	meteringd_records_test_init.StartTestService(t)
	restPort := tests.StartObsidian(t)

	hwId := "TestAGHwId00003"
	csn := test_utils.StartMockGwAccessControl(t, []string{hwId})

	testUrlRoot := fmt.Sprintf(
		"http://localhost:%d%s/networks", restPort, obsidian.RestRoot)

	// Test Register Network
	registerNetworkTestCase := tests.Testcase{
		Name:                      "Register Network",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s?requested_id=meteringd_records_obsidian_test_network", testUrlRoot),
		Payload:                   `{"name":"This Is A Test Network Name"}`,
		Skip_payload_verification: true,
	}
	_, networkId, _ := tests.RunTest(t, registerNetworkTestCase)
	sId := "IMSI12333333333"
	json.Unmarshal([]byte(networkId), &networkId)

	// Test Register AG
	registerAGTestCase := tests.Testcase{
		Name:     "Register AG",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/gateways", testUrlRoot, networkId),
		Payload:  fmt.Sprintf(`{"hardware_id":"%s", "key": {"key_type": "ECHO"}}}`, hwId),
		Expected: fmt.Sprintf(`"%s"`, hwId),
	}
	tests.RunTest(t, registerAGTestCase)

	// Test Add Subscriber
	addSubTestCase := tests.Testcase{
		Name:   "Add Subscriber",
		Method: "POST",
		Url:    fmt.Sprintf("%s/%s/subscribers", testUrlRoot, networkId),
		Payload: fmt.Sprintf(`{"id":"%s",
           "lte":{"state":"ACTIVE",
           "auth_algo":"MILENAGE",
           "auth_key":"AAAAAAAAAAAAAAAAAAAAAA==",
           "auth_opc":"AAECAwQFBgcICQoLDA0ODw=="}}`,
			sId),
		Expected: fmt.Sprintf(`"%s"`, sId),
	}
	tests.RunTest(t, addSubTestCase)

	// Create fake gateway context ids
	id := &orcprotos.Identity{}
	idGw := orcprotos.Identity_Gateway{
		HardwareId: hwId,
		NetworkId:  networkId,
		LogicalId:  fmt.Sprintf(`"%s"`, hwId),
	}
	id.SetGateway(&idGw)

	// Create flow for sub on gateway
	recordId := &protos.FlowRecord_ID{Id: "test"}
	record := &protos.FlowRecord{Id: recordId, Sid: sId, BytesTx: 1554, BytesRx: 1553, PktsTx: 1234, PktsRx: 5432}
	tbl := &protos.FlowTable{}
	tbl.Flows = append(tbl.Flows, record)
	err := UpdateFlowsTest(csn[0], tbl)
	assert.NoError(t, err)

	// Test Listing All Subscriber Flow Records
	expectedRecord := &models.FlowRecord{}
	err = expectedRecord.FromProto(record)
	assert.NoError(t, err)
	marshaledRecord, err := expectedRecord.MarshalBinary()
	assert.NoError(t, err)
	expected := string(marshaledRecord)
	listFlowRecordsTestCase := tests.Testcase{
		Name:   "List All Subscriber Flow Records",
		Method: "GET",
		Url: fmt.Sprintf("%s/%s/subscribers/%s/flow_records",
			testUrlRoot, networkId, sId),
		Payload:  "",
		Expected: fmt.Sprintf("[%s]", expected),
	}
	tests.RunTest(t, listFlowRecordsTestCase)

	// Test Get Flow Records
	getFlowRecordTestCase := tests.Testcase{
		Name:     "Get Flow Record",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/flow_records/test", testUrlRoot, networkId),
		Payload:  "",
		Expected: fmt.Sprintf(`{"bytes_rx":1553,"bytes_tx":1554,"pkts_rx":5432,"pkts_tx":1234,"subscriber_id":"%s"}`, sId),
	}
	tests.RunTest(t, getFlowRecordTestCase)

}
