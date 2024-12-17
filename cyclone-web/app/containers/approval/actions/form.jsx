import React from 'react';
import { renderFormDetail } from 'common/utils';
import { APPROVAL_STATUS, APPROVAL_TYPE, APPROVAL_ACTION } from 'common/enums';
import { get, put, del } from 'common/xFetch2';
import { Form, Input, notification, Timeline, Icon, Button } from 'antd';
const { TextArea } = Input;
const FormItem = Form.Item;
import InfoTable from './info_table';
import { getPermissonBtn } from 'common/utils';

class MyForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: {}
    };
  }

  getData = () => {
    get(`/api/cloudboot/v1/approvals/${this.props.approval_id}`).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      this.setState({ data: res.content || {} });
    });
  };

  componentDidMount() {
    this.getData();
  }

  revoke = () => {
    del(`/api/cloudboot/v1/approvals/${this.props.approval_id}`).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: res.message });
      this.getData();
      this.props.reload();
      // this.props.onCancel();
    });
  };

  approval = (action, step) => {
    this.props.form.validateFields((err, values) => {
      const data = {
        approval_id: this.props.approval_id,
        approval_step_id: values.approval_step_id,
        remark: values.remark,
        action: action
      };
      put(`/api/cloudboot/v1/approvals/${this.props.approval_id}/step/${values.approval_step_id}`, data).then(res => {
        if (res.status !== 'success') {
          return notification.error({ message: res.message });
        }
        notification.success({ message: res.message });
        this.getData();
        this.props.reload();
        // this.props.onCancel();
      });
    });
  };

  render() {
    const { getFieldDecorator } = this.props.form;

    const data = this.state.data;

    let ifFinished = false;
    let ifApproval = false;
    let ifRevoke = false;
    let ifApply = false;
    let currentStep = 0;

    let step0 = {};
    let step1 = {};

    if (data.steps) {

      step0 = data.steps[0];

      //有实施的步骤
      if (data.steps.length > 1) {
        step1 = data.steps[1];
        ifApply = true;
      }

      //判断审批的最后一个步骤是否完成
      if (data.steps[data.steps.length - 1 ].end_time) {
        ifFinished = true;
      }

      //是否可以撤销审批（还未开始审批的单且是本人）
      if (!step0.end_time && data.initiator && data.initiator.id === this.props.userInfo.id) {
        ifRevoke = true;
      }

      //判断当前步骤
      if (!step0.end_time) {
        currentStep = 0;
        //审批人与登录用户相同则可以审批
        if (step0.approver && step0.approver.id === this.props.userInfo.id) {
          ifApproval = true;
        }
      }
      if (data.steps.length > 1 && step0.end_time && !step1.end_time) {
        currentStep = 1;
        //实施人与登录用户相同则可以审批
        if (step1.approver && step1.approver.id === this.props.userInfo.id) {
          ifApproval = true;
        }
      }
    }

    return (
      <div>
        <h3 className='detail-title detail-title-drawer'>申请信息</h3>
        <div className='detail-info detail-info-24'>
          {renderFormDetail([
            {
              label: '申请类型',
              value: APPROVAL_TYPE[data.type]
            },
            {
              label: '审批状态',
              value: APPROVAL_STATUS[data.status] && APPROVAL_STATUS[data.status]
            },
            {
              label: '备注',
              value: data.remark
            }
          ], 12)}
          <InfoTable type={data.type} data={JSON.parse(data.front_data || '[]')}/>
        </div>
        <h3 className='detail-title detail-title-drawer'>审批历史</h3>
        <div className='detail-info detail-info-24'>
          <Timeline>
            <Timeline.Item color='blue'>
              <p>申请</p>
              <p>
                <span><Icon type='user' />{data.initiator && data.initiator.name}</span>
                <span style={{ marginLeft: 8 }}><Icon type='clock-circle' />{data.start_time}</span>
              </p>
            </Timeline.Item>

            <Timeline.Item color='blue'>
              <p>审批
                {!step0.end_time ? <span className='approval_yellow'>待审批</span> :
                <span className='approval_yellow' style={{ color: step0.action === 'agree' ? '#6bc646' : '#ff3700' }}>{APPROVAL_ACTION[step0.action]}</span>
                }
              </p>
              <p>
                <span><Icon type='user' />{step0.approver && step0.approver.name}</span>
                {
                  step0.end_time &&
                  <span style={{ marginLeft: 8 }}><Icon type='clock-circle' />{step0.end_time}</span>
                }
              </p>
              <p>{step0.remark}</p>
              {
                ifApproval && data.status !== 'revoked' && currentStep === 0 &&
                <Form>
                  <FormItem >
                    {getFieldDecorator(`remark`, {
                    })(
                      <TextArea rows={4} placeholder='备注'/>
                    )}
                  </FormItem>
                  <FormItem>
                    {getFieldDecorator(`approval_step_id`, {
                      initialValue: step0.id
                    })(
                      <Input hidden={true}/>
                    )}
                  </FormItem>
                </Form>
              }
            </Timeline.Item>
            {
              ifApply &&
              <Timeline.Item color='blue'>
                <p>实施
                  {!step1.end_time ? <span className='approval_yellow'>待实施</span> :
                  <span className='approval_yellow' style={{ color: step1.action === 'agree' ? '#6bc646' : '#ff3700' }}>{APPROVAL_ACTION[step1.action]}</span>
                  }
                </p>
                <p>
                  <span><Icon type='user' />{step1.approver && step1.approver.name}</span>
                  {
                    step1.end_time &&
                    <span style={{ marginLeft: 8 }}><Icon type='clock-circle' />{step1.end_time}</span>
                  }
                </p>
                <p>{step1.remark}</p>
                {
                  ifApproval && data.status !== 'revoked' && data.is_rejected !== 'yes' && currentStep === 1 &&
                  <Form>
                    <FormItem >
                      {getFieldDecorator(`remark`, {
                      })(
                        <TextArea rows={4} placeholder='备注'/>
                      )}
                    </FormItem>
                    <FormItem>
                      {getFieldDecorator(`approval_step_id`, {
                        initialValue: step1.id
                      })(
                        <Input hidden={true}/>
                      )}
                    </FormItem>
                  </Form>
                }
              </Timeline.Item>
            }
            {
              ifFinished && <Timeline.Item color='blue'>
                <p>{APPROVAL_STATUS[data.status] && APPROVAL_STATUS[data.status]}</p>
              </Timeline.Item>
            }
            {
              data.status === 'revoked' &&
              <Timeline.Item color='red'>
                <p>{APPROVAL_STATUS[data.status] && APPROVAL_STATUS[data.status]}</p>
                <p><Icon type='clock-circle' />{data.end_time}</p>
              </Timeline.Item>
            }
          </Timeline>
        </div>
        {
          data.status !== 'revoked' &&
          <div className='drawer-button'>
            {ifApproval && data.is_rejected !== 'yes' && <Button
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_approval_agree')}
              onClick={() => this.approval('agree', currentStep)} type='primary' style={{ marginRight: 8 }}
            >通过</Button>}
            {ifApproval && data.is_rejected !== 'yes' && <Button
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_approval_disagree')}
              onClick={() => this.approval('reject', currentStep)} style={{ marginRight: 8 }}
            >不通过</Button>}
            {ifRevoke && <Button
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_approval_revoke')}
              onClick={this.revoke}
            >撤销</Button>}
          </div>
        }
      </div>
    );
  }
}

export default Form.create()(MyForm);
