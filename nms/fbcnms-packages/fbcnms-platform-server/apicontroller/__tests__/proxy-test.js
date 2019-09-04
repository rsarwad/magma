/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {networkIdFilter, networksResponseDecorator} from '../routes';

// $FlowIgnore Ignoring error for tests.
const testNetworkIdFilter = async req => await networkIdFilter(req);

const testNetworksResponseDecorator = async (
  proxyResponse,
  proxyResponseData,
  userRequest,
  {},
) =>
  await networksResponseDecorator(
    // $FlowIgnore Ignoring error for tests.
    proxyResponse,
    proxyResponseData,
    // $FlowIgnore Ignoring error for tests.
    userRequest,
    // $FlowIgnore Ignoring error for tests.
    {}, // userResponse
  );

const createData = (
  controllerNetworks: Array<string>,
  orgNetworks: Array<string>,
  userNetworks: Array<string>,
  isSuperUser: boolean,
) => {
  const proxyResponse = {statusCode: 200};
  const proxyResponseData = new Buffer(
    JSON.stringify(controllerNetworks),
    'utf8',
  );

  const userRequest = {
    organization: async () => ({networkIDs: orgNetworks}),
    user: {isSuperUser, networkIDs: userNetworks},
  };

  return {proxyResponse, proxyResponseData, userRequest};
};

describe('Proxy test', () => {
  describe('networkIdFilter', () => {
    const params = {
      networkID: 'test',
    };

    it('allows superuser access to org network', async () => {
      const req = {
        params,
        organization: () => ({
          networkIDs: ['test'],
        }),
        user: {
          isSuperUser: true,
          networkIDs: [],
        },
      };
      const isAllowed = await testNetworkIdFilter(req);
      expect(isAllowed).toBe(true);
    });

    it('disallows superuser access to non-org network', async () => {
      const req = {
        params,
        organization: () => ({
          networkIDs: ['not-test'],
        }),
        user: {
          isSuperUser: true,
          networkIDs: [],
        },
      };
      const isAllowed = await testNetworkIdFilter(req);
      expect(isAllowed).toBe(false);
    });

    it('allows user with access to network on specific org', async () => {
      const req = {
        params,
        organization: () => ({
          networkIDs: ['test'],
        }),
        user: {
          isSuperUser: false,
          networkIDs: ['test'],
        },
      };
      const isAllowed = await testNetworkIdFilter(req);
      expect(isAllowed).toBe(true);
    });

    it('disallows user with access to network but not org', async () => {
      const req = {
        params,
        organization: () => ({
          networkIDs: ['not-test'],
        }),
        user: {
          isSuperUser: false,
          networkIDs: ['test'],
        },
      };
      const isAllowed = await testNetworkIdFilter(req);
      expect(isAllowed).toBe(false);
    });

    it('disallows user without access to network', async () => {
      const req = {
        params,
        organization: () => ({
          networkIDs: ['test', 'not-test'],
        }),
        user: {
          isSuperUser: false,
          networkIDs: ['not-test'],
        },
      };
      const isAllowed = await testNetworkIdFilter(req);
      expect(isAllowed).toBe(false);
    });
  });

  describe('networksResponseDecorator', () => {
    it('allows normal user access to her network only', async () => {
      const {proxyResponse, proxyResponseData, userRequest} = createData(
        ['network1', 'network2'],
        ['network2'],
        ['network2'],
        false, // isSuperUser
      );

      const result = await testNetworksResponseDecorator(
        proxyResponse,
        proxyResponseData,
        userRequest,
        {}, // userResponse
      );

      expect(JSON.parse(result)).toEqual(['network2']);
    });

    it('allows super user access to all org networks', async () => {
      const {proxyResponse, proxyResponseData, userRequest} = createData(
        ['network1', 'network2', 'network3'],
        ['network1', 'network2'],
        [],
        true, // isSuperUser
      );

      const result = await testNetworksResponseDecorator(
        proxyResponse,
        proxyResponseData,
        userRequest,
        {}, // userResponse
      );

      expect(JSON.parse(result)).toEqual(['network1', 'network2']);
    });

    it('denies normal user access to a network not in the org', async () => {
      const {proxyResponse, proxyResponseData, userRequest} = createData(
        ['network1', 'network2', 'network3'],
        ['network1', 'network2'],
        ['network1', 'network3', 'network4'],
        false, // isSuperUser
      );

      const result = await testNetworksResponseDecorator(
        proxyResponse,
        proxyResponseData,
        userRequest,
        {}, // userResponse
      );

      expect(JSON.parse(result)).toEqual(['network1']);
    });

    it('denies super user access to a network not in the controller', async () => {
      const {proxyResponse, proxyResponseData, userRequest} = createData(
        ['network1', 'network2', 'network3'],
        ['network1', 'network2', 'network4'],
        [],
        true, // isSuperUser
      );

      const result = await testNetworksResponseDecorator(
        proxyResponse,
        proxyResponseData,
        userRequest,
        {}, // userResponse
      );

      expect(JSON.parse(result)).toEqual(['network1', 'network2']);
    });
  });
});
