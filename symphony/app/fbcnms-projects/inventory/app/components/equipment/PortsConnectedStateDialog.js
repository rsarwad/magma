/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  AddLinkMutationResponse,
  AddLinkMutationVariables,
} from '../../mutations/__generated__/AddLinkMutation.graphql';
import type {Equipment, EquipmentPort} from '../../common/Equipment';
import type {Link} from '../../common/Equipment';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {Property} from '../../common/Property';
import type {
  RemoveLinkMutationResponse,
  RemoveLinkMutationVariables,
} from '../../mutations/__generated__/RemoveLinkMutation.graphql';
import type {WithSnackbarProps} from 'notistack';

import AddLinkMutation from '../../mutations/AddLinkMutation';
import Alert from '@fbcnms/ui/components/Alert/Alert';
import Dialog from '@material-ui/core/Dialog';
import PortsConnectDialog from './PortsConnectDialog';
import React from 'react';
import RemoveLinkMutation from '../../mutations/RemoveLinkMutation';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {toPropertyInput} from '../../common/Property';
import {withSnackbar} from 'notistack';

type Props = {
  workOrderId: ?string,
  equipment: Equipment,
  port: EquipmentPort,
  link: ?Link,
  mode: 'connect' | 'disconnect',
  open: boolean,
  onClose: () => void,
} & WithSnackbarProps;

class PortsConnectedStateDialog extends React.Component<Props> {
  render() {
    const {open, onClose, mode, link} = this.props;
    if (mode === 'disconnect') {
      const linkServices = link?.services ?? [];
      const deleteMsg = (
        <span>
          Are you sure you want to disconnect these ports?
          {linkServices.length > 0 && (
            <span>
              <br />
              The link between them is used by some services and deleting it can
              potentially break them
            </span>
          )}
        </span>
      );
      return (
        <Alert
          cancelLabel="Cancel"
          confirmLabel="Confirm"
          open={open}
          message={deleteMsg}
          onConfirm={() => this.disconnectPorts()}
          onCancel={onClose}
        />
      );
    }

    return (
      <Dialog
        open={this.props.open}
        onClose={this.props.onClose}
        maxWidth={false}
        fullWidth={true}>
        <PortsConnectDialog
          equipment={this.props.equipment}
          port={this.props.port}
          onConnectPorts={this.onConnectPorts}
        />
      </Dialog>
    );
  }

  onConnectPorts = (targetPort: EquipmentPort, properties: Array<Property>) => {
    ServerLogger.info(LogEvents.CONNECT_PORTS_CLICKED);
    const variables: AddLinkMutationVariables = {
      input: {
        sides: [
          {
            equipment: this.props.equipment.id,
            port: this.props.port.definition.id,
          },
          {
            equipment: targetPort.parentEquipment.id,
            port: targetPort.definition.id,
          },
        ],
        workOrder: this.props.workOrderId,
        properties: toPropertyInput(properties),
      },
    };
    const callbacks: MutationCallbacks<AddLinkMutationResponse> = {
      onCompleted: (_, errors) => {
        if (errors && errors[0]) {
          this.props.enqueueSnackbar(errors[0].message, {
            children: key => (
              <SnackbarItem
                id={key}
                message={errors[0].message}
                variant="error"
              />
            ),
          });
        }
        this.props.onClose();
      },
      onError: () => {
        this.props.onClose();
      },
    };
    const updater = store => {
      const {port} = this.props;
      const newLink = store.getRootField('addLink');
      if (newLink == null) {
        return;
      }
      if (port.id.includes('@tmp')) {
        const linkPorts = newLink.getLinkedRecords('ports');
        if (linkPorts == null) {
          return;
        }
        linkPorts.forEach(lp => lp?.setLinkedRecord(newLink, 'link'));
        const newPort = linkPorts.find(
          lp =>
            lp?.getLinkedRecord('definition')?.getDataID() ===
            port.definition.id,
        );
        if (newPort == null) {
          return;
        }
        const equipmentProxy = store.get(this.props.equipment.id);
        if (equipmentProxy == null) {
          return;
        }
        const eqPorts = equipmentProxy.getLinkedRecords('ports') || [];
        equipmentProxy.setLinkedRecords([...eqPorts, newPort], 'ports');
      } else {
        const linkPorts = newLink?.getLinkedRecords('ports');
        linkPorts?.map(lp => lp && lp.setLinkedRecord(newLink, 'link'));
      }
    };
    AddLinkMutation(variables, callbacks, updater);
  };

  disconnectPorts = () => {
    ServerLogger.info(LogEvents.DISCONNECT_PORTS_CLICKED);
    const variables: RemoveLinkMutationVariables = {
      id: this.props.link?.id || '',
      workOrderId: this.props.workOrderId,
    };
    const callbacks: MutationCallbacks<RemoveLinkMutationResponse> = {
      onCompleted: () => {
        this.props.onClose();
      },
      onError: () => {
        this.props.onClose();
      },
    };
    const updater = store => {
      const sourcePortProxy = store.get(this.props.port.id);
      if (sourcePortProxy == null) {
        return;
      }
      if (this.props.workOrderId) {
        const deletedLink = store.getRootField('removeLink');
        deletedLink != null &&
          sourcePortProxy.setLinkedRecord(deletedLink, 'link');
      } else {
        sourcePortProxy.setValue(null, 'connectedPort');
        sourcePortProxy.setValue(null, 'link');
      }
    };

    RemoveLinkMutation(variables, callbacks, updater);
  };
}

export default withSnackbar(PortsConnectedStateDialog);
