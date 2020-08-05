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

package models

import (
	"magma/orc8r/cloud/go/obsidian/models"

	"github.com/go-openapi/strfmt"
)

const (
	lteAuthKeyLength = 16
	lteAuthOpcLength = 16
)

func (m *LteSubscription) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}

	authKeyLen := len([]byte(m.AuthKey))
	if authKeyLen != lteAuthKeyLength {
		return models.ValidateErrorf("expected lte auth key to be %d bytes but got %d bytes", lteAuthKeyLength, authKeyLen)
	}

	// OPc is optional, but if it's provided it should be 16 bytes
	authOpcLen := len([]byte(m.AuthOpc))
	if authOpcLen > 0 && authOpcLen != lteAuthOpcLength {
		return models.ValidateErrorf("expected lte auth opc to be %d bytes but got %d bytes", lteAuthOpcLength, authOpcLen)
	}

	return nil
}

func (m *MutableSubscriber) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	if err := m.Lte.ValidateModel(); err != nil {
		return err
	}
	return nil
}

func (m *IcmpStatus) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}
