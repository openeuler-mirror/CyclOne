import React from 'react';
import Popup from 'components/popup/draw';
import MyForm from './form';

export default async function action(options) {

  Popup.open({
    title: '网络设备详情',
    onCancel: () => {
      Popup.close();
    },
    content: (
      <div>
        <MyForm
          type={options.type}
          room={options.room}
          cabinet={options.cabinet}
          id={options.records.id}
          showSubmit={false}
          onCancel={() => {
            Popup.close();
          }}
        />
      </div>
    )
  });
}
