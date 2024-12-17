import React from 'react';
import { get } from 'common/xFetch2';
import {
  Table,
  Button,
  Pagination,
  notification,
  Badge
} from 'antd';
const ButtonGroup = Button.Group;
import actions from '../actions';
import TableControlCell from 'components/TableControlCell';
import { NET_STATUS_COLOR } from "common/enums";
import { getPermissonBtn } from 'common/utils';

class MyTable extends React.Component {

  reload = () => {
    this.props.dispatch({
      type: 'database-network/table-data/reload'
    });
    this.props.dispatch({
      type: 'database-network/table-data/set/selectedRows',
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

  //下载导入模板
  downloadImportTemplate = () => {
    window.open('assets/files/network-areas.xlsx');
  };


  getRowSelection = () => {
    const selectedRowKeys = this.props.tableData.selectedRowKeys;
    return {
      selectedRowKeys,
      onChange: (selectedRowKeys, selectedRows) => {
        this.props.dispatch({
          type: 'database-network/table-data/set/selectedRows',
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
      type: `database-network/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `database-network/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };


  getColumns = () => {
    return [
      {
        title: '网络区域名称',
        dataIndex: 'name',
        render: (text, record) => <a onClick={() => this.execAction('detail', record)}>{text}</a>
      },
      {
        title: '机房管理单元',
        dataIndex: 'server_room',
        render: (text, record) => <a onClick={() => this.execAction('room_detail', { id: text.id })}>{text.name}</a>
      },
      {
        title: '关联物理区域',
        dataIndex: 'physical_area',
        render: (text) => <span>{(text || []).map(it => it.name).join('，')}</span>
      },
      {
        title: '状态',
        dataIndex: 'status',
        render: type => {
          const color = NET_STATUS_COLOR[type] ? NET_STATUS_COLOR[type][0] : 'transparent';
          const word = NET_STATUS_COLOR[type] ? NET_STATUS_COLOR[type][1] : '';
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
        title: '操作',
        dataIndex: 'operate',
        render: (text, record) => {
          const commands = [
            {
              name: '编辑',
              command: '_update',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_network_area_update')
            },
            {
              name: '删除',
              command: '_delete',
              type: 'danger',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_network_area_delete')
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


  render() {
    const { tableData } = this.props;
    const { loading, pagination, selectedRows } = tableData;
    return (
      <div>
        <div className='operate_btns'>
          <Button
            onClick={() => this.execAction('_create')}
            type='primary'
            style={{ marginRight: 8 }}
            icon='plus'
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_network_area_create')}
          >
            新增
          </Button>
          <Button.Group style={{ marginRight: 8 }}>
            <Button
              onClick={() => this.batchExecAction('production')}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_network_area_production')}
            >
              投产
            </Button>
            {/*<Button*/}
            {/*onClick={() => this.batchExecAction('offline')}*/}
            {/*disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_network_area_offline')}*/}
            {/*>*/}
            {/*下线*/}
            {/*</Button>*/}
          </Button.Group>
          <span>
            已选 { selectedRows.length } 项
          </span>
          <span className='pull-right'>
            <ButtonGroup style={{ marginRight: 8 }}>
              <Button
                onClick={() => this.downloadImportTemplate()}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_network_area_download')}
              >
                下载导入模板
              </Button>
              <Button
                onClick={() => this.execAction('_import')}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_network_area_import')}
              >
                导入
              </Button>
            </ButtonGroup>
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
