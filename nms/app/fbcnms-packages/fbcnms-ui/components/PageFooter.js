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
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
    position: 'fixed',
    bottom: 0,
    left: 0,
    right: 0,
    backgroundColor: theme.palette.common.white,
    boxShadow: '0px -1px 4px 0px rgba(0,0,0,0.11)',
    padding: '14px 24px',
  },
  spacer: {
    flexGrow: 1,
  },
}));

type Props = {
  children: Array<React.Element<*>> | React.Element<*>,
};

const PageFooter = (props: Props) => {
  const {children} = props;
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <div className={classes.spacer} />
      {children}
    </div>
  );
};

export default PageFooter;
