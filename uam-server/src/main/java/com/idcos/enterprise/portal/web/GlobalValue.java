/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2018 All Rights Reserved.
 */
package com.idcos.enterprise.portal.web;

import com.alibaba.fastjson.JSONObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.annotation.PostConstruct;
import java.util.List;

/**
 * 统一管理使用配置文件的配置项
 *
 * @author Xizhao.Dai
 * @version GlobalValue.java, v1 2018/9/6 上午1:17 Xizhao.Dai Exp $$
 */
public class GlobalValue {
    private static final Logger logger = LoggerFactory.getLogger(GlobalValue.class);

    private String ssoLoginUrl;

    private String secretKey;

    private String appId;

    private List<String> ignoreUris;

    private String decryptKey;

    private String accessTimeout;

    private Boolean isMultiTenant;

    public GlobalValue(String ssoLoginUrl, String secretKey, String appId, List<String> ignoreUris, String decryptKey,
                       String accessTimeout, Boolean isMultiTenant) {
        this.ssoLoginUrl = ssoLoginUrl;
        this.secretKey = secretKey;
        this.appId = appId;
        this.ignoreUris = ignoreUris;
        this.decryptKey = decryptKey;
        this.accessTimeout = accessTimeout;
        this.isMultiTenant = isMultiTenant;
    }

    @PostConstruct
    public void afterInit() {
        logger.info("========== UAM_GlobalValue ==========");
        logger.info(JSONObject.toJSONString(this));
    }

    public String getSsoLoginUrl() {
        return ssoLoginUrl;
    }

    public String getSecretKey() {
        return secretKey;
    }

    public String getAppId() {
        return appId;
    }

    public List<String> getIgnoreUris() {
        return ignoreUris;
    }

    public String getDecryptKey() {
        return decryptKey;
    }

    public String getAccessTimeout() {
        return accessTimeout;
    }

    public Boolean getMultiTenant() {
        return isMultiTenant;
    }
}