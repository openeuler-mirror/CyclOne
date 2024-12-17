import { CAB_STATUS_COLOR, CAB_TYPE, YES_NO } from "common/enums";
import { Badge } from 'antd';

export const getCabColumns = () => {
  return [
    {
      title: '数据中心',
      dataIndex: 'idc',
      render: (text) => <span>{text.name}</span>
    },
    {
      title: '机房管理单元',
      dataIndex: 'server_room',
      render: (text) => <span>{text.name}</span>
    },
    {
      title: '网络区域',
      dataIndex: 'network_area',
      render: (text, record) => <span>{text.name}</span>
    },
    {
      title: '机架编号',
      dataIndex: 'number'
    },
    {
      title: '机架状态',
      dataIndex: 'status',
      width: 100,
      render: type => {
        const color = CAB_STATUS_COLOR[type] ? CAB_STATUS_COLOR[type][0] : 'transparent';
        const word = CAB_STATUS_COLOR[type] ? CAB_STATUS_COLOR[type][1] : '';
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
      title: '类型',
      dataIndex: 'type',
      render: (text) => <span>{CAB_TYPE[text]}</span>
    },
    {
      title: '是否启用',
      dataIndex: 'is_enabled',
      render: (text) => <span className={`yes_no_status ${text === 'yes' ? 'yes_status' : 'no_status'}`}>{YES_NO[text]}</span>
    },
    {
      title: '是否开电',
      dataIndex: 'is_powered',
      render: (text) => <span className={`yes_no_status ${text === 'yes' ? 'yes_status' : 'no_status'}`}>{YES_NO[text]}</span>
    },
    {
      title: '峰值功率/W',
      dataIndex: 'max_power'
    },
    {
      title: '机架高度/U',
      dataIndex: 'height'
    },
    {
      title: '机位总数',
      dataIndex: 'usite_count'
    }
  ];
};
