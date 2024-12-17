import React from 'react';
import Popup from 'components/popup';
import FileUpload from 'components/upload';

export default function action(options) {
  const getColumns = () => {
    return [
      {
        title: '固资编号',
        dataIndex: 'fixed_asset_number'
      },
      {
        title: '序列号',
        dataIndex: 'sn'
      },
      {
        title: '设备型号',
        dataIndex: 'model'
      },
      {
        title: 'CPU架构',
        dataIndex: 'arch'
      },
      {
        title: '用途',
        dataIndex: 'usage'
      },
      {
        title: '设备类型',
        dataIndex: 'category'
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
        title: '机位编号',
        dataIndex: 'server_usite_number'
      },
      {
        title: '硬件备注',
        dataIndex: 'hardware_remark'
      },
      {
        title: 'RAID结构',
        dataIndex: 'raid_remark'
      },
      {
        title: '厂商',
        dataIndex: 'vendor'
      },
      {
        title: '启用时间',
        dataIndex: 'started_at'
      },
      {
        title: '上架时间',
        dataIndex: 'onshelve_at'
      },
      {
        title: '带外',
        dataIndex: 'oob_init'
      },
      {
        title: '内网IP',
        dataIndex: 'intranet_ip'
      },
      {
        title: '外网IP',
        dataIndex: 'extranet_ip'
      },
      {
        title: '操作系统',
        dataIndex: 'os'
      },
      {
        title: '运营状态',
        dataIndex: 'operation_status'
      }
    ];
  };


  Popup.open({
    title: '导入到机房',
    width: 1500,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <FileUpload
        importApi='/api/cloudboot/v1/devices/stock/imports'
        uploadApi='/api/cloudboot/v1/devices/stock/upload'
        previewApi='/api/cloudboot/v1/devices/stock/imports/previews'
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
