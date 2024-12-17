import React from 'react';
import { get } from 'common/xFetch2';
import {
  Table,
  Button,
  Pagination,
  Tabs,
  Modal
} from 'antd';
import actions from '../actions';
import TableControlCell from 'components/TableControlCell';
const TabPane = Tabs.TabPane;
import { BOOT_MODE, OS_LIFECYCLE } from 'common/enums';
import { getPermissonBtn } from 'common/utils';

class MyTable extends React.Component {
  state = { visible: false };
  componentDidMount() {
    this.reload();
  }
  reload = () => {
    //获取系统模板列表
    this.props.dispatch({
      type: 'template-system/systemConfig-table/get'
    });
  };

  execAction = (name, record) => {
    if (actions[name]) {
      actions[name]({
        record,
        type: name,
        osFamily: this.props.osFamily,
        reload: () => {
          this.reload();
        }
      });
    }
  };

  showModal = () => {
    this.setState({
      visible: true,
    });
  };

  handleOk = e => {
    //console.log(e);
    this.setState({
      visible: false,
    });
  };

  handleCancel = e => {
    //console.log(e);
    this.setState({
      visible: false,
    });
  };

  getColumns = () => {
    return [
      {
        title: '名称',
        dataIndex: 'name',
        width: 200,
        render: (text, record) => {
          return <a onClick={() => this.execAction('systemDetail', record)}>{text}</a>;
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
              command: 'copySystem',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_system_template_create')
            },
            {
              name: '修改',
              command: 'editSystem',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_system_template_update')
            },
            {
              name: '删除',
              command: 'deleteSystem',
              type: 'danger',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_system_template_delete')
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

  onOsChange = (key) => {
    this.props.dispatch({
      type: 'template-system/systemConfig-table/search',
      payload: {
        family: key === 'all' ? null : key
      }
    });
  };

  changePage = page => {
    this.props.dispatch({
      type: `template-system/systemConfig-table/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `template-system/systemConfig-table/change-page-size`,
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
  render() {
    const { tableData, osFamily } = this.props;
    const { pagination } = tableData;
    return (
      <div>
        <div className='operate_btns'>
          <Button
            onClick={() => this.execAction('addSystem', {})}
            type='primary'
            icon='plus'
            style={{ marginRight: 8 }}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_system_template_create')}
          >
            新增配置
          </Button>
          <Button onClick={this.showModal}>
          说明
          </Button>          
        </div>
        <div className='clearfix' />
        <Modal
          title="装机配置说明"
          visible={this.state.visible}
          onOk={this.handleOk}
          onCancel={this.handleCancel}
          footer={null}
          width={800}        >
            <p>1.x86_64架构CPU支持与厂商组合请求适用的bootos名称：bootos_x86_64_[cpuvendor]</p>
            <p>2.当前支持的x86_64架构与厂商组合：bootos_x86_64_[intel|hygon] </p>
            <p>3.不在支持列表的CPU厂商默认请求适用的bootos名称：bootos_x86_64</p>
            <p>4.arm64架构CPU默认请求适用的bootos名称：bootos_arm64</p>
            <p>5.不在支持列表的CPU架构默认请求适用的bootos名称：bootos_default</p>
            <p>6.windows server模板默认请求适用的bootos名称：winpe2012_x86_64</p>
        </Modal>
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
