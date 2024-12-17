/**
 * This class is generated by jOOQ
 */
package com.idcos.enterprise.jooq;


import com.idcos.enterprise.jooq.tables.PortalDept;
import com.idcos.enterprise.jooq.tables.PortalDeptRoleRel;
import com.idcos.enterprise.jooq.tables.PortalGroupRoleRel;
import com.idcos.enterprise.jooq.tables.PortalGroupUserRel;
import com.idcos.enterprise.jooq.tables.PortalPermission;
import com.idcos.enterprise.jooq.tables.PortalResource;
import com.idcos.enterprise.jooq.tables.PortalRole;
import com.idcos.enterprise.jooq.tables.PortalSysDict;
import com.idcos.enterprise.jooq.tables.PortalTenant;
import com.idcos.enterprise.jooq.tables.PortalToken;
import com.idcos.enterprise.jooq.tables.PortalUser;
import com.idcos.enterprise.jooq.tables.PortalUserGroup;

import javax.annotation.Generated;


/**
 * Convenience access to all tables in clouduam
 */
@Generated(
	value = {
		"http://www.jooq.org",
		"jOOQ version:3.7.2"
	},
	comments = "This class is generated by jOOQ"
)
@SuppressWarnings({ "all", "unchecked", "rawtypes" })
public class Tables {

	/**
	 * 部门信息表
	 */
	public static final PortalDept PORTAL_DEPT = com.idcos.enterprise.jooq.tables.PortalDept.PORTAL_DEPT;

	/**
	 * 部门与角色的关系表
	 */
	public static final PortalDeptRoleRel PORTAL_DEPT_ROLE_REL = com.idcos.enterprise.jooq.tables.PortalDeptRoleRel.PORTAL_DEPT_ROLE_REL;

	/**
	 * 用户组与角色关系表
	 */
	public static final PortalGroupRoleRel PORTAL_GROUP_ROLE_REL = com.idcos.enterprise.jooq.tables.PortalGroupRoleRel.PORTAL_GROUP_ROLE_REL;

	/**
	 * 用户组用户关系表
	 */
	public static final PortalGroupUserRel PORTAL_GROUP_USER_REL = com.idcos.enterprise.jooq.tables.PortalGroupUserRel.PORTAL_GROUP_USER_REL;

	/**
	 * 授权信息表
	 */
	public static final PortalPermission PORTAL_PERMISSION = com.idcos.enterprise.jooq.tables.PortalPermission.PORTAL_PERMISSION;

	/**
	 * 权限资源表
	 */
	public static final PortalResource PORTAL_RESOURCE = com.idcos.enterprise.jooq.tables.PortalResource.PORTAL_RESOURCE;

	/**
	 * 角色信息表
	 */
	public static final PortalRole PORTAL_ROLE = com.idcos.enterprise.jooq.tables.PortalRole.PORTAL_ROLE;

	/**
	 * uam系统参数表
	 */
	public static final PortalSysDict PORTAL_SYS_DICT = com.idcos.enterprise.jooq.tables.PortalSysDict.PORTAL_SYS_DICT;

	/**
	 * 租户信息表
	 */
	public static final PortalTenant PORTAL_TENANT = com.idcos.enterprise.jooq.tables.PortalTenant.PORTAL_TENANT;

	/**
	 * token信息表
	 */
	public static final PortalToken PORTAL_TOKEN = com.idcos.enterprise.jooq.tables.PortalToken.PORTAL_TOKEN;

	/**
	 * 用户信息表
	 */
	public static final PortalUser PORTAL_USER = com.idcos.enterprise.jooq.tables.PortalUser.PORTAL_USER;

	/**
	 * 用户组信息表
	 */
	public static final PortalUserGroup PORTAL_USER_GROUP = com.idcos.enterprise.jooq.tables.PortalUserGroup.PORTAL_USER_GROUP;
}