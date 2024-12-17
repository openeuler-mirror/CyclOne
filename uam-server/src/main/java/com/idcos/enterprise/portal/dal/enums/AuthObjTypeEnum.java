package com.idcos.enterprise.portal.dal.enums;

/**
 * @author
 * @version PortalFilterUtil.java, v1 2017/11/18 下午6:22  Exp $$
 */
public enum AuthObjTypeEnum {
    /**
     * ROLE
     */
    ROLE("ROLE"),
    /**
     * NETSEG
     */
    NETSEG("NETSEG");

    private String code;

    public String getCode() {
        return code;
    }

    public void setCode(String code) {
        this.code = code;
    }

    AuthObjTypeEnum(String code) {
        this.code = code;
    }
}
