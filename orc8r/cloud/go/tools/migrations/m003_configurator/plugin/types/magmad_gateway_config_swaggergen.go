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

// Code generated by go-swagger; DO NOT EDIT.

package types

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// MagmadGatewayConfig magmad gateway config
// swagger:model magmad_gateway_config
type MagmadGatewayConfig struct {

	// autoupgrade enabled
	AutoupgradeEnabled bool `json:"autoupgrade_enabled,omitempty"`

	// autoupgrade poll interval
	AutoupgradePollInterval int32 `json:"autoupgrade_poll_interval,omitempty"`

	// checkin interval
	CheckinInterval int32 `json:"checkin_interval,omitempty"`

	// checkin timeout
	CheckinTimeout int32 `json:"checkin_timeout,omitempty"`

	// dynamic services
	DynamicServices []string `json:"dynamic_services"`

	// feature flags
	FeatureFlags map[string]bool `json:"feature_flags,omitempty"`

	// ID of tier within network that gateway is grouped into
	// Min Length: 1
	Tier string `json:"tier,omitempty"`
}

// Validate validates this magmad gateway config
func (m *MagmadGatewayConfig) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTier(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *MagmadGatewayConfig) validateTier(formats strfmt.Registry) error {

	if swag.IsZero(m.Tier) { // not required
		return nil
	}

	if err := validate.MinLength("tier", "body", string(m.Tier), 1); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *MagmadGatewayConfig) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MagmadGatewayConfig) UnmarshalBinary(b []byte) error {
	var res MagmadGatewayConfig
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
