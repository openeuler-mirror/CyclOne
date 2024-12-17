import React from 'react';
import { Input, Select, Button, Form, notification } from 'antd';
import { post, put } from 'common/xFetch2';
const FormItem = Form.Item;
import { formItemLayout, tailFormItemLayout } from 'common/enums';
import Popup from 'components/popup';
const { Option } = Select;

export default function action(options) {
  const type = options.type;
  const typeMap = {
    _create: { name: '新增参数规则', method: post, url: '/api/cloudboot/v1/device-setting-rules' },
    _update: { name: '修改参数规则', method: put, url: `/api/cloudboot/v1/device-setting-rules` }
  };
  const onSubmit = (values) => {
    if (type == '_update') {
      values.id = options.records.id;
    }
    typeMap[type].method(typeMap[type].url, values).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: '操作成功' });
      options.reload();
      Popup.close();
    });
  };

  Popup.open({
    title: `${typeMap[type].name}`,
    width: 600,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <WrapperForm
        initialValue={options.records}
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
    const initialValue = this.props.initialValue;
    return <Form onSubmit={this.handleSubmit}>
      <FormItem {...formItemLayout} label='规则前件' >
        {getFieldDecorator('condition', {
          initialValue: initialValue.condition,
          rules: [
            {
              required: true,
              message: '请输入参数规则前件'
            }
          ]
        })(
          <Input.TextArea rows={3}/>
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='规则推论' >
        {getFieldDecorator('action', {
          initialValue: initialValue.action,
          rules: [
            {
              required: true,
              message: '请输入规则推论'
            }
          ]
        })(
          <Input />
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='规则类别' >
        {getFieldDecorator('rule_category', {
          initialValue: initialValue.rule_category,
          rules: [
            {
              required: true,
              message: '请输入规则类别'
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
