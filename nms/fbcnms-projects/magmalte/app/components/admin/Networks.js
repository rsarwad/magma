/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Button from '@material-ui/core/Button';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import NetworkDialog from './NetworkDialog';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import {MagmaAPIUrls} from '../../common/MagmaAPI';

import {Route} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {sortBy} from 'lodash';
import {useAxios, useRouter} from '@fbcnms/ui/hooks';
import {useCallback, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';

const useStyles = makeStyles({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  paper: {
    margin: '10px',
  },
});

function Networks() {
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const {relativePath, relativeUrl, history} = useRouter();
  const [networks, setNetworks] = useState(null);
  const {error, isLoading} = useAxios({
    url: MagmaAPIUrls.networks(),
    onResponse: useCallback(
      res => setNetworks(sortBy(res.data, [n => n.toLowerCase()])),
      [],
    ),
  });

  if (error || isLoading || !networks) {
    return <LoadingFiller />;
  }

  const rows = networks.map(network => (
    <TableRow key={network}>
      <TableCell>{network}</TableCell>
      <TableCell>
        <IconButton
          onClick={() => history.push(relativeUrl(`/edit/${network}`))}>
          <EditIcon />
        </IconButton>
      </TableCell>
    </TableRow>
  ));

  const closeDialog = () => history.push(relativeUrl(''));
  return (
    <div className={classes.paper}>
      <div className={classes.header}>
        <div />
        <NestedRouteLink to="/new">
          <Button variant="contained" color="primary">
            Add Network
          </Button>
        </NestedRouteLink>
      </div>
      <Paper className={classes.tableRoot} elevation={2}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Network ID</TableCell>
              <TableCell />
            </TableRow>
          </TableHead>
          <TableBody>{rows}</TableBody>
        </Table>
      </Paper>
      <Route
        path={relativePath('/new')}
        render={() => (
          <NetworkDialog
            onClose={closeDialog}
            onSave={networkID => {
              setNetworks([...networks, networkID]);
              enqueueSnackbar('Network created successfully', {
                variant: 'success',
              });
              closeDialog();
            }}
          />
        )}
      />
      <Route
        path={relativePath('/edit/:networkID')}
        render={() => (
          <NetworkDialog
            onClose={closeDialog}
            onSave={_ => {
              enqueueSnackbar('Network updated successfully', {
                variant: 'success',
              });
              closeDialog();
            }}
          />
        )}
      />
    </div>
  );
}

export default Networks;
