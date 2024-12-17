export const IDC_USAGE = {
  'production': '生产',
  'disaster_recovery': '容灾',
  'pre_production': '准生产',
  'testing': '测试'
};
export const IDC_STATUS = {
  'under_construction': '建设中',
  'accepted': '已验收',
  'production': '已投产',
  'abolished': '已裁撤'
};
export const IDC_STATUS_COLOR = {
  'under_construction': [ '#3A9EEF', '建设中' ],
  'accepted': [ '#8BD46D', '已验收' ],
  'production': [ '#594DCF', '已投产' ],
  'abolished': [ '#F39821', '已裁撤' ]
};
export const NET_STATUS = {
  'nonproduction': '未投产',
  'production': '已投产',
  'offline': '已下线(回收)'
};
export const NET_STATUS_COLOR = {
  'nonproduction': [ '#F39821', '未投产' ],
  'production': [ '#8BD46D', '已投产' ],
  'offline': [ '#929EA6', '已下线(回收)' ]
};
export const CAB_TYPE = {
  'server': '通用服务器',
  'kvm_server': '虚拟化服务器',
  'network_device': '网络设备',
  'reserved': '预留'
};
export const YES_NO = {
  'yes': '是',
  'no': '否',
  'disabled': '不可用'
};

export const OOB_STATUS = {
  unknown: '未知',
  'yes': '正常',
  'no': '异常'
};
export const OOB_ACCESSIBLE = {
  unknown: '未知',
  yes: '正常',
  no: '异常'
};

export const OOB_STATUS_COLOR = {
  'no': [ '#F39821', '异常' ],
  'yes': [ '#8BD46D', '正常' ],
  'unknown': [ '#929EA6', '未知' ],
  '': [ '#929EA6', '未知' ]
};

export const POWER_STATUS = {
  'power_on': '开电',
  'power_off': '关电'
};
export const CAB_STATUS = {
  'under_construction': '建设中',
  'not_enabled': '未启用',
  'enabled': '已启用',
  'offline': '已下线',
  'locked': '已锁定'
};
export const CAB_STATUS_COLOR = {
  'under_construction': [ '#3A9EEF', '建设中' ],
  'not_enabled': [ '#F39821', '未启用' ],
  'enabled': [ '#8BD46D', '已启用' ],
  'offline': [ '#929EA6', '已下线' ],
  'locked': [ '#5c0011', '已锁定']
};
export const USITE_STATUS = {
  'free': '空闲',
  'pre_occupied': '预占用',
  'used': '已使用',
  'disabled': '不可用'
};
export const USITE_PORT_RATE = {
  'GE': 'GE',
  '10GE': '10GE',
  '25GE': '25GE',
  '40GE': '40GE'
};
export const USITE_STATUS_COLOR = {
  'free': [ '#8BD46D', '空闲' ],
  'pre_occupied': [ '#3A9EEF', '预占用' ],
  'used': [ '#594DCF', '已使用' ],
  'disabled': [ '#929EA6', '不可用' ]
};
export const IP_NETWORK_CATEGORY = {
  'ilo': '服务器ILO',
  'tgw_intranet': '服务器TGW内网',
  'tgw_extranet': '服务器TGW外网',
  'intranet': '服务器普通内网',
  'extranet': '服务器普通外网',
  'v_intranet': '服务器虚拟化内网',
  'v_extranet': '服务器虚拟化外网'
};
export const IP_VERSION = {
  'ipv4': 'IPv4',
  'ipv6': 'IPv6'
};
export const OPERATION_STATUS = {
  'run_with_alarm': '运营中(需告警)',
  'run_without_alarm': '运营中(无需告警)',
  'reinstalling': '重装中',
  'in_store': '库房中',
  'moving': '搬迁中',
  'pre_retire': '待退役',
  'retiring': '退役中',  
  'retired': '已退役',
  'pre_deploy': '待部署',
  'on_shelve': '已上架',
  'recycling': '回收中',
  'maintaining': '维护中',
  'pre_move': '待搬迁'
};
export const OPERATION_STATUS_COLOR = {
  'run_with_alarm': [ '#F39821', '运营中(需告警)' ],
  'run_without_alarm': [ '#3A9EEF', '运营中(无需告警)' ],
  'reinstalling': [ '#ff3700', '重装中' ],
  'moving': [ '#FC5CAD', '搬迁中' ],
  'in_store': [ '#FC5CAD50', '库房中' ],
  'pre_retire': [ '#BB67D4', '待退役' ],
  'retiring': [ '#BB67D4', '退役中' ], 
  'retired': [ '#929EA6', '已退役' ],
  'pre_deploy': [ '#594DCF', '待部署' ],
  'on_shelve': [ '#8BD46D', '已上架' ],
  'recycling': [ '#7D89E6', '回收中' ],
  'maintaining': [ '#F39821', '维护中' ],
  'pre_move': [ '#594DCF70', '待搬迁' ]
};
export const NETWORK_DEVICE_TYPE = {
  'switch': '交换机'
};
export const NETWORK_DEVICE_STATUS = {
  '运营中': '运营中',
  '待启用': '待启用',
};
export const NETWORK_IPS_CATEGORY = {
  'pxe': 'PXEIP',
  'business': '业务IP'
};
export const NETWORK_IPS_SCOPE = {
  "intranet": '内网',
  "extranet": '外网'
};

export const INSTALL_TYPE = {
  'pxe': 'pxe安装',
  'image': '镜像安装'
};
export const BOOT_MODE = {
  uefi: 'UEFI',
  legacy_bios: 'BIOS'
};
export const OS_LIFECYCLE = {
	"testing":"Testing",
	"active_default":"Active(Default)",
	"active":"Active",
	"containment":"Containment",
	"end_of_life":"EOL"
};
export const DEVICE_INSTALL_STATUS = {
  'pre_install': '等待部署',
  'installing': '正在部署',
  'failure': '部署失败',
  'success': '部署成功'
};
export const DEVICE_INSTALL_STATUS_COLOR = {
  success: [ '#6bc646', '部署成功' ],
  failure: [ '#ff3700', '部署失败' ],
  pre_install: [ '#00c3ff', '等待部署' ],
  installing: [ '#8E44AD', '正在部署' ]
};
export const PRIVILEGE_LEVEL = {
  '0': '0(noaccess)',
  '1': '1(callback)',
  '2': '2(user)',
  '3': '3(operator)',
  '4': '4(administrator)',
  '5': '5(oem)'
};

export const NICSIDE = {
  'inside': '内置',
  'outside': '外置'
};
export const RUNNING_STTAUS = {
  'running': '正在巡检',
  'done': '巡检完成'
};
export const INSPECTION_RESULT = {
  nominal: '正常',
  critical: '异常',
  unknown: '未知',
  warning: '警告'
};
export const INSPECTION_RESULT_COLOR = {
  nominal: [ '#6bc646', '正常' ],
  critical: [ '#ff3700', '异常' ],
  unknown: [ '#c4cdd7', '未知' ],
  warning: [ '#ff9e0c', '警告' ]
};

export const APPROVAL_STATUS_COLOR = {
  completed: [ '#6bc646', '已完成' ],
  revoked: [ '#c4cdd7', '已撤销' ],
  approval: [ '#ff9e0c', '进行中' ],
  failure: [ '#ff3700', '失败' ]
};

export const APPROVAL_STATUS = {
  completed: '已完成',
  revoked: '已撤销',
  approval: '进行中',
  failure: '失败'
};

export const APPROVAL_TYPE = {
  cabinet_power_off: '机架关电',
  cabinet_offline: '机架下线',
  device_os_reinstallation: '物理机操作系统重装',
  device_migration: '物理机搬迁',
  device_retirement: '物理机退役(报废)',
  idc_abolish: '数据中心裁撤',
  server_room_abolish: '机房裁撤',
  network_area_offline: '网络区域下线',
  ip_unassign: 'IP回收',
  device_restart: '物理机重启',
  device_power_off: '物理机关电',
  device_recycle_pre_retire: '回收退役',
  device_recycle_pre_move: '回收搬迁',
  device_recycle_reinstall: '回收重装'
};
export const APPROVAL_ACTION = {
  agree: '通过',
  reject: '不通过'
};
export const TIME_FORMAT = 'YYYY-MM-DD HH:mm:ss';

export const DEPLOY_STATUS = {
  'pre_deploy': '等待部署',
  'pre_online': '等待上线',
  'online': '已上线',
  'offline': '已下线'
};

export const PAGE_TYPE = {
  edit: '编辑配置',
  create: '新增配置',
  detail: '查看配置'
};
export const API_STATUS = {
  'success': '成功',
  'failure': '失败'
};
export const API_STATUS_COLOR = {
  success: [ '#6bc646', '成功' ],
  failure: [ '#ff3700', '失败' ]
};

export const ORDER_STATUS = {
  purchasing: '采购中',
  partly_arrived: '部分到货',
  all_arrived: '全部到货',
  canceled: '已取消',
  finished: '已完成'
};
export const BUILTIN = {
  yes: '是',
  no: '否'
};
export const ARCH = {
  x86_64: 'x86_64',
  aarch64: 'aarch64',
  ppc64: 'ppc64'
};
export const TASK_STATUS = {
  running: '运行中',
  paused: '已暂停',
  stoped: '已结束',
  deleted: '已删除'
};

export const TASK_STATUS_COLOR = {
  stoped: [ '#6bc646', '已结束' ],
  deleted: [ '#ff3700', '已删除' ],
  paused: [ '#ff9e0c', '已暂停' ],
  running: [ '#00c3ff', '运行中' ]
};
export const TASK_CATEGORY = {
  inspection: '硬件巡检任务',
  installation_timeout: '装机超时检查任务',
  release_ip: '释放IP任务',
  auto_deploy: '自动部署任务',
  mailsend: '定时邮件任务',
  update_device_lifecycle: '更新设备维保状态',  
  xmdb_rabbitmq: 'xmdb信息同步'
};

export const TASK_RATE = {
  immediately: '立刻',
  fixed_rate: '定时'
};

export const DEVICE_SETTING_RULE_CATEGORY = {
  os: '操作系统',
  raid: '阵列结构',
  network: '网络配置'
};

export const DEVICE_MAINTENANCE_SERVICE_STATUS = {
  under_warranty: '在保',
  out_of_warranty: '过保',
  inactive: '未激活'
};

export const formItemLayout = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 5 }
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 18 }
  }
};
export const formItemLayout_page = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 2 }
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 22 }
  }
};
export const tailFormItemLayout = {
  wrapperCol: {
    xs: {
      span: 24,
      offset: 0
    },
    sm: {
      span: 10,
      offset: 13
    }
  }
};
export const tableFormItemLayout = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 1 }
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 23 }
  }
};
export const tailFormItemLayout_page = {
  wrapperCol: {
    xs: {
      span: 24,
      offset: 0
    },
    sm: {
      span: 10,
      offset: 2
    }
  }
};

export function getSearchList(data) {
  return Object.keys(data).map(key => {
    return {
      value: key, label: data[key]
    };
  });
}


export function getEnumLabel(enums = {}, key) {
  return enums[key] || key;
}

export function getEnumOptions(enums) {
  const keys = Object.keys(enums);
  const options = keys.map(key => {
    return {
      label: enums[key],
      value: key
    };
  });
  options.unshift({
    label: '全部',
    value: ''
  });
  return options;
}
