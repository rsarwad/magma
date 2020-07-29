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
 * @format
 */

module.exports = {
  collectCoverageFrom: [
    '**/fbcnms-projects/**/*.js',
    '**/fbcnms-packages/**/*.js',
    '!**/__mocks__/**',
    '!**/__tests__/**',
    '!**/fbcnms-packages/fbcnms-ui/stories/**',
    '!**/thrift/gen-nodejs/**',
    '!**/node_modules/**',
    '!**/fbcnms-packages/fbcnms-test/**',
  ],

  coverageReporters: ['json', 'html'],
  modulePathIgnorePatterns: [],
  projects: [
    {
      name: 'server',
      testEnvironment: 'node',
      testMatch: [
        '<rootDir>/__tests__/*.js',
        '<rootDir>/fbcnms-projects/**/server/**/__tests__/*.js',
        '<rootDir>/fbcnms-packages/fbcnms-auth/**/__tests__/*.js',
        '<rootDir>/fbcnms-packages/fbcnms-express-middleware/**/__tests__/*.js',
        '<rootDir>/fbcnms-packages/fbcnms-platform-server/**/__tests__/*.js',
        '<rootDir>/fbcnms-projects/platform-server/**/__tests__/*.js',
        // run app/server shared tests in both node and jsdom environments
        '<rootDir>/fbcnms-packages/fbcnms-util/**/__tests__/*.js',
        '<rootDir>/fbcnms-projects/**/shared/**/__tests__/*.js',
      ],
      transform: {
        '^.+\\.js$': 'babel-jest',
      },
      transformIgnorePatterns: ['/node_modules/(?!@fbcnms)'],
    },
    {
      moduleNameMapper: {
        '\\.(jpg|jpeg|png|gif|eot|otf|webp|svg|ttf|woff|woff2|mp4|webm|wav|mp3|m4a|aac|oga)$':
          '<rootDir>/__mocks__/fileMock.js',
        '\\.(css|less)$': 'identity-obj-proxy',
      },
      name: 'app',
      setupFiles: [require.resolve('@fbcnms/babel-register/polyfill')],
      testEnvironment: 'jsdom',
      testMatch: [
        '<rootDir>/fbcnms-projects/**/app/**/__tests__/*.js',
        '<rootDir>/fbcnms-packages/fbcnms-ui/**/__tests__/*.js',
        '<rootDir>/fbcnms-packages/fbcnms-alarms/(components|hooks)/__tests__/*.js',
        // run app/server shared tests in both node and jsdom environments
        '<rootDir>/fbcnms-packages/fbcnms-util/**/__tests__/*.js',
        '<rootDir>/fbcnms-packages/fbcnms-mobileapp/**/__tests__/*.js',
        '<rootDir>/fbcnms-projects/**/shared/**/__tests__/*.js',
      ],
      transform: {
        '^.+\\.js$': 'babel-jest',
      },
      transformIgnorePatterns: ['/node_modules/(?!@fbcnms)'],
    },
  ],
  testEnvironment: 'jsdom',
  testPathIgnorePatterns: ['/node_modules/'],
};
