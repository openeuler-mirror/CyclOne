/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.dal.enums;

/**
 * 用户组类型
 *
 * @author Dana
 * @version PortalUserGroupTypeEnum.java, v1 2017/12/10 上午11:56 Dana Exp $$
 */
public enum PortalUserGroupTypeEnum {
    /**
     * 默认用户组
     */
    DEFAULT("default", "默认用户组");
    /**
     * 枚举code
     */
    private String code;

    /**
     * 枚举描述
     */
    private String description;

    PortalUserGroupTypeEnum(String code, String description) {
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
