/**
 * yunjikeji Inc.
 * Copyright (c) 2004-2015 All Rights Reserved.
 */
package com.idcos.enterprise.portal.dal.enums;

/**
 * 业务实体枚举信息，表述具体的实体对象
 *
 * @author yanlv
 * @version $Id: DeviceActiveEnum.java, v 0.1 2015年3月18日 下午8:12:24 yanlv Exp $
 */
public enum BusinessIdentityEnum {
    /**
     * 角色信息
     */
    PORTAL_ROLE("PORTAL_ROLE", "portalRoleRepository", "角色信息"),
    /**
     * 用户信息
     */
    PORTAL_USER("PORTAL_USER", "portalUserRepository", "用户信息"),
    /**
     * 用户组信息
     */
    PORTAL_USER_GROUP("PORTAL_USER_GROUP", "portalUserGroupRepository", "用户组信息"),
    /**
     * 授权信息信息
     */
    PORTAL_PERMISSION("PORTAL_PERMISSION", "portalPermissionRepository", "授权信息信息"),
    /**
     * 权限资源管理
     */
    PORTAL_RESOURCE("PORTAL_RESOURCE", "portalResourceRepository", "权限资源管理"),

    ;
    /**
     * 枚举code
     */
    private String code;

    /**
     * 实体 repository 对应的bean name
     */
    private String beanName;

    /**
     * 枚举描述
     */
    private String description;

    private BusinessIdentityEnum(String code, String beanName, String description) {
        this.code = code;
        this.beanName = beanName;
        this.description = description;
    }

    public String getCode() {
        return code;
    }

    public void setCode(String code) {
        this.code = code;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getBeanName() {
        return beanName;
    }

    public void setBeanName(String beanName) {
        this.beanName = beanName;
    }

}
