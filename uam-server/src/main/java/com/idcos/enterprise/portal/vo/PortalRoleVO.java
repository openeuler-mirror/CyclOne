

package com.idcos.enterprise.portal.vo;

import com.idcos.cloud.core.common.BaseLinkVO;

import java.util.Arrays;
import java.util.Date;

// auto generated imports

/**
 * 返回结果对象 {vo.className}
 * <p>第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖</p>
 *
 * @author yanlv
 * @version PortalRoleVO.java, v 1.1 2015-06-06 10:09:55 yanlv Exp $
 */

public class PortalRoleVO extends BaseLinkVO {

    private static final long serialVersionUID = 1L;

    private String id;
    private String code;
    private String name;
    private String remark;
    private Date gmtCreate;
    private String isActive;
    private String[] selGroups;
    private String[] selUsers;
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

    public String getRemark() {
        return remark;
    }

    public void setRemark(String remark) {
        this.remark = remark;
    }

    public String getIsActive() {
        return isActive;
    }

    public void setIsActive(String isActive) {
        this.isActive = isActive;
    }

    public String[] getSelGroups() {
        return selGroups;
    }

    public void setSelGroups(String[] selGroups) {
        this.selGroups = selGroups;
    }

    public String[] getSelUsers() {
        return selUsers;
    }

    public void setSelUsers(String[] selUsers) {
        this.selUsers = selUsers;
    }

    public Date getGmtCreate() {
        return gmtCreate;
    }

    public void setGmtCreate(Date gmtCreate) {
        this.gmtCreate = gmtCreate;
    }

    @Override
    public String toString() {
        return "PortalRoleVO{" +
                "id='" + id + '\'' +
                ", code='" + code + '\'' +
                ", name='" + name + '\'' +
                ", remark='" + remark + '\'' +
                ", gmtCreate=" + gmtCreate +
                ", isActive='" + isActive + '\'' +
                ", selGroups=" + Arrays.toString(selGroups) +
                ", selUsers=" + Arrays.toString(selUsers) +
                ", tenantId='" + tenantId + '\'' +
                '}';
    }
}