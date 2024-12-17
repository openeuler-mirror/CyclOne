package com.idcos.enterprise.sso.service;

import com.idcos.enterprise.portal.dal.enums.LoginType;
import com.idcos.enterprise.portal.form.LoginForm;

/**
 * @author xizhao
 */
public interface LoginService {

    /**
     * action
     *
     * @return LoginType
     */
    LoginType action();

    /**
     * login
     *
     * @param form
     * @return
     */
    LoginResult login(LoginForm form);

    /**
     * afterLogin
     *
     * @param customer
     * @return
     */
    String afterLogin(String customer);
}
