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
import type {Gateway} from './GatewayUtils';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import Button from '@material-ui/core/Button';
import Check from '@material-ui/icons/Check';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import Divider from '@material-ui/core/Divider';
import Fade from '@material-ui/core/Fade';
import FormField from './FormField';
import Input from '@material-ui/core/Input';
import React from 'react';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import moment from 'moment';
import {GatewayStatus} from './GatewayUtils';
import {MagmaAPIUrls, updateGatewayName} from '../common/MagmaAPI';

import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

type Props = ContextRouter &
  WithAlert &
  WithStyles<typeof styles> & {
    onClose: () => void,
    onSave: (gatewayID: string) => void,
    gateway: Gateway,
  };

type State = {
  name: string,
  showRebootCheck: boolean,
  showRestartCheck: boolean,
};

const styles = {
  input: {
    width: '100%',
  },
  divider: {
    margin: '10px 0',
  },
};

class GatewaySummaryFields extends React.Component<Props, State> {
  state = {
    name: this.props.gateway.name,
    showRebootCheck: false,
    showRestartCheck: false,
  };

  render() {
    const {gateway} = this.props;
    return (
      <>
        <DialogContent>
          <FormField label="Name">
            <Input
              className={this.props.classes.input}
              value={this.state.name}
              onChange={this.nameChanged}
              placeholder="E.g. Gateway 1234"
            />
          </FormField>
          <FormField label="Hardware UUID">{gateway.hardware_id}</FormField>
          <FormField label="Gateway ID">{gateway.logicalID}</FormField>
          <FormField label="Last Checkin">
            {moment(parseInt(gateway.lastCheckin, 10)).fromNow()}
          </FormField>
          <FormField label="Version">{gateway.version}</FormField>
          <FormField label="VPN IP">{gateway.vpnIP}</FormField>
          <FormField label="RF Transmitter">
            <GatewayStatus
              isGrey={false}
              isActive={gateway.enodebRFTXEnabled}
            />
            {gateway.enodebRFTXEnabled ? '' : 'Not '}
            Allowed
            {'  '}
            <GatewayStatus
              isGrey={gateway.isBackhaulDown}
              isActive={gateway.enodebConnected && gateway.enodebRFTXOn}
            />
            {gateway.enodebRFTXOn ? '' : 'Not '}
            Connected
          </FormField>
          <FormField label="GPS synchronized">
            <GatewayStatus
              isGrey={gateway.isBackhaulDown}
              isActive={gateway.enodebConnected && gateway.gpsConnected}
            />
            {gateway.gpsConnected ? '' : 'Not '}
            Synced
          </FormField>
          <FormField label="MME">
            <GatewayStatus
              isGrey={gateway.isBackhaulDown}
              isActive={gateway.enodebConnected && gateway.mmeConnected}
            />
            {gateway.mmeConnected ? '' : 'Not '}
            Connected
          </FormField>
          <Divider className={this.props.classes.divider} />
          <Typography variant="subtitle1">Commands</Typography>
          <FormField label="Reboot Gateway">
            <Button onClick={this.handleRebootGateway} color="primary">
              Reboot
            </Button>
            <Fade in={this.state.showRebootCheck} timeout={500}>
              <Check style={{verticalAlign: 'middle'}} htmlColor="green" />
            </Fade>
          </FormField>
          <FormField label="">
            <Button onClick={this.handleRestartServices} color="primary">
              Restart services
            </Button>
            <Fade in={this.state.showRestartCheck} timeout={500}>
              <Check style={{verticalAlign: 'middle'}} htmlColor="green" />
            </Fade>
          </FormField>
        </DialogContent>
        <DialogActions>
          <Button onClick={this.props.onClose} color="primary">
            Cancel
          </Button>
          <Button onClick={this.onSave} color="primary" variant="contained">
            Save
          </Button>
        </DialogActions>
      </>
    );
  }

  onSave = () => {
    const {match, gateway} = this.props;
    const id = gateway.logicalID;

    updateGatewayName(id, this.state.name, match)
      .then(() => this.props.onSave(id))
      .catch(error => this.props.alert(error.response.data.message));
  };

  nameChanged = ({target}: SyntheticInputEvent<*>) =>
    this.setState({name: target.value});

  handleRebootGateway = () => {
    const {match, gateway} = this.props;
    const id = gateway.logicalID;
    const commandName = 'reboot';

    axios
      .post(MagmaAPIUrls.command(match, id, commandName))
      .then(_resp => {
        this.props.alert('Successfully initiated reboot');
        this.setState({showRebootCheck: true}, () => {
          setTimeout(() => this.setState({showRebootCheck: false}), 5000);
        });
      })
      .catch(error =>
        this.props.alert('Reboot failed: ' + error.response.data.message),
      );
  };

  handleRestartServices = () => {
    const {match, gateway} = this.props;
    const id = gateway.logicalID;
    const commandName = 'restart_services';

    axios
      .post(MagmaAPIUrls.command(match, id, commandName), [])
      .then(_resp => {
        this.props.alert('Successfully initiated service restart');
        this.setState({showRestartCheck: true}, () => {
          setTimeout(() => this.setState({showRestartCheck: false}), 5000);
        });
      })
      .catch(error =>
        this.props.alert(
          'Restart services failed: ' + error.response.data.message,
        ),
      );
  };
}

export default withStyles(styles)(withRouter(withAlert(GatewaySummaryFields)));
