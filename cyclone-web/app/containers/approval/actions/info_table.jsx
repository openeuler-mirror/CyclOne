import React from 'react';
import {
  Table,
  Tooltip,
  Col
} from 'antd';
import { getCabColumns } from "../apply_pages/pages/columns/getCabColumns";
import { getIdcColumns } from "../apply_pages/pages/columns/getIdcColumns";
import { getRoomColumns } from "../apply_pages/pages/columns/getRoomColumns";
import { getNetworkColumns } from "../apply_pages/pages/columns/getNetworkColumns";
import { getIpColumns } from "../apply_pages/pages/columns/getIpColumns";
import { getColumns } from 'containers/device/common/colums';
import { YES_NO } from 'common/enums';


export default class MyTable extends React.Component {

  getMoveColumns = () => {
    return [
      {
        title: '固资编号',
        dataIndex: 'fixed_asset_number',
        render: (text) => <Tooltip placement="top" title={text}>{text}</Tooltip>
      }, {
        title: '序列号',
        dataIndex: 'sn',
        render: (text) => <Tooltip placement="top" title={text}>{text}</Tooltip>
      }, {
        title: '设备型号',
        dataIndex: 'model'
      }, {
        title: '数据中心',
        dataIndex: 'idc',
        render: (t, record) => {
          return <div className='table-noWrap'>
            <div>原：{record.old_idc}</div>
            <div>目标：{t}</div>
          </div>;
        }
      }, {
        title: '机房',
        dataIndex: 'server_room_name',
        render: (t, record) => {
          return <div className='table-noWrap'>
            <div>原：{!record.in_store && record.old_server_room}</div>
            <div>目标：{t}</div>
          </div>;
        }
      }, {
        title: '机架',
        dataIndex: 'server_cabinet_number',
        render: (t, record) => {
          return <div className='table-noWrap'>
            <div>原：{!record.in_store && record.old_server_cabinet}</div>
            <div>目标：{t}</div>
          </div>;
        }
      }, {
        title: '机位',
        dataIndex: 'server_usite_number',
        render: (t, record) => {
          return <div className='table-noWrap'>
            <div>原：{!record.in_store && record.old_server_usite}</div>
            <div>目标：{t}</div>
          </div>;
        }
      },
      {
        title: '虚拟库房',
        dataIndex: 'store_room_name',
        render: (t, record) => {
          return <div className='table-noWrap'>
            <div>原：{record.in_store && record.old_server_room}</div>
            <div>目标：{t}</div>
          </div>;
        }
      }, {
        title: '虚拟货架',
        dataIndex: 'virtual_cabinet_number',
        render: (t, record) => {
          return <div className='table-noWrap'>
            <div>原：{record.in_store && record.old_server_cabinet}</div>
            <div>目标：{t}</div>
          </div>;
        }
      }];
  };
  getReinstallColumns = () => {
    return [
      {
        title: '固资编号',
        dataIndex: 'fixed_asset_number'
      }, {
        title: '序列号',
        dataIndex: 'sn'
      }, {
        title: '操作系统',
        dataIndex: 'os_template_name'
      }, {
        title: 'RAID类型',
        dataIndex: 'hardware_template_name'
      }, {
        title: '是否分配外网IP',
        dataIndex: 'need_extranet_ip',
        render: (text) => <span className={`yes_no_status ${text === 'yes' ? 'yes_status' : 'no_status'}`}>{YES_NO[text]}</span>
      }, {
        title: '设备型号',
        dataIndex: 'model'
      }, {
        title: '机位编号',
        dataIndex: 'server_usite_number'
      }];
  };

  render() {
    if (!Array.isArray(this.props.data)) {
      return <Col span={24} />;
    }
    const type = this.props.type;
    let title = '';
    let columns = [];
    if (type === 'cabinet_power_off' || type === 'cabinet_offline') {
      columns = getCabColumns();
      title = '选择的机架';
    } else if (type === 'device_migration') {
      columns = this.getMoveColumns();
      title = '物理机搬迁信息';
    } else if (type === 'device_os_reinstallation' || type == 'device_recycle_reinstall') {
      columns = this.getReinstallColumns();
      title = '物理机重装信息';
    } else if (type === 'idc_abolish') {
      columns = getIdcColumns();
      title = '选择的数据中心';
    } else if (type === 'server_room_abolish') {
      columns = getRoomColumns();
      title = '选择的机房';
    } else if (type === 'network_area_offline') {
      columns = getNetworkColumns();
      title = '选择的网络区域';
    } else if (type === 'ip_unassign') {
      columns = getIpColumns();
      title = '选择的IP';
    } else {
      columns = getColumns(null, false);
      title = '选择的设备';
    }
    return (
      <div>
        <Col span={24}>
          {title}:
        </Col>
        <Table
          rowKey={'id'}
          columns={columns}
          pagination={false}
          dataSource={this.props.data}
        />
      </div>
    );
  }
}
