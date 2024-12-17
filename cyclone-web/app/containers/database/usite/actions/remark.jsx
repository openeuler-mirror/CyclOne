import React from 'react';
import Popup from 'components/popup';
import { put } from 'common/xFetch2';
import { notification } from 'antd';
import FormGenerator from 'components/idcos-form/FormGenerator';

export default function action(options) {
  const record = options.records;
  const ids = record.map(item => item.id);

  const formSchema = {
    name: 'form',
    id: 'form',
    elements: [
      {
        id: 'remark',
        name: 'remark',
        label: '备注',
        type: 'Textarea',
        rules: [
          {
            required: true,
            message: '请输入备注内容'
          }
        ]
      }
    ]
  };

  const onSubmit = (values) => {

    put('/api/cloudboot/v1/server-usites/remark', { remark: values.remark, ids: ids }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: res.message });
      Popup.close();
      options.reload();
    });
  };

  Popup.open({
    title: '更新备注信息',
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
