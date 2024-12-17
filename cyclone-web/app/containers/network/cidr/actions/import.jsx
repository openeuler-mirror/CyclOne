import React from 'react';
import Popup from 'components/popup';
import FileUpload from 'components/upload';
import { IP_NETWORK_CATEGORY } from 'common/enums';

export default function action(options) {
  const getColumns = () => {
    return [
      {
        title: '机房管理单元名称',
        dataIndex: 'server_room_name'
      },
      {
        title: '网段名',
        dataIndex: 'cidr'
      },
      {
        title: '类型',
        dataIndex: 'category',
        render: (text) => <span>{IP_NETWORK_CATEGORY[text]}</span>
      },
      {
        title: '掩码',
        dataIndex: 'netmask'
      },
      {
        title: '网关',
        dataIndex: 'gateway'
      },
      {
        title: 'IP池',
        dataIndex: 'ip_pool'
      },
      {
        title: 'PXE-IP池',
        dataIndex: 'pxe_pool'
      },
      {
        title: '覆盖交换机',
        dataIndex: 'switchs'
      },
      {
        title: 'VLAN',
        dataIndex: 'vlan'
      },
      {
        title: 'IP版本',
        dataIndex: 'version'
      },      
    ];
  };


  Popup.open({
    title: '导入数据',
    width: 1500,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <FileUpload
        importApi='/api/cloudboot/v1/ip-networks/imports'
        uploadApi='/api/cloudboot/v1/ip-networks/upload'
        previewApi='/api/cloudboot/v1/ip-networks/imports/previews'
        getColumns={getColumns}
        onSuccess={() => {
          Popup.close();
          options.reload();
        }}
        onCancel={() => {
          Popup.close();
        }}
      />
    )
  });
}
