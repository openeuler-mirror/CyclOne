import React from 'react';
import { get } from 'common/xFetch2';
import {
  Table,
  Tooltip,
  Button,
  Pagination,
  Badge,
  notification
} from 'antd';
import { Link } from 'react-router';
import actions from '../actions';
import { INSPECTION_RESULT_COLOR, RUNNING_STTAUS } from 'common/enums';
import { post } from 'common/xFetch2';
import { getPermissonBtn } from 'common/utils';

class MyTable extends React.Component {

  reload = () => {
    this.props.dispatch({
      type: 'device-inspection-list/table-data/reload'
    });
    this.props.dispatch({
      type: 'device-inspection-list/table-data/reset'
    });
  };

  execAction = (name, record) => {
    if (actions[name]) {
      actions[name]({
        record,
        reload: () => {
          this.reload();
        }
      });
    }
  };
  renderStatus = (record, key) => {
    let type = '';
    record.result.map(data => {
      if (data.type === key) {
        type = data.status;
      }
    });
    const color = INSPECTION_RESULT_COLOR[type] ? INSPECTION_RESULT_COLOR[type][0] : 'transparent';
    const word = INSPECTION_RESULT_COLOR[type] ? INSPECTION_RESULT_COLOR[type][1] : type;
    return (
      <div>
        <Badge
          dot={true}
          style={{
            background: color
          }}
        />{' '}
        &nbsp;&nbsp;
        <Link to={`/device/inspection/detail/${record.sn}?type=${key}`}>{word}</Link>
      </div>
    );
  };

  getColumns = () => {
    return [
      {
        title: '固资编号',
        dataIndex: 'fixed_asset_number',
        width: 100,
        render: (text) => {
          return <Tooltip placement="top" title={text}>{text}</Tooltip>
        }
      },
      {
        title: '序列号',
        dataIndex: 'sn',
        width: 100,
        render: (text, record) => {
          return <Tooltip placement="top" title={text}>
            <Link to={`/device/detail/${text}`}>{text}</Link>
          </Tooltip>
        }
      },
      {
        title: '内网IP',
        dataIndex: 'intranet_ip',
        width: 100,
        render: (text) => {
          return <Tooltip placement="top" title={text}>{text}</Tooltip>
        }
      },
      {
        title: '巡检开始时间',
        dataIndex: 'start_time',
        width: 100,
        render: (text) => <Tooltip placement="top" title={text}>{text}</Tooltip>
      },
      {
        title: '巡检结束时间',
        dataIndex: 'end_time',
        width: 100,
        render: (text) => <Tooltip placement="top" title={text}>{text}</Tooltip>
      },
      {
        title: '温度',
        dataIndex: 'temperature',
        width: 100,
        render: (text, record) => this.renderStatus(record, 'temperature')
      },
      {
        title: '电压',
        dataIndex: 'voltage',
        width: 100,
        render: (text, record) => this.renderStatus(record, 'voltage')
      },
      {
        title: '风扇',
        dataIndex: 'fan',
        width: 100,
        render: (text, record) => this.renderStatus(record, 'fan')
      },
      {
        title: '内存',
        dataIndex: 'memory',
        width: 100,
        render: (text, record) => this.renderStatus(record, 'memory')
      },
      {
        title: '电源',
        dataIndex: 'power_supply',
        width: 100,
        render: (text, record) => this.renderStatus(record, 'power_supply')
      },
      {
        title: '运行状态',
        dataIndex: 'running_status',
        width: 100,
        render: (text, record) => {
          return (<span>{RUNNING_STTAUS[text]}</span>);
        }
      },
      {
        title: '健康状况',
        dataIndex: 'health_status',
        width: 100,
        render: (type, record) => {
          const color = INSPECTION_RESULT_COLOR[type] ? INSPECTION_RESULT_COLOR[type][0] : 'transparent';
          const word = INSPECTION_RESULT_COLOR[type] ? INSPECTION_RESULT_COLOR[type][1] : type;
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
      },
      {
        title: '错误信息',
        dataIndex: 'error',
        width: 100,
        render: (text) => <Tooltip placement="top" title={text}>{text}</Tooltip>
      },
      {
        title: '操作',
        dataIndex: 'operate',
        width: 100,
        render: (text, record) => {
          return <Link to={`/device/inspection/detail/${record.sn}`}>详情</Link>;
        }
      }
    ];
  };


  changePage = page => {
    this.props.dispatch({
      type: `device-inspection-list/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `device-inspection-list/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };

  getRowSelection = () => {
    const selectedRowKeys = this.props.tableData.selectedRowKeys;
    return {
      selectedRowKeys,
      onChange: (selectedRowKeys, selectedRows) => {
        this.props.dispatch({
          type: 'device-inspection-list/table-data/set/selectedRows',
          payload: {
            selectedRowKeys,
            selectedRows
          }
        });
      }
    };
  };

  inspect = (type) => {
    const { tableData } = this.props;
    const selectedRows = tableData.selectedRows || [];
    let sns = selectedRows.map(s => s.sn);
    if (type === 'all') {
      sns = [];
    } else if (selectedRows.length <= 0) {
      return notification.error({ message: '请选择设备' });
    }

    post('/api/cloudboot/v1/jobs/inspections', { sn: sns, rate: 'immediately' }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: res.message });
      this.reload();
    });
  };

  render() {
    const { tableData } = this.props;
    const { loading, pagination } = tableData;
    return (
      <div>
        <div className='operate_btns'>
          <Button
            onClick={() => this.execAction('addTask')}
            icon='plus'
            type='primary'
            style={{ marginRight: 8 }}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_inspection_addTask')}
          >
            新建巡检任务
          </Button>
          <Button.Group>
            <Button
              onClick={() => this.inspect()}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_inspection_inspect')}
            >
              重新巡检
            </Button>
            <Button
              onClick={() => this.inspect('all')}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_inspection_inspect_all')}
            >
              巡检全部
            </Button>
          </Button.Group>
        </div>
        <div className='clearfix' />
        <Table
          rowKey={'id'}
          columns={this.getColumns()}
          dataSource={tableData.list}
          loading={loading}
          pagination={false}
          defaultPageSize={3}
          rowSelection={this.getRowSelection()}
        />
        <div>
          <Pagination
            showTotal={(total) => `共 ${total} 条`}
            showQuickJumper={true}
            showSizeChanger={true}
            current={pagination.page}
            pageSize={pagination.pageSize}
            total={pagination.total}
            onShowSizeChange={this.changePageSize}
            onChange={this.changePage}
          />
        </div>
      </div>
    );
  }
}

export default MyTable;
