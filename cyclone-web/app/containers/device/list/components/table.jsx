import React from 'react';
import {
  Table,
  Button,
  Pagination,
  Dropdown,
  Menu,
  Icon,
  notification
} from 'antd';
const ButtonGroup = Button.Group;
import { hashHistory } from 'react-router';
import actions from '../actions';
import { getColumns, plainOptions } from "containers/device/common/colums";
import { getPermissonBtn } from 'common/utils';
import { renderDisplayMore } from 'common/utils';

class MyTable extends React.Component {
  state = {
    checkedList: [],
    indeterminate: true,
    checkAll: false
  };

  reload = () => {
    this.props.dispatch({
      type: 'device-list/table-data/reload'
    });
    this.props.dispatch({
      type: 'device-list/table-data/set/selectedRows',
      payload: {
        selectedRows: [],
        selectedRowKeys: []
      }
    });
  };
  批量操作入口
  batchExecAction = (name) => {
    let flag = false;
    const { tableData } = this.props;
    const selectedRows = tableData.selectedRows || [];
    if (name == 'reAccess') {
      this.execAction(name, selectedRows);
      return;
    }
    if (selectedRows.length < 1) {
      return notification.error({ message: '请至少选择一条数据' });
    }
    if (name === 'networkBoot') {
      selectedRows.forEach(tr => {
        if (tr.operation_status === 'on_shelve') {
          notification.error({ message: `物理机【${tr.sn}】已上架, 不支持从网卡启动` });
          flag = true;
        }
      });
      if (flag) {
        return;
      }
    }
    this.execAction(name, selectedRows);
  };

  //获取带外状态
  changePowerStatus = (record) => {
    this.props.dispatch({
      type: 'device-list/table-data/power-status/change',
      payload: record
    });
  };


  //操作入口
  execAction = (name, records) => {
    if (actions[name]) {
      console.log(name);
      actions[name]({
        records,
        initialValue: records,
        type: name,
        room: this.props.room,
        reload: () => {
          this.reload();
        }
      });
    }
  };


  osInstall = () => {
    const { selectedRows } = this.props.tableData;
    if (selectedRows.length < 1) {
      return notification.error({ message: '请至少选择一台设备' });
    }
    hashHistory.push({ pathname: '/device/entry', state: selectedRows });
  };


  //下载导入模板
  downloadImportTemplate = (key) => {
    if (key === 'device') {
      window.open('assets/files/existing_devices.xlsx');

    } else if (key === 'store_room') {
      window.open('assets/files/device_store_import.xlsx');
    }
  };


  getRowSelection = () => {
    const selectedRowKeys = this.props.tableData.selectedRowKeys;
    return {
      selectedRowKeys,
      onChange: (selectedRowKeys, selectedRows) => {
        this.props.dispatch({
          type: 'device-list/table-data/set/selectedRows',
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
      type: `device-list/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `device-list/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };

  exportTemp = () => {
    const query = this.props.tableData.query;
    let keys = Object.keys(query);
    const { tableData } = this.props;
    const selectedRowKeys = tableData.selectedRowKeys || [];
    keys = keys
      .map(key => {
        return `${key}=${query[key]}`;
      }) 
      .join('&');
    window.open(`/api/cloudboot/v1/devices/export?${keys}&id=${selectedRowKeys}`);
  };

  render() {
    const { tableData } = this.props;
    const { loading, pagination, selectedRows } = tableData;
    const menu_excel = (
      <Menu onClick={(e) => this.downloadImportTemplate(e.key)}>
        <Menu.Item disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_download')}
          key='device'
        >存量设备导入模板</Menu.Item>
        <Menu.Item disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_download')}
          key='store_room'
        >设备导入到库房模板</Menu.Item>
      </Menu>
    );
    const menu = (
      <Menu onClick={(e) => this.execAction(e.key)}>

        <Menu.Item disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_import')}
          key='_importStoreRoom'
        >设备导入到库房</Menu.Item>
        <Menu.Item disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_import')}
          key='_import'
        >存量设备导入</Menu.Item>
      </Menu>
    );
    return (
      <div>
        <div className='operate_btns'>
          <ButtonGroup>
            <Button
              onClick={() => this.batchExecAction('editStatus')}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_update_status')}
            >
           批量修改状态
            </Button>
            <Button
              onClick={() => this.batchExecAction('editUsage')}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_update_usage')}
            >
           批量修改用途
            </Button>
            <Button
              onClick={() => this.batchExecAction('editHardwareRemark')}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_update_usage')}
            >
           批量修改硬件备注
            </Button>            
          </ButtonGroup>
          <span className='pull-right'>
            <ButtonGroup style={{ marginRight: 8, marginBottom: 8 }} >
              <Dropdown overlay={menu_excel}>
                <Button style={{ borderBottomRightRadius: 0, borderTopRightRadius: 0 }}>下载导入模板<Icon type='down' /></Button>
              </Dropdown>
              <Dropdown overlay={menu}>
                <Button style={{ borderBottomLeftRadius: 0, borderTopLeftRadius: 0, marginLeft: -1 }}>导入 <Icon type='down' /></Button>
              </Dropdown>
              <Button
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_export')}
                onClick={this.exportTemp}
              >
                导出
              </Button>
            </ButtonGroup>
            <Button
              style={{ marginRight: 8 }}
              onClick={() => {
                this.reload();
              }}
              icon='reload'
            >
            </Button>
            {renderDisplayMore(this, plainOptions)}
          </span>
        </div>
        <div className='clearfix' />
        <Table
          rowKey={'id'}
          columns={getColumns(this, true)}
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
