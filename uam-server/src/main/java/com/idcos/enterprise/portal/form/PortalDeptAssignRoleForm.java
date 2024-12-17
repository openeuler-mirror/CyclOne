/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.form;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import org.hibernate.validator.constraints.NotBlank;
import org.springframework.web.bind.annotation.RequestBody;

import java.util.List;

/**
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月04 下午3:48 souakiragen Exp $
 */
@ApiModel(value = "为部门分配角色")
public class PortalDeptAssignRoleForm {

    @ApiModelProperty(name = "id", value = "部门id", required = true)
    @NotBlank(message = "部门id不能为空")
    private String id;

    @ApiModelProperty(name = "roleIds", value = "角色id列表", required = false)
    private List<String> roleIds;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public List<String> getRoleIds() {
        return roleIds;
    }

    public void setRoleIds(List<String> roleIds) {
        this.roleIds = roleIds;
    }
}
