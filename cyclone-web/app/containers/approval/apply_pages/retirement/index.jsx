import React from 'react';
import { connect } from 'react-redux';
import Layout from 'components/layout/page-layout';
import { Input, notification, Button, Form, Select } from 'antd';
import { getBreadcrumb } from 'common/utils';
import { hashHistory } from 'react-router';
import { formItemLayout_page, tailFormItemLayout_page } from 'common/enums';
import addDevice from '../add-device';
import DeviceTable from '../device_table';
import { post } from 'common/xFetch2';
const { TextArea } = Input;
import { getPermissonBtn } from 'common/utils';

class Container extends React.Component {

  state = {
    dataSource: [],
    loading: false
  };

  handleSubmit = (e) => {
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (err) {
        notification.error({ message: '还有未填写完整的项' });
      }
      values.approvers = [ values.approvers0, values.approvers1 ];
      delete values.approvers0;
      delete values.approvers1;
      post('/api/cloudboot/v1/approvals/devices/retirements', values).then(res => {
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
    setFieldsValue({ sns: data.map(it => it.sn) });
    setFieldsValue({ front_data: JSON.stringify(data) });
  };
  render() {
    const { getFieldDecorator } = this.props.form;

    return (
      <Layout>
        <div style={{ marginTop: -10 }}>
          {getBreadcrumb('物理机退役')}
        </div>
        <div>
          <div className='operate_btns'>
            <Button icon='plus' onClick={() => addDevice({
              query: {
                operation_status: 'pre_retire'
              },
              limit: 50,
              handleDeviceSubmit: (tableData, onSuccess) => {
                const selectedRows = tableData.selectedRows || [];
                this.setFormValue(selectedRows);
                onSuccess();
              }
            })} type='primary'
            >
            添加设备
            </Button>
          </div>
          <Form>
            <Form.Item
              label='已选设备'
              {...formItemLayout_page}
            >
              {getFieldDecorator('sns', {
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
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_approval_physical_machine_retirement')}
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
