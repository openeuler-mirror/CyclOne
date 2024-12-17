import React from 'react';
import Popup from 'components/popup/draw';
import FormGenerator from 'components/idcos-form/FormGenerator';
import { get } from 'common/xFetch2';
import { formSchema } from './formSchema';
import { notification } from 'antd';

export default async function action(options) {
  const record = options.records;
  const res = await get(`/api/cloudboot/v1/network-areas/${record.id}`);
  if (res.status !== 'success') {
    return notification.error({ message: res.message });
  }
  let initialValue = { ...res.content };

  //增加和修改属性
  initialValue.server_room_id = res.content.server_room.id;
  initialValue.physical_area = res.content.physical_area.map(it => it && it.name);
  initialValue.idc_name = res.content.idc.name;

  //查看详情的机房信息不加载机房列表（模拟）
  const room = {
    loading: false,
    data: [initialValue.server_room]
  };

  Popup.open({
    title: '网络区域详情',
    onCancel: () => {
      Popup.close();
    },
    content: (
      <div>
        <FormGenerator
          initialValue={initialValue}
          schema={formSchema(true, { room })}
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
