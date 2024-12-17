import React from 'react';
import { Form, Radio, Button, notification } from 'antd';
const FormItem = Form.Item;
const RadioGroup = Radio.Group;
import DeviceTable from 'containers/device/common/device';
import { put } from 'common/xFetch2';
import { NETWORK_IPS_SCOPE, getSearchList } from "common/enums";

export const formItemLayout = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 3 }
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 20 }
  }
};
export const tailFormItemLayout = {
  wrapperCol: {
    xs: {
      span: 24,
      offset: 0
    },
    sm: {
      span: 11,
      offset: 12
    }
  }
};
class AssignIPv6Form extends React.Component {
  handleSubmit = (e) => {
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (err) {
        return notification.error({ message: '还有未填写的项' });
      }
      const sn = values.sn ? values.sn[0].sn : '';
      put('/api/cloudboot/v1/ips/assignsv6', { sn: sn, scope: values.scope }).then(res => {
        if (res.status !== 'success') {
          return notification.error({ message: res.message });
        }
        notification.success({ message: res.message });
        this.props.onSuccess();
      });
    });
  };

  render() {
    const { getFieldDecorator } = this.props.form;
    //const ip_network = this.props.initialValue.ip_network || {};
    //const scope = ip_network.category.indexOf('intranet') > -1 ? 'intranet' : 'extranet';
    return (
      <Form onSubmit={this.handleSubmit}>
        <FormItem
          {...formItemLayout}
          label='IP作用范围'
        >
          {getFieldDecorator('scope', {
            rules: [{
              required: true, message: '请选择内外网'
            }]
          })(
            <RadioGroup>
              {getSearchList(NETWORK_IPS_SCOPE).map(it => <Radio value={it.value}>{it.label}</Radio>)}
            </RadioGroup>
          )}
        </FormItem>
        <FormItem
          {...formItemLayout}
          label='关联设备'
        >
          {getFieldDecorator('sn', {
            rules: [{
              required: true, message: '请选择关联设备'
            }]
          })(
            <DeviceTable checkType='radio' hideButton={true} form={this.props.form} />
          )}
        </FormItem>
        <FormItem
          {...tailFormItemLayout}
        >
          <div className='pull-right'>
            <Button onClick={() => this.props.onCancel()}>取消</Button>
            <Button
              style={{ marginLeft: 8 }}
              type='primary'
              htmlType='submit'
            >
                确定
            </Button>
          </div>
        </FormItem>
      </Form>
    );
  }
}

const WrappedAssignIPv6Form = Form.create()(AssignIPv6Form);

export default WrappedAssignIPv6Form;
