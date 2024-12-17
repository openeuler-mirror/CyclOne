import React from 'react';
import {
  Table,
  Button,
  Pagination,
  notification,
  Menu,
  Badge,
  Tooltip,
  Icon
} from 'antd';
const ButtonGroup = Button.Group;

import actions from '../actions';
import { CAB_STATUS_COLOR, CAB_TYPE, YES_NO } from "common/enums";
import TableControlCell from 'components/TableControlCell';
import { renderDisplayMore } from 'common/utils';
import { getPermissonBtn } from 'common/utils';


const plainOptions = [
  { value: 'enable_time', name: '启用时间' },
  { value: 'power_on_time', name: '开电时间' },
  { value: 'power_off_time', name: '关电时间' },
  { value: 'idc.name', name: '数据中心' },
  { value: 'height', name: '机架高度' },
  { value: 'usite_count', name: '机位总数' },
];

class MyTable extends React.Component {
  state = {
    checkedList: [],
    indeterminate: false,
    checkAll: false
  };
  reload = () => {
    this.props.dispatch({
      type: 'database-cabinet/table-data/reload'
    });
    this.props.dispatch({
      type: 'database-cabinet/table-data/set/selectedRows',
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

  //单个操作入口
  oneExecAction = (name) => {
    const { tableData } = this.props;
    const selectedRows = tableData.selectedRows || [];
    if (selectedRows.length !== 1) {
      return notification.error({ message: '请选择一条数据' });
    }
    this.execAction(name, selectedRows[0]);
  };

  //操作入口
  execAction = (name, records) => {
    if (actions[name]) {
      actions[name]({
        records,
        initialValue: records,
        type: name,
        network: this.props.network,
        reload: () => {
          this.reload();
        }
      });
    }
  };

  //下载导入模板
  downloadImportTemplate = () => {
    window.open('assets/files/server-cabinets.xlsx');
  };


  getRowSelection = () => {
    const selectedRowKeys = this.props.tableData.selectedRowKeys;
    return {
      selectedRowKeys,
      onChange: (selectedRowKeys, selectedRows) => {
        this.props.dispatch({
          type: 'database-cabinet/table-data/set/selectedRows',
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
      type: `database-cabinet/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `database-cabinet/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };


  getColumns = () => {
    const columns = [
      {
        title: '机架编号',
        dataIndex: 'number',
        width:80,
        render: (text, record) => <a onClick={() => this.execAction('detail', record)}>{text}</a>
      },
      {
        title: '网络区域',
        dataIndex: 'network_area',
        render: (text, record) => <a onClick={() => this.execAction('network_detail', { id: text.id })}>{text.name}</a>
      },
      {
        title: '机房管理单元',
        dataIndex: 'server_room',
        render: (text) => <Tooltip placement="top" title={text.name}>{text.name}</Tooltip>
      },
      //{
      //  title: '数据中心',
      //  dataIndex: 'idc',
      //  render: (text) => <span>{text.name}</span>
      //},
      {
        title: '机架状态',
        dataIndex: 'status',
        width: 80,
        render: type => {
          const color = CAB_STATUS_COLOR[type] ? CAB_STATUS_COLOR[type][0] : 'transparent';
          const word = CAB_STATUS_COLOR[type] ? CAB_STATUS_COLOR[type][1] : '';
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
        title: '类型',
        dataIndex: 'type',
        render: (text) => <span>{CAB_TYPE[text]}</span>
      },
      {
        title: '是否启用',
        dataIndex: 'is_enabled',
        render: (text) => <span className={`yes_no_status ${text === 'yes' ? 'yes_status' : 'no_status'}`}>{YES_NO[text]}</span>
      },
      {
        title: '是否开电',
        dataIndex: 'is_powered',
        render: (text) => <span className={`yes_no_status ${text === 'yes' ? 'yes_status' : 'no_status'}`}>{YES_NO[text]}</span>
      },
      {
        title: '峰值功率/W',
        dataIndex: 'max_power'
      },
      {
        title: '电流/A',
        dataIndex: 'current'
      },
      {
        title: '网络速率/G',
        dataIndex: 'network_rate'
      },
      //{
      //  title: '机架高度/U',
      //  dataIndex: 'height'
      //},
      //{
      //  title: '机位总数',
      //  dataIndex: 'usite_count'
      //},
      {
        title: '更新时间',
        dataIndex: 'updated_at',
        width: 150
      },      
      {
        title: '备注',
        dataIndex: 'remark'
      },    
      {
        title: '操作',
        dataIndex: 'operate',
        width: 100,
        render: (text, record) => {
          const commands = [
            {
              name: '编辑',
              command: '_update',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_server_cabinet_update')
            },
            {
              name: '删除',
              command: '_delete',
              type: 'danger',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_server_cabinet_delete')
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
    if (this.state.checkedList.length > 0) {
      this.state.checkedList.forEach(data => {
        columns.push({
          title: data.name,
          dataIndex: data.value
        });
      });
    }
    return columns;
  };


  render() {
    const { tableData } = this.props;
    const { loading, pagination, selectedRows } = tableData;
    const menu = (
      <Menu onClick={(e) => this.batchExecAction(e.key)}>
        {/*<Menu.Item key='online'>验收</Menu.Item>*/}
        {/*<Menu.Item key='offline'>改造</Menu.Item>*/}
        {/*<Menu.Item key='applyUsers'>拆分</Menu.Item>*/}
        {/*<Menu.Item key='setBlack'>合并</Menu.Item>*/}
        <Menu.Item key='offline' disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_cabinet_offline')}>下线</Menu.Item>
      </Menu>
    );
    return (
      <div>
        <div className='operate_btns'>
          <Button
            onClick={() => this.execAction('_create')}
            type='primary'
            style={{ marginRight: 8 }}
            icon='plus'
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_cabinet_create')}
          >
            新增
          </Button>
          <Button.Group style={{ marginRight: 8 }}>
            <Button
              onClick={() => this.batchExecAction('enabled')}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_cabinet_enabled')}
            >
              启用
            </Button>
            <Button
              onClick={() => this.batchExecAction('powerOn')}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_cabinet_powerOn')}
            >
              开电
            </Button>
            <Tooltip placement="top"  title='机架置为[已锁定]并将关联的[空闲][预占用]机位置为[不可用]'>
              <Button
                onClick={() => this.batchExecAction('locked')}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_cabinet_locked')}
              >
                锁定<Icon type='question-circle-o' />
              </Button>
            </Tooltip>
            <Button
            onClick={() => this.batchExecAction('remark')}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_cabinet_remark')}
            >
            备注
            </Button>            
            <Button
            style={{ marginRight: 8 }}
            onClick={() => this.batchExecAction('changeType')}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_cabinet_type')}
            >
            更新机架类型
            </Button>
          </Button.Group>
          <span>
            已选 { selectedRows.length } 项
          </span>
          <span className='pull-right'>
            <ButtonGroup style={{ marginRight: 8 }}>
              <Button
                onClick={() => this.downloadImportTemplate()}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_cabinet_download')}
              >
                下载导入模板
              </Button>
              <Button
                onClick={() => this.execAction('_import')}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_cabinet_import')}
              >
                导入
              </Button>
            </ButtonGroup>
            {renderDisplayMore(this, plainOptions)}
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
