import React from 'react';
import { connect } from 'react-redux';
import { notification, Input, Button, Form, Select } from 'antd';
import { hashHistory } from 'react-router';
import { formItemLayout_page, tailFormItemLayout_page } from 'common/enums';
import { post } from 'common/xFetch2';
const { TextArea } = Input;
import { getPermissonBtn } from 'common/utils';

class Container extends React.Component {

  state = {
    loading: false,
    type: this.props.type //type传的是recycle
  };

  handleSubmit = (e) => {
    e.preventDefault();
    this.setFormValue(this.props.dataSource);

    this.props.form.validateFields((err, values) => {

      let isErr = false;
      values.settings.forEach(data => {
        if (!data.os_template_name) {
          isErr = true;
          return;
        }
        if (!data.hardware_template_name) {
          isErr = true;
          return;
        }
        isErr = false;
      });

      if (err || isErr) {
        return notification.error({ message: '还有未填写完整的项' });
      }

      values.approvers = [values.approvers0];
      delete values.approvers0;
      let postUrl = '/api/cloudboot/v1/approvals/devices/os-reinstallations';
      if (this.props.type === 'recycle') {
        values = { approval_type: 'device_recycle_reinstall', ...values };
        postUrl = '/api/cloudboot/v1/approvals/devices/recycle';
      }
      post(postUrl, values).then(res => {
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
    const { setFieldsValue } = this.props.form;
    setFieldsValue({ settings: data.map(it => {
      return {
        "sn": it.sn,
        install_type: it.install_type,
        hardware_template_name: it.hardware_template_name,
        os_template_name: it.os_template_name,
        need_extranet_ip: it.need_extranet_ip || 'no',
        need_intranet_ipv6: it.need_intranet_ipv6 || 'no',
        need_extranet_ipv6: it.need_extranet_ipv6 || 'no'
      };
    }) });
    setFieldsValue({ front_data: JSON.stringify(data.map(i => {
      return {
        fixed_asset_number: i.fixed_asset_number,
        sn: i.sn,
        hardware_template_name: i.hardware_template_name,
        os_template_name: i.os_template_name,
        need_extranet_ip: i.need_extranet_ip || 'no',
        need_intranet_ipv6: i.need_intranet_ipv6 || 'no',
        need_extranet_ipv6: i.need_extranet_ipv6 || 'no',        
        model: i.model,
        server_usite_number: i.server_usite.number
      };
    })) });
  };

  render() {
    const { getFieldDecorator } = this.props.form;

    return (
      <div>
        <Form>
          <Form.Item >
            {getFieldDecorator('front_data', {
            })(
              <Input hidden={true}/>
            )}
          </Form.Item>
          <Form.Item >
            {getFieldDecorator('settings', {
            })(
              <Input hidden={true}/>
            )}
          </Form.Item>
          <Form.Item
            label='备注'
            {...formItemLayout_page}
          >
            {getFieldDecorator('remark', {
            })(
              <TextArea rows={4} style={{ width: 400 }}/>
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
