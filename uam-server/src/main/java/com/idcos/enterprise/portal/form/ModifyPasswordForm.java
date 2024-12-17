/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.form;

import org.hibernate.validator.constraints.NotBlank;

/**
 * @author Dana
 * @version ModifyPasswordForm.java, v1 2017/12/7 下午1:55 Dana Exp $$
 */
public class ModifyPasswordForm {
    /**
     * 用户ID.
     */
    @NotBlank(message = "用户Id")
    private String userId;

    /**
     * 用户老的密码
     */
    @NotBlank(message = "原始密码")
    private String oldPassword;

    /**
     * 用户新的密码。
     */
    @NotBlank(message = "新密码")
    private String newPassword;

    public String getUserId() {
        return userId;
    }

    public void setUserId(String userId) {
        this.userId = userId;
    }

    public String getOldPassword() {
        return oldPassword;
    }

    public void setOldPassword(String oldPassword) {
        this.oldPassword = oldPassword;
    }

    public String getNewPassword() {
        return newPassword;
    }

    public void setNewPassword(String newPassword) {
        this.newPassword = newPassword;
    }
}