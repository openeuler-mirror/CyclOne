import React from 'react';
import Popup from 'components/popup';
import { put } from 'common/xFetch2';
import { notification } from 'antd';
import FormGenerator from 'components/idcos-form/FormGenerator';
import { USITE_STATUS, getSearchList } from 'common/enums';

export default function action(options) {
  const record = options.records;
  const ids = record.map(item => item.id);

  const formSchema = {
    name: 'form',
    id: 'form',
    elements: [
      {
        id: 'status',
        name: 'status',
        label: '状态',
        type: 'Select',
        options: getSearchList(USITE_STATUS),
        rules: [
          {
            required: true,
            message: '请选择机位状态'
          }
        ]
      }
    ]
  };

  const onSubmit = (values) => {

    put('/api/cloudboot/v1/server-usites/status', { status: values.status, ids: ids }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: res.message });
      Popup.close();
      options.reload();
    });
  };

  Popup.open({
    title: '更新机位状态',
    onCancel: () => {
      Popup.close();
    },
    content: (
      <div>
        <FormGenerator
          initialValue={record}
          schema={formSchema}
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
