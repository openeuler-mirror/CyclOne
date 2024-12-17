import { IDC_STATUS_COLOR } from "common/enums";
import { Badge } from 'antd';

export const getRoomColumns = () => {
  return [
    {
      title: '机房管理单元',
      dataIndex: 'name'
    },
    {
      title: '数据中心',
      dataIndex: 'idc',
      render: (text, record) => <span>{text.name}</span>
    },
    {
      title: '所属一级机房',
      dataIndex: 'first_server_room',
      render: (text, record) => <span>{text.name}</span>
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 90,
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
      title: '城市',
      dataIndex: 'city'
    },
    {
      title: '地址',
      dataIndex: 'address'
    },
    {
      title: '机架数',
      dataIndex: 'cabinet_count'
    },
    {
      title: '机房负责人',
      dataIndex: 'server_room_manager'
    },
    {
      title: '供应商负责人',
      dataIndex: 'vendor_manager'
    },
    {
      title: '网络资产负责人',
      dataIndex: 'network_asset_manager'
    },
    {
      title: '7*24小时保障电话',
      dataIndex: 'support_phone_number'
    }
  ];
};
