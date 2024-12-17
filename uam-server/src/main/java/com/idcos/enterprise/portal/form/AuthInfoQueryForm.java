/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.form;

import org.hibernate.validator.constraints.NotBlank;

/**
 * @author Dana
 * @version getAuthInfoForm.java, v1 2018/1/10 下午5:35 Dana Exp $$
 */
public class AuthInfoQueryForm {
    /**
     * 登录Id.
     */
    @NotBlank(message = "登录Id")
    private String loginId;

    /**
     * 用户密码
     */
    @NotBlank(message = "密码")
    private String password;

    /**
     * 租户。
     */
    @NotBlank(message = "租户")
    private String tenantId;

    /**
     * 系统名。
     */
    @NotBlank(message = "系统名")
    private String appId;

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

    public String getAppId() {
        return appId;
    }

    public void setAppId(String appId) {
        this.appId = appId;
    }
}