import React from 'react';
import { get, getWithArgs } from 'common/xFetch2';
import {
  Form,
  Button,
  notification,
  Input,
  InputNumber,
  Select
} from 'antd';
const FormItem = Form.Item;
import { formItemLayout, tailFormItemLayout } from 'common/enums';
import { getSearchList, USITE_STATUS } from "common/enums";
const Option = Select.Option;
import Port from './port';


class MyForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      initialValue: {},
      cabinet: [],
      network_devices: [],
      physicalAreas: []
    };
  }

  componentDidMount() {
    
    const { type, id } = this.props;
    if (type === 'detail' || type === '_update') {
      get(`/api/cloudboot/v1/server-usites/${id}`).then(res => {
        if (res.status !== 'success') {
          return notification.error({ message: res.message });
        }
        const content = res.content || {};
        // 编辑时需要根据机架拉取对应网络区域的所有枚举值
        getWithArgs('/api/cloudboot/v1/server-cabinets', { page: 1, page_size: 1000, server_room_id: content.server_room.id , server_cabinet_id: content.server_cabinet.id }).then(res => {
          console.log(res)
          if (res.status !== 'success') {
            return notification.error({ message: res.message });
          }
          res.content.records.map((it) => {
            get(`/api/cloudboot/v1/physical-areas?network_area_id=${it.network_area.id}`).then(res => {
              this.setState({
                physicalAreas: res.content.list || []
              });
            });
          });
        });

        this.setState({
          initialValue: {
            ...content,
            server_cabinet_id: content.server_cabinet.id,
            server_room_id: content.server_room.id,
            //查看回显字段
            idc_name: content.idc.name
          },
          //机架库
          cabinet: [content.server_cabinet]
        });




      });
    }
    if (type === '_update' || type === '_create') {
      //不指定机架ID
      getWithArgs(`/api/cloudboot/v1/network/devices`, { page: 1,
        page_size: 10000
        // server_cabinet_id: this.props.server_cabinet_id
      }).then(res => {
        if (res.status !== 'success') {
          return notification.error({ message: res.message });
        }
        this.setState({
          network_devices: res.content.records
        });
      });
    }
   
  }

  getNetworkAreaId = (v, options) => {
    const id = options.props.neworkAreaId;
    get(`/api/cloudboot/v1/physical-areas?network_area_id=${id}`).then(res => {
      this.setState({
        physicalAreas: res.content.list || []
      });
    });
  }

  getParamValue = (values, label) => {
    const allKeys = Object.keys(values);
    const target_keys = allKeys.filter(
      key => key.indexOf(label) >= 0
    );
    let target = [];
    target_keys.map(key => {
      const id = key.replace(`${label}-name-`, '');
      if (values[`${label}-name-${id}`] || values[`${label}-port-${id}`]) {
        target.push({
          name: values[`${label}-name-${id}`],
          port: values[`${label}-port-${id}`]
        });
      }
    });
    return target;
  };

  handleSubmit = ev => {
    ev && ev.preventDefault();
    ev && ev.stopPropagation();
    this.props.form.validateFields((error, values) => {

      values.extranet_switches = this.getParamValue(values, 'extranet_switches');
      values.oobnet_switches = this.getParamValue(values, 'oobnet_switches');
      values.intranet_switches = this.getParamValue(values, 'intranet_switches');

      const allKeys = Object.keys(values);
      const delete_keys = allKeys.filter(
        key => key.indexOf('-name-') >= 0 || key.indexOf('-port-') >= 0
      );
      delete_keys.forEach(key => {
        delete values[key];
      });
      if (error) {
        notification.warning({
          message: '还有未填写完整的选项'
        });
        return;
      }
      this.props.onSubmit(values);
    });
  };

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
    const { room, showSubmit } = this.props;
    const { initialValue, network_devices } = this.state;
    const { getFieldDecorator } = this.props.form;

    return (
      <div>
        <Form onSubmit={this.handleSubmit}>
          <FormItem {...formItemLayout} label='机位编号' >
            {getFieldDecorator('number', {
              initialValue: initialValue.number,
              rules: [
                {
                  required: true,
                  message: '请填写机位编号'
                }
              ]
            })(<Input disabled={!showSubmit}/>)}
          </FormItem>
          {
            !showSubmit &&
            <FormItem {...formItemLayout} label='数据中心' >
              {getFieldDecorator('idc_name', {
                initialValue: initialValue.idc_name
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
                !room.loading &&
                room.data.map((os, index) => <Option value={os.id} key={os.id}>{os.name}</Option>)
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
            })(<Select disabled={!showSubmit} onChange={this.getNetworkAreaId} showSearch={true}
              filterOption={(input, option) =>
                option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
              }
            >
              {
                this.state.cabinet.map((os, index) => <Option neworkAreaId={os.network_area ? os.network_area.id : null} value={os.id} key={os.id}>{os.number}</Option>)
              }
            </Select>)}
          </FormItem>
          <FormItem {...formItemLayout} label='外网交换机_端口'>
            {getFieldDecorator('extranet_switches', {
              initialValue: initialValue.extranet_switches
            })(
              <Port network_devices={network_devices} label='外网交换机' required={false} name='extranet_switches' data={initialValue.extranet_switches} form={this.props.form} disabled={!showSubmit}/>
            )}
          </FormItem>
          <FormItem {...formItemLayout} label={<span className='ant-form-item-required'>内网交换机_端口</span>}>
            {getFieldDecorator('intranet_switches', {
              initialValue: initialValue.intranet_switches
            })(
              <Port network_devices={network_devices} label='内网交换机' required={true} name='intranet_switches' data={initialValue.intranet_switches} form={this.props.form} disabled={!showSubmit}/>
            )}
          </FormItem>
          <FormItem {...formItemLayout} label={<span className='ant-form-item-required'>管理交换机_端口</span>}>
            {getFieldDecorator('oobnet_switches', {
              initialValue: initialValue.oobnet_switches
            })(
              <Port network_devices={network_devices} label='管理交换机' required={true} name='oobnet_switches' data={initialValue.oobnet_switches} form={this.props.form} disabled={!showSubmit}/>
            )}
          </FormItem>
          <FormItem {...formItemLayout} label='内外网端口速率' >
            {getFieldDecorator('la_wa_port_rate', {
              initialValue: initialValue.la_wa_port_rate,
              rules: [
                {
                  required: true,
                  message: '请填写端口速率'
                }
              ]              
            })(
              <Select>
              <Option value='GE'>{'GE'}</Option>
              <Option value='10GE'>{'10GE'}</Option>
              <Option value='25GE'>{'25GE'}</Option>
              <Option value='40GE'>{'40GE'}</Option>
              </Select>
            )}
          </FormItem>
          <FormItem {...formItemLayout} label='起始U数' >
            {getFieldDecorator('beginning', {
              initialValue: initialValue.beginning,
              rules: [
                {
                  required: true,
                  message: '请填写起始U数'
                }
              ]
            })(<InputNumber disabled={!showSubmit} style={{ width: '100%' }}/>)}
          </FormItem>
          <FormItem {...formItemLayout} label='机位高度' >
            {getFieldDecorator('height', {
              initialValue: initialValue.height,
              rules: [
                {
                  required: true,
                  message: '请输入机位高度'
                }
              ]
            })(
              <InputNumber disabled={!showSubmit} min={1} precision={0} style={{ width: '100%' }}/>
            )}
          </FormItem>
          <FormItem {...formItemLayout} label='物理区域' >
            {getFieldDecorator('physical_area', {
              initialValue: initialValue.physical_area,
              rules: [
                {
                  required: true,
                  message: '请输入物理区域'
                }
              ]
            })(<Select disabled={!showSubmit}>
              {
                (this.state.physicalAreas || []).map(it => <Option value={it.name}>{it.name}</Option>)
              }
            </Select>)}
          </FormItem>
          {
            !showSubmit &&
            <FormItem {...formItemLayout} label='机位状态' >
              {getFieldDecorator('status', {
                initialValue: initialValue.status
              })(<Select disabled={!showSubmit}>
                {
                  getSearchList(USITE_STATUS).map((it) => <Option value={it.value} key={it.value}>{it.label}</Option>)
                }
              </Select>)}
            </FormItem>
          }
          {
            !showSubmit &&
            <FormItem {...formItemLayout} label='创建时间' >
              {getFieldDecorator('created_at', {
                initialValue: initialValue.created_at
              })(<Input disabled={!showSubmit}/>)}
            </FormItem>
          }
          {
            !showSubmit &&
            <FormItem {...formItemLayout} label='创建用户' >
              {getFieldDecorator('creator', {
                initialValue: initialValue.creator
              })(<Input disabled={!showSubmit}/>)}
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
          <FormItem {...formItemLayout} label='备注' >
            {getFieldDecorator('remark', {
              initialValue: initialValue.remark
            })(<Input.TextArea disabled={!showSubmit} />)}
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
