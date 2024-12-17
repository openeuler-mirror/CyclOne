

package com.idcos.enterprise.portal.vo;

// auto generated imports

import com.idcos.cloud.core.common.BaseVO;

/**
 * 返回结果对象 {vo.className}
 * <p>第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖</p>
 *
 * @author yanlv
 * @version PortalPermissionVO.java, v 1.1 2015-06-09 10:00:37 yanlv Exp $
 */

public class PortalPermissionVO extends BaseVO {

    private static final long serialVersionUID = 1L;

    private String id;

    private String authResType;

    private String authResId;

    private String authResName;

    private String authObjType;

    private String authObjId;

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

    public String getAuthResType() {
        return authResType;
    }

    public void setAuthResType(String authResType) {
        this.authResType = authResType;
    }

    public String getAuthResId() {
        return authResId;
    }

    public void setAuthResId(String authResId) {
        this.authResId = authResId;
    }

    public String getAuthObjType() {
        return authObjType;
    }

    public void setAuthObjType(String authObjType) {
        this.authObjType = authObjType;
    }

    public String getAuthObjId() {
        return authObjId;
    }

    public void setAuthObjId(String authObjId) {
        this.authObjId = authObjId;
    }

    public String getAuthResName() {
        return authResName;
    }

    public void setAuthResName(String authResName) {
        this.authResName = authResName;
    }

    @Override
    public String toString() {
        return "PortalPermissionVO{" +
                "id='" + id + '\'' +
                ", authResType='" + authResType + '\'' +
                ", authResId='" + authResId + '\'' +
                ", authResName='" + authResName + '\'' +
                ", authObjType='" + authObjType + '\'' +
                ", authObjId='" + authObjId + '\'' +
                ", tenantId='" + tenantId + '\'' +
                '}';
    }
}