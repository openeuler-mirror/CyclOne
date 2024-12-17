import React from 'react';
import {
  Table,
  Tooltip,
  Pagination,
  Button,
  notification
} from 'antd';
import actions from '../actions';
import TableControlCell from 'components/TableControlCell';
import { getPermissonBtn } from 'common/utils';
import { BUILTIN } from 'common/enums';

class MyTable extends React.Component {

  reload = () => {
    this.props.dispatch({
      type: 'order-deviceCategory/table-data/reload'
    });
    this.props.dispatch({
      type: 'order-deviceCategory/table-data/set/selectedRows',
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
          type: 'order-deviceCategory/table-data/set/selectedRows',
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
      type: `order-deviceCategory/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `order-deviceCategory/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };

  getColumns = () => {
    return [
      {
        title: '设备类型',
        dataIndex: 'category'
      },
      {
        title: '硬件配置',
        dataIndex: 'hardware',
        render: (text) => <Tooltip placement="top" title={text}>{text}</Tooltip>
      },
      {
        title: '处理器生产商',
        dataIndex: 'central_processor_manufacture'
      },
      {
        title: '处理器架构',
        dataIndex: 'central_processor_arch'
      },
      {
        title: '功率/W',
        dataIndex: 'power'
      },            
      {
        title: '设备高度/U',
        dataIndex: 'unit'
      },
      {
        title: '是否金融信创生态产品',
        dataIndex: 'is_fiti_eco_product',
        render: (text) => <span className={`yes_no_status ${text === 'yes' ? 'yes_status' : 'no_status'}`}>{BUILTIN[text]}</span>
      },      
      {
        title: '备注',
        dataIndex: 'remark'
      },
      {
        title: '操作',
        dataIndex: 'operate',
        render: (text, record) => {
          const commands = [
            {
              name: '修改',
              command: '_update',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_device_category_update')
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
            onClick={() => this.execAction('_create', {})}
            type='primary'
            icon='plus'
            style={{ marginRight: 8 }}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_device_category_create')}
          >
            新增
          </Button>
          <Button
            type='danger'
            onClick={() => this.batchExecAction('_delete')}
            style={{ marginRight: 8 }}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_device_category_delete')}
          >
            删除
          </Button>
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
