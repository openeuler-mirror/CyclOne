package com.idcos.enterprise.portal.dal.enums;

/**
 * @Author: Dai
 * @Date: 2018/10/16 上午12:06
 * @Description:
 */
public enum SourceTypeEnum {
    /**
     * 本系统
     */
    NATIVE("native", "本系统"),

    /**
     * WEBANK
     */
    WEBANK("webank", "微众");

    /**
     * 枚举code
     */
    private String code;

    /**
     * 枚举描述
     */
    private String description;

    private SourceTypeEnum(String code, String description) {
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
