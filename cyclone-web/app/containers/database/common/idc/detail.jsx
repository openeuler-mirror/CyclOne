import React from 'react';
import Popup from 'components/popup/draw';
import FormGenerator from 'components/idcos-form/FormGenerator';
import { get } from 'common/xFetch2';
import { formSchema } from './formSchema';
import { notification } from 'antd';

export default async function action(options) {
  const record = options.records;
  const res = await get(`/api/cloudboot/v1/idcs/${record.id}`);
  if (res.status !== 'success') {
    return notification.error({ message: res.message });
  }
  let initialValue = { ...res.content };
  initialValue.first_server_room = (res.content.first_server_room || []).map(it => it && it.name);

  Popup.open({
    title: '数据中心详情',
    onCancel: () => {
      Popup.close();
    },
    content: (
      <div>
        <FormGenerator
          initialValue={initialValue}
          schema={formSchema(true)}
          showCancel={false}
          hideReset={true}
          showSubmit={false}
          onCancel={() => {
            Popup.close();
          }}
        />
      </div>
    )
  });
}
