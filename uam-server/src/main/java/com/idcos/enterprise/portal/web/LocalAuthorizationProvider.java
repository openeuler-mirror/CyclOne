/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2016 All Rights Reserved.
 */
package com.idcos.enterprise.portal.web;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.alibaba.fastjson.JSON;
import com.idcos.common.web.model.BasicAuthorizationInfo;
import com.idcos.common.web.service.AuthorizationProvider;
import com.idcos.enterprise.portal.biz.common.PortalResponse;
import com.idcos.enterprise.portal.manager.auto.PortalRestfulService;

/**
 * @author zhouqin
 * @version com.idcos.automate.web.LocalAuthorizationProvider.java, v 1.1 5/13/16 zhouqin Exp $
 */
@Service
public class LocalAuthorizationProvider implements AuthorizationProvider {

    @Autowired
    private GlobalValue globalValue;

    @Autowired
    private PortalRestfulService portalRestfulService;

    @Override
    public BasicAuthorizationInfo authorityContent(String accessToken) {
        Claims claims = Jwts.parser().setSigningKey(globalValue.getSecretKey()).parseClaimsJws(accessToken).getBody();

        String userId = (String) claims.get("userId");

        PortalResponse portalResponse = portalRestfulService.queryAuthority(userId);

        BasicAuthorizationInfo basicAuthorizationInfo = JSON.parseObject(JSON.toJSONString(portalResponse.getContent()),
                BasicAuthorizationInfo.class);

        return basicAuthorizationInfo;
    }

}
