import React from 'react';
import { get } from 'common/xFetch2';
import {
  Table,
  Button,
  Pagination,
  notification,
  Badge
} from 'antd';
import actions from '../actions';
import TableControlCell from 'components/TableControlCell';
import { IDC_USAGE, IDC_STATUS_COLOR } from "common/enums";
import { getPermissonBtn } from 'common/utils';
import { Link } from 'react-router';
const ButtonGroup = Button.Group;


class MyTable extends React.Component {

  //下载导入模板
  downloadImportTemplate = () => {
    window.open('assets/files/store_room_import.xlsx');
  };


  reload = () => {
    this.props.dispatch({
      type: 'database-store/table-data/reload'
    });
    this.props.dispatch({
      type: 'database-store/table-data/set/selectedRows',
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
        idc: this.props.idc,
        type: name,
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
          type: 'database-store/table-data/set/selectedRows',
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
      type: `database-store/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `database-store/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };


  getColumns = () => {

    return [
      {
        title: '库房管理单元',
        dataIndex: 'name',
        render: (text, record) => <Link to={`/database/store-room/${record.id}`}>{text}</Link>
      },
      {
        title: '数据中心',
        dataIndex: 'idc',
        render: (text, record) => <a onClick={() => this.execAction('idc_detail', { id: text.id })}>{text.name}</a>
      },
      {
        title: '一级机房',
        dataIndex: 'first_server_room',
        render: (text) => <span>{text ? text.name : ''}</span>
      },
      // {
      //   title: '状态',
      //   dataIndex: 'status',
      //   render: type => {
      //     const color = IDC_STATUS_COLOR[type] ? IDC_STATUS_COLOR[type][0] : 'transparent';
      //     const word = IDC_STATUS_COLOR[type] ? IDC_STATUS_COLOR[type][1] : '';
      //     return (
      //       <div>
      //         <Badge
      //           dot={true}
      //           style={{
      //             background: color
      //           }}
      //         />{' '}
      //         &nbsp;&nbsp; {word}
      //       </div>
      //     );
      //   }
      // },
      {
        title: '城市',
        dataIndex: 'city'
      },
      {
        title: '地址',
        dataIndex: 'address'
      },
      {
        title: '虚拟货架数',
        dataIndex: 'cabinet_count',
        render: (text, record) => <Link to={`/database/store-room/${record.id}`}>{text}</Link>
      },
      {
        title: '库房负责人',
        dataIndex: 'store_room_manager'
      },
      {
        title: '供应商负责人',
        dataIndex: 'vendor_manager'
      },
      {
        title: '操作',
        dataIndex: 'operate',
        render: (text, record) => {
          const commands = [
            {
              name: '编辑',
              command: '_update',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_store_room_update')
            },
            {
              name: '删除',
              command: '_delete',
              type: 'danger',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_store_room_delete')
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
    const { loading, pagination } = tableData;
    return (
      <div>
        <div className='operate_btns'>
          <Button
            onClick={() => this.execAction('_create')}
            type='primary'
            style={{ marginRight: 8 }}
            icon='plus'
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_store_room_create')}
          >
            新增
          </Button>
          <span className='pull-right'>
            <ButtonGroup style={{ marginRight: 8 }}>
              <Button
                onClick={() => this.downloadImportTemplate()}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_cabinet_download')}
              >
                下载导入模板
              </Button>
              <Button
                onClick={() => this.execAction('_import')}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_store_room_import')}
              >
                导入
              </Button>
            </ButtonGroup>
          </span>
        </div>
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
