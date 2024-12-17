/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.dal.entity;

import com.idcos.cloud.core.common.BaseVO;

import javax.persistence.*;
import java.io.Serializable;
import java.util.Date;

/**
 * @author Dana
 * @version PortalToken.java, v1 2017/11/30 上午1:01 Dana Exp $$
 */
@Entity
@Table(name = "PORTAL_TOKEN",
        indexes = {@Index(columnList = "TOKEN_CRC")})
public class PortalToken extends BaseVO implements Serializable {
    private static final long serialVersionUID = 741231858441822655L;

    //========== properties ==========
    /**
     * This property corresponds to db column <tt>ID</tt>.
     * 字段备注:<tt>主键</tt>
     * 字段类型:<tt>varchar</tt>
     * 字段长度:<tt>64</tt>
     * 可否为空:<tt>不可为空</tt>
     */
    @Id
    @Column(name = "ID")
    private String id;

    /**
     * This property corresponds to db column <tt>NAME</tt>.
     * 字段备注:<tt>token值</tt>
     * 字段类型:<tt>text</tt>
     * 字段长度:<tt>未定义</tt>
     * 可否为空:<tt>可为空</tt>
     */
    @Column(name = "NAME")
    private String name;

    /**
     * This property corresponds to db column <tt>TOKEN_CRC</tt>.
     * 字段备注:<tt>token串的crc哈希值</tt>
     * 字段类型:<tt>int unsigned</tt>
     * 字段长度:<tt>未定义</tt>
     * 可否为空:<tt>不可为空</tt>
     */
    @Column(name = "TOKEN_CRC")
    private long tokenCrc;

    /**
     * This property corresponds to db column <tt>LOGIN_ID</tt>.
     * 字段备注:<tt>登录名</tt>
     * 字段类型:<tt>varchar</tt>
     * 字段长度:<tt>64</tt>
     * 可否为空:<tt>不可为空</tt>
     */
    @Column(name = "LOGIN_ID")
    private String loginId;

    /**
     * This property corresponds to db column <tt>TENANT_ID</tt>.
     * 字段备注:<tt>租户code</tt>
     * 字段类型:<tt>varchar</tt>
     * 字段长度:<tt>64</tt>
     * 可否为空:<tt>不可为空</tt>
     */
    @Column(name = "TENANT_ID")
    private String tenantId;

    /**
     * This property corresponds to db column <tt>IS_ACTIVE</tt>.
     * 字段备注:<tt>是否可用</tt>
     * 字段类型:<tt>char</tt>
     * 可否为空:<tt>不可为空</tt>
     */
    @Column(name = "IS_ACTIVE")
    private String isActive;

    /**
     * This property corresponds to db column <tt>EXPIRE_TIME</tt>.
     * 字段备注:<tt>token过期时间</tt>
     * 字段类型:<tt>DATETIME</tt>
     * 可否为空:<tt>不可为空</tt>
     */
    @Column(name = "EXPIRE_TIME")
    private Date expireTime;

    /**
     * This property corresponds to db column <tt>GMT_CREATE</tt>.
     * 字段备注:<tt>创建日期</tt>
     * 字段类型:<tt>DATETIME</tt>
     * 可否为空:<tt>可为空</tt>
     */
    @Column(name = "GMT_CREATE")
    private Date gmtCreate;

    /**
     * This property corresponds to db column <tt>GMT_MODIFIED</tt>.
     * 字段备注:<tt>修改日期</tt>
     * 字段类型:<tt>DATETIME</tt>
     * 可否为空:<tt>不可为空</tt>
     */
    @Column(name = "GMT_MODIFIED")
    private Date gmtModified;

    /**
     * This property corresponds to db column <tt>REMARK</tt>.
     * 字段备注:<tt>备注</tt>
     * 字段类型:<tt>DATETIME</tt>
     * 可否为空:<tt>可为空</tt>
     */
    @Column(name = "REMARK")
    private String remark;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public long getTokenCrc() {
        return tokenCrc;
    }

    public void setTokenCrc(long tokenCrc) {
        this.tokenCrc = tokenCrc;
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

    public String getIsActive() {
        return isActive;
    }

    public void setIsActive(String isActive) {
        this.isActive = isActive;
    }

    public Date getExpireTime() {
        return expireTime;
    }

    public void setExpireTime(Date expireTime) {
        this.expireTime = expireTime;
    }

    public Date getGmtCreate() {
        return gmtCreate;
    }

    public void setGmtCreate(Date gmtCreate) {
        this.gmtCreate = gmtCreate;
    }

    public Date getGmtModified() {
        return gmtModified;
    }

    public void setGmtModified(Date gmtModified) {
        this.gmtModified = gmtModified;
    }

    public String getRemark() {
        return remark;
    }

    public void setRemark(String remark) {
        this.remark = remark;
    }
}