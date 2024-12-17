import React from 'react';
import Popup from 'components/popup';
import FormGenerator from 'components/idcos-form/FormGenerator';
import { post, put } from 'common/xFetch2';
import { notification } from 'antd';
import { formSchema } from '../../common/idc/formSchema';


export default function action(options) {
  const record = options.records;

  const type = options.type;
  const typeMap = {
    _create: '新增数据中心',
    _update: '编辑数据中心'
  };

  let initialValue = { ...record };

  if (type === '_update') {
    //一级机房数据回显
    initialValue.first_server_room = (record.first_server_room || []).map(it => it.name);
  }

  const onSubmit = (values) => {

    let method = post;
    let url = '/api/cloudboot/v1/idcs';
    if (type === '_update') {
      method = put;
      url = `/api/cloudboot/v1/idcs/${record.id}`;
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
    width: 600,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <div>
        <FormGenerator
          initialValue={initialValue}
          schema={formSchema(false)}
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
