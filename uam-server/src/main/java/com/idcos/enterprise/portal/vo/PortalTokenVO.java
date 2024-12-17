/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.vo;

import com.idcos.cloud.core.common.BaseLinkVO;

import java.util.Date;

/**
 * @author Dana
 * @version PortalTokenVO.java, v1 2017/11/30 下午4:05 Dana Exp $$
 */
public class PortalTokenVO extends BaseLinkVO {
    private static final long serialVersionUID = 1L;

    private String id;

    private String name;

    private String loginId;

    private String tenantId;

    private String isActive;

    private Date expireTime;

    private Date gmtCreate;

    private Date gmtModified;

    private String remark;

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

    public String getLoginId() {
        return loginId;
    }

    public void setLoginId(String loginId) {
        this.loginId = loginId;
    }

    public String getTenantId() {
        return tenantId;
    }

    public void setTenantId(String tenantId) {
        this.tenantId = tenantId;
    }

    public String getIsActive() {
        return isActive;
    }

    public void setIsActive(String isActive) {
        this.isActive = isActive;
    }

    public Date getExpireTime() {
        return expireTime;
    }

    public void setExpireTime(Date expireTime) {
        this.expireTime = expireTime;
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

    public String getRemark() {
        return remark;
    }

    public void setRemark(String remark) {
        this.remark = remark;
    }

    @Override
    public String toString() {
        return "PortalTokenVO{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", loginId='" + loginId + '\'' +
                ", tenantId='" + tenantId + '\'' +
                ", isActive='" + isActive + '\'' +
                ", expireTime=" + expireTime +
                ", gmtCreate=" + gmtCreate +
                ", gmtModified=" + gmtModified +
                ", remark='" + remark + '\'' +
                '}';
    }
}