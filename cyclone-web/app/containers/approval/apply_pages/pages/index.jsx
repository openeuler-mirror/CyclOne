import React from 'react';
import { connect } from 'react-redux';
import Layout from 'components/layout/page-layout';
import { Input, notification, Button, Form, Select } from 'antd';
import { getBreadcrumb } from 'common/utils';
import { hashHistory } from 'react-router';
import { formItemLayout_page, tailFormItemLayout_page } from 'common/enums';
import addData from './add-data';
import FormTable from './form_table';
import { post } from 'common/xFetch2';
const { TextArea } = Input;
// import { getPermissonBtn } from 'common/utils';
import { TYPE } from './config';

class Container extends React.Component {

  state = {
    dataSource: [],
    loading: false,
    name: TYPE[this.props.params.type].name,
    category: TYPE[this.props.params.type].category,
    tableDataUrl: TYPE[this.props.params.type].tableDataUrl,
    tableColumn: TYPE[this.props.params.type].tableColumn,
    tableQuery: TYPE[this.props.params.type].tableQuery,
    submitUrl: TYPE[this.props.params.type].submitUrl,
    searchKey: TYPE[this.props.params.type].searchKey,
    submitData: TYPE[this.props.params.type].submitData,
    limit: TYPE[this.props.params.type].limit,
    modalKey: TYPE[this.props.params.type].modalKey
  };

  handleSubmit = (e) => {
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (err) {
        notification.error({ message: '还有未填写完整的项' });
        return;
      }
      values.approvers = [values.approvers];
      if (this.state.submitData) {
        values = { ...this.state.submitData, ...values };
      }
      post(this.state.submitUrl, values).then(res => {
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
    if (this.state.modalKey == 'sns') {
      setFieldsValue({ sns: data.map(it => it.sn) });
    } else if (this.state.modalKey == 'ids') {
      setFieldsValue({ ids: data.map(it => it.id) });
    }
    setFieldsValue({ front_data: JSON.stringify(data) });
  };

  addDataSource = () => {
    addData({
      tableQuery: this.state.tableQuery,
      category: this.state.category,
      tableColumn: this.state.tableColumn,
      tableDataUrl: this.state.tableDataUrl,
      searchKey: this.state.searchKey,
      idc: this.props.idc,
      room: this.props.room,
      networkArea: this.props.networkArea,
      limit: this.state.limit || 10,
      handleDeviceSubmit: (tableData, onSuccess) => {
        const selectedRows = tableData.selectedRows || [];
        this.setFormValue(selectedRows);
        onSuccess();
      }
    });
  };
  render() {
    const { getFieldDecorator } = this.props.form;

    return (
      <Layout>
        <div style={{ marginTop: -10 }}>
          {getBreadcrumb(this.state.name)}
        </div>
        <div>
          <div className='operate_btns'>
            <Button icon='plus' onClick={this.addDataSource} type='primary'>
            添加{this.state.category}
            </Button>
          </div>
          <Form>
            <Form.Item
              label={`已选${this.state.category}`}
              {...formItemLayout_page}
            >
              {getFieldDecorator(`${this.state.modalKey}`, {
                rules: [{
                  required: true,
                  message: `请选择${this.state.category}`
                }]
              })(
                <FormTable
                  tableColumn={this.state.tableColumn}
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
              {getFieldDecorator('front_data', {
              })(
                <Input hidden={true}/>
              )}
            </Form.Item>
          
            <Form.Item
              label='审批人'
              {...formItemLayout_page}
            >
              {getFieldDecorator('approvers', {
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
          </Form>
          <Form.Item {...tailFormItemLayout_page}>
            <Button
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
    userList: state.getIn([ 'global', 'userList' ]).toJS(),
    idc: state.getIn([ 'global', 'idc' ]).toJS(),
    room: state.getIn([ 'global', 'room' ]).toJS(),
    networkArea: state.getIn([ 'global', 'networkArea' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}


export default connect(mapStateToProps, mapDispatchToProps)(Form.create()(Container));
