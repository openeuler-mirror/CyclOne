import React from 'react';
import Popup from 'components/popup';
import FileUpload from 'components/upload';

export default function action(options) {
  const getColumns = () => {
    return [
      {
        title: '序列号SN',
        dataIndex: 'sn'
      },
      {
        title: '目标IDC名称',
        dataIndex: 'dst_idc_name'
      },
      {
        title: '目标机房管理单元',
        dataIndex: 'dst_server_room_name'
      },
      {
        title: '目标机架编号',
        dataIndex: 'dst_cabinet_number'
      },
      {
        title: '目标机位编号',
        dataIndex: 'dst_usite_number'
      },
      {
        title: '目标库房管理单元',
        dataIndex: 'dst_store_room_name'
      },
      {
        title: '目标虚拟货架编号',
        dataIndex: 'dst_virtual_cabinet_number'
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
        // importApi='/api/cloudboot/v1/approvals/devices/migrations/imports'
        uploadApi='/api/cloudboot/v1/approvals/devices/migrations/upload'
        previewApi='/api/cloudboot/v1/approvals/devices/migrations/imports/previews'
        getColumns={getColumns}
        onSuccess={(data) => {
          Popup.close();
          options.reload(data);
        }}
        onCancel={() => {
          Popup.close();
        }}
      />
    )
  });
}
