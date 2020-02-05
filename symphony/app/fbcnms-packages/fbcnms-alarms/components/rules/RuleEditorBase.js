/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 *
 * Base component for rule editors to render. Handles rendering common elements
 * such as receiver config and label editor.
 */

import * as React from 'react';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CardHeader from '@material-ui/core/CardHeader';
import Editor from '../common/Editor';
import Grid from '@material-ui/core/Grid';
import LabelsEditor from './LabelsEditor';
import RuleContext from './RuleContext';
import SelectReceiver from '../prometheus/Receivers/SelectReceiver';
import SelectRuleType from './SelectRuleType';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import {useAlertRuleReceiver, useForm} from '../hooks';
import type {ApiUtil} from '../AlarmsApi';
import type {Props as EditorProps} from '../common/Editor';
import type {Labels} from '../AlarmAPIType';

type Props = EditorProps & {
  onChange: (form: RuleEditorBaseFields) => void,
  initialState: ?RuleEditorBaseFields,
  apiUtil: ApiUtil,
};

// Fields for inputs which are standard between different rule editors
export type RuleEditorBaseFields = {
  name: string,
  description: string,
  labels: Labels,
};

export default function RuleEditorBase({
  isNew,
  apiUtil,
  children,
  initialState,
  onChange,
  onSave,
  ...props
}: Props) {
  const ruleContext = React.useContext(RuleContext);
  const {formState, handleInputChange, updateFormState} = useForm({
    initialState: initialState || defaultState(),
    onFormUpdated: onChange,
  });
  const {receiver, setReceiver, saveReceiver} = useAlertRuleReceiver({
    ruleName: formState?.name || '',
    apiUtil,
  });

  const handleSave = React.useCallback(async () => {
    await onSave();
    await saveReceiver();
  }, [saveReceiver, onSave]);

  const handleLabelsChange = React.useCallback(
    (labels: Labels) => {
      updateFormState({
        labels,
      });
    },
    [updateFormState],
  );

  return (
    <Editor
      {...props}
      title={formState?.name ?? 'New Alert Rule'}
      description="Configure rules to fire alerts"
      isNew={isNew}
      onSave={handleSave}>
      <Grid container item spacing={4}>
        <Grid container direction="column" item xs={7} spacing={4}>
          <Grid item>
            <Card>
              <CardHeader title="Details" />
              <CardContent>
                <Grid container direction="column" spacing={2}>
                  {isNew && (
                    <Grid item xs={6}>
                      <SelectRuleType
                        ruleMap={ruleContext.ruleMap}
                        value={ruleContext.ruleType}
                        onChange={ruleContext.selectRuleType}
                      />
                    </Grid>
                  )}
                  {children}
                </Grid>
              </CardContent>
            </Card>
          </Grid>
          <Grid item>
            <Card>
              <CardHeader title="Summary" />
              <CardContent>
                <Grid item>
                  <TextField
                    disabled={!isNew}
                    required
                    label="Rule Name"
                    placeholder="Ex: Link down"
                    fullWidth
                    value={formState.name}
                    onChange={handleInputChange(val => ({name: val}))}
                  />
                </Grid>
                <Grid item>
                  <TextField
                    disabled={!isNew}
                    required
                    label="Description"
                    placeholder="Ex: The link is down"
                    fullWidth
                    value={formState.description}
                    onChange={handleInputChange(val => ({description: val}))}
                  />
                </Grid>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
        <Grid container direction="column" item spacing={4} xs={5}>
          <Grid item>
            <Card>
              <CardHeader
                title={
                  <>
                    <Typography variant="h5" gutterBottom>
                      Notify
                    </Typography>
                    <Typography
                      color="textSecondary"
                      gutterBottom
                      variant="body2">
                      Team or service to notify when this alert fires
                    </Typography>
                  </>
                }
              />
              <CardContent>
                <SelectReceiver
                  label="Send Notification To:"
                  fullWidth
                  apiUtil={apiUtil}
                  receiver={receiver}
                  onChange={setReceiver}
                />
              </CardContent>
            </Card>
          </Grid>
          <Grid item>
            <LabelsEditor
              labels={formState.labels}
              onChange={handleLabelsChange}
            />
          </Grid>
        </Grid>
      </Grid>
    </Editor>
  );
}

function defaultState(): RuleEditorBaseFields {
  return {
    name: '',
    description: '',
    labels: {},
  };
}
