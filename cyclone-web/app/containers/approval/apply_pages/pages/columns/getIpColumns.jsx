import { NETWORK_IPS_SCOPE, YES_NO } from "common/enums";

export const getIpColumns = () => {
  return [
    {
      title: 'IP地址',
      dataIndex: 'ip'
    },
    {
      title: 'IP类别',
      dataIndex: 'ip_network',
      render: (text) => {
        const scope = text.category.indexOf('intranet') > -1 ? 'intranet' : 'extranet';
        return <span>{NETWORK_IPS_SCOPE[scope]}</span>;
      }
    },
    {
      title: '网段名称',
      dataIndex: 'ip_network',
      render: (text, record) => <span>{text.cidr}</span>
    },
    {
      title: '网段网关',
      dataIndex: 'gateway',
      render: (text, record) => <span>{record.ip_network.gateway}</span>
    },
    {
      title: '是否被使用',
      dataIndex: 'is_used',
      render: (text) => <span className={`yes_no_status ${text === 'yes' ? 'yes_status' : 'no_status'}`}>{YES_NO[text]}</span>
    },
    {
      title: 'IP作用范围',
      dataIndex: 'scope',
      render: (text) => <span>{NETWORK_IPS_SCOPE[text]}</span>
    },
    {
      title: '关联设备',
      dataIndex: 'sn'
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at'
    }
  ];
};
