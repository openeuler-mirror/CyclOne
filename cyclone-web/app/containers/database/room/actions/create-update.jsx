import React from 'react';
import Popup from 'components/popup';
import { post, put } from 'common/xFetch2';
import { notification } from 'antd';
import MyForm from '../../common/room/form';


export default function action(options) {
  const record = options.records || {};
  const type = options.type;
  const typeMap = {
    _create: '新增机房信息',
    _update: '编辑机房信息'
  };

  const onSubmit = (values) => {

    let method = post;
    let url = '/api/cloudboot/v1/server-rooms';
    if (type === '_update') {
      method = put;
      url = `/api/cloudboot/v1/server-rooms/${record.id}`;
    }
    values.idc_id = JSON.parse(values.idc_data || '{}').id;
    delete values.idc_data;
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
        <MyForm
          id={record.id}
          idc={options.idc}
          type={type}
          showSubmit={true}
          onSubmit={(values) => onSubmit(values)}
          onCancel={() => {
            Popup.close();
          }}
        />
      </div>
    )
  });
}
