import React from 'react';
import {
  Table,
  Button,
  notification,
  Tabs,
  Tooltip
} from 'antd';
import { createTableStore } from 'utils/table-reducer';
import handleAction from './sync-actions/index';
import asyncActions from './async-actions/index';
import M from 'immutable';
import { BOOT_MODE, INSTALL_TYPE, OS_LIFECYCLE } from "common/enums";
import Layout from 'components/layout/page-layout-tabs';

const TabPane = Tabs.TabPane;

function listFilter(target, key) {
  const result = []
  target.map(item =>{
      if (item["install_type"] == key) {
        result.push(item)
      };
  });
  return result
};

export default class OperationTarget extends React.Component {
  componentDidMount() {
    this.dispatch({
      type: 'system/table-data/get'
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
  getColumns = () => {
    return [
      {
        title: '名称',
        dataIndex: 'name',
        width: 300
      },
      {
        title: '操作系统',
        dataIndex: 'family',
        width: 100
      },
      {
        title: '安装方式',
        dataIndex: 'install_type',
        width: 100,
        render: (text) => <span>{INSTALL_TYPE[text]}</span>
      },
      {
        title: '启动方式',
        dataIndex: 'boot_mode',
        width: 100,
        render: (text) => <span>{BOOT_MODE[text]}</span>
      },
      {
        title: '架构',
        dataIndex: 'arch',
        width: 100
      }, 
      {
        title: '生命周期',
        dataIndex: 'os_lifecycle',
        width: 100,
        render: (text) => <Tooltip placement="top"  title='默认仅提供Active版本'><span>{OS_LIFECYCLE[text]}</span></Tooltip>
      },
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

  // 改变标签时通过传参过滤得到对应安装方式的数据
  onTabChange = (key) => {
    const installType = key
    this.setState({installType})
  }

  render() {
    return (
      <div className='host'>
        <div className='panel'>
          <div className='panel-body'>
            <Layout title='安装方式'>
              <Tabs onChange={this.onTabChange} >
                <TabPane tab={<span>镜像</span>} key='image'>{this.renderBody()}</TabPane>
                <TabPane tab={<span>PXE</span>} key='pxe'>{this.renderBody()}</TabPane>
              </Tabs>
            </Layout>
          </div>
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
    const { list, loading, selectedRowKeys, selectedRows } = tableData;
    var filteredList = []
    
    // 默认展示镜像安装方式的数据
    if (this.state.installType) {
      filteredList = listFilter(list,this.state.installType)  
    } else {
      filteredList = listFilter(list,'image')
    }
    
    const rowSelection = {
      type: 'radio',
      selectedRowKeys,
      onChange: (selectedRowKeys, selectedRows) => {
        this.dispatch({
          type: 'system/table-data/set/selectedRows',
          payload: {
            selectedRowKeys,
            selectedRows
          }
        });
        if (this.props.bunchEdit) {
          this.props.dispatch({
            type: 'bunchEdit/system/data',
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
          dataSource={filteredList}
          loading={loading}
        />
      </div>
    );
  }

  handleSubmit = () => {
    const tableData = this.state.data.get('tableData').toJS();
    const { selectedRows } = tableData;
    if (selectedRows.length < 1) {
      return notification.error({ message: '请选择操作系统' });
    }
    this.props.handleSubmit(selectedRows[0]);
  };
}
