/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import 'jest-dom/extend-expect';
import DashboardAlertTable from '../DashboardAlertTable';
import MagmaAPIBindings from '@fbcnms/magma-api';
import React from 'react';
import axiosMock from 'axios';
import {MemoryRouter, Route} from 'react-router-dom';
import {cleanup, fireEvent, render, wait} from '@testing-library/react';
import type {gettable_alert, prom_firing_alert} from '@fbcnms/magma-api';

afterEach(cleanup);

const tbl_alert: gettable_alert = {
  name: 'null_receiver',
};
const mockAlertSt: Array<prom_firing_alert> = [
  {
    annotations: {
      description: 'TestMetric1 Description',
      summary: 'TestMetric1 Minor Alert',
    },
    endsAt: '2020-05-14T18:55:25.844Z',
    fingerprint: '0de443c4dd7af53e',
    receivers: tbl_alert,
    startsAt: '2020-05-14T18:27:25.844Z',
    status: {inhibitedBy: [], silencedBy: [], state: 'active'},
    updatedAt: '2020-05-14T18:52:31.971Z',
    generatorURL:
      'http://1521e855c607:9090/graph?g0.expr=TestMetric9+%3E+1\u0026g0.tab=1',
    labels: {
      alertname: 'TestAlert1',
      instance: '192.168.0.124:2112',
      job: 'myapp',
      monitor: 'master',
      networkID: 'test',
      severity: 'critical',
    },
  },
  {
    annotations: {
      description: 'TestMetric2 Description',
      summary: 'TestMetric2 Major Alert',
    },
    endsAt: '2020-05-14T18:55:25.844Z',
    fingerprint: '0de443c4dd7af53e',
    receivers: tbl_alert,
    startsAt: '2020-05-14T18:27:25.844Z',
    status: {inhibitedBy: [], silencedBy: [], state: 'active'},
    updatedAt: '2020-05-14T18:52:31.971Z',
    generatorURL:
      'http://1521e855c607:9090/graph?g0.expr=TestMetric9+%3E+1\u0026g0.tab=1',
    labels: {
      alertname: 'TestAlert2',
      instance: '192.168.0.124:2112',
      job: 'myapp',
      monitor: 'master',
      networkID: 'test',
      severity: 'major',
    },
  },
  {
    annotations: {
      description: 'TestMetric3 Description',
      summary: 'TestMetric3 Critical Alert',
    },
    endsAt: '2020-05-14T18:55:25.844Z',
    fingerprint: '0de443c4dd7af53e',
    receivers: tbl_alert,
    startsAt: '2020-05-14T18:27:25.844Z',
    status: {inhibitedBy: [], silencedBy: [], state: 'active'},
    updatedAt: '2020-05-14T18:52:31.971Z',
    generatorURL:
      'http://1521e855c607:9090/graph?g0.expr=TestMetric9+%3E+1\u0026g0.tab=1',
    labels: {
      alertname: 'TestAlert3',
      instance: '192.168.0.124:2112',
      job: 'myapp',
      monitor: 'master',
      networkID: 'test',
      severity: 'minor',
    },
  },
  {
    annotations: {
      description: 'TestMetric4 Description',
      summary: 'TestMetric1 Other Alert',
    },
    endsAt: '2020-05-14T18:55:25.844Z',
    fingerprint: '0de443c4dd7af53e',
    receivers: tbl_alert,
    startsAt: '2020-05-14T18:27:25.844Z',
    status: {inhibitedBy: [], silencedBy: [], state: 'active'},
    updatedAt: '2020-05-14T18:52:31.971Z',
    generatorURL:
      'http://1521e855c607:9090/graph?g0.expr=TestMetric9+%3E+1\u0026g0.tab=1',
    labels: {
      alertname: 'TestAlert4',
      instance: '192.168.0.124:2112',
      job: 'myapp',
      monitor: 'master',
      networkID: 'test',
      severity: 'normal',
    },
  },
];

jest.mock('axios');
jest.mock('@fbcnms/magma-api');
jest.mock('@fbcnms/ui/hooks/useSnackbar');

describe('<DashboardAlertTable />', () => {
  beforeEach(() => {
    MagmaAPIBindings.getNetworksByNetworkIdAlerts.mockResolvedValue(
      mockAlertSt,
    );
  });

  afterEach(() => {
    axiosMock.get.mockClear();
  });

  const Wrapper = () => (
    <MemoryRouter initialEntries={['/nms/mynetwork']} initialIndex={0}>
      <Route path="/nms/:networkId" component={DashboardAlertTable} />
    </MemoryRouter>
  );

  it('renders', async () => {
    const {getByTestId, getByText} = render(<Wrapper />);
    await wait();
    expect(MagmaAPIBindings.getNetworksByNetworkIdAlerts).toHaveBeenCalledTimes(
      1,
    );

    // each sections have only one row so querying rowID 0
    const rowIdx = 0;
    // check if the default is critical alert sections
    expect(getByTestId('alertName' + rowIdx)).toHaveTextContent('TestAlert1');
    fireEvent.click(getByText('Critical'));
    expect(getByTestId('alertName' + rowIdx)).toHaveTextContent('TestAlert1');

    fireEvent.click(getByText('Major'));
    expect(getByTestId('alertName' + rowIdx)).toHaveTextContent('TestAlert2');

    fireEvent.click(getByText('Minor'));
    expect(getByTestId('alertName' + rowIdx)).toHaveTextContent('TestAlert3');

    fireEvent.click(getByText('Other'));
    expect(getByTestId('alertName' + rowIdx)).toHaveTextContent('TestAlert4');

    expect(getByText('Alerts (4)')).toBeInTheDocument();
  });
});
