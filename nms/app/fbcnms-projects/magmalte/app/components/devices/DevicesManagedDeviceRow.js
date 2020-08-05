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
 * @flow
 * @format
 */

import type {FullDevice} from './DevicesUtils';

import Alert from '@fbcnms/ui/components/Alert/Alert';
import DeleteIcon from '@material-ui/icons/Delete';
import DevicesState from './DevicesState';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import React from 'react';
import ShowChartIcon from '@material-ui/icons/ShowChart';
import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';

import nullthrows from '@fbcnms/util/nullthrows';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

const useStyles = makeStyles(_ => ({
  actionsCell: {
    textAlign: 'right',
  },
  deviceCell: {
    paddingBottom: '15px',
    paddingLeft: '50px',
    paddingRight: '15px',
    paddingTop: '15px',
  },
  iconButton: {
    // color: theme.palette.secondary.light,
    padding: '5px',
  },
  subrowCell: {
    paddingBottom: '15px',
    paddingLeft: '75px',
    paddingRight: '15px',
    paddingTop: '15px',
  },
  tableCell: {
    padding: '15px',
  },
  tableRow: {
    height: 'auto',
    whiteSpace: 'nowrap',
    verticalAlign: 'top',
  },
}));

type Props = {
  enableDeviceEditing?: boolean,
  device: FullDevice,
  deviceID: string,
  onDeleteDevice: string => void,
};

export default function DevicesManagedDeviceRow(props: Props) {
  const classes = useStyles();
  const {match} = useRouter();

  const {device} = props;

  const [confirmDialog, setConfirmDialog] = useState(false);

  const deleteDevice = () => {
    MagmaV1API.deleteSymphonyByNetworkIdDevicesByDeviceId({
      networkId: nullthrows(match.params.networkId),
      deviceId: props.deviceID,
    }).then(() => {
      setConfirmDialog(false);
      props.onDeleteDevice(props.deviceID);
    });
  };

  return (
    <>
      <Alert
        open={confirmDialog}
        message={`Are you sure you want to delete device "${props.deviceID}"?`}
        confirmLabel={'Yes'}
        cancelLabel={'No'}
        onConfirm={() => deleteDevice()}
        onCancel={() => setConfirmDialog(false)}
      />
      <TableRow className={classes.tableRow}>
        <TableCell className={classes.subrowCell}>{props.deviceID}</TableCell>
        <TableCell>{device && <DevicesState device={device} />}</TableCell>
        <TableCell>{device.managingAgentId || '<none>'}</TableCell>

        <TableCell className={classes.actionsCell}>
          <NestedRouteLink to={`/metrics/${props.deviceID}`}>
            <IconButton className={classes.iconButton}>
              <ShowChartIcon />
            </IconButton>
          </NestedRouteLink>
          {props.enableDeviceEditing && (
            <NestedRouteLink to={`/edit_device/${props.deviceID}`}>
              <IconButton className={classes.iconButton}>
                <EditIcon />
              </IconButton>
            </NestedRouteLink>
          )}
          <IconButton
            className={classes.iconButton}
            onClick={() => setConfirmDialog(true)}>
            <DeleteIcon />
          </IconButton>
        </TableCell>
      </TableRow>
    </>
  );
}
