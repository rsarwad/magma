/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {AddEditWorkOrderTypeCard_workOrderType} from '../configure/__generated__/AddEditWorkOrderTypeCard_workOrderType.graphql';
import type {
  AddWorkOrderCardTypeQuery,
  AddWorkOrderCardTypeQueryResponse,
} from './__generated__/AddWorkOrderCardTypeQuery.graphql';
import type {
  AddWorkOrderMutationResponse,
  AddWorkOrderMutationVariables,
} from '../../mutations/__generated__/AddWorkOrderMutation.graphql';
import type {ChecklistCategoriesMutateStateActionType} from '../checklist/ChecklistCategoriesMutateAction';
import type {ChecklistCategoriesStateType} from '../checklist/ChecklistCategoriesMutateState';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {WorkOrder} from '../../common/WorkOrder';

import AddWorkOrderMutation from '../../mutations/AddWorkOrderMutation';
import Breadcrumbs from '@fbcnms/ui/components/Breadcrumbs';
import CheckListCategoryExpandingPanel from '../checklist/checkListCategory/CheckListCategoryExpandingPanel';
import ChecklistCategoriesMutateDispatchContext from '../checklist/ChecklistCategoriesMutateDispatchContext';
import ExpandingPanel from '@fbcnms/ui/components/ExpandingPanel';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import FormSaveCancelPanel from '@fbcnms/ui/components/design-system/Form/FormSaveCancelPanel';
import Grid from '@material-ui/core/Grid';
import LocationTypeahead from '../typeahead/LocationTypeahead';
import NameDescriptionSection from '@fbcnms/ui/components/NameDescriptionSection';
import ProjectTypeahead from '../typeahead/ProjectTypeahead';
import PropertyValueInput from '../form/PropertyValueInput';
import React, {useCallback, useReducer, useState} from 'react';
import Select from '@fbcnms/ui/components/design-system/Select/Select';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import TextField from '@material-ui/core/TextField';
import UserTypeahead from '../typeahead/UserTypeahead';
import nullthrows from '@fbcnms/util/nullthrows';
import {FormContextProvider} from '../../common/FormContext';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {convertChecklistCategoriesStateToInput} from '../checklist/ChecklistUtils';
import {generateTempId, getGraphError} from '../../common/EntUtils';
import {
  getInitialPropertyFromType,
  toMutablePropertyType,
} from '../../common/PropertyType';
import {
  getInitialStateFromChecklistDefinitions,
  reducer,
} from '../checklist/ChecklistCategoriesMutateReducer';
import {graphql} from 'relay-runtime';
import {makeStyles} from '@material-ui/styles';
import {priorityValues, statusValues} from '../../common/WorkOrder';
import {sortPropertiesByIndex, toPropertyInput} from '../../common/Property';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useHistory, useRouteMatch} from 'react-router';
import {useLazyLoadQuery} from 'react-relay/hooks';

const useStyles = makeStyles(theme => ({
  root: {
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    padding: '40px 32px',
  },
  contentRoot: {
    display: 'flex',
    flexDirection: 'column',
    position: 'relative',
    flexGrow: 1,
    overflow: 'auto',
  },
  cards: {
    flexGrow: 1,
    overflow: 'hidden',
    overflowY: 'auto',
  },
  card: {
    display: 'flex',
    flexDirection: 'column',
  },
  input: {
    width: '250px',
    paddingBottom: '24px',
  },
  gridInput: {
    display: 'inline-flex',
  },
  nameHeader: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: '24px',
  },
  breadcrumbs: {
    flexGrow: 1,
  },
  separator: {
    borderBottom: `1px solid ${theme.palette.grey[100]}`,
    margin: '0 0 24px -24px',
    paddingBottom: '24px',
    width: 'calc(100% + 48px)',
  },
  separator: {
    borderBottom: `1px solid ${theme.palette.grey[100]}`,
    margin: '0 0 24px -24px',
    paddingBottom: '24px',
    width: 'calc(100% + 48px)',
  },
  dense: {
    paddingTop: '9px',
    paddingBottom: '9px',
    height: '14px',
  },
  cancelButton: {
    marginRight: '8px',
  },
}));

const workOrderTypeQuery = graphql`
  query AddWorkOrderCardTypeQuery($workOrderTypeId: ID!) {
    workOrderType: node(id: $workOrderTypeId) {
      __typename
      ... on WorkOrderType {
        id
        name
        description
        propertyTypes {
          id
          name
          type
          nodeType
          index
          stringValue
          intValue
          booleanValue
          floatValue
          latitudeValue
          longitudeValue
          rangeFromValue
          rangeToValue
          isEditable
          isMandatory
          isInstanceProperty
          isDeleted
          category
        }
        checkListCategoryDefinitions {
          id
          title
          description
          checklistItemDefinitions {
            id
            title
            type
            index
            enumValues
            enumSelectionMode
            helpText
          }
        }
      }
    }
  }
`;

type Props = $ReadOnly<{|
  workOrderTypeId: string,
|}>;

const AddWorkOrderCard = (props: Props) => {
  const {workOrderTypeId} = props;
  const classes = useStyles();

  const {
    workOrderType,
  }: AddWorkOrderCardTypeQueryResponse = useLazyLoadQuery<AddWorkOrderCardTypeQuery>(
    workOrderTypeQuery,
    {workOrderTypeId},
  );

  const [workOrder, setWorkOrder] = useState<?WorkOrder>(
    workOrderType?.__typename === 'WorkOrderType'
      ? {
          id: generateTempId(),
          workOrderType: null,
          workOrderTypeId: workOrderType.id,
          name: workOrderType.name,
          description: workOrderType.description,
          locationId: null,
          location: null,
          properties:
            workOrderType.propertyTypes
              ?.filter(Boolean)
              .filter(propertyType => !propertyType.isDeleted)
              .map(propType =>
                getInitialPropertyFromType(toMutablePropertyType(propType)),
              )
              .sort(sortPropertiesByIndex) ?? [],
          workOrders: [],
          owner: {id: '', email: ''},
          creationDate: '',
          installDate: '',
          status: 'PENDING',
          priority: 'NONE',
          equipmentToAdd: [],
          equipmentToRemove: [],
          linksToAdd: [],
          linksToRemove: [],
          files: [],
          images: [],
          assignedTo: null,
          projectId: null,
          checkListCategories: [],
        }
      : null,
  );

  const enqueueSnackbar = useEnqueueSnackbar();
  const history = useHistory();
  const match = useRouteMatch();

  const [editingCategories, dispatch] = useReducer<
    ChecklistCategoriesStateType,
    ChecklistCategoriesMutateStateActionType,
    ?$ElementType<
      AddEditWorkOrderTypeCard_workOrderType,
      'checkListCategoryDefinitions',
    >,
  >(
    reducer,
    workOrderType?.__typename === 'WorkOrderType'
      ? workOrderType.checkListCategoryDefinitions
      : null,
    getInitialStateFromChecklistDefinitions,
  );

  const _enqueueError = useCallback(
    (message: string) => {
      enqueueSnackbar(message, {
        children: key => (
          <SnackbarItem id={key} message={message} variant="error" />
        ),
      });
    },
    [enqueueSnackbar],
  );

  const _saveWorkOrder = () => {
    const {
      name,
      description,
      locationId,
      projectId,
      assignedTo,
      status,
      priority,
      properties,
    } = nullthrows(workOrder);
    const workOrderTypeId = nullthrows(workOrder?.workOrderTypeId);
    const variables: AddWorkOrderMutationVariables = {
      input: {
        name,
        description,
        locationId,
        workOrderTypeId,
        assigneeId: assignedTo?.id,
        projectId,
        status,
        priority,
        properties: toPropertyInput(properties),
        checkListCategories: convertChecklistCategoriesStateToInput(
          editingCategories,
        ),
      },
    };

    const callbacks: MutationCallbacks<AddWorkOrderMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          _enqueueError(errors[0].message);
        } else {
          // navigate to main page
          history.push(match.url);
        }
      },
      onError: (error: Error) => {
        _enqueueError(getGraphError(error));
      },
    };
    ServerLogger.info(LogEvents.SAVE_PROJECT_BUTTON_CLICKED, {
      source: 'workOrder_details',
    });
    AddWorkOrderMutation(variables, callbacks);
  };

  const _setWorkOrderDetail = (
    key:
      | 'name'
      | 'description'
      | 'assignedTo'
      | 'projectId'
      | 'locationId'
      | 'priority'
      | 'status',
    value,
  ) => {
    setWorkOrder(prevWorkOrder => {
      if (!prevWorkOrder) {
        return;
      }
      return {...prevWorkOrder, [`${key}`]: value};
    });
  };

  const _propertyChangedHandler = index => property =>
    // eslint-disable-next-line no-warning-comments
    // $FlowFixMe - known techdebt with Property/PropertyType flow definitions
    setWorkOrder(prevWorkOrder => {
      if (!prevWorkOrder) {
        return;
      }
      return {
        ...prevWorkOrder,
        properties: [
          ...prevWorkOrder.properties.slice(0, index),
          // eslint-disable-next-line no-warning-comments
          // $FlowFixMe - known techdebt with Property/PropertyType flow definitions
          property,
          ...prevWorkOrder.properties.slice(index + 1),
        ],
      };
    });

  const _checkListCategoryChangedHandler = updatedCategories => {
    setWorkOrder(prevWorkOrder => {
      if (!prevWorkOrder) {
        return;
      }
      return {
        ...prevWorkOrder,
        checkListCategories: updatedCategories,
      };
    });
  };

  const navigateToMainPage = () => {
    ServerLogger.info(LogEvents.WORK_ORDERS_SEARCH_NAV_CLICKED, {
      source: 'work_order_details',
    });
    history.push(match.url);
  };

  if (workOrder == null) {
    return null;
  }

  return (
    <div className={classes.root}>
      <FormContextProvider
        permissions={{
          entity: 'workorder',
          action: 'create',
        }}>
        <div className={classes.nameHeader}>
          <Breadcrumbs
            className={classes.breadcrumbs}
            breadcrumbs={[
              {
                id: 'workOrders',
                name: 'WorkOrders',
                onClick: () => navigateToMainPage(),
              },
              {
                id: `new_workOrder_` + Date.now(),
                name: 'New WorkOrder',
              },
            ]}
            size="large"
          />
          <FormSaveCancelPanel
            onCancel={navigateToMainPage}
            onSave={_saveWorkOrder}
          />
        </div>
        <div className={classes.contentRoot}>
          <div className={classes.cards}>
            <Grid container spacing={2}>
              <Grid item xs={8} sm={8} lg={8} xl={8}>
                <ExpandingPanel title="Details">
                  <NameDescriptionSection
                    name={workOrder.name}
                    description={workOrder.description}
                    onNameChange={value => _setWorkOrderDetail('name', value)}
                    onDescriptionChange={value =>
                      _setWorkOrderDetail('description', value)
                    }
                  />
                  <div className={classes.separator} />
                  <Grid container spacing={2}>
                    <Grid item xs={12} sm={6} lg={4} xl={4}>
                      <FormField label="Project">
                        <ProjectTypeahead
                          className={classes.gridInput}
                          margin="dense"
                          onProjectSelection={project =>
                            _setWorkOrderDetail('projectId', project?.id)
                          }
                        />
                      </FormField>
                    </Grid>
                    {workOrder.workOrderType && (
                      <Grid item xs={12} sm={6} lg={4} xl={4}>
                        <FormField label="Type">
                          <TextField
                            disabled
                            variant="outlined"
                            margin="dense"
                            className={classes.gridInput}
                            value={workOrder.workOrderType.name}
                          />
                        </FormField>
                      </Grid>
                    )}
                    <Grid item xs={12} sm={6} lg={4} xl={4}>
                      <FormField label="Priority">
                        <Select
                          options={priorityValues}
                          selectedValue={workOrder.priority}
                          onChange={value =>
                            _setWorkOrderDetail('priority', value)
                          }
                        />
                      </FormField>
                    </Grid>
                    <Grid item xs={12} sm={6} lg={4} xl={4}>
                      <FormField label="Status">
                        <Select
                          options={statusValues}
                          selectedValue={workOrder.status}
                          onChange={value => {
                            _setWorkOrderDetail('status', value);
                          }}
                        />
                      </FormField>
                    </Grid>
                    <Grid item xs={12} sm={6} lg={4} xl={4}>
                      <FormField label="Location">
                        <LocationTypeahead
                          headline={null}
                          className={classes.gridInput}
                          margin="dense"
                          onLocationSelection={location =>
                            _setWorkOrderDetail(
                              'locationId',
                              location?.id ?? null,
                            )
                          }
                        />
                      </FormField>
                    </Grid>
                    {workOrder.properties
                      .filter(property => !property.propertyType.isDeleted)
                      .map((property, index) => (
                        <Grid
                          key={property.id}
                          item
                          xs={12}
                          sm={6}
                          lg={4}
                          xl={4}>
                          <PropertyValueInput
                            required={!!property.propertyType.isMandatory}
                            disabled={!property.propertyType.isInstanceProperty}
                            label={property.propertyType.name}
                            className={classes.gridInput}
                            margin="dense"
                            inputType="Property"
                            property={property}
                            headlineVariant="form"
                            fullWidth={true}
                            onChange={_propertyChangedHandler(index)}
                          />
                        </Grid>
                      ))}
                  </Grid>
                </ExpandingPanel>
                <ChecklistCategoriesMutateDispatchContext.Provider
                  value={dispatch}>
                  <CheckListCategoryExpandingPanel
                    categories={editingCategories}
                  />
                </ChecklistCategoriesMutateDispatchContext.Provider>
              </Grid>
              <Grid item xs={4} sm={4} lg={4} xl={4}>
                <ExpandingPanel title="Team">
                  <FormField className={classes.input} label="Assignee">
                    <UserTypeahead
                      onUserSelection={user =>
                        _setWorkOrderDetail('assignedTo', user)
                      }
                      margin="dense"
                    />
                  </FormField>
                </ExpandingPanel>
              </Grid>
            </Grid>
          </div>
        </div>
      </FormContextProvider>
    </div>
  );
};

export default AddWorkOrderCard;
