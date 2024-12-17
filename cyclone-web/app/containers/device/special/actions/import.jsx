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
        title: '厂商',
        dataIndex: 'vendor'
      },
      {
        title: '型号',
        dataIndex: 'model'
      },
      {
        title: '设备类型',
        dataIndex: 'category'
      },
      {
        title: '负责人',
        dataIndex: 'owner'
      },
      {
        title: '维保服务起始日期',
        dataIndex: 'maintenance_service_date_begin'
      },
      {
        title: '保修期（月数）',
        dataIndex: 'maintenance_months'
      },
      {
        title: '机房管理单元名称',
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
        title: '硬件说明',
        dataIndex: 'hardware_remark'
      },
      {
        title: '操作系统名称',
        dataIndex: 'os_release_name'
      },      
      {
        title: '是否分配内网IPv4',
        dataIndex: 'need_intranet_ip'
      },
      {
        title: '是否分配外网IPv4',
        dataIndex: 'need_extranet_ip'
      },
      {
        title: '关联订单编号',
        dataIndex: 'order_number'
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
        importApi='/api/cloudboot/v1/special-devices/imports'
        uploadApi='/api/cloudboot/v1/special-devices/upload'
        previewApi='/api/cloudboot/v1/special-devices/imports/previews'
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
