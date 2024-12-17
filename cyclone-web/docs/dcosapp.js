[{
  attributes: {
    0 ipAddress: "IP地址",
    0 macAddress: "mac地址",
    0 networkId: "网络ID",
    appInfoId: null,
    ciId: null,
    id: null
  },
  ci: "AppDevNicVO.json",
  name: "应用部署实例网卡"
}, {
  attributes: {
    0 hostName: "主机名",
    appInfoId: null,
    ciId: null,
    id: null
  },
  ci: "AppDevVO.json",
  name: "应用部署实例"
}, {
  attributes: {
    0 ipAddress: "服务IP(VIP)",
    appInfoId: null,
    ciId: null,
    id: null
  },
  ci: "AppDpUnitLbVO.json",
  name: "应用部署单元负载均衡配置"
}, {
  attributes: {
    0 addressType: "服务地址类型",
    appInfoId: null,
    ciId: null,
    id: null
  },
  ci: "AppDpUnitNetServiceVO.json",
  name: "应用部署单元网络服务"
}, {
  attributes: {
    0 rel_DcosAppLogicUnit_Id: "应用逻辑单元",
    1 name: "名称",
    10 devNum: "部署实例数量",
    2 rel_DcosIdc_Id: "数据中心",
    3 rel_DcosAppDpUnitType_Id: "部署单元类型",
    4 rel_DcosComputeResPools_Id: "计算资源池",
    5 rel_DcosComputeSysImgs_Id: "计算资源系统镜像",
    6 rel_DcosComputeFlavors_Id: "计算规格",
    7 rel_DcosNetArea_Id: "网络分区",
    8 rel_DcosNetSecureArea_Id: "安全分区",
    appInfoId: null,
    ciId: null,
    id: null,
    rel_DcosAppDpUnitType_Name: null,
    rel_DcosAppLogicUnit_Name: null,
    rel_DcosComputeFlavors_Name: null,
    rel_DcosComputeResPools_Name: null,
    rel_DcosComputeSysImgs_Name: null,
    rel_DcosIdc_Name: null,
    rel_DcosNetArea_Name: null,
    rel_DcosNetSecureArea_Name: null
  },
  ci: "AppDpUnitVO.json",
  name: "应用部署单元"
}, {
  attributes: {
    0 ownerId: "应用负责人ID",
    0 rel_DcosAppInfoStatus_Id: "应用状态",
    1 name: "应用名称",
    10 remark: "备注",
    2 code: "应用编码",
    3 rel_DcosTenant_Id: "租户",
    4 rel_DcosEnvironment_Id: "环境",
    5 rel_DcosAppInfoType_Id: "应用类型",
    6 rel_DcosHaLevel_Id: "高可用等级",
    7 rel_DcosContinuityLevel_Id: "连续性等级",
    8 rel_DcosProtectLevel_Id: "等级保护等级",
    9 ownerName: "应用负责人",
    appInfoId: null,
    ciId: null,
    id: null,
    rel_DcosAppInfoStatus_Name: null,
    rel_DcosAppInfoType_Name: null,
    rel_DcosContinuityLevel_Name: null,
    rel_DcosEnvironment_Name: null,
    rel_DcosHaLevel_Name: null,
    rel_DcosProtectLevel_Name: null,
    rel_DcosTenant_Name: null
  },
  ci: "AppInfoVO.json",
  name: "应用信息"
}, {
  attributes: {
    0 rel_DcosAppInfo_Id: "应用信息",
    1 name: "名称",
    2 code: "编码",
    3 rel_DcosAppLogicUnitType_Id: "应用逻辑单元类型",
    appInfoId: null,
    ciId: null,
    id: null,
    rel_DcosAppInfo_Name: null,
    rel_DcosAppLogicUnitType_Name: null
  },
  ci: "AppLogicUnitVO.json",
  name: "应用逻辑单元"
}, {
  attributes: {
    0 rel_DcosAppService_Id: "应用服务",
    1 rel_DcosLbModel_Id: "负载均衡模式",
    2 rel_DcosLbAlgorithm_Id: "负载均衡算法",
    3 rel_DcosLbSessionPersistenceType_Id: "会话保持类型",
    4 sessionPersistenceTime: "会话保持时间",
    5 rel_DcosLbHealthCheckModel_Id: "康检查模式",
    6 monitorCheckUrl: "健康检查地址",
    7 rel_DcosLbRecordSourceIp_Id: "是否记录源IP",
    8 rel_DcosLbSslOffload_Id: "否需网络执行ssl卸载",
    appInfoId: null,
    ciId: null,
    id: null,
    rel_DcosAppService_Name: null,
    rel_DcosLbAlgorithm_Name: null,
    rel_DcosLbHealthCheckModel_Name: null,
    rel_DcosLbModel_Name: null,
    rel_DcosLbRecordSourceIp_Name: null,
    rel_DcosLbSessionPersistenceType_Name: null,
    rel_DcosLbSslOffload_Name: null
  },
  ci: "AppServiceLbVO.json",
  name: "应用服务负载均衡配置"
}, {
  attributes: {
    0 rel_DcosAppLogicUnit_Id: "应用逻辑单元",
    1 name: "服务名称",
    2 domainName: "服务域名",
    3 port: "服务端口",
    4 rel_DcosAppServiceModel_Id: "服务模式",
    5 rel_DcosAppServiceProtocolType_Id: "协议类型",
    6 rel_DcosIsPersistentConnection_Id: "长短连接",
    7 timeOut: "长连接超时时间(h)",
    8 rel_DcosDomainIsGSLB_Id: "域名是否全局负载",
    9 rel_DcosAppServiceUseLB_Id: "是否使用负载均衡",
    appInfoId: null,
    ciId: null,
    id: null,
    rel_DcosAppLogicUnit_Name: null,
    rel_DcosAppServiceModel_Name: null,
    rel_DcosAppServiceProtocolType_Name: null,
    rel_DcosAppServiceUseLB_Name: null,
    rel_DcosDomainIsGSLB_Name: null,
    rel_DcosIsPersistentConnection_Name: null
  },
  ci: "AppServiceVO.json",
  name: "应用服务"
}, {
  attributes: {
    0 rel_DcosAppLogicUnit_Id: "应用逻辑单元",
    1 rel_DcosAppStorageType_Id: "存储类型",
    2 mountPoint: "挂载点",
    3 size: "容量",
    4 rel_DcosAppStorageLevel_Id: "存储规格",
    5 rel_DcosAppStorageFileSystemType_Id: "文件系统类型 ",
    6 deviceName: "设备名",
    7 ownerUser: "属主",
    8 ownerGroup: "属组",
    9 remark: "备注",
    appInfoId: null,
    ciId: null,
    id: null,
    rel_DcosAppLogicUnit_Name: null,
    rel_DcosAppStorageFileSystemType_Name: null,
    rel_DcosAppStorageLevel_Name: null,
    rel_DcosAppStorageType_Name: null
  },
  ci: "AppStorageVO.json",
  name: "应用存储卷"
}, {
  attributes: {
    0 rel_DcosAppLogicUnit_Id: "应用逻辑单元",
    1 gid: "GID",
    2 name: "用户组名称",
    appInfoId: null,
    ciId: null,
    id: null,
    rel_DcosAppLogicUnit_Name: null
  },
  ci: "AppSysUserGroupVO.json",
  name: "操作系统用户组"
}, {
  attributes: {
    0 rel_DcosAppLogicUnit_Id: "应用逻辑单元",
    1 gid: "GID",
    2 uid: "UID",
    3 name: "用户名",
    4 home: "Home目录",
    appInfoId: null,
    ciId: null,
    id: null,
    rel_DcosAppLogicUnit_Name: null
  },
  ci: "AppSysUserVO.json",
  name: "操作系统用户"
}]