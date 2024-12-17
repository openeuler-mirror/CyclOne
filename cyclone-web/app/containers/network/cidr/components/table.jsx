import React from 'react';
import {
  Table,
  Tooltip,
  Icon,
  Button,
  Pagination,
  notification,
  Typography
} from 'antd';
import actions from '../actions';
import TableControlCell from 'components/TableControlCell';
import { IP_NETWORK_CATEGORY, IP_VERSION } from 'common/enums';
import { getPermissonBtn } from 'common/utils';


class MyTable extends React.Component {

  reload = () => {
    this.props.dispatch({
      type: 'network-cidr/table-data/reload'
    });
    this.props.dispatch({
      type: 'network-cidr/table-data/set/selectedRows',
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
        device: this.props.device,
        reload: () => {
          this.reload();
        }
      });
    }
  };


  getRowSelection = () => {
    const selectedRowKeys = this.props.tableData.selectedRowKeys;
    return {
      selectedRowKeys,
      onChange: (selectedRowKeys, selectedRows) => {
        this.props.dispatch({
          type: 'network-cidr/table-data/set/selectedRows',
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
      type: `network-cidr/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `network-cidr/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };


  getColumns = () => {
    return [
      {
        title: '网段名称',
        dataIndex: 'cidr',
        render: (text, record) => <a onClick={() => this.execAction('detail', record)}>{text}</a>
      },
      {
        title: '网段掩码',
        dataIndex: 'netmask'
      },
      {
        title: '网段网关',
        dataIndex: 'gateway'
      },
      {
        title: 'IP资源池',
        dataIndex: 'ip_pool',
        render: (text) => <Tooltip placement="top" title={text}>{text}</Tooltip>
      },
      {
        title: 'PXE资源池',
        dataIndex: 'pxe_pool',
        render: (text) => <Tooltip placement="top" title={text}>{text}</Tooltip>
      },
      {
        title: '机房管理单元',
        dataIndex: 'server_room',
        render: (text, record) => <Tooltip placement="top" title={text.name}>
            <a onClick={() => this.execAction('room_detail', { id: text.id })}>
              {text.name}
            </a>
          </Tooltip>
      },
      {
        title: '网络区域',
        dataIndex: 'network_area'
      },
      {
        title: '网段类别',
        dataIndex: 'category',
        render: (text) => <span>{IP_NETWORK_CATEGORY[text]}</span>
      },
      {
        title: '覆盖交换机',
        dataIndex: 'switchs',
        render: (text) => <Tooltip placement="top" title={(text || []).map(it => it.name + '(' + it.fixed_asset_number + ')').join('+')}>
          {(text || []).map(it => it.fixed_asset_number).join('+')}
        </Tooltip>
      },
      {
        title: 'VLAN',
        dataIndex: 'vlan'
      },
      {
        title: 'IP版本',
        dataIndex: 'version',
        render: (text) => <span>{IP_VERSION[text]}</span>  
      },    
      {
        title: '操作',
        dataIndex: 'operate',
        width: 100,
        render: (text, record) => {
          const commands = [
            {
              name: '编辑',
              command: '_update',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_ip_network_update')
            },
            {
              name: '删除',
              command: '_delete',
              type: 'danger',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_ip_network_delete')
            }
          ];
          return (
            <TableControlCell
              commands={commands}
              record={record}
              execCommand={command => {
                this.execAction(command, record);
              }}
            />
          );
        }
      }
    ];
  };

  //下载导入模板
  downloadImportTemplate = () => {
    window.open('assets/files/network_segment_import.xlsx');
  };

  render() {
    const { tableData } = this.props;
    const { loading, pagination } = tableData;
    return (
      <div>
        <div className='operate_btns'>
          <Button
            onClick={() => this.execAction('_create')}
            type='primary'
            style={{ marginRight: 8 }}
            icon='plus'
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_ip_network_create')}
          >
            新增
          </Button>        
          <span className='pull-right'>
            <Button.Group style={{ marginRight: 8 }}>
              <Button
                onClick={() => this.downloadImportTemplate()}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_download')}
              >
                下载导入模板
              </Button>
              <Button
                onClick={() => this.execAction('_import')}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_import')}
              >
                导入
              </Button>
              <Button
                onClick={() => this.batchExecAction('_batchdelete')}
                type='danger'
                style={{ marginLeft: 8 }}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_ip_network_delete')}
              >
                删除
              </Button>           
            </Button.Group>
          </span>
        </div>
        <div className='clearfix' />
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
