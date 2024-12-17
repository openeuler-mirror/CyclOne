import React from 'react';
import {
  Table,
  Tooltip,
  Icon,
  Pagination,
  Button,
  Form,
} from 'antd';
import actions from '../actions';
import TableControlCell from 'components/TableControlCell';
import { renderDisplayMore } from 'common/utils';
import { getPermissonBtn } from 'common/utils';


const FormItem = Form.Item;

const plainOptions = [
  { value: 'idc.name', name: '数据中心' },
  { value: 'os', name: '操作系统' },
  { value: 'vendor', name: '厂商' },
  { value: 'model', name: '型号' },
  { value: 'type', name: '类型' },
];

class MyTable extends React.Component {
  state = {
    checkedList: [],
    indeterminate: false,
    checkAll: false
  };
  reload = () => {
    this.props.dispatch({
      type: 'network-device/table-data/reload'
    });
    this.props.dispatch({
      type: 'network-device/table-data/set/selectedRows',
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

  //操作入口
  execAction = (name, records) => {
    if (actions[name]) {
      actions[name]({
        records,
        initialValue: records,
        idc: this.props.idc,
        type: name,
        reload: () => {
          this.reload();
        }
      });
    }
  };

  getRowSelection = () => {
    const selectedRowKeys = this.props.tableData.selectedRowKeys;
    return {
      selectedRowKeys,
      onChange: (selectedRowKeys, selectedRows) => {
        this.props.dispatch({
          type: 'network-device/table-data/set/selectedRows',
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
      type: `network-device/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `network-device/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };

  getColumns = () => {
    const columns = [
      {
        title: '固资编号',
        dataIndex: 'fixed_asset_number'
      },
      {
        title: '序列号',
        dataIndex: 'sn',
        render: (text, record) => <a onClick={() => this.execAction('_detail', record)}>{text}</a>
      },
      {
        title: '名称',
        dataIndex: 'name',
        render: (name) => <Tooltip placement="topLeft" title={name}>{name}</Tooltip>
      },
      {
        title: '机房管理单元',
        dataIndex: 'server_room',
        render: (text, record) => <Tooltip placement="top" title={text.name}>
            <a onClick={() => this.execAction('room_detail', { id: text.id })}>
              {text.name}
            </a>
          </Tooltip>
      },
      {
        title: '机架编号',
        dataIndex: 'server_cabinet',
        render: (text, record) => <a onClick={() => this.execAction('cabinet_detail', { id: text.id })}>{text.number}</a>
      },
      {
        title: '用途',
        dataIndex: 'usage'
      },
      {
        title: '状态',
        dataIndex: 'status'
      },      
      {
        title: 'TOR',
        dataIndex: 'tor',
        render: (text) => <Tooltip placement="top" title={text}>{text}</Tooltip>
      },
      {
        title: '操作',
        dataIndex: 'operate',
        width: 80,
        render: (text, record) => {
          const commands = [
            {
              name: '删除',
              command: '_delete',
              type: 'danger',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_network_device_delete')
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

  //下载导入模板
  downloadImportTemplate = () => {
    window.open('assets/files/network_device_import.xlsx');
  };

  render() {
    const { tableData } = this.props;
    const { loading, pagination } = tableData;

    return (
      <div>
        <div className='operate_btns'>
          <Button
            onClick={() => this.execAction('_create')}
            type='primary'
            icon='plus'
            style={{ marginRight: 8 }}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_network_device_create')}
          >
            新增
          </Button>
          <span className='pull-right'>
            <Button.Group style={{ marginRight: 8 }}>
              <Button
                onClick={() => this.downloadImportTemplate()}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_download')}
              >
                下载导入模板
              </Button>
              <Button
                onClick={() => this.execAction('_import')}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_import')}
              >
                导入
              </Button>
              <Button
                onClick={() => this.batchExecAction('_batchdelete')}
                type='danger'
                style={{ marginLeft: 8 }}
                disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_network_device_delete')}
              >
                删除
              </Button>              
            </Button.Group>
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
