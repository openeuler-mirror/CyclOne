import { NET_STATUS_COLOR } from "common/enums";
import { Badge } from 'antd';

export const getNetworkColumns = () => {
  return [
    {
      title: '网络区域名称',
      dataIndex: 'name'
    },
    {
      title: '机房管理单元',
      dataIndex: 'server_room',
      render: (t) => t.name
    },
    {
      title: '关联物理区域',
      dataIndex: 'physical_area',
      render: (text) => <span>{(text || []).map(it => it.name).join('，')}</span>
    },
    {
      title: '状态',
      dataIndex: 'status',
      render: type => {
        const color = NET_STATUS_COLOR[type] ? NET_STATUS_COLOR[type][0] : 'transparent';
        const word = NET_STATUS_COLOR[type] ? NET_STATUS_COLOR[type][1] : '';
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
};
