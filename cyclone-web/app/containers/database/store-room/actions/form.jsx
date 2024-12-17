import React from 'react';
import { get } from 'common/xFetch2';
import {
  Form,
  Button,
  notification,
  Input,
  Select
} from 'antd';
const FormItem = Form.Item;
import { formItemLayout, tailFormItemLayout } from 'common/enums';
const Option = Select.Option;

class MyForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      initialValue: {}
    };
  }

  componentDidMount() {
    const { type, id } = this.props;
    if (type.indexOf('detail') !== -1 || type === '_update') {
      get(`/api/cloudboot/v1/store-room/${id}`).then(res => {
        if (res.status !== 'success') {
          return notification.error({ message: res.message });
        }
        const idc = res.content.idc || {};
        this.setState({
          initialValue: {
            ...res.content,
            idc_name: idc.name,
            idc_data: JSON.stringify({ first_server_room: idc.first_server_room, id: idc.id }),
            first_server_room: res.content.first_server_room.name
          }
        });
      });
    }
  }

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

  getRoom = (v) => {
    //清空已选一级机房的值
    const { setFieldsValue } = this.props.form;
    if (v) {
      setFieldsValue({ first_server_room: null });
    }
  };

  render() {
    const { idc, showSubmit, type } = this.props;
    const { initialValue } = this.state;
    const { getFieldDecorator, getFieldValue } = this.props.form;
    const idc_data = JSON.parse(getFieldValue('idc_data') || initialValue.idc_data || '{}');

    return (
      <div>
        <Form onSubmit={this.handleSubmit}>

          {
            type.indexOf('detail') !== -1 ?
              <FormItem {...formItemLayout} label='数据中心' >
                {getFieldDecorator('idc_name', {
                  initialValue: initialValue.idc_name,
                  rules: [
                    {
                      required: true,
                      message: '请选择数据中心'
                    }
                  ]
                })(
                  <Input disabled={!showSubmit} />
                )}
              </FormItem> :
              <FormItem {...formItemLayout} label='数据中心' >
                {getFieldDecorator('idc_data', {
                  initialValue: initialValue.idc_data,
                  rules: [
                    {
                      required: true,
                      message: '请选择数据中心'
                    }
                  ]
                })(
                  <Select disabled={!showSubmit} onChange={(v) => this.getRoom(v)}>
                    {
                      !idc.loading &&
                      idc.data.map((os, index) => <Option value={JSON.stringify({ first_server_room: os.first_server_room, id: os.id })} key={os.id}>{os.name}</Option>)
                    }
                  </Select>
                )}
              </FormItem>
          }



          <FormItem {...formItemLayout} label='所属一级机房' >
            {getFieldDecorator('first_server_room', {
              initialValue: initialValue.first_server_room,
              rules: [
                {
                  required: true,
                  message: '请选择所属一级机房'
                }
              ]
            })(
              <Select disabled={!showSubmit} >
                {
                  (idc_data.first_server_room || []).map(item => <Option value={item.name} key={item.id}>{item.name}</Option>)
                }
              </Select>
            )}
          </FormItem>
          <FormItem {...formItemLayout} label='库房管理单元' >
            {getFieldDecorator('name', {
              initialValue: initialValue.name,
              rules: [
                {
                  required: true,
                  message: '请输入库房管理单元'
                }
              ]
            })(
              <Input disabled={!showSubmit} />
            )}
          </FormItem>
          <FormItem {...formItemLayout} label='城市' >
            {getFieldDecorator('city', {
              initialValue: initialValue.city,
              rules: [
                {
                  required: true,
                  message: '请输入城市'
                }
              ]
            })(
              <Input disabled={!showSubmit} />
            )}
          </FormItem>
          <FormItem {...formItemLayout} label='地址' >
            {getFieldDecorator('address', {
              initialValue: initialValue.address,
              rules: [
                {
                  required: true,
                  message: '请输入地址'
                }
              ]
            })(<Input disabled={!showSubmit} />)}
          </FormItem>
          <FormItem {...formItemLayout} label='库房负责人' >
            {getFieldDecorator('store_room_manager', {
              initialValue: initialValue.store_room_manager,
              rules: [
                {
                  required: true,
                  message: '请输入库房负责人'
                }
              ]
            })(<Input disabled={!showSubmit} />)}
          </FormItem>
          <FormItem {...formItemLayout} label='供应商负责人' >
            {getFieldDecorator('vendor_manager', {
              initialValue: initialValue.vendor_manager,
              rules: [
                {
                  required: true,
                  message: '请输入供应商负责人'
                }
              ]
            })(<Input disabled={!showSubmit} />)}
          </FormItem>
          <FormItem {...tailFormItemLayout}>
            <div className='pull-right'>
              <Button onClick={() => this.props.onCancel()}>取消</Button>
              {
                showSubmit &&
                <Button
                  style={{ marginLeft: 8 }}
                  type='primary'
                  htmlType='submit'
                >
                  提交
                </Button>
              }
            </div>
          </FormItem>
        </Form>
      </div>
    );
  }
}

export default Form.create()(MyForm);
