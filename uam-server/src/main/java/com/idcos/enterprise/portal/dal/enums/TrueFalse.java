package com.idcos.enterprise.portal.dal.enums;

/**
 * @author admin
 * @date 16/03/2017
 */
public enum TrueFalse {
    /**
     * TRUE
     */
    TRUE("true", "TRUE"),
    /**
     * FALSE
     */
    FALSE("false", "FALSE");

    /**
     * 枚举编码
     */
    private String code;

    /**
     * 枚举描述信息
     */
    private String description;

    /**
     * @param code
     * @param description
     */
    private TrueFalse(String code, String description) {

        this.code = code;
        this.description = description;
    }

    public String getCode() {
        return code;
    }

    public String getDescription() {
        return description;
    }

}
