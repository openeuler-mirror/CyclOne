/**
 * 云霁科技网络技术有限公司 idcos.com Inc.
 * Copyright (c) 2004-2015 All Rights Reserved.
 */
package com.idcos.enterprise.portal.dal.enums;

/**
 * @author jiaohuizhe
 * @version $Id: FormInputTypeEnum.java, v 0.1 2015年10月12日 下午3:41:00 jiaohuizhe Exp $
 */
public enum FormInputTypeEnum {
    /**
     * 文本框
     */
    INPUT("INPUT", "文本框"),
    /**
     * 下拉框
     */
    SELECT("SELECT", "下拉框"),
    /**
     * 复选框
     */
    CHECKBOX("CHECKBOX", "复选框"),
    /**
     * 单选框
     */
    RADIO("RADIO", "单选框"),
    /**
     * 文本域
     */
    TEXTAREA("TEXTAREA", "文本域"),
    ;

    /**
     * 枚举code
     */
    private String code;

    /**
     * 枚举描述
     */
    private String description;

    private FormInputTypeEnum(String code, String description) {
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
