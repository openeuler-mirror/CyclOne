import React from 'react';
import {
  Table,
  Pagination,
  Badge
} from 'antd';
import actions from '../actions';
import { APPROVAL_TYPE, APPROVAL_STATUS_COLOR } from 'common/enums';

class MyTable extends React.Component {

  //操作入口
  execAction = (name, record) => {
    if (actions[name]) {
      actions[name]({
        record,
        approval_id: record.id,
        userInfo: this.props.userInfo,
        reload: () => {
          this.props.reload();
        }
      });
    }
  };

  changePage = page => {
    this.props.dispatch({
      type: `approval/${this.props.type}-table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `approval/${this.props.type}-table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };


  getColumns = () => {
    return [
      {
        title: '审批标题',
        dataIndex: 'title',
        render: (text, record) => <a onClick={() => this.execAction('_detail', record)}>{text}</a>
      },
      {
        title: '审批类型',
        dataIndex: 'type',
        render: (t) => APPROVAL_TYPE[t]
      },
      {
        title: '发起时间',
        dataIndex: 'start_time'
      },
      {
        title: '完成时间',
        dataIndex: 'end_time'
      },
      {
        title: '状态',
        dataIndex: 'status',
        render: (type, record) => {
          const color = APPROVAL_STATUS_COLOR[type] ? APPROVAL_STATUS_COLOR[type][0] : 'transparent';
          const word = APPROVAL_STATUS_COLOR[type] ? APPROVAL_STATUS_COLOR[type][1] : '';
          return (
            <div>
              <Badge
                dot={true}
                style={{
                  background: color
                }}
              />{' '}
              &nbsp;&nbsp; {word}
              {/*{record.is_rejected === 'yes' ? <span style={{ color: '#ff3700' }}>（审批被拒）</span> : ''}*/}
            </div>
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
