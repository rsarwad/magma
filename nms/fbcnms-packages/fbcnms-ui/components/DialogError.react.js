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
import ErrorOutlineIcon from '@material-ui/icons/ErrorOutline';
import Typography from '@material-ui/core/Typography';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  root: {
    display: 'flex',
    flexDirection: 'row',
    borderColor: 'red',
    borderWidth: '1px',
    borderRadius: '4px',
    borderStyle: 'solid',
    margin: 'auto 8px',
    alignItems: 'center',
  },
  errorIcon: {
    margin: '2px 8px',
    color: 'red',
  },
}));

type Props = {
  message: string,
};

const DialogError = ({message}: Props) => {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <ErrorOutlineIcon className={classes.errorIcon} />
      <Typography variant="subtitle2" color="error">
        {message}
      </Typography>
    </div>
  );
};

export default DialogError;
