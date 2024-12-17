/**
 * 云霁科技网络技术有限公司 idcos.com Inc.
 * Copyright (c) 2004-2015 All Rights Reserved.
 */
package com.idcos.enterprise.portal.dal.enums;

/**
 * @author jiaohuizhe
 * @version $Id: FormInputValueTypeEnum.java, v 0.1 2015年10月12日 下午3:41:00 jiaohuizhe Exp $
 */
public enum FormInputValueTypeEnum {
    /**
     * 文本
     */
    TEXT("TEXT", "文本"),
    /**
     * 整数
     */
    INTEGER("INTEGER", "整数"),
    /**
     * 数值
     */
    NUMBER("NUMBER", "数值"),
    /**
     * 日期
     */
    DATE("DATE", "日期"),
    /**
     * JSON
     */
    JSON("JSON", "JSON"),

    ;

    /**
     * 枚举code
     */
    private String code;

    /**
     * 枚举描述
     */
    private String description;

    private FormInputValueTypeEnum(String code, String description) {
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
