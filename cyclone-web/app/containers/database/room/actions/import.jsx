import React from 'react';
import Popup from 'components/popup';
import FileUpload from 'components/upload';

export default function action(options) {
  const getColumns = () => {
    return [
      {
        title: '机房管理单元',
        dataIndex: 'name'
      },
      {
        title: '数据中心',
        dataIndex: 'idc_name'
      },
      {
        title: '所属一级机房',
        dataIndex: 'first_server_room_name'
      },
      {
        title: '城市',
        dataIndex: 'city'
      },
      {
        title: '机房地址',
        dataIndex: 'address'
      },
      {
        title: '机房负责人',
        dataIndex: 'server_room_manager'
      },
      {
        title: '供应商负责人',
        dataIndex: 'vendor_manager'
      },
      {
        title: '网络资产负责人',
        dataIndex: 'network_asset_manager'
      },
      {
        title: '7*24小时保障电话',
        dataIndex: 'support_phone_number'
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
        importApi='/api/cloudboot/v1/server-rooms/imports'
        uploadApi='/api/cloudboot/v1/server-rooms/upload'
        previewApi='/api/cloudboot/v1/server-rooms/imports/previews'
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
