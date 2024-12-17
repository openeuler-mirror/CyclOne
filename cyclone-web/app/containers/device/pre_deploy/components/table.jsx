import React from 'react';
import {
  Table,
  Button,
  Pagination,
  notification,
} from 'antd';
const ButtonGroup = Button.Group;
import { hashHistory } from 'react-router';
import actions from '../../list/actions';
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
      type: 'device-pre_deploy/table-data/reload'
    });
    this.props.dispatch({
      type: 'device-pre_deploy/table-data/set/selectedRows',
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

  //获取带外状态
  changePowerStatus = (record) => {
    this.props.dispatch({
      type: 'device-pre_deploy/table-data/power-status/change',
      payload: record
    });
  };


  //操作入口
  execAction = (name, records) => {
    if (actions[name]) {
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

  oobAction = (name, ifSelected) => {
    let flag = false;
    const { tableData } = this.props;
    const records = tableData.selectedRows;
    if (!ifSelected && records.length < 1) {
      return notification.error({ message: '请至少选择一台设备' });
    }
    if (actions[name]) {
      actions[name]({
        records,
        userInfo: this.props.userInfo,
        type: name,
        reload: () => {
          this.reload();
          this.props.dispatch({
            type: 'device-setting/table-data/set/selectedRows',
            payload: {
              selectedRows: [],
              selectedRowKeys: []
            }
          });
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
  downloadImportTemplate = () => {
    window.open('assets/files/devices.xlsx');
  };


  getRowSelection = () => {
    const selectedRowKeys = this.props.tableData.selectedRowKeys;
    return {
      selectedRowKeys,
      onChange: (selectedRowKeys, selectedRows) => {
        this.props.dispatch({
          type: 'device-pre_deploy/table-data/set/selectedRows',
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
      type: `device-pre_deploy/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `device-pre_deploy/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };

  render() {
    const { tableData } = this.props;
    const { loading, pagination, selectedRows } = tableData;

    return (
      <div>
        <div className='operate_btns'>
          <Button
            onClick={() => this.osInstall()}
            type='primary'
            style={{ marginRight: 8 }}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_predeploy_physical_machine_osInstall')}
          >
            申请上架部署
          </Button>
          <span>
            已选 { selectedRows.length } 项
          </span>
          <span className='pull-right'>
            <ButtonGroup style={{ marginRight: 8 }}>
              <Button
                onClick={() => this.downloadImportTemplate()}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_predeploy_physical_machine_download')}
              >
                下载导入模板
              </Button>
              <Button
                onClick={() => this.execAction('pre_import')}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_predeploy_physical_machine_import')}
              >
                导入
              </Button>
            </ButtonGroup>
            <Button
            onClick={() => this.oobAction('reAccess')}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_oob_re_access')}
            style={{ marginRight: 8 }}
            >
              纳管带外
            </Button>      
            <Button
              onClick={() => {
                this.reload();
              }}
              icon='reload'
              style={{ marginRight: 8 }}
            >
            </Button>

            {renderDisplayMore(this, plainOptions)}

            <Button
            onClick={() => this.batchExecAction('_delete')}
            type='danger'
            style={{ marginLeft: 8 }}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_predeploy_physical_machine_delete')}
          >
            删除
          </Button>

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
