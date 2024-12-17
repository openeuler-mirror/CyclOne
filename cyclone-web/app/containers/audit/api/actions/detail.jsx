import React from 'react';
import Popup from 'components/popup/draw';
import Form from './form';

export default function action(options) {

  Popup.open({
    title: '操作详情',
    onCancel: () => {
      Popup.close();
    },
    content: (
      <Form initialValue={options.initialValue}>
      </Form>
    )
  });
}
