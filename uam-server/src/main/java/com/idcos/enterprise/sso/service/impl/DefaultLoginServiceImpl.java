package com.idcos.enterprise.sso.service.impl;

import com.idcos.enterprise.portal.dal.entity.PortalUser;
import com.idcos.enterprise.portal.dal.enums.LoginType;
import com.idcos.enterprise.portal.form.LoginForm;
import com.idcos.enterprise.sso.CheckLoginService;
import com.idcos.enterprise.sso.service.LoginResult;
import com.idcos.enterprise.sso.service.LoginService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

/**
 * @Author: Dai
 * @Date: 2018/10/28 4:02 AM
 * @Description:
 */
@Service
public class DefaultLoginServiceImpl implements LoginService {
    private static final Logger logger = LoggerFactory.getLogger(DefaultLoginServiceImpl.class);
    @Autowired
    private CheckLoginService checkLoginService;

    @Override
    public LoginType action() {
        return LoginType.DEFAULT;
    }

    /**
     * 默认登录逻辑
     *
     * @param form
     * @return
     */
    @Override
    public LoginResult login(LoginForm form) {
        PortalUser portalUser = checkLoginService.checkUser(form.getLoginId(), form.getTenantId());
        checkLoginService.checkPassword(portalUser, form.getPassword());
        logger.info("===============用户" + form.getLoginId() + "@" + form.getTenantId() + "开始登入rbac系统================");
        LoginResult result = new LoginResult();
        result.setStatus(true);
        return result;
    }

    @Override
    public String afterLogin(String customer) {
        return null;
    }
}
