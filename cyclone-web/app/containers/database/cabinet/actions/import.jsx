import React from 'react';
import Popup from 'components/popup';
import FileUpload from 'components/upload';

export default function action(options) {
  const getColumns = () => {
    return [
      {
        title: '机架编号',
        dataIndex: 'number'
      },
      {
        title: '网络区域',
        dataIndex: 'network_area'
      },
      {
        title: '机房管理单元',
        dataIndex: 'server_room'
      },
      {
        title: '类型',
        dataIndex: 'type'
      },
      {
        title: '峰值功率',
        dataIndex: 'max_power'
      },
      {
        title: '网络速率',
        dataIndex: 'network_rate'
      },
      {
        title: '电流',
        dataIndex: 'current'
      },
      {
        title: '可用功率',
        dataIndex: 'available_power'
      },
      {
        title: '机架高度',
        dataIndex: 'height'
      },
      {
        title: '备注',
        dataIndex: 'remark'
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
        importApi='/api/cloudboot/v1/server-cabinets/imports'
        uploadApi='/api/cloudboot/v1/server-cabinets/upload'
        previewApi='/api/cloudboot/v1/server-cabinets/imports/previews'
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
