import React from 'react';
import Popup from 'components/popup';
import { post, put } from 'common/xFetch2';
import { notification } from 'antd';
import MyForm from '../../common/usite/form';

export default function action(options) {
  const record = options.records || {};
  const type = options.type;
  const typeMap = {
    _create: '新增机位信息',
    _update: '编辑机位信息'
  };

  const onSubmit = (values) => {

    let method = post;
    let url = '/api/cloudboot/v1/server-usites';
    if (type === '_update') {
      method = put;
      url = `/api/cloudboot/v1/server-usites/${record.id}`;
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
        <MyForm
          type={options.type}
          dispatch={options.dispatch}
          room={options.room}
          cabinet={options.cabinet}
          server_cabinet_id={(record.server_cabinet || {}).id}
          id={record.id}
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
