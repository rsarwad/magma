/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import Checkbox from '@fbcnms/ui/components/design-system/Checkbox/Checkbox';
import classNames from 'classnames';
import useSideEffectCallback from './useSideEffectCallback';
import {
  HierarchyContextProvider,
  useHierarchyContext,
} from './HierarchyContext';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useEffect, useState} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
  },
  children: {
    marginLeft: '24px',
  },
}));

export const HIERARCHICAL_RELATION = {
  BI_DIRECTIONAL_INCLUSIVE: 'BI_DIRECTIONAL_INCLUSIVE',
  PARENT_REQUIRED: 'PARENT_REQUIRED',
};
type HierarchicalRelationType = $Keys<typeof HIERARCHICAL_RELATION>;

type SubTreeProps = $ReadOnly<{|
  title: React.Node,
  disabled?: ?boolean,
  onChange: (?boolean) => void,
  hierarchicalRelation: HierarchicalRelationType,
  children?: React.Node,
  className?: ?string,
|}>;

function CheckboxSubTree(props: SubTreeProps) {
  const {
    onChange,
    hierarchicalRelation,
    title,
    disabled,
    className,
    children,
  } = props;
  const classes = useStyles();

  const hierarchyContext = useHierarchyContext();

  const propagateValue = useSideEffectCallback(() => {
    const allChildren = hierarchyContext.childrenValues;
    const hasTrueChild = allChildren.findKey(child => child === true) != null;

    let aggregatedValue;

    if (hierarchicalRelation === HIERARCHICAL_RELATION.PARENT_REQUIRED) {
      if (hasTrueChild) {
        aggregatedValue = true;
      } else {
        return;
      }
    } else if (
      hierarchicalRelation === HIERARCHICAL_RELATION.BI_DIRECTIONAL_INCLUSIVE
    ) {
      const hasFalseChild =
        allChildren.findKey(child => child === false) != null;
      const hasNullChild = allChildren.findKey(child => child == null) != null;

      const childTypesCount =
        (hasFalseChild ? 1 : 0) +
        (hasTrueChild ? 1 : 0) +
        (hasNullChild ? 1 : 0);

      if (childTypesCount === 0) {
        if (hierarchyContext.parentValue == null) {
          aggregatedValue = false;
        } else {
          return;
        }
      } else if (childTypesCount > 1 || hasNullChild) {
        aggregatedValue = null;
      } else if (hasFalseChild) {
        aggregatedValue = false;
      } else if (hasTrueChild) {
        aggregatedValue = true;
      } else {
        return;
      }
    }

    if (aggregatedValue != hierarchyContext.parentValue) {
      onChange(aggregatedValue);
    }
  });

  useEffect(
    () => {
      propagateValue();
    }, // eslint-disable-next-line react-hooks/exhaustive-deps
    [hierarchyContext.childrenValues],
  );

  return (
    <div className={classNames(classes.root, className)}>
      <Checkbox
        checked={hierarchyContext.parentValue === true}
        disabled={disabled}
        indeterminate={
          hierarchyContext.parentValue == null &&
          !hierarchyContext.childrenValues.isEmpty()
        }
        title={title}
        onChange={status => onChange(status === 'checked')}
      />
      <div className={classes.children}>{children}</div>
    </div>
  );
}

type Props = $ReadOnly<{|
  ...SubTreeProps,
  id: string,
  value?: ?boolean,
  onChange?: ?(?boolean) => void,
  hierarchicalRelation?: ?HierarchicalRelationType,
|}>;

export default function HierarchicalCheckbox(props: Props) {
  const {
    id,
    value: propValue,
    hierarchicalRelation: hierarchicalRelationProp,
    onChange,
    ...subTreeProps
  } = props;
  const [value, setValue] = useState<?boolean>(null);
  const hierarchyContext = useHierarchyContext();

  const hierarchicalRelation =
    hierarchicalRelationProp ?? HIERARCHICAL_RELATION.BI_DIRECTIONAL_INCLUSIVE;

  const updateValueInContext = hierarchyContext.setChildValue;
  const isRegistered = hierarchyContext.childrenValues.has(id);

  const updateMyValue = useCallback(
    newValue => {
      setValue(newValue);
      updateValueInContext(id, newValue);
    },
    [updateValueInContext, id],
  );

  useEffect(() => {
    updateMyValue(propValue);
  }, [propValue, updateMyValue]);

  useEffect(
    () => {
      if (isRegistered) {
        if (
          hierarchyContext.parentValue != null &&
          hierarchyContext.parentValue != value
        ) {
          if (
            hierarchicalRelation ===
              HIERARCHICAL_RELATION.BI_DIRECTIONAL_INCLUSIVE ||
            (hierarchicalRelation === HIERARCHICAL_RELATION.PARENT_REQUIRED &&
              hierarchyContext.parentValue !== true)
          ) {
            updateMyValue(hierarchyContext.parentValue);
            if (onChange) {
              onChange(hierarchyContext.parentValue);
            }
          }
        }
      }
    }, // eslint-disable-next-line react-hooks/exhaustive-deps
    [isRegistered, hierarchyContext.parentValue, updateMyValue],
  );

  return (
    <HierarchyContextProvider parentValue={value}>
      <CheckboxSubTree
        {...subTreeProps}
        hierarchicalRelation={hierarchicalRelation}
        onChange={newValue => {
          updateMyValue(newValue);
          if (onChange) {
            onChange(newValue);
          }
        }}
      />
    </HierarchyContextProvider>
  );
}
