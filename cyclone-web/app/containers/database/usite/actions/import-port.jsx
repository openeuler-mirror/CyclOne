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
        title: '管理网交换机设备名称-端口',
        dataIndex: 'oobnet_switch_name_port'
      },
      {
        title: '内网交换机设备名称-端口',
        dataIndex: 'intranet_switch_name_port'
      },
      {
        title: '外网交换机设备名称-端口',
        dataIndex: 'extranet_switch_name_port'
      },
      {
        title: '内外网端口速率',
        dataIndex: 'la_wa_port_rate'
      }      
    ];
  };
  Popup.open({
    title: '机位关联端口导入',
    width: 1000,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <FileUpload
        importApi='/api/cloudboot/v1/server-usites/ports/imports'
        uploadApi='/api/cloudboot/v1/server-usites/ports/upload'
        previewApi='/api/cloudboot/v1/server-usites/ports/imports/previews'
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
