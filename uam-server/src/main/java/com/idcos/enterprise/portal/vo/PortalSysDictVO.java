/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.vo;

import com.idcos.cloud.core.common.BaseVO;

/**
 * @author Dana
 * @version PortalSysDict.java, v1 2017/12/5 上午9:48 Dana Exp $$
 */
public class PortalSysDictVO extends BaseVO {
    private static final long serialVersionUID = 1L;

    private String typeCode;

    private String code;

    private String value;

    private String tenantId;

    private String remark;

    @Override
    public String toString() {
        return super.toString();
    }

    public String getTypeCode() {
        return typeCode;
    }

    public void setTypeCode(String typeCode) {
        this.typeCode = typeCode;
    }

    public String getCode() {
        return code;
    }

    public void setCode(String code) {
        this.code = code;
    }

    public String getValue() {
        return value;
    }

    public void setValue(String value) {
        this.value = value;
    }

    public String getTenantId() {
        return tenantId;
    }

    public void setTenantId(String tenantId) {
        this.tenantId = tenantId;
    }

    public String getRemark() {
        return remark;
    }

    public void setRemark(String remark) {
        this.remark = remark;
    }
}