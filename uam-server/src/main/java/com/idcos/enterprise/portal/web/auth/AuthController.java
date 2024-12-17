/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2016 All Rights Reserved.
 */
package com.idcos.enterprise.portal.web.auth;

import com.alibaba.fastjson.JSONObject;
import com.idcos.common.service.vo.CommonRestResult;
import com.idcos.common.web.model.JwtToken;
import com.idcos.enterprise.portal.web.GlobalValue;
import io.swagger.annotations.Api;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.authentication.AnonymousAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

import java.util.Map;

/**
 * @author zhouqin
 * @version com.idcos.automate.web.controller.auth.AuthController.java, v 1.1 2/29/16 zhouqin Exp $
 */
@RestController
@RequestMapping("/auth")
@Api(tags = "02.认证接口", description = "AuthController")
public class AuthController {

    @Autowired
    private GlobalValue globalValue;

    @RequestMapping(method = RequestMethod.GET)
    public CommonRestResult<?> currentUser() {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        if (!authentication.isAuthenticated() || authentication instanceof AnonymousAuthenticationToken) {
            CommonRestResult<Map> userResult = new CommonRestResult<>();
            userResult.setStatus("AUTH_FAILED");
            userResult.setMessage("未登录状态");
            return userResult;
        }
        JwtToken jwtToken = (JwtToken) authentication;
        return new CommonRestResult<>(jwtToken);
    }

    /**
     * 判断是否是多租户
     *
     * @return
     */
    @RequestMapping(value = "/isMultiTenant", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    public JSONObject isMultiTenant() {
        JSONObject json = new JSONObject();
        json.put("status", "success");
        json.put("message", "isMultiTenant");
        json.put("content", globalValue.getMultiTenant());
        return json;
    }

}
