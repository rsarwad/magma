/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ContextRouter} from 'react-router-dom';
import type {WifiGateway} from './WifiUtils';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import AddIcon from '@material-ui/icons/Add';
import Button from '@material-ui/core/Button';
import ChevronRight from '@material-ui/icons/ChevronRight';
import ClipboardLink from '@fbcnms/ui/components/ClipboardLink';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import ExpandMore from '@material-ui/icons/ExpandMore';
import IconButton from '@material-ui/core/IconButton';
import InfoIcon from '@material-ui/icons/Info';
import LinkIcon from '@material-ui/icons/Link';
import MagmaV1API from '@fbcnms/magmalte/app/common/MagmaV1API';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import React from 'react';
import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';
import Tooltip from '@material-ui/core/Tooltip';
import axios from 'axios';
import url from 'url';
import {groupBy} from 'lodash';

import WifiDeviceDetails, {InfoRow} from './WifiDeviceDetails';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {GatewayStatus} from '@fbcnms/magmalte/app/components/GatewayUtils';
import {MagmaAPIUrls} from '@fbcnms/magmalte/app/common/MagmaAPI';
import {meshesURL} from './WifiUtils';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  actionsCell: {
    textAlign: 'right',
  },
  gatewayCell: {
    paddingBottom: '15px',
    paddingLeft: '75px',
    paddingRight: '15px',
    paddingTop: '15px',
  },
  deviceWarning: {
    color: 'red',
    paddingLeft: 40,
  },
  iconButton: {
    color: theme.palette.secondary.light,
    padding: '5px',
  },
  meshButton: {
    margin: 0,
    textTransform: 'none',
  },
  meshCell: {
    padding: '5px',
  },
  meshID: {
    color: theme.palette.primary.dark,
    fontWeight: 'bolder',
  },
  meshIconButton: {
    color: theme.palette.primary.dark,
    padding: '5px',
  },
  tableCell: {
    padding: '15px',
  },
  tableRow: {
    height: 'auto',
    whiteSpace: 'nowrap',
    verticalAlign: 'top',
  },
});

type Props = ContextRouter &
  WithStyles<typeof styles> &
  WithAlert & {
    enableDeviceEditing?: boolean,
    meshID: string,
    gateways: WifiGateway[],
    onDeleteMesh: string => void,
    onDeleteDevice: WifiGateway => void,
  };

const EXPANDED_STATE_TYPES = {
  none: 0,
  device: 1,
  neighbors: 2,
  fullDump: 3,
  configs: 4,
};

const MESH_ID_PARAM = 'meshID';
const DEVICE_ID_PARAM = 'deviceID';
const EXPANDED_STATE_PARAM = 'expandedState';

type State = {
  expanded: boolean,
  expandedGateways: {[key: string]: 0 | 1 | 2 | 3 | 4},
};

class WifiMeshRow extends React.Component<Props, State> {
  constructor(props) {
    super(props);
    const queryParams = new URLSearchParams(this.props.location.search);
    const initialState = {
      expanded: false,
      expandedGateways: {},
    };
    if (queryParams.get(MESH_ID_PARAM) === this.props.meshID) {
      initialState.expanded = true;
      const deviceID = queryParams.get(DEVICE_ID_PARAM);
      if (deviceID != null) {
        const expandedState = parseInt(queryParams.get(EXPANDED_STATE_PARAM));
        initialState.expandedGateways = {[deviceID]: expandedState};
      }
    }
    this.state = initialState;
  }

  handleToggleAllDevices = () => {
    const {gateways} = this.props;

    // determine old max state by getting max() of all the states
    const maxState = gateways
      .map(gateway => this.state.expandedGateways[gateway.id])
      .reduce((max, state) => (state ? Math.max(max, state) : max), 0);

    // calculate next state
    const nextState = (maxState + 1) % Object.keys(EXPANDED_STATE_TYPES).length;

    // assign same next state to all gateways
    if (nextState === 0) {
      // no need to set any gateways for 0/unexpanded state
      this.setState({expandedGateways: {}});
    } else {
      const expandedGateways = gateways
        .map(gateway => gateway.id)
        .reduce((expandedGateways, id) => {
          expandedGateways[id] = nextState;
          return expandedGateways;
        }, {});
      this.setState({expandedGateways});
    }
  };

  render() {
    const {meshID, gateways, classes} = this.props;

    let gatewayRows;
    let gatewayVersions;
    if (this.state.expanded) {
      gatewayRows = gateways.map(gateway => (
        <TableRow className={this.props.classes.tableRow} key={gateway.id}>
          <TableCell className={classes.gatewayCell}>{gateway.info}</TableCell>
          <TableCell className={classes.tableCell}>
            {status}
            <GatewayStatus isGrey={!gateway.status} isActive={!!gateway.up} />
            <Tooltip
              title="Click to toggle device info"
              enterDelay={400}
              placement={'right'}>
              <span onClick={() => this.expandGateway(gateway.id)}>
                {gateway.id}
              </span>
            </Tooltip>
            {gateway.coordinates.includes(NaN) && (
              <span className={classes.deviceWarning}>
                {' '}
                Please configure Lat/Lng
              </span>
            )}
            {gateway.status &&
              gateway.status.meta &&
              gateway.status.meta['validation_status'] !== 'passed' && (
                <span className={classes.deviceWarning}>
                  {' '}
                  Please check image validation status
                </span>
              )}

            {!!this.state.expandedGateways[gateway.id] && gateway.status && (
              <WifiDeviceDetails
                device={gateway}
                hideHeader={true}
                showConfigs={
                  this.state.expandedGateways[gateway.id] ===
                  EXPANDED_STATE_TYPES.configs
                }
                showDevice={
                  this.state.expandedGateways[gateway.id] ===
                  EXPANDED_STATE_TYPES.device
                }
                showNeighbors={
                  this.state.expandedGateways[gateway.id] ===
                  EXPANDED_STATE_TYPES.neighbors
                }
                showFullDump={
                  this.state.expandedGateways[gateway.id] ===
                  EXPANDED_STATE_TYPES.fullDump
                }
              />
            )}
          </TableCell>
          <TableCell className={classes.actionsCell}>
            <ClipboardLink title="Copy link to this device">
              {({copyString}) => (
                <IconButton
                  className={classes.iconButton}
                  onClick={() =>
                    copyString(this.buildLinkURL(meshID, gateway.id))
                  }>
                  <LinkIcon />
                </IconButton>
              )}
            </ClipboardLink>
            <Tooltip title="Click to toggle device info" enterDelay={400}>
              <IconButton
                className={classes.iconButton}
                onClick={() => this.expandGateway(gateway.id)}>
                <InfoIcon />
              </IconButton>
            </Tooltip>
            {this.props.enableDeviceEditing && (
              <NestedRouteLink to={`/${meshID}/edit_device/${gateway.id}`}>
                <IconButton className={classes.iconButton}>
                  <EditIcon />
                </IconButton>
              </NestedRouteLink>
            )}
            <IconButton
              className={classes.iconButton}
              onClick={() => this.showDeviceDeleteDialog(gateway)}>
              <DeleteIcon />
            </IconButton>
          </TableCell>
        </TableRow>
      ));

      // construct version list per mesh
      const versionGroups = groupBy(gateways, device => {
        if (device.versionParsed) {
          if (device.versionParsed.fbpkg !== 'none') {
            return device.versionParsed.fbpkg;
          } else {
            return device.versionParsed.hash;
          }
        }
        return device.version || 'UNKNOWN';
      });

      gatewayVersions = Object.keys(versionGroups).map(version => (
        <div key={version}>
          <Tooltip
            title={`${versionGroups[version].length} device(s) with ${versionGroups[version][0].version}`}
            enterDelay={100}
            key={version}>
            <span>
              {version}: <b>{versionGroups[version].length}</b>
            </span>
          </Tooltip>
        </div>
      ));
    }

    return (
      <>
        <TableRow className={this.props.classes.tableRow}>
          <TableCell className={this.props.classes.meshCell}>
            <IconButton
              className={classes.meshIconButton}
              onClick={gateways.length == 0 ? null : this.onToggleExpand}>
              {this.state.expanded ? <ExpandMore /> : <ChevronRight />}
            </IconButton>
            <span className={classes.meshID}>{meshID}</span>
          </TableCell>
          <TableCell className={this.props.classes.meshCell}>
            {gateways.length > 0 && (
              <>
                <InfoRow
                  label="Up"
                  data={`${gateways.filter(gateway => gateway.up).length} of ${
                    gateways.length
                  }`}
                />
                {this.state.expanded && (
                  <>
                    <Tooltip
                      title="Click to toggle device info"
                      enterDelay={400}
                      placement={'right'}>
                      <Button
                        size="small"
                        className={classes.meshButton}
                        onClick={this.handleToggleAllDevices}>
                        toggle info
                      </Button>
                    </Tooltip>
                    {gatewayVersions}
                  </>
                )}
              </>
            )}
          </TableCell>
          <TableCell className={this.props.classes.actionsCell}>
            <ClipboardLink title="Copy link to this mesh">
              {({copyString}) => (
                <IconButton
                  className={classes.iconButton}
                  onClick={() => copyString(this.buildLinkURL(meshID))}>
                  <LinkIcon />
                </IconButton>
              )}
            </ClipboardLink>
            <NestedRouteLink to={`/add_device/${meshID}`}>
              <IconButton className={classes.meshIconButton}>
                <AddIcon />
              </IconButton>
            </NestedRouteLink>
            <NestedRouteLink to={`/edit_mesh/${meshID}`}>
              <IconButton className={classes.meshIconButton}>
                <EditIcon />
              </IconButton>
            </NestedRouteLink>
            <IconButton
              className={classes.meshIconButton}
              onClick={this.showMeshDeleteDialog}>
              <DeleteIcon />
            </IconButton>
          </TableCell>
        </TableRow>
        {gatewayRows}
      </>
    );
  }

  onToggleExpand = () => this.setState({expanded: !this.state.expanded});
  expandGateway = id => {
    const expandedGateways = {
      ...this.state.expandedGateways,
      [id]:
        ((this.state.expandedGateways[id] | 0) + 1) %
        Object.keys(EXPANDED_STATE_TYPES).length,
    };

    this.setState({expandedGateways});
  };

  showMeshDeleteDialog = () => {
    this.props
      .confirm(
        `Are you sure you want to delete mesh "${this.props.meshID}" and all its devices (count: ${this.props.gateways.length})?`,
      )
      .then(async confirmed => {
        if (!confirmed) {
          return;
        }

        const requests = this.props.gateways.map(device =>
          MagmaV1API.deleteNetworksByNetworkIdGatewaysByGatewayId({
            networkId: nullthrows(this.props.match.params.networkId),
            gatewayId: device.id,
          }),
        );
        const requestsc = this.props.gateways.map(device =>
          axios.delete(
            MagmaAPIUrls.gatewayConfigsForType(
              this.props.match,
              device.id,
              'wifi',
            ),
          ),
        );

        // delete all devices and wifi configs
        await axios.all([...requests, ...requestsc]);

        // delete mesh
        await axios.delete(
          meshesURL(this.props.match) + '/' + this.props.meshID,
        );

        this.props.onDeleteMesh(this.props.meshID);
      });
  };

  showDeviceDeleteDialog = (device: WifiGateway) => {
    this.props
      .confirm(`Are you sure you want to delete "${device.id}"?`)
      .then(async confirmed => {
        if (!confirmed) {
          return;
        }

        // delete all parts
        await axios.all([
          MagmaV1API.deleteNetworksByNetworkIdGatewaysByGatewayId({
            networkId: nullthrows(this.props.match.params.networkId),
            gatewayId: device.id,
          }),
          axios.delete(
            MagmaAPIUrls.gatewayConfigsForType(
              this.props.match,
              device.id,
              'wifi',
            ),
          ),
        ]);

        this.props.onDeleteDevice(device);
      });
  };

  buildLinkURL = (meshID: string, deviceID: ?string = null): string => {
    const query: {[string]: string | number} = {[MESH_ID_PARAM]: meshID};
    if (deviceID) {
      query[DEVICE_ID_PARAM] = deviceID;
      query[EXPANDED_STATE_PARAM] = this.state.expandedGateways[deviceID] ?? 1;
    }
    const {protocol, host, pathname} = window.location;
    return url.format({protocol, host, pathname, query});
  };
}

export default withStyles(styles)(withRouter(withAlert(WifiMeshRow)));
