import React from 'react';
import { get } from 'common/xFetch2';
import {
  Table,
  Pagination,
  Badge
} from 'antd';
import TableControlCell from 'components/TableControlCell';
import actions from '../actions';
import { TASK_STATUS_COLOR, TASK_RATE, TASK_CATEGORY, BUILTIN } from 'common/enums';
import { Link } from 'react-router';
import { getPermissonBtn } from 'common/utils';

class MyTable extends React.Component {

  reload = () => {
    this.props.dispatch({
      type: 'task-list/table-data/reload'
    });
  };

  execAction = (name, record) => {
    if (actions[name]) {
      actions[name]({
        record,
        type: name,
        reload: () => {
          this.reload();
        }
      });
    }
  };


  getColumns = () => {
    return [
      {
        title: '标题',
        dataIndex: 'title',
        render: (t, record) => <Link to={`task/detail/${record.id}`}>{t}</Link>
      },
      {
        title: '是否内置',
        dataIndex: 'builtin',
        render: (text) => <span className={`yes_no_status ${text === 'yes' ? 'yes_status' : 'no_status'}`}>{BUILTIN[text]}</span>
      },
      {
        title: '类别',
        dataIndex: 'category',
        render: (t) => <span>{TASK_CATEGORY[t]}</span>
      },
      {
        title: '执行频率',
        dataIndex: 'rate',
        render: (t) => <span>{TASK_RATE[t]}</span>
      },
      {
        title: '状态',
        dataIndex: 'status',
        render: (type, record) => {
          const color = TASK_STATUS_COLOR[type] ? TASK_STATUS_COLOR[type][0] : 'transparent';
          const word = TASK_STATUS_COLOR[type] ? TASK_STATUS_COLOR[type][1] : type;
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
        title: '创建人',
        dataIndex: 'creator',
        render: (t) => <span>{t.name}</span>
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
        title: '操作',
        dataIndex: 'manufacturer',
        render: (text, record) => {
          let commands = [
            {
              name: '暂停',
              command: '_pause',
              disabled: record.rate === 'immediately' || record.status === 'paused' || !getPermissonBtn(this.props.userInfo.permissions, 'button_task_pause')
            },
            {
              name: '继续',
              command: '_continue',
              disabled: record.status !== 'paused' || !getPermissonBtn(this.props.userInfo.permissions, 'button_task_continue')
            },
            {
              name: '删除',
              command: '_delete',
              type: 'danger',
              disabled: record.builtin === 'yes' || !getPermissonBtn(this.props.userInfo.permissions, 'button_task_delete')
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


  changePage = page => {
    this.props.dispatch({
      type: `task-list/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `task-list/table-data/change-page-size`,
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
          type: 'task-list/table-data/set/selectedRows',
          payload: {
            selectedRowKeys,
            selectedRows
          }
        });
      }
    };
  };


  render() {
    const { tableData } = this.props;
    const { loading, pagination } = tableData;
    return (
      <div>
        <Table
          rowKey={'id'}
          columns={this.getColumns()}
          dataSource={tableData.list}
          loading={loading}
          pagination={false}
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
          />
        </div>
      </div>
    );
  }
}

export default MyTable;
