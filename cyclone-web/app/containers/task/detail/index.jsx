import React from 'react';
import { connect } from 'react-redux';
import Layout from 'components/layout/page-layout';
import { Row, Popover, Icon, Col, Table, Badge, Form } from 'antd';
import { renderFormDetail, getBreadcrumb } from 'common/utils';
import { TASK_CATEGORY, BUILTIN, TASK_STATUS, TASK_RATE, ARCH, OPERATION_STATUS_COLOR } from 'common/enums';
import { getColumns } from 'containers/device/common/colums';
import { hashHistory } from 'react-router';
import Crontab from 'components/crontab';

class Container extends React.Component {
  state = {
    id: this.props.params.id
  };

  componentDidMount() {
    this.reload();
  }

  reload = () => {
    this.props.dispatch({
      type: 'task-detail/detail-info/get',
      payload: { id: this.state.id }
    });
  };


  getColumns = () => {
    return [
      {
        title: '序列号',
        dataIndex: 'sn'
      },
      {
        title: '机房管理单元',
        dataIndex: 'server_room',
        render: (text) => {

          return <span>{text.name}</span>;
        }
      },
      {
        title: '机架编号',
        dataIndex: 'server_cabinet',
        render: (text) => {
          if (text && text !== 'null') {

            return <span>{text.number}</span>;
          }
        }
      },
      {
        title: '机位编号',
        dataIndex: 'server_usite',
        render: (text) => {
          if (text && text !== 'null') {
            return <span>{text.number}</span>;
          }
        }
      },
      {
        title: '内网 IP',
        dataIndex: 'intranet_ip'
      },
      {
        title: '外网 IP',
        dataIndex: 'extranet_ip'
      },
      {
        title: '系统',
        dataIndex: 'os'
      },
      {
        title: '用途',
        dataIndex: 'usage'
      },
      {
        title: '设备类型',
        dataIndex: 'category'
      },
      {
        title: '运营状态',
        dataIndex: 'operation_status',
        render: type => {
          const color = OPERATION_STATUS_COLOR[type] ? OPERATION_STATUS_COLOR[type][0] : 'transparent';
          const word = OPERATION_STATUS_COLOR[type] ? OPERATION_STATUS_COLOR[type][1] : '';
          return (
            <div>
              <Badge
                dot={true}
                style={{
                  background: color
                }}
              />{' '}
              &nbsp;&nbsp; {word}
            </div>
          );
        }
      }
    ];
  };
  render() {
    const { detailInfo, device } = this.props.data;
    const data = detailInfo.data || {};
    return (
      <Layout>
        <div className='page-header page-header-50'>
          {getBreadcrumb(this.state.id, 'task/list')}
        </div>
        <div className='page-body'>
          <div className='detail-body' style={{ marginTop: 0 }}>
            <h3 className='detail-title'>
                基本信息
            </h3>
            <div className='detail-info'>
              <Row>
                {renderFormDetail([
                  {
                    label: '任务ID',
                    value: data.id
                  },
                  {
                    label: '标题',
                    value: data.title
                  },
                  {
                    label: '类别',
                    value: TASK_CATEGORY[data.category]
                  },
                  {
                    label: '内置',
                    value: BUILTIN[data.builtin]
                  },
                  {
                    label: '状态',
                    value: TASK_STATUS[data.status]
                  },
                  {
                    label: '执行频率',
                    value: TASK_RATE[data.rate]
                  },
                  {
                    label: '创建人',
                    value: data.creator ? data.creator.name : ''
                  },
                  {
                    label: '创建时间',
                    value: data.created_at
                  },
                  {
                    label: '修改时间',
                    value: data.updated_at
                  }
                ])}
              </Row>
              {
                data.rate === 'fixed_rate' && <Row>
                  <Col span={8}>
                    <span className='panel-label'>cron表达式：</span>
                    <span className='panel-value'>{data.cron}
                      {
                        data.cron_render &&
                        <Popover
                          placement='right'
                          title={'查看Crontab表达式详情'}
                          content={
                            <Crontab
                              initialValue={{ crontab: data.cron, crontabUi: data.cron_render }}
                              form={this.props.form}
                              handleClick={() => console.log()}
                            />
                          }
                          trigger='click'
                        >
                          <Icon
                            type='info-circle-o'
                            style={{
                              fontSize: 16,
                              color: '#f6ae68',
                              marginLeft: 8
                            }}
                          />
                        </Popover>
                      }
                    </span>
                  </Col>
                </Row>
              }
            </div>
            {
              data.category === 'inspection' && <div>
                <h3 className='detail-title' style={{ marginBottom: 8 }}>
                  已选设备
                </h3>
                <Table
                  rowKey={'id'}
                  columns={this.getColumns()}
                  pagination={{ showQuickJumper: true, showSizeChanger: true }}
                  dataSource={device.data}
                  loading={device.loading}
                />
              </div>
            }
          </div>
        </div>
      </Layout>
    );
  }
}
function mapStateToProps(state) {
  return {
    data: state.get('task-detail').toJS(),
    userInfo: state.getIn([ 'global', 'userData' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Form.create()(Container));
