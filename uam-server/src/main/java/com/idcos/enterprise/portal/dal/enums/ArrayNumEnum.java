/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2018 All Rights Reserved.
 */
package com.idcos.enterprise.portal.dal.enums;

/**
 * @author Xizhao.Dai
 * @version ArrayNumEnum.java, v1 2018/4/29 下午5:35 Xizhao.Dai Exp $$
 */
public enum ArrayNumEnum {
    /**
     * 0
     */
    ZERO(0, "0"),
    /**
     * 1
     */
    ONE(1, "1"),
    /**
     * 2
     */
    TWO(2, "2"),
    /**
     * 3
     */
    THREE(3, "3"),
    /**
     * 4
     */
    FOUR(4, "4"),
    /**
     * 5
     */
    FIVE(5, "5"),
    /**
     * 6
     */
    SIX(6, "6"),
    /**
     * 7
     */
    SEVEN(7, "7"),
    /**
     * 8
     */
    EIGHT(8, "8"),
    /**
     * 9
     */
    NINE(9, "9"),
    /**
     * 10
     */
    TEN(10, "10");

    /**
     * 枚举code
     */
    private int code;

    /**
     * 枚举描述
     */
    private String description;

    ArrayNumEnum(int code, String description) {
        this.code = code;

        this.description = description;
    }

    public int getCode() {
        return code;
    }

    public void setCode(int code) {
        this.code = code;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }
}
