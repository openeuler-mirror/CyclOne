import React from 'react';
import {
  Select,
  Table,
  Pagination,
  Button,
  notification
} from 'antd';
import { createTableStore } from 'utils/table-reducer';
import handleAction from './sync-actions/index';
import asyncActions from './async-actions/index';
import M from 'immutable';
const Option = Select.Option;

export default class OperationTarget extends React.Component {
  componentDidMount() {
    this.dispatch({
      type: 'hardware/table-data/get'
    });
    const { record } = this.props;
    if (record.hardware_tpl_id && record.hardware_tpl_id !== 0) {
      this.dispatch({
        type: 'hardware/table-data/set/selectedRows',
        payload: {
          selectedRowKeys: [record.hardware_tpl_id],
          selectedRows: [{ id: record.hardware_tpl_id, name: record.hardware_tpl_name }]
        }
      });
    }
  }

  constructor(props) {
    super(props);
    this.state = {
      data: M.fromJS({
        tableData: createTableStore()
      })
    };
  }
  getState = () => {
    return this.state.data;
  };
  getColumns = () => {
    return [
      {
        title: '名称',
        dataIndex: 'name',
        width: 200
      }
      // {
      //   title: '厂商',
      //   dataIndex: 'vendor',
      //   width: 200
      // },
      // {
      //   title: '产品型号',
      //   dataIndex: 'model_name',
      //   width: 200
      // },
      // {
      //   title: '配置内容',
      //   dataIndex: 'data_config',
      //   width: 200,
      //   render: (text, record) => {
      //     const data = record.data || [];
      //     return <span>{[...new Set(data.map(d => d.category))].join(' | ')}</span>;
      //   }
      // }
    ];
  };


  dispatch = action => {
    if (asyncActions[action.type]) {
      asyncActions[action.type](
        this.state.data,
        action,
        this.dispatch,
        this.getState
      );
      return;
    }

    const data = handleAction(this.state.data, action);
    this.setState({
      data
    });
  };
  render() {
    return (
      <div className='host'>
        <div className='panel'>
          <div className='panel-body'>{this.renderBody()}</div>
          { !this.props.bunchEdit && <div className=' panel-footer'>
            <Button type='primary' onClick={this.handleSubmit}>
              确定
            </Button>
          </div>}
        </div>
      </div>
    );
  }


  renderBody() {
    const tableData = this.state.data.get('tableData').toJS();
    const { list, pagination, loading, selectedRowKeys, selectedRows } = tableData;
    const rowSelection = {
      type: 'radio',
      selectedRowKeys,
      onChange: (selectedRowKeys, selectedRows) => {
        this.dispatch({
          type: 'hardware/table-data/set/selectedRows',
          payload: {
            selectedRowKeys,
            selectedRows
          }
        });
        if (this.props.bunchEdit) {
          this.props.dispatch({
            type: 'bunchEdit/hardware/data',
            payload: selectedRows[0]
          });
        }
      }
    };

    return (
      <div className='node-body'>
        <Table
          rowKey={'id'}
          rowSelection={rowSelection}
          scroll={{ y: 200 }}
          columns={this.getColumns()}
          pagination={false}
          dataSource={list}
          loading={loading}
        />
        <Pagination
          showSizeChanger={true}
          current={pagination.page}
          pageSize={pagination.pageSize}
          total={pagination.total}
          showTotal={(total) => `共 ${total} 条`}
          onShowSizeChange={(page, pageSize) => {
            this.dispatch({
              type: 'hardware/table-data/change-page-size',
              payload: {
                page,
                pageSize
              }
            });
          }}
          onChange={page => {
            this.dispatch({
              type: 'hardware/table-data/change-page',
              payload: {
                page
              }
            });
          }}
        />
      </div>
    );
  }

  handleSubmit = ev => {
    const tableData = this.state.data.get('tableData').toJS();
    const { selectedRows } = tableData;
    //单次编辑时选择硬件配置
    if (!this.props.bunchEdit && selectedRows.length < 1) {
      return notification.error({ message: '请选择RAID类型' });
    }
    if (this.props.bunchEdit && selectedRows.length > 0) {
      this.props.dispatch({
        type: 'bunchEdit/hardware/data',
        payload: selectedRows[0]
      });
    }
    this.props.handleSubmit(selectedRows[0]);
  };
}
