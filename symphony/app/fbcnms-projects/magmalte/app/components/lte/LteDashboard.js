/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import AppBar from '@material-ui/core/AppBar';
import EnodebKPIs from '../EnodebKPIs';
import EventAlertChart from '../EventAlertChart';
import GatewayKPIs from '../GatewayKPIs';
import Grid from '@material-ui/core/Grid';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React, {useState} from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '@fbcnms/ui/components/design-system/Text';
import moment from 'moment';
import {Alarm, GpsFixed, NetworkCheck} from '@material-ui/icons';
import {DateTimePicker} from '@material-ui/pickers';
import {Redirect, Route, Switch} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
  },
  topBar: {
    backgroundColor: theme.palette.magmalte.background,
    padding: '20px 40px 20px 40px',
  },
  tabBar: {
    backgroundColor: theme.palette.magmalte.appbar,
    padding: '0 0 0 20px',
  },
  tabs: {
    color: 'white',
  },
  tab: {
    fontSize: '18px',
    textTransform: 'none',
  },
  tabLabel: {
    padding: '20px 0 20px 0',
  },
  tabIconLabel: {
    verticalAlign: 'middle',
    margin: '0 5px 3px 0',
  },
  // TODO: remove this when we actually fill out the grid sections
  contentPlaceholder: {
    padding: '50px 0',
  },
  input: {
    color: 'white',
  },
}));

function LteDashboard() {
  const classes = useStyles();
  const {relativePath, relativeUrl} = useRouter();

  // datetime picker
  const [startDate, setStartDate] = useState(moment().subtract(3, 'hours'));
  const [endDate, setEndDate] = useState(moment());

  return (
    <>
      <div className={classes.topBar}>
        <Text color="light" weight="medium">
          Dashboard
        </Text>
      </div>

      <AppBar position="static" color="default" className={classes.tabBar}>
        <Grid container>
          <Grid item xs={6}>
            <Tabs
              value={0}
              indicatorColor="primary"
              TabIndicatorProps={{style: {height: '5px'}}}
              textColor="inherit"
              className={classes.tabs}>
              <Tab
                key="Network"
                component={NestedRouteLink}
                label={<NetworkDashboardTabLabel />}
                to="/network"
                className={classes.tab}
              />
            </Tabs>
          </Grid>
          <Grid item xs={6}>
            <Grid container justify="flex-end" alignItems="center" spacing={1}>
              <Grid item>
                <Text color="light">Filter By Date</Text>
              </Grid>
              <Grid item>
                <DateTimePicker
                  autoOk
                  variant="inline"
                  inputVariant="outlined"
                  maxDate={endDate}
                  disableFuture
                  value={startDate}
                  inputProps={{className: classes.input}}
                  onChange={setStartDate}
                />
              </Grid>
              <Grid item>
                <Text color="light">To</Text>
              </Grid>
              <Grid item>
                <DateTimePicker
                  autoOk
                  variant="inline"
                  inputVariant="outlined"
                  disableFuture
                  value={endDate}
                  inputProps={{className: classes.input}}
                  onChange={setEndDate}
                />
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </AppBar>

      <Switch>
        <Route
          path={relativePath('/network')}
          render={props => (
            <LteNetworkDashboard {...props} startEnd={[startDate, endDate]} />
          )}
        />
        <Redirect to={relativeUrl('/network')} />
      </Switch>
    </>
  );
}

function LteNetworkDashboard({startEnd}: {startEnd: [moment, moment]}) {
  const classes = useStyles();

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={3}>
        <Grid item xs={12}>
          <Paper>
            <EventAlertChart startEnd={startEnd} />
          </Paper>
        </Grid>

        <Grid item xs={12}>
          <Text>
            <Alarm /> Alerts (92)
          </Text>
          <Paper>
            <div className={classes.contentPlaceholder}>
              Alert Table Goes Here
            </div>
          </Paper>
        </Grid>

        <Grid item xs={6}>
          <Paper elevation={2}>
            <GatewayKPIs />
          </Paper>
        </Grid>

        <Grid item xs={6}>
          <Paper elevation={2}>
            <EnodebKPIs />
          </Paper>
        </Grid>

        <Grid item xs={12}>
          <Text>
            <GpsFixed /> Events (388)
          </Text>
          <Paper>
            <div className={classes.contentPlaceholder}>
              Events Table Goes Here
            </div>
          </Paper>
        </Grid>
      </Grid>
    </div>
  );
}

function NetworkDashboardTabLabel() {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      <NetworkCheck className={classes.tabIconLabel} /> Network
    </div>
  );
}

export default LteDashboard;
