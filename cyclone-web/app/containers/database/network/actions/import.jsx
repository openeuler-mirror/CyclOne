import React from 'react';
import Popup from 'components/popup';
import FileUpload from 'components/upload';
import { NET_STATUS } from "common/enums";

export default function action(options) {
  const getColumns = () => {
    return [
      {
        title: '网络区域名称',
        dataIndex: 'name'
      },
      {
        title: '机房管理单元',
        dataIndex: 'server_room_name'
      },
      {
        title: '关联物理区域',
        dataIndex: 'physical_area'
      },
      {
        title: '状态',
        dataIndex: 'status',
        render: (text) => <span>{NET_STATUS[text]}</span>
      }
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
        importApi='/api/cloudboot/v1/network-areas/imports'
        uploadApi='/api/cloudboot/v1/network-areas/upload'
        previewApi='/api/cloudboot/v1/network-areas/imports/previews'
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
