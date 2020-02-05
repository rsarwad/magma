/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ButtonProps} from '../Button';
import type {OptionProps} from './SelectMenu';

import * as React from 'react';
import BasePopoverTrigger from '../ContexualLayer/BasePopoverTrigger';
import Button from '../Button';
import SelectMenu from './SelectMenu';
import classNames from 'classnames';
import emptyFunction from '@fbcnms/util/emptyFunction';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  menu: {
    margin: '8px 0px',
  },
  menuDockRight: {
    position: 'absolute',
    right: '0',
  },
});

type Props<TValue> = {
  className?: string,
  menuDockRight?: boolean,
  children: React.Node,
  options: Array<OptionProps<TValue>>,
  onChange?: (value: TValue) => void | (() => void),
  leftIcon?: React$ComponentType<SvgIconExports>,
  rightIcon?: React$ComponentType<SvgIconExports>,
  searchable?: boolean,
  onOptionsFetchRequested?: (searchTerm: string) => void,
  ...ButtonProps,
};

const PopoverMenu = <TValue>({
  className,
  children,
  leftIcon,
  rightIcon,
  menuDockRight,
  onChange,
  variant,
  skin,
  disabled,
  ...selectMenuProps
}: Props<TValue>) => {
  const classes = useStyles();
  return (
    <BasePopoverTrigger
      popover={
        <SelectMenu
          {...selectMenuProps}
          onChange={onChange || emptyFunction}
          size="normal"
          className={classNames(classes.menu, {
            [classes.menuDockRight]: menuDockRight,
          })}
        />
      }>
      {(onShow, _onHide, contextRef) => (
        <Button
          onClick={onShow}
          ref={contextRef}
          variant={variant}
          skin={skin || 'regular'}
          disabled={disabled}
          className={className}
          leftIcon={leftIcon}
          rightIcon={rightIcon}>
          {children}
        </Button>
      )}
    </BasePopoverTrigger>
  );
};

export default PopoverMenu;
