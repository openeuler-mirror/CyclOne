import React from 'react';
import { get, del } from 'common/xFetch2';
import { INSTALL_TYPE } from "common/enums";
import {
  Table,
  Tooltip,
  Button,
  Icon,
  Pagination,
  notification,
  Badge,
  Progress
} from 'antd';
import { Link } from 'react-router';
import actions from '../actions';
import { DEVICE_INSTALL_STATUS_COLOR, DEVICE_INSTALL_TYPE } from 'common/enums';
import showExecuteResult from 'components/execute-result-modal';
import { getPermissonBtn } from 'common/utils';
import oobAction from '../../common/oob';
import { renderDisplayMore } from 'common/utils';

const plainOptions = [
  { value: 'dhcp_token', name: 'DHCP_TOKEN' , render: (text) => { return <Tooltip placement="top" title={text}>{text}</Tooltip>}},
  { value: 'tor', name: 'TOR' , render: (text) => { return <Tooltip placement="top" title={text}>{text}</Tooltip>}},
  { value: 'install_type', name: '安装类型' , render: (text) => { return <span>{INSTALL_TYPE[text]}</span>}},
];

class MyTable extends React.Component {
  state = {
    checkedList: [],
    indeterminate: true,
    checkAll: false
  };
  
  reload = () => {
    this.props.dispatch({
      type: 'device-setting/table-data/reload'
    });
    this.props.dispatch({
      type: 'device-setting/statistics/get'
    });
  };

  execAction = (name, ifSelected) => {
    let flag = false;
    const { tableData } = this.props;
    const records = tableData.selectedRows;
    if (!ifSelected && records.length < 1) {
      return notification.error({ message: '请至少选择一台设备' });
    }
    if (name === 'reInstall') {
      records.forEach(tr => {
        if (tr.status === 'success') {
          flag = true;
          notification.error({ message: `设备【${tr.sn}】已部署成功，重新部署需走审批流程` });
        }
      });
      if (flag) {
        return;
      }
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
          this.props.loadNum();
        }
      });
    }
  };

  getLog = (record, self) => {
    let logInfo = '';
    get(`/api/cloudboot/v1/devices/${record.id}/installations/logs`).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      const data = res.content.logs || [];
      const logoData = [];
      if (data.length > 0) {
        data.forEach(item => {
          logoData.push(item.updated_at + ': ' + '【' + item.title + '】 ' + item.content);
        });
        logInfo = logoData.join('\n');
        self.setState({ logInfo: logInfo });
      }
    });
  };

  release = (record) => {
    del(`/api/cloudboot/v1/devices/${record.sn}/limiters/tokens`).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: res.message });
    });
  };

  getColumns = (self) => {
    const { status } = this.props;

    let columns = [
      {
        title: '固资编号',
        dataIndex: 'device',
        render: (text) => <span>{text.fixed_asset_number}</span>
      },
      {
        title: '序列号',
        dataIndex: 'sn',
        render: (text, record) => {
          return <Link to={`/device/detail/${text}`}>{text}</Link>;
        }
      },
      {
        title: '设备类型',
        dataIndex: 'device',
        render: (text) => <span>{text.Category}</span>
      },
      {
        title: '机房管理单元',
        dataIndex: 'server_room',
        render: (text) => <Tooltip placement="top" title={text.name}>{text.name}</Tooltip>
      },
      {
        title: '机架编号',
        dataIndex: 'server_cabinet',
        render: (text) => <span>{text.number}</span>
      },
      {
        title: '机位编号',
        dataIndex: 'server_usite',
        render: (text) => <span>{text.number}</span>
      },
      {
        title: '物理区域',
        dataIndex: 'paysical_area',
        render: (text,record) => <Tooltip placement="top" title={record.server_usite ? record.server_usite.physical_area : ''}>
          {record.server_usite ? record.server_usite.physical_area : ''}
        </Tooltip>
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
        title: 'RAID类型',
        dataIndex: 'hardware_template',
        render: (text) => <Tooltip placement="top" title={text.name}>{text.name}</Tooltip>
      },
      {
        title: '操作系统',
        dataIndex: 'image_template',
        render: (text, record) => <span>{text.id !== 0 ? text.name : record.system_template.name}</span>
      },         
      {
        title: '部署状态',
        dataIndex: 'status',
        width: 100,
        render: type => {
          const color = DEVICE_INSTALL_STATUS_COLOR[type] ? DEVICE_INSTALL_STATUS_COLOR[type][0] : 'white';
          const word = DEVICE_INSTALL_STATUS_COLOR[type] ? DEVICE_INSTALL_STATUS_COLOR[type][1] : '';
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
        title: '部署进度',
        dataIndex: 'install_progress',
        width: 100,
        render: (text, record) => <Progress percent={(text * 100).toFixed(0)} size='small' status='active' />
      },
      {
        title: '部署日志',
        dataIndex: 'log',
        width: 50,
        render: (text, record) => <a onClick={ev => {
          showExecuteResult({
            record: record,
            getData: this.getLog,
            title: '查看部署日志',
            reload: this.reload
          });
        }}
        >查看</a>
      },
      {
         title: '带外',
         dataIndex: 'oob',
         render: (text, record) => <a href='javascript:;' onClick={() => oobAction(record.sn)}>查看</a>
       },  
    ];
    //if (status === 'failure') {
    //  columns.push({
    //    title: '操作',
    //    dataIndex: 'operate',
    //    width: 100,
    //    render: (text, record) => <a disabled={!record.dhcp_token} onClick={() => this.release(record)}>释放</a>
    //  });
    //}
    if (self && self.state.checkedList.length > 0) {
      self.state.checkedList.forEach(data => {
        const content = {
          title: data.name,
          dataIndex: data.value
        };
        if (data.render) {
          content.render = (text, record) => data.render(text, record, self);
        }
        columns.push(content);
      });
    }    
    return columns;
  };


  changePage = page => {
    this.props.dispatch({
      type: `device-setting/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `device-setting/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };

  getRowSelection = () => {
    const selectedRowKeys = this.props.tableData.selectedRowKeys;
    return {
      selectedRowKeys,
      onChange: (selectedRowKeys, selectedRows) => {
        this.props.dispatch({
          type: 'device-setting/table-data/set/selectedRows',
          payload: {
            selectedRowKeys,
            selectedRows
          }
        });
      }
    };
  };

  render() {
    const { tableData, status } = this.props;
    const { loading, pagination } = tableData;
    return (
      <div>
        <div className='operate_btns'>
          {
            (status === 'failure') &&
            <Button
              onClick={() => this.execAction('reInstall')}
              type='primary'
              style={{ marginRight: 8 }}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_device_setting_reInstall')}
            >
              重新部署
            </Button>
          }
          {
            (status === 'pre_install' || status === 'installing') &&
            <Button
              onClick={() => this.execAction('cancelInstall')}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_device_setting_cancelInstall')}
              style={{ marginRight: 8 }}
            >
              取消部署
            </Button>
          }
          <span className='pull-right'>
          <Button
              onClick={() => {
                this.reload();
              }}
              icon='reload'
              style={{ marginBottom: 8 , marginRight: 8 }}
            >
            </Button>
          {
            (status === 'failure') &&                    
            <Button
            onClick={() => this.execAction('reAccess')}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_oob_re_access')}
            >
              重新纳管带外
            </Button>            
          }
          {
            (status === 'failure') &&          
            <Tooltip placement="top"  title='重启[待部署]物理机并从网卡启动进入PXE'>
            <Button
              onClick={() => this.execAction('networkBoot')}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_networkBoot')}
            >
              从网卡启动<Icon type='question-circle-o' />
            </Button>
            </Tooltip>
          }          
          {
            (status === 'failure') &&
            <Tooltip placement="top"  title='仅设置部署状态为成功，注意确认系统部署情况'>
            <Button
              onClick={() => this.execAction('finshInstall')}
              style={{ marginRight: 8 }}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_device_setting_reInstall')}
            >
              完成部署
            </Button>
            </Tooltip>
          }
          {
            (status === 'failure') &&
            <Tooltip placement="top"  title='释放部署队列令牌 DHCP_TOKEN'>
            <Button
              onClick={() => this.execAction('limitersToken')}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_device_setting_delete')}
              style={{ marginRight: 8 }}
            >
              一键释放
            </Button>
            </Tooltip>
          }
          {renderDisplayMore(this, plainOptions)}
          {
            (status === 'failure') &&
            <Button
              style={{ marginLeft: 16 }}
              type='danger'
              onClick={() => this.execAction('powerOff')}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_powerOff')}
            >
              关电
            </Button>
          }
          {
            (status === 'failure') &&
            <Tooltip placement="top"  title='删除部署记录并回收IP,注意操作风险'>
            <Button
              style={{ marginRight: 8, marginLeft: 16}}
              type='danger'
              onClick={() => this.execAction('deleteDevice')}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_device_setting_delete')}
            >
              删除
            </Button>
            </Tooltip>
          }                   
          </span>
        </div>
        <div className='clearfix' />
        <Table
          rowKey={'id'}
          columns={this.getColumns(this)}
          dataSource={tableData.list}
          loading={loading}
          pagination={false}
          defaultPageSize={3}
          rowSelection={this.getRowSelection()}
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
