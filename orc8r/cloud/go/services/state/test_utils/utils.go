package test_utils

import (
	"encoding/json"
	"testing"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/service/middleware/unary/test_utils"
	"magma/orc8r/cloud/go/services/state"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func ReportGatewayStatus(t *testing.T, ctx context.Context, req *models.GatewayStatus) {
	client, err := state.GetStateClient()
	assert.NoError(t, err)

	serializedGWStatus, err := serde.Serialize(state.SerdeDomain, orc8r.GatewayStateType, req)
	assert.NoError(t, err)
	states := []*protos.State{
		{
			Type:     orc8r.GatewayStateType,
			DeviceID: req.HardwareID,
			Value:    serializedGWStatus,
		},
	}
	_, err = client.ReportStates(
		ctx,
		&protos.ReportStatesRequest{States: states},
	)
	assert.NoError(t, err)
}

func GetContextWithCertificate(t *testing.T, hwID string) context.Context {
	csn := test_utils.StartMockGwAccessControl(t, []string{hwID})
	return metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs(identity.CLIENT_CERT_SN_KEY, csn[0]))
}

func GetGWStatusViaHTTPNoError(t *testing.T, url string, networkID string, key string) {
	status, response, err := tests.SendHttpRequest("GET", url, "")
	assert.NoError(t, err)
	assert.Equal(t, 200, status)
	expected, err := state.GetGatewayStatus(networkID, key)
	assert.NoError(t, err)
	expectedJSON, err := json.Marshal(*expected)
	assert.NoError(t, err)
	assert.Equal(t, string(expectedJSON), response)
}

func GetGWStatusExpectNotFound(t *testing.T, url string) {
	status, _, err := tests.SendHttpRequest("GET", url, "")
	assert.NoError(t, err)
	assert.Equal(t, 404, status)
}
