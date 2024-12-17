/**
 * yunjikeji Inc.
 * Copyright (c) 2004-2015 All Rights Reserved.
 */
package com.idcos.enterprise.portal.dal.enums;


/**
 * 信息是否激活
 *
 * @author yanlv
 * @version $Id: DeviceActiveEnum.java, v 0.1 2015年3月18日 下午8:12:24 yanlv Exp $
 */
public enum IsActiveEnum {
    /**
     * 已激活
     */
    HAS_ACTIVE("Y", "已激活"),
    /**
     * 未激活
     */
    NO_ACTIVE("N", "未激活"),

    ;

    /**
     * 枚举code
     */
    private String code;

    /**
     * 枚举描述
     */
    private String description;

    private IsActiveEnum(String code, String description) {
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


