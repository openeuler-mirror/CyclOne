package com.idcos.enterprise.portal.dal.enums;

/**
 * @author
 * @version PortalFilterUtil.java, v1 2017/11/18 下午6:22  Exp $$
 */
public enum RbacMenuEnum {
    /**
     * 用户管理
     */
    USER("USER", "用户管理"),
    /**
     * 用户组管理
     */
    GROUP("GROUP", "用户组管理"),
    /**
     * 角色管理
     */
    ROLE("ROLE", "角色管理"),
    /**
     * 权限资源管理
     */
    RESROUCE("RESOURCE", "权限资源管理"),
    /**
     * 权限资源分配
     */
    PERMISSION("PERMISSION", "权限资源分配"),
    ;

    private String code;

    private String name;

    public String getCode() {
        return code;
    }

    public void setCode(String code) {
        this.code = code;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    RbacMenuEnum(String code, String name) {
        this.code = code;
        this.name = name;
    }
}
