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
import LinearProgress from '@material-ui/core/LinearProgress';
import React from 'react';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import grey from '@material-ui/core/colors/grey';
import {MagmaAPIUrls} from '../common/MagmaAPI';

import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = _theme => ({
  input: {
    margin: '10px 0',
    width: '100%',
  },
  divider: {
    margin: '10px 0',
  },
});

type Props = ContextRouter &
  WithAlert &
  WithStyles<typeof styles> & {
    onClose?: () => void,
    gatewayID: string,
    showRestartCommand: boolean,
    showRebootEnodebCommand: boolean,
    showPingCommand: boolean,
    showGenericCommand: boolean,
  };

type State = {
  showRebootCheck: boolean,
  showRestartCheck: boolean,
  enodebSerial: string,
  showRebootEnodebProgress: boolean,
  rebootEnodebResponse: string,
  pingHosts: string,
  pingPackets: string,
  pingResponse: string,
  showPingProgress: boolean,
  genericCommandName: string,
  genericParams: string,
  genericResponse: string,
  showGenericProgress: boolean,
};

function CommandResponse(props) {
  return (
    <pre
      style={{
        backgroundColor: grey[100],
        fontSize: '12px',
        color: grey[900],
      }}>
      {props.showProgressBar && <LinearProgress />}
      <code>{props.response}</code>
    </pre>
  );
}

class GatewayCommandFields extends React.Component<Props, State> {
  state = {
    showRebootCheck: false,
    showRestartCheck: false,
    enodebSerial: '',
    showRebootEnodebProgress: false,
    rebootEnodebResponse: '',
    pingHosts: '',
    pingPackets: '',
    pingResponse: '',
    showPingProgress: false,
    genericCommandName: '',
    genericParams: '{\n}',
    genericResponse: '',
    showGenericProgress: false,
  };

  render() {
    return (
      <>
        <DialogContent>
          <Typography className={this.props.classes.title} variant="subtitle1">
            Reboot
          </Typography>
          <FormField
            label="Reboot Device"
            tooltip="Reboot the Magma gateway server">
            <Button onClick={this.handleRebootGateway} color="primary">
              Reboot
            </Button>
            <Fade in={this.state.showRebootCheck} timeout={500}>
              <Check style={{verticalAlign: 'middle'}} htmlColor="green" />
            </Fade>
          </FormField>
          <div style={this.props.showRestartCommand ? {} : {display: 'none'}}>
            <FormField
              label="Restart Services"
              tooltip="Restart all MagmaD services on this gateway">
              <Button onClick={this.handleRestartServices} color="primary">
                Restart Services
              </Button>
              <Fade in={this.state.showRestartCheck} timeout={500}>
                <Check style={{verticalAlign: 'middle'}} htmlColor="green" />
              </Fade>
            </FormField>
          </div>
          <div
            style={this.props.showRebootEnodebCommand ? {} : {display: 'none'}}>
            <Divider className={this.props.classes.divider} />
            <Typography
              className={this.props.classes.title}
              variant="subtitle1">
              Reboot eNodeB
            </Typography>
            <FormField label="eNodeB Serial ID">
              <Input
                className={this.props.classes.input}
                value={this.state.enodebSerial}
                onChange={this.enodebSerialChanged}
                placeholder="Leave empty to reboot every connected eNodeB"
              />
            </FormField>
            <FormField label="">
              <Button onClick={this.handleRebootEnodeb} color="primary">
                Reboot
              </Button>
            </FormField>
            <FormField label="">
              <CommandResponse
                response={this.state.rebootEnodebResponse}
                showProgressBar={this.state.showRebootEnodebProgress}
              />
            </FormField>
          </div>
          <div style={this.props.showPingCommand ? {} : {display: 'none'}}>
            <Divider className={this.props.classes.divider} />
            <Typography
              className={this.props.classes.title}
              variant="subtitle1">
              Ping
            </Typography>
            <FormField label="Host(s) (one per line)">
              <Input
                className={this.props.classes.input}
                value={this.state.pingHosts}
                onChange={this.pingHostsChanged}
                placeholder="E.g. example.com"
                multiline={true}
              />
            </FormField>
            <FormField label="Packets (default 4)">
              <Input
                className={this.props.classes.input}
                value={this.state.pingPackets}
                onChange={this.pingPacketsChanged}
                placeholder="E.g. 4"
                type="number"
              />
            </FormField>
            <FormField label="">
              <Button onClick={this.handlePing} color="primary">
                Ping
              </Button>
            </FormField>
            <FormField label="">
              <CommandResponse
                response={this.state.pingResponse}
                showProgressBar={this.state.showPingProgress}
              />
            </FormField>
          </div>
          <div style={this.props.showGenericCommand ? {} : {display: 'none'}}>
            <Divider className={this.props.classes.divider} />
            <Typography
              className={this.props.classes.title}
              variant="subtitle1">
              Generic
            </Typography>
            <FormField label="Command">
              <Input
                className={this.props.classes.input}
                value={this.state.genericCommandName}
                onChange={this.genericCommandNameChanged}
                placeholder="Command name"
              />
            </FormField>
            <FormField label="Parameters">
              <Input
                className={this.props.classes.input}
                value={this.state.genericParams}
                onChange={this.genericParamsChanged}
                multiline={true}
                style={{fontFamily: 'monospace', fontSize: '14px'}}
              />
            </FormField>
            <FormField label="">
              <Button onClick={this.handleGeneric} color="primary">
                Execute
              </Button>
            </FormField>
            <FormField label="">
              <CommandResponse
                response={this.state.genericResponse}
                showProgressBar={this.state.showGenericProgress}
              />
            </FormField>
          </div>
        </DialogContent>
        {this.props.onClose && (
          <DialogActions>
            <Button onClick={this.props.onClose} color="primary">
              Close
            </Button>
          </DialogActions>
        )}
      </>
    );
  }

  handleRebootGateway = () => {
    const {match, gatewayID} = this.props;
    const commandName = 'reboot';

    axios
      .post(MagmaAPIUrls.command(match, gatewayID, commandName))
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
    const {match, gatewayID} = this.props;
    const commandName = 'restart_services';

    axios
      .post(MagmaAPIUrls.command(match, gatewayID, commandName), [])
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

  handlePing = () => {
    const {match, gatewayID} = this.props;
    const commandName = 'ping';

    const hosts = this.state.pingHosts.split('\n').filter(host => host);
    const packets = parseInt(this.state.pingPackets);
    const params = {
      hosts,
      packets,
    };

    this.setState({showPingProgress: true});
    axios
      .post(MagmaAPIUrls.command(match, gatewayID, commandName), params)
      .then(resp => {
        this.setState({pingResponse: JSON.stringify(resp.data, null, 2)});
      })
      .catch(error =>
        this.props.alert('Ping failed: ' + error.response.data.message),
      )
      .finally(() => this.setState({showPingProgress: false}));
  };

  handleRebootEnodeb = () => {
    const {match, gatewayID} = this.props;
    const commandName = 'generic';

    const enodebSerial = this.state.enodebSerial;
    const params =
      enodebSerial.length > 0
        ? {
            command: 'reboot_enodeb',
            params: {shell_params: [enodebSerial]},
          }
        : {
            command: 'reboot_all_enodeb',
            params: {},
          };

    this.setState({showRebootEnodebProgress: true});
    axios
      .post(MagmaAPIUrls.command(match, gatewayID, commandName), params)
      .then(resp => {
        this.setState({
          rebootEnodebResponse: JSON.stringify(resp.data, null, 2),
        });
      })
      .catch(error =>
        this.props.alert(
          'Reboot eNodeB failed: ' + error.response.data.message,
        ),
      )
      .finally(() => this.setState({showRebootEnodebProgress: false}));
  };

  handleGeneric = () => {
    const {match, gatewayID} = this.props;
    const commandName = 'generic';

    const genericCommandName = this.state.genericCommandName;
    let genericCommandParams = {};
    try {
      genericCommandParams = JSON.parse(this.state.genericParams);
    } catch (e) {
      this.props.alert('Error parsing params: ' + e);
      return;
    }
    const params = {
      command: genericCommandName,
      params: genericCommandParams,
    };

    this.setState({showGenericProgress: true});
    axios
      .post(MagmaAPIUrls.command(match, gatewayID, commandName), params)
      .then(resp => {
        this.setState({genericResponse: JSON.stringify(resp.data, null, 2)});
      })
      .catch(error =>
        this.props.alert(
          'Generic command failed: ' + error.response.data.message,
        ),
      )
      .finally(() => this.setState({showGenericProgress: false}));
  };

  enodebSerialChanged = ({target}) =>
    this.setState({enodebSerial: target.value});

  pingHostsChanged = ({target}) => this.setState({pingHosts: target.value});
  pingPacketsChanged = ({target}) => this.setState({pingPackets: target.value});

  genericCommandNameChanged = ({target}) =>
    this.setState({genericCommandName: target.value});
  genericParamsChanged = ({target}) => {
    this.setState({genericParams: target.value});
  };
}

export default withStyles(styles)(withRouter(withAlert(GatewayCommandFields)));
