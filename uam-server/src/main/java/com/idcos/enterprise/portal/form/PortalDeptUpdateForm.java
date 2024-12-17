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
 * @version $Id: , v 0.1 2017年11月04 下午2:30 souakiragen Exp $
 */
@ApiModel(value = "修改部门")
public class PortalDeptUpdateForm extends PortalDeptAddForm {

    @ApiModelProperty(name = "id", value = "部门id", notes = "要修改的部门id", required = true)
    @NotBlank(message = "部门id不能为空")
    private String id;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }
}
