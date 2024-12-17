import React from 'react';
import Popup from 'components/popup';
import FileUpload from 'components/upload';

export default function action(options) {
  const getColumns = () => {
    return [
      {
        title: '序列号',
        dataIndex: 'sn'
      },
      {
        title: '厂商',
        dataIndex: 'vendor'
      },
      {
        title: '设备型号',
        dataIndex: 'model'
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
        title: '库房管理单元',
        dataIndex: 'store_room_name'
      },
      {
        title: '虚拟货架编号',
        dataIndex: 'virtual_cabinet_number'
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
        title: '带外',
        dataIndex: 'oob_init'
      },
      {
        title: '关联订单号',
        dataIndex: 'order_number'
      },
      {
        title: '负责人',
        dataIndex: 'owner'
      },
      {
        title: '是否租赁',
        dataIndex: 'is_rental'
      },
      {
        title: '资产归属',
        dataIndex: 'asset_belongs'
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
        title: '维保服务供应商',
        dataIndex: 'maintenance_service_provider'
      },
      {
        title: '维保服务内容',
        dataIndex: 'maintenance_service'
      },
      {
        title: '物流服务内容',
        dataIndex: 'logistics_service'
      },
      {
        title: '启用时间',
        dataIndex: 'started_at'
      },
      {
        title: '上架时间',
        dataIndex: 'onshelve_at'
      },           
    ];
  };


  Popup.open({
    title: '导入到库房',
    width: 2000,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <FileUpload
        importApi='/api/cloudboot/v1/devices/store/imports'
        uploadApi='/api/cloudboot/v1/devices/store/upload'
        previewApi='/api/cloudboot/v1/devices/store/imports/previews'
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
