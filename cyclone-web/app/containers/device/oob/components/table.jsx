import React from 'react';
import {
  Table,
  Pagination,
  Button,
  Icon,
  notification,
  Badge,
  Divider,
  Tooltip,
} from 'antd';
import actions from '../actions';
import { getPermissonBtn } from 'common/utils';
import { Link } from 'react-router';
import copy from 'copy-to-clipboard';
import { OOB_ACCESSIBLE, OOB_STATUS_COLOR } from "common/enums";
import { get } from 'https';


class MyTable extends React.Component {
  state = {
    checkedList: [],
    indeterminate: true,
    checkAll: false
  };

  reload = () => {
    this.props.dispatch({
      type: 'device-oob/table-data/reload'
    });
    this.props.dispatch({
      type: 'device-oob/table-data/set/selectedRows',
      payload: {
        selectedRows: [],
        selectedRowKeys: []
      }
    });
  };
  getOobUser = (record) => {
    get(`/api/cloudboot/v1/devices/${record.sn}/oob-user`).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message || '获取失败' });
      }
      notification.success({ message: '操作成功' });
    });
  }

  exportOob = () => {
    const query = this.props.tableData.query;
    let keys = Object.keys(query);
    const { tableData } = this.props;
    const selectedRowKeys = tableData.selectedRowKeys || [];
    keys = keys
      .map(key => {
        return `${key}=${query[key]}`;
      }) 
      .join('&');
    window.open(`/api/cloudboot/v1/devices/oob/export?${keys}&id=${selectedRowKeys}`);
  };

    //批量操作入口
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
        if (tr.operation_status !== 'pre_deploy') {
          notification.error({ message: `物理机【${tr.sn}】运营状态非【待部署】, 不支持从网卡启动` });
          flag = true;
        }
      });
      if (flag) {
        return;
      }
    }
    this.execAction(name, selectedRows);
  };
    

  //操作入口
  execAction = (name, records) => {
    if (actions[name]) {
      actions[name]({
        records,
        type: name,
        initialValue: records,
        reload: () => {
          this.reload();
        }
      });
    }
  };

  changePage = page => {
    this.props.dispatch({
      type: `device-oob/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `device-oob/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };

  copyPsd = (t) => {
    if (copy(t)) {
      notification.success({ message: '复制成功' });
    } else {
      notification.error({ message: '复制失败' });
    }
  }

  getColumns = () => {

    const columns = [
      {
        title: '固资编号',
        dataIndex: 'fixed_asset_number'
      },
      {
        title: '序列号',
        dataIndex: 'sn',
        render: (text) => {
          return <Link to={`/device/detail/${text}`}>{text}</Link>;
        }
      },
      {
        title: '内网 IP',
        dataIndex: 'intranet_ip'
      },
      {
        title: '带外IP',
        dataIndex: 'oob_ip',
        render: (t) => {
          return <a href={`http://${t}`} target='_blank'>{t}</a>;
        }
      },
      {
        title: '带外用户名',
        dataIndex: 'oob_user'
      },
      {
        title: '带外密码',
        dataIndex: 'oob_password',
        render: (t) => {
          if (t) {
            return <span>{t} &nbsp;<Icon onClick={() => this.copyPsd(t)} style={{ float: 'right' }} type='copy'/></span>;

          }
        }
      },
      {
        title: '带外状态',
        dataIndex: 'oob_accessible',
        render: type => {
          const color = OOB_STATUS_COLOR[type] ? OOB_STATUS_COLOR[type][0] : 'transparent';
          const word = OOB_STATUS_COLOR[type] ? OOB_STATUS_COLOR[type][1] : '';
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
      },
      {
        title: '操作',
        dataIndex: 'operate',
        render: (text, record) => {
          return (
            <div>
              <a disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_update_oob')} href='javascript:;' onClick={() => this.execAction('update', record)}>修改</a>
              <Divider type='vertical'/>
              <a href='javascript:;' onClick={() => this.getOobUser(record)}>刷新</a>
            </div>
          );
        }
      }
    ];
    return columns;
  };

  getRowSelection = () => {
    const selectedRowKeys = this.props.tableData.selectedRowKeys;
    return {
      selectedRowKeys,
      onChange: (selectedRowKeys, selectedRows) => {
        this.props.dispatch({
          type: 'device-oob/table-data/set/selectedRows',
          payload: {
            selectedRowKeys,
            selectedRows
          }
        });
      }
    };
  };

  render() {
    const { tableData } = this.props;
    const { loading, pagination, selectedRows } = tableData;
    return (
      <div>
        <div className='operate_btns'>
          <Button.Group style={{ marginRight: 8 }}>
            <Button
              type='primary'
              onClick={() => this.batchExecAction('powerOn')}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_powerOn')}
            >
              开电
            </Button>
            <Tooltip placement="top"  title='重启[待部署]物理机并从网卡启动进入PXE'>
            <Button
              onClick={() => this.batchExecAction('networkBoot')}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_networkBoot')}
            >
              从网卡启动<Icon type='question-circle-o' />
            </Button>
            </Tooltip>
            <Button
              onClick={() => this.batchExecAction('reAccess')}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_oob_re_access')}
            >
              重新纳管带外
            </Button>
          </Button.Group>
          <span>
            已选 {selectedRows.length} 项
          </span>
          <Button
            onClick={this.exportOob}
            style={{ float: 'right', marginBottom: 8 }}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_oob_export')}
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
