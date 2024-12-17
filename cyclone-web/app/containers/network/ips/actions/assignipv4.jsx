import React from 'react';
import Popup from 'components/popup';
import Task from './components/assign-ipv4';

export default function action(options) {

  Popup.open({
    title: '为选定设备分配一个IPv4地址',
    width: 1000,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <Task
        //id={options.initialValue[0].id}
        //initialValue={options.initialValue[0]}
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
