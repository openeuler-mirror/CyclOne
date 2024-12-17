import React from 'react';
import Popup from 'components/popup';
import { post, put } from 'common/xFetch2';
import { notification } from 'antd';
import MyForm from './form';

export default function action(options) {
  const record = options.records || {};
  const type = options.type;
  const typeMap = {
    _create: '新增网络设备信息',
    _update: '编辑网络设备信息'
  };

  const onSubmit = (values) => {

    //若id=0，则新增。若id>0，则修改。
    // if (type === '_create') {
    //   values.id = 0;
    // } else {
    //   values.id = record.id;
    // }

    post('/api/cloudboot/v1/network/devices', values).then(res => {
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
          idc={options.idc}
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
