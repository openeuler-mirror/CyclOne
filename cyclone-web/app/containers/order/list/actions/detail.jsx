import React from 'react';
import Popup from 'components/popup/draw';
import MyForm from './form';

export default async function action(options) {

  Popup.open({
    title: '订单详情',
    onCancel: () => {
      Popup.close();
    },
    content: (
      <div>
        <MyForm
          type={options.type}
          id={options.records.id}
          physicalArea={options.physicalArea}
          deviceCategory={options.deviceCategory}
          showSubmit={false}
          onCancel={() => {
            Popup.close();
          }}
        />
      </div>
    )
  });
}
