import React from 'react';
import Popup from 'components/popup';
import FormGenerator from 'components/idcos-form/FormGenerator';

import { put } from 'common/xFetch2';
import { notification } from 'antd';

export default function action(options) {
  let newPassword;

  const onSubmit = (values) => {

    put('/api/cloudboot/v1/users/password', values).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: res.message });
      Popup.close();
      options.reload();
    });
  };
  const formSchema = {
    name: 'form',
    id: 'form',
    elements: [
      {
        id: 'old_password',
        name: 'old_password',
        label: '旧密码',
        type: 'TextInput',
        inputType: 'password',
        rules: [{
          required: true,
          message: '请输入旧密码'
        }]
      },
      {
        id: 'new_password',
        name: 'new_password',
        label: '新密码',
        type: 'TextInput',
        inputType: 'password',
        placeholder: '必须包含大小写字母，数字和特殊字符',
        rules: [{
          required: true,
          message: '请输入新密码'
        }, {
          validator: (rule, value, callback) => {
            newPassword = value;
            if (value) {
              if (!(/[a-z]+/.test(value) && /[A-Z]+/.test(value) && /[\d]/.test(value) && /[~!@#\$%\^&\*\(\)_\-\+=\|\\\{\}\[\]\"';:\/\?\.>,<]/.test(value))) {
                callback('必须包含大小写字母，数字和特殊字符');
              }
            }
            callback();
          }
        }]
      }, {
        id: '',
        name: '',
        label: '确认新密码',
        type: 'TextInput',
        inputType: 'password',
        rules: [{
          required: true,
          message: '请确认新密码'
        }, {
          validator: (rule, value, callback) => {
            if (value !== newPassword) {
              callback([
                new Error(
                  '两次输入不一致'
                )
              ]);
            }
            callback();
          }
        }]
      }]
  };

  Popup.open({
    title: '修改密码',
    width: 700,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <div>
        <FormGenerator
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
