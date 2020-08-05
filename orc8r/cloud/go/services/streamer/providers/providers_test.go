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

package providers_test

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/mocks"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	orchestrator_test_init "magma/orc8r/cloud/go/services/orchestrator/test_init"
	"magma/orc8r/cloud/go/services/streamer"
	streamer_test_init "magma/orc8r/cloud/go/services/streamer/test_init"
	"magma/orc8r/cloud/go/services/streamer/test_utils/mconfig/factory"
	"magma/orc8r/cloud/go/services/streamer/test_utils/mconfig/test_protos"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/mock"
	assert "github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestMconfigStreamer_Configurator(t *testing.T) {
	configurator_test_init.StartTestService(t)
	streamer_test_init.StartTestService(t)
	orchestrator_test_init.StartTestService(t)

	// Set up mock mconfig builders (legacy and new)
	configurator.ClearMconfigBuilders(t)
	mockBuilder := &mocks.MconfigBuilder{}
	mockBuilder.On("Build", "n1", "gw1", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			out := args[4].(map[string]proto.Message)
			out["new_builder"] = &test_protos.Message1{Field: "hello"}
		})
	configurator.RegisterMconfigBuilders(mockBuilder)

	mockLegacyBuilder := &mockMconfigBuilder{
		retVal: map[string]proto.Message{
			"builder1_1": &test_protos.Message1{Field: "hello"},
			"builder1_2": &test_protos.Message2{Field1: "hello", Field2: "world"},
		},
	}
	factory.ClearMconfigBuilders(t)
	factory.RegisterMconfigBuilders(mockLegacyBuilder)

	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "gw1", PhysicalID: "hw1"})
	assert.NoError(t, err)

	conn, err := registry.GetConnection(streamer.ServiceName)
	assert.NoError(t, err)
	grpcClient := protos.NewStreamerClient(conn)

	// Make normal call for config updates
	extraArgs := &protos.GatewayConfigsDigest{Md5HexDigest: "useless_digest"}
	serializedExtraArgs, _ := ptypes.MarshalAny(extraArgs)
	stream, err := grpcClient.GetUpdates(
		context.Background(),
		&protos.StreamRequest{GatewayId: "hw1", StreamName: "configs", ExtraArgs: serializedExtraArgs},
	)
	assert.NoError(t, err)

	expectedProtos := map[string]proto.Message{
		"new_builder": &test_protos.Message1{Field: "hello"},
	}
	expected := make(map[string]*any.Any, len(expectedProtos))
	for k, v := range expectedProtos {
		anyV, err := ptypes.MarshalAny(v)
		assert.NoError(t, err)
		expected[k] = anyV
	}
	actualMarshaled, err := stream.Recv()
	assert.NoError(t, err)
	actual := &protos.GatewayConfigs{}
	err = protos.Unmarshal(actualMarshaled.Updates[0].Value, actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual.ConfigsByKey)

	// Make optimized call for config updates--when passed config digest
	// matches provider's digest, empty update batch is returned
	extraArgs = &protos.GatewayConfigsDigest{Md5HexDigest: actual.Metadata.Digest.Md5HexDigest}
	serializedExtraArgs, _ = ptypes.MarshalAny(extraArgs)
	stream, err = grpcClient.GetUpdates(
		context.Background(),
		&protos.StreamRequest{GatewayId: "hw1", StreamName: "configs", ExtraArgs: serializedExtraArgs},
	)
	assert.NoError(t, err)

	actualMarshaled, err = stream.Recv()
	assert.NoError(t, err)
	assert.Empty(t, actualMarshaled.Updates)

	mockBuilder.AssertExpectations(t)
}

type mockMconfigBuilder struct {
	retVal map[string]proto.Message
	retErr error
}

func (builder *mockMconfigBuilder) Build(networkId string, gatewayId string) (map[string]proto.Message, error) {
	return builder.retVal, builder.retErr
}
