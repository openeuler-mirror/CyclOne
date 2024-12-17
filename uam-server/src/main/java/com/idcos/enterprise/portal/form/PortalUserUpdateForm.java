/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.form;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import org.apache.commons.lang3.StringUtils;
import org.hibernate.validator.constraints.NotBlank;

/**
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月03 下午8:02 souakiragen Exp $
 */
@ApiModel("更新用户")
public class PortalUserUpdateForm {
    /**
     * id。
     */
    @ApiModelProperty(name = "id", value = "员工Id", notes = "员工的Id", required = true)
    @NotBlank(message = "Id不能为空")
    private String id;
    /**
     * 登录名。
     */
    @ApiModelProperty(name = "loginId", value = "员工登录名", notes = "员工的登录名，不能重复", required = true)
    @NotBlank(message = "登录名不能为空")
    private String loginId;

    /**
     * 中文名。
     */
    @NotBlank(message = "用户名不能为空")
    @ApiModelProperty(name = "loginId", value = "员工中文名", notes = "中文名称，可能会重复。", required = true)
    private String name;


    /**
     * 部门ID。
     */
    @ApiModelProperty(name = "deptId", value = "部门ID", notes = "部门内部ID。", required = false)
    private String deptId;

    /**
     * 状态。
     */
    @ApiModelProperty(name = "status", value = "状态", notes = "状态。", required = false)
    private String status;

    /**
     * 职务。
     */
    @ApiModelProperty(name = "title", value = "职务", notes = "职务", required = false)
    private String title;

    /**
     * 移动电话
     */
    @ApiModelProperty(name = "mobile1", value = "移动电话1", notes = "移动电话1", required = false)
    private String mobile1;

    /**
     * 移动电话
     */
    @ApiModelProperty(name = "mobile2", value = "移动电话2", notes = "移动电话2", required = false)
    private String mobile2;

    /**
     * rtx
     */
    @ApiModelProperty(name = "rtx", value = "RTX", notes = "RTX", required = false)
    private String rtx;

    /**
     * 办公电话
     */
    @ApiModelProperty(name = "officeTel1", value = "办公电话1", notes = "办公电话1", required = false)
    private String officeTel1;

    /**
     * 办公电话
     */
    @ApiModelProperty(name = "officeTel2", value = "办公电话2", notes = "办公电话2", required = false)
    private String officeTel2;

    /**
     * 员工类型
     */
    @ApiModelProperty(name = "employeeType", value = "员工类型", notes = "员工类型", required = false)
    private String employeeType;
    /**
     * 微信。
     */
    @ApiModelProperty(name = "weixin", value = "微信号", notes = "微信号", required = false)
    private String weixin;

    /**
     * 邮箱。
     */
    @ApiModelProperty(name = "email", value = "邮箱", notes = "电子邮箱地址。", required = false)
    private String email;

    /**
     * 备注
     */
    @ApiModelProperty(name = "remark", value = "备注", notes = "备注", required = false)
    private String remark = "";

    /**
     * 租户
     */
    @NotBlank(message = "租户名不能为空")
    @ApiModelProperty(name = "tenantId", value = "租户Id", notes = "租户Id", required = true)
    private String tenantId;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getLoginId() {
        return loginId;
    }

    /**
     * 去除首尾空格
     *
     * @param loginId
     */
    public void setLoginId(String loginId) {
        this.loginId = StringUtils.trim(loginId);
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

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
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

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
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
}
