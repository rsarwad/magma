/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handlers_test

import (
	"testing"

	"magma/orc8r/cloud/go/obsidian/swagger"
	"magma/orc8r/cloud/go/obsidian/swagger/handlers"
	"magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	swagger_lib "magma/orc8r/cloud/go/swagger"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/registry"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func Test_GenerateCombinedSpecHandler(t *testing.T) {
	e := echo.New()
	testURLRoot := "/magma/v1"

	commonTag := swagger_lib.TagDefinition{Name: "Tag Common"}
	commonSpec := swagger_lib.Spec{
		Tags: []swagger_lib.TagDefinition{commonTag},
	}
	yamlCommon := marshalToYAML(t, commonSpec)

	// Success with no registered servicers
	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		Handler:        handlers.GetGenerateCombinedSpecHandler(yamlCommon),
		ExpectedStatus: 200,
		ExpectedResult: commonSpec,
	}
	tests.RunUnitTest(t, e, tc)

	tags := []swagger_lib.TagDefinition{
		{Name: "Tag 1"},
		{Name: "Tag 2"},
		{Name: "Tag 3"},
	}
	services := []string{"test_spec_service1", "test_spec_service2", "test_spec_service3"}
	expectedSpec := swagger_lib.Spec{
		Tags: []swagger_lib.TagDefinition{tags[0], tags[1], tags[2], commonTag},
	}

	// Clean up registry
	defer registry.RemoveServicesWithLabel(orc8r.SwaggerSpecLabel)

	for i, s := range services {
		registerServicer(t, s, tags[i])
	}

	// Success with registered servicers
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		Handler:        handlers.GetGenerateCombinedSpecHandler(yamlCommon),
		ExpectedStatus: 200,
		ExpectedResult: expectedSpec,
	}
	tests.RunUnitTest(t, e, tc)
}

func registerServicer(t *testing.T, service string, tag swagger_lib.TagDefinition) {
	labels := map[string]string{
		orc8r.SwaggerSpecLabel: "true",
	}

	srv, lis := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, service, labels, nil)
	spec := swagger_lib.Spec{Tags: []swagger_lib.TagDefinition{tag}}

	yamlSpec := marshalToYAML(t, spec)
	protos.RegisterSwaggerSpecServer(srv.GrpcServer, swagger.NewSpecServicer(yamlSpec))

	go srv.RunTest(lis)
}

func marshalToYAML(t *testing.T, spec swagger_lib.Spec) string {
	data, err := yaml.Marshal(&spec)
	assert.NoError(t, err)
	return string(data)
}
