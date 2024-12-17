import React from 'react';
import { notification } from 'antd';
import { post } from 'common/xFetch2';
import Popup from 'components/popup';
import MyForm from './form';

export default function action(options) {

  const onSubmit = (values) => {
    values.need_extranet_ip = values.need_extranet_ip ? 'yes' : 'no';
    values.need_intranet_ip = values.need_intranet_ip ? 'yes' : 'no';
    
    if (values.maintenance_service_date_begin) {
      values.maintenance_service_date_begin = moment(values.maintenance_service_date_begin).format('YYYY-MM-DD');
    }

    post('/api/cloudboot/v1/special-device', values).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: '操作成功' });
      options.reload(); 
      Popup.close();
    });
  };

  Popup.open({
    title: `新增特殊设备`,
    width: 700,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <MyForm
        {...options}
        //initialValue={options.records}
        username={options.username}
        onSubmit={onSubmit}
        onCancel={() => {
          Popup.close();
        }}
      />
    )
  });
}