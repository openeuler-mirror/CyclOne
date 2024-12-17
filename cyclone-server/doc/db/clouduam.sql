-- --------------------------------------------------------
-- 主机:                           x
-- 服务器版本:                        5.7.38-log - MySQL Community Server (GPL)
-- 服务器操作系统:                      Linux
-- HeidiSQL 版本:                  11.0.0.6063
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;


-- 导出 clouduam 的数据库结构
CREATE DATABASE IF NOT EXISTS `clouduam` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci */;
USE `clouduam`;

-- 导出  表 clouduam.PORTAL_DEPT 结构
CREATE TABLE IF NOT EXISTS `PORTAL_DEPT` (
  `ID` varchar(64) NOT NULL,
  `CODE` varchar(255) DEFAULT NULL COMMENT '部门编码',
  `DISPLAY_NAME` varchar(64) DEFAULT NULL COMMENT '部门名称',
  `PARENT_ID` varchar(64) DEFAULT NULL COMMENT '父项目id',
  `STATUS` varchar(64) DEFAULT NULL COMMENT '部门状态',
  `MANAGER_ID` varchar(64) DEFAULT NULL,
  `REMARK` varchar(255) DEFAULT NULL COMMENT '备注',
  `TENANT_ID` varchar(64) DEFAULT NULL COMMENT '租户id',
  `SOURCE_TYPE` varchar(64) NOT NULL DEFAULT 'native' COMMENT '数据同步来源来源，默认值native表示本系统，如果来自第三方，则需重新定义该值',
  `GMT_CREATE` datetime DEFAULT NULL COMMENT '创建日期',
  `GMT_MODIFIED` datetime DEFAULT NULL COMMENT '修改日期',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='部门信息表';

-- 正在导出表  clouduam.PORTAL_DEPT 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `PORTAL_DEPT` DISABLE KEYS */;
/*!40000 ALTER TABLE `PORTAL_DEPT` ENABLE KEYS */;

-- 导出  表 clouduam.PORTAL_DEPT_ROLE_REL 结构
CREATE TABLE IF NOT EXISTS `PORTAL_DEPT_ROLE_REL` (
  `ID` varchar(64) NOT NULL COMMENT 'id',
  `DEPT_ID` varchar(64) NOT NULL COMMENT '部门id',
  `ROLE_ID` varchar(64) NOT NULL COMMENT '角色id',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='部门与角色的关系表';

-- 正在导出表  clouduam.PORTAL_DEPT_ROLE_REL 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `PORTAL_DEPT_ROLE_REL` DISABLE KEYS */;
/*!40000 ALTER TABLE `PORTAL_DEPT_ROLE_REL` ENABLE KEYS */;

-- 导出  表 clouduam.PORTAL_GROUP_ROLE_REL 结构
CREATE TABLE IF NOT EXISTS `PORTAL_GROUP_ROLE_REL` (
  `ID` varchar(64) NOT NULL COMMENT '关系ID',
  `ROLE_ID` varchar(64) NOT NULL COMMENT '角色ID',
  `GROUP_ID` varchar(64) NOT NULL COMMENT '用户组ID',
  `TENANT` varchar(64) DEFAULT NULL COMMENT '租户code',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户组与角色关系表';

-- 正在导出表  clouduam.PORTAL_GROUP_ROLE_REL 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `PORTAL_GROUP_ROLE_REL` DISABLE KEYS */;
REPLACE INTO `PORTAL_GROUP_ROLE_REL` (`ID`, `ROLE_ID`, `GROUP_ID`, `TENANT`) VALUES
	('25352137-9405-11ee-a6bb-b4055d569ecd', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', '245c438e-d66a-11e8-bf55-5ce04de286f6', 'default');
/*!40000 ALTER TABLE `PORTAL_GROUP_ROLE_REL` ENABLE KEYS */;

-- 导出  表 clouduam.PORTAL_GROUP_USER_REL 结构
CREATE TABLE IF NOT EXISTS `PORTAL_GROUP_USER_REL` (
  `ID` varchar(64) NOT NULL COMMENT '关系ID',
  `USER_ID` varchar(64) NOT NULL COMMENT '用户ID',
  `GROUP_ID` varchar(64) NOT NULL COMMENT '用户组ID',
  `TENANT` varchar(64) DEFAULT NULL COMMENT '租户code'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户组用户关系表';

-- 正在导出表  clouduam.PORTAL_GROUP_USER_REL 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `PORTAL_GROUP_USER_REL` DISABLE KEYS */;
REPLACE INTO `PORTAL_GROUP_USER_REL` (`ID`, `USER_ID`, `GROUP_ID`, `TENANT`) VALUES
	('253513dd-9405-11ee-a6bb-b4055d569ecd', '59df5960cd6ac35f53135b31', '245c438e-d66a-11e8-bf55-5ce04de286f6', 'default');
/*!40000 ALTER TABLE `PORTAL_GROUP_USER_REL` ENABLE KEYS */;

-- 导出  表 clouduam.PORTAL_PERMISSION 结构
CREATE TABLE IF NOT EXISTS `PORTAL_PERMISSION` (
  `ID` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL COMMENT '授权信息ID',
  `APP_ID` varchar(64) NOT NULL COMMENT '应用系统名称',
  `AUTH_RES_TYPE` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL COMMENT '权限资源类型',
  `AUTH_RES_ID` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL COMMENT '权限资源ID',
  `AUTH_RES_NAME` varchar(64) DEFAULT NULL,
  `AUTH_OBJ_ID` varchar(64) NOT NULL COMMENT '授权对象ID',
  `AUTH_OBJ_TYPE` varchar(64) NOT NULL COMMENT '授权对象类型',
  `TENANT` varchar(64) DEFAULT NULL COMMENT '租户code',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='授权信息表';

-- 正在导出表  clouduam.PORTAL_PERMISSION 的数据：~209 rows (大约)
/*!40000 ALTER TABLE `PORTAL_PERMISSION` DISABLE KEYS */;
REPLACE INTO `PORTAL_PERMISSION` (`ID`, `APP_ID`, `AUTH_RES_TYPE`, `AUTH_RES_ID`, `AUTH_RES_NAME`, `AUTH_OBJ_ID`, `AUTH_OBJ_TYPE`, `TENANT`) VALUES
	('008b6fe2-08db-4b24-8501-9f04a31a6122', 'cloudboot', 'cloudboot_button_permission', 'button_ip_unassign', '取消IP分配', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('008cc5ea-fba6-48f4-b42e-f525e3523c47', 'cloudboot', 'cloudboot_button_permission', 'button_oob_re_access', '重新纳管带外', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('0164e8f6-7259-4167-b3e3-88279b54e9bc', 'cloudboot', 'cloudboot_menu_permission', 'menu_oob_info', '带外信息', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('01b6c6da-c94c-4236-a32d-6b80a3c692bd', 'cloudboot', 'cloudboot_menu_permission', 'menu_special_device', '特殊设备列表', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('02c414c2-13ca-4b4c-9cd2-15cd3dab0b7a', 'cloudboot', 'cloudboot_button_permission', 'button_order_export', '导出订单', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('03baba80-3339-46b6-a633-125a0bd9963f', 'cloudboot', 'cloudboot_button_permission', 'button_device_setting_reInstall', '重新部署', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('03cb37a5-6ae7-4c44-a7e8-76fcea3840de', 'cloudboot', 'cloudboot_menu_permission', 'menu_user_management', '用户管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('07fbd788-5ef5-40b8-9de8-e02567652cac', 'cloudboot', 'cloudboot_button_permission', 'button_approval_idc_abolish', '数据中心裁撤审批', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('08d22ddd-2923-4420-9358-326a9e70847d', 'cloudboot', 'cloudboot_button_permission', 'button_server_cabinet_delete', '删除机架', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('0ba142cf-a559-4c6c-9337-b0675b48b5ee', 'cloudboot', 'cloudboot_button_permission', 'button_server_room_production', '投产机房信息', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('0bb9ac7a-4d09-4ad8-8f10-444f67accc19', 'cloudboot', 'cloudboot_button_permission', 'button_idc_create', '新增数据中心', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('0bcf082f-a528-47cf-84db-57e6dedba8d9', 'cloudboot', 'cloudboot_menu_permission', 'menu_audit', '操作审计', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('0c49acec-6168-4317-ba9d-7b1c72bac16f', 'cloudboot', 'cloudboot_menu_permission', 'menu_device_setting_rule', '装机参数规则列表', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('0d44ec17-a614-46a3-b76e-9c79c680608a', 'cloudboot', 'cloudboot_button_permission', 'button_device_category_delete', '删除设备类型', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('0d4f345f-9c0f-4de8-9cf5-9b35a8492356', 'cloudboot', 'cloudboot_button_permission', 'button_server_room_delete', '删除机房信息', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('0ef22348-90eb-4294-8768-e9e1ab6a6cc2', 'cloudboot', 'cloudboot_button_permission', 'button_predeploy_physical_machine_reBoot', '待部署物理机重启', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('10825a54-f82e-46f1-a25d-59343b7b405a', 'cloudboot', 'cloudboot_menu_permission', 'menu_physical_machine', '物理机列表', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('1152433a-b872-42a7-87d7-e7ba7d114757', 'cloudboot', 'cloudboot_button_permission', 'button_server_cabinet_update', '修改机架', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('150b8408-98b5-4da0-874c-f22d4a56c2dd', 'cloudboot', 'cloudboot_button_permission', 'button_special_device_import', '导入特殊设备', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('155d0f55-2aee-4664-b915-553c6320df20', 'cloudboot', 'cloudboot_button_permission', 'button_dhcp_token_release', '释放DHCP令牌', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('18938e6c-f63d-41e6-a4b7-5e2970817b39', 'cloudboot', 'cloudboot_menu_permission', 'menu_network_area', '网络区域', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('1aa701ce-cca9-4c28-a544-79e2663b1990', 'cloudboot', 'cloudboot_button_permission', 'button_virtual_cabinet_delete', '删除虚拟货架', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('1abf64ac-9442-4bf3-9e1b-d640c7c47b91', 'cloudboot', 'cloudboot_button_permission', 'button_network_device_delete', '删除网络设备', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('1afe76fa-feb2-456f-8662-94b17f2fa496', 'cloudboot', 'cloudboot_button_permission', 'button_predeploy_physical_machine_download', '下载待部署物理机模板', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('1bbb97ce-3f2c-43b2-90de-3bf08331a836', 'cloudboot', 'cloudboot_menu_permission', 'menu_ip', 'IP分配', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('1c00804f-8fac-4995-8682-92beb3bb4e9b', 'cloudboot', 'cloudboot_button_permission', 'button_approval_physical_machine_recycle_reInstall', '物理机回收重装审批', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('1c74c354-4e22-4461-84a6-78ad6e4f1db7', 'cloudboot', 'cloudboot_button_permission', 'menu_order', '订单列表', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('1d1bcf39-bdf4-41c4-956b-5e98192a4980', 'cloudboot', 'cloudboot_menu_permission', 'menu_template_management', '配置管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('1fe902b4-e65c-4ca7-98ee-e2b128105fee', 'cloudboot', 'cloudboot_button_permission', 'button_approval_agree', '审批通过', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('205bb1e9-7d31-41e1-bb8f-75a374e64c62', 'cloudboot', 'cloudboot_button_permission', 'menu_idc', '数据中心', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('2072266b-3c07-4c2d-a147-22afcc9250bf', 'cloudboot', 'cloudboot_button_permission', 'button_approval_cabinet_powerOff', '机架关电审批', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('22d9c296-8b1b-4652-9a46-1d8672047014', 'cloudboot', 'cloudboot_menu_permission', 'menu_device_setting', '部署列表', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('22e271ac-187b-4b48-a6d7-d6340b6b64db', 'cloudboot', 'cloudboot_button_permission', 'button_predeploy_physical_machine_import', '导入待部署物理机模板', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('240429ef-1927-4325-a82b-40c470be8249', 'cloudboot', 'cloudboot_button_permission', 'menu_special_device', '特殊设备列表', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('253532ff-9405-11ee-a6bb-b4055d569ecd', 'clouduam', 'CLOUDUAM_MENU', 'user', '用户管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('253533e2-9405-11ee-a6bb-b4055d569ecd', 'clouduam', 'CLOUDUAM_MENU', 'userGroup', '用户组管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('25353448-9405-11ee-a6bb-b4055d569ecd', 'clouduam', 'CLOUDUAM_MENU', 'role', '角色管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('2535346e-9405-11ee-a6bb-b4055d569ecd', 'clouduam', 'CLOUDUAM_MENU', 'resource', '权限资源分配', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('25353497-9405-11ee-a6bb-b4055d569ecd', 'clouduam', 'CLOUDUAM_MENU', 'permission', '权限资源管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('2957f3ff-9ebd-42b3-83c9-046695577a0e', 'cloudboot', 'cloudboot_button_permission', 'button_server_usite_create', '新增机位', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('2ca6fd26-7543-45fd-b758-b8d5067dd6d2', 'cloudboot', 'cloudboot_button_permission', 'button_ip_network_create', '新增IP网段管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('2d2f6b4a-9730-4dd5-b41d-4d3ef073c121', 'cloudboot', 'cloudboot_button_permission', 'button_approval_physical_machine_recycle_move', '物理机回收搬迁审批', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('2d412d01-5495-46c4-bf7d-3111f6449951', 'cloudboot', 'cloudboot_button_permission', 'menu_network_device', '网络设备', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('2dae1e2e-620f-4a6a-b3ab-adc3ed5b382d', 'cloudboot', 'cloudboot_button_permission', 'button_physical_machine_update_usage', '批量修改物理机用途', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('2ed1422a-3e1c-4396-9bd8-e837b394ea0e', 'cloudboot', 'cloudboot_button_permission', 'button_ip_network_import_download', '下载导入IP网段模板', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('31afdabc-601a-44ed-8968-e7d1f4602ed9', 'cloudboot', 'cloudboot_button_permission', 'button_device_setting_rule_delete', '删除装机参数规则', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('31b7f964-20ce-4384-b1af-aff12327a820', 'cloudboot', 'cloudboot_button_permission', 'button_hardware_template_create', '新建/克隆硬件配置', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('31fa3965-a937-48c3-80d3-ce1506437c53', 'CloudUam', 'CLOUDUAM_MENU', 'user.operate', '用户操作（编辑、删除等）', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('32010218-5f6f-45a5-aa7e-88e0f7385776', 'cloudboot', 'cloudboot_button_permission', 'button_server_room_download', '下载机房信息模板', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('32115f99-4f19-4fc5-a81b-4d9bad6c9520', 'cloudboot', 'cloudboot_button_permission', 'button_approval_cabinet_offline', '机架下线审批', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('32abfdcb-fb4b-4be0-8da7-e382b18142d3', 'cloudboot', 'cloudboot_button_permission', 'button_special_device_create', '新增特殊设备', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('32f1943b-4d57-4ee0-9418-93a3a0f25482', 'cloudboot', 'cloudboot_menu_permission', 'menu_server_usite', '机位', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('339cba02-b66b-4217-9120-4b80ad6c8e88', 'cloudboot', 'cloudboot_menu_permission', 'menu_inspection', '硬件巡检', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('34a4a2a3-fd72-4b47-aca1-15ca37adda0f', 'cloudboot', 'cloudboot_button_permission', 'button_device_category_update', '修改设备类型', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('37c9c6ac-90f6-4f50-a908-c1636aaa5b32', 'cloudboot', 'cloudboot_button_permission', 'menu_ip', 'IP分配', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('38986e16-5054-476f-bac0-0ecfa6aa4d88', 'CloudUam', 'CLOUDUAM_MENU', 'user.import', '导入用户', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('391c2c73-95af-43d4-a6d3-11869e6958b0', 'cloudboot', 'cloudboot_button_permission', 'button_predeploy_physical_machine_powerOff', '待部署物理机关电', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('3c0885e6-2074-4d2f-956a-b9ba2810cba4', 'cloudboot', 'cloudboot_button_permission', 'button_network_area_offline', '下线网络区域', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('3dc38034-4b11-4d1f-9cb9-9b7f81f8c922', 'cloudboot', 'cloudboot_button_permission', 'button_server_usite_remark', '更新机位备注', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('3e291a9d-b318-49a4-8226-d00ef829684a', 'cloudboot', 'cloudboot_button_permission', 'button_inspection_inspect_all', '巡检全部', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('3f7a94db-7e79-4bb4-ba33-acf9cf773199', 'cloudboot', 'cloudboot_button_permission', 'menu_approval', '审批管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('3f942298-1499-44e6-977d-6cef05fa2b14', 'cloudboot', 'cloudboot_button_permission', 'button_idc_update', '修改数据中心', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('3fcfaa9b-666c-428d-9348-60b33f639e1c', 'cloudboot', 'cloudboot_button_permission', 'button_device_setting_delete', '删除部署物理机', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('43e403f5-15ba-480b-bc8b-fc9372bc14db', 'cloudboot', 'cloudboot_button_permission', 'button_store_room_update', '修改库房管理单元', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('4614567f-a8cb-4f2c-9ddb-224609f56775', 'cloudboot', 'cloudboot_button_permission', 'button_order_create', '新增订单', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('4620b95a-53c7-4b3e-94ff-5109da2b2b02', 'cloudboot', 'cloudboot_button_permission', 'button_network_device_sync', '同步网络设备', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('476ce033-812a-4f35-a46e-e9bd490d8172', 'cloudboot', 'cloudboot_button_permission', 'button_device_setting_rule_create', '新增装机参数规则', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('4846cca0-014b-402c-875f-082a9b4297bd', 'cloudboot', 'cloudboot_button_permission', 'button_mirror_template_create', '新建镜像配置', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('49375f78-6eb0-4de6-bf95-bca359f4a25e', 'cloudboot', 'cloudboot_button_permission', 'menu_network_area', '网络区域', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('4b3d54ed-61cd-493f-951a-8ed4ec67be40', 'cloudboot', 'cloudboot_button_permission', 'button_approval_physical_machine_move_download', '下载物理机搬迁审批导入模板', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('4b95e082-72a5-493c-91e7-b7a6f01287bc', 'cloudboot', 'cloudboot_button_permission', 'button_device_store_import', '物理机导入到库房', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('4bb746c0-957c-4b3c-a34f-d31eadc3f2b5', 'cloudboot', 'cloudboot_button_permission', 'button_idc_production', '投产数据中心', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('4d3fcda6-589d-4756-be83-dcaa98522863', 'cloudboot', 'cloudboot_button_permission', 'button_predeploy_physical_machine_osInstall', '待部署物理机申请上架部署', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('516793ab-212c-46d3-8e97-b07c4c17b250', 'cloudboot', 'cloudboot_menu_permission', 'menu_hardware_template', '硬件配置', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('5392455e-da22-4f3e-89a2-2abcb25807ce', 'cloudboot', 'cloudboot_button_permission', 'menu_physical_machine', '物理机列表', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('54f551e2-fdd1-434d-a4d7-5bc14e2a06b2', 'cloudboot', 'cloudboot_button_permission', 'menu_predeploy_physical_machine', '待部署物理机', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('5600c39d-627c-4fe0-bd6b-1df004fdc3a2', 'cloudboot', 'cloudboot_button_permission', 'button_network_area_create', '新增网络区域', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('575a3d6e-cfdd-48f9-b774-89dac2edf2bd', 'cloudboot', 'cloudboot_button_permission', 'menu_hardware_template', '硬件配置', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('57fcda67-9e95-4054-9185-74a2a4c4b2f6', 'cloudboot', 'cloudboot_menu_permission', 'menu_physical_machine_management', '物理机管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('58c4528b-e91f-4ce9-83a6-6fea56396499', 'cloudboot', 'cloudboot_button_permission', 'button_server_cabinet_locked', '锁定机架', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('5b794b4d-9b9d-4ab5-b9ac-82ae52cc65e8', 'cloudboot', 'cloudboot_menu_permission', 'menu_system_template', '装机配置', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('5da07f63-f95a-4aed-872b-a42615903d7a', 'cloudboot', 'cloudboot_button_permission', 'button_approval_physical_machine_retirement', '物理机退役审批', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('60ce6a29-d4e0-4575-a020-49a0f68d9e5a', 'cloudboot', 'cloudboot_button_permission', 'button_server_cabinet_type', '更新机架类型', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('60d5faf0-89e7-4af5-9d50-47ae40839cc4', 'cloudboot', 'cloudboot_menu_permission', 'menu_store_room', '库房', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('61603ed4-e5cd-42d9-8cc6-61168f731aed', 'cloudboot', 'cloudboot_button_permission', 'button_device_category_create', '新增设备类型', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('61a20339-6fd4-477f-beaa-1870ee23ceb4', 'cloudboot', 'cloudboot_button_permission', 'button_server_usite_import', '导入机位', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('63ddc1bd-560f-4e76-b842-08f79e245830', 'cloudboot', 'cloudboot_button_permission', 'menu_template_management', '配置管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('640219d8-9476-44f9-b5c7-97d9e248f5d6', 'cloudboot', 'cloudboot_menu_permission', 'menu_server_cabinet', '机架', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('64474fe0-b314-4061-af03-a13280b1917c', 'cloudboot', 'cloudboot_button_permission', 'button_inspection_addTask', '新建巡检任务', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('64e1e28c-0bd2-4a14-b8f8-480cdff222ae', 'cloudboot', 'cloudboot_button_permission', 'menu_server_usite', '机位', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('6519decd-69b1-488f-92f4-058a5976a31b', 'cloudboot', 'cloudboot_button_permission', 'button_order_confirm', '确认订单', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('65a01cf7-4fc5-4a85-878a-0b69b24116ad', 'cloudboot', 'cloudboot_button_permission', 'button_idc_delete', '删除数据中心', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('6612aab8-1273-41a2-afe1-3e98da57bee0', 'cloudboot', 'cloudboot_button_permission', 'button_mirror_template_delete', '删除镜像配置', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('665bd3f2-f777-4895-aefc-c76449dbc157', 'CloudUam', 'CLOUDUAM_MENU', 'user.list', '查看用户列表', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('687b8c82-b495-4b9e-89d9-7f07e775c132', 'cloudboot', 'cloudboot_button_permission', 'menu_ip_network', 'IP网段', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('69e011c4-3628-4562-9cf5-fe4e0685b28d', 'cloudboot', 'cloudboot_button_permission', 'button_approval_physical_machine_reInstall', '物理机重装审批', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('69eb7e0b-c6f4-4578-a27c-8d44eea78335', 'cloudboot', 'cloudboot_button_permission', 'button_server_room_abolished', '裁撤机房信息', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('6a24c0ef-9d6e-458e-9ec9-cd270c65c119', 'cloudboot', 'cloudboot_button_permission', 'button_approval_revoke', '审批撤销', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('6b1eb99c-1715-4a2b-824d-0a5e519b7227', 'cloudboot', 'cloudboot_button_permission', 'button_physical_machine_download', '存量物理机模板下载', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('6c28d1c9-3d04-4670-9c76-8be1d7674c0d', 'cloudboot', 'cloudboot_button_permission', 'button_idc_abolished', '裁撤数据中心', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('6d8a9bd8-5aac-4bc7-86e3-fee1a3783f62', 'cloudboot', 'cloudboot_button_permission', 'button_predeploy_physical_machine_delete', '删除物理机', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('7008b44b-339c-46bf-860c-1f4268f76ec7', 'cloudboot', 'cloudboot_button_permission', 'button_system_template_update', '修改PXE配置', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('716bf4b2-e6cd-433d-960d-27bdb7360ab2', 'cloudboot', 'cloudboot_button_permission', 'menu_oob_info', '带外管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('757c0a4c-675d-4849-8c15-ca4f170bf1e8', 'cloudboot', 'cloudboot_button_permission', 'button_idc_accepted', '验收数据中心', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('765e444d-eb62-4665-83dc-72a32dea9e9b', 'cloudboot', 'cloudboot_button_permission', 'menu_device_category', '设备类型列表', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('76de34ce-6f1a-4c3b-b9ed-94be42051521', 'cloudboot', 'cloudboot_button_permission', 'button_server_room_import', '导入机房信息', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('795a0407-377c-4a8f-a6c3-cff29d575ff4', 'cloudboot', 'cloudboot_button_permission', 'menu_inspection', '硬件巡检', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('7a6db539-c1a9-4c53-a1b7-0f3cd0b04938', 'cloudboot', 'cloudboot_button_permission', 'button_server_cabinet_powerOn', '开电机架', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('7b2546bf-c4a2-4d4c-b004-f14e4b08515c', 'cloudboot', 'cloudboot_button_permission', 'button_task_pause', '暂停', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('7c00dea9-f913-4c8f-9c4c-7157b1df730a', 'cloudboot', 'cloudboot_button_permission', 'button_dhcp_token_batch_release', '批量释放DHCP令牌', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('7c13e51f-f585-4303-80a3-e3542038a3c6', 'cloudboot', 'cloudboot_button_permission', 'button_network_area_delete', '删除网络区域', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('7e4000f2-0b60-4c46-a942-cef6bab2d582', 'cloudboot', 'cloudboot_button_permission', 'button_server_usite_delete', '删除机位', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('7eb259e0-16d5-49db-ac39-2f0753cc6950', 'cloudboot', 'cloudboot_button_permission', 'button_network_area_download', '下载网络区域模板', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('7fdf964c-c7d6-46b7-9040-b752ee299101', 'cloudboot', 'cloudboot_button_permission', 'button_predeploy_physical_machine_powerOn', '待部署物理机开电', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('80ee7836-07a1-47ed-b597-3d9f0e4ed101', 'cloudboot', 'cloudboot_button_permission', 'button_network_area_update', '修改网络区域', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('81063b9d-c0e2-44b6-91f8-d2f8ec8e2680', 'cloudboot', 'cloudboot_menu_permission', 'menu_home', '概览', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('81a74a60-6034-4820-8d6d-ae6a27e14b09', 'cloudboot', 'cloudboot_menu_permission', 'menu_order', '订单列表', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('82764f72-74c7-4d60-9d06-95dcf134fbd0', 'cloudboot', 'cloudboot_button_permission', 'button_task_continue', '继续', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('838694b5-f558-49ca-8a6c-cd3fc784783b', 'cloudboot', 'cloudboot_button_permission', 'menu_device_setting_rule', '装机参数规则列表', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('862f2fe0-87a5-42b0-ba7e-abede83ecfdc', 'cloudboot', 'cloudboot_button_permission', 'button_server_cabinet_download', '下载机架模板', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('87190dd4-5cc4-4191-bec4-d5dff399ca34', 'cloudboot', 'cloudboot_button_permission', 'button_approval_physical_machine_power_off', '物理机关电审批', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('8781d800-2095-471c-bfa2-2894621f782f', 'cloudboot', 'cloudboot_menu_permission', 'menu_server_room', '机房', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('886d9326-9a33-4164-a403-00863135d653', 'cloudboot', 'cloudboot_button_permission', 'button_network_device_import', '导入网络设备', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('8dcbcf24-f4b1-449c-a2b6-d74d66736b9c', 'cloudboot', 'cloudboot_button_permission', 'button_hardware_template_delete', '删除硬件配置', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('8f248e02-c9f7-4d01-b8ce-47b46ae29992', 'cloudboot', 'cloudboot_button_permission', 'button_server_usite_import_port', '导入机位端口', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('91b24f3f-84d1-4f14-a5f5-283ab5657257', 'cloudboot', 'cloudboot_button_permission', 'button_physical_machine_update_oob', '修改物理机带外', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('9249de2c-9b47-43bb-84a4-6c12c931fcec', 'cloudboot', 'cloudboot_button_permission', 'button_network_device_create', '新增网络设备', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('9306325f-6741-4ae6-a39a-6211e4d2ff43', 'cloudboot', 'cloudboot_button_permission', 'button_task_delete', '删除', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('938e1772-5779-4c7f-b1cf-6d1093932218', 'CloudUam', 'CLOUDUAM_MENU', 'user.create', '新建用户', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('9411d788-b8dc-4d82-8dc7-461bb7345f7f', 'cloudboot', 'cloudboot_button_permission', 'button_approval_physical_machine_move', '物理机搬迁审批', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('966c2b5a-2fb6-49f1-a3be-733fce3d4224', 'cloudboot', 'cloudboot_button_permission', 'button_ip_network_update', '修改IP网段管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('974029b0-0201-48f3-9b04-dc38ae6baf86', 'cloudboot', 'cloudboot_button_permission', 'button_device_store_import_download', '物理机导入到库房模板下载', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('9aa3b53a-07be-4753-b9e0-fecbaa25eabd', 'cloudboot', 'cloudboot_button_permission', 'button_system_template_delete', '删除PXE配置', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('9d83ac48-c61c-42d9-8494-578bd0397c29', 'cloudboot', 'cloudboot_button_permission', 'button_store_room_delete', '删除库房管理单元', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('9f0ab8c1-e9b0-4fce-ad42-930b514623c3', 'cloudboot', 'cloudboot_button_permission', 'button_ip_unassign', 'IP回收审批', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('a0a756ff-d1a1-4bf2-b457-96fb9596b391', 'cloudboot', 'cloudboot_menu_permission', 'menu_task_management', '任务管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('a0cbe0b0-b511-434c-ab9d-6e88153195a6', 'CloudUam', 'CLOUDUAM_MENU', 'userGroup.list', '查看用户组列表', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('a189e7b4-b9e2-464c-82d3-73db66eb53f2', 'cloudboot', 'cloudboot_button_permission', 'button_inspection_inspect', '重新巡检', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('a1b46bd7-3bb0-4bec-85d4-e7b1771dfe13', 'cloudboot', 'cloudboot_menu_permission', 'menu_idc', '数据中心', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('a31b593c-3204-4d07-9b3d-a85407d12fa3', 'cloudboot', 'cloudboot_button_permission', 'menu_task_management', '任务管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('a3428f07-0193-4ed5-a989-8d5aa80c64ca', 'cloudboot', 'cloudboot_button_permission', 'menu_system_template', '系统/镜像配置', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('a455c488-8c61-4e90-9285-37c151a79970', 'cloudboot', 'cloudboot_button_permission', 'button_server_usite_update', '修改机位', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('a478ca5d-35f2-4e85-b142-320216d3b2e9', 'cloudboot', 'cloudboot_button_permission', 'button_order_cancel', '取消订单', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('a5119e97-b83d-4808-aa71-15da5dce1902', 'cloudboot', 'cloudboot_button_permission', 'button_network_device_import_download', '下载导入网络设备模板', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('a7858387-79ef-4a7f-bd30-a408a305f3fb', 'CloudUam', 'CLOUDUAM_MENU', 'user.export', '导出用户', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('a93aad67-70b5-4e67-9df0-a27c2bb7c491', 'cloudboot', 'cloudboot_button_permission', 'button_network_area_production', '投产网络区域', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('aa12db28-1197-4e27-b50e-b9c1edd39ee9', 'cloudboot', 'cloudboot_button_permission', 'button_server_usite_status', '更新机位状态', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('ab3d8496-269c-4c70-a543-1d7db58b4675', 'cloudboot', 'cloudboot_button_permission', 'button_ip_network_import', '导入IP网段', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('af32d53c-a412-4dd0-ac5e-2e9a2cfd8cda', 'cloudboot', 'cloudboot_button_permission', 'button_special_device_import_download', '导入特殊设备模板下载', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('b06fde8d-bbd6-46e8-8ff1-a72d4718f227', 'cloudboot', 'cloudboot_button_permission', 'button_system_template_create', '新建PXE配置', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('b2a6cad0-18e1-4540-9483-d4c9157cc035', 'cloudboot', 'cloudboot_button_permission', 'button_approval_server_room_abolish', '机房裁撤审批', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('b4417088-e113-4d78-9a26-4b9330e0d085', 'cloudboot', 'cloudboot_button_permission', 'button_server_cabinet_create', '新增机架', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('b4fc2ee4-e727-42ac-ac9b-01b237d7f5b7', 'CloudUam', 'CLOUDUAM_MENU', 'user.deptOperate', '部门操作（编辑、删除等）', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('b585b3b6-46ba-45de-aa59-0669e206b852', 'cloudboot', 'cloudboot_button_permission', 'menu_order_management', '订单管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('b59b503b-9f81-4a85-bb99-349a2a56f0d5', 'cloudboot', 'cloudboot_button_permission', 'button_physical_oob_export', '导出带外', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('b5ccf06c-1a9c-45a2-bc47-0b89360ea8de', 'cloudboot', 'cloudboot_menu_permission', 'menu_predeploy_physical_machine', '待部署物理机', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('b68c62c5-df30-4da6-9e9e-cbbfedb99208', 'cloudboot', 'cloudboot_button_permission', 'button_physical_machine_networkBoot', '从网卡启动物理机', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('b696c8ac-fca4-40b0-99f8-4a2081083aec', 'cloudboot', 'cloudboot_button_permission', 'button_mirror_template_update', '修改镜像配置', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('b69ca84f-87e6-4397-a561-24fd98e02431', 'cloudboot', 'cloudboot_button_permission', 'button_hardware_template_update', '修改硬件配置', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('b89cc9ec-e55b-420d-b259-0af77cdd01c4', 'cloudboot', 'cloudboot_button_permission', 'button_physical_machine_powerOn', '开电物理机', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('b93e49c1-19e6-490e-9f7e-71b770a221ec', 'cloudboot', 'cloudboot_button_permission', 'menu_network_management', '网络管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('b9d74fa4-15d2-42b1-8e0a-74b13f37792a', 'cloudboot', 'cloudboot_button_permission', 'button_physical_machine_update_status', '批量修改物理机状态', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('bb28fecb-d15a-4756-b681-41ac5c46869a', 'cloudboot', 'cloudboot_menu_permission', 'menu_approval', '审批管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('bc709f25-9666-4b4c-bf46-7a869a51547f', 'cloudboot', 'cloudboot_button_permission', 'button_server_room_accepted', '验收机房信息', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('bd56f7fd-33b5-47d9-87f4-d1cee768bf65', 'cloudboot', 'cloudboot_button_permission', 'button_physical_machine_import', '存量物理机导入', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('bd5701c3-ef3c-401d-bd24-9d88c41fe65e', 'cloudboot', 'cloudboot_button_permission', 'button_special_device_delete', '新增特殊设备', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('bf3cc067-efd9-4ee9-bab6-dcbc7f24cb9c', 'cloudboot', 'cloudboot_button_permission', 'button_physical_machine_export', '物理机导出', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('c1ff4521-8694-4ee9-a963-9deb268bf38e', 'cloudboot', 'cloudboot_button_permission', 'menu_server_room', '机房', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('c239d77a-1e4c-4270-a407-702803a5ffe3', 'cloudboot', 'cloudboot_menu_permission', 'menu_network_management', '网络管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('c26d5101-2a99-4a5b-b06e-a3882ce6b4af', 'cloudboot', 'cloudboot_menu_permission', 'menu_device_category', '设备类型列表', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('c294b623-40aa-4ed1-a8bb-c7d44993a5f3', 'cloudboot', 'cloudboot_button_permission', 'button_server_cabinet_remark', '更新机架备注', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('c430ad2d-0a26-4f14-a200-e484e76dcc40', 'cloudboot', 'cloudboot_button_permission', 'button_device_setting_rule_update', '修改装机参数规则', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('c4ad9da6-601b-4c1a-b79a-f0395f5ac6d9', 'cloudboot', 'cloudboot_menu_permission', 'menu_ip_network', 'IP网段', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('c59434ea-3c84-4bfd-b4b7-d52f62594029', 'cloudboot', 'cloudboot_button_permission', 'menu_idc_management', '数据中心管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('c6f4b165-304e-45ae-807e-566b5b5de858', 'cloudboot', 'cloudboot_button_permission', 'button_server_room_update', '修改机房信息', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('c733c10a-096e-4ea6-918c-b2523c461c06', 'cloudboot', 'cloudboot_button_permission', 'button_network_area_import', '导入网络区域', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('c764de6e-a574-4d2d-adfd-b5984155abc1', 'cloudboot', 'cloudboot_menu_permission', 'menu_audit_api', '接口调用记录', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('c7a0c8ea-bc4c-4407-8ffc-164931667816', 'cloudboot', 'cloudboot_button_permission', 'button_store_room_create', '新增库房管理单元', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('c7ae3bc8-9ce1-4f91-8e34-d0eac33f3540', 'cloudboot', 'cloudboot_button_permission', 'button_server_cabinet_powerOff', '关电机架', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('c88d2e66-bdb3-4d20-af25-8307787cab08', 'cloudboot', 'cloudboot_button_permission', 'button_approval_physical_machine_recycle_retire', '物理机回收退役审批', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('c99e93f1-ad57-4948-b05c-7e6039d5aa88', 'CloudUam', 'CLOUDUAM_MENU', 'login', '登录', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('cb380632-c71a-48c6-a28b-e82ee465b046', 'cloudboot', 'cloudboot_button_permission', 'button_server_usite_delete_port', '删除机位端口', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('cb95c48b-2531-440f-aa54-e9a448bae621', 'cloudboot', 'cloudboot_menu_permission', 'menu_idc_management', '数据中心管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('cc0bc233-e679-4bd0-8481-3d74255640ce', 'cloudboot', 'cloudboot_button_permission', 'menu_device_setting', '部署列表', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('d1350200-1aaf-45ab-8fa7-646513d5fabf', 'CloudUam', 'CLOUDUAM_MENU', 'userGroup.operate', '用户组操作（编辑、删除等）', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('d1bd69c2-8de1-4e2b-980a-9414ab6f4fde', 'cloudboot', 'cloudboot_button_permission', 'button_predeploy_physical_machine_networkBoot', '待部署物理机从网卡启动', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('d33f7d03-871a-4b5e-917e-922206a01ff4', 'cloudboot', 'cloudboot_button_permission', 'button_order_delete', '删除订单', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('d51a3796-bbe9-4701-a5ce-91bd556bbea0', 'cloudboot', 'cloudboot_button_permission', 'button_physical_machine_update', '修改物理机', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('d6bf9f40-537c-47c9-a900-0fb9c4d13e03', 'cloudboot', 'cloudboot_button_permission', 'button_approval_physical_machine_move_import', '导入物理机搬迁审批', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('d9a53b9f-dc12-4ab4-9478-2ab446dba580', 'cloudboot', 'cloudboot_button_permission', 'menu_server_cabinet', '机架', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('dadb1dda-faf8-4ddd-82a7-d0390286585d', 'cloudboot', 'cloudboot_menu_permission', 'menu_order_management', '订单管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('dadba20a-db4e-4e4c-b6a8-a7f69dae08ee', 'cloudboot', 'cloudboot_menu_permission', 'menu_audit_log', '操作记录', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('dc00afa8-2bde-4131-a9dc-9b130222de81', 'cloudboot', 'cloudboot_button_permission', 'button_server_cabinet_enabled', '启用机架', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('dc7576d2-fbd6-49cb-a66f-792b03084b7b', 'cloudboot', 'cloudboot_button_permission', 'button_server_cabinet_offline', '下线机架', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('ddc0eb42-138c-4175-8061-ecdd8c093e2e', 'cloudboot', 'cloudboot_button_permission', 'button_approval_disagree', '审批拒绝', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('df0fd248-e289-43d2-8c32-804cdfbcb4b3', 'cloudboot', 'cloudboot_button_permission', 'button_server_cabinet_import', '导入机架', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('df1649b2-b577-4a2c-b678-f125de056819', 'cloudboot', 'cloudboot_button_permission', 'button_store_room_import', '导入库房管理单元', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('df7b5fa6-db35-47b4-a3b7-893127c01f9e', 'cloudboot', 'cloudboot_button_permission', 'menu_physical_machine_management', '物理机管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('dfab5389-07b5-489e-95f4-792e323fc075', 'cloudboot', 'cloudboot_button_permission', 'button_virtual_cabinet_create', '新增虚拟货架', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('e00db732-81d6-42a9-bf4c-5a844be768cf', 'cloudboot', 'cloudboot_button_permission', 'button_ip_assign', 'IP分配', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('e118501a-25b2-48fc-a66b-94914d772e24', 'cloudboot', 'cloudboot_button_permission', 'button_server_usite_download', '下载机位模板', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('e2d5ee81-2420-4c90-a637-6e9b40d021ee', 'cloudboot', 'cloudboot_button_permission', 'button_approval_network_area_offline', '网络区域下线审批', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('e3b0a4d4-8b45-42f9-b382-e3eb13b55fd0', 'cloudboot', 'cloudboot_button_permission', 'menu_store_room', '库房信息管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('e714da44-a718-4659-8d7f-8f204ac23c29', 'cloudboot', 'cloudboot_button_permission', 'button_device_setting_cancelInstall', '取消部署', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('e7aa158c-c14b-423e-aed3-74287537b5aa', 'cloudboot', 'cloudboot_button_permission', 'button_server_room_create', '新增机房信息', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('e8a5e5c6-c4d4-4893-be0c-6bcf7e26b400', 'cloudboot', 'cloudboot_button_permission', 'button_ip_network_delete', '删除IP网段管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('ec0bc4b0-0569-4a9a-933e-16d4becc9007', 'cloudboot', 'cloudboot_button_permission', 'button_approval_physical_machine_restart', '物理机重启审批', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('edd85dfd-89aa-46e3-a225-0ea6e5e6f36d', 'cloudboot', 'cloudboot_menu_permission', 'menu_network_device', '网络设备', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('ee45c3b2-748c-4628-89d2-8e95ac33d7c6', 'cloudboot', 'cloudboot_button_permission', 'button_server_usite_download_port', '下载机位端口模板', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('f0ba296a-0082-42b2-99cb-3b44023543a3', 'CloudUam', 'CLOUDUAM_MENU', 'userGroup.create', '新建用户组', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('f4ef743c-d1d1-455e-9d5a-c53482857722', 'CloudUam', 'CLOUDUAM_MENU', 'tenant', '租户管理', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('f5c51fbd-ad26-4a98-903a-5e1c57ad17b6', 'cloudboot', 'cloudboot_button_permission', 'button_ip_network_sync', '同步IP网段', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default'),
	('fc6cc7ed-399d-44e5-95de-d0f663aaa131', 'CloudUam', 'CLOUDUAM_MENU', 'user.deptTree', '查看部门树', '3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'ROLE', 'default');
/*!40000 ALTER TABLE `PORTAL_PERMISSION` ENABLE KEYS */;

-- 导出  表 clouduam.PORTAL_RESOURCE 结构
CREATE TABLE IF NOT EXISTS `PORTAL_RESOURCE` (
  `ID` varchar(64) NOT NULL,
  `APP_ID` varchar(64) NOT NULL COMMENT '应用系统名称',
  `CODE` varchar(64) NOT NULL COMMENT '权限资源类型',
  `NAME` varchar(128) NOT NULL COMMENT '权限资源名称',
  `URL` varchar(128) NOT NULL COMMENT '权限资源URL',
  `REMARK` varchar(256) NOT NULL COMMENT '备注',
  `IS_ACTIVE` varchar(1) NOT NULL DEFAULT 'Y',
  `TENANT` varchar(64) DEFAULT NULL COMMENT '租户code',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='权限资源表';

-- 正在导出表  clouduam.PORTAL_RESOURCE 的数据：~3 rows (大约)
/*!40000 ALTER TABLE `PORTAL_RESOURCE` DISABLE KEYS */;
REPLACE INTO `PORTAL_RESOURCE` (`ID`, `APP_ID`, `CODE`, `NAME`, `URL`, `REMARK`, `IS_ACTIVE`, `TENANT`) VALUES
	('2534b1fb-9405-11ee-a6bb-b4055d569ecd', 'CloudUam', 'CLOUDUAM_MENU', 'CLOUDUAM菜单权限', 'http://127.0.0.1:8092/uam/permission/menuCodes', 'CLOUDUAM菜单权限', 'Y', 'default'),
	('2de4ca83-75a9-4ec3-b3db-bea54ec60f23', 'cloudboot', 'cloudboot_menu_permission', 'cloudboot菜单权限集合', 'http://127.0.0.1:8083/api/cloudboot/v1/permissions/codes?type=cloudboot_menu_permission', 'cloudboot菜单权限集合', 'Y', 'default'),
	('7fa1a8c8-b568-46c7-92d7-cd712e239cc8', 'cloudboot', 'cloudboot_button_permission', 'cloudboot按钮权限集合', 'http://127.0.0.1:8083/api/cloudboot/v1/permissions/codes?type=cloudboot_button_permission', 'cloudboot按钮权限集合', 'Y', 'default');
/*!40000 ALTER TABLE `PORTAL_RESOURCE` ENABLE KEYS */;

-- 导出  表 clouduam.PORTAL_ROLE 结构
CREATE TABLE IF NOT EXISTS `PORTAL_ROLE` (
  `ID` varchar(64) NOT NULL COMMENT '角色ID',
  `CODE` varchar(64) DEFAULT NULL COMMENT '角色编码',
  `NAME` varchar(64) NOT NULL COMMENT '角色名称',
  `REMARK` varchar(256) DEFAULT NULL COMMENT '角色描述',
  `CREATE_USER` varchar(64) DEFAULT NULL COMMENT '创建人',
  `GMT_CREATE` datetime DEFAULT NULL COMMENT '创建时间',
  `GMT_MODIFIED` datetime DEFAULT NULL COMMENT '最后修改时间',
  `IS_ACTIVE` char(1) DEFAULT NULL COMMENT '是否有效',
  `TENANT` varchar(64) DEFAULT NULL COMMENT '租户code',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='角色信息表';

-- 正在导出表  clouduam.PORTAL_ROLE 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `PORTAL_ROLE` DISABLE KEYS */;
REPLACE INTO `PORTAL_ROLE` (`ID`, `CODE`, `NAME`, `REMARK`, `CREATE_USER`, `GMT_CREATE`, `GMT_MODIFIED`, `IS_ACTIVE`, `TENANT`) VALUES
	('3cb5fde4-d66a-11e8-bf55-5ce04de286f6', 'SUPER_ADMIN', '超级管理员', '超级管理员', 'system', '2023-12-06 15:00:31', '2023-12-06 15:00:31', 'Y', 'default');
/*!40000 ALTER TABLE `PORTAL_ROLE` ENABLE KEYS */;

-- 导出  表 clouduam.PORTAL_SYS_DICT 结构
CREATE TABLE IF NOT EXISTS `PORTAL_SYS_DICT` (
  `TYPE_CODE` varchar(64) NOT NULL COMMENT '系统字典类型编码',
  `CODE` varchar(64) NOT NULL COMMENT '系统字典编码',
  `VALUE` text COMMENT '参数值',
  `TENANT_ID` varchar(64) NOT NULL COMMENT '租户code',
  `REMARK` text COMMENT '说明',
  PRIMARY KEY (`TENANT_ID`,`CODE`,`TYPE_CODE`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='uam系统参数表';

-- 正在导出表  clouduam.PORTAL_SYS_DICT 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `PORTAL_SYS_DICT` DISABLE KEYS */;
REPLACE INTO `PORTAL_SYS_DICT` (`TYPE_CODE`, `CODE`, `VALUE`, `TENANT_ID`, `REMARK`) VALUES
	('userGroupType', 'default', '默认', 'default', '用户组类型default，值为默认');
/*!40000 ALTER TABLE `PORTAL_SYS_DICT` ENABLE KEYS */;

-- 导出  表 clouduam.PORTAL_TENANT 结构
CREATE TABLE IF NOT EXISTS `PORTAL_TENANT` (
  `ID` varchar(64) NOT NULL,
  `NAME` varchar(128) DEFAULT NULL COMMENT '租户编码',
  `DISPLAY_NAME` varchar(128) DEFAULT NULL COMMENT '租户名称',
  `GMT_CREATE` datetime DEFAULT NULL COMMENT '创建日期',
  `GMT_MODIFIED` datetime NOT NULL COMMENT '修改日期',
  `IS_ACTIVE` char(1) DEFAULT NULL COMMENT '是否可用',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='租户信息表';

-- 正在导出表  clouduam.PORTAL_TENANT 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `PORTAL_TENANT` DISABLE KEYS */;
REPLACE INTO `PORTAL_TENANT` (`ID`, `NAME`, `DISPLAY_NAME`, `GMT_CREATE`, `GMT_MODIFIED`, `IS_ACTIVE`) VALUES
	('59df4f30e4b088cbabc4cb4d', 'default', '管理租户', '2023-12-06 15:00:31', '2023-12-06 15:00:31', 'Y');
/*!40000 ALTER TABLE `PORTAL_TENANT` ENABLE KEYS */;

-- 导出  表 clouduam.PORTAL_TOKEN 结构
CREATE TABLE IF NOT EXISTS `PORTAL_TOKEN` (
  `ID` varchar(64) NOT NULL,
  `NAME` text NOT NULL COMMENT 'token值',
  `TOKEN_CRC` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'token串的crc哈希值',
  `LOGIN_ID` varchar(64) DEFAULT NULL COMMENT '登录名',
  `TENANT_ID` varchar(64) NOT NULL COMMENT '租户code',
  `IS_ACTIVE` char(1) DEFAULT NULL COMMENT '是否可用',
  `EXPIRE_TIME` datetime NOT NULL COMMENT 'token过期时间',
  `GMT_CREATE` datetime DEFAULT NULL COMMENT '创建日期',
  `GMT_MODIFIED` datetime NOT NULL COMMENT '修改日期',
  `REMARK` longtext COMMENT '备注',
  PRIMARY KEY (`ID`),
  KEY `token_hash_crc` (`TOKEN_CRC`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='token信息表';

-- 正在导出表  clouduam.PORTAL_TOKEN 的数据：~11 rows (大约)
/*!40000 ALTER TABLE `PORTAL_TOKEN` DISABLE KEYS */;
REPLACE INTO `PORTAL_TOKEN` (`ID`, `NAME`, `TOKEN_CRC`, `LOGIN_ID`, `TENANT_ID`, `IS_ACTIVE`, `EXPIRE_TIME`, `GMT_CREATE`, `GMT_MODIFIED`, `REMARK`) VALUES
	('13e136a7-cfd7-4d51-afa7-a1c9c3c13246', 'eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiLnrqHnkIblkZgiLCJ1c2VySWQiOiI1OWRmNTk2MGNkNmFjMzVmNTMxMzViMzEiLCJuYW1lIjoi566h55CG5ZGYIiwibG9naW5JZCI6ImFkbWluIiwibG9naW5OYW1lIjoiYWRtaW4iLCJ0ZW5hbnRJZCI6ImRlZmF1bHQiLCJ0aW1lb3V0IjoyMTYwMCwiZXhwIjoxNzAxOTMyNDYxLCJjcmVhdFRpbWUiOjE3MDE5MTA4NjExMzcsInRlbmFudE5hbWUiOiLnrqHnkIbnp5_miLcifQ.zhdFeCAHJFlQHgXNkKACS4Zd3VmTQJtDBwHTIH_qFZQ', 4002311361, 'admin', 'default', 'Y', '2023-12-07 15:01:01', '2023-12-07 09:01:01', '2023-12-07 09:01:01', NULL),
	('20801c97-a049-4054-994c-c3b502ada8c1', 'eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiLnrqHnkIblkZgiLCJ1c2VySWQiOiI1OWRmNTk2MGNkNmFjMzVmNTMxMzViMzEiLCJuYW1lIjoi566h55CG5ZGYIiwibG9naW5JZCI6ImFkbWluIiwibG9naW5OYW1lIjoiYWRtaW4iLCJ0ZW5hbnRJZCI6ImRlZmF1bHQiLCJ0aW1lb3V0IjoyMTYwMCwiZXhwIjoxNzAyNDY5NzgyLCJjcmVhdFRpbWUiOjE3MDI0NDgxODI1MzAsInRlbmFudE5hbWUiOiLnrqHnkIbnp5_miLcifQ.tt_KhA4owFJsGZApElY6nxkUR_hftNHS0ihnlN_CxjI', 2690808154, 'admin', 'default', 'Y', '2023-12-13 20:16:22', '2023-12-13 14:16:23', '2023-12-13 14:16:23', NULL),
	('25354401-9405-11ee-a6bb-b4055d569ecd', 'eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiLnrqHnkIblkZgiLCJ1c2VySWQiOiI1OWRmNTk2MGNkNmFjMzVmNTMxMzViMzEiLCJuYW1lIjoi566h55CG5ZGYIiwibG9naW5JZCI6ImFkbWluIiwibG9naW5OYW1lIjoiYWRtaW4iLCJ0ZW5hbnRJZCI6ImRlZmF1bHQiLCJ0aW1lb3V0IjoyNDAwMCwiZXhwIjoxODkzNDkyMDAwLCJjcmVhdFRpbWUiOjE1MjgzNDEwNjA0ODcsInRlbmFudE5hbWUiOiLnrqHnkIbnp5_miLcifQ._Iq6xfjaTfM7OXSJu-pyQycfK3ycGqlV-bm9_mSXh6g', 3062456733, 'admin', 'default', 'Y', '2030-01-01 18:00:00', '2023-12-06 15:00:31', '2023-12-06 15:00:31', '系统配置的永久token'),
	('2c6c79c9-b69a-4c93-a43f-263489cf7aab', 'eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiLnrqHnkIblkZgiLCJ1c2VySWQiOiI1OWRmNTk2MGNkNmFjMzVmNTMxMzViMzEiLCJuYW1lIjoi566h55CG5ZGYIiwibG9naW5JZCI6ImFkbWluIiwibG9naW5OYW1lIjoiYWRtaW4iLCJ0ZW5hbnRJZCI6ImRlZmF1bHQiLCJ0aW1lb3V0IjoyMTYwMCwiZXhwIjoxNzAyMTA4MzMyLCJjcmVhdFRpbWUiOjE3MDIwODY3MzI3NzMsInRlbmFudE5hbWUiOiLnrqHnkIbnp5_miLcifQ.hXFHQC5fVnfcJlHxuq5Dv6yexl52T11ihHWYsCrFYvE', 2455366502, 'admin', 'default', 'Y', '2023-12-09 15:52:12', '2023-12-09 09:52:13', '2023-12-09 09:52:13', NULL),
	('2e63431b-f8e2-43ca-9a0b-1181a728cd9b', 'eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiLnrqHnkIblkZgiLCJ1c2VySWQiOiI1OWRmNTk2MGNkNmFjMzVmNTMxMzViMzEiLCJuYW1lIjoi566h55CG5ZGYIiwibG9naW5JZCI6ImFkbWluIiwibG9naW5OYW1lIjoiYWRtaW4iLCJ0ZW5hbnRJZCI6ImRlZmF1bHQiLCJ0aW1lb3V0IjoyMTYwMCwiZXhwIjoxNzAxODg4NzM1LCJjcmVhdFRpbWUiOjE3MDE4NjcxMzU1NzIsInRlbmFudE5hbWUiOiLnrqHnkIbnp5_miLcifQ.R2dyTcDsD5MPWru2bYSN84l2VJ-RQoJIDfhyKCJDbTA', 3527168430, 'admin', 'default', 'Y', '2023-12-07 02:52:15', '2023-12-06 20:52:16', '2023-12-06 20:52:16', NULL),
	('3f17d062-4a74-47be-9ca4-ec41ef390f9f', 'eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiLnrqHnkIblkZgiLCJ1c2VySWQiOiI1OWRmNTk2MGNkNmFjMzVmNTMxMzViMzEiLCJuYW1lIjoi566h55CG5ZGYIiwibG9naW5JZCI6ImFkbWluIiwibG9naW5OYW1lIjoiYWRtaW4iLCJ0ZW5hbnRJZCI6ImRlZmF1bHQiLCJ0aW1lb3V0IjoyMTYwMCwiZXhwIjoxNzAyNTM1ODAzLCJjcmVhdFRpbWUiOjE3MDI1MTQyMDMzMDYsInRlbmFudE5hbWUiOiLnrqHnkIbnp5_miLcifQ.9t_Ymtik-kQy-ih2WXCV_Qna1M7IvRSYK_rrMOuf_yA', 2463626145, 'admin', 'default', 'Y', '2023-12-14 14:36:43', '2023-12-14 08:36:43', '2023-12-14 08:36:43', NULL),
	('6e07257a-871f-4f9a-ac90-fb99865d3b9d', 'eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiLnrqHnkIblkZgiLCJ1c2VySWQiOiI1OWRmNTk2MGNkNmFjMzVmNTMxMzViMzEiLCJuYW1lIjoi566h55CG5ZGYIiwibG9naW5JZCI6ImFkbWluIiwibG9naW5OYW1lIjoiYWRtaW4iLCJ0ZW5hbnRJZCI6ImRlZmF1bHQiLCJ0aW1lb3V0IjoyMTYwMCwiZXhwIjoxNzAxOTMwNjY0LCJjcmVhdFRpbWUiOjE3MDE5MDkwNjQyNjcsInRlbmFudE5hbWUiOiLnrqHnkIbnp5_miLcifQ.12FSyVvSfIZ_SG1gKzOSOf3PEPIXoT1KIxaek5YX4Yo', 560336943, 'admin', 'default', 'Y', '2023-12-07 14:31:04', '2023-12-07 08:31:04', '2023-12-07 08:31:04', NULL),
	('74ccc709-797f-437c-9476-4a29bd42204b', 'eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiLnrqHnkIblkZgiLCJ1c2VySWQiOiI1OWRmNTk2MGNkNmFjMzVmNTMxMzViMzEiLCJuYW1lIjoi566h55CG5ZGYIiwibG9naW5JZCI6ImFkbWluIiwibG9naW5OYW1lIjoiYWRtaW4iLCJ0ZW5hbnRJZCI6ImRlZmF1bHQiLCJ0aW1lb3V0IjoyMTYwMCwiZXhwIjoxNzAyNTY2NjE5LCJjcmVhdFRpbWUiOjE3MDI1NDUwMTkxMTYsInRlbmFudE5hbWUiOiLnrqHnkIbnp5_miLcifQ.osjpGjFxzIrx60oZ6SzvLl2e9ZS1IVjJ429rXcAnejQ', 2560445608, 'admin', 'default', 'Y', '2023-12-14 23:10:19', '2023-12-14 17:10:19', '2023-12-14 17:10:19', NULL),
	('7f53711c-8e55-412c-badd-d7de7827e9bc', 'eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiLnrqHnkIblkZgiLCJ1c2VySWQiOiI1OWRmNTk2MGNkNmFjMzVmNTMxMzViMzEiLCJuYW1lIjoi566h55CG5ZGYIiwibG9naW5JZCI6ImFkbWluIiwibG9naW5OYW1lIjoiYWRtaW4iLCJ0ZW5hbnRJZCI6ImRlZmF1bHQiLCJ0aW1lb3V0IjoyMTYwMCwiZXhwIjoxNzAxOTMyOTEyLCJjcmVhdFRpbWUiOjE3MDE5MTEzMTIwMjAsInRlbmFudE5hbWUiOiLnrqHnkIbnp5_miLcifQ.-ikjTKeA0-1wy1_BjCYjNRKCmdyue87OEWMbfJqG9iM', 3014900625, 'admin', 'default', 'Y', '2023-12-07 15:08:32', '2023-12-07 09:08:32', '2023-12-07 09:08:32', NULL),
	('afbd4797-c54e-4c9d-97ec-edad10eb511b', 'eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiLnrqHnkIblkZgiLCJ1c2VySWQiOiI1OWRmNTk2MGNkNmFjMzVmNTMxMzViMzEiLCJuYW1lIjoi566h55CG5ZGYIiwibG9naW5JZCI6ImFkbWluIiwibG9naW5OYW1lIjoiYWRtaW4iLCJ0ZW5hbnRJZCI6ImRlZmF1bHQiLCJ0aW1lb3V0IjoyMTYwMCwiZXhwIjoxNzAxODg5MTU1LCJjcmVhdFRpbWUiOjE3MDE4Njc1NTUzMzcsInRlbmFudE5hbWUiOiLnrqHnkIbnp5_miLcifQ.wWe-5k7FPoTyBY0ShdBgl9L5TpvrLKxAtF1Vzpe0QEI', 3784738346, 'admin', 'default', 'Y', '2023-12-07 02:59:15', '2023-12-06 20:59:15', '2023-12-06 20:59:15', NULL),
	('bbc56d0e-547f-4feb-a4c8-11aaa25f59e8', 'eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiLnrqHnkIblkZgiLCJ1c2VySWQiOiI1OWRmNTk2MGNkNmFjMzVmNTMxMzViMzEiLCJuYW1lIjoi566h55CG5ZGYIiwibG9naW5JZCI6ImFkbWluIiwibG9naW5OYW1lIjoiYWRtaW4iLCJ0ZW5hbnRJZCI6ImRlZmF1bHQiLCJ0aW1lb3V0IjoyMTYwMCwiZXhwIjoxNzAxODg4NzI1LCJjcmVhdFRpbWUiOjE3MDE4NjcxMjU1NTIsInRlbmFudE5hbWUiOiLnrqHnkIbnp5_miLcifQ.tJZKswRYzikHKuzYpRQe7IG2Qrx7GNrGpAwMaQ5Bhxc', 3261116001, 'admin', 'default', 'Y', '2023-12-07 02:52:05', '2023-12-06 20:52:06', '2023-12-06 20:52:06', NULL);
/*!40000 ALTER TABLE `PORTAL_TOKEN` ENABLE KEYS */;

-- 导出  表 clouduam.PORTAL_USER 结构
CREATE TABLE IF NOT EXISTS `PORTAL_USER` (
  `ID` varchar(64) NOT NULL COMMENT '用户ID',
  `NAME` varchar(64) DEFAULT NULL,
  `LOGIN_ID` varchar(64) DEFAULT NULL COMMENT '登录名',
  `DEPT_ID` varchar(64) DEFAULT NULL COMMENT '部门id',
  `TITLE` varchar(64) DEFAULT NULL COMMENT '标签',
  `MOBILE1` varchar(64) DEFAULT NULL,
  `WEIXIN` varchar(64) DEFAULT NULL COMMENT '标签',
  `EMAIL` varchar(256) DEFAULT NULL COMMENT '邮箱',
  `IS_ACTIVE` char(1) DEFAULT NULL COMMENT '是否有效',
  `STATUS` varchar(64) NOT NULL DEFAULT 'INIT' COMMENT '用户类型：INIT、ENABLED、DISABLED、LOCKED',
  `REMARK` varchar(256) DEFAULT NULL COMMENT '备注',
  `TENANT_ID` varchar(64) DEFAULT NULL COMMENT '租户编码',
  `SOURCE_TYPE` varchar(64) NOT NULL DEFAULT 'native' COMMENT '数据同步来源来源，默认值native表示本系统，如果来自第三方，则需重新定义该值',
  `LAST_MODIFIED_TIME` datetime DEFAULT NULL COMMENT '最后登录时间',
  `CREATE_TIME` datetime DEFAULT NULL COMMENT '创建时间',
  `LAST_LOGIN_TIME` datetime DEFAULT NULL COMMENT '最后修改时间',
  `PASSWORD` varchar(255) DEFAULT NULL COMMENT '密码',
  `SALT` varchar(255) DEFAULT NULL COMMENT '密码盐',
  `MOBILE2` varchar(64) DEFAULT NULL,
  `RTX` varchar(32) DEFAULT NULL COMMENT 'RTX沟通工具',
  `OFFICE_TEL1` varchar(64) DEFAULT NULL,
  `OFFICE_TEL2` varchar(64) DEFAULT NULL,
  `EMPLOYEE_TYPE` varchar(64) DEFAULT NULL,
  `CREATE_USER` varchar(64) DEFAULT NULL COMMENT '创建人',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户信息表';

-- 正在导出表  clouduam.PORTAL_USER 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `PORTAL_USER` DISABLE KEYS */;
REPLACE INTO `PORTAL_USER` (`ID`, `NAME`, `LOGIN_ID`, `DEPT_ID`, `TITLE`, `MOBILE1`, `WEIXIN`, `EMAIL`, `IS_ACTIVE`, `STATUS`, `REMARK`, `TENANT_ID`, `SOURCE_TYPE`, `LAST_MODIFIED_TIME`, `CREATE_TIME`, `LAST_LOGIN_TIME`, `PASSWORD`, `SALT`, `MOBILE2`, `RTX`, `OFFICE_TEL1`, `OFFICE_TEL2`, `EMPLOYEE_TYPE`, `CREATE_USER`) VALUES
	('59df5960cd6ac35f53135b31', '管理员', 'admin', NULL, NULL, NULL, NULL, NULL, 'Y', 'ENABLED', NULL, 'default', 'native', '2023-12-06 15:00:31', '2023-12-06 15:00:31', '2023-12-06 15:00:31', 'c10c1a40b6f8956f057cd95ce35519ae', 'yWw1FKf6BZk=', NULL, NULL, NULL, NULL, NULL, 'system');
/*!40000 ALTER TABLE `PORTAL_USER` ENABLE KEYS */;

-- 导出  表 clouduam.PORTAL_USER_GROUP 结构
CREATE TABLE IF NOT EXISTS `PORTAL_USER_GROUP` (
  `ID` varchar(64) NOT NULL COMMENT '用户组ID',
  `NAME` varchar(128) DEFAULT NULL COMMENT '用户组名称',
  `TYPE` varchar(64) DEFAULT NULL COMMENT '用户组类型',
  `REMARK` varchar(256) DEFAULT NULL COMMENT '备注',
  `CREATE_USER` varchar(64) DEFAULT NULL COMMENT '创建人',
  `GMT_CREATE` datetime DEFAULT NULL COMMENT '创建时间',
  `GMT_MODIFIED` datetime NOT NULL COMMENT '最后修改时间',
  `IS_ACTIVE` char(1) DEFAULT NULL COMMENT '是否有效',
  `TENANT` varchar(64) DEFAULT NULL COMMENT '租户code',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户组信息表';

-- 正在导出表  clouduam.PORTAL_USER_GROUP 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `PORTAL_USER_GROUP` DISABLE KEYS */;
REPLACE INTO `PORTAL_USER_GROUP` (`ID`, `NAME`, `TYPE`, `REMARK`, `CREATE_USER`, `GMT_CREATE`, `GMT_MODIFIED`, `IS_ACTIVE`, `TENANT`) VALUES
	('245c438e-d66a-11e8-bf55-5ce04de286f6', '超级用户组', 'default', '超级用户组', 'system', '2023-12-06 15:00:31', '2023-12-06 15:00:31', 'Y', 'default');
/*!40000 ALTER TABLE `PORTAL_USER_GROUP` ENABLE KEYS */;

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
