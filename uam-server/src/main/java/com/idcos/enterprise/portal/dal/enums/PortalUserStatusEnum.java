/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.dal.enums;

/**
 * 用户状态
 *
 * @author Dana
 * @version PortalUserStatusEnum.java, v1 2017/12/4 下午9:07 Dana Exp $$
 */
public enum PortalUserStatusEnum {
    /**
     * 初始状态
     */
    INIT("INIT", "初始状态"),
    /**
     * 激活状态
     */
    ENABLED("ENABLED", "激活状态"),
    /**
     * 禁用状态
     */
    DISABLED("DISABLED", "禁用状态"),
    /**
     * 锁定状态
     */
    LOCKED("LOCKED", "锁定状态");

    /**
     * 枚举code
     */
    private String code;

    /**
     * 枚举描述
     */
    private String description;

    PortalUserStatusEnum(String code, String description) {
        this.code = code;

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
}
