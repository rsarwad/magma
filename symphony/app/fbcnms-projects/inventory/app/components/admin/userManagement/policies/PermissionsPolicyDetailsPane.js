/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {PermissionsPolicy} from '../utils/UserManagementUtils';

import * as React from 'react';
import FormFieldTextInput from '../utils/FormFieldTextInput';
import Grid from '@material-ui/core/Grid';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    backgroundColor: symphony.palette.white,
    paddingBottom: '16px',
  },
  nameField: {
    marginRight: '8px',
  },
  descriptionField: {
    marginTop: '8px',
  },
}));

type Props = $ReadOnly<{
  policy: PermissionsPolicy,
  onChange: PermissionsPolicy => void,
  className?: ?string,
}>;

export default function PermissionsPolicyDetailsPane(props: Props) {
  const {policy, className, onChange} = props;
  const classes = useStyles();

  return (
    <div className={classNames(classes.root, className)}>
      <ViewContainer header={{title: <fbt desc="">Policy Details</fbt>}}>
        <Grid container>
          <Grid item xs={12} sm={6} lg={6} xl={6}>
            <FormFieldTextInput
              className={classes.nameField}
              label={`${fbt('Policy Name', '')}`}
              validationId="name"
              value={policy.name}
              onValueChanged={name => {
                onChange({
                  ...policy,
                  name,
                });
              }}
            />
          </Grid>
          <Grid item xs={12}>
            <FormFieldTextInput
              className={classes.descriptionField}
              label={`${fbt('Policy Description', '')}`}
              value={policy.description || ''}
              onValueChanged={description => {
                onChange({
                  ...policy,
                  description,
                });
              }}
            />
          </Grid>
        </Grid>
      </ViewContainer>
    </div>
  );
}
