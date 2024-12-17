package com.idcos.enterprise.portal.biz.common.utils;

import com.idcos.common.web.model.JwtToken;
import com.idcos.enterprise.portal.dal.entity.JwtUser;
import com.idcos.enterprise.portal.web.GlobalValue;
import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Service;

/**
 * @author GuanBin
 * @version CurrentUser.java, v1 2017/9/25 下午2:52 GuanBin Exp $$
 */
@Service
public class CurrentUser {

    @Autowired
    private GlobalValue globalValue;

    private static final Logger LOGGER = LoggerFactory.getLogger(CurrentUser.class);

    public JwtUser getUser() {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        JwtToken jwtToken = (JwtToken) authentication;
        JwtUser jwtUser = null;

        try {
            Claims claims = (Claims) Jwts.parser().setSigningKey(globalValue.getSecretKey())
                    .parseClaimsJws(jwtToken.getToken()).getBody();
            jwtUser = new JwtUser();
            jwtUser.setExp(claims.get("exp") == null ? "" : claims.get("exp").toString());
            jwtUser.setLoginId(claims.get("loginId") == null ? "" : claims.get("loginId").toString());
            jwtUser.setName(claims.get("name") == null ? "" : claims.get("name").toString());
            jwtUser.setUserId(claims.get("userId") == null ? "" : claims.get("userId").toString());
            jwtUser.setUserName(claims.get("userName") == null ? "" : claims.get("userName").toString());
            jwtUser.setTimeout(claims.get("timeout") == null ? "" : claims.get("timeout").toString());
            jwtUser.setTenantId(claims.get("tenantId") == null ? "" : claims.get("tenantId").toString());
            jwtUser.setTenantName(claims.get("tenantName") == null ? "" : claims.get("tenantName").toString());
        } catch (Exception e) {
            LOGGER.error("用户解析异常:" + e);
        }

        return jwtUser;
    }
}
