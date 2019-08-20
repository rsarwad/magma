/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {BreadcrumbData} from './Breadcrumb.react';

import Breadcrumb from './Breadcrumb.react';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import Popover from '@material-ui/core/Popover';
import React, {useState} from 'react';
import Typography from '@material-ui/core/Typography';
import classNames from 'classnames';
import {gray8} from '@fbcnms/ui/theme/colors';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  breadcrumbs: {
    display: 'flex',
    alignItems: 'flex-start',
    minWidth: '200px',
  },
  moreIcon: {
    display: 'flex',
    alignItems: 'center',
  },
  moreIconButton: {
    cursor: 'pointer',
    color: gray8,
    '&:hover': {
      color: theme.palette.primary.main,
    },
  },
  collapsedBreadcrumbsList: {
    minWidth: '100px',
  },
  subtext: {
    fontSize: theme.typography.pxToRem(11),
    color: theme.palette.text.secondary,
    marginLeft: '8px',
  },
  largeText: {
    fontSize: '20px',
    lineHeight: '24px',
    fontWeight: 500,
  },
  smallText: {
    fontSize: '14px',
    lineHeight: '18px',
    fontWeight: 500,
  },
  slash: {
    color: gray8,
    margin: '0 6px',
  },
}));

const MAX_NUM_BREADCRUMBS = 3;

type Props = {
  breadcrumbs: Array<BreadcrumbData>,
  size?: 'default' | 'small' | 'large',
};

const Breadcrumbs = (props: Props) => {
  const {breadcrumbs, size} = props;
  const classes = useStyles();

  const [isBreadcrumbsMenuOpen, toggleBreadcrumbsMenuOpen] = useState(false);
  const [anchorEl, setAnchorEl] = React.useState(null);

  let collapsedBreadcrumbs = [];
  if (breadcrumbs.length > MAX_NUM_BREADCRUMBS) {
    collapsedBreadcrumbs = breadcrumbs.slice(1, breadcrumbs.length - 2);
  }
  const hasCollapsedBreadcrumbs = collapsedBreadcrumbs.length > 0;
  const startBreadcrumbs = hasCollapsedBreadcrumbs
    ? breadcrumbs.slice(0, 1)
    : [];
  const endBreadcrumbs = breadcrumbs.slice(
    collapsedBreadcrumbs.length + (hasCollapsedBreadcrumbs ? 1 : 0),
  );

  const textClass = size === 'small' ? classes.smallText : classes.largeText;

  return (
    <div className={classes.breadcrumbs}>
      {startBreadcrumbs.map(b => (
        <Breadcrumb
          key={b.id}
          data={b}
          isLastBreadcrumb={false}
          size={size}
          onClick={b.onClick}
        />
      ))}
      {hasCollapsedBreadcrumbs && (
        <div className={classes.moreIcon}>
          <Typography
            className={classNames([classes.moreIconButton, textClass])}
            onClick={e => {
              toggleBreadcrumbsMenuOpen(true);
              setAnchorEl(e.currentTarget);
            }}>
            {'...'}
          </Typography>
          <Typography className={classNames([classes.slash, textClass])}>
            {'/'}
          </Typography>
        </div>
      )}
      {endBreadcrumbs.map((b, i) => (
        <Breadcrumb
          key={b.id}
          data={b}
          isLastBreadcrumb={i === endBreadcrumbs.length - 1}
          size={size}
        />
      ))}
      <Popover
        open={isBreadcrumbsMenuOpen}
        anchorEl={anchorEl}
        onClose={() => {
          toggleBreadcrumbsMenuOpen(false);
          setAnchorEl(null);
        }}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'center',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'center',
        }}>
        <List className={classes.collapsedBreadcrumbsList}>
          {collapsedBreadcrumbs.map(b => (
            <ListItem
              key={`list_item_${b.id}`}
              button
              onClick={() => {
                b.onClick && b.onClick(b.id);
                toggleBreadcrumbsMenuOpen(false);
                setAnchorEl(null);
              }}>
              <Typography>{b.name}</Typography>
              <Typography className={classes.subtext}>{b.subtext}</Typography>
            </ListItem>
          ))}
        </List>
      </Popover>
    </div>
  );
};

Breadcrumbs.defaultProps = {
  size: 'default',
  showTypes: true,
};

export default Breadcrumbs;
