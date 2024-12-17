/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.vo;

import java.io.Serializable;
import java.util.List;
import java.util.Map;

/**
 * @author Dana
 * @version AuthInfoVO.java, v1 2017/11/18 下午6:22 Dana Exp $$
 */
public class AuthInfoVO implements Serializable {
    private static final long serialVersionUID = 1L;

    private String id;

    private String loginId;

    private String name;

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

    private String deptId;

    private String deptName;

    private String tenantId;

    private String tenantName;

    private Map<String, List<String>> permissions;

    private List<String> userGroups;

    private List<String> roleIds;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getLoginId() {
        return loginId;
    }

    public void setLoginId(String loginId) {
        this.loginId = loginId;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
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

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
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

    public Map<String, List<String>> getPermissions() {
        return permissions;
    }

    public void setPermissions(Map<String, List<String>> permissions) {
        this.permissions = permissions;
    }

    public List<String> getUserGroups() {
        return userGroups;
    }

    public void setUserGroups(List<String> userGroups) {
        this.userGroups = userGroups;
    }

    public List<String> getRoleIds() {
        return roleIds;
    }

    public void setRoleIds(List<String> roleIds) {
        this.roleIds = roleIds;
    }

    @Override
    public String toString() {
        return "AuthInfoVO{" + "id='" + id + '\'' + ", loginId='" + loginId + '\'' + ", name='" + name + '\''
                + ", email='" + email + '\'' + ", mobile1='" + mobile1 + '\'' + ", mobile2='" + mobile2 + '\''
                + ", officeTel1='" + officeTel1 + '\'' + ", officeTel2='" + officeTel2 + '\'' + ", employeeType='"
                + employeeType + '\'' + ", title='" + title + '\'' + ", weixin='" + weixin + '\'' + ", remark='" + remark
                + '\'' + ", isActive='" + isActive + '\'' + ", status='" + status + '\'' + ", deptId='" + deptId + '\''
                + ", deptName='" + deptName + '\'' + ", tenantId='" + tenantId + '\'' + ", tenantName='" + tenantName
                + '\'' + ", permissions=" + permissions + ", userGroups=" + userGroups + ", roleIds=" + roleIds + '}';
    }
}