/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.dal.enums;

/**
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月04 下午2:49 souakiragen Exp $
 */
public enum StatusEnum {
    /**
     * 激活
     */
    Y("1", "激活"),
    /**
     * 删除
     */
    N("0", "删除");

    /**
     * 枚举code
     */
    private String code;

    /**
     * 枚举描述
     */
    private String description;

    private StatusEnum(String code, String description) {
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
