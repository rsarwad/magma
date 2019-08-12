/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {WithStyles} from '@material-ui/core';

import ErrorIcon from '@material-ui/icons/Error';
import React from 'react';
import Typography from '@material-ui/core/Typography';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  root: {
    padding: '8px',
    display: 'flex',
    alignItems: 'center',
  },
  errorIcon: {
    marginRight: '8px',
  },
};

type Props = {
  children: any,
} & WithStyles<typeof styles>;

type State = {
  error: ?Error,
};

class ErrorBoundary extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {error: null};
  }

  componentDidCatch(error: Error) {
    this.setState({
      error: error,
    });
    // TODO: Log JS errors here
  }

  render() {
    const {classes} = this.props;
    if (this.state.error) {
      return (
        <div className={classes.root}>
          <ErrorIcon size="small" className={classes.errorIcon} />
          <Typography variant="body1">Oops, something went wrong.</Typography>
        </div>
      );
    }
    return this.props.children;
  }
}

export default withStyles(styles)(ErrorBoundary);
