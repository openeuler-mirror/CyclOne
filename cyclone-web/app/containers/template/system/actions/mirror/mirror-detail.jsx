import React from 'react';
import Popup from 'components/popup/draw';
import { get } from 'common/xFetch2';
import MyForm from './form';
import { notification } from 'antd';

export default async function action(options) {
  //使用接口方便其他入口使用
  const res = await get(`/api/cloudboot/v1/image-templates/${options.record.id}`);
  if (res.status !== 'success') {
    return notification.error({ message: res.message });
  }
  const initialValue = res.content;

  Popup.open({
    title: '查看镜像配置详情',
    onCancel: () => {
      Popup.close();
    },
    content: (
      <div>
        <MyForm
          initialValue={initialValue}
          {...options}
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
