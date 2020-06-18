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
import Button from '../Button';
import Checkbox from '../Checkbox/Checkbox';
import DialogRoot from './DialogRoot';
import IconButton from '../IconButton';
import Strings from '@fbcnms/strings/Strings';
import Text from '../Text';
import {CloseIcon} from '../Icons';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    padding: '24px',
    minWidth: '480px',
    minHeight: '210px',
    maxWidth: '600px',
    maxHeight: '600px',
    display: 'flex',
    flexDirection: 'column',
    boxSizing: 'border-box',
  },
  titleContainer: {
    display: 'flex',
    flexDirection: 'row',
    marginBottom: '16px',
  },
  titleText: {
    flexGrow: 1,
    maxWidth: '560px',
    overflow: 'hidden',
    marginRight: '16px',
  },
  checkboxContainer: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    flexGrow: 1,
  },
  content: {
    flexGrow: 1,
    marginBottom: '40px',
  },
  footer: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'flex-end',
  },
  cancelButton: {
    marginRight: '8px',
  },
}));

export type DialogSkin = 'primary' | 'red';

export type MessageDialogProps = $ReadOnly<{|
  title: React.Node,
  message: React.Node,
  verificationCheckbox?: ?{
    label: React.Node,
    isMandatory?: ?boolean,
  },
  cancelLabel?: React.Node,
  confirmLabel?: React.Node,
  skin?: DialogSkin,
  onCancel?: () => void,
  onClose: () => void,
  onConfirm?: (?boolean) => void,
|}>;

export type MessageDialogComponentProps = $ReadOnly<{|
  ...MessageDialogProps,
  hidden?: boolean,
|}>;

const MessageDialog = ({
  title,
  message,
  onClose,
  verificationCheckbox,
  cancelLabel = Strings.common.cancelButton,
  confirmLabel = Strings.common.okButton,
  onCancel,
  onConfirm,
  hidden,
  skin = 'primary',
}: MessageDialogComponentProps) => {
  const classes = useStyles();
  const [checkboxChecked, setCheckboxChecked] = useState(false);

  useEffect(() => {
    setCheckboxChecked(false);
  }, [hidden]);

  return (
    <DialogRoot className={classes.root} onClose={onClose} hidden={hidden}>
      <div className={classes.titleContainer}>
        <Text className={classes.titleText} weight="medium">
          {title}
        </Text>
        <IconButton skin="gray" icon={CloseIcon} onClick={onClose} />
      </div>
      <div className={classes.content}>
        <Text>{message}</Text>
      </div>
      <div className={classes.footer}>
        {verificationCheckbox && (
          <div className={classes.checkboxContainer}>
            <Checkbox
              checked={checkboxChecked}
              title={verificationCheckbox.label}
              onChange={selection =>
                setCheckboxChecked(selection === 'checked' ? true : false)
              }
            />
          </div>
        )}
        {cancelLabel && (
          <Button
            skin="gray"
            onClick={onCancel}
            className={classes.cancelButton}>
            {cancelLabel}
          </Button>
        )}
        {confirmLabel && (
          <Button
            onClick={() =>
              onConfirm &&
              onConfirm(verificationCheckbox == null ? null : checkboxChecked)
            }
            skin={skin}
            disabled={
              verificationCheckbox?.isMandatory === true && !checkboxChecked
            }>
            {confirmLabel}
          </Button>
        )}
      </div>
    </DialogRoot>
  );
};

export default MessageDialog;
