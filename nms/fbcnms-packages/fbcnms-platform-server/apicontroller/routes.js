/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ExpressResponse} from 'express';
import type {FBCNMSRequest} from '@fbcnms/auth/access';

const express = require('express');
const proxy = require('express-http-proxy');
const HttpsProxyAgent = require('https-proxy-agent');
const url = require('url');
const {apiCredentials, API_HOST, NETWORK_FALLBACK} = require('../config');
import auditLoggingDecorator from './auditLoggingDecorator';

import {intersection} from 'lodash';

const router = express.Router();

let agent = null;
if (process.env.HTTPS_PROXY) {
  const options = url.parse(process.env.HTTPS_PROXY);
  agent = new HttpsProxyAgent(options);
}
const PROXY_OPTIONS = {
  https: true,
  memoizeHost: false,
  proxyReqOptDecorator: (proxyReqOpts, _originalReq) => {
    return {
      ...proxyReqOpts,
      agent: agent,
      cert: apiCredentials().cert,
      key: apiCredentials().key,
      rejectUnauthorized: false,
    };
  },
  proxyReqPathResolver: req =>
    req.originalUrl.replace(/^\/nms\/apicontroller/, ''),
};

export async function networkIdFilter(req: FBCNMSRequest): Promise<boolean> {
  // If not using organizations, always allow
  if (!req.organization) {
    return true;
  }

  const organization = await req.organization();

  // If the request isn't an organization network, block
  // the request
  const isOrganizationAllowed = containsNetworkID(
    organization.networkIDs,
    req.params.networkID,
  );
  if (!isOrganizationAllowed) {
    return false;
  }

  // super users on standalone deployments
  // have access to all proxied API requests
  // for the organization
  if (req.user.isSuperUser) {
    return true;
  }
  return containsNetworkID(req.user.networkIDs, req.params.networkID);
}

export async function networksResponseDecorator(
  proxyRes: ExpressResponse,
  proxyResData: Buffer,
  userReq: FBCNMSRequest,
  userRes: ExpressResponse,
) {
  let networkIds;
  if (
    (proxyRes.statusCode === 403 || proxyRes.statusCode === 401) &&
    NETWORK_FALLBACK.length > 0
  ) {
    // Temporary hack -- if you don't have a root magma cert,
    // it will return a 403.
    userRes.statusCode = 200;
    networkIds = NETWORK_FALLBACK;
  } else {
    networkIds = JSON.parse(proxyResData.toString('utf8'));
  }

  let result = networkIds;
  if (userReq.organization) {
    const organization = await userReq.organization();
    if (userReq.user.isSuperUser) {
      // if this is a Super User, they have access to all networks in the org
      // that are also available in the Magma controller
      result = intersection(organization.networkIDs, networkIds);
    } else {
      // otherwise, the list of networks is further restricted to what the user
      // is allowed to see
      result = intersection(
        organization.networkIDs,
        networkIds,
        userReq.user.networkIDs,
      );
    }
  }

  return JSON.stringify(result);
}

const containsNetworkID = function(
  allowedNetworkIDs: string[],
  networkID: string,
): boolean {
  return (
    allowedNetworkIDs.indexOf(networkID) !== -1 ||
    // Remove secondary condition after T34404422 is addressed. Reason:
    //   Request needs to be lower cased otherwise calling
    //   MagmaAPIUrls.gateways() potentially returns missing devices.
    allowedNetworkIDs
      .map(id => id.toString().toLowerCase())
      .indexOf(networkID.toString().toLowerCase()) !== -1
  );
};

router.use(
  /^\/magma\/networks$/,
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    userResDecorator: networksResponseDecorator,
  }),
);

router.use(
  '/magma/networks/:networkID',
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    filter: networkIdFilter,
    userResDecorator: auditLoggingDecorator,
  }),
);

router.use(
  '/magma/v1/networks/:networkID',
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    filter: networkIdFilter,
    userResDecorator: auditLoggingDecorator,
  }),
);

router.use(
  '/magma/v1/lte/:networkID',
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    filter: networkIdFilter,
    userResDecorator: auditLoggingDecorator,
  }),
);

router.use(
  '/magma/channels/:channel',
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    filter: (req, _res) => req.method === 'GET',
  }),
);

router.use('', (req: FBCNMSRequest, res: ExpressResponse) => {
  res.status(404).send('Not Found');
});

export default router;
