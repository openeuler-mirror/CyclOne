import React from 'react';
import { get } from 'common/xFetch2';
import {
  Table,
  Button,
  Pagination,
  Tabs
} from 'antd';
import actions from '../actions';
import TableControlCell from 'components/TableControlCell';
const TabPane = Tabs.TabPane;
import { BOOT_MODE, OS_LIFECYCLE } from 'common/enums';
import { getPermissonBtn } from 'common/utils';

class MyTable extends React.Component {

  componentDidMount() {
    this.reload();
  }
  reload = () => {
    //获取镜像模板列表
    this.props.dispatch({
      type: 'template-system/mirrorInstallTpl-table/get'
    });
  };

  execAction = (name, record) => {
    if (actions[name]) {
      actions[name]({
        record,
        userInfo: this.props.userInfo,
        osFamily: this.props.osFamily,
        type: name,
        reload: () => {
          this.reload();
        }
      });
    }
  };

  getColumns = () => {
    return [
      {
        title: '名称',
        dataIndex: 'name',
        width: 200,
        render: (text, record) => {
          return <a onClick={() => this.execAction('mirrorDetail', record)}>{text}</a>;
        }
      },
      {
        title: '启动模式',
        dataIndex: 'boot_mode',
        width: 50,
        render: (text, record) => {
          return BOOT_MODE[text];
        }
      },
      {
        title: '创建时间',
        dataIndex: 'created_at',
        width: 50
      },
      {
        title: '修改时间',
        dataIndex: 'updated_at',
        width: 50
      },
      {
        title: '架构',
        dataIndex: 'arch',
        width: 50
      },      
      {
        title: '生命周期',
        dataIndex: 'os_lifecycle',
        width: 50,
        render: (text, record) => {
          return OS_LIFECYCLE[text];
        }
      },      
      {
        title: '操作',
        dataIndex: 'SystemName',
        width: 100,
        render: (text, record) => {
          let commands = [
            {
              name: '克隆',
              command: 'copyMirror',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_mirror_template_create')
            },
            {
              name: '修改',
              command: 'editMirror',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_mirror_template_update')
            },
            {
              name: '删除',
              command: 'deleteMirror',
              type: 'danger',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_mirror_template_delete')
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


  changePage = page => {
    this.props.dispatch({
      type: `template-system/mirrorInstallTpl-table/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `template-system/mirrorInstallTpl-table/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };


  getContent = () => {
    const { tableData } = this.props;
    const { loading } = tableData;
    return (
      <Table
        rowKey={'id'}
        //scroll={{ y: 'calc(100vh - 320px)' }}
        columns={this.getColumns()}
        dataSource={tableData.list}
        loading={loading}
        pagination={false}
      />
    );
  };
  onOsChange = (key) => {
    this.props.dispatch({
      type: 'template-system/mirrorInstallTpl-table/search',
      payload: {
        family: key === 'all' ? null : key
      }
    });
  };
  render() {
    const { tableData, osFamily } = this.props;
    const { pagination } = tableData;
    return (
      <div>
        <div className='operate_btns'>
          <Button
            onClick={() => this.execAction('addMirror', {})}
            type='primary'
            icon='plus'
            style={{ marginRight: 8 }}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_mirror_template_create')}
          >
            新增配置
          </Button>
        </div>
        <div className='clearfix' />
        <Tabs
          defaultActiveKey='0'
          tabPosition='left'
          onChange={this.onOsChange}
        >
          <TabPane tab={'全部'} key={'all'}>{this.getContent()}</TabPane>
          {
            !osFamily.loading &&
            osFamily.data.map((os, index) => <TabPane tab={os.name} key={os.name}>{this.getContent()}</TabPane>)
          }
        </Tabs>
        <div>
          <Pagination
            showQuickJumper={true}
            showSizeChanger={true}
            current={pagination.page}
            pageSize={pagination.pageSize}
            total={pagination.total}
            onShowSizeChange={this.changePageSize}
            onChange={this.changePage}
            showTotal={(total) => `共 ${total} 条`}
          />
        </div>
      </div>
    );
  }
}

export default MyTable;
