import React from 'react';
import { Input, Icon, Tooltip, Button, Form, notification, InputNumber, Radio } from 'antd';
import { post, put } from 'common/xFetch2';
const FormItem = Form.Item;
import { formItemLayout, tailFormItemLayout, getSearchList, BUILTIN } from 'common/enums';
import Popup from 'components/popup';
const RadioGroup = Radio.Group;

export default function action(options) {
  const type = options.type;
  const typeMap = {
    _create: { name: '新增设备类型', method: post, url: '/api/cloudboot/v1/device-category' },
    _update: { name: '修改设备类型', method: put, url: `/api/cloudboot/v1/device-category` }
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
      <FormItem {...formItemLayout} label='设备类型' help='相同硬件配置，设备类型唯一'>
        {getFieldDecorator('category', {
          initialValue: initialValue.category,
          rules: [
            {
              required: true,
              message: '请输入设备类型'
            }
          ]
        })(
          <Input />
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='硬件配置' >
        {getFieldDecorator('hardware', {
          initialValue: initialValue.hardware,
          rules: [
            {
              required: true,
              message: '请输入硬件配置'
            }
          ]
        })(
          <Input />
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='处理器生产商' help='如：Intel(R) Corporation\HiSilicon'>
        {getFieldDecorator('central_processor_manufacture', {
          initialValue: initialValue.central_processor_manufacture,
          rules: [
            {
              required: true,
              message: '请输入处理器生产商'
            }
          ]
        })(
          <Input />
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='处理器架构' help='如: x86_64\aarch64'>
        {getFieldDecorator('central_processor_arch', {
          initialValue: initialValue.central_processor_arch,
          rules: [
            {
              required: true,
              message: '请输入处理器架构'
            }
          ]
        })(
          <Input />
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='功率' >
        {getFieldDecorator('power', {
          initialValue: initialValue.power,
          rules: [
            {
              required: true,
              message: '请输入功率'
            }
          ]
        })(
          <Input />
        )}
      </FormItem>            
      <FormItem {...formItemLayout} label='设备U数' >
            {getFieldDecorator('unit', {
              initialValue: initialValue.unit,
              rules: [
                {
                  required: true,
                  message: '请输入设备U数'
                }
              ]
            })(<InputNumber style={{ width: '100%' }} />)}
      </FormItem>
      <FormItem {...formItemLayout} label='是否金融信创生态产品'>
            {getFieldDecorator('is_fiti_eco_product', {
              initialValue: initialValue.is_fiti_eco_product,
              rules: [{
                required: true, message: '请选择是否金融信创生态产品'
              }]
            })(
              <RadioGroup>
                {getSearchList(BUILTIN).map(it => <Radio value={it.value}>{it.label}</Radio>)}
              </RadioGroup>
            )}
          </FormItem>
      <FormItem {...formItemLayout} label='备注' >
        {getFieldDecorator('remark', {
          initialValue: initialValue.remark
        })(
          <Input.TextArea />
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
