package com.idcos;

import java.util.List;

import org.springframework.beans.factory.annotation.Autowire;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.web.support.SpringBootServletInitializer;
import org.springframework.context.annotation.Bean;
import org.springframework.core.env.Environment;

import com.idcos.enterprise.portal.ext.ExtValue;
import com.idcos.enterprise.portal.web.GlobalValue;

/**
 * Uam启动类
 *
 * @author pengganyu
 * @version $Id: UamAppliction.java, v 0.1 2016年5月8日 下午12:14:16 pengganyu Exp $
 */
@SpringBootApplication
public class UamAppliction extends SpringBootServletInitializer {

    @Autowired
    private Environment env;

    @Bean(autowire = Autowire.BY_NAME)
    public GlobalValue globalValue() {
        String ssoLoginUrl;
        String secretKey;
        String appId;
        String decryptKey;
        String accessTimeout;
        /**
         * 使用getRequiredProperty，缺少该配置项，则启动报错
         */
        try {
            ssoLoginUrl = env.getRequiredProperty("sso.login.url");
        } catch (IllegalStateException e) {
            throw new RuntimeException("项目缺少必须配置的配置项[sso.login.url]，请检查添加后再重新启动");
        }
        try {
            secretKey = env.getRequiredProperty("secret.key");
        } catch (IllegalStateException e) {
            throw new RuntimeException("项目缺少必须配置的配置项[secret.key]，请检查添加后再重新启动");
        }
        try {
            appId = env.getRequiredProperty("sso.appID");
        } catch (IllegalStateException e) {
            throw new RuntimeException("项目缺少必须配置的配置项[sso.appID]，请检查添加后再重新启动");
        }
        try {
            decryptKey = env.getRequiredProperty("decryptKey");
        } catch (IllegalStateException e) {
            throw new RuntimeException("项目缺少必须配置的配置项[decryptKey]，请检查添加后再重新启动");
        }
        try {
            accessTimeout = env.getRequiredProperty("access.timeout");
        } catch (IllegalStateException e) {
            throw new RuntimeException("项目缺少必须配置的配置项[access.timeout]，请检查添加后再重新启动");
        }

        List<String> ignoreUris = env.getProperty("ignore.ignoreUris", List.class);
        /**
         * 不写该配置项，则默认值为true；配置项值.equalsIgnoreCase("true")，值为true;其余均为false
         */
        Boolean isMultiTenant = Boolean.parseBoolean(env.getProperty("isMultiTenant", "true"));

        return new GlobalValue(ssoLoginUrl, secretKey, appId, ignoreUris, decryptKey, accessTimeout, isMultiTenant);
    }

    @Bean(autowire = Autowire.BY_NAME)
    public ExtValue getExtValue() {
        String cron = env.getProperty("ext.user.sync.cron", "");
        String domain = env.getProperty("ext.user.server.domain", "");
        String ip = env.getProperty("ext.user.server.ip", "");
        String port = env.getProperty("ext.user.server.port", "");
        String baseDn = env.getProperty("ext.user.server.baseDn", "");
        String basePassword = env.getProperty("ext.user.server.basePassword", "");
        String baseDept = env.getProperty("ext.user.server.baseDept", "");
        String extAppId = env.getProperty("ext.user.server.extAppId", "");
        String extAppToken = env.getProperty("ext.user.server.extAppToken", "");
        return new ExtValue(cron, domain, ip, port, baseDn, basePassword, baseDept, extAppId, extAppToken);
    }


    public static void main(String[] args) {
        SpringApplication.run(UamAppliction.class, args);
    }

}
