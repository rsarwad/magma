/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
 * @format
 */
import AppBar from '@material-ui/core/AppBar';
import DashboardAlertTable from '../DashboardAlertTable';
import DashboardKPIs from '../DashboardKPIs';
import EventAlertChart from '../EventAlertChart';
import EventsTable from '../../views/events/EventsTable';
import Grid from '@material-ui/core/Grid';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import React, {useState} from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '../../theme/design-system/Text';
import moment from 'moment';

import {DateTimePicker} from '@material-ui/pickers';
import {NetworkCheck, People} from '@material-ui/icons';
import {Redirect, Route, Switch} from 'react-router-dom';
import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
  },
  topBar: {
    backgroundColor: colors.primary.mirage,
    padding: '20px 40px 20px 40px',
    color: colors.primary.white,
  },
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    padding: `0 ${theme.spacing(5)}px`,
  },
  tabs: {
    color: colors.primary.white,
  },
  tab: {
    fontSize: '18px',
    textTransform: 'none',
  },
  tabLabel: {
    padding: '16px 0 16px 0',
    display: 'flex',
    alignItems: 'center',
  },
  tabIconLabel: {
    marginRight: '8px',
  },
  dateTimeText: {
    color: colors.primary.selago,
  },
}));

function LteDashboard() {
  const classes = useStyles();
  const {relativePath, relativeUrl} = useRouter();

  // datetime picker
  const [startDate, setStartDate] = useState(moment().subtract(3, 'days'));
  const [endDate, setEndDate] = useState(moment());

  return (
    <>
      <div className={classes.topBar}>
        <Text variant="body2">Dashboard</Text>
      </div>

      <AppBar position="static" color="default" className={classes.tabBar}>
        <Grid container direction="row" justify="flex-end" alignItems="center">
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
                label={<DashboardTabLabel label="Network" />}
                to="/network"
                className={classes.tab}
              />
            </Tabs>
          </Grid>
          <Grid item xs={6}>
            <Grid container justify="flex-end" alignItems="center" spacing={2}>
              <Grid item>
                <Text variant="body3" className={classes.dateTimeText}>
                  Filter By Date
                </Text>
              </Grid>
              <DateTimePicker
                autoOk
                variant="inline"
                inputVariant="outlined"
                maxDate={endDate}
                disableFuture
                value={startDate}
                onChange={setStartDate}
              />
              <Grid item>
                <Text variant="body3" className={classes.dateTimeText}>
                  to
                </Text>
              </Grid>
              <DateTimePicker
                autoOk
                variant="inline"
                inputVariant="outlined"
                disableFuture
                value={endDate}
                onChange={setEndDate}
              />
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
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <EventAlertChart startEnd={startEnd} />
        </Grid>

        <Grid item xs={12}>
          <DashboardAlertTable />
        </Grid>

        <Grid item xs={12}>
          <DashboardKPIs />
        </Grid>
        <Grid item xs={12}>
          <EventsTable eventStream="NETWORK" sz="md" />
        </Grid>
      </Grid>
    </div>
  );
}

function DashboardTabLabel(props) {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      {props.label === 'Subscribers' ? (
        <People className={classes.tabIconLabel} />
      ) : props.label === 'Network' ? (
        <NetworkCheck className={classes.tabIconLabel} />
      ) : null}
      {props.label}
    </div>
  );
}

export default LteDashboard;
