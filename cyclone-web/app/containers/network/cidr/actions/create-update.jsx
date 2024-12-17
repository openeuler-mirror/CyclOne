import React from 'react';
import Popup from 'components/popup';
import FormGenerator from 'components/idcos-form/FormGenerator';
import { post, put } from 'common/xFetch2';
import { notification } from 'antd';
import { formSchema } from './formSchema';


export default function action(options) {
  const record = options.records;
  const type = options.type;

  let initialValue = { ...record };

  if (type === '_update') {
    initialValue.server_room_id = record.server_room.id;
    initialValue.idc_name = record.idc.name;
    //switchs回显
    initialValue.switchs = (record.switchs || []).map(it => it.fixed_asset_number);
  }

  const typeMap = {
    _create: '新增IP网段',
    _update: '编辑IP网段'
  };

  const onSubmit = (values) => {

    let method = post;
    let url = '/api/cloudboot/v1/ip-networks';
    if (type === '_update') {
      method = put;
      url = `/api/cloudboot/v1/ip-networks/${record.id}`;
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
          schema={formSchema(false, options, type)}
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
