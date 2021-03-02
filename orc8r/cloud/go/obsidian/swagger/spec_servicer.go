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

package swagger

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"magma/orc8r/cloud/go/obsidian/swagger/protos"

	"github.com/golang/glog"
)

type specServicer struct {
	spec string
}

// NewSpecServicerFromFile intializes a spec servicer
// given a service name.
func NewSpecServicerFromFile(service string) protos.SwaggerSpecServer {
	service = strings.ToLower(service)
	path := getSpecPath(service)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		// Swallowing ReadFile error because the service should
		// continue to run even if it can't find its Swagger spec file.
		glog.Errorf("Error retrieving Swagger Spec of service %s: %+v", service, err)
		return NewSpecServicer("")
	}
	return NewSpecServicer(string(data))
}

// NewSpecServicer constructs a spec servicer.
func NewSpecServicer(spec string) protos.SwaggerSpecServer {
	return &specServicer{spec}
}

func (s *specServicer) GetSpec(ctx context.Context, request *protos.GetSpecRequest) (*protos.GetSpecResponse, error) {
	return &protos.GetSpecResponse{SwaggerSpec: s.spec}, nil
}

// getSpecPath returns the filepath on the production image
// that contains the service's Swagger spec
func getSpecPath(service string) string {
	specDir := "/etc/magma/swagger/specs"
	specPath := filepath.Join(specDir, fmt.Sprintf("%s.swagger.v1.yml", service))
	return specPath
}
