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
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import classNames from 'classnames';
import {gray8} from '@fbcnms/ui/theme/colors';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
  },
  upperSection: {
    display: 'flex',
    flexDirection: 'column',
  },
  subtext: {
    fontSize: theme.typography.pxToRem(13),
  },
  slash: {
    color: gray8,
    margin: '0 6px',
  },
  breadcrumbName: {
    whiteSpace: 'nowrap',
    fontWeight: 500,
    color: theme.palette.blueGrayDark,
    cursor: 'pointer',
  },
  parentBreadcrumb: {
    color: gray8,
    '&:hover': {
      color: theme.palette.primary.main,
    },
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    maxWidth: '100px',
    display: 'inline-block',
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
}));

export type BreadcrumbData = {
  id: string,
  name: string,
  subtext?: ?string | React.Node,
  onClick?: ?(id: string) => void,
};

type Props = {
  data: BreadcrumbData,
  isLastBreadcrumb: boolean,
  size?: 'default' | 'small' | 'large',
};

const Breadcrumb = (props: Props) => {
  const {data, isLastBreadcrumb, size} = props;
  const {id, name, subtext, onClick} = data;
  const classes = useStyles();
  const textClass = size === 'small' ? classes.smallText : classes.largeText;
  return (
    <div key={id} className={classes.root}>
      <div className={classes.upperSection}>
        <Tooltip
          placement="top"
          title={
            typeof subtext === 'string' ? (
              <Typography className={classes.subtext}>{subtext}</Typography>
            ) : (
              subtext ?? ''
            )
          }>
          <Typography
            className={classNames({
              [classes.breadcrumbName]: true,
              [classes.parentBreadcrumb]: !isLastBreadcrumb,
              [textClass]: true,
            })}
            onClick={() => onClick && onClick(id)}>
            {name}
          </Typography>
        </Tooltip>
      </div>
      {!isLastBreadcrumb && (
        <Typography className={classNames([classes.slash, textClass])}>
          {'/'}
        </Typography>
      )}
    </div>
  );
};

Breadcrumb.defaultProps = {
  size: 'default',
};

export default Breadcrumb;
