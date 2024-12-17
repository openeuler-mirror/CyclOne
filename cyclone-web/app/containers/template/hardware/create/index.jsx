import React from 'react';
import { connect } from 'react-redux';
import Layout from 'components/layout/page-layout';
import { Row, Col, Button, Form, Input, Select, notification, Spin, Alert } from 'antd';
import { renderFormDetail, getBreadcrumb } from 'common/utils';
import { PAGE_TYPE } from 'common/enums';
const FormItem = Form.Item;
const Option = Select.Option;
import { post, put } from 'common/xFetch2';
import Box from './component/box';
import { hashHistory } from 'react-router';
import { getFormData } from './component/data';

class Container extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      cards: null
    };
  }

  componentDidMount() {
    const { id } = this.props.params;
    //获取详情
    if (id === 'new') {
      this.props.dispatch({
        type: 'hardware/template/create'
      });
    } else {
      this.props.dispatch({
        type: 'hardware/template/get',
        payload: id
      });
    }

    // //获取待固件升级包
    // this.props.dispatch({
    //   type: 'hardware/firmwares/get'
    // });
    // //获取待固件类型
    // this.props.dispatch({
    //   type: 'hardware/dictionaries/get'
    // });
  }

  componentWillUnmount() {
    this.props.dispatch({
      type: 'hardware/template/clear'
    });
  }

  reload = () => {
    hashHistory.push('template/hardware/list');
  };

  //保存
  handleSubmit = () => {
    this.props.form.validateFields((err, values) => {
      if (err) {
        return notification.error({ message: '还有信息未填写完整' });
      }
      const data = getFormData(this.state.cards, values);
      const postData = {
        vendor: values.vendor,
        model_name: values.model_name,
        name: values.name,
        data: data
      };
      //新增 编辑
      const { params, location } = this.props;
      if (location.pathname.split('/')[3] === 'create') {
        post('/api/cloudboot/v1/hardware-templates', postData).then(res => {
          if (res.status === 'success') {
            notification.success({ message: '新增成功' });
            this.reload();
          } else {
            notification.error({ message: res.message });
          }
        });
      } else {
        put(`/api/cloudboot/v1/hardware-templates/${params.id}`, postData).then(res => {
          if (res.status === 'success') {
            notification.success({ message: '更新成功' });
            this.reload();
          } else {
            notification.error({ message: res.message });
          }
        });
      }

    });
  };

  getOrder = (cards) => {
    this.setState({ cards });
  };

  render() {
    const { template, firmwares, dictionaries, location } = this.props;
    const { data, loading } = template;
    const { getFieldDecorator } = this.props.form;
    const pageType = location.pathname.split('/')[3];
    return (
      <Layout>
        {getBreadcrumb(PAGE_TYPE[pageType])}
        <div className='detail-body'>
          <div className='detail-info'>
            <Row>
              <Form>
                <Row gutter={24}>
                  <Col span={8}>
                    <FormItem label={'名称'}>
                      {getFieldDecorator('name', {
                        initialValue: data.name,
                        rules: [{
                          required: true,
                          message: '名称不能为空'
                        }]
                      })(
                        <Input disabled={pageType === 'detail'} />
                      )}
                    </FormItem>
                  </Col>
                  {/*<Col span={8}>*/}
                  {/*<FormItem label={'厂商'}>*/}
                  {/*{getFieldDecorator('vendor', {*/}
                  {/*initialValue: data.vendor,*/}
                  {/*rules: [{*/}
                  {/*required: true,*/}
                  {/*message: '厂商不能为空'*/}
                  {/*}]*/}
                  {/*})(*/}
                  {/*<Input disabled={pageType === 'detail'} />*/}
                  {/*)}*/}
                  {/*</FormItem>*/}
                  {/*</Col>*/}
                  {/*<Col span={8}>*/}
                  {/*<FormItem label={'产品型号'}>*/}
                  {/*{getFieldDecorator('model_name', {*/}
                  {/*initialValue: data.model_name,*/}
                  {/*rules: [{*/}
                  {/*required: true,*/}
                  {/*message: '产品型号不能为空'*/}
                  {/*}]*/}
                  {/*})(*/}
                  {/*<Input disabled={pageType === 'detail'} />*/}
                  {/*)}*/}
                  {/*</FormItem>*/}
                  {/*</Col>*/}
                </Row>
              </Form>
            </Row>
          </div>
          <h3 style={{ marginBottom: 4 }}>配置项</h3>
          {
            pageType !== 'detail' &&
            <Alert
              message='拖动序号图标可进行排序'
              type='info'
              showIcon={true}
              closable={true}
            />
          }
          <div className='configure-panel' style={{ marginTop: 8 }}>
            {
              !loading ?
                <Box
                  form={this.props.form}
                  getOrder={this.getOrder}
                  firmwares={firmwares.data}
                  dictionaries={dictionaries.data}
                  cards={data.data}
                  disabled={pageType === 'detail'}
                /> :
                <Spin />
            }
          </div>
          <div className='pull-right'>
            <Button onClick={this.reload}>取消</Button>
            { pageType !== 'detail' && <Button style={{ marginLeft: 8 }} type='primary' onClick={this.handleSubmit}>保存</Button> }
          </div>
        </div>
      </Layout>
    );
  }
}
function mapStateToProps(state) {
  return {
    firmwares: state.get('hardware-template-create').toJS().firmwares,
    dictionaries: state.get('hardware-template-create').toJS().dictionaries,
    template: state.get('hardware-template-create').toJS().template,
    userInfo: state.getIn([ 'global', 'userData' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Form.create()(Container));
