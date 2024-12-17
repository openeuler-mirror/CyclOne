import React from 'react';
import { getWithArgs } from 'common/xFetch2';
import {
  Form,
  Button,
  notification,
  Input,
  Select,
  DatePicker,
  TreeSelect,
  InputNumber,
  Radio
} from 'antd';
const FormItem = Form.Item;
const RadioGroup = Radio.Group;
import { formItemLayout, tailFormItemLayout, getSearchList, BUILTIN } from 'common/enums';
const Option = Select.Option;
import { get } from 'common/xFetch2';
import renderTreeData from 'utils/get-tree-select-data';
const { SHOW_CHILD } = TreeSelect;
import moment from 'moment';

class MyForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      idc_id: 0,
      physical_area: '',
      initialValue: {
        'expected_arrival_date': moment().add(45, 'days'),
        'maintenance_service_date_begin': moment().add(45, 'days'),
        'maintenance_service_date_end': moment().add(5, 'years'),
      },
      room: [],
      treeData: []
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
      get(`/api/cloudboot/v1/order/${id}`).then(res => {
        if (res.status !== 'success') {
          return notification.error({ message: res.message });
        }
        const content = res.content || {};
        if (content.expected_arrival_date) {
          content.expected_arrival_date = moment(content.expected_arrival_date);
        }
        if (content.maintenance_service_date_begin) {
          content.maintenance_service_date_begin = moment(content.maintenance_service_date_begin);
        }
        if (content.maintenance_service_date_end) {
          content.maintenance_service_date_end = moment(content.maintenance_service_date_end);
        }
        this.setState({
          initialValue: {
            ...content,
            idc_name: content.idc.name,
            server_room_id: content.server_room.id,
            idc_id: content.idc.id
          },
          //机房库
          room: [content.server_room]
        });
      });
    }
  }


  //选择数据中心，获取机房列表
  getRoom = (v) => {
    this.setState({
      idc_id: v
    });
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


  //选择机房获取树
  getCabinet = (v) => {
    const { getFieldValue } = this.props.form;
    const server_room_id = getFieldValue('server_room_id');

    getWithArgs('/api/cloudboot/v1/server-usites/tree', { physical_area: v, idc_id: this.state.idc_id, usite_status: 'free', server_room_id: server_room_id }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      this.setState({
        treeData: renderTreeData(res.content.roots)
      });
    });
  };

  render() {
    const { idc, showSubmit, physicalArea, deviceCategory } = this.props;
    const { initialValue, treeData } = this.state;
    const { getFieldDecorator } = this.props.form;
    const tProps = {
      treeData,
      labelInValue: true,
      treeCheckable: true,
      showCheckedStrategy: SHOW_CHILD,
      searchPlaceholder: '请选择机位',
      style: {
        width: '100%'
      }
    };
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
                    idc.data.map((os, index) => <Option value={os.value} key={os.value}>{os.label}</Option>)
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
                })(<Input disabled={!showSubmit} />)}
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
            })(<Select disabled={!showSubmit}>
              {
                this.state.room.map((os, index) => <Option value={os.id} key={os.id}>{os.name}</Option>)
              }
            </Select>)}
          </FormItem>
          <FormItem {...formItemLayout} label='物理区域' >
            {getFieldDecorator('physical_area', {
              initialValue: initialValue.physical_area,
              rules: [
                {
                  required: true,
                  message: '请选择物理区域'
                }
              ]
            })(<Select disabled={!showSubmit} onChange={(value) => this.getCabinet(value)}>
              {
                !physicalArea.loading &&
                physicalArea.data.map((os, index) => <Option value={os.name} key={os.name}>{os.name}</Option>)
              }
            </Select>)}
          </FormItem>
          {
            showSubmit ? <FormItem {...formItemLayout} label='可用机位' >
              {getFieldDecorator('pre_occupied_usites', {
                initialValue: initialValue.pre_occupied_usites,
                rules: [
                  {
                    required: true,
                    message: '请选择可用机位'
                  }
                ]
              })(<TreeSelect disabled={!showSubmit} {...tProps} />)}
            </FormItem> : <FormItem {...formItemLayout} label='可用机位' >
                {getFieldDecorator('pre_occupied_usites', {
                  initialValue: initialValue.pre_occupied_usites,
                  rules: [
                    {
                      required: true,
                      message: '请选择可用机位'
                    }
                  ]
                })(<Input.TextArea disabled={!showSubmit} />)}
              </FormItem>
          }

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
          <FormItem {...formItemLayout} label='设备类型' >
            {getFieldDecorator('category', {
              initialValue: initialValue.category,
              rules: [
                {
                  required: true,
                  message: '请选择设备类型'
                }
              ]
            })(
              <Select disabled={!showSubmit}>
                {
                  !deviceCategory.loading &&
                  deviceCategory.data.map((os, index) => <Option value={os.category} key={os.category}>{os.category}</Option>)
                }
              </Select>
            )}
          </FormItem>
          <FormItem {...formItemLayout} label='数量' >
            {getFieldDecorator('amount', {
              initialValue: initialValue.amount,
              rules: [
                {
                  required: true,
                  message: '请输入名称'
                }
              ]
            })(<InputNumber style={{ width: '100%' }} disabled={!showSubmit} />)}
          </FormItem>
          <FormItem {...formItemLayout} label='预计到货时间' >
            {getFieldDecorator('expected_arrival_date', {
              initialValue: initialValue.expected_arrival_date,
              rules: [
                {
                  required: true,
                  message: '请输入预计到货时间'
                }
              ]
            })(<DatePicker style={{ width: '100%' }} disabled={!showSubmit} />)}
          </FormItem>
          <FormItem {...formItemLayout} label='资产归属' >
            {getFieldDecorator('asset_belongs', {
              initialValue: initialValue.asset_belongs,
              rules: [
                {
                  required: true,
                  message: '请输入资产归属方'
                }
              ]
            })(<Input disabled={!showSubmit} placeholder='请输入资产归属方'/>)}
          </FormItem>
          <FormItem {...formItemLayout} label='负责人' >
            {getFieldDecorator('owner', {
              initialValue: initialValue.owner,
              rules: [
                {
                  required: true,
                  message: '请输入设备负责人企业微信号，多个负责人使用逗号分隔',
                }
              ]
            })(<Input disabled={!showSubmit} placeholder='请输入设备负责人企业微信号，多个负责人使用逗号分隔' />)}
          </FormItem>
          <FormItem
            {...formItemLayout}
            label='是否租赁'
          >
            {getFieldDecorator('is_rental', {
              initialValue: initialValue.is_rental,
              rules: [{
                required: true, message: '请选择是否租赁'
              }]
            })(
              <RadioGroup>
                {getSearchList(BUILTIN).map(it => <Radio value={it.value}>{it.label}</Radio>)}
              </RadioGroup>
            )}
          </FormItem>
          <FormItem {...formItemLayout} label='维保服务供应商' >
            {getFieldDecorator('maintenance_service_provider', {
              initialValue: initialValue.maintenance_service_provider,
              rules: [
                {
                  required: true,
                  message: '请输入维保服务供应商'
                }
              ]
            })(<Input disabled={!showSubmit} />)}
          </FormItem>
          <FormItem {...formItemLayout} label='维保服务内容' >
            {getFieldDecorator('maintenance_service', {
              initialValue: initialValue.maintenance_service
            })(<Input.TextArea disabled={!showSubmit} autoSize={{ minRows: 2}} />)}
          </FormItem>
          <FormItem {...formItemLayout} label='物流服务内容' >
            {getFieldDecorator('logistics_service', {
              initialValue: initialValue.logistics_service
            })(<Input.TextArea disabled={!showSubmit} autoSize={{ minRows: 2}} />)}
          </FormItem>
          <FormItem {...formItemLayout} label='维保服务起始日期' >
            {getFieldDecorator('maintenance_service_date_begin', {
              initialValue: initialValue.maintenance_service_date_begin,
              rules: [
                {
                  required: true,
                  message: '请输入预计到货时间'
                }
              ]
            })(<DatePicker style={{ width: '100%' }} disabled={!showSubmit} />)}
          </FormItem>
          <FormItem {...formItemLayout} label='维保服务截止日期' >
            {getFieldDecorator('maintenance_service_date_end', {
              initialValue: initialValue.maintenance_service_date_end,
              rules: [
                {
                  required: true,
                  message: '请输入预计到货时间'
                }
              ]
            })(<DatePicker style={{ width: '100%' }} disabled={!showSubmit} />)}
          </FormItem>          
          <FormItem {...formItemLayout} label='备注' >
            {getFieldDecorator('remark', {
              initialValue: initialValue.remark
            })(<Input disabled={!showSubmit} />)}
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
              })(<Input disabled={!showSubmit} />)}
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
