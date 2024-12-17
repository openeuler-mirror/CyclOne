package com.idcos.enterprise.portal.vo;

import com.idcos.cloud.core.common.BaseVO;

/**
 * @author GuanBin
 * @version PortalDeptVO.java, v1 2017/9/26 下午3:05 GuanBin Exp $$
 */
public class PortalDeptTreeVO extends BaseVO {
    private static final long serialVersionUID = 1L;

    private String id;
    private String code;
    private String name;
    private String pid;
    private String sourceType;
    private String status;
    private String managerId;
    private String remark;
    private String tenantId;

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

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getPid() {
        return pid;
    }

    public void setPid(String pid) {
        this.pid = pid;
    }

    public String getSourceType() {
        return sourceType;
    }

    public void setSourceType(String sourceType) {
        this.sourceType = sourceType;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
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
        return "PortalDeptTreeVO{" +
                "id='" + id + '\'' +
                ", code='" + code + '\'' +
                ", name='" + name + '\'' +
                ", pid='" + pid + '\'' +
                ", sourceType='" + sourceType + '\'' +
                ", status='" + status + '\'' +
                ", managerId='" + managerId + '\'' +
                ", remark='" + remark + '\'' +
                ", tenantId='" + tenantId + '\'' +
                '}';
    }
}
