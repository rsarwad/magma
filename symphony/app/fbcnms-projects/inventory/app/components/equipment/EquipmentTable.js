/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {EquipmentTable_equipments} from './__generated__/EquipmentTable_equipments.graphql';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {
  RemoveEquipmentMutationResponse,
  RemoveEquipmentMutationVariables,
} from '../../mutations/__generated__/RemoveEquipmentMutation.graphql';
import type {TableRowDataType} from '@fbcnms/ui/components/design-system/Table/Table';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';

import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@fbcnms/ui/components/design-system/Button';
import CommonStrings from '@fbcnms/strings/Strings';
import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';
import FormActionWithPermissions from '../../common/FormActionWithPermissions';
import IconButton from '@fbcnms/ui/components/design-system/IconButton';
import React, {useCallback, useContext, useMemo} from 'react';
import RemoveEquipmentMutation from '../../mutations/RemoveEquipmentMutation';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import Table from '@fbcnms/ui/components/design-system/Table/Table';
import fbt from 'fbt';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {DeleteIcon} from '@fbcnms/ui/components/design-system/Icons';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {capitalize} from '@fbcnms/util/strings';
import {createFragmentContainer, graphql} from 'react-relay';
import {lowerCase} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {withSnackbar} from 'notistack';

const useStyles = makeStyles(() => ({
  root: {
    width: '100%',
    marginTop: '24px',
    overflowX: 'auto',
  },
  table: {
    minWidth: 70,
    marginBottom: '12px',
  },
  addButton: {
    paddingLeft: '16px',
    paddingRight: '16px',
  },
  futureState: {
    textTransform: 'capitalize',
    maxWidth: '50px',
  },
  iconColumn: {
    width: '36px',
  },
}));

type Props = $ReadOnly<{|
  equipments: EquipmentTable_equipments,
  selectedWorkOrderId: ?string,
  onEquipmentSelected: RelayEquipment => void,
  onWorkOrderSelected: (workOrderId: string) => void,
  ...WithSnackbarProps,
  ...WithAlert,
|}>;

export type RelayEquipment = $ElementType<EquipmentTable_equipments, number>;

const getEquipmentStatus = row =>
  row.equipment.futureState
    ? `${capitalize(lowerCase(row.equipment.workOrder?.status))} ${lowerCase(
        row.equipment.futureState,
      )}`
    : '';
const getIsEquipmentDeviceActive = (eq: TableRowDataType<RelayEquipment>) =>
  eq.device?.up;
const getEquipmentType = row => row.equipment.equipmentType.name || '';
const getEquipmentName = row => row.equipment.name || '';

const handleDelete = (props: Props) => (
  equipment: TableRowDataType<RelayEquipment>,
) => {
  const {alert, confirm, enqueueSnackbar, selectedWorkOrderId} = props;
  ServerLogger.info(LogEvents.DELETE_EQUIPMENT_CLICKED);
  confirm({
    title: <fbt desc="">Delete Equipment?</fbt>,
    message: (
      <fbt desc="">
        By removing{' '}
        <fbt:param name="equipment name">{equipment.name}</fbt:param> from this
        location, all information related to this equipment, like links and
        sub-positions, will be deleted.
      </fbt>
    ),
    checkboxLabel: <fbt desc="">I understand</fbt>,
    cancelLabel: CommonStrings.common.cancelButton,
    confirmLabel: CommonStrings.common.deleteButton,
    skin: 'red',
  }).then(confirmed => {
    if (confirmed) {
      const variables: RemoveEquipmentMutationVariables = {
        id: equipment.id,
        work_order_id: selectedWorkOrderId,
      };

      const cbs: MutationCallbacks<RemoveEquipmentMutationResponse> = {
        onCompleted: (_, errors) => {
          if (errors && errors[0]) {
            enqueueSnackbar(errors[0].message, {
              children: key => (
                <SnackbarItem
                  id={key}
                  message={errors[0].message}
                  variant="error"
                />
              ),
            });
          }
        },
        onError: e => {
          alert(e);
        },
      };

      RemoveEquipmentMutation(variables, cbs, store => {
        if (!selectedWorkOrderId) {
          store.delete(equipment.id);
        }
      });
    }
  });
};

const EquipmentTable = (props: Props) => {
  const {equipments, onEquipmentSelected, onWorkOrderSelected} = props;
  const classes = useStyles();
  const {isFeatureEnabled} = useContext(AppContext);

  const onDelete = useCallback(handleDelete(props), [props]);

  const data: $ReadOnlyArray<
    TableRowDataType<{|equipment: RelayEquipment|}>,
  > = useMemo(
    () =>
      equipments == null
        ? []
        : equipments.filter(Boolean).map(equipment => ({
            key: equipment.id,
            equipment,
          })),
    [equipments],
  );

  const equipmetStatusEnabled = isFeatureEnabled('planned_equipment');
  const equipmentLiveStatusEnabled = isFeatureEnabled('equipment_live_status');

  const columns = useMemo(() => {
    const colsToReturn = [
      {
        key: 'name',
        title: <fbt desc="">Name</fbt>,
        getSortingValue: getEquipmentName,
        render: row => (
          <Button
            variant="text"
            useEllipsis={true}
            onClick={() => onEquipmentSelected(row.equipment)}>
            {equipmentLiveStatusEnabled ? (
              <DeviceStatusCircle
                isGrey={!getIsEquipmentDeviceActive(row.equipment)}
                isActive={!!getIsEquipmentDeviceActive(row.equipment)}
              />
            ) : null}
            {getEquipmentName(row)}
          </Button>
        ),
      },
      {
        key: 'type',
        title: <fbt desc="">Type</fbt>,
        getSortingValue: getEquipmentType,
        render: getEquipmentType,
      },
    ];
    if (equipmetStatusEnabled) {
      colsToReturn.push({
        key: 'status',
        title: <fbt desc="">Status</fbt>,
        getSortingValue: getEquipmentStatus,
        render: row => (
          <Button
            variant="text"
            useEllipsis={true}
            onClick={() =>
              onWorkOrderSelected(nullthrows(row?.equipment.workOrder?.id))
            }>
            {getEquipmentStatus(row)}
          </Button>
        ),
      });
    }
    colsToReturn.push({
      key: 'delete_action',
      title: null,
      titleClassName: classes.iconColumn,
      className: classes.iconColumn,
      render: row => (
        <FormActionWithPermissions
          permissions={{entity: 'equipment', action: 'delete'}}>
          <IconButton
            icon={DeleteIcon}
            onClick={() => onDelete(row.equipment)}
          />
        </FormActionWithPermissions>
      ),
    });
    return colsToReturn;
  }, [
    classes.iconColumn,
    equipmentLiveStatusEnabled,
    equipmetStatusEnabled,
    onDelete,
    onEquipmentSelected,
    onWorkOrderSelected,
  ]);

  if (data.length === 0) {
    return null;
  }

  return (
    <Table
      variant="embedded"
      dataRowsSeparator="border"
      className={classes.table}
      columns={columns}
      data={data}
    />
  );
};

export default withAlert(
  withSnackbar(
    createFragmentContainer(EquipmentTable, {
      equipments: graphql`
        fragment EquipmentTable_equipments on Equipment @relay(plural: true) {
          id
          name
          futureState
          equipmentType {
            id
            name
          }
          workOrder {
            id
            status
          }
          device {
            up
          }
          services {
            id
          }
        }
      `,
    }),
  ),
);
