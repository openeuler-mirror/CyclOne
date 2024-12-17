import React from 'react';
import { Form, Input, Radio, Button, notification } from 'antd';
const FormItem = Form.Item;
const RadioGroup = Radio.Group;
import DeviceTable from 'containers/device/common/device';
import { post } from 'common/xFetch2';
import Crontab from 'components/crontab';
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
class RegistrationForm extends React.Component {
  state = {
    crontabUi: {}
  };

  handleSubmit = (e) => {
    e.preventDefault();
    this.props.form.validateFields((err, values) => {

      //指定类型
      if (values.type === 'define') {
        values.cron = this.state.crontabUi.crontabExpression;
        values.cron_render = JSON.stringify(this.state.crontabUi);
      }

      if (err) {
        return notification.error({ message: '还有未填写的项' });
      }

      const selectedRows = values.sn || [];
      const sns = selectedRows.map(s => s.sn);
      post('/api/cloudboot/v1/jobs/inspections', { sn: sns, rate: values.rate, cron: values.cron, cron_render: values.cron_render }).then(res => {
        if (res.status !== 'success') {
          return notification.error({ message: res.message });
        }
        notification.success({ message: res.message });
        this.props.onSuccess();
      });
    });
  };

  render() {
    const { getFieldDecorator, getFieldValue } = this.props.form;
    const rate = getFieldValue('rate');
    return (
      <Form onSubmit={this.handleSubmit}>
        <FormItem
          {...formItemLayout}
          label='任务类型'
        >
          {getFieldDecorator('rate', {
            initialValue: 'immediately',
            rules: [{
              required: true,
              message: '请选择任务执行类型'
            }]
          })(
            <RadioGroup>
              <Radio value='immediately'>立即执行</Radio>
              <Radio value='fixed_rate'>定时执行</Radio>
            </RadioGroup>
          )}
        </FormItem>
        {
          rate === 'fixed_rate' &&
          <FormItem
            {...formItemLayout}
            label={<span className='ant-form-item-required'>cron表达式</span>}
          >
            <Crontab
              form={this.props.form}
              handleClick={(values) => this.setState({ crontabUi: values })}
            />
          </FormItem>
        }
        <FormItem
          {...formItemLayout}
          label='未指定则全量'
        >
          {getFieldDecorator('sn', {
            rules: [{
              required: false, message: '请选择设备列表'
            }]
          })(
            <DeviceTable hideButton={true} form={this.props.form} />
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

const WrappedRegistrationForm = Form.create()(RegistrationForm);

export default WrappedRegistrationForm;
