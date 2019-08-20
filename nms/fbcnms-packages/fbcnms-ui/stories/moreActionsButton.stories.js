/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import MoreActionsButton from '../components/MoreActionsButton.react';
import React from 'react';
import {storiesOf} from '@storybook/react';

storiesOf('MoreActionsButton', module).add('string', () => (
  <MoreActionsButton
    variant="primary"
    items={[
      {name: 'Item 1', onClick: () => window.alert('clicked item #1')},
      {name: 'Item 2', onClick: () => window.alert('clicked item #2')},
      {name: 'Item 3', onClick: () => window.alert('clicked item #3')},
    ]}
  />
));
