

package com.idcos.enterprise.portal.form;

// auto generated imports

import com.idcos.cloud.biz.common.BaseForm;
import com.idcos.cloud.core.dal.common.query.DataQueryField;
import com.idcos.cloud.core.dal.common.query.OperatorEnum;

/**
 * 表单对象PortalRoleQueryByPageForm
 * <p>form由代码自动生成框架自动生成，不可进行编辑</p>
 *
 * @author
 * @version PortalRoleQueryByPageForm.java, v 1.1 2015-10-30 15:00:49  Exp $
 */

public class PortalRoleQueryByPageForm extends BaseForm {

    //========== properties ==========
    /**
     * 名称
     */
    @DataQueryField(operator = OperatorEnum.LIKE)
    private String name;

    /**
     * 编码
     */
    @DataQueryField
    private String code;

    /**
     * 备注
     */
    @DataQueryField
    private String remark;

    //========== getters and setters ==========

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getCode() {
        return code;
    }

    public void setCode(String code) {
        this.code = code;
    }

    public String getRemark() {
        return remark;
    }

    public void setRemark(String remark) {
        this.remark = remark;
    }

}