import React from 'react';
import Popup from 'components/popup';
import FileUpload from 'components/upload';

export default function action(options) {
  const getColumns = () => {
    return [
      {
        title: '机房管理单元',
        dataIndex: 'server_room_name'
      },
      {
        title: '机架编号',
        dataIndex: 'server_cabinet_number'
      },
      {
        title: '机位编号',
        dataIndex: 'number'
      },
      {
        title: '机位高度',
        dataIndex: 'height'
      },
      {
        title: '起始U数',
        dataIndex: 'beginning'
      },
      {
        title: '物理区域',
        dataIndex: 'physical_area'
      },
      {
        title: '状态',
        dataIndex: 'status'
      },
      {
        title: '备注',
        dataIndex: 'remark'
      }
    ];
  };

  Popup.open({
    title: '机位导入',
    width: 1000,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <FileUpload
        importApi='/api/cloudboot/v1/server-usites/imports'
        uploadApi='/api/cloudboot/v1/server-usites/upload'
        previewApi='/api/cloudboot/v1/server-usites/imports/previews'
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
