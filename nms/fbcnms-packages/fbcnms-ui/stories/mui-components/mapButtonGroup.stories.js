/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
declare module '@material-ui/icons/AccessAlarms' {
  declare module.exports: React$ComponentType<SvgIconExports>;
}

import ListAltIcon from '@material-ui/icons/ListAlt';
import MapButtonGroup from '@fbcnms/ui/components/map/MapButtonGroup';
import MapIcon from '@material-ui/icons/Map';
import React from 'react';
import Typography from '@material-ui/core/Typography';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles({
  text: {
    fontSize: '12px',
    LineHeight: '14px',
  },
});

const AddMapButtonGroup = () => {
  return (
    <MapButtonGroup
      onIconClicked={() => {}}
      buttons={[
        {item: <ListAltIcon />, id: 'list'},
        {item: <MapIcon />, id: 'map'},
      ]}
    />
  );
};
const AddThreeMapButton = () => {
  return (
    <MapButtonGroup
      onIconClicked={() => {}}
      buttons={[
        {item: <ListAltIcon />, id: 'list'},
        {item: <MapIcon />, id: 'map'},
        {item: <MapIcon />, id: 'map2'},
      ]}
    />
  );
};
const AddTextButton = () => {
  const classes = useStyles();
  return (
    <MapButtonGroup
      onIconClicked={() => {}}
      buttons={[
        {
          item: <Typography className={classes.text}>Status</Typography>,
          id: 'status',
        },
        {
          item: <Typography className={classes.text}>Technician</Typography>,
          id: 'Technician',
        },
      ]}
    />
  );
};

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/MapButtonGroup`, module)
  .add('two', () => {
    return <AddMapButtonGroup />;
  })
  .add('three', () => {
    return <AddThreeMapButton />;
  })
  .add('text', () => {
    return <AddTextButton />;
  });
