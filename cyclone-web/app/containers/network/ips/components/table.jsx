import React from 'react';
import {
  Table,
  Button,
  notification,
  Pagination,
  Tooltip,
  Icon
} from 'antd';
import { YES_NO, IP_NETWORK_CATEGORY, NETWORK_IPS_CATEGORY, NETWORK_IPS_SCOPE, IP_VERSION } from 'common/enums';
import actions from '../actions';
import { Link } from 'react-router';
import TableControlCell from 'components/TableControlCell';
import { getPermissonBtn } from 'common/utils';

class MyTable extends React.Component {

  reload = () => {
    this.props.dispatch({
      type: 'network-ips/table-data/reload'
    });
    this.props.dispatch({
      type: 'network-ips/table-data/set/selectedRows',
      payload: {
        selectedRows: [],
        selectedRowKeys: []
      }
    });
  };

  //批量操作入口
  batchExecAction = (name) => {
    const { tableData } = this.props;
    const selectedRows = tableData.selectedRows || [];
    if (selectedRows.length < 1) {
      return notification.error({ message: '请至少选择一条数据' });
    }
    this.execAction(name, selectedRows);
  };

  //操作入口
  execAction = (name, records) => {
    if (actions[name]) {
      actions[name]({
        records,
        initialValue: records,
        type: name,
        room: this.props.room,
        reload: () => {
          this.reload();
        }
      });
    }
  };

  exportIP = () => {
    const query = this.props.tableData.query;
    let keys = Object.keys(query);
    const { tableData } = this.props;
    const selectedRowKeys = tableData.selectedRowKeys || [];
    keys = keys
      .map(key => {
        return `${key}=${query[key]}`;
      }) 
      .join('&');
    window.open(`/api/cloudboot/v1/ips/export?${keys}&id=${selectedRowKeys}`);
  };

  getRowSelection = () => {
    const selectedRowKeys = this.props.tableData.selectedRowKeys;
    return {
      selectedRowKeys,
      onChange: (selectedRowKeys, selectedRows) => {
        this.props.dispatch({
          type: 'network-ips/table-data/set/selectedRows',
          payload: {
            selectedRowKeys,
            selectedRows
          }
        });
      }
    };
  };

  changePage = page => {
    this.props.dispatch({
      type: `network-ips/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `network-ips/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };


  getColumns = () => {
    return [
      {
        title: 'IP地址',
        dataIndex: 'ip'
      },
      {
        title: '网段名称',
        dataIndex: 'ip_network',
        render: (text, record) => <a onClick={() => this.execAction('cidr_detail', { id: record.ip_network.id })}>{text.cidr}</a>
      },
      {
        title: '网段网关',
        dataIndex: 'gateway',
        render: (text, record) => <span>{record.ip_network.gateway}</span>
      },
      {
        title: '网段掩码',
        dataIndex: 'netmask',
        render: (text, record) => <span>{record.ip_network.netmask}</span>
      },
      {
        title: '网段类别',
        dataIndex: 'ip_network',
        render: (text, record) => {
          return <span>{IP_NETWORK_CATEGORY[record.ip_network.category]}</span>;
        }
      },
      {
        title: 'IP版本',
        dataIndex: 'ip_network',
        render: (text, record) => {
          return <span>{IP_VERSION[record.ip_network.version]}</span>;
        }
      },               
      {
        title: '是否被使用',
        dataIndex: 'is_used',
        render: (text) => <span className={`yes_no_status ${text === 'yes' ? 'yes_status' : 'no_status'}`}>{YES_NO[text]}</span>
      },
      {
        title: 'IP作用范围',
        dataIndex: 'scope',
        render: (text) => <span>{NETWORK_IPS_SCOPE[text]}</span>
      },
      {
        title: 'IP用途',
        dataIndex: 'category',
        render: (text) => <span>{NETWORK_IPS_CATEGORY[text]}</span>
      },      
      {
        title: '关联设备序列号',
        dataIndex: 'sn',
        render: (text, record) => <Link to={`/device/detail/${text}`}>{text}</Link>
      },
      {
        title: '固资编号',
        dataIndex: 'fixed_asset_number'
      },
      {
        title: '更新时间',
        dataIndex: 'updated_at'
      },
      {
        title: '操作',
        dataIndex: 'operate',
        width: 150,
        render: (text, record) => {
          const commands = [
            {
              name: '分配IP',
              command: 'assign',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_ip_assign') || (record.is_used === 'yes' || record.is_used === 'disabled')
            },
            {
              name: '禁用',
              command: '_disabled',
              type: 'danger',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_ip_unassign') || (record.is_used === 'yes' || record.is_used === 'disabled')
            }
          ];
          return (
            <TableControlCell
              commands={commands}
              record={record}
              execCommand={command => {
                this.execAction(command, [record]);
              }}
            />
          );
        }
      }
    ];
  };


  render() {
    const { tableData } = this.props;
    const { loading, pagination } = tableData;
    return (
      <div>
        <span className='pull-right'>
          <Button
            onClick={this.exportIP}
            style={{ float: 'right', marginRight: 30 }}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_network_device_create')}
          >
            导出
          </Button>          
          <Button
            onClick={() => this.batchExecAction('_disabled')}
            type='danger'
            style={{ marginBottom: 8  }}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_network_device_create')}
          >
            批量禁用
          </Button>
        </span>
        <Tooltip placement="top"  title='为选中的设备自动分配一个内网或外网IP'>
        <Button
          onClick={() => this.execAction('assignipv4')}
          type='primary'
          style={{ marginBottom: 8 }}
          disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_network_device_create')}
        >
          分配IPv4<Icon type='question-circle-o' />
        </Button>
        </Tooltip>
        <Tooltip placement="top"  title='为选中的设备自动分配一个内网或外网IP'>
        <Button
          onClick={() => this.execAction('assignipv6')}
          type='primary'
          style={{ marginLeft: 18 }}
          disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_network_device_create')}
        >
          分配IPv6<Icon type='question-circle-o' />
        </Button>
        </Tooltip>
        <Table
          rowKey={'id'}
          columns={this.getColumns()}
          pagination={false}
          dataSource={tableData.list}
          rowSelection={this.getRowSelection()}
          loading={loading}
        />
        <div>
          <Pagination
            showQuickJumper={true}
            showSizeChanger={true}
            current={pagination.page}
            pageSize={pagination.pageSize}
            total={pagination.total}
            onShowSizeChange={this.changePageSize}
            onChange={this.changePage}
            showTotal={(total) => `共 ${total} 条`}
          />
        </div>
      </div>
    );
  }
}

export default MyTable;
