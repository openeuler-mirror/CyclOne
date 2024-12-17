import React from 'react';
import {
  Alert,
  Table,
  Pagination,
  Button,
  Input,
  Select,
  notification
} from 'antd';
import { createTableStore } from 'utils/table-reducer';
import handleAction from './sync-actions/index';
import asyncActions from './async-actions/index';
import M from 'immutable';
const { Option } = Select;


export default class MyTable extends React.Component {
  componentDidMount() {
    this.dispatch({
      type: 'approval-page/table-data/search',
      payload: this.props.tableQuery
    });
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

  dispatch = action => {
    if (asyncActions[action.type]) {
      asyncActions[action.type](
        this.state.data,
        action,
        this.dispatch,
        this.getState,
        this.props.tableDataUrl
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
        type: 'approval-page/table-data/search',
        payload: {
          ...tableQuery,
          [key]: value,
          ...this.props.query
        }
      });
    }
    this.dispatch({
      type: 'approval-page/table-data/search',
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
          type: 'approval-page/table-data/set/selectedRows',
          payload: {
            selectedRowKeys,
            selectedRows
          }
        });
        if (this.props.form) {
          const { setFieldsValue } = this.props.form;
          setFieldsValue({ ids: selectedRowKeys });
        }
      }
    };

    return (
      <div className='node-body'>
        <div>
          <Alert
            message={this.props.checkType === 'radio' ? `选择的${this.props.category}： ${selectedRows[0] ? selectedRows[0][this.props.searchKey] : ''}` : `已选择： ${selectedRowKeys.length} 项`}
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
          {
            this.props.searchKey.map((item, index) => {
              if (item.type === 'select') {
                return <span style={{ marginRight: index === this.props.searchKey.length - 1 ? 0 : 8 }}>
                  <Select
                    mode='multiple'
                    style={{ width: 200 }}
                    placeholder={item.label}
                    onChange={value => {
                      this.onSearch(value, tableQuery, item.key);
                    }}
                  >
                    {
                      item.async ? (this.props[item.async].data || []).map(it => <Option value={it.value} key={it.value}>{it.label}</Option>) :
                     item.children.map(it => <Option value={it.value} key={it.value}>{it.label}</Option>
                     )
                   }
                  </Select>
                </span>;
              }
              return <span style={{ marginRight: index === this.props.searchKey.length - 1 ? 0 : 8 }}>
                <Input.Search
                  placeholder={item.label}
                  style={{ width: 200 }}
                  onChange={e => {
                    this.onSearch(e.target.value, tableQuery, item.key);
                  }}
                />
              </span>;
            })
          }
          {/*<span>*/}
          {/*<Input.Search*/}
          {/*placeholder={`搜索${this.props.category}...`}*/}
          {/*style={{ width: 200 }}*/}
          {/*onChange={ev => {*/}
          {/*if (this.props.query) {*/}
          {/*return this.dispatch({*/}
          {/*type: 'approval-page/table-data/search',*/}
          {/*payload: {*/}
          {/*[this.props.searchKey]: ev.target.value,*/}
          {/*...this.props.query*/}
          {/*}*/}
          {/*});*/}
          {/*}*/}
          {/*this.dispatch({*/}
          {/*type: 'approval-page/table-data/search',*/}
          {/*payload: {*/}
          {/*[this.props.searchKey]: ev.target.value*/}
          {/*}*/}
          {/*});*/}
          {/*}}*/}
          {/*/>*/}
          {/*</span>*/}
        </div>

        <Table
          rowKey={'id'}
          rowSelection={rowSelection}
          columns={this.props.tableColumn()}
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
              type: 'approval-page/table-data/change-page-size',
              payload: {
                page,
                pageSize
              }
            });
          }}
          onChange={page => {
            this.dispatch({
              type: 'approval-page/table-data/change-page',
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
    if (selectedRows.length < 1) {
      return notification.error({ message: `请选择${this.props.category}` });
    }
    let limit = this.props.limit || Infinity;
    if (selectedRows.length > limit) {
      return notification.error({ message: `选择的${this.props.category}数不能超过${limit}` });
    }
    this.props.handleSubmit(tableData);
  };
}
