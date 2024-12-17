

package com.idcos.enterprise.portal.vo;

// auto generated imports

import com.idcos.cloud.core.common.BaseVO;

/**
 * PortalResourceVO
 *
 * @author pengganyu
 * @version $Id: PortalResourceVO.java, v 0.1 2016年5月10日 上午9:46:42 pengganyu Exp $
 */
public class PortalResourceVO extends BaseVO {

    private static final long serialVersionUID = 1L;

    /**
     * 权限资源主键
     */
    private String id;

    /**
     * 权限资源类型
     */
    private String code;

    /**
     * 权限资源名称
     */
    private String name;

    /**
     * 权限资源URL
     */
    private String url;

    /**
     * 系统名称
     */
    private String appId;

    /**
     * 备注
     */
    private String remark;


    private String tenantId;

    public String getTenantId() {
        return tenantId;
    }

    public void setTenantId(String tenantId) {
        this.tenantId = tenantId;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getCode() {
        return code;
    }

    public void setCode(String code) {
        this.code = code;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getUrl() {
        return url;
    }

    public void setUrl(String url) {
        this.url = url;
    }

    public String getAppId() {
        return appId;
    }

    public void setAppId(String appId) {
        this.appId = appId;
    }

    public String getRemark() {
        return remark;
    }

    public void setRemark(String remark) {
        this.remark = remark;
    }

    @Override
    public String toString() {
        return "PortalResourceVO{" +
                "id='" + id + '\'' +
                ", code='" + code + '\'' +
                ", name='" + name + '\'' +
                ", url='" + url + '\'' +
                ", appId='" + appId + '\'' +
                ", remark='" + remark + '\'' +
                ", tenantId='" + tenantId + '\'' +
                '}';
    }
}