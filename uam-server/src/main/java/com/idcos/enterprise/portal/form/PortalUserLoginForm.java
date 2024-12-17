
package com.idcos.enterprise.portal.form;

import com.idcos.cloud.biz.common.BaseForm;

/**
 * 用户登陆form请求
 *
 * @author yanlv
 * @version $Id: PortalUserLoginForm.java, v 0.1 2015年7月16日 上午9:52:04 yanlv Exp $
 */
public class PortalUserLoginForm extends BaseForm {
    private String username;
    private String password;
    private String rememberMe;
    private String source;

    public String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public String getRememberMe() {
        return rememberMe;
    }

    public void setRememberMe(String rememberMe) {
        this.rememberMe = rememberMe;
    }

    public String getSource() {
        return source;
    }

    public void setSource(String source) {
        this.source = source;
    }

}
