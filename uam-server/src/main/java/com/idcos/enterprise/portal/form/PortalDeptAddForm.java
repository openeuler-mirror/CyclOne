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
 * @version $Id: , v 0.1 2017年11月04 上午11:20 souakiragen Exp $
 */
@ApiModel
public class PortalDeptAddForm {

    @NotBlank(message = "部门编码不能为空")
    @ApiModelProperty(name = "code", value = "部门编码", notes = "部门编码,不能为空", dataType = "String", required = true)
    private String code;

    @NotBlank(message = "部门名称不能为空")
    @ApiModelProperty(name = "displayName", value = "部门名称", notes = "部门名称,不能为空", dataType = "String", required = true)
    private String displayName;

    @ApiModelProperty(name = "parentId", value = "父级部门ID", notes = "属于哪个部门", dataType = "String", required = false)
    private String parentId;

    @ApiModelProperty(name = "remark", value = "备注", notes = "备注", dataType = "String", required = false)
    private String remark;

    @NotBlank(message = "租户不能为空")
    @ApiModelProperty(name = "tenantId", value = "租户id", notes = "租户id,不能为空", dataType = "String", required = true)
    private String tenantId;

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
