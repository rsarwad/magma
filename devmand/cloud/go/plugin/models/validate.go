/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"github.com/go-openapi/strfmt"
)

func (m *SymphonyAgent) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *ManagedDevices) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *SymphonyDevice) ValidateModel() error {
	return m.Validate(strfmt.Default)
}
