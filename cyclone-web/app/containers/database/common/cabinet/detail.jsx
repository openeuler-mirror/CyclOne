import React from 'react';
import Popup from 'components/popup/draw';
import FormGenerator from 'components/idcos-form/FormGenerator';
import { get } from 'common/xFetch2';
import { formSchema } from './formSchema';
import { notification } from 'antd';

export default async function action(options) {
  const record = options.records;
  const res = await get(`/api/cloudboot/v1/server-cabinets/${record.id}`);
  if (res.status !== 'success') {
    return notification.error({ message: res.message });
  }
  let initialValue = { ...res.content };

  //增加属性
  initialValue.network_area_id = res.content.network_area.id;
  initialValue.idc_name = res.content.idc.name;

  //查看详情的网络信息不加载网络列表（模拟）
  initialValue.network_area.server_room = res.content.server_room;
  const network = {
    loading: false,
    data: [initialValue.network_area]
  };


  Popup.open({
    title: '机架详情',
    onCancel: () => {
      Popup.close();
    },
    content: (
      <div>
        <FormGenerator
          initialValue={initialValue}
          schema={formSchema(true, { network })}
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
