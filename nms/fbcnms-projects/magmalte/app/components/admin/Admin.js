/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import AdminMain from './AdminMain';
import AppContext from '@fbcnms/ui/context/AppContext';
import AssignmentIcon from '@material-ui/icons/Assignment';
import AuditLog from './AuditLog';
import NavListItem from '@fbcnms/ui/components/NavListItem.react';
import Networks from './Networks';
import Paper from '@material-ui/core/Paper';
import PeopleIcon from '@material-ui/icons/People';
import SecuritySettings from '../SecuritySettings';
import SignalCellularAlt from '@material-ui/icons/SignalCellularAlt';
import UsersSettings from '../UsersSettings';
import {Redirect, Route, Switch} from 'react-router-dom';

import {makeStyles} from '@material-ui/styles';
import {useFeatureFlag, useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  paper: {
    margin: theme.spacing(3),
    padding: theme.spacing(),
  },
}));

function NavItems() {
  const {relativeUrl} = useRouter();
  const auditLogEnabled = useFeatureFlag(AppContext, 'audit_log_view');
  const networkManagementEnabled = useFeatureFlag(
    AppContext,
    'magma_network_management',
  );

  return (
    <>
      <NavListItem
        label="Users"
        path={relativeUrl('/users')}
        icon={<PeopleIcon />}
      />
      {auditLogEnabled && (
        <NavListItem
          label="Audit Log"
          path={relativeUrl('/audit_log')}
          icon={<AssignmentIcon />}
        />
      )}
      {networkManagementEnabled && (
        <NavListItem
          label="Networks"
          path={relativeUrl('/networks')}
          icon={<SignalCellularAlt />}
        />
      )}
    </>
  );
}

function NavRoutes() {
  const classes = useStyles();
  const {relativeUrl} = useRouter();
  return (
    <Switch>
      <Route path={relativeUrl('/users')} component={UsersSettings} />
      <Route path={relativeUrl('/audit_log')} component={AuditLog} />
      <Route path={relativeUrl('/networks')} component={Networks} />
      <Route
        path={relativeUrl('/settings')}
        render={() => (
          <Paper className={classes.paper}>
            <SecuritySettings />
          </Paper>
        )}
      />
      <Redirect to={relativeUrl('/users')} />
    </Switch>
  );
}

export default () => (
  <AdminMain navRoutes={() => <NavRoutes />} navItems={() => <NavItems />} />
);
