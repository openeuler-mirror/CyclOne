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
import { NETWORK_DEVICE_TYPE, ORDER_STATUS } from 'common/enums';
import { getPermissonBtn } from 'common/utils';
import { post } from 'common/xFetch2';

class MyTable extends React.Component {

  reload = () => {
    this.props.dispatch({
      type: 'order-list/table-data/reload'
    });
    this.props.dispatch({
      type: 'order-list/table-data/set/selectedRows',
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
    const selectedRowKeys = tableData.selectedRowKeys || [];
    if (name === '_export') {
      const query = this.props.tableData.query;
      let keys = Object.keys(query);
      keys = keys
        .map(key => {
          return `${key}=${query[key]}`;
        }) 
        .join('&');
      window.open(`/api/cloudboot/v1/orders/export?${keys}&id=${selectedRowKeys}`);
      return;
    }
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
        physicalArea: this.props.physicalArea,
        deviceCategory: this.props.deviceCategory,
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
          type: 'order-list/table-data/set/selectedRows',
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
      type: `order-list/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `order-list/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };

  getColumns = () => {
    return [
      {
        title: '订单号',
        dataIndex: 'number',
        //render: (text, record) => <a onClick={() => this.execAction('_detail', { id: record.id })}>{text}</a>
        render: (text, record) => <Tooltip placement="top" title={text}>
          <a onClick={() => this.execAction('_detail', { id: record.id })}>{text}</a>
        </Tooltip>
      },
      {
        title: '日期',
        dataIndex: 'created_at'
      },
      {
        title: '数据中心',
        dataIndex: 'idc',
        render: (text, record) => <a onClick={() => this.execAction('idc_detail', { id: text.id })}>{text.name}</a>
      },
      {
        title: '机房管理单元',
        dataIndex: 'server_room',
        //render: (text, record) => <a onClick={() => this.execAction('room_detail', { id: text.id })}>{text.name}</a>
        render: (text, record) => <Tooltip placement="top" title={text.name}>
          <a onClick={() => this.execAction('room_detail', { id: text.id })}>
            {text.name}
          </a>
        </Tooltip>
      },
      {
        title: '物理区域',
        dataIndex: 'physical_area',
        render: (text) => <Tooltip placement="top" title={text}>{text}</Tooltip>
      },
      {
        title: '设备类型',
        dataIndex: 'category'
      },
      {
        title: '数量',
        dataIndex: 'amount'
      },
      {
        title: '剩余到货数',
        dataIndex: 'left_amount'
      },
      {
        title: '预计到货时间',
        dataIndex: 'expected_arrival_date'
      },
      {
        title: '备注',
        dataIndex: 'remark'
      },
      {
        title: '状态',
        dataIndex: 'status',
        render: (T) => ORDER_STATUS[T]
      },
      {
        title: '操作',
        dataIndex: 'operate',
        render: (text, record) => {
          const commands = [
            // {
            //   name: '修改',
            //   command: '_update',
            //   disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_network_device_delete')
            // },
            {
              name: '取消',
              command: '_cancel',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_order_cancel') || (record.status === 'canceled' || record.status === 'finished')
            },
            {
              name: '确认',
              command: '_confirm',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_order_confirm') || (record.status === 'canceled' || record.status === 'finished')
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
    const { loading, pagination, selectedRowKeys } = tableData;
    return (
      <div>
        <div className='operate_btns'>
          <Button
            onClick={() => this.execAction('_create')}
            type='primary'
            icon='plus'
            style={{ marginRight: 8 }}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_order_create')}
          >
            新增
          </Button>
          <Button
            type='danger'
            onClick={() => this.batchExecAction('_delete')}
            style={{ marginRight: 8 }}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_order_delete')}
          >
            删除
          </Button>
          <Button
            onClick={() => this.batchExecAction('_export')}
            style={{ marginRight: 8 }}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_order_export')}
          >
            导出
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
