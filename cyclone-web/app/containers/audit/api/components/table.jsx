import React from 'react';
import {
  Table,
  Tooltip,
  Pagination,
  Badge
} from 'antd';
import actions from '../actions';
import TableControlCell from 'components/TableControlCell';
import { API_STATUS_COLOR } from 'common/enums';
import moment from 'moment';
import { TIME_FORMAT } from 'common/enums';

class MyTable extends React.Component {

  reload = () => {
    this.props.dispatch({
      type: 'audit-api/table-data/reload'
    });
    this.props.dispatch({
      type: 'audit-api/table-data/set/selectedRows',
      payload: {
        selectedRows: [],
        selectedRowKeys: []
      }
    });
  };


  //操作入口
  execAction = (name, records) => {
    if (actions[name]) {
      actions[name]({
        records,
        initialValue: records,
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
          type: 'audit-api/table-data/set/selectedRows',
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
      type: `audit-api/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `audit-api/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };


  getColumns = () => {
    return [
      {
        title: '接口地址',
        dataIndex: 'api',
        render: (text) => <Tooltip placement="top" title={text}>{text}</Tooltip>
      },
      {
        title: '接口描述',
        dataIndex: 'description'
      },
      {
        title: '请求方法',
        dataIndex: 'method'
      },
      {
        title: '操作者',
        dataIndex: 'operator'
      },
      {
        title: '操作时间',
        dataIndex: 'created_at',
        render: (t) => <span>{moment(t).format(TIME_FORMAT)}</span>
      },
      {
        title: '耗时（s）',
        dataIndex: 'time'
      },
      {
        title: '执行状态',
        dataIndex: 'status',
        render: type => {
          const color = API_STATUS_COLOR[type] ? API_STATUS_COLOR[type][0] : 'transparent';
          const word = API_STATUS_COLOR[type] ? API_STATUS_COLOR[type][1] : '';
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
        width: 70,
        render: (text, record) => {
          const commands = [
            {
              name: '详情',
              command: '_detail'
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
