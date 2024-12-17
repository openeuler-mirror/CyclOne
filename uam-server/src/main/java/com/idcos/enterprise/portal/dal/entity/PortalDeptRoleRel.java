package com.idcos.enterprise.portal.dal.entity;

import com.idcos.cloud.core.common.BaseVO;
import com.idcos.cloud.core.dal.common.ColumnMeta;
import org.hibernate.annotations.GenericGenerator;

import javax.persistence.*;
import java.io.Serializable;

/**
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月04 下午3:28 souakiragen Exp $
 */
@Entity
@Table(name = "PORTAL_DEPT_ROLE_REL")
public class PortalDeptRoleRel extends BaseVO implements Serializable {

    private static final long serialVersionUID = 741231858441822677L;

    //========== properties ==========
    /**
     * 数据库字段 <tt>ID</tt>.
     */
    @Id
    @GeneratedValue(strategy = GenerationType.AUTO, generator = "system-uuid")
    @GenericGenerator(name = "system-uuid", strategy = "uuid2")
    @Column(name = "ID")
    @ColumnMeta(name = "ID", length = 64, description = "关系ID", pk = true, unique = false, nullable = false, scale = 0)
    private String id;

    /**
     * 数据库字段 <tt>DEPT_ID</tt>.
     */
    @Column(name = "DEPT_ID")
    @ColumnMeta(name = "DEPT_ID", length = 64, description = "部门ID", pk = false, unique = false, nullable = false, scale = 0)
    private String deptId;

    /**
     * 数据库字段 <tt>ROLE_ID</tt>.
     */
    @Column(name = "ROLE_ID")
    @ColumnMeta(name = "ROLE_ID", length = 64, description = "用户组ID", pk = false, unique = false, nullable = false, scale = 0)
    private String roleId;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getDeptId() {
        return deptId;
    }

    public void setDeptId(String deptId) {
        this.deptId = deptId;
    }

    public String getRoleId() {
        return roleId;
    }

    public void setRoleId(String roleId) {
        this.roleId = roleId;
    }
}
