import React from 'react';
import Popup from 'components/popup';
import FormGenerator from 'components/idcos-form/FormGenerator';
import { post, put } from 'common/xFetch2';
import { notification } from 'antd';
import { formSchema } from '../../common/network/formSchema';


export default function action(options) {
  const record = options.records;
  const type = options.type;

  let initialValue = { ...record };

  if (type === '_update') {
    initialValue.server_room_id = record.server_room.id;
    initialValue.physical_area = record.physical_area.map(it => it && it.name);
    initialValue.idc_name = record.idc.name;
  }
  const typeMap = {
    _create: '新增网络区域',
    _update: '编辑网络区域'
  };

  const onSubmit = (values) => {

    let method = post;
    let url = '/api/cloudboot/v1/network-areas';
    if (type === '_update') {
      method = put;
      url = `/api/cloudboot/v1/network-areas/${record.id}`;
    }

    method(url, values).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: res.message });
      Popup.close();
      options.reload();
    });
  };

  Popup.open({
    title: `${typeMap[type]}`,
    width: 700,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <div>
        <FormGenerator
          initialValue={initialValue}
          schema={formSchema(false, options)}
          showCancel={true}
          hideReset={true}
          onSubmit={(values) => onSubmit(values)}
          onCancel={() => {
            Popup.close();
          }}
        />
      </div>
    )
  });
}
