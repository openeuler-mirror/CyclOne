import React from 'react';
import { getWithArgs } from 'common/xFetch2';
import {
  Form,
  Button,
  notification,
  Input,
  Select
} from 'antd';
const FormItem = Form.Item;
import { formItemLayout, tailFormItemLayout } from 'common/enums';
import { getSearchList, NETWORK_DEVICE_TYPE, NETWORK_DEVICE_STATUS} from "common/enums";
const Option = Select.Option;
import { get } from 'common/xFetch2';

class MyForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      initialValue: {},
      cabinet: [],
      room: []
    };
  }

  handleSubmit = ev => {
    ev && ev.preventDefault();
    ev && ev.stopPropagation();
    this.props.form.validateFields((error, values) => {
      if (error) {
        notification.warning({
          message: '还有未填写完整的选项'
        });
        return;
      }
      this.props.onSubmit(values);
    });
  };

  componentDidMount() {
    const { type, id } = this.props;
    if (type === '_detail' || type === '_update') {
      get(`/api/cloudboot/v1/network/devices/${id}`).then(res => {
        if (res.status !== 'success') {
          return notification.error({ message: res.message });
        }
        const content = res.content || {};
        this.setState({
          initialValue: {
            ...content,
            idc_name: content.idc.name,
            server_cabinet_id: content.server_cabinet.id,
            server_room_id: content.server_room.id
          },
          //机架库机房库
          cabinet: [content.server_cabinet],
          room: [content.server_room]
        });
      });
    }
  }


  //选择数据中心，获取机房列表
  getRoom = (v) => {
    getWithArgs('/api/cloudboot/v1/server-rooms', { page: 1, page_size: 100, idc_id: v }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      this.setState({
        room: res.content.records
      });

      //清空已选机房的值
      const { setFieldsValue, getFieldValue } = this.props.form;
      const room = getFieldValue('server_room_id');
      if (room) {
        setFieldsValue({ server_room_id: null });
      }
      //清空已选机架的值
      const cabinet = getFieldValue('server_cabinet_id');
      if (cabinet) {
        setFieldsValue({ server_cabinet_id: null });
      }
    });
  };

  //选择机房获取机架列表
  getCabinet = (v) => {
    getWithArgs('/api/cloudboot/v1/server-cabinets', { page: 1, page_size: 1000, server_room_id: v }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      this.setState({
        cabinet: res.content.records
      });

      //清空已选机架的值
      const { setFieldsValue, getFieldValue } = this.props.form;
      const cabinet = getFieldValue('server_cabinet_id');
      if (cabinet) {
        setFieldsValue({ server_cabinet_id: null });
      }
    });
  };

  render() {
    const { idc, showSubmit } = this.props;
    const { initialValue } = this.state;
    const { getFieldDecorator } = this.props.form;

    return (
      <div>
        <Form onSubmit={this.handleSubmit}>
          {
            showSubmit ?
              <FormItem {...formItemLayout} label='数据中心' >
                {getFieldDecorator('idc_id', {
                  initialValue: initialValue.idc_id,
                  rules: [
                    {
                      required: true,
                      message: '请选择数据中心'
                    }
                  ]
                })(<Select disabled={!showSubmit} onChange={(value) => this.getRoom(value)}>
                  {
                    !idc.loading &&
                    idc.data.map((os, index) => <Option value={os.id} key={os.id}>{os.name}</Option>)
                  }
                </Select>)}
              </FormItem> :
              <FormItem {...formItemLayout} label='数据中心' >
                {getFieldDecorator('idc_name', {
                  initialValue: initialValue.idc_name,
                  rules: [
                    {
                      required: true,
                      message: '请选择数据中心'
                    }
                  ]
                })(<Input disabled={!showSubmit}/>)}
              </FormItem>
          }
          <FormItem {...formItemLayout} label='机房管理单元' >
            {getFieldDecorator('server_room_id', {
              initialValue: initialValue.server_room_id,
              rules: [
                {
                  required: true,
                  message: '请选择所属机房'
                }
              ]
            })(<Select disabled={!showSubmit} onChange={(value) => this.getCabinet(value)}>
              {
                this.state.room.map((os, index) => <Option value={os.id} key={os.id}>{os.name}</Option>)
              }
            </Select>)}
          </FormItem>
          <FormItem {...formItemLayout} label='机架编号' >
            {getFieldDecorator('server_cabinet_id', {
              initialValue: initialValue.server_cabinet_id,
              rules: [
                {
                  required: true,
                  message: '请选择所属机架'
                }
              ]
            })(<Select disabled={!showSubmit}>
              {
                this.state.cabinet.map((os, index) => <Option value={os.id} key={os.id}>{os.number}</Option>)
              }
            </Select>)}
          </FormItem>
          <FormItem {...formItemLayout} label='固资编号' >
            {getFieldDecorator('fixed_asset_number', {
              initialValue: initialValue.fixed_asset_number,
              rules: [
                {
                  required: true,
                  message: '请填写固资编号'
                }
              ]
            })(<Input disabled={!showSubmit} />)}
          </FormItem>
          <FormItem {...formItemLayout} label='序列号' >
            {getFieldDecorator('sn', {
              initialValue: initialValue.sn,
              rules: [
                {
                  required: true,
                  message: '请输入序列号'
                }
              ]
            })(
              <Input disabled={!showSubmit} />
            )}
          </FormItem>
          <FormItem {...formItemLayout} label='名称' >
            {getFieldDecorator('name', {
              initialValue: initialValue.name,
              rules: [
                {
                  required: true,
                  message: '请输入名称'
                }
              ]
            })(<Input disabled={!showSubmit} />)}
          </FormItem>
          <FormItem {...formItemLayout} label='产品型号' >
            {getFieldDecorator('model', {
              initialValue: initialValue.model,
              rules: [
                {
                  required: true,
                  message: '请输入产品型号'
                }
              ]
            })(<Input disabled={!showSubmit} />)}
          </FormItem>
          <FormItem {...formItemLayout} label='厂商' >
            {getFieldDecorator('vendor', {
              initialValue: initialValue.vendor,
              rules: [
                {
                  required: true,
                  message: '请输入厂商'
                }
              ]
            })(<Input disabled={!showSubmit} />)}
          </FormItem>
          <FormItem {...formItemLayout} label='TOR' >
            {getFieldDecorator('tor', {
              initialValue: initialValue.tor,
              rules: [
                {
                  required: true,
                  message: '请输入TOR'
                }
              ]
            })(<Input disabled={!showSubmit} />)}
          </FormItem>
          <FormItem {...formItemLayout} label='操作系统' >
            {getFieldDecorator('os', {
              initialValue: initialValue.os,
              rules: [
                {
                  required: true,
                  message: '请输入操作系统'
                }
              ]
            })(<Input disabled={!showSubmit} />)}
          </FormItem>
          <FormItem {...formItemLayout} label='类型' >
            {getFieldDecorator('type', {
              initialValue: initialValue.type,
              rules: [
                {
                  required: true,
                  message: '请选择类型'
                }
              ]
            })(<Select disabled={!showSubmit}>
              {
                getSearchList(NETWORK_DEVICE_TYPE).map((it) => <Option value={it.value} key={it.value}>{it.label}</Option>)
              }
            </Select>)}
          </FormItem>
          <FormItem {...formItemLayout} label='用途' >
            {getFieldDecorator('usage', {
              initialValue: initialValue.usage,
              rules: [
                {
                  required: true,
                  message: '请输入用途'
                }
              ]
            })(<Input disabled={!showSubmit} />)}
          </FormItem>
          <FormItem {...formItemLayout} label='状态' >
            {getFieldDecorator('status', {
              initialValue: initialValue.status,
              rules: [
                {
                  required: true,
                  message: '请输入状态'
                }
              ]
            })(
              <Select disabled={!showSubmit}>
                {
                  getSearchList(NETWORK_DEVICE_STATUS).map((it) => <Option value={it.value} key={it.value}>{it.label}</Option>)
                }
              </Select>
            )}
          </FormItem>          
          {
            !showSubmit &&
            <FormItem {...formItemLayout} label='创建时间' >
              {getFieldDecorator('created_at', {
                initialValue: initialValue.created_at
              })(<Input disabled={!showSubmit} />)}
            </FormItem>
          }
          {
            !showSubmit &&
            <FormItem {...formItemLayout} label='更新时间' >
              {getFieldDecorator('updated_at', {
                initialValue: initialValue.updated_at
              })(<Input disabled={!showSubmit}/>)}
            </FormItem>
          }
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
