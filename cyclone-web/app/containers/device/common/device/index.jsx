import React from 'react';
import {
  Alert,
  Table,
  Pagination,
  Button,
  Input,
  notification
} from 'antd';
import { createTableStore } from 'utils/table-reducer';
import handleAction from './sync-actions/index';
import asyncActions from './async-actions/index';
import M from 'immutable';
import { getColumns } from "containers/device/common/colums";
import { getWithArgs } from 'common/xFetch2';

export default class MyTable extends React.Component {
  componentDidMount() {
    this.dispatch({
      type: 'device/table-data/search',
      payload: this.props.query
    });
  }

  constructor(props) {
    super(props);
    this.state = {
      checkedList: [],
      data: M.fromJS({
        tableData: createTableStore()
      })
    };
  }
  getState = () => {
    return this.state.data;
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
          {
            !this.props.hideButton &&
            <div className=' panel-footer'>
              <Button type='primary' onClick={this.handleSubmit}>
                确定
              </Button>
            </div>
          }

        </div>
      </div>
    );
  }

  onSearch = (value, tableQuery, key) => {
    if (this.props.query) {
      return this.dispatch({
        type: 'device/table-data/search',
        payload: {
          ...tableQuery,
          [key]: value,
          ...this.props.query
        }
      });
    }
    this.dispatch({
      type: 'device/table-data/search',
      payload: {
        ...tableQuery,
        [key]: value
      }
    });
  };


  renderBody() {
    const tableData = this.state.data.get('tableData').toJS();
    const tableQuery = tableData.query;
    const { list, pagination, loading, selectedRowKeys, selectedRows } = tableData;
    const rowSelection = {
      type: this.props.checkType,
      selectedRowKeys,
      onChange: (selectedRowKeys, selectedRows) => {
        this.dispatch({
          type: 'device/table-data/set/selectedRows',
          payload: {
            selectedRowKeys,
            selectedRows
          }
        });
        if (this.props.form) {
          const { setFieldsValue } = this.props.form;
          setFieldsValue({ sn: selectedRows });
        }
      }
    };

    return (
      <div className='node-body'>
        <div>
          <Alert
            message={this.props.checkType === 'radio' ? `选择的设备： ${selectedRows[0] ? selectedRows[0].sn : ''}` : `已选择： ${selectedRowKeys.length} 项`}
            type='info'
            showIcon={true}
            style={{
              marginBottom: 8,
              overflow: 'hidden',
              wordBreak: 'break-all'
            }}
          />
        </div>
        <div
          className='search-form'
          style={{ textAlign: 'right', marginBottom: 8 }}
        >
          <span style={{ marginRight: 8 }}>
            <Input.Search
              placeholder='序列号'
              style={{ width: 200 }}
              onChange={e => {
                this.onSearch(e.target.value, tableQuery, 'sn');
              }}
            />
          </span>
          <span style={{ marginRight: 8 }}>
            <Input.Search
              placeholder='固资编号'
              style={{ width: 200 }}
              onChange={e => {
                this.onSearch(e.target.value, tableQuery, 'fixed_asset_number');
              }}
            />
          </span>
          <span>
            <Input.Search
              placeholder='内网IP'
              style={{ width: 200 }}
              onChange={e => {
                this.onSearch(e.target.value, tableQuery, 'intranet_ip');
              }}
            />
          </span>
        </div>

        <Table
          rowKey={'id'}
          rowSelection={rowSelection}
          columns={getColumns(this, [])} //兼容物理设备列表
          pagination={false}
          dataSource={list}
          loading={loading}
        />
        <Pagination
          showTotal={(total) => `共 ${total} 条`}
          showSizeChanger={true}
          current={pagination.page}
          pageSize={pagination.pageSize}
          total={pagination.total}
          onShowSizeChange={(page, pageSize) => {
            this.dispatch({
              type: 'device/table-data/change-page-size',
              payload: {
                page,
                pageSize
              }
            });
          }}
          onChange={page => {
            this.dispatch({
              type: 'device/table-data/change-page',
              payload: {
                page
              }
            });
          }}
        />
      </div>
    );
  }


  handleSubmit = async (ev) => {
    const tableData = this.state.data.get('tableData').toJS();
    const { selectedRows } = tableData;
    if (selectedRows.length < 1) {
      return notification.error({ message: '请选择设备' });
    }
    let limit = this.props.limit || Infinity;
    if (selectedRows.length > limit) {
      return notification.error({ message: `选择的设备数不能超过${limit}` });
    }
    // if (this.props.getServerRoom) {
    //   try {
    //     await Promise.all(selectedRows.map((it, i) => new Promise((resolve, reject) => {
    //       getWithArgs('/api/cloudboot/v1/server-rooms', { page: 1, page_size: 100, idc_id: it.idc.id }).then(res => {
    //         if (res.status !== 'success') {
    //           notification.error({ message: res.message });
    //           reject();
    //         }
    //
    //         it.rooms = res.content.records || [];
    //         resolve();
    //       });
    //     })));
    //   } catch (err) {
    //     notification(err);
    //   }
    // }
    this.props.handleSubmit(tableData);
  };
}
