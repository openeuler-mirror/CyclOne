/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.dal.entity;

import com.idcos.cloud.core.common.BaseVO;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import javax.persistence.*;
import java.io.Serializable;

/**
 * @author Dana
 * @version PortalSysDict.java, v1 2017/11/24 上午1:40 Dana Exp $$
 */
@Entity
@Table(name = "PORTAL_SYS_DICT")
@IdClass(PortalSysDict.PortalSysDictPrimarkKey.class)
@EntityListeners(AuditingEntityListener.class)
public class PortalSysDict extends BaseVO implements Serializable {

    private static final long serialVersionUID = 741231858441822677L;

    //========== properties ==========
    /**
     * This property corresponds to db column <tt>TYPE_CODE</tt>.
     * 字段备注:<tt>系统字典类型编码</tt>
     * 字段类型:<tt>varchar</tt>
     * 字段长度:<tt>64</tt>
     * 可否为空:<tt>不可为空</tt>
     */
    @Id
    @Column(name = "TYPE_CODE")
    private String typeCode;

    /**
     * This property corresponds to db column <tt>CODE</tt>.
     * 字段备注:<tt>系统字典编码</tt>
     * 字段类型:<tt>varchar</tt>
     * 字段长度:<tt>64</tt>
     * 可否为空:<tt>不可为空</tt>
     */
    @Id
    @Column(name = "CODE")
    private String code;

    /**
     * This property corresponds to db column <tt>VALUE</tt>.
     * 字段备注:<tt>参数值</tt>
     * 字段类型:<tt>text</tt>
     * 可否为空:<tt>可为空</tt>
     */
    @Column(name = "VALUE")
    private String value;

    /**
     * This property corresponds to db column <tt>TENANT_ID</tt>.
     * 字段备注:<tt>租户</tt>
     * 字段类型:<tt>varchar</tt>
     * 字段长度:<tt>64</tt>
     * 可否为空:<tt>不可为空</tt>
     */
    @Id
    @Column(name = "TENANT_ID")
    private String tenantId;

    /**
     * This property corresponds to db column <tt>REMARK</tt>.
     * 字段备注:<tt>说明</tt>
     * 字段类型:<tt>text</tt>
     * 可否为空:<tt>可为空</tt>
     */
    @Column(name = "REMARK")
    private String remark;

    public String getTypeCode() {
        return typeCode;
    }

    public void setTypeCode(String typeCode) {
        this.typeCode = typeCode;
    }

    public String getCode() {
        return code;
    }

    public void setCode(String code) {
        this.code = code;
    }

    public String getValue() {
        return value;
    }

    public void setValue(String value) {
        this.value = value;
    }

    public String getTenantId() {
        return tenantId;
    }

    public void setTenantId(String tenantId) {
        this.tenantId = tenantId;
    }

    public String getRemark() {
        return remark;
    }

    public void setRemark(String remark) {
        this.remark = remark;
    }

    public static class PortalSysDictPrimarkKey implements Serializable {
        private static final long serialVersionUID = 741231858441822666L;

        private String typeCode;

        private String code;

        private String tenantId;

        public String getTypeCode() {
            return typeCode;
        }

        public void setTypeCode(String typeCode) {
            this.typeCode = typeCode;
        }

        public String getCode() {
            return code;
        }

        public void setCode(String code) {
            this.code = code;
        }

        public String getTenantId() {
            return tenantId;
        }

        public void setTenantId(String tenantId) {
            this.tenantId = tenantId;
        }
    }
}