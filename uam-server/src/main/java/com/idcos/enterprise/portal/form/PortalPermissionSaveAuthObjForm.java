

package com.idcos.enterprise.portal.form;

// auto generated imports

import org.hibernate.validator.constraints.NotBlank;

import com.idcos.cloud.biz.common.BaseForm;

/**
 * 表单对象PortalPermissionSaveAuthObjForm
 * <p>第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖</p>
 *
 * @author yanlv
 * @version PortalPermissionSaveAuthObjForm.java, v 1.1 2015-06-09 10:00:37 yanlv Exp $
 */

public class PortalPermissionSaveAuthObjForm extends BaseForm {

    //========== properties ==========
    /**
     * 授权资源类型
     */
    @NotBlank(message = "授权资源类型不能为空")
    private String authResType;

    /**
     * 授权资源ID
     */
    @NotBlank(message = "授权资源ID不能为空")
    private String authResId;

    /**
     * 授权资源ID
     */
    @NotBlank(message = "授权资源名称不能为空")
    private String authResName;

    /**
     * 授权对象类型
     */
    @NotBlank(message = "授权对象类型不能为空")
    private String authObjType;

    /**
     * null
     */
    private String authObjIds;

    //========== getters and setters ==========

    public String getAuthResType() {
        return authResType;
    }

    public void setAuthResType(String authResType) {
        this.authResType = authResType;
    }

    public String getAuthResName() {
        return authResName;
    }

    public void setAuthResName(String authResName) {
        this.authResName = authResName;
    }

    public String getAuthResId() {
        return authResId;
    }

    public void setAuthResId(String authResId) {
        this.authResId = authResId;
    }

    public String getAuthObjType() {
        return authObjType;
    }

    public void setAuthObjType(String authObjType) {
        this.authObjType = authObjType;
    }

    public String getAuthObjIds() {
        return authObjIds;
    }

    public void setAuthObjIds(String authObjIds) {
        this.authObjIds = authObjIds;
    }

}