package com.idcos.enterprise.portal.vo;

import com.idcos.cloud.core.common.BaseVO;

/**
 * @author GuanBin
 * @version PortalDeptVO.java, v1 2017/9/26 下午3:05 GuanBin Exp $$
 */
public class PortalDeptVO extends BaseVO {
    private static final long serialVersionUID = 1L;

    private String id;
    private String code;
    private String displayName;
    private String parentId;
    private String status;
    private String sourceType;
    private String managerId;
    private String remark;
    private String tenantId;
    private String tenant;

    public String getTenant() {
        return tenant;
    }

    public void setTenant(String tenant) {
        this.tenant = tenant;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getCode() {
        return code;
    }

    public void setCode(String code) {
        this.code = code;
    }

    public String getDisplayName() {
        return displayName;
    }

    public void setDisplayName(String displayName) {
        this.displayName = displayName;
    }

    public String getParentId() {
        return parentId;
    }

    public void setParentId(String parentId) {
        this.parentId = parentId;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public String getSourceType() {
        return sourceType;
    }

    public void setSourceType(String sourceType) {
        this.sourceType = sourceType;
    }

    public String getManagerId() {
        return managerId;
    }

    public void setManagerId(String managerId) {
        this.managerId = managerId;
    }

    public String getRemark() {
        return remark;
    }

    public void setRemark(String remark) {
        this.remark = remark;
    }

    public String getTenantId() {
        return tenantId;
    }

    public void setTenantId(String tenantId) {
        this.tenantId = tenantId;
    }

    @Override
    public String toString() {
        return "PortalDeptVO{" +
                "id='" + id + '\'' +
                ", code='" + code + '\'' +
                ", displayName='" + displayName + '\'' +
                ", parentId='" + parentId + '\'' +
                ", status='" + status + '\'' +
                ", managerId='" + managerId + '\'' +
                ", remark='" + remark + '\'' +
                ", tenantId='" + tenantId + '\'' +
                ", tenant='" + tenant + '\'' +
                '}';
    }
}
