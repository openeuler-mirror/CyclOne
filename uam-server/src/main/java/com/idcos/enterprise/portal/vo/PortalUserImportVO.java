/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.vo;

import com.idcos.cloud.core.common.BaseVO;

/**
 * @author Dana
 * @version PortalUserImportVO.java, v1 2017/12/22 上午9:45 Dana Exp $$
 */
public class PortalUserImportVO extends BaseVO {
    private static final long serialVersionUID = 1L;

    private String id;

    private String loginId;

    private String name;

    private String deptId;

    private String deptFullName;

    private String title;

    private String password;

    private String weixin;

    private String email;

    private String mobile1;

    private String mobile2;

    private String officeTel1;

    private String officeTel2;

    private String employeeType;

    private String tenantId;

    private String sourceType;

    private String remark;

    private String dataStatus;

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

    public String getDeptId() {
        return deptId;
    }

    public void setDeptId(String deptId) {
        this.deptId = deptId;
    }

    public String getDeptFullName() {
        return deptFullName;
    }

    public void setDeptFullName(String deptFullName) {
        this.deptFullName = deptFullName;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public String getWeixin() {
        return weixin;
    }

    public void setWeixin(String weixin) {
        this.weixin = weixin;
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

    public String getTenantId() {
        return tenantId;
    }

    public void setTenantId(String tenantId) {
        this.tenantId = tenantId;
    }

    public String getSourceType() {
        return sourceType;
    }

    public void setSourceType(String sourceType) {
        this.sourceType = sourceType;
    }

    public String getRemark() {
        return remark;
    }

    public void setRemark(String remark) {
        this.remark = remark;
    }

    public String getDataStatus() {
        return dataStatus;
    }

    public void setDataStatus(String dataStatus) {
        this.dataStatus = dataStatus;
    }

    @Override
    public String toString() {
        return super.toString();
    }
}