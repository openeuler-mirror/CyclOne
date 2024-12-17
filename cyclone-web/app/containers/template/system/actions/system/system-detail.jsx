import React from 'react';
import Popup from 'components/popup/draw';
import { get } from 'common/xFetch2';
import MyForm from './form';
import { notification } from 'antd';

export default async function action(options) {

  //使用接口方便入口使用
  const res = await get(`/api/cloudboot/v1/system-templates/${options.record.id}`);
  if (res.status !== 'success') {
    return notification.error({ message: res.message });
  }
  const initialValue = res.content;

  Popup.open({
    title: '查看PXE配置详情',
    onCancel: () => {
      Popup.close();
    },
    content: (
      <div>
        <MyForm
          {...options}
          initialValue={initialValue}
          id={options.record.id}
          showSubmit={false}
          onCancel={() => {
            Popup.close();
          }}
        />
      </div>
    )
  });
}
