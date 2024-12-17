import React from 'react';
import Popup from 'components/popup/draw';
import Form from './form';

export default function action(options) {
  Popup.open({
    title: options.record.title || '审批详情',
    onCancel: () => {
      Popup.close();
    },
    content: (
      <Form onCancel={() => Popup.close()} reload={options.reload} userInfo={options.userInfo} initialValue={options.record} approval_id={options.approval_id}>
      </Form>
    )
  });
}
