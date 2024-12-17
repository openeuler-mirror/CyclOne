import React from 'react';
import Popup from 'components/popup';
import { post } from 'common/xFetch2';
import { notification } from 'antd';
import FormGenerator from 'components/idcos-form/FormGenerator';

const formSchema = (disabled) => {
  let elements = [
    {
      id: 'number',
      name: 'number',
      label: '货架编号',
      type: 'TextInput',
      disabled: disabled,
      rules: [
        {
          required: true,
          message: '请填写货架编号'
        }
      ]
    },
    {
      id: 'remark',
      name: 'remark',
      label: '备注',
      type: 'Textarea',
      disabled: disabled
    }
  ];
  return {
    name: 'form',
    id: 'form',
    elements: elements
  };
};
export default function action(options) {
  const onSubmit = (values) => {
    values.store_room_id = Number(options.id);
    post('/api/cloudboot/v1/virtual-cabinet', values).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: res.message });
      Popup.close();
      options.reload();
    });
  };

  Popup.open({
    title: `新增虚拟货架`,
    width: 600,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <FormGenerator
        schema={formSchema(false)}
        showCancel={true}
        hideReset={true}
        onSubmit={(values) => onSubmit(values)}
        onCancel={() => {
          Popup.close();
        }}
      />
    )
  });
}
