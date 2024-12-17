package com.idcos.enterprise.sso.service.impl;

import static com.idcos.enterprise.portal.dal.enums.LoginType.DEFAULT;
import static com.idcos.enterprise.portal.dal.enums.LoginType.WEBANK;

import java.util.*;

import javax.servlet.http.HttpServletResponse;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryCallback;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryTemplate;
import com.idcos.enterprise.portal.dal.entity.PortalUser;
import com.idcos.enterprise.portal.dal.enums.LoginType;
import com.idcos.enterprise.portal.dal.enums.SourceTypeEnum;
import com.idcos.enterprise.portal.form.LoginForm;
import com.idcos.enterprise.sso.CheckLoginService;
import com.idcos.enterprise.sso.service.LoginResult;
import com.idcos.enterprise.sso.service.LoginService;

/**
 * @Author: Dai
 * @Date: 2018/10/8 下午5:25
 * @Description:
 */
@Service
public class LoginServiceImpl {
    private BusinessQueryTemplate businessQueryTemplate;
    private CheckLoginService checkLoginService;

    private final Map<LoginType, LoginService> loginServiceMap = new HashMap<>();

    @Autowired
    public LoginServiceImpl(BusinessQueryTemplate businessQueryTemplate, CheckLoginService checkLoginService,
        List<LoginService> loginServices) {
        this.businessQueryTemplate = businessQueryTemplate;
        this.checkLoginService = checkLoginService;
        for (LoginService loginService : loginServices) {
            loginServiceMap.put(loginService.action(), loginService);
        }
    }

    public CommonResult<?> login(final LoginForm loginForm) {
        return businessQueryTemplate.process(new BusinessQueryCallback<LoginResult>() {
            @Override
            public LoginResult doQuery() {
                return loginServiceMap.get(DEFAULT).login(loginForm);
            }

            @Override
            public void checkParam() {}
        });
    }

    public CommonResult<?> afterLogin(final HttpServletResponse response, final String customer) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Void>() {

            @Override
            public Void doQuery() {
                return null;
            }

            @Override
            public void checkParam() {}
        });
    }
}
