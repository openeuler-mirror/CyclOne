package com.idcos.enterprise.portal.vo;

import com.idcos.cloud.core.common.BaseVO;

import java.util.Date;

/**
 * @author GuanBin
 * @version PortalDeptVO.java, v1 2017/9/26 下午3:05 GuanBin Exp $$
 */
public class PortalTenantVO extends BaseVO {
    private static final long serialVersionUID = 1L;

    private String id;
    private String tenantId;
    private String tenantName;
    private Date gmtCreate;
    private Date gmtModified;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getTenantId() {
        return tenantId;
    }

    public void setTenantId(String tenantId) {
        this.tenantId = tenantId;
    }

    public String getTenantName() {
        return tenantName;
    }

    public void setTenantName(String tenantName) {
        this.tenantName = tenantName;
    }

    public Date getGmtCreate() {
        return gmtCreate;
    }

    public void setGmtCreate(Date gmtCreate) {
        this.gmtCreate = gmtCreate;
    }

    public Date getGmtModified() {
        return gmtModified;
    }

    public void setGmtModified(Date gmtModified) {
        this.gmtModified = gmtModified;
    }

    @Override
    public String toString() {
        return "PortalTenantVO{" +
                "id='" + id + '\'' +
                ", tenantId='" + tenantId + '\'' +
                ", tenantName='" + tenantName + '\'' +
                ", gmtCreate=" + gmtCreate +
                ", gmtModified=" + gmtModified +
                '}';
    }
}
