import React from 'react';
import Popup from 'components/popup';
import { post, put } from 'common/xFetch2';
import { notification } from 'antd';
import MyForm from './form';
import moment from 'moment';

export default function action(options) {
  const record = options.records || {};
  const type = options.type;
  const typeMap = {
    _create: { name: '新增订单', method: post },
    _update: { name: '编辑订单', method: put }
  };
  const onSubmit = (values) => {
    //修改。
    if (type === '_update') {
      values.id = record.id;
    }
    if (values.pre_occupied_usites.length < values.amount) {
      return notification.error({ message: '机房可用机位数' + values.pre_occupied_usites.length + '小于订单总数' + values.amount });
    }
    if (values.pre_occupied_usites) {
      values.pre_occupied_usites = JSON.stringify(values.pre_occupied_usites);
    }
    if (values.expected_arrival_date) {
      values.expected_arrival_date = moment(values.expected_arrival_date).format('YYYY-MM-DD');
    }
    if (values.maintenance_service_date_begin) {
      values.maintenance_service_date_begin = moment(values.maintenance_service_date_begin).format('YYYY-MM-DD');
    }
    if (values.maintenance_service_date_end) {
      values.maintenance_service_date_end = moment(values.maintenance_service_date_end).format('YYYY-MM-DD');
    }

    typeMap[type].method('/api/cloudboot/v1/order', values).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: res.message });
      Popup.close();
      options.reload();
    });
  };

  Popup.open({
    title: `${typeMap[type].name}`,
    width: 700,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <div>
        <MyForm
          id={record.id}
          type={options.type}
          dispatch={options.dispatch}
          idc={options.idc}
          physicalArea={options.physicalArea}
          deviceCategory={options.deviceCategory}
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
