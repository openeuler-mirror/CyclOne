import React from 'react';
import { connect } from 'react-redux';
import Layout from 'components/layout/page-layout';
import { notification, Input, Button, Form, Select } from 'antd';
import { getBreadcrumb } from 'common/utils';
import { hashHistory } from 'react-router';
import { formItemLayout_page, tailFormItemLayout_page } from 'common/enums';
import addDevice from '../add-device';
import DeviceTable from './device_edit';
import { post } from 'common/xFetch2';
const { TextArea } = Input;
import { getPermissonBtn } from 'common/utils';
import action from './import';

class Container extends React.Component {

  state = {
    dataSource: [],
    loading: false
  };

  //下载导入模板
  downloadImportTemplate = () => {
    window.open('assets/files/device_migration_import.xlsx');
  };
  importData = () => {
    action({
      reload: (data) => {
        this.setFormValue(data);
      }
    });
  };


  handleSubmit = (e) => {
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      let isErr = false;
      values.data.forEach(data => {
        if (!data.dst_idc_id) {
          isErr = true;
          return;
        }

        if (!data.dst_server_room_id && !data.dst_cabinet_id && !data.dst_usite_id) {
          if (!data.dst_store_room_id || !data.dst_virtual_cabinet_id) {
            isErr = true;
            return;
          }
        }

        if (!data.dst_store_room_id && !data.dst_virtual_cabinet_id) {
          if (!data.dst_server_room_id || !data.dst_cabinet_id || !data.dst_usite_id) {
            isErr = true;
            return;
          }
        }

        isErr = false;
      });

      if (err || isErr) {
        return notification.error({ message: '还有未填写完整的项' });
      }

      values.approvers = [ values.approvers0, values.approvers1 ];
      delete values.approvers0;
      delete values.approvers1;
      post('/api/cloudboot/v1/approvals/devices/migrations', values).then(res => {
        if (res.status !== 'success') {
          return notification.error({ message: res.message });
        }
        notification.success({ message: res.message });
        //操作成功置灰保存按钮，防止多次点击提交
        this.setState({ loading: true });
        hashHistory.push('/approval');
      });

    });
  };


  setFormValue = (data) => {
    this.setState({ dataSource: data });
    const { setFieldsValue } = this.props.form;
    setFieldsValue({ data: data.map(it => {
      return {
        "sn": it.sn,
        "dst_server_room_id": it.dst_server_room_id,
        "dst_cabinet_id": it.dst_cabinet_id,
        "dst_usite_id": it.dst_usite_id,
        "dst_idc_id": it.dst_idc_id,
        "dst_store_room_id": it.dst_store_room_id,
        "dst_virtual_cabinet_id": it.dst_virtual_cabinet_id
      };
    }) });
    console.log(data);
    setFieldsValue({ front_data: JSON.stringify(data.map(i => {
      return {
        fixed_asset_number: i.fixed_asset_number,
        sn: i.sn,
        model: i.model,
        old_idc: i.idc ? i.idc.name : '',
        old_server_room: i.server_room ? i.server_room.name : '',
        old_server_usite: i.server_usite ? i.server_usite.number : '',
        old_server_cabinet: i.server_cabinet ? i.server_cabinet.number : '',
        old_store_room: i.store_room ? i.store_room.name : '',
        old_virtual_cabinets: i.virtual_cabinets ? i.virtual_cabinets.number : '',
        idc: i.dst_idc_name,
        server_room_name: i.dst_server_room_name,
        server_cabinet_number: i.dst_cabinet_number,
        server_usite_number: i.dst_usite_number,
        store_room_name: i.dst_store_room_name,
        virtual_cabinet_number: i.dst_virtual_cabinet_number,
        in_store: i.operation_status == 'in_store'
      };
    })) });
  };
  render() {
    const { getFieldDecorator } = this.props.form;
    return (
      <Layout>
        <div style={{ marginTop: -10 }}>
          {getBreadcrumb('物理机搬迁')}
        </div>
        <div>
          <div className='operate_btns'>
            <Button icon='plus' style={{ marginRight: 8 }} onClick={() => addDevice({
              // getServerRoom: true,
              limit: 50,
              query: {
                operation_status: 'on_shelve,pre_deploy,in_store,pre_move'
              },
              handleDeviceSubmit: (tableData, onSuccess) => {
                const selectedRows = tableData.selectedRows || [];
                this.setFormValue(selectedRows);
                onSuccess();
              }
            })} type='primary'
            >
            添加设备
            </Button>
            <Button.Group>
              <Button
                onClick={() => this.downloadImportTemplate()}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_approval_physical_machine_move_download')}
              >
                下载导入模板
              </Button>
              <Button
                onClick={this.importData}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_approval_physical_machine_move_import')}
              >
                导入
              </Button>
            </Button.Group>
          </div>
          <Form>
            <Form.Item
              label='已选设备'
              {...formItemLayout_page}
            >
              {getFieldDecorator('data', {
                rules: [{
                  required: true,
                  message: '请选择设备'
                }]
              })(
                <DeviceTable
                  dataSource={this.state.dataSource}
                  setFormValue={this.setFormValue}
                  form={this.props.form}
                />
              )}
            </Form.Item>
            {/*<Form.Item*/}
            {/*label='申请标题'*/}
            {/*{...formItemLayout_page}*/}
            {/*>*/}
            {/*{getFieldDecorator('title', {*/}
            {/*})(*/}
            {/*<Input style={{ width: 400 }}/>*/}
            {/*)}*/}
            {/*</Form.Item>*/}
            <Form.Item
              label='备注'
              {...formItemLayout_page}
            >
              {getFieldDecorator('remark', {
              })(
                <TextArea rows={4} style={{ width: 400 }}/>
              )}
            </Form.Item>
            <Form.Item >
              {getFieldDecorator('front_data', {
              })(
                <Input hidden={true}/>
              )}
            </Form.Item>
            <Form.Item
              label='审批人'
              {...formItemLayout_page}
            >
              {getFieldDecorator('approvers0', {
                rules: [{
                  required: true,
                  message: '请选择审批人'
                }]
              })(
                <Select style={{ width: 400 }}>
                  {
                    (this.props.userList.data || []).map(it => <Option disabled={it.id === this.props.userInfo.id} value={it.id}>{it.name}</Option>)
                  }
                </Select>
              )}
            </Form.Item>
            <Form.Item
              label='实施人'
              {...formItemLayout_page}
            >
              {getFieldDecorator('approvers1', {
                rules: [{
                  required: true,
                  message: '请选择实施人'
                }]
              })(
                <Select style={{ width: 400 }}>
                  {
                    (this.props.userList.data || []).map(it => <Option value={it.id}>{it.name}</Option>)
                  }
                </Select>
              )}
            </Form.Item>
          </Form>
          <Form.Item {...tailFormItemLayout_page}>
            <Button
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_approval_physical_machine_move')}
              loading={this.state.loading} onClick={this.handleSubmit} type='primary' style={{ marginRight: 8 }}
            >
              提交
            </Button>
            <Button onClick={() => hashHistory.push('/approval')}>
              取消
            </Button>
          </Form.Item>
        </div>
      </Layout>
    );
  }
}

function mapStateToProps(state) {
  return {
    userInfo: state.getIn([ 'global', 'userData' ]).toJS(),
    userList: state.getIn([ 'global', 'userList' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}


export default connect(mapStateToProps, mapDispatchToProps)(Form.create()(Container));
