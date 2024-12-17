import React from 'react';
import { Input, Select, Button, Form, notification } from 'antd';
import { post, put } from 'common/xFetch2';
const FormItem = Form.Item;
import { formItemLayout, tailFormItemLayout } from 'common/enums';
import Popup from 'components/popup';
const { Option } = Select;
import { OPERATION_STATUS, getSearchList } from 'common/enums';

export default function action(options) {

  const onSubmit = (values) => {
    const data = options.records.map(it => {
      return {
        "fixed_asset_number": it.fixed_asset_number,
        "sn": it.sn,
        "usage": values.usage
      };
    });
    put('/api/cloudboot/v1/devices', { devices: data }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: '操作成功' });
      options.reload();
      Popup.close();
    });
  };

  Popup.open({
    title: `批量修改用途`,
    width: 600,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <WrapperForm
        initialValue={options.record}
        username={options.username}
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
    return <Form onSubmit={this.handleSubmit}>
      <FormItem {...formItemLayout} label='用途' >
        {getFieldDecorator('usage', {
        })(
          <Input.TextArea/>
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
