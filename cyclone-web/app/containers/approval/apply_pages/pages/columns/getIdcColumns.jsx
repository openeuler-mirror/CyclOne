import { IDC_USAGE, IDC_STATUS_COLOR } from "common/enums";
import { Badge } from 'antd';

export const getIdcColumns = () => {
  return [
    {
      title: '数据中心',
      dataIndex: 'name'
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
    }
  ];
};
