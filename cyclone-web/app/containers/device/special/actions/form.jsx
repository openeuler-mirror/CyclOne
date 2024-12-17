import React from 'react';
import { Input, Checkbox, Select, Button, Form, notification, DatePicker, InputNumber } from 'antd';
const FormItem = Form.Item;
import { formItemLayout, tailFormItemLayout } from 'common/enums';
const { Option } = Select;
import { getWithArgs } from 'common/xFetch2';
import moment from 'moment';


class MyForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      cabinet: [],
      usite: [],
      initialValue: {
        'maintenance_service_date_begin': moment().add(45, 'days'),
        'maintenance_months': 60,
      },
    };
  }
  //state = {
  //  cabinet: [],
  //  usite: [],
  //  initialValue: {
  //    'maintenance_service_date_begin': moment().add(45, 'days'),
  //    'maintenance_months': 60
  //  }
  //};
 
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

  clearCabinet = () => {
    const { setFieldsValue, getFieldValue } = this.props.form;
    const cabinet = getFieldValue('server_cabinet_number');
    if (cabinet) {
      setFieldsValue({ server_cabinet_number: null });
    }
    this.setState({
      cabinet: []
    });
  };
  clearUsite = () => {
    const { setFieldsValue, getFieldValue } = this.props.form;
    const usite = getFieldValue('server_usite_number');
    if (usite) {
      setFieldsValue({ server_usite_number: null });
    }
    this.setState({
      usite: []
    });
  };
  //选择机房获取机架列表
  getCabinet = (v, page, clear) => {
    this.getCabinetPage(v, page);
    if (clear) {
      this.clearCabinet();
      this.clearUsite();
    }
  };

  getCabinetPage = (v, page) => {
    getWithArgs('/api/cloudboot/v1/server-cabinets', { page: page, page_size: 1000, server_room_id: v }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      this.setState((preState) => {
        return {
          cabinet: [ ...preState.cabinet, ...res.content.records ]
        };
      }, () => {
        if (res.content.total_pages > page) {
          this.getCabinetPage(v, page + 1);
        }
      });
    });
  };

  //选择机架获取机位列表
  getUsite = (v, page, clear) => {
    this.getUsitePage(v, page);
    if (clear) {
      this.clearUsite();
    }
  };

  getUsitePage = (v, page) => {
    getWithArgs('/api/cloudboot/v1/server-usites', { page: page, page_size: 100, server_cabinet_id: v, status: 'free,pre_occupied' }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      this.setState((preState) => {
        return {
          usite: [ ...preState.usite, ...res.content.records ]
        };
      }, () => {
        if (res.content.total_pages > page) {
          this.getUsitePage(v, page + 1);
        }
      });
    });
  };



  render() {
    const { getFieldDecorator } = this.props.form;
    const { room } = this.props;
    const { initialValue } = this.state;
    console.log(initialValue);

    return <Form onSubmit={this.handleSubmit}>
      <FormItem {...formItemLayout} label='序列号' >
        {getFieldDecorator('sn', {
          rules: [
            {
              required: true
            }
          ]

        })(
          <Input />
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='型号' >
        {getFieldDecorator('model', {
          rules: [
            {
              required: true
            }
          ]
        })(
          <Input />
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='厂商' >
        {getFieldDecorator('vendor', {
          rules: [
            {
              required: true
            }
          ]
        })(
          <Input />
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='硬件描述' >
        {getFieldDecorator('hardware_remark', {
          rules: [
            {
              required: true
            }
          ]
        })(
          <Input />
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='设备类型' >
        {getFieldDecorator('category', {
          rules: [
            {
              required: true,
            }
          ]
        })(
          <Select>
            <Option value='堡垒机'>{'堡垒机'}</Option>
            <Option value='加密机'>{'加密机'}</Option>
            <Option value='前置机'>{'前置机'}</Option>
            <Option value='SSL服务器'>{'SSL服务器'}</Option>
            <Option value='签名服务器'>{'签名服务器'}</Option>
            <Option value='X86服务器'>{'X86服务器'}</Option>
            <Option value='GPU服务器'>{'GPU服务器'}</Option>
            <Option value='授时服务器'>{'授时服务器'}</Option>
            <Option value='其他'>{'其他'}</Option>
        </Select>
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='负责人' >
        {getFieldDecorator('owner', {
          rules: [
            {
              required: true,
              message: '请输入设备负责人企业微信号，多个负责人使用逗号分隔',
            }
          ]
        })(
          <Input placeholder='请输入设备负责人企业微信号，多个负责人使用逗号分隔'/>
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='维保服务起始日期' >
            {getFieldDecorator('maintenance_service_date_begin', {
              initialValue: initialValue.maintenance_service_date_begin,
              rules: [
                {
                  required: false,
                  message: '默认'
                }
              ]
            })(<DatePicker style={{ width: '100%' }} />)}
      </FormItem>
          <FormItem {...formItemLayout} label='保修期（月数）' >
            {getFieldDecorator('maintenance_months', {
              initialValue: initialValue.maintenance_months,
              rules: [
                {
                  required: true,
                  message: '保修期（月数）'
                }
              ]
            })(<InputNumber style={{ width: '100%' }} />)}
      </FormItem>                        
      <FormItem {...formItemLayout} label='订单编号' >
        {getFieldDecorator('order_number', {
          rules: [
            {
              required: false
            }
          ]
        })(
          <Input />
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='操作系统名称' >
        {getFieldDecorator('os_release_name', {
          rules: [
            {
              required: true,
              message: '请输入操作系统名称，如：openEuler 20.03 (LTS)',
            }
          ]
        })(
          <Input placeholder='请输入操作系统名称，如：openEuler 20.03 (LTS)'/>
        )}
      </FormItem>      
      <FormItem {...formItemLayout} label='是否分配内网IP' >
        {getFieldDecorator('need_intranet_ip', {

        })(
          <Checkbox />
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='是否分配外网IP' >
        {getFieldDecorator('need_extranet_ip', {

        })(
          <Checkbox />
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='机房名称' >
        {getFieldDecorator('server_room_id', {
          rules: [
            {
              required: true
            }
          ]
        })(
          <Select
            showSearch={true}
            filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
            onChange={(value) => this.getCabinet(value, 1, true)}
          >
            {
              !room.loading && room.data.map((os, index) => <Option value={os.value} key={os.value}>{os.label}</Option>)
            }
          </Select>
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='机架编号' >
        {getFieldDecorator('server_cabinet_id', {
          rules: [
            {
              required: true
            }
          ]
        })(
          <Select
            showSearch={true}
            filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
            onChange={(value) => this.getUsite(value, 1, true)}
          >
            {
              this.state.cabinet.map((os, index) => <Option value={os.id} key={os.id}>{os.number}</Option>)
            }
          </Select>
        )}
      </FormItem>
      <FormItem {...formItemLayout} label='机位编号' >
        {getFieldDecorator('server_usite_id', {
          rules: [
            {
              required: true
            }
          ]
        })(
          <Select
            showSearch={true}
            filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
          >
            {
              this.state.usite.map((os, index) => <Option value={os.id} key={os.id}>{os.number}</Option>)
            }
          </Select>
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
export default Form.create()(MyForm);

