import { getIdcColumns } from './columns/getIdcColumns';
import { getRoomColumns } from './columns/getRoomColumns';
import { getNetworkColumns } from './columns/getNetworkColumns';
import { getCabColumns } from './columns/getCabColumns';
import { getIpColumns } from './columns/getIpColumns';
import { getColumns } from 'containers/device/common/colums';


export const TYPE = {
  cabinet_power_off: {
    name: '机架关电',
    category: '机架',
    tableDataUrl: '/api/cloudboot/v1/server-cabinets',
    tableColumn: getCabColumns,
    tableQuery: { is_powered: 'yes' },
    searchKey: [{ key: 'idc_id', label: '数据中心', type: 'select', async: 'idc' },
    { key: 'server_room_name', label: '机房' },
    { key: 'number', label: '机架编号' },
    { key: 'network_area_name', label: '网络区域' }],
    modalKey: 'ids',
    submitUrl: '/api/cloudboot/v1/approvals/server-cabinets/poweroffs'
  },
  cabinet_offline: {
    name: '机架下线',
    category: '机架',
    searchKey: [{ key: 'number', label: '机架编号' }, { key: 'idc_id', label: '数据中心', type: 'select', async: 'idc' },
    { key: 'server_room_name', label: '机房' },
    { key: 'network_area_name', label: '网络区域' }],
    tableDataUrl: '/api/cloudboot/v1/server-cabinets',
    tableColumn: getCabColumns,
    tableQuery: { is_powered: 'no' },
    modalKey: 'ids',
    submitUrl: '/api/cloudboot/v1/approvals/server-cabinets/offlines'
  },
  // device_os_reinstallation: '物理机操作系统重装',
  // device_migration: '物理机搬迁',
  // device_retirement: '物理机退役(报废)',
  // device_recycle_reinstall: '回收重装',
  idc_abolish: {
    name: '数据中心裁撤',
    category: '数据中心',
    searchKey: [{ key: 'name', label: '名称' }],
    tableDataUrl: '/api/cloudboot/v1/idcs',
    tableColumn: getIdcColumns,
    tableQuery: { status: 'under_construction,accepted,production' },
    submitUrl: '/api/cloudboot/v1/approvals/idc/abolish',
    modalKey: 'ids'
  },
  server_room_abolish: {
    name: '机房裁撤',
    category: '机房',
    searchKey: [{ key: 'name', label: '名称' }],
    tableDataUrl: '/api/cloudboot/v1/server-rooms',
    tableColumn: getRoomColumns,
    tableQuery: {},
    submitUrl: '/api/cloudboot/v1/approvals/server-room/abolish',
    modalKey: 'ids'
  },
  network_area_offline: {
    name: '网络区域下线',
    category: '网络区域',
    searchKey: [{ key: 'name', label: '名称' }],
    tableDataUrl: '/api/cloudboot/v1/network-areas',
    tableColumn: getNetworkColumns,
    tableQuery: {},
    submitUrl: '/api/cloudboot/v1/approvals/network-area/offline',
    modalKey: 'ids'
  },
  ip_unassign: {
    name: 'IP回收',
    category: 'IP',
    searchKey: [{ key: 'ip', label: 'IP' }, { key: 'fixed_asset_number', label: '固资编号' }],
    tableDataUrl: '/api/cloudboot/v1/ips',
    tableColumn: getIpColumns,
    tableQuery: { category: 'business', is_used: 'disabled,yes' },
    submitUrl: '/api/cloudboot/v1/approvals/ip/unassign',
    modalKey: 'ids'
  },
  device_power_off: {
    name: '物理机关电',
    category: '物理机',
    searchKey: [{ key: 'sn', label: '序列号' }, { key: 'fixed_asset_number', label: '固资编号' }, { key: 'intranet_ip', label: '内网IP' }],
    tableDataUrl: '/api/cloudboot/v1/devices',
    tableColumn: getColumns,
    tableQuery: {},
    submitUrl: '/api/cloudboot/v1/approvals/device/poweroff',
    modalKey: 'sns'
  },
  device_restart: {
    name: '物理机重启',
    category: '物理机',
    searchKey: [{ key: 'sn', label: '序列号' }, { key: 'fixed_asset_number', label: '固资编号' }, { key: 'intranet_ip', label: '内网IP' }],
    tableDataUrl: '/api/cloudboot/v1/devices',
    tableColumn: getColumns,
    tableQuery: {},
    submitUrl: '/api/cloudboot/v1/approvals/device/restart',
    modalKey: 'sns'
  },
  device_recycle_pre_retire: {
    name: '回收待退役',
    category: '物理机',
    limit: 50,
    searchKey: [{ key: 'sn', label: '序列号' }, { key: 'fixed_asset_number', label: '固资编号' }, { key: 'intranet_ip', label: '内网IP' }],
    tableDataUrl: '/api/cloudboot/v1/devices',
    tableColumn: getColumns,
    tableQuery: { operation_status: 'recycling' },
    submitUrl: '/api/cloudboot/v1/approvals/devices/recycle',
    submitData: { approval_type: 'device_recycle_pre_retire' },
    modalKey: 'sns'
  },
  device_recycle_pre_move: {
    name: '回收待搬迁',
    category: '物理机',
    limit: 50,
    searchKey: [{ key: 'sn', label: '序列号' }, { key: 'fixed_asset_number', label: '固资编号' }, { key: 'intranet_ip', label: '内网IP' }],
    tableDataUrl: '/api/cloudboot/v1/devices',
    tableColumn: getColumns,
    tableQuery: { operation_status: 'recycling' },
    submitUrl: '/api/cloudboot/v1/approvals/devices/recycle',
    submitData: { approval_type: 'device_recycle_pre_move' },
    modalKey: 'sns'
  }
};

