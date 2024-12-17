/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.form;

import org.hibernate.validator.constraints.NotBlank;

/**
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月07 上午10:38 souakiragen Exp $
 */
public class PortalTenantUpdateForm {

    @NotBlank(message = "id不能为空")
    private String id;
    @NotBlank(message = "tenantId不能为空")
    private String tenantId;
    @NotBlank(message = "租户名称不能为空")
    private String tenantName;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
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
}
