

package com.idcos.enterprise.portal.form;

// auto generated imports

import org.hibernate.validator.constraints.NotBlank;

import com.idcos.cloud.biz.common.BaseForm;

/**
 * 表单对象PortalResourceUpdateForm
 * <p>第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖</p>
 *
 * @author yanlv
 * @version PortalResourceUpdateForm.java, v 1.1 2015-06-09 10:00:37 yanlv Exp $
 */

public class PortalResourceUpdateForm extends BaseForm {

    //========== properties ==========
    /**
     * 权限资源类型
     */
    @NotBlank(message = "权限资源类型不能为空")
    private String code;

    /**
     * 权限资源名称
     */
    @NotBlank(message = "权限资源名称不能为空")
    private String name;

    /**
     * 权限资源URL
     */
    @NotBlank(message = "权限资源URL不能为空")
    private String url;

    /**
     * 系统名称
     */
    @NotBlank(message = "系统名称不能为空")
    private String appId;

    /**
     * 备注
     */
    private String remark;

    //========== getters and setters ==========

    public String getCode() {
        return code;
    }

    public void setCode(String code) {
        this.code = code;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getUrl() {
        return url;
    }

    public void setUrl(String url) {
        this.url = url;
    }

    public String getAppId() {
        return appId;
    }

    public void setAppId(String appId) {
        this.appId = appId;
    }

    public String getRemark() {
        return remark;
    }

    public void setRemark(String remark) {
        this.remark = remark;
    }

}