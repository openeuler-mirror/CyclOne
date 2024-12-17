/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.form;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import org.hibernate.validator.constraints.NotBlank;

/**
 * @author Dana
 * @version PortalGroupAndUserCreateForm.java, v1 2017/12/25 上午11:43 Dana Exp $$
 */
@ApiModel(value = "新建用户组(包括用户)")
public class PortalGroupAndUserCreateForm {

    @ApiModelProperty(name = "groupName", value = "用户组名", required = true)
    @NotBlank(message = "用户组名")
    private String groupName;

    @ApiModelProperty(name = "groupType", value = "用户组类型")
    private String groupType;

    @ApiModelProperty(name = "groupRemark", value = "用户组备注")
    private String groupRemark;

    @ApiModelProperty(name = "tenantId", value = "租户", required = true)
    @NotBlank(message = "租户")
    private String tenantId;

    @ApiModelProperty(name = "loginIds", value = "用户登录id列表")
    private String loginIds;

    public String getGroupName() {
        return groupName;
    }

    public void setGroupName(String groupName) {
        this.groupName = groupName;
    }

    public String getGroupType() {
        return groupType;
    }

    public void setGroupType(String groupType) {
        this.groupType = groupType;
    }

    public String getGroupRemark() {
        return groupRemark;
    }

    public void setGroupRemark(String groupRemark) {
        this.groupRemark = groupRemark;
    }

    public String getTenantId() {
        return tenantId;
    }

    public void setTenantId(String tenantId) {
        this.tenantId = tenantId;
    }

    public String getLoginIds() {
        return loginIds;
    }

    public void setLoginIds(String loginIds) {
        this.loginIds = loginIds;
    }
}
