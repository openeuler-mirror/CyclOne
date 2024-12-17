import React from 'react';
import { get } from 'common/xFetch2';
import {
  Table,
  Button,
  Pagination,
  notification,
  Badge
} from 'antd';
const ButtonGroup = Button.Group;
import actions from '../actions';
import TableControlCell from 'components/TableControlCell';
import { IDC_USAGE, IDC_STATUS_COLOR } from "common/enums";
import { getPermissonBtn } from 'common/utils';

class MyTable extends React.Component {

  reload = () => {
    this.props.dispatch({
      type: 'database-idc/table-data/reload'
    });
    this.props.dispatch({
      type: 'database-idc/table-data/set/selectedRows',
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
        type: name,
        reload: () => {
          this.reload();
        }
      });
    }
  };

  //下载导入模板
  downloadImportTemplate = () => {
    window.open('assets/files/idcs.xlsx');
  };


  getRowSelection = () => {
    const selectedRowKeys = this.props.tableData.selectedRowKeys;
    return {
      selectedRowKeys,
      onChange: (selectedRowKeys, selectedRows) => {
        this.props.dispatch({
          type: 'database-idc/table-data/set/selectedRows',
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
      type: `database-idc/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `database-idc/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };


  getColumns = () => {
    return [
      {
        title: '数据中心',
        dataIndex: 'name',
        render: (text, record) => <a onClick={() => this.execAction('detail', record)}>{text}</a>
      },
      {
        title: '用途',
        dataIndex: 'usage',
        render: (text) => <span>{IDC_USAGE[text]}</span>
      },
      {
        title: '一级机房',
        dataIndex: 'first_server_room',
        render: (text) => <span>{(text || []).map(item => item.name).join('，')}</span>
      },
      {
        title: '供应商',
        dataIndex: 'vendor'
      },
      {
        title: '状态',
        dataIndex: 'status',
        render: type => {
          const color = IDC_STATUS_COLOR[type] ? IDC_STATUS_COLOR[type][0] : 'transparent';
          const word = IDC_STATUS_COLOR[type] ? IDC_STATUS_COLOR[type][1] : '';
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
        dataIndex: 'updated_at'
      },
      {
        title: '操作',
        dataIndex: 'operate',
        render: (text, record) => {
          const commands = [
            {
              name: '编辑',
              command: '_update',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_idc_update')
            },
            {
              name: '删除',
              command: '_delete',
              type: 'danger',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_idc_delete')
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
  };


  render() {
    const { tableData } = this.props;
    const { loading, pagination, selectedRows } = tableData;
    return (
      <div>
        <div className='operate_btns'>
          <Button
            onClick={() => this.execAction('_create')}
            type='primary'
            style={{ marginRight: 8 }}
            icon='plus'
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_idc_create')}
          >
            新增
          </Button>
          <Button.Group style={{ marginRight: 8 }}>
            <Button
              onClick={() => this.batchExecAction('accepted')}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_idc_accepted')}
            >
              验收
            </Button>
            <Button
              onClick={() => this.batchExecAction('production')}
              disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_idc_production')}
            >
              投产
            </Button>
            {/*<Button*/}
            {/*onClick={() => this.batchExecAction('abolished')}*/}
            {/*disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_idc_abolished')}*/}
            {/*>*/}
            {/*裁撤*/}
            {/*</Button>*/}
          </Button.Group>
          <span>
            已选 { selectedRows.length } 项
          </span>
          {/*<span className='pull-right'>*/}
          {/*<ButtonGroup style={{ marginRight: 8 }}>*/}
          {/*<Button*/}
          {/*onClick={() => this.downloadImportTemplate()}*/}
          {/*>*/}
          {/*下载导入模板*/}
          {/*</Button>*/}
          {/*<Button*/}
          {/*onClick={() => this.execAction('_import')}*/}
          {/*>*/}
          {/*导入*/}
          {/*</Button>*/}
          {/*</ButtonGroup>*/}
          {/*</span>*/}
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
