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
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  toggleContainer: {
    display: 'flex',
    alignItems: 'flex-end',
    justifyContent: 'flex-start',
    background: theme.palette.background.default,
    padding: 0,
    border: 0,
    margin: '10px',
    position: 'absolute',
    bottom: '0',
    left: 0,
    borderStyle: 'solid',
    borderRadius: '4px',
  },
}));

type Props = {children: ?React.Node};

const MapToggleContainer = (props: Props) => {
  const classes = useStyles();
  return <div className={classes.toggleContainer}>{props.children}</div>;
};

export default MapToggleContainer;
