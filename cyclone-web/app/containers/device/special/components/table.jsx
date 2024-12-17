import React from 'react';
import {
  Table,
  Button,
  Badge,
  Pagination,
  notification,
  Tooltip
} from 'antd';
const ButtonGroup = Button.Group;
import { Link } from 'react-router';
import actions from '../actions';
import { plainOptions } from "containers/device/common/colums";
import { getPermissonBtn } from 'common/utils';
import { renderDisplayMore } from 'common/utils';
import { OPERATION_STATUS_COLOR } from "common/enums";


class MyTable extends React.Component {
  state = {
    checkedList: [],
    indeterminate: true,
    checkAll: false
  };

  reload = () => {
    this.props.dispatch({
      type: 'device-special/table-data/reload'
    });
    this.props.dispatch({
      type: 'device-special/table-data/set/selectedRows',
      payload: {
        selectedRows: [],
        selectedRowKeys: []
      }
    });
  };

  //获取带外状态
  changePowerStatus = (record) => {
    this.props.dispatch({
      type: 'device-special/table-data/power-status/change',
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

  //批量操作入口
  batchExecAction = (name) => {
    const { tableData } = this.props;
    const selectedRows = tableData.selectedRows || [];
    if (selectedRows.length < 1) {
      return notification.error({ message: '请至少选择一条数据' });
    }
    this.execAction(name, selectedRows);
  };


  //下载导入模板
  downloadImportTemplate = () => {
    window.open('assets/files/special_device_import.xlsx');
  };


  getRowSelection = () => {
    const selectedRowKeys = this.props.tableData.selectedRowKeys;
    return {
      selectedRowKeys,
      onChange: (selectedRowKeys, selectedRows) => {
        this.props.dispatch({
          type: 'device-special/table-data/set/selectedRows',
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
      type: `device-special/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `device-special/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };
  getColumns = () => {
    return [
      {
        title: '固资编号',
        dataIndex: 'fixed_asset_number',
        render: (text) => <Tooltip placement="top" title={text}>{text}</Tooltip>
      },      
      {
        title: '序列号SN',
        dataIndex: 'sn',
        render: (text, record) => <Link to={`/device/detail/${text}`}>{text}</Link>
      },
      {
        title: '厂商',
        dataIndex: 'vendor'
      },
      {
        title: '型号',
        dataIndex: 'model'
      },
      {
        title: '用途',
        dataIndex: 'usage'
      },      
      {
        title: '设备类型',
        dataIndex: 'category'
      },      
      {
        title: '机房管理单元名称',
        dataIndex: 'server_room_name',
        render: (t, record) => record.server_room ? record.server_room.name : ''
      },
      {
        title: '机架编号',
        dataIndex: 'server_cabinet_number',
        render: (t, record) => record.server_cabinet ? record.server_cabinet.number : ''
      },
      {
        title: '机位编号',
        dataIndex: 'server_usite_number',
        render: (t, record) => record.server_usite ? record.server_usite.number : ''
      },
      {
        title: '硬件说明',
        dataIndex: 'hardware_remark'
      },
      {
        title: '内网IP',
        dataIndex: 'intranet_ip'
      },
      {
        title: '外网IP',
        dataIndex: 'extranet_ip'
      },
      {
        title: '运营状态',
        dataIndex: 'operation_status',
        render: type => {
          const color = OPERATION_STATUS_COLOR[type] ? OPERATION_STATUS_COLOR[type][0] : 'transparent';
          const word = OPERATION_STATUS_COLOR[type] ? OPERATION_STATUS_COLOR[type][1] : '';
          return (
            <div>
              <Badge
                dot={true}
                style={{
                  background: color
                }}
              />{' '}
              &nbsp;&nbsp; {word}
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
        <div className='operate_btns'>
          <Button
            type='primary'
            style={{ marginRight: 8 }}
            onClick={() => this.execAction('create')}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_special_device_create')}
          >
            新增
            </Button>
          <Button
            type='danger'
            style={{ marginRight: 8 }}
            onClick={() => this.batchExecAction('_delete')}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_special_device_delete')}
          >
            删除
            </Button>
          <span className='pull-right'>
            <ButtonGroup style={{ marginRight: 8 }}>
              <Button
                onClick={() => this.downloadImportTemplate()}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_special_device_import_download')}
              >
                下载导入模板
              </Button>
              <Button
                onClick={() => this.execAction('_import')}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_special_device_import')}
              >
                导入
              </Button>
            </ButtonGroup>
          </span>
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
