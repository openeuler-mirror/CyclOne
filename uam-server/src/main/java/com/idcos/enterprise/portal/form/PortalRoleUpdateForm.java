

package com.idcos.enterprise.portal.form;

// auto generated imports

import com.idcos.cloud.biz.common.BaseForm;
import org.hibernate.validator.constraints.NotBlank;


/**
 * 表单对象PortalRoleUpdateForm
 * <p>form由代码自动生成框架自动生成，不可进行编辑</p>
 *
 * @author
 * @version PortalRoleUpdateForm.java, v 1.1 2015-10-30 15:00:49  Exp $
 */

public class PortalRoleUpdateForm extends BaseForm {

    //========== properties ==========
    /**
     * 名称
     */
    @NotBlank(message = "名称不能为空")
    private String name;

    /**
     * 编码
     */
    @NotBlank(message = "编码不能为空")
    private String code;

    /**
     * 备注
     */
    private String remark;

    /**
     * 编辑类型
     */
    private String editType;


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


    public String getEditType() {
        return editType;
    }

    public void setEditType(String editType) {
        this.editType = editType;
    }


}