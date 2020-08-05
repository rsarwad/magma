/**
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
 *
 * @flow strict-local
 * @format
 */
import type {SectionsConfigs} from '@fbcnms/magmalte/app/components/layout/Section';

import CellWifiIcon from '@material-ui/icons/CellWifi';
import FEGConfigure from './FEGConfigure';
import FEGGateways from './FEGGateways';
import React from 'react';
import SettingsCellIcon from '@material-ui/icons/SettingsCell';

export function getFEGSections(): SectionsConfigs {
  const sections = [
    {
      path: 'gateways',
      label: 'Gateways',
      icon: <CellWifiIcon />,
      component: FEGGateways,
    },
    {
      path: 'configure',
      label: 'Configure',
      icon: <SettingsCellIcon />,
      component: FEGConfigure,
    },
  ];

  return [
    'gateways', // landing path
    sections,
  ];
}
