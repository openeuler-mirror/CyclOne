import React from 'react';
import {
  Table,
  Pagination
} from 'antd';
import actions from '../actions';
import TableControlCell from 'components/TableControlCell';
import moment from 'moment';
import { TIME_FORMAT } from 'common/enums';


class MyTable extends React.Component {

  reload = () => {
    this.props.dispatch({
      type: 'audit-log/table-data/reload'
    });
    this.props.dispatch({
      type: 'audit-log/table-data/set/selectedRows',
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
          type: 'audit-log/table-data/set/selectedRows',
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
      type: `audit-log/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `audit-log/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };


  getColumns = () => {
    return [
      {
        title: '操作时间',
        dataIndex: 'created_at',
        render: (t) => <span>{moment(t).format(TIME_FORMAT)}</span>
      },
      {
        title: '操作人',
        dataIndex: 'operator'
      },
      {
        title: '操作类型',
        dataIndex: 'category_name'
      },
      {
        title: '路由',
        dataIndex: 'url'
      },
      {
        title: '请求方式',
        dataIndex: 'http_method'
      },
      {
        title: '操作',
        dataIndex: 'operate',
        width: 70,
        render: (text, record) => {
          const commands = [
            {
              name: '详情',
              command: '_compare'
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
