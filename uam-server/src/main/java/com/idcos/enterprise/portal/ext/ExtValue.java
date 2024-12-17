/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2018 All Rights Reserved.
 */
package com.idcos.enterprise.portal.ext;

import javax.annotation.PostConstruct;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.alibaba.fastjson.JSONObject;


/**
 * 统一管理对接外部系统的配置项
 *
 * @author Xizhao.Dai
 * @version ExtValue.java, v1 2018/9/6 上午1:17 Xizhao.Dai Exp $$
 */
public class ExtValue {
    private static final Logger logger = LoggerFactory.getLogger(ExtValue.class);

    private String cron;

    private String domain;

    private String ip;

    private String port;

    private String baseDn;

    private String basePassword;

    private String baseDept;

    private String extAppId;

    private String extAppToken;

    public ExtValue(String cron, String domain, String ip, String port,
        String baseDn, String basePassword, String baseDept, String extAppId, String extAppToken) {
        this.cron = cron;
        this.domain = domain;
        this.ip = ip;
        this.port = port;
        this.baseDn = baseDn;
        this.basePassword = basePassword;
        this.baseDept = baseDept;
        this.extAppId = extAppId;
        this.extAppToken = extAppToken;
    }

    @PostConstruct
    public void afterInit() {
        logger.info("========== ExtValue ==========");
        logger.info(JSONObject.toJSONString(this));
    }

    public String getCron() {
        return cron;
    }

    public String getDomain() {
        return domain;
    }

    public String getIp() {
        return ip;
    }

    public String getPort() {
        return port;
    }

    public String getBaseDn() {
        return baseDn;
    }

    public String getBasePassword() {
        return basePassword;
    }

    public String getBaseDept() {
        return baseDept;
    }

    public String getExtAppId() {
        return extAppId;
    }

    public String getExtAppToken() {
        return extAppToken;
    }

}