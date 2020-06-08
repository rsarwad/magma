/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {Equipment, EquipmentPosition} from '../../common/Equipment';
import type {EquipmentType} from '../../common/EquipmentType';
import type {TabProps} from '@fbcnms/ui/components/design-system/Tabs/TabsBar';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import EquipmentTypesList from '../EquipmentTypesList';
import LocationEquipments from '../location/LocationEquipments';
import MoveEquipmentToPositionMutation from '../../mutations/MoveEquipmentToPositionMutation';
import TabsBar from '@fbcnms/ui/components/design-system/Tabs/TabsBar';
import fbt from 'fbt';
import usePermissions from '../admin/userManagement/utils/usePermissions';
import {createFragmentContainer, graphql} from 'react-relay';
import {last} from 'lodash';

import nullthrows from '@fbcnms/util/nullthrows';
import {useCallback} from 'react';
import {useMemo, useState} from 'react';

type Props = {
  open: boolean,
  onClose: () => void,
  onEquipmentTypeSelected: (equipmentType: EquipmentType) => void,
  parentEquipment: Equipment,
  position: EquipmentPosition,
};

type ViewTab = {|
  id: string,
  tab: TabProps,
  view: React.Node,
|};

const AddToEquipmentDialog = (props: Props) => {
  const locations = props.parentEquipment.locationHierarchy;
  const parentLocation = nullthrows(last(locations));
  const [selectedEquipmentType, setSelectedEquipmentType] = useState(null);
  const [selectedEquipment, setSelectedEquipment] = useState(null);

  const permissions = usePermissions();
  const missingAddNewPermissionsMessage = useMemo(
    () =>
      permissions.check({
        entity: 'equipment',
        action: 'create',
      }),
    [permissions],
  );
  const allowAddingNew = !missingAddNewPermissionsMessage;

  const addExistingView = useMemo(
    () => (
      <LocationEquipments
        locationId={parentLocation.id}
        onSelect={equipment => setSelectedEquipment(equipment)}
      />
    ),
    [parentLocation.id],
  );

  const tabBars: Array<ViewTab> = useMemo(
    () => [
      {
        id: 'new',
        tab: {
          label: fbt('NEW EQUIPMENT', ''),
        },
        view: (
          <EquipmentTypesList
            onSelect={type => setSelectedEquipmentType(type)}
          />
        ),
      },
      {
        id: 'existing',
        tab: {
          label: fbt('EXISTING EQUIPMENT', ''),
        },
        view: addExistingView,
      },
    ],
    [addExistingView],
  );
  const [activeTabBar, setActiveTabBar] = useState<number>(
    allowAddingNew ? 0 : 1,
  );

  const clearState = useCallback(() => {
    setSelectedEquipment(null);
    setSelectedEquipmentType(null);
  }, []);

  const onClose = useCallback(() => {
    clearState();
    if (props.onClose) {
      props.onClose();
    }
  }, [clearState, props]);

  const updateActiveTabIndex = useCallback(
    newActiveTabIndex => {
      clearState();
      setActiveTabBar(newActiveTabIndex);
    },
    [clearState],
  );

  return (
    <Dialog maxWidth="sm" open={props.open} onClose={onClose}>
      <DialogTitle>
        {allowAddingNew ? (
          <TabsBar
            spread={true}
            tabs={tabBars.map(tabBar => tabBar.tab)}
            activeTabIndex={activeTabBar}
            onChange={updateActiveTabIndex}
          />
        ) : (
          <fbt desc="">Select an existing equipment</fbt>
        )}
      </DialogTitle>
      <DialogContent>{tabBars[activeTabBar].view}</DialogContent>
      <DialogActions>
        <Button onClick={onClose} skin="regular">
          Cancel
        </Button>
        <Button
          disabled={
            (tabBars[activeTabBar].id === 'new' &&
              selectedEquipmentType === null) ||
            (tabBars[activeTabBar].id === 'existing' &&
              selectedEquipment === null)
          }
          onClick={() => {
            if (tabBars[activeTabBar].id === 'new') {
              props.onEquipmentTypeSelected(nullthrows(selectedEquipmentType));
            } else {
              MoveEquipmentToPositionMutation(
                {
                  parent_equipment_id: props.parentEquipment.id,
                  position_definition_id: props.position.definition.id,
                  equipment_id: nullthrows(selectedEquipment).id,
                },
                {
                  onCompleted: onClose,
                  onError: () => {},
                },
              );
            }
          }}>
          OK
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default createFragmentContainer(AddToEquipmentDialog, {
  parentEquipment: graphql`
    fragment AddToEquipmentDialog_parentEquipment on Equipment {
      id
      locationHierarchy {
        id
      }
    }
  `,
});
