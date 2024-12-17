/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.form;

import io.swagger.annotations.ApiModelProperty;
import org.hibernate.validator.constraints.NotBlank;

/**
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月07 上午9:14 souakiragen Exp $
 */
public class PortalTenantAddForm {

    @ApiModelProperty(name = "tenantId", value = "租户id", notes = "租户id", required = true)
    @NotBlank(message = "租户id不能为空")
    private String tenantId;
    @NotBlank(message = "租户的名称不能为空")
    @ApiModelProperty(name = "tenantName", value = "租户的名称", notes = "租户的名称", required = true)
    private String tenantName;

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
}
