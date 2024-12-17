import React from 'react';
import Popup from 'components/popup';
import { post, put } from 'common/xFetch2';
import { notification } from 'antd';
import FormGenerator from 'components/idcos-form/FormGenerator';
import { formSchema } from '../../common/cabinet/formSchema';

export default function action(options) {
  const record = options.records;
  const type = options.type;
  if (type === '_update') {
    record.network_area_id = record.network_area.id;
  }
  const typeMap = {
    _create: '新增机架信息',
    _update: '编辑机架信息'
  };

  const onSubmit = (values) => {

    let method = post;
    let url = '/api/cloudboot/v1/server-cabinets';
    if (type === '_update') {
      method = put;
      url = `/api/cloudboot/v1/server-cabinets/${record.id}`;
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
          initialValue={record}
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
