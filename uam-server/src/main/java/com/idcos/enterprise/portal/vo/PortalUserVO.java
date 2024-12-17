package com.idcos.enterprise.portal.vo;

import com.idcos.cloud.core.common.BaseLinkVO;

import java.util.Arrays;
import java.util.Date;

/**
 * 用户界面展示类
 *
 * @author jiaohuizhe
 * @version $Id: PortalUserVO.java, v 0.1 2015年5月8日 上午11:04:01 jiaohuizhe Exp $
 */
public class PortalUserVO extends BaseLinkVO {

    private static final long serialVersionUID = 1L;

    private String id;
    private String name;
    private String deptId;
    private String deptName;
    private String email;
    private String mobile1;
    private String mobile2;
    private String rtx;
    private String officeTel1;
    private String officeTel2;
    private String employeeType;
    private String title;
    private String weixin;
    private String remark;
    private String isActive;
    private String status;
    private String tenantId;
    private String[] selGroups;
    private String[] selRoles;
    private String loginId;
    private String sourceType;
    private Date lastLoginTime;
    private Date lastModifiedTime;

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

    public String getDeptId() {
        return deptId;
    }

    public void setDeptId(String deptId) {
        this.deptId = deptId;
    }

    public String getDeptName() {
        return deptName;
    }

    public void setDeptName(String deptName) {
        this.deptName = deptName;
    }

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }

    public String getMobile1() {
        return mobile1;
    }

    public void setMobile1(String mobile1) {
        this.mobile1 = mobile1;
    }

    public String getMobile2() {
        return mobile2;
    }

    public void setMobile2(String mobile2) {
        this.mobile2 = mobile2;
    }

    public String getRtx() {
        return rtx;
    }

    public void setRtx(String rtx) {
        this.rtx = rtx;
    }

    public String getOfficeTel1() {
        return officeTel1;
    }

    public void setOfficeTel1(String officeTel1) {
        this.officeTel1 = officeTel1;
    }

    public String getOfficeTel2() {
        return officeTel2;
    }

    public void setOfficeTel2(String officeTel2) {
        this.officeTel2 = officeTel2;
    }

    public String getEmployeeType() {
        return employeeType;
    }

    public void setEmployeeType(String employeeType) {
        this.employeeType = employeeType;
    }

    public String getWeixin() {
        return weixin;
    }

    public void setWeixin(String weixin) {
        this.weixin = weixin;
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

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public String getTenantId() {
        return tenantId;
    }

    public void setTenantId(String tenantId) {
        this.tenantId = tenantId;
    }

    public String[] getSelGroups() {
        return selGroups;
    }

    public void setSelGroups(String[] selGroups) {
        this.selGroups = selGroups;
    }

    public String[] getSelRoles() {
        return selRoles;
    }

    public void setSelRoles(String[] selRoles) {
        this.selRoles = selRoles;
    }

    public String getLoginId() {
        return loginId;
    }

    public void setLoginId(String loginId) {
        this.loginId = loginId;
    }

    public Date getLastLoginTime() {
        return lastLoginTime;
    }

    public void setLastLoginTime(Date lastLoginTime) {
        this.lastLoginTime = lastLoginTime;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getSourceType() {
        return sourceType;
    }

    public void setSourceType(String sourceType) {
        this.sourceType = sourceType;
    }

    public Date getLastModifiedTime() {
        return lastModifiedTime;
    }

    public void setLastModifiedTime(Date lastModifiedTime) {
        this.lastModifiedTime = lastModifiedTime;
    }

    @Override
    public String toString() {
        return super.toString();
    }
}