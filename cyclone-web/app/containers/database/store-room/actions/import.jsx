import React from 'react';
import Popup from 'components/popup';
import FileUpload from 'components/upload';

export default function action(options) {
  const getColumns = () => {
    return [
      {
        title: '数据中心',
        dataIndex: 'idc_name'
      },
      {
        title: '一级机房名称',
        dataIndex: 'first_server_room'
      },
      {
        title: '库房名称',
        dataIndex: 'name'
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
        title: '库房负责人',
        dataIndex: 'store_room_manager'
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
        importApi='/api/cloudboot/v1/store-room/imports'
        uploadApi='/api/cloudboot/v1/store-room/upload'
        previewApi='/api/cloudboot/v1/store-room/imports/previews'
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
