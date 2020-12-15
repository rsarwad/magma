/*
 Copyright 2020 The Magma Authors.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package servicers_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"magma/lte/cloud/go/lte"
	lte_plugin "magma/lte/cloud/go/plugin"
	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/ha/servicers"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/analytics/query_api/mocks"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	state_test_init "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/cloud/go/storage"
	orc8r_protos "magma/orc8r/lib/go/protos"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	//_ = flag.Set("alsologtostderr", "true") // uncomment to view logs during test
}

func TestHAServicer_GetEnodebOffloadState(t *testing.T) {
	assert.NoError(t, plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{}))
	assert.NoError(t, plugin.RegisterPluginForTests(t, &lte_plugin.LteOrchestratorPlugin{}))
	configurator_test_init.StartTestService(t)
	state_test_init.StartTestService(t)
	mockPromClient := &mocks.PrometheusAPI{}
	servicer := servicers.NewHAServicer(mockPromClient)

	testNetworkId := "n1"
	testGwHwId1 := "hw1"
	testGwId1 := "g1"
	testGwHwId2 := "hw2"
	testGwId2 := "g2"
	testGwPool := "pool1"
	enbSn := "enb1"
	err := configurator.CreateNetwork(configurator.Network{ID: testNetworkId}, serdes.Network)
	assert.NoError(t, err)

	// Initialize HA network topology
	_, err = configurator.CreateEntity(
		testNetworkId,
		configurator.NetworkEntity{
			Type:   lte.CellularEnodebEntityType,
			Key:    enbSn,
			Config: newDefaultUnmanagedEnodebConfig(),
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	_, err = configurator.CreateEntities(
		testNetworkId,
		[]configurator.NetworkEntity{
			{
				Type: lte.CellularGatewayEntityType, Key: testGwId1,
				Config:       newDefaultGatewayConfig(1, 255),
				Associations: []storage.TypeAndKey{{Type: lte.CellularEnodebEntityType, Key: enbSn}},
			},
			{
				Type: orc8r.MagmadGatewayType, Key: testGwId1,
				Name: "foobar", Description: "foo bar",
				PhysicalID:   testGwHwId1,
				Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: testGwId1}},
			},
			{
				Type: lte.CellularGatewayEntityType, Key: testGwId2,
				Config:       newDefaultGatewayConfig(2, 1),
				Associations: []storage.TypeAndKey{{Type: lte.CellularEnodebEntityType, Key: enbSn}},
			},
			{
				Type: orc8r.MagmadGatewayType, Key: testGwId2,
				Name: "foobar2", Description: "foo bar",
				PhysicalID:   testGwHwId2,
				Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: testGwId2}},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	_, err = configurator.CreateEntities(
		testNetworkId,
		[]configurator.NetworkEntity{
			{
				Type: lte.CellularGatewayPoolEntityType, Key: testGwPool,
				Config: &lte_models.CellularGatewayPoolConfigs{
					MmeGroupID: 1,
				},
				Associations: []storage.TypeAndKey{
					{Type: lte.CellularGatewayEntityType, Key: testGwId1},
					{Type: lte.CellularGatewayEntityType, Key: testGwId2},
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	// Initialize network state for given devices
	gwStatus := &models.GatewayStatus{
		CheckinTime: uint64(time.Now().Unix()),
		HardwareID:  testGwHwId1,
	}
	ctx1 := test_utils.GetContextWithCertificate(t, testGwHwId1)
	test_utils.ReportGatewayStatus(t, ctx1, gwStatus)
	gwStatus2 := &models.GatewayStatus{
		CheckinTime: uint64(time.Now().Unix()),
		HardwareID:  testGwHwId2,
	}
	ctx2 := test_utils.GetContextWithCertificate(t, testGwHwId2)
	test_utils.ReportGatewayStatus(t, ctx2, gwStatus2)

	enbState := getDefaultEnodebState(testGwId1)
	reportEnodebState(t, ctx1, enbSn, enbState)

	metric1 := model.Metric{}
	metric1["networkID"] = "n1"
	metric1["gatewayID"] = "g1"
	throughputVec := model.Vector{{
		Metric: metric1,
		Value:  10000000,
	}}
	mockPromClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(throughputVec, nil, nil).Once()

	// First test primary healthy
	ctx := orc8r_protos.NewGatewayIdentity(testGwHwId2, testNetworkId, testGwId2).NewContextWithIdentity(context.Background())
	res, err := servicer.GetEnodebOffloadState(ctx, &protos.GetEnodebOffloadStateRequest{})
	assert.NoError(t, err)
	expectedRes := &protos.GetEnodebOffloadStateResponse{
		EnodebOffloadStates: map[uint32]protos.GetEnodebOffloadStateResponse_EnodebOffloadState{
			138: protos.GetEnodebOffloadStateResponse_PRIMARY_CONNECTED_AND_SERVING_UES,
		},
	}
	assert.Equal(t, expectedRes, res)
	mockPromClient.AssertExpectations(t)

	// Now simulate failed primary not checking in
	stateTooOld := time.Now().Add(-time.Second * 600)
	clock.SetAndFreezeClock(t, stateTooOld)
	test_utils.ReportGatewayStatus(t, ctx1, gwStatus)

	res, err = servicer.GetEnodebOffloadState(ctx, &protos.GetEnodebOffloadStateRequest{})
	assert.NoError(t, err)
	expectedRes = &protos.GetEnodebOffloadStateResponse{
		EnodebOffloadStates: map[uint32]protos.GetEnodebOffloadStateResponse_EnodebOffloadState{},
	}
	assert.Equal(t, expectedRes, res)
	clock.UnfreezeClock(t)
	test_utils.ReportGatewayStatus(t, ctx1, gwStatus)

	// Simulate too old of ENB state
	clock.SetAndFreezeClock(t, stateTooOld)
	reportEnodebState(t, ctx1, enbSn, enbState)

	res, err = servicer.GetEnodebOffloadState(ctx, &protos.GetEnodebOffloadStateRequest{})
	assert.NoError(t, err)
	expectedRes = &protos.GetEnodebOffloadStateResponse{
		EnodebOffloadStates: map[uint32]protos.GetEnodebOffloadStateResponse_EnodebOffloadState{
			138: protos.GetEnodebOffloadStateResponse_NO_OP,
		},
	}
	assert.Equal(t, expectedRes, res)
	clock.UnfreezeClock(t)

	// ENB not connected
	enbState.EnodebConnected = swag.Bool(false)
	reportEnodebState(t, ctx1, enbSn, enbState)

	res, err = servicer.GetEnodebOffloadState(ctx, &protos.GetEnodebOffloadStateRequest{})
	assert.NoError(t, err)
	expectedRes = &protos.GetEnodebOffloadStateResponse{
		EnodebOffloadStates: map[uint32]protos.GetEnodebOffloadStateResponse_EnodebOffloadState{
			138: protos.GetEnodebOffloadStateResponse_NO_OP,
		},
	}
	assert.Equal(t, expectedRes, res)

	// Connected but could not query metrics
	enbState.EnodebConnected = swag.Bool(true)
	reportEnodebState(t, ctx1, enbSn, enbState)
	mockPromClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil, fmt.Errorf("error")).Once()

	res, err = servicer.GetEnodebOffloadState(ctx, &protos.GetEnodebOffloadStateRequest{})
	assert.NoError(t, err)
	expectedRes = &protos.GetEnodebOffloadStateResponse{
		EnodebOffloadStates: map[uint32]protos.GetEnodebOffloadStateResponse_EnodebOffloadState{
			138: protos.GetEnodebOffloadStateResponse_PRIMARY_CONNECTED,
		},
	}
	assert.Equal(t, expectedRes, res)
	mockPromClient.AssertExpectations(t)

	// Connected but no throughput
	throughputVec2 := model.Vector{{
		Metric: metric1,
		Value:  0,
	}}
	mockPromClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(throughputVec2, nil, nil).Once()

	res, err = servicer.GetEnodebOffloadState(ctx, &protos.GetEnodebOffloadStateRequest{})
	assert.NoError(t, err)
	expectedRes = &protos.GetEnodebOffloadStateResponse{
		EnodebOffloadStates: map[uint32]protos.GetEnodebOffloadStateResponse_EnodebOffloadState{
			138: protos.GetEnodebOffloadStateResponse_PRIMARY_CONNECTED,
		},
	}
	assert.Equal(t, expectedRes, res)
	mockPromClient.AssertExpectations(t)

	// Back to connected with users
	mockPromClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(throughputVec, nil, nil).Once()
	res, err = servicer.GetEnodebOffloadState(ctx, &protos.GetEnodebOffloadStateRequest{})
	assert.NoError(t, err)
	expectedRes = &protos.GetEnodebOffloadStateResponse{
		EnodebOffloadStates: map[uint32]protos.GetEnodebOffloadStateResponse_EnodebOffloadState{
			138: protos.GetEnodebOffloadStateResponse_PRIMARY_CONNECTED_AND_SERVING_UES,
		},
	}
	assert.Equal(t, expectedRes, res)
	mockPromClient.AssertExpectations(t)
}

func reportEnodebState(t *testing.T, ctx context.Context, enodebSerial string, req *lte_models.EnodebState) {
	client, err := state.GetStateClient()
	assert.NoError(t, err)

	serializedEnodebState, err := serde.Serialize(req, lte.EnodebStateType, serdes.State)
	assert.NoError(t, err)
	states := []*orc8r_protos.State{
		{
			Type:     lte.EnodebStateType,
			DeviceID: enodebSerial,
			Value:    serializedEnodebState,
		},
	}
	_, err = client.ReportStates(
		ctx,
		&orc8r_protos.ReportStatesRequest{States: states},
	)
	assert.NoError(t, err)
}

func newDefaultGatewayConfig(mmeCode uint32, mmeRelCap uint32) *lte_models.GatewayCellularConfigs {
	return &lte_models.GatewayCellularConfigs{
		Ran: &lte_models.GatewayRanConfigs{
			Pci:             260,
			TransmitEnabled: swag.Bool(true),
		},
		Epc: &lte_models.GatewayEpcConfigs{
			NatEnabled: swag.Bool(true),
			IPBlock:    "192.168.128.0/24",
		},
		NonEpsService: &lte_models.GatewayNonEpsConfigs{
			CsfbMcc:              "001",
			CsfbMnc:              "01",
			Lac:                  swag.Uint32(1),
			CsfbRat:              swag.Uint32(0),
			Arfcn2g:              nil,
			NonEpsServiceControl: swag.Uint32(0),
		},
		DNS: &lte_models.GatewayDNSConfigs{
			DhcpServerEnabled: swag.Bool(true),
			EnableCaching:     swag.Bool(false),
			LocalTTL:          swag.Int32(0),
		},
		HeConfig: &lte_models.GatewayHeConfig{},
		Pooling: lte_models.CellularGatewayPoolRecords{
			{
				GatewayPoolID:       "pool1",
				MmeCode:             mmeCode,
				MmeRelativeCapacity: mmeRelCap,
			},
		},
	}
}

func newDefaultUnmanagedEnodebConfig() *lte_models.EnodebConfig {
	ip := strfmt.IPv4("192.168.0.124")
	return &lte_models.EnodebConfig{
		ConfigType: "UNMANAGED",
		UnmanagedConfig: &lte_models.UnmanagedEnodebConfiguration{
			CellID:    swag.Uint32(138),
			Tac:       swag.Uint32(1),
			IPAddress: &ip,
		},
	}
}

func getDefaultEnodebState(gwID string) *lte_models.EnodebState {
	return &lte_models.EnodebState{
		MmeConnected:       swag.Bool(true),
		EnodebConnected:    swag.Bool(true),
		IPAddress:          "10.0.0.1",
		ReportingGatewayID: gwID,
		EnodebConfigured:   swag.Bool(true),
		GpsConnected:       swag.Bool(true),
		GpsLatitude:        swag.String("foo"),
		GpsLongitude:       swag.String("bar"),
		OpstateEnabled:     swag.Bool(true),
		PtpConnected:       swag.Bool(true),
		RfTxOn:             swag.Bool(true),
		RfTxDesired:        swag.Bool(true),
		FsmState:           swag.String("abc"),
	}
}
