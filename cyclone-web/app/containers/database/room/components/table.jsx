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
import TableControlCell from 'components/TableControlCell';
import actions from '../actions';
import { IDC_STATUS_COLOR } from "common/enums";
import { renderDisplayMore } from 'common/utils';
const plainOptions = [
  { value: 'vendor_manager', name: '供应商负责人' },
  { value: 'network_asset_manager', name: '网络资产负责人' },
  { value: 'support_phone_number', name: '7*24小时保障电话' }
];
import { getPermissonBtn } from 'common/utils';

class MyTable extends React.Component {

  state = {
    checkedList: [],
    indeterminate: false,
    checkAll: false
  };

  reload = () => {
    this.props.dispatch({
      type: 'database-room/table-data/reload'
    });
    this.props.dispatch({
      type: 'database-room/table-data/set/selectedRows',
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
    if (selectedRows.length <= 0) {
      return notification.error({ message: '请至少选择一条数据' });
    }
    this.execAction(name, selectedRows);
  };

  //单个操作入口
  oneExecAction = (name) => {
    const { tableData } = this.props;
    const selectedRows = tableData.selectedRows || [];
    if (selectedRows.length !== 1) {
      return notification.error({ message: '请选择一条数据' });
    }
    this.execAction(name, selectedRows[0]);
  };

  //操作入口
  execAction = (name, records) => {
    if (actions[name]) {
      actions[name]({
        records,
        initialValue: records,
        type: name,
        idc: this.props.idc,
        reload: () => {
          this.reload();
        }
      });
    }
  };

  //下载导入模板
  downloadImportTemplate = () => {
    window.open('assets/files/server-rooms.xlsx');
  };


  getRowSelection = () => {
    const selectedRowKeys = this.props.tableData.selectedRowKeys;
    return {
      selectedRowKeys,
      onChange: (selectedRowKeys, selectedRows) => {
        this.props.dispatch({
          type: 'database-room/table-data/set/selectedRows',
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
      type: `database-room/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `database-room/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };


  getColumns = () => {
    const columns = [
      {
        title: '机房管理单元',
        dataIndex: 'name',
        render: (text, record) => <a onClick={() => this.execAction('detail', record)}>{text}</a>
      },
      {
        title: '数据中心',
        dataIndex: 'idc',
        render: (text, record) => <a onClick={() => this.execAction('idc_detail', { id: text.id })}>{text.name}</a>
      },
      {
        title: '所属一级机房',
        dataIndex: 'first_server_room',
        render: (text, record) => <span>{text.name}</span>
      },
      {
        title: '状态',
        dataIndex: 'status',
        width: 90,
        render: type => {
          const color = IDC_STATUS_COLOR[type] ? IDC_STATUS_COLOR[type][0] : 'transparent';
          const word = IDC_STATUS_COLOR[type] ? IDC_STATUS_COLOR[type][1] : '';
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
        title: '城市',
        dataIndex: 'city'
      },
      {
        title: '地址',
        dataIndex: 'address'
      },
      {
        title: '机架数',
        dataIndex: 'cabinet_count'
      },
      {
        title: '机房负责人',
        dataIndex: 'server_room_manager'
      },
      {
        title: '供应商负责人',
        dataIndex: 'vendor_manager'
      },
      {
        title: '网络资产负责人',
        dataIndex: 'network_asset_manager'
      },
      {
        title: '7*24小时保障电话',
        dataIndex: 'support_phone_number'
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
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_server_room_update')
            },
            {
              name: '删除',
              command: '_delete',
              type: 'danger',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_server_room_delete')
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
    // if (this.state.checkedList.length > 0) {
    //   this.state.checkedList.forEach(data => {
    //     columns.push({
    //       title: data.name,
    //       dataIndex: data.value
    //     });
    //   });
    // }
    return columns;
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
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_room_create')}
          >
            新增
          </Button>
          <Button.Group style={{ marginRight: 8 }}>
            <Button
              onClick={() => this.batchExecAction('accepted')}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_room_accepted')}
            >
              验收
            </Button>
            <Button
              onClick={() => this.batchExecAction('production')}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_room_production')}
            >
              投产
            </Button>
            {/*<Button*/}
            {/*onClick={() => this.batchExecAction('abolished')}*/}
            {/*disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_room_abolished')}*/}
            {/*>*/}
            {/*裁撤*/}
            {/*</Button>*/}
          </Button.Group>
          <span>
            已选 { selectedRows.length } 项
          </span>
          <span className='pull-right'>
            <ButtonGroup style={{ marginRight: 8 }}>
              <Button
                onClick={() => this.downloadImportTemplate()}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_room_download')}
              >
                下载导入模板
              </Button>
              <Button
                onClick={() => this.execAction('_import')}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_room_import')}
              >
                导入
              </Button>
            </ButtonGroup>
            {/*{renderDisplayMore(this, plainOptions)}*/}
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
