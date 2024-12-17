package com.idcos.enterprise.portal.form;

import org.hibernate.validator.constraints.NotBlank;

import com.idcos.cloud.biz.common.BaseForm;

/**
 * 用户表单提交类
 *
 * @author jiaohuizhe
 * @version $Id: PortalUserForm.java, v 0.1 2015年5月8日 上午11:03:29 jiaohuizhe Exp $
 */
public class PortalUserAllocateGroupForm extends BaseForm {

    @NotBlank(message = "用户id不能为空")
    private String id;

    private String selGroups;

    private String tenantId;

    public String getTenantId() {
        return tenantId;
    }

    public void setTenantId(String tenantId) {
        this.tenantId = tenantId;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getSelGroups() {
        return selGroups;
    }

    public void setSelGroups(String selGroups) {
        this.selGroups = selGroups;
    }

}