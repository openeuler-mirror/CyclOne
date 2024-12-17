package com.idcos.enterprise.portal.vo;

import com.idcos.cloud.core.common.BaseLinkVO;

import java.util.Arrays;
import java.util.Date;

/**
 * 用户组界面展示类
 *
 * @author jiaohuizhe
 * @version $Id: PortalUserGroupVO.java, v 0.1 2015年5月8日 上午11:03:56 jiaohuizhe Exp $
 */
public class PortalUserGroupVO extends BaseLinkVO {

    private static final long serialVersionUID = 1L;

    private String id;

    private String name;

    private String type;

    private String remark;

    private String isActive;

    private Date gmtCreate;

    private Date gmtModified;

    private String[] selRoles;

    private String[] selUsers;

    private String tenant;

    private String tenantName;

    public String getTenantName() {
        return tenantName;
    }

    public void setTenantName(String tenantName) {
        this.tenantName = tenantName;
    }

    public String getTenantId() {
        return tenant;
    }

    public void setTenantId(String tenantId) {
        this.tenant = tenant;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public String getRemark() {
        return remark;
    }

    public void setRemark(String remark) {
        this.remark = remark;
    }

    public String getIsActive() {
        return isActive;
    }

    public void setIsActive(String isActive) {
        this.isActive = isActive;
    }

    public String[] getSelRoles() {
        return selRoles;
    }

    public void setSelRoles(String[] selRoles) {
        this.selRoles = selRoles;
    }

    public String[] getSelUsers() {
        return selUsers;
    }

    public void setSelUsers(String[] selUsers) {
        this.selUsers = selUsers;
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
        return "PortalUserGroupVO{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", type='" + type + '\'' +
                ", remark='" + remark + '\'' +
                ", isActive='" + isActive + '\'' +
                ", gmtCreate=" + gmtCreate +
                ", gmtModified=" + gmtModified +
                ", selRoles=" + Arrays.toString(selRoles) +
                ", selUsers=" + Arrays.toString(selUsers) +
                ", tenant='" + tenant + '\'' +
                ", tenantName='" + tenantName + '\'' +
                '}';
    }
}