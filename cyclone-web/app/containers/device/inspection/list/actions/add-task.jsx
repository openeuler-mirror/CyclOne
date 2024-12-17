import React from 'react';
import Popup from 'components/popup';
import Task from './components/task';
import { post } from 'common/xFetch2';

export default function action(options) {

  Popup.open({
    title: '新建巡检任务',
    width: 1000,
    onCancel: () => {
      Popup.close();
    },

    content: (
      <Task
        onCancel={() => {
          Popup.close();
        }}
        onSuccess={() => {
          Popup.close();
          options.reload();
        }}
      />

    )
  });

}
