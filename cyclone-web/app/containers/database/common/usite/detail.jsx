import React from 'react';
import Popup from 'components/popup/draw';
import MyForm from "./form";
import { post, put } from 'common/xFetch2';

export default function action(options) {
  Popup.open({
    title: '机位详情',
    onCancel: () => {
      Popup.close();
    },
    content: (
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
    )
  });
}
