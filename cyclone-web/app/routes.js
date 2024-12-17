// These are the pages you can go to.
// They are all wrapped in the App component, which should contain the navbar etc
// See http://blog.mxstbr.com/2016/01/react-apps-with-pages for more information
// about the code splitting business

import Err404 from "containers/common/err404/get-component";
import Authenticated from "containers/common/authenticated/get-component";
import Homepage from "containers/homepage/get-component";
import Develop from "containers/common/develop/get-component";

import DatabaseIdc from 'containers/database/idc/get-component';
import DatabaseRoom from 'containers/database/room/get-component';
import DatabaseNetwork from 'containers/database/network/get-component';
import DatabaseCabinet from 'containers/database/cabinet/get-component';
import DatabaseUsite from 'containers/database/usite/get-component';
import DatabaseStoreRoom from 'containers/database/store-room/get-component';
import DatabaseStoreRoomDetail from 'containers/database/store-room-detail/get-component';

import NetworkCidr from 'containers/network/cidr/get-component';
import NetworkIps from 'containers/network/ips/get-component';
import NetworkDevice from 'containers/network/device/get-component';

import DeviceList from 'containers/device/list/get-component';
import DeviceOob from 'containers/device/oob/get-component';
import DeviceSpecial from 'containers/device/special/get-component';
import DevicePreDeploy from 'containers/device/pre_deploy/get-component';
import DeviceSetting from 'containers/device/setting/get-component';
import DeviceSettingRule from 'containers/device/setting_rule/get-component';
import DeviceEntry from 'containers/device/entry/get-component';
import DeviceDetail from 'containers/device/detail/get-component';
import DeviceInspectionList from 'containers/device/inspection/list/get-component';
import DeviceInspectionDetail from 'containers/device/inspection/detail/get-component';

import TemplateHardwareList from 'containers/template/hardware/list/get-component';
import TemplateHardwareCreate from 'containers/template/hardware/create/get-component';
import TemplateSystem from 'containers/template/system/get-component';

import AuditLog from 'containers/audit/log/get-component';
import AuditApi from 'containers/audit/api/get-component';

import Approval from 'containers/approval/get-component';
import ApprovalRetirement from 'containers/approval/apply_pages/retirement/get-component';
import ApprovalMove from 'containers/approval/apply_pages/move/get-component';
import ApprovalReInstall from 'containers/approval/apply_pages/reInstall/get-component';
import ApprovalPages from 'containers/approval/apply_pages/pages/get-component';

import OrderList from 'containers/order/list/get-component';
import OrderDeviceCategory from 'containers/order/deviceCategory/get-component';

import TaskList from 'containers/task/list/get-component';
import TaskDetail from 'containers/task/detail/get-component';

import errorLoading from 'common/load-route-error';
import { getAsyncInjectors } from './utils/asyncInjectors';

export default function createRoutes(store) {
  const { injectReducer, injectSagas } = getAsyncInjectors(store);
  const options = {
    injectReducer,
    injectSagas,
    errorLoading
  };
  return [
    {
      path: '/',
      name: '门户',
      getComponent: Homepage(options)
    },
    {
      path: "/test",
      name: "测试",
      getComponent: Develop(options)
    },
    {
      path: "/authenticated",
      name: "权限",
      getComponent: Authenticated()
    },
    {
      path: "database/idc",
      name: "数据中心消息管理",
      getComponent: DatabaseIdc(options)
    },
    {
      path: "database/room",
      name: "机房信息管理",
      getComponent: DatabaseRoom(options)
    },
    {
      path: "database/network",
      name: "网络区域信息管理",
      getComponent: DatabaseNetwork(options)
    },
    {
      path: "database/cabinet",
      name: "机架信息管理",
      getComponent: DatabaseCabinet(options)
    },
    {
      path: "database/usite",
      name: "机位信息管理",
      getComponent: DatabaseUsite(options)
    },
    {
      path: "database/store-room",
      name: "库房信息管理",
      getComponent: DatabaseStoreRoom(options)
    },
    {
      path: "database/store-room/:id",
      name: "库房管理详情",
      getComponent: DatabaseStoreRoomDetail(options)
    },
    {
      path: "network/cidr",
      name: "IP网段管理",
      getComponent: NetworkCidr(options)
    },
    {
      path: "network/ips",
      name: "IP分配管理",
      getComponent: NetworkIps(options)
    },
    {
      path: "network/device",
      name: "网络设备",
      getComponent: NetworkDevice(options)
    },
    {
      path: "device/list",
      name: "物理机列表",
      getComponent: DeviceList(options)
    },
    {
      path: "device/oob",
      name: "带外",
      getComponent: DeviceOob(options)
    },
    {
      path: "device/special",
      name: "特殊设备",
      getComponent: DeviceSpecial(options)
    },
    {
      path: "device/pre_deploy",
      name: "待部署列表",
      getComponent: DevicePreDeploy(options)
    },
    {
      path: "device/detail/:sn",
      name: "设备详情",
      getComponent: DeviceDetail(options)
    },
    {
      path: "device/setting",
      name: "装机列表",
      getComponent: DeviceSetting(options)
    },
    {
      path: "device/setting_rule",
      name: "部署参数规则",
      getComponent: DeviceSettingRule(options)
    },    
    {
      path: "device/entry",
      name: "上架部署",
      getComponent: DeviceEntry(options)
    },
    {
      path: "device/inspection/list",
      name: "硬件巡检",
      getComponent: DeviceInspectionList(options)
    },
    {
      path: "device/inspection/detail/:sn",
      name: "硬件巡检详情",
      getComponent: DeviceInspectionDetail(options)
    },
    {
      path: "template/hardware/list",
      name: "硬件模板",
      getComponent: TemplateHardwareList(options)
    },
    {
      path: "template/hardware/create/:id",
      name: "新建硬件配置模板",
      getComponent: TemplateHardwareCreate(options)
    },
    {
      path: "template/hardware/edit/:id",
      name: "编辑硬件配置模板",
      getComponent: TemplateHardwareCreate(options)
    },
    {
      path: "template/hardware/detail/:id",
      name: "查看硬件配置模板",
      getComponent: TemplateHardwareCreate(options)
    },
    {
      path: "template/system",
      name: "系统模板",
      getComponent: TemplateSystem(options)
    },
    {
      path: "audit/log",
      name: "操作记录",
      getComponent: AuditLog(options)
    },
    {
      path: "audit/api",
      name: "接口调用记录",
      getComponent: AuditApi(options)
    },
    {
      path: "approval",
      name: "审批",
      getComponent: Approval(options)
    },
    {
      path: "approval/pages/:type",
      name: "审批单",
      getComponent: ApprovalPages(options)
    },
    {
      path: "approval/retire",
      name: "物理机退役",
      getComponent: ApprovalRetirement(options)
    },
    {
      path: "approval/move",
      name: "物理机搬迁",
      getComponent: ApprovalMove(options)
    },
    {
      path: "approval/reinstall",
      name: "物理机重装",
      getComponent: ApprovalReInstall(options)
    },
    {
      path: "order/list",
      name: "订单列表",
      getComponent: OrderList(options)
    },
    {
      path: "order/deviceCategory",
      name: "设备类型",
      getComponent: OrderDeviceCategory(options)
    },
    {
      path: "task/list",
      name: "任务列表",
      getComponent: TaskList(options)
    },
    {
      path: "task/detail/:id",
      name: "任务管理详情",
      getComponent: TaskDetail(options)
    },
    {
      path: "*",
      name: "notfound",
      getComponent: Err404(options)
    }
  ];
}
