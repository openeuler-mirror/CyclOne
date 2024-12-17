import React from 'react';
import { Link } from 'react-router';
import { OPERATION_STATUS_COLOR, OOB_ACCESSIBLE, POWER_STATUS, OOB_STATUS_COLOR } from "common/enums";
import { Badge, Icon ,Tooltip} from 'antd';

import idcAction from 'containers/database/common/idc/detail';
import roomAction from 'containers/database/common/room/detail';
import cabinetAction from 'containers/database/common/cabinet/detail';
import usiteAction from 'containers/database/common/usite/detail';
import oobAction from './oob';

export const plainOptions = [
  { value: 'arch', name: '硬件架构' },
  { value: 'tor', name: 'TOR' },
  {
    value: 'idc', name: '数据中心', render: (text, record, snLink) => {
      if (snLink) {
        return <a onClick={() => idcAction({ records: { id: text.id }, type: 'idc_detail' })}>{text.name}</a>;
      }
      return <span>{text.name}</span>;
    }
  },
  {
    value: 'store_room',
    name: '库房管理单元',
    render: (text) => {
      return <Tooltip placement="top" title={text.name}>{text.name}</Tooltip>;
    }
  },
  {
    value: 'virtual_cabinets',
    name: '货架',
    render: (text) => {
      return <Tooltip placement="top" title={text.name}>{text.number}</Tooltip>;
    }
  },
  { value: 'model', name: '设备型号' },
  { value: 'vendor', name: '厂商' },
  {
    value: 'power_status', name: '电源状态', render: (text, record, snLink, self) => {
      return <div>
        <span className={`yes_no_status ${text === 'power_on' ? 'yes_status' : 'no_status'}`}>{POWER_STATUS[text] || '无'}</span>
        &nbsp;&nbsp;&nbsp;
        {
          snLink && <a>
            <Icon
              onClick={() => {
                self.changePowerStatus(record);
              }}
              type='reload'
            >
            </Icon>
          </a>
        }
      </div>;
    }
  }
];

export const getColumns = (self, snLink) => {

  const columns = [
    {
      title: '固资编号',
      dataIndex: 'fixed_asset_number',
      render: (text) => <Tooltip placement="top" title={text}>{text}</Tooltip>
    },
    {
      title: '序列号',
      dataIndex: 'sn',
      render: (text, record) => {
        if (snLink) {
          return <Tooltip placement="top" title={text}>
            <Link to={`/device/detail/${text}`}>{text}</Link>
          </Tooltip>
        }
        return (text) => <Tooltip placement="top" title={text}>{text}</Tooltip>;
      }
    },
    {
      title: '机房管理单元',
      dataIndex: 'server_room',
      render: (text) => {
        if (snLink) {
          return <Tooltip placement="top" title={text.name}>
              <a onClick={() => roomAction({ records: { id: text.id }, type: 'room_detail' })}>
                {text.name}
              </a>
            </Tooltip>;
        }
        return <Tooltip placement="top" title={text.name}>{text.name}</Tooltip>;
      }
    },
    {
      title: '机架编号',
      dataIndex: 'server_cabinet',
      render: (text) => {
        if (text && text !== 'null') {
          if (snLink) {
            return <a onClick={() => cabinetAction({ records: { id: text.id }, type: 'cabinet_detail' })}>{text.number}</a>;
          }
          return <span>{text.number}</span>;
        }
      }
    },
    {
      title: '机位编号',
      dataIndex: 'server_usite',
      render: (text) => {
        if (text && text !== 'null') {
          if (snLink) {
            return <a onClick={() => usiteAction({ records: { id: text.id }, type: 'detail', room: { loading: false, data: [text] } })}>{text.number}</a>;
          }
          return <span>{text.number}</span>;
        }
      }
    },
    {
      title: '物理区域',
      dataIndex: 'physical_area',
      //render: (t, record) => record.server_usite ? record.server_usite.physical_area : ''
      render: (text,record) => <Tooltip placement="top" title={record.server_usite ? record.server_usite.physical_area : ''}>
        {record.server_usite ? record.server_usite.physical_area : ''}
      </Tooltip>
    },
    {
      title: '内网 IP',
      dataIndex: 'intranet_ip'
    },
    {
      title: '外网 IP',
      dataIndex: 'extranet_ip'
    },
    {
      title: '系统',
      dataIndex: 'os'
    },
    {
       title: '带外',
       dataIndex: 'oob',
       render: (text, record) => <a href='javascript:;' onClick={() => oobAction(record.sn)}>查看</a>
     },
     {
       title: '带外状态',
       dataIndex: 'oob_accessible',
      render: type => {
        const color = OOB_STATUS_COLOR[type] ? OOB_STATUS_COLOR[type][0] : 'transparent';
        const word = OOB_STATUS_COLOR[type] ? OOB_STATUS_COLOR[type][1] : '';
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
      title: '用途',
      dataIndex: 'usage'
    },
    {
      title: '设备类型',
      dataIndex: 'category',
      //render: (t) => {
      //  if (t === 'SpecialDev') {
      //    return '特殊设备';
      //  }
      //  return t;
      //}
    },
    {
      title: '运营状态',
      dataIndex: 'operation_status',
      render: type => {
        const color = OPERATION_STATUS_COLOR[type] ? OPERATION_STATUS_COLOR[type][0] : 'transparent';
        const word = OPERATION_STATUS_COLOR[type] ? OPERATION_STATUS_COLOR[type][1] : '';
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
    }
  ];
  if (self && self.state.checkedList.length > 0) {
    self.state.checkedList.forEach(data => {
      const content = {
        title: data.name,
        dataIndex: data.value
      };
      if (data.render) {
        content.render = (text, record) => data.render(text, record, snLink, self);
      }
      columns.push(content);
    });
  }
  return columns;
};

