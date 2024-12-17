package cloudbootserver

import (
	"github.com/go-chi/chi"

	"idcos.io/cloudboot/server/cloudbootserver/route"
)

func registerRoutes(mux *chi.Mux) {
	// 数据中心
	mux.Post("/api/cloudboot/v1/idcs", route.SaveIDC)
	mux.Put("/api/cloudboot/v1/idcs/status", route.UpdateIDCStatus)
	mux.Delete("/api/cloudboot/v1/idcs/{id}", route.RemoveIDCByID)
	mux.Put("/api/cloudboot/v1/idcs/{id}", route.UpdateIDC)
	mux.Get("/api/cloudboot/v1/idcs", route.GetIDCPage)
	mux.Get("/api/cloudboot/v1/idcs/{id}", route.GetIDCByID)

	// 机房
	mux.Post("/api/cloudboot/v1/server-rooms", route.SaveServerRoom)
	mux.Put("/api/cloudboot/v1/server-rooms/{id}", route.UpdateServerRoom)
	mux.Put("/api/cloudboot/v1/server-rooms/status", route.UpdateServerRoomStatus)
	mux.Delete("/api/cloudboot/v1/server-rooms/{id}", route.RemoveServerRoomByID)
	mux.Get("/api/cloudboot/v1/server-rooms", route.GetServerRoomPage)
	mux.Get("/api/cloudboot/v1/server-rooms/{id}", route.GetServerRoomByID)
	mux.Post("/api/cloudboot/v1/server-rooms/upload", route.UploadServerRoom)
	mux.Post("/api/cloudboot/v1/server-rooms/imports/previews", route.ImportServerRoomPriview)
	mux.Post("/api/cloudboot/v1/server-rooms/imports", route.ImportServerRoom)

	// 库房
	mux.Post("/api/cloudboot/v1/store-room", route.SaveStoreRoom)
	mux.Delete("/api/cloudboot/v1/store-room/{id}", route.RemoveStoreRoom)
	mux.Put("/api/cloudboot/v1/store-room", route.SaveStoreRoom)
	mux.Get("/api/cloudboot/v1/store-rooms", route.GetStoreRooms)
	mux.Get("/api/cloudboot/v1/store-room/{id}", route.GetStoreRoom)
	mux.Post("/api/cloudboot/v1/store-room/upload", route.UploadStoreRoom)
	mux.Post("/api/cloudboot/v1/store-room/imports/previews", route.ImportStoreRoomPriview)
	mux.Post("/api/cloudboot/v1/store-room/imports", route.ImportStoreRoom)

	//虚拟货架
	mux.Post("/api/cloudboot/v1/virtual-cabinet", route.SaveVirtualCabinet)
	mux.Delete("/api/cloudboot/v1/virtual-cabinet/{id}", route.RemoveVirtualCabinet)
	mux.Get("/api/cloudboot/v1/virtual-cabinets", route.GetVirtualCabinets)

	// 网络区域
	mux.Post("/api/cloudboot/v1/network-areas", route.SaveNetworkArea)
	mux.Put("/api/cloudboot/v1/network-areas/{id}", route.UpdateNetworkArea)
	mux.Put("/api/cloudboot/v1/network-areas/status", route.UpdateNetworkAreasStatus)
	mux.Delete("/api/cloudboot/v1/network-areas/{id}", route.RemoveNetworkAreaByID)
	mux.Get("/api/cloudboot/v1/network-areas", route.GetNetworkAreaPage)
	mux.Get("/api/cloudboot/v1/network-areas/{id}", route.GetNetworkAreaByID)
	mux.Post("/api/cloudboot/v1/network-areas/upload", route.UploadNetworkArea)
	mux.Post("/api/cloudboot/v1/network-areas/imports/previews", route.ImportNetworkAreaPriview)
	mux.Post("/api/cloudboot/v1/network-areas/imports", route.ImportNetworkArea)

	// 机架(柜)
	mux.Post("/api/cloudboot/v1/server-cabinets", route.SaveServerCabinet)
	mux.Put("/api/cloudboot/v1/server-cabinets/{id}", route.UpdateServerCabinet)
	mux.Put("/api/cloudboot/v1/server-cabinets/type", route.BatchUpdateServerCabinetType)
	mux.Put("/api/cloudboot/v1/server-cabinets/remark", route.BatchUpdateServerCabinetRemark)
	mux.Delete("/api/cloudboot/v1/server-cabinets/{id}", route.RemoveServerCabinetByID)
	mux.Get("/api/cloudboot/v1/server-cabinets", route.GetServerCabinetPage)
	mux.Get("/api/cloudboot/v1/server-cabinets/{id}", route.GetServerCabinetByID)
	mux.Post("/api/cloudboot/v1/server-cabinets/power", route.PowerOnServerCabinetByID)
	mux.Delete("/api/cloudboot/v1/server-cabinets/{id}/power", route.PowerOffServerCabinetByID)
	mux.Post("/api/cloudboot/v1/server-cabinets/upload", route.UploadServerCabinet)
	mux.Post("/api/cloudboot/v1/server-cabinets/imports/previews", route.ImportServerCabinetPriview)
	mux.Post("/api/cloudboot/v1/server-cabinets/imports", route.ImportServerCabinet)
	mux.Put("/api/cloudboot/v1/server-cabinets/status", route.UpdateServerCabinetStatus) //将这个接口拆成多个独立的接口
	mux.Put("/api/cloudboot/v1/server-cabinets/status/accept", route.AcceptServerCabinet)
	mux.Put("/api/cloudboot/v1/server-cabinets/status/enabled", route.EnableServerCabinet)
	//mux.Put("/api/cloudboot/v1/server-cabinets/status/offline", route.OfflineServerCabinet)
	mux.Put("/api/cloudboot/v1/server-cabinets/status/reconstruct", route.ReconstructServerCabinet)

	// 机位(U位)
	mux.Post("/api/cloudboot/v1/server-usites", route.SaveServerUSite)
	mux.Put("/api/cloudboot/v1/server-usites/{id}", route.UpdateServerUSite)
	mux.Delete("/api/cloudboot/v1/server-usites/{id}", route.DeleteServerUSite)
	mux.Get("/api/cloudboot/v1/server-usites", route.GetServerUSitePage)
	mux.Get("/api/cloudboot/v1/server-usites/{id}", route.GetServerUSiteByID)
	mux.Get("/api/cloudboot/v1/server-usites/tree", route.GetUsiteTree)
	mux.Delete("/api/cloudboot/v1/server-usites/{id}/ports", route.DeleteServerUSitePort)
	mux.Put("/api/cloudboot/v1/server-usites/status", route.BatchUpdateServerUSitesStatus)
	mux.Put("/api/cloudboot/v1/server-usites/remark", route.BatchUpdateServerUSitesRemark)
	mux.Put("/api/cloudboot/v1/server-usites/usitestatus", route.BatchUpdateServerUSitesStatusByCond)
	mux.Post("/api/cloudboot/v1/server-usites/upload", route.UploadServerUSite)
	mux.Post("/api/cloudboot/v1/server-usites/imports/previews", route.ImportServerUSitePriview)
	mux.Post("/api/cloudboot/v1/server-usites/imports", route.ImportServerUSite)
	mux.Post("/api/cloudboot/v1/server-usites/ports/upload", route.UploadServerUSitePort)
	mux.Post("/api/cloudboot/v1/server-usites/ports/imports/previews", route.ImportServerUSitePortsPriview)
	mux.Post("/api/cloudboot/v1/server-usites/ports/imports", route.ImportServerUSitePort)
	mux.Get("/api/cloudboot/v1/physical-areas", route.GetPhysicalAreas)
	// 物理机
	mux.Post("/api/cloudboot/v1/devices/collections", route.SaveCollectedDevice)
	mux.Get("/api/cloudboot/v1/devices/{sn}/collections", route.GetDeviceBySN)
	mux.Get("/api/cloudboot/v1/devices/{sn}/lifecycle", route.GetDeviceLifecycleBySN)
	mux.Put("/api/cloudboot/v1/devices/{sn}/lifecycle", route.UpdateDeviceLifecycleBySN)
	mux.Put("/api/cloudboot/v1/devices/lifecycles", route.BatchUpdateDeviceLifecycleBySN)
	mux.Post("/api/cloudboot/v1/devices", route.SaveNewDevices) // API新增设备
	mux.Get("/api/cloudboot/v1/devices", route.GetDevicePage)
	mux.Get("/api/cloudboot/v1/devices/tors", route.GetDevicePageByTor)
	mux.Post("/api/cloudboot/v1/devices/upload", route.UploadDevices) //导入到机架
	mux.Post("/api/cloudboot/v1/devices/imports/previews", route.ImportDevicesPreview)
	mux.Post("/api/cloudboot/v1/devices/imports", route.ImportDevices)
	mux.Post("/api/cloudboot/v1/devices/store/upload", route.UploadDevices2Store) //导入到库房
	mux.Post("/api/cloudboot/v1/devices/store/imports/previews", route.ImportDevices2StorePreview)
	mux.Post("/api/cloudboot/v1/devices/store/imports", route.ImportDevices2Store)
	mux.Post("/api/cloudboot/v1/devices/stock/upload", route.UploadStockDevices) //存量设备导入
	mux.Post("/api/cloudboot/v1/devices/stock/imports/previews", route.ImportStockDevicesPreview)
	mux.Post("/api/cloudboot/v1/devices/stock/imports", route.ImportStockDevices)
	mux.Get("/api/cloudboot/v1/devices/{sn}/combined", route.GetCombinedDeviceBySN)
	mux.Put("/api/cloudboot/v1/device", route.UpdateDevice)
	mux.Put("/api/cloudboot/v1/device/operation/status", route.UpdateDeviceOperationStatus)
	mux.Put("/api/cloudboot/v1/devices", route.UpdateDevices) //批量修改
	mux.Delete("/api/cloudboot/v1/devices", route.DeleteDevices)
	mux.Get("/api/cloudboot/v1/devices/export", route.ExportCombinedDevices)
	mux.Get("/api/cloudboot/v1/devices/query-params/{param_name}", route.GetDeviceQuerys)
	mux.Post("/api/cloudboot/v1/devices/move", route.BatchMoveDevices) //设备搬迁
	mux.Post("/api/cloudboot/v1/devices/retire", route.BatchRetireDevices) //设备退役


	//特殊设备
	mux.Post("/api/cloudboot/v1/special-device", route.SaveSpecialDevice)
	mux.Post("/api/cloudboot/v1/special-devices/upload", route.UploadSpecialDevices)
	mux.Post("/api/cloudboot/v1/special-devices/imports/previews", route.ImportSpecialDevicesPreview)
	mux.Post("/api/cloudboot/v1/special-devices/imports", route.ImportSpecialDevices)

	// 网段
	mux.Get("/api/cloudboot/v1/ip-networks", route.GetIPNetworkPage)
	mux.Get("/api/cloudboot/v1/ip-networks/{id}", route.GetIPNetworkByID)
	mux.Post("/api/cloudboot/v1/ip-networks", route.SaveIPNetwork)
	mux.Put("/api/cloudboot/v1/ip-networks/{id}", route.UpdateIPNetwork)
	mux.Delete("/api/cloudboot/v1/ip-networks/{id}", route.RemoveIPNetworkByID)
	mux.Delete("/api/cloudboot/v1/ip-networks", route.RemoveIPNetworks)
	mux.Get("/api/cloudboot/v1/ips", route.GetIPSPage)
	mux.Get("/api/cloudboot/v1/ips/export", route.ExportIP)      //导出IP信息
	mux.Put("/api/cloudboot/v1/ips/assigns", route.AssignIP)
	mux.Put("/api/cloudboot/v1/ips/assignsv4", route.AssignIPv4)
	mux.Put("/api/cloudboot/v1/ips/assignsv6", route.AssignIPv6)
	mux.Put("/api/cloudboot/v1/ips/status/disable", route.DisableIP)
	//mux.Put("/api/cloudboot/v1/ips/unassigns", route.UnassignIP)
	mux.Post("/api/cloudboot/v1/ip-networks/upload", route.UploadIPNetworks)
	mux.Post("/api/cloudboot/v1/ip-networks/imports/previews", route.ImportIPNetworksPreview)
	mux.Post("/api/cloudboot/v1/ip-networks/imports", route.ImportIPNetworks)

	// 网络设备
	mux.Get("/api/cloudboot/v1/network/devices", route.GetNetworkDevicePage)
	mux.Get("/api/cloudboot/v1/network/devices/{id}", route.GetNetworkDeviceByID)
	mux.Delete("/api/cloudboot/v1/network/devices/{id}", route.DeleteNetworkDeviceByID)
	mux.Delete("/api/cloudboot/v1/network/devices", route.RemoveNetworkDevices)
	mux.Post("/api/cloudboot/v1/network/devices", route.SaveNetworkDevice)
	mux.Post("/api/cloudboot/v1/network/devices/upload", route.UploadNetworkDevices)
	mux.Post("/api/cloudboot/v1/network/devices/imports/previews", route.ImportNetworkDevicesPreview)
	mux.Post("/api/cloudboot/v1/network/devices/imports", route.ImportNetworkDevices)

	// 设备装机过程/状态
	mux.Get("/api/cloudboot/v1/devices/{sn}/is-in-install-list", route.IsInInstallList)
	mux.Post("/api/cloudboot/v1/devices/{sn}/installations/progress", route.ReportInstallProgress)
	mux.Get("/api/cloudboot/v1/devices/{sn}/installations/status", route.GetInstallationStatus)

	// 装机参数
	mux.Post("/api/cloudboot/v1/devices/settings", route.SaveDeviceSettings)
	mux.Put("/api/cloudboot/v1/devices/settings", route.UpdateDeviceSetting)
	mux.Get("/api/cloudboot/v1/devices/settings", route.GetDeviceSettingPage)
	mux.Get("/api/cloudboot/v1/devices/{sn}/settings", route.GetDeviceSettingBySN)
	mux.Get("/api/cloudboot/v1/devices/{sn}/settings/networks", route.GetNetworkSettingBySN)
	mux.Get("/api/cloudboot/v1/devices/{sn}/settings/os-users", route.GetOSUserSettingsBySN)
	mux.Get("/api/cloudboot/v1/devices/{sn}/settings/hardwares", route.GetHardwareSettingBySN)
	mux.Get("/api/cloudboot/v1/devices/{sn}/settings/hardwareinfo", route.GetHardwareInfoBySN)
	mux.Get("/api/cloudboot/v1/devices/{sn}/settings/system-template", route.GetSystemTemplateBySN)
	mux.Get("/api/cloudboot/v1/devices/{sn}/settings/image-template", route.GetImageTemplateBySN)
	mux.Get("/api/cloudboot/v1/devices/{sn}/pxe", route.GetPXE)
	mux.Post("/api/cloudboot/v1/devices/{sn}/pxe", route.GenPXE)
	mux.Get("/api/cloudboot/v1/devices/installations/statistics", route.CountDeviceInstallStatic)
	mux.Put("/api/cloudboot/v1/devices/installations/reinstalls", route.Reinstalls)
	mux.Put("/api/cloudboot/v1/devices/installations/autoreinstalls", route.AutoReinstalls)  // 调用规则引擎自动生成装机参数并发起部署
	mux.Put("/api/cloudboot/v1/devices/installations/cancels", route.CancelInstalls)
	mux.Put("/api/cloudboot/v1/devices/installations/setinstallsok", route.SetInstallsOK)
	mux.Delete("/api/cloudboot/v1/devices/settings", route.RemoveDeviceSettings)
	mux.Post("/api/cloudboot/v1/devices/{sn}/centos6/uefi/pxe", route.GenPXE4CentOS6UEFI)
	mux.Post("/api/cloudboot/v1/devices/installations/os-reinstallations", route.SaveDeviceSettingsAndReinstalls)  // 保存装机参数并发起部署，完成后恢复运营状态
	mux.Post("/api/cloudboot/v1/devices/settings/save", route.SaveDeviceSettingsWithoutInstalls)

	// 硬件模板
	mux.Get("/api/cloudboot/v1/hardware-templates", route.GetHardwareTpls)
	mux.Get("/api/cloudboot/v1/hardware-templates/{id}", route.GetHardwareTemplateByID)
	mux.Delete("/api/cloudboot/v1/hardware-templates/{id}", route.RemoveHardTemplateByID)
	mux.Put("/api/cloudboot/v1/hardware-templates/{id}", route.UpdateHardwareTemplate)
	mux.Post("/api/cloudboot/v1/hardware-templates", route.SaveHardwareTemplate)

	// 操作系统安装模板
	mux.Get("/api/cloudboot/v1/os-templates", route.GetTemplatesByCond)

	// 镜像模板
	mux.Post("/api/cloudboot/v1/image-templates", route.SaveImageTemplate)
	mux.Put("/api/cloudboot/v1/image-templates/{id}", route.UpdateImageTemplateByID)
	mux.Get("/api/cloudboot/v1/image-templates/{id}", route.GetImageTemplateByID)
	mux.Delete("/api/cloudboot/v1/image-templates/{id}", route.RemoveImageTemplateByID)
	mux.Get("/api/cloudboot/v1/image-templates", route.GetImageTemplatePage)

	// 系统模板
	mux.Post("/api/cloudboot/v1/system-templates", route.SaveSystemTemplate)
	mux.Put("/api/cloudboot/v1/system-templates/{id}", route.UpdateSystemTemplateByID)
	mux.Get("/api/cloudboot/v1/system-templates/{id}", route.GetSystemTemplateByID)
	mux.Delete("/api/cloudboot/v1/system-templates/{id}", route.RemoveSystemTemplateByID)
	mux.Get("/api/cloudboot/v1/system-templates", route.GetSystemTemplatePage)

	// 操作系统安装日志
	mux.Get("/api/cloudboot/v1/devices/{device_setting_id}/installations/logs", route.GetDeviceLogByDeviceSettingID)

	// 带外管理
	mux.Get("/api/cloudboot/v1/devices/{sn}/oob-user", route.GetOOBUserBySN)
	mux.Post("/api/cloudboot/v1/devices/power", route.OOBPowerOn)
	mux.Put("/api/cloudboot/v1/devices/power/pxe/restart", route.OOBPowerPxeRestart)
	mux.Put("/api/cloudboot/v1/devices/power/restart", route.OOBPowerRestart)
	mux.Get("/api/cloudboot/v1/devices/{sn}/power/status", route.DevicePowerStatus)
	mux.Delete("/api/cloudboot/v1/devices/power", route.OOBPowerOff)
	mux.Put("/api/cloudboot/v1/devices/{sn}/oob/password", route.UpdateOOBPasswordBySN)
	mux.Put("/api/cloudboot/v1/devices/oob/re-access", route.ReAccessOOB) //重新纳管带外
	mux.Get("/api/cloudboot/v1/devices/oob/export", route.ExportOOB)      //导出带外信息
	mux.Get("/api/cloudboot/v1/devices/{sn}/oob/log", route.GetOOBlogBySN)
	mux.Get("/api/cloudboot/v1/devices/oob/inspection", route.OOBInspectionOperate)

	// 任务
	mux.Post("/api/cloudboot/v1/jobs/inspections", route.AddInspectionJob)
	mux.Get("/api/cloudboot/v1/jobs", route.GetJobPage)
	mux.Get("/api/cloudboot/v1/jobs/{job_id}", route.GetJobByID)
	mux.Delete("/api/cloudboot/v1/jobs/{job_id}", route.RemoveJob)
	mux.Put("/api/cloudboot/v1/jobs/{job_id}/pausing", route.PauseJob)
	mux.Put("/api/cloudboot/v1/jobs/{job_id}/unpausing", route.UnpauseJob)
	// 数据字典
	mux.Get("/api/cloudboot/v1/dictionaries", route.GetDataDicts)

	// 硬件巡检
	mux.Get("/api/cloudboot/v1/devices/inspections/statistics", route.GetInspectionStatistics)
	mux.Get("/api/cloudboot/v1/devices/inspections", route.GetInspectionPage)
	mux.Get("/api/cloudboot/v1/devices/inspections/records", route.GetInspectionRecordsPage)
	mux.Get("/api/cloudboot/v1/devices/{sn}/inspections", route.GetInspectionBySN)
	mux.Get("/api/cloudboot/v1/devices/{sn}/inspections/start-times", route.GetInspectionStartTimesBySN)

	// 权限
	mux.Get("/api/cloudboot/v1/permissions/codes", route.PermissionCodeTree)
	// 用户
	mux.Get("/api/cloudboot/v1/users", route.GetUserPage)
	mux.Get("/api/cloudboot/v1/users/info", route.GetUserByToken)

	mux.Put("/api/cloudboot/v1/users/password", route.ChangeUserPassword)

	// 系统配置
	mux.Get("/api/cloudboot/v1/system/login/settings", route.GetSystemLoginSetting)

	// 操作记录
	mux.Get("/api/cloudboot/v1/operate/log", route.GetOperateLogPage)
	mux.Get("/api/cloudboot/v1/api/log", route.GetAPILogPage)

	// 审批
	mux.Post("/api/cloudboot/v1/approvals/idc/abolish", route.SubmitIDCAbolishApproval)
	mux.Post("/api/cloudboot/v1/approvals/server-room/abolish", route.SubmitServerRoomAbolishApproval)
	mux.Post("/api/cloudboot/v1/approvals/network-area/offline", route.SubmitNetAreaOfflineApproval)
	mux.Post("/api/cloudboot/v1/approvals/ip/unassign", route.SubmitIPUnassignApproval) //IP回收
	mux.Post("/api/cloudboot/v1/approvals/device/poweroff", route.SubmitDevicePowerOffApproval)
	mux.Post("/api/cloudboot/v1/approvals/device/restart", route.SubmitDeviceRestartApproval)
	mux.Post("/api/cloudboot/v1/approvals/server-cabinets/offlines", route.SubmitCabinetOfflineApproval)
	mux.Post("/api/cloudboot/v1/approvals/server-cabinets/poweroffs", route.SubmitCabinetPowerOffApproval)
	mux.Post("/api/cloudboot/v1/approvals/devices/migrations", route.SubmitDeviceMigrationApproval)
	mux.Post("/api/cloudboot/v1/approvals/devices/migrations/upload", route.UploadMigrationApproval)
	mux.Post("/api/cloudboot/v1/approvals/devices/migrations/imports/previews", route.ImportMigrationApprovalPriview)
	mux.Post("/api/cloudboot/v1/approvals/devices/migrations/imports", route.ImportMigrationApproval)
	mux.Post("/api/cloudboot/v1/approvals/devices/retirements", route.SubmitDeviceRetirementApproval)
	mux.Post("/api/cloudboot/v1/approvals/devices/os-reinstallations", route.SubmitDeviceOSReInstallationApproval)
	mux.Post("/api/cloudboot/v1/approvals/devices/recycle", route.SubmitDeviceRecycleApproval)
	mux.Put("/api/cloudboot/v1/approvals/{approval_id}/step/{approval_step_id}", route.Approve)
	mux.Get("/api/cloudboot/v1/users/initiated/approvals", route.GetMyApprovalPage)
	mux.Get("/api/cloudboot/v1/users/pending/approvals", route.GetApproveByMePage)
	mux.Get("/api/cloudboot/v1/users/approved/approvals", route.GetApprovedByMePage)
	mux.Delete("/api/cloudboot/v1/approvals/{approval_id}", route.RevokeApproval)
	mux.Get("/api/cloudboot/v1/approvals/{approval_id}", route.GetApprovalByID)

	// dhcp ip token
	mux.Delete("/api/cloudboot/v1/devices/{sn}/limiters/tokens", route.ReturnLimiterTokenBySN)
	mux.Delete("/api/cloudboot/v1/devices/limiters/tokens", route.ReturnLimiterTokens) //一键（批量）释放

	// ServerConf
	mux.Get("/api/cloudboot/v1/samba/settings", route.GetSambaConf)
	// 组件日志
	// 组件
	mux.Post("/api/cloudboot/v1/devices/{sn}/components/hw-server/logs", route.SaveHWServerLog)
	mux.Post("/api/cloudboot/v1/devices/{sn}/components/peconfig/logs", route.SavePEConfigLog)
	mux.Post("/api/cloudboot/v1/devices/{sn}/components/cloudboot-agent/logs", route.SaveCloudbootAgentLog)
	mux.Post("/api/cloudboot/v1/devices/{sn}/components/winconfig/logs", route.SaveWinConfigLog)
	mux.Post("/api/cloudboot/v1/devices/{sn}/components/network-config/logs", route.SaveOSConfigLog)
	mux.Post("/api/cloudboot/v1/devices/{sn}/components/image-clone/logs", route.SaveImageCloneLog)

	//订单
	mux.Post("/api/cloudboot/v1/order", route.SaveOrder)
	mux.Delete("/api/cloudboot/v1/orders", route.RemoveOrders)
	mux.Get("/api/cloudboot/v1/order/{id}", route.GetOrderByID)
	mux.Put("/api/cloudboot/v1/order", route.SaveOrder)
	mux.Put("/api/cloudboot/v1/order/status", route.UpdateOrderStatus)
	mux.Get("/api/cloudboot/v1/orders", route.GetOrderPage)
	mux.Get("/api/cloudboot/v1/orders/export", route.ExportOrders)

	//设备类型，用于维护设备类型和硬件配置的关系
	mux.Post("/api/cloudboot/v1/device-category", route.SaveDeviceCategory)
	mux.Delete("/api/cloudboot/v1/device-categories", route.RemoveDeviceCategorys)
	mux.Put("/api/cloudboot/v1/device-category", route.SaveDeviceCategory)
	mux.Get("/api/cloudboot/v1/device-categories", route.GetDeviceCategoryPage)
	mux.Get("/api/cloudboot/v1/device-category/{id}", route.GetDeviceCategoryByID)
	mux.Get("/api/cloudboot/v1/device-categories/query-params/{param_name}", route.GetDeviceCategoryQuerys)

	// 规则表，规则引擎生成装机参数的规则记录
	mux.Get("/api/cloudboot/v1/device-setting-rules", route.GetDeviceSettingRulePage)
	mux.Post("/api/cloudboot/v1/device-setting-rules", route.SaveDeviceSettingRule)
	mux.Put("/api/cloudboot/v1/device-setting-rules", route.SaveDeviceSettingRule)
	mux.Delete("/api/cloudboot/v1/device-setting-rules", route.RemoveDeviceSettingRules)
}