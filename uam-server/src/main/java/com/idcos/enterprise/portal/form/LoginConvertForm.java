

package com.idcos.enterprise.portal.form;

// auto generated imports

import com.idcos.cloud.biz.common.BaseForm;
import org.hibernate.validator.constraints.NotBlank;

/**
 * 表单对象LoginConvertForm
 * <p>第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖</p>
 *
 * @author yanlv
 * @version LoginConvertForm.java, v 1.1 2015-06-09 10:00:37 yanlv Exp $
 */

public class LoginConvertForm extends BaseForm {

    //========== properties ==========

    private String userInfo;

    private String password;

    private String tenant;


    //========== getters and setters ==========


    public String getUserInfo() {
        return userInfo;
    }

    public void setUserInfo(String userInfo) {
        this.userInfo = userInfo;
    }

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public String getTenant() {
        return tenant;
    }

    public void setTenant(String tenant) {
        this.tenant = tenant;
    }
}