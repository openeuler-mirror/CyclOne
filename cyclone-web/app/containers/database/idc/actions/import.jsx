import React from 'react';
import Popup from 'components/popup';
import FileUpload from 'components/upload';

export default function action(options) {
  const getColumns = () => {
    return [
      {
        title: '数据中心',
        dataIndex: 'name'
      },
      {
        title: '用途',
        dataIndex: 'usage'
      },
      {
        title: '一级机房',
        dataIndex: 'first_server_room'
      },
      {
        title: '供应商',
        dataIndex: 'vendor'
      },
      {
        title: '状态',
        dataIndex: 'status'
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
        importApi='/api/cloudboot/v1/idcs/imports'
        uploadApi='/api/cloudboot/v1/idcs/upload'
        previewApi='/api/cloudboot/v1/idcs/imports/previews'
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
