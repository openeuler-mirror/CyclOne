import React from 'react';
import { Input, Button, Form, notification } from 'antd';
import { put } from 'common/xFetch2';
const FormItem = Form.Item;
import { formItemLayout, tailFormItemLayout } from 'common/enums';
import Popup from 'components/popup';

export default function action(options) {

  const onSubmit = (values) => {
    put(`/api/cloudboot/v1/devices/${options.records.sn}/oob/password`, values).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: '操作成功' });
      options.reload();
      Popup.close();
    });
  };

  Popup.open({
    title: `修改带外信息`,
    width: 600,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <WrapperForm
        {...options}
        onSubmit={onSubmit}
        onCancel={() => {
          Popup.close();
        }}
      />
    )
  });
}

class MyForm extends React.Component {
  handleSubmit = ev => {
    ev && ev.preventDefault();
    ev && ev.stopPropagation();
    this.props.form.validateFields({ force: true }, (error, values) => {
      const postValues = {

      };
      if (error) {
        notification.warning({
          message: '还有未填写完整的选项'
        });
        return;
      }
      this.props.onSubmit(values);
    });
  };

  render() {
    const { getFieldDecorator } = this.props.form;
    const { initialValue } = this.props;

    return <Form onSubmit={this.handleSubmit}>
      <FormItem {...formItemLayout} label='序列号' >
        {getFieldDecorator('sn', {
          initialValue: initialValue.sn
        })(
          <Input disabled={true} />
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='带外IP' >
        {getFieldDecorator('oob_ip', {
          initialValue: initialValue.oob_ip
        })(
          <Input disabled={true} />
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='带外用户名' >
        {getFieldDecorator('oob_user_name', {
          initialValue: initialValue.oob_user
        })(
          <Input/>
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='带外旧密码' >
        {getFieldDecorator('oob_password_old', {
          initialValue: initialValue.oob_password
        })(
          <Input/>
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='带外新密码' >
        {getFieldDecorator('oob_password_new', {
          rules: [
            {
              required: true
            }
          ]
        })(
          <Input />
        )}
      </FormItem>
      <FormItem {...tailFormItemLayout}>
        <div className='pull-right'>
          <Button onClick={() => this.props.onCancel()}>取消</Button>
          <Button
            style={{ marginLeft: 8 }}
            type='primary'
            htmlType='submit'
          >
            确认
          </Button>
        </div>
      </FormItem>
    </Form>;
  }
}
const WrapperForm = Form.create()(MyForm);
