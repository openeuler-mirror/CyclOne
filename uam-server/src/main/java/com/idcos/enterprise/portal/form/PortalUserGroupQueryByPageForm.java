

package com.idcos.enterprise.portal.form;

// auto generated imports

import com.idcos.cloud.biz.common.BaseForm;
import com.idcos.cloud.core.dal.common.query.DataQueryField;
import com.idcos.cloud.core.dal.common.query.OperatorEnum;

/**
 * 表单对象PortalUserGroupQueryByPageForm
 * <p>form由代码自动生成框架自动生成，不可进行编辑</p>
 *
 * @author
 * @version PortalUserGroupQueryByPageForm.java, v 1.1 2015-10-28 14:17:42  Exp $
 */

public class PortalUserGroupQueryByPageForm extends BaseForm {

    //========== properties ==========
    /**
     * 名称
     */
    @DataQueryField(operator = OperatorEnum.LIKE)
    private String name;

    /**
     * 类型
     */
    @DataQueryField
    private String type;

    /**
     * 租户
     */
    @DataQueryField(name = "tenant")
    private String tenantId;

    /**
     * 备注
     */
    @DataQueryField
    private String remark;

    /**
     * 所属角色
     */
    @DataQueryField
    private String selectRoles;

    /**
     * 包含用户
     */
    @DataQueryField
    private String selectUsers;

    //========== getters and setters ==========


    public String getTenantId() {
        return tenantId;
    }

    public void setTenantId(String tenantId) {
        this.tenantId = tenantId;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getRemark() {
        return remark;
    }

    public void setRemark(String remark) {
        this.remark = remark;
    }

    public String getSelectRoles() {
        return selectRoles;
    }

    public void setSelectRoles(String selectRoles) {
        this.selectRoles = selectRoles;
    }

    public String getSelectUsers() {
        return selectUsers;
    }

    public void setSelectUsers(String selectUsers) {
        this.selectUsers = selectUsers;
    }

}