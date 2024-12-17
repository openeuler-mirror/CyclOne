

package com.idcos.enterprise.portal.form;

// auto generated imports

import com.idcos.cloud.biz.common.BaseForm;
import org.hibernate.validator.constraints.NotBlank;


/**
 * 表单对象PortalPermissionQueryAuthResForm
 * <p>第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖</p>
 *
 * @author yanlv
 * @version PortalPermissionQueryAuthResForm.java, v 1.1 2015-06-09 10:00:37 yanlv Exp $
 */

public class PortalPermissionQueryAuthResForm extends BaseForm {

    //========== properties ==========
    /**
     * 授权资源类型
     */
    @NotBlank(message = "授权资源类型不能为空")
    private String authResType;

    /**
     * 授权对象类型
     */
    @NotBlank(message = "授权对象类型不能为空")
    private String authObjType;

    /**
     * 授权对象ID
     */
    @NotBlank(message = "授权对象ID不能为空")
    private String authObjId;


    //========== getters and setters ==========

    public String getAuthResType() {
        return authResType;
    }

    public void setAuthResType(String authResType) {
        this.authResType = authResType;
    }


    public String getAuthObjType() {
        return authObjType;
    }

    public void setAuthObjType(String authObjType) {
        this.authObjType = authObjType;
    }


    public String getAuthObjId() {
        return authObjId;
    }

    public void setAuthObjId(String authObjId) {
        this.authObjId = authObjId;
    }


}