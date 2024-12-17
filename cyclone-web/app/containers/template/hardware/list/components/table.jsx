import React from 'react';
import { get } from 'common/xFetch2';
import {
  Table
} from 'antd';
import { Link } from 'react-router';
import TableControlCell from 'components/TableControlCell';
import { getPermissonBtn } from 'common/utils';

class MyTable extends React.Component {

  getColumns = () => {
    return [
      {
        title: '名称',
        dataIndex: 'name',
        render: (text, record) => {
          return <Link to={`/template/hardware/detail/${record.id}`}>{text}</Link>;
        }
      },
      // {
      //   title: '厂商',
      //   dataIndex: 'vendor',
      //   width: 200
      // },
      // {
      //   title: '产品型号',
      //   dataIndex: 'model_name',
      //   width: 200
      // },
      // {
      //   title: '配置内容',
      //   dataIndex: 'data_config',
      //   width: 200,
      //   render: (text, record) => {
      //     const data = record.data || [];
      //     return <span>{[...new Set(data.map(d => d.category))].join(' | ')}</span>;
      //   }
      // },
      {
        title: '创建时间',
        dataIndex: 'created_at'
      },
      {
        title: '修改时间',
        dataIndex: 'updated_at'
      },
      {
        title: '操作',
        dataIndex: 'operation',
        render: (text, record) => {
          let commands = [
            {
              name: '克隆',
              command: 'copyTemplate',
              disabled: !getPermissonBtn(this.props.userInfo.permissions, 'button_hardware_template_create')
            },
            {
              name: '修改',
              command: 'editTemplate',
              disabled: record.builtin === 'yes' || !getPermissonBtn(this.props.userInfo.permissions, 'button_hardware_template_update')
            ,
              message: '系统内置'
            },
            {
              name: '删除',
              command: 'deleteTemplate',
              type: 'danger',
              disabled: record.builtin === 'yes' || !getPermissonBtn(this.props.userInfo.permissions, 'button_hardware_template_delete'),
              message: '系统内置'
            }
          ];
          return (
            <TableControlCell
              commands={commands}
              record={record}
              execCommand={command => {
                this.props.execAction(command, record);
              }}
            />
          );
        }
      }
    ];
  };

  render() {
    const { dataSource, loading } = this.props;
    return (
      <div>
        <Table
          rowKey={'id'}
          columns={this.getColumns()}
          dataSource={dataSource}
          loading={loading}
          pagination={false}
        />
      </div>
    );
  }
}

export default MyTable;
