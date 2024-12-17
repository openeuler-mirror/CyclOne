import React from 'react';
import {
  Table,
  Button,
  Pagination,
  notification,
  Dropdown,
  Menu,
  Icon,
  Badge,
  Tooltip
} from 'antd';
import actions from '../actions';
import { USITE_STATUS_COLOR } from "common/enums";
import TableControlCell from 'components/TableControlCell';
import { getPermissonBtn } from 'common/utils';
import { renderDisplayMore } from 'common/utils';

const plainOptions = [
  { value: 'idc.name', name: '数据中心' },
  { value: 'beginning', name: '起始U数' },
  { value: 'creator', name: '创建用户' },
  { value: 'created_at', name: '创建时间' },
];

class MyTable extends React.Component {
  state = {
    checkedList: [],
    indeterminate: false,
    checkAll: false
  };
  reload = () => {
    this.props.dispatch({
      type: 'database-usite/table-data/reload'
    });
    this.props.dispatch({
      type: 'database-usite/table-data/set/selectedRows',
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
    if (selectedRows.length <= 0) {
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
        room: this.props.room,
        reload: () => {
          this.reload();
        }
      });
    }
  };

  //下载导入模板
  downloadImportTemplate = (key) => {
    if (key === 'usite') {
      window.open(`assets/files/server-usites.xlsx`);
    } else if (key === 'usite_port') {
      window.open(`assets/files/server-usites-port.xlsx`);
    }
  };

  getRowSelection = () => {
    const selectedRowKeys = this.props.tableData.selectedRowKeys;
    return {
      selectedRowKeys,
      onChange: (selectedRowKeys, selectedRows) => {
        this.props.dispatch({
          type: 'database-usite/table-data/set/selectedRows',
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
      type: `database-usite/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `database-usite/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };


  getColumns = () => {
    const columns = [
      {
        title: '机位编号',
        dataIndex: 'number',
        width:80,
        render: (text, record) => <a onClick={() => this.execAction('detail', record)}>{text}</a>
      },
      {
        title: '机架编号',
        dataIndex: 'server_cabinet',
        width:80,
        render: (text, record) => <a onClick={() => this.execAction('cabinet_detail', { id: text.id })}>{text.number}</a>
      },
      {
        title: '机房管理单元',
        dataIndex: 'server_room',
        render: (text) => <Tooltip placement="top" title={text.name}>{text.name}</Tooltip>
      },
      //{
      //  title: '数据中心',
      //  dataIndex: 'idc',
      //  render: (text, record) => <span>{text.name}</span>
      //},
      {
        title: '物理区域',
        dataIndex: 'physical_area'
      },
      {
        title: '内外网端口速率',
        dataIndex: 'la_wa_port_rate',
        width:120,
      },
      {
        title: '机位高度',
        dataIndex: 'height',
        width:80,
      },
      {
        title: '机位状态',
        dataIndex: 'status',
        width:80,
        render: type => {
          const color = USITE_STATUS_COLOR[type] ? USITE_STATUS_COLOR[type][0] : 'transparent';
          const word = USITE_STATUS_COLOR[type] ? USITE_STATUS_COLOR[type][1] : '';
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
        width: 180,
        render: (text, record) => {
          const commands = [
            {
              name: '编辑',
              command: '_update',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_server_usite_update')
            },
            {
              name: '删除',
              command: '_delete',
              type: 'danger',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_server_usite_delete')
            },
            {
              name: '端口删除',
              command: 'deletePort',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_server_usite_delete_port')
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
    const menu_excel = (
      <Menu onClick={(e) => this.downloadImportTemplate(e.key)}>
        <Menu.Item disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_usite_download')}
          key='usite'
        >机位模板</Menu.Item>
        <Menu.Item disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_usite_download_port')}
          key='usite_port'
        >机位关联端口模板</Menu.Item>
      </Menu>
    );
    const menu = (
      <Menu onClick={(e) => this.execAction(e.key)}>
        <Menu.Item disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_usite_import')}
          key='_import'
        >机位导入</Menu.Item>
        <Menu.Item disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_usite_import_port')}
          key='_importPort'
        >机位关联端口导入</Menu.Item>
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
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_usite_create')}
          >
            新增
          </Button>
          <Button
            onClick={() => this.batchExecAction('remark')}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_usite_remark')}
            >
            备注
            </Button>           
          <Button
            style={{ marginRight: 8 }}
            onClick={() => this.batchExecAction('changeStatus')}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_server_usite_status')}
          >
            更新机位状态
          </Button>
          <span>
            已选 { selectedRows.length } 项
          </span>
          <span className='pull-right'>
            <Dropdown overlay={menu_excel}>
              <Button style={{ borderBottomRightRadius: 0, borderTopRightRadius: 0 }}>下载导入模板<Icon type='down' /></Button>
            </Dropdown>
            <Dropdown overlay={menu}>
              <Button style={{ borderBottomLeftRadius: 0, borderTopLeftRadius: 0, marginLeft: -1 }}>导入 <Icon type='down' /></Button>
            </Dropdown>
          </span>
          <span className='pull-right'>{renderDisplayMore(this, plainOptions)}</span>
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
