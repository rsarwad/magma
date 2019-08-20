/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import DialogError from '../components/DialogError.react';
import React from 'react';
import {storiesOf} from '@storybook/react';

storiesOf('DialogError', module).add('default', () => (
  <DialogError message={'This is an error message!'} />
));
