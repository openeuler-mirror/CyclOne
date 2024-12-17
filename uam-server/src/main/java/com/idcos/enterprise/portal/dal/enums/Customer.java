package com.idcos.enterprise.portal.dal.enums;

/**
 * @Author: Dai
 * @Date: 2018/10/28 3:20 AM
 * @Description:
 */
public enum Customer {
    /**
     * 通用客户
     */
    DEFAULT("default", "通用客户");

    /**
     * 枚举code
     */
    private String code;

    /**
     * 枚举描述
     */
    private String description;

    private Customer(String code, String description) {
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

    public static Customer getByValue(String value) {
        for (Customer customer : values()) {
            if (customer.getCode().equals(value)) {
                return customer;
            }
        }
        //默认default用户
        return DEFAULT;
    }
}
