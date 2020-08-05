/**
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
 * @flow
 * @format
 */

const {
  DEV_MODE,
  LOG_FORMAT,
  LOG_LEVEL,
} = require('@fbcnms/platform-server/config');

// This must be done before any module imports to configure
// logging correctly
const logging = require('@fbcnms/logging');
logging.configure({
  LOG_FORMAT,
  LOG_LEVEL,
});

const {
  appMiddleware,
  csrfMiddleware,
  organizationMiddleware,
  sessionMiddleware,
  webpackSmartMiddleware,
} = require('@fbcnms/express-middleware');
const connectSession = require('connect-session-sequelize');
const express = require('express');
const passport = require('passport');
const path = require('path');
const fbcPassport = require('@fbcnms/auth/passport').default;
const session = require('express-session');
const {sequelize} = require('@fbcnms/sequelize-models');
const OrganizationLocalStrategy = require('@fbcnms/auth/strategies/OrganizationLocalStrategy')
  .default;
const OrganizationSamlStrategy = require('@fbcnms/auth/strategies/OrganizationSamlStrategy')
  .default;

const {access, configureAccess} = require('@fbcnms/auth/access');
const {
  AccessRoles: {SUPERUSER, USER},
} = require('@fbcnms/auth/roles');

import type {ExpressResponse} from 'express';
import type {FBCNMSRequest} from '@fbcnms/auth/access';

const devMode = process.env.NODE_ENV !== 'production';

// Create Sequelize Store
const SessionStore = connectSession(session.Store);
const sequelizeSessionStore = new SessionStore({db: sequelize});

// FBC express initialization
const app = express<FBCNMSRequest, ExpressResponse>();
app.set('trust proxy', 1);
app.use(organizationMiddleware());
app.use(appMiddleware());
app.use(
  sessionMiddleware({
    devMode,
    sessionStore: sequelizeSessionStore,
    sessionToken:
      process.env.SESSION_TOKEN || 'fhcfvugnlkkgntihvlekctunhbbdbjiu',
  }),
);
app.use(passport.initialize());
app.use(passport.session()); // must be after sessionMiddleware

fbcPassport.use();
passport.use('local', OrganizationLocalStrategy());
passport.use(
  'saml',
  OrganizationSamlStrategy({
    urlPrefix: '/user',
  }),
);

// Views
app.set('views', path.join(__dirname, '..', 'views'));
app.set('view engine', 'pug');

// Routes
app.use(
  webpackSmartMiddleware({
    devMode: DEV_MODE,
    devWebpackConfig: require('../config/webpack.config.js'),
    distPath: require('../config/paths').distPath,
  }),
);
app.use('/user', require('@fbcnms/auth/express').unprotectedUserRoutes());

app.use(configureAccess({loginUrl: '/user/login'}));

// Grafana uses its own CSRF, so we don't need to handle it on our side.
// Grafana can access all metrics of an org, so it must be restricted
// to superusers
app.use(
  '/grafana',
  access(SUPERUSER),
  require('@fbcnms/platform-server/grafana/routes').default,
);

app.use('/', csrfMiddleware(), access(USER), require('./main/routes').default);

export default app;
