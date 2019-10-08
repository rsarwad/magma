/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React from 'react';

export type FormFieldContextValue = {
  disabled: boolean,
  hasError: boolean,
};

const FormFieldContext = React.createContext<FormFieldContextValue>({
  disabled: false,
  hasError: false,
});

export function useFormField() {
  return React.useContext(FormFieldContext);
}

export default FormFieldContext;
