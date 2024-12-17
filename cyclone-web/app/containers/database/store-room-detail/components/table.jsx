import React from 'react';
import {
  Table,
  Button,
  Pagination,
  Badge
} from 'antd';
import actions from '../actions';
import TableControlCell from 'components/TableControlCell';
import { CAB_STATUS_COLOR } from "common/enums";
import { getPermissonBtn } from 'common/utils';

class MyTable extends React.Component {

  //操作入口
  execAction = (name, records) => {
    if (actions[name]) {
      actions[name]({
        records,
        id: this.props.id,
        reload: () => {
          this.props.reload();
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
          type: 'database-store-detail/table-data/set/selectedRows',
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
      type: `database-store-detail/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `database-store-detail/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };


  getColumns = () => {
    return [
      {
        title: '货架编号',
        dataIndex: 'number'
      },
      {
        title: '创建者',
        dataIndex: 'creator'
      },
      {
        title: '创建时间',
        dataIndex: 'created_at'
      },
      {
        title: '修改时间',
        dataIndex: 'updated_at'
      },
      {
        title: '备注',
        dataIndex: 'remark'
      },
      {
        title: '状态',
        dataIndex: 'status',
        width: 100,
        render: type => {
          const color = CAB_STATUS_COLOR[type] ? CAB_STATUS_COLOR[type][0] : 'transparent';
          const word = CAB_STATUS_COLOR[type] ? CAB_STATUS_COLOR[type][1] : '';
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
              name: '删除',
              command: '_delete',
              type: 'danger',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_virtual_cabinet_delete')
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
            style={{ marginTop: 8 }}
            icon='plus'
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_virtual_cabinet_create')}
          >
            新增
          </Button>
        </div>
        <Table
          rowKey={'id'}
          columns={this.getColumns()}
          pagination={false}
          dataSource={tableData.list}
          // rowSelection={this.getRowSelection()}
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
