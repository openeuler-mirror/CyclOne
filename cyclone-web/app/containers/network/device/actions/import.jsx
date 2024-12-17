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
        title: '机房管理单元',
        dataIndex: 'server_room_name'
      },
      {
        title: '机架编号',
        dataIndex: 'server_cabinet_number'
      },
      {
        title: '固资编号',
        dataIndex: 'fixed_asset_number'
      },
      {
        title: '序列号SN',
        dataIndex: 'sn'
      },
      {
        title: '设备名称',
        dataIndex: 'name'
      },
      {
        title: '型号',
        dataIndex: 'model'
      },
      {
        title: '厂商',
        dataIndex: 'vendor'
      },
      {
        title: '操作系统',
        dataIndex: 'os'
      },
      {
        title: '类型',
        dataIndex: 'type'
      },
      {
        title: 'TOR',
        dataIndex: 'tor'
      },
      {
        title: '用途',
        dataIndex: 'usage'
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
        importApi='/api/cloudboot/v1/network/devices/imports'
        uploadApi='/api/cloudboot/v1/network/devices/upload'
        previewApi='/api/cloudboot/v1/network/devices/imports/previews'
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
