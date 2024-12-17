/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2018 All Rights Reserved.
 */
package com.idcos.enterprise.portal.web;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.idcos.common.web.model.JwtToken;
import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.filter.GenericFilterBean;

import javax.servlet.FilterChain;
import javax.servlet.ServletException;
import javax.servlet.ServletRequest;
import javax.servlet.ServletResponse;
import javax.servlet.http.Cookie;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.io.PrintWriter;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import static com.idcos.enterprise.portal.UamConstant.SLASH;

/**
 * @author Xizhao.Dai
 * @version JwtFilter.java, v1 2018/9/5 下午9:17 Xizhao.Dai Exp $$
 */
public class JwtFilter extends GenericFilterBean {
    private final static Logger LOGGER = LoggerFactory.getLogger(JwtFilter.class);

    private String secretKey;

    private Map loginUrl = new HashMap();

    private List<String> ignoreUris = new ArrayList<>();

    public JwtFilter(String ssoLoginUrl, String secretKey) {
        loginUrl.put("ssoWebUrl", ssoLoginUrl);
        this.secretKey = secretKey;
        LOGGER.debug("[JwtFilter]注入secretKey:" + secretKey);
    }

    @Override
    public void doFilter(ServletRequest servletRequest, ServletResponse servletResponse,
                         FilterChain filterChain) throws IOException, ServletException {

        final HttpServletRequest request = (HttpServletRequest) servletRequest;
        HttpServletResponse response = (HttpServletResponse) servletResponse;

        if (verifyIgnoreAddressAndUri(servletRequest, servletResponse, filterChain, request)) {
            return;
        }

        String authHeader = request.getHeader("access-token");
        String accessToken = null;
        if (StringUtils.isNotBlank(authHeader) && authHeader.startsWith("Bearer ")) {
            //优先判断header中的token
            accessToken = authHeader.substring(7);
            if (StringUtils.isNotBlank(accessToken) && !StringUtils.equals(accessToken, "null")) {
                LOGGER.info("从header中获取到accessToken:" + accessToken);
            } else {
                accessToken = null;
            }
        }
        if (StringUtils.isBlank(accessToken)) {
            //header中没有的情况下,从cookie获取
            Cookie[] cookies = request.getCookies();
            if (cookies == null || cookies.length == 0) {
                LOGGER.debug("没有cookie");
            } else {
                for (Cookie cookie : cookies) {
                    if ("access-token".equals(cookie.getName()) && StringUtils.isNotBlank(cookie.getValue())) {
                        accessToken = cookie.getValue();
                        LOGGER.info("从cookie中获取到accessToken:" + accessToken);
                        break;
                    }
                }
            }
        }
        if (StringUtils.isNotEmpty(accessToken)) {

            JwtToken jwtToken = new JwtToken(accessToken);
            if (jwtToken.verify(secretKey)) {
                SecurityContextHolder.getContext().setAuthentication(jwtToken);
            }
        } else {
            setLoginFailureStatus(response);
            return;
        }

        filterChain.doFilter(servletRequest, servletResponse);
    }

    private ObjectMapper mapper = new ObjectMapper();

    private void setLoginFailureStatus(HttpServletResponse response) throws IOException {
        response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
        response.setContentType("application/json;charset=UTF-8");
        RestResponse rsp = new RestResponse();
        rsp.setContent(loginUrl);
        rsp.setStatus("fail");
        rsp.setMessage("会话超时或没有登录,请重新登录。");
        rsp.setStatusCode(ResponseStatusCode.RSP_STATUS_LOGINOUT);
        PrintWriter writer = response.getWriter();
        writer.write(mapper.writeValueAsString(rsp));
        writer.close();
    }

    private boolean verifyIgnoreAddressAndUri(ServletRequest servletRequest, ServletResponse servletResponse,
                                              FilterChain filterChain,
                                              HttpServletRequest request) throws IOException, ServletException {
        String uri = request.getRequestURI();
        if ("".equals(uri.trim()) || SLASH.equals(uri.trim())) {
            filterChain.doFilter(servletRequest, servletResponse);
            return true;
        }

        for (String ignoreUri : ignoreUris) {
            if (uri.startsWith(ignoreUri)) {
                filterChain.doFilter(servletRequest, servletResponse);
                return true;
            }
            if (uri.matches(".*\\.js$") || uri.matches(".*\\.css$") || uri.matches(".*\\.jpg$")
                    || uri.matches(".*\\.png$") || uri.matches(".*\\.svg$") || uri.matches(".*\\.jpeg$") || uri.matches(".*\\.json$")
                    || uri.matches(".*\\.txt$") || uri.matches(".*\\.ts$") || uri.matches(".*\\.ico$")
                    || uri.matches(".*\\.ttf$") || uri.matches(".*\\.woff$")) {
                filterChain.doFilter(servletRequest, servletResponse);
                return true;
            }
        }
        return false;
    }

    public void setIgnoreUris(List<String> ignoreUris) {
        if (ignoreUris != null && ignoreUris.size() >= 0) {
            this.ignoreUris.addAll(ignoreUris);
        }
    }
}