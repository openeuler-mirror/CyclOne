

package com.idcos.enterprise.portal.form;

import org.hibernate.validator.constraints.NotBlank;

// auto generated imports

import com.idcos.cloud.biz.common.BaseForm;

/**
 * 表单对象LoginForm
 * <p>第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖</p>
 *
 * @author yanlv
 * @version LoginForm.java, v 1.1 2015-06-09 10:00:37 yanlv Exp $
 */

public class LoginForm extends BaseForm {

    //========== properties ==========
    /**
     * loginId
     */
    @NotBlank(message = "登录id")
    private String loginId;

    /**
     * password
     */
    @NotBlank(message = "密码")
    private String password;


    /**
     * 租户id
     */
    @NotBlank(message = "tenant")
    private String tenantId;


    //========== getters and setters ==========


    public String getLoginId() {
        return loginId;
    }

    public void setLoginId(String loginId) {
        this.loginId = loginId;
    }

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public String getTenantId() {
        return tenantId;
    }

    public void setTenantId(String tenantId) {
        this.tenantId = tenantId;
    }
}