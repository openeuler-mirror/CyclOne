/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2018 All Rights Reserved.
 */
package com.idcos.enterprise.portal.vo;

import com.idcos.cloud.core.common.BaseVO;

/**
 * @author Xizhao.Dai
 * @version JwtTokenVO.java, v1 2018/4/19 下午9:15 Xizhao.Dai Exp $$
 */
public class JwtTokenVO extends BaseVO {
    private String userId;

    private String name;

    private String loginId;

    private String tenantId;

    private String tenantName;

    private String expireTime;

    private String creatTime;

    public String getUserId() {
        return userId;
    }

    public void setUserId(String userId) {
        this.userId = userId;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getLoginId() {
        return loginId;
    }

    public void setLoginId(String loginId) {
        this.loginId = loginId;
    }

    public String getTenantId() {
        return tenantId;
    }

    public void setTenantId(String tenantId) {
        this.tenantId = tenantId;
    }

    public String getTenantName() {
        return tenantName;
    }

    public void setTenantName(String tenantName) {
        this.tenantName = tenantName;
    }

    public String getExpireTime() {
        return expireTime;
    }

    public void setExpireTime(String expireTime) {
        this.expireTime = expireTime;
    }

    public String getCreatTime() {
        return creatTime;
    }

    public void setCreatTime(String creatTime) {
        this.creatTime = creatTime;
    }

    @Override
    public String toString() {
        return "JwtTokenVO{" + "userId='" + userId + '\'' + ", name='" + name + '\'' + ", loginId='" + loginId + '\''
                + ", tenantId='" + tenantId + '\'' + ", tenantName='" + tenantName + '\'' + ", expireTime='" + expireTime
                + '\'' + ", creatTime='" + creatTime + '\'' + '}';
    }
}