/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.form;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import org.hibernate.validator.constraints.NotBlank;

/**
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月02 下午5:36 souakiragen Exp $
 */
@ApiModel("分页查询")
public class PortalUserQueryPageListForm {

    /**
     * 租户id
     */
    @ApiModelProperty(name = "tenantId", value = "租户id", notes = "租户id", required = true)
    @NotBlank(message = "租户不能为空")
    private String tenantId;

    /**
     * 部门id
     */
    @ApiModelProperty(name = "deptId", value = "部门id", notes = "部门id")
    private String deptId;

    /**
     * 用户名称
     */
    @ApiModelProperty(name = "name", value = "用户名称", notes = "用户名称，查询条件")
    private String name;

    /**
     * 页码
     */
    @ApiModelProperty(name = "pageNo", value = "页码", notes = "第几页")
    private int pageNo;

    /**
     * 分页大小
     */
    @ApiModelProperty(name = "pageSize", value = "分页大小", notes = "每页的记录数")
    private int pageSize;

    public String getTenantId() {
        return tenantId;
    }

    public void setTenantId(String tenantId) {
        this.tenantId = tenantId;
    }

    public String getDeptId() {
        return deptId;
    }

    public void setDeptId(String deptId) {
        this.deptId = deptId;
    }

    public int getPageNo() {
        return pageNo;
    }

    public void setPageNo(int pageNo) {
        this.pageNo = pageNo;
    }

    public int getPageSize() {
        return pageSize;
    }

    public void setPageSize(int pageSize) {
        this.pageSize = pageSize;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }
}
