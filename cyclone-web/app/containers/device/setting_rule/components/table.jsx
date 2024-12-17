import React from 'react';
import {
  Table,
  Tooltip,
  Button,
  Pagination,
  notification,
  Modal
} from 'antd';
import actions from '../actions';
import TableControlCell from 'components/TableControlCell';
import { getPermissonBtn } from 'common/utils';

class MyTable extends React.Component {
  state = { visible: false };
  reload = () => {
    this.props.dispatch({
      type: 'device-setting-rules/table-data/reload'
    });
    this.props.dispatch({
      type: 'device-setting-rules/table-data/set/selectedRows',
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

  getColumns = () => {

    let columns = [
      {
        title: '规则前件',
        dataIndex: 'condition',
        width: 800,
        render: (text) => {
          return <Tooltip placement="top" title={text}>{text}</Tooltip>
        }
      },
      {
        title: '规则推论',
        dataIndex: 'action',
      },
      {
        title: '规则类别',
        dataIndex: 'rule_category',
      },
      {
        title: '操作',
        dataIndex: 'operate',
        render: (text, record) => {
          const commands = [
            {
              name: '修改',
              command: '_update',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_device_setting_rule_update')
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
      },
    ];

    return columns;
  };

  changePage = page => {
    this.props.dispatch({
      type: `device-setting-rules/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `device-setting-rules/table-data/change-page-size`,
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
          type: 'device-setting-rules/table-data/set/selectedRows',
          payload: {
            selectedRowKeys,
            selectedRows
          }
        });
      }
    };
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


  render() {
    const { tableData } = this.props;
    const { loading, pagination } = tableData;

    return (
      <div>
        <div className='operate_btns'>
          <Button
            onClick={() => this.execAction('_create', {})}
            type='primary'
            icon='plus'
            style={{ marginRight: 8 }}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_device_setting_rule_create')}
          >
            新增
          </Button>
          <Button onClick={this.showModal}>
            说明
          </Button>        
          <span className='pull-right'>
          <Button
            type='danger'
            onClick={() => this.batchExecAction('_delete')}
            style={{ marginRight: 8 }}
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_device_setting_rule_delete')}
          >
            删除
          </Button>
            <Button
              onClick={() => {
                this.reload();
              }}
              icon='reload'
              style={{ marginBottom: 8 }}
            >
            </Button>
          </span>
        </div>
        <div className='clearfix' />
        <Modal
          title="规则说明"
          visible={this.state.visible}
          onOk={this.handleOk}
          onCancel={this.handleCancel}
          footer={null}
          width={800}        >
            <p>1.元规则：获取设备属性[category-设备类型|vendor-厂商（小写英文）|physical_area-物理区域|is_fiti_eco_product-是否金融信创生态产品]值+逻辑判断操作符[equal|contains|in]进行推导</p>
            <p>2.多个元规则通过逻辑操作符[and|or]进行组合，优先级：or 大于 and</p>
            <p>3.相同类型[os|raid|network]的规则按排列先后顺序匹配推导</p>
            <p>4.当所有类型[os|raid|network]的规则都推导成功时返回完整的装机参数进入部署</p>
        </Modal>
        <Table
          rowKey={'id'}
          columns={this.getColumns()}
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
