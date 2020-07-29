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
import type {EnodebInfo} from '../../components/lte/EnodebUtils';
import type {lte_gateway} from '@fbcnms/magma-api';

import AccessAlarmIcon from '@material-ui/icons/AccessAlarm';
import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import DashboardIcon from '@material-ui/icons/Dashboard';
import EventsTable from '../../views/events/EventsTable';
import GatewayConfig from './GatewayDetailConfig';
import GatewayDetailStatus from './GatewayDetailStatus';
import GatewayLogs from './GatewayLogs';
import GatewaySummary from './GatewaySummary';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';
import Grid from '@material-ui/core/Grid';
import ListAltIcon from '@material-ui/icons/ListAlt';
import MyLocationIcon from '@material-ui/icons/MyLocation';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '../../theme/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';

import {CardTitleRow} from '../../components/layout/CardTitleRow';
import {GatewayJsonConfig} from './GatewayDetailConfig';
import {GetCurrentTabPos} from '../../components/TabUtils.js';
import {Redirect, Route, Switch} from 'react-router-dom';
import {colors, typography} from '../../theme/default';
import {magmaEventTypes} from '../../views/events/EventsTable';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

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
  appBarBtn: {
    color: colors.primary.white,
    background: colors.primary.comet,
    fontFamily: typography.button.fontFamily,
    fontWeight: typography.button.fontWeight,
    fontSize: typography.button.fontSize,
    lineHeight: typography.button.lineHeight,
    letterSpacing: typography.button.letterSpacing,

    '&:hover': {
      background: colors.primary.mirage,
    },
  },
  appBarBtnSecondary: {
    color: colors.primary.white,
  },
  paper: {
    textAlign: 'center',
    padding: theme.spacing(10),
  },
}));

export function GatewayDetail({
  lteGateways,
  enbInfo,
}: {
  lteGateways: {[string]: lte_gateway},
  enbInfo: {[string]: EnodebInfo},
}) {
  const classes = useStyles();
  const {relativePath, relativeUrl, match} = useRouter();
  const gatewayId: string = nullthrows(match.params.gatewayId);
  const [gwInfo, setGwInfo] = useState(lteGateways[gatewayId]);
  const gwEnbs =
    gwInfo.connected_enodeb_serials?.reduce(
      (enbs: {[string]: EnodebInfo}, serial: string) => {
        if (enbInfo[serial]) {
          enbs[serial] = enbInfo[serial];
        }
        return enbs;
      },
      {},
    ) || {};

  return (
    <>
      <div className={classes.topBar}>
        <Text variant="body2">Equipment/{gatewayId}</Text>
      </div>

      <AppBar position="static" color="default" className={classes.tabBar}>
        <Grid container direction="row" justify="flex-end" alignItems="center">
          <Grid item xs={8}>
            <Tabs
              value={GetCurrentTabPos(match.url, [
                'overview',
                'event',
                'log',
                'alert',
                'config',
              ])}
              indicatorColor="primary"
              TabIndicatorProps={{style: {height: '5px'}}}
              textColor="inherit"
              className={classes.tabs}>
              <Tab
                key="Overview"
                component={NestedRouteLink}
                label={<OverviewTabLabel />}
                to="/overview"
                className={classes.tab}
              />
              <Tab
                key="Event"
                component={NestedRouteLink}
                label={<EventTabLabel />}
                to="/event"
                className={classes.tab}
              />
              <Tab
                key="Log"
                component={NestedRouteLink}
                label={<LogTabLabel />}
                to="/logs"
                className={classes.tab}
              />
              <Tab
                key="Alert"
                component={NestedRouteLink}
                label={<AlertTabLabel />}
                to="/alert"
                className={classes.tab}
              />
              <Tab
                key="Config"
                component={NestedRouteLink}
                label={<ConfigTabLabel />}
                to="/config"
                className={classes.tab}
              />
            </Tabs>
          </Grid>
          <Grid
            item
            xs={4}
            direction="row"
            justify="flex-end"
            alignItems="center">
            <Grid container justify="flex-end" alignItems="center" spacing={2}>
              <Grid item>
                <Button variant="text" className={classes.appBarBtnSecondary}>
                  Secondary Action
                </Button>
              </Grid>
              <Grid item>
                <Button variant="contained" className={classes.appBarBtn}>
                  Reboot
                </Button>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </AppBar>

      <Switch>
        <Route
          path={relativePath('/config/json')}
          render={() => (
            <GatewayJsonConfig
              gwInfo={gwInfo}
              onSave={gateway => {
                setGwInfo(gateway);
              }}
            />
          )}
        />
        <Route
          path={relativePath('/config')}
          render={() => (
            <GatewayConfig
              gwInfo={gwInfo}
              enbInfo={gwEnbs}
              onSave={gateway => {
                setGwInfo(gateway);
              }}
            />
          )}
        />
        <Route
          path={relativePath('/event')}
          render={() => (
            <EventsTable
              eventTypes={magmaEventTypes.GATEWAY}
              eventKey={gwInfo.device.hardware_id}
              sz="lg"
            />
          )}
        />
        <Route
          path={relativePath('/overview')}
          render={() => <GatewayOverview gwInfo={gwInfo} enbInfo={gwEnbs} />}
        />
        <Route path={relativePath('/logs')} component={GatewayLogs} />
        <Redirect to={relativeUrl('/overview')} />
      </Switch>
    </>
  );
}

function GatewayOverview({gwInfo}: {gwInfo: lte_gateway}) {
  const classes = useStyles();
  const {match} = useRouter();
  const gatewayId: string = nullthrows(match.params.gatewayId);

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12} md={6}>
          <Grid container spacing={4} direction="column">
            <Grid item xs={12} alignItems="center">
              <CardTitleRow icon={CellWifiIcon} label={gatewayId} />
              <GatewaySummary gwInfo={gwInfo} />
            </Grid>
            <Grid item xs={12} alignItems="center">
              <CardTitleRow icon={MyLocationIcon} label="Events" />
              <EventsTable
                eventTypes={magmaEventTypes.GATEWAY}
                eventKey={gwInfo.device.hardware_id}
                sz="sm"
              />
            </Grid>
          </Grid>
        </Grid>
        <Grid item xs={12} md={6}>
          <Grid container spacing={4} direction="column">
            <Grid item>
              <CardTitleRow icon={GraphicEqIcon} label="Status" />
              <GatewayDetailStatus gwInfo={gwInfo} />
            </Grid>
            <Grid item>
              <CardTitleRow
                icon={SettingsInputAntennaIcon}
                label="Connected eNodeBs"
              />
              <Paper className={classes.paper} elevation={0}>
                <Text variant="body2">Connected eNodeB Information</Text>
              </Paper>
            </Grid>
            <Grid item>
              <CardTitleRow icon={PeopleIcon} label="Subscribers" />
              <Paper className={classes.paper} elevation={0}>
                <Text variant="body2">Subscribers data</Text>
              </Paper>
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

function OverviewTabLabel() {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      <DashboardIcon className={classes.tabIconLabel} /> Overview
    </div>
  );
}

function ConfigTabLabel() {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      <SettingsIcon className={classes.tabIconLabel} /> Config
    </div>
  );
}

function EventTabLabel() {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      <MyLocationIcon className={classes.tabIconLabel} /> Event
    </div>
  );
}

function LogTabLabel() {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      <ListAltIcon className={classes.tabIconLabel} /> Logs
    </div>
  );
}

function AlertTabLabel() {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      <AccessAlarmIcon className={classes.tabIconLabel} /> Alerts
    </div>
  );
}

export default GatewayDetail;
