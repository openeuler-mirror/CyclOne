/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.dal.repository;

import com.idcos.cloud.core.dal.common.jpa.BaseRepository;
import com.idcos.enterprise.portal.dal.entity.PortalDeptRoleRel;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;

import java.util.List;

/**
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月04 下午3:37 souakiragen Exp $
 */
public interface PortalDeptRoleRelRepository extends BaseRepository<PortalDeptRoleRel, String> {
    /**
     * 根据部门id删除
     *
     * @param deptId
     * @return
     */
    @Query("delete from PortalDeptRoleRel dr where dr.deptId = ?1")
    @Modifying
    int deleteByDeptId(String deptId);

    /**
     * 根据部门id查找
     *
     * @param deptId
     * @return
     */
    @Query("select r from PortalDeptRoleRel r where r.deptId = ?1")
    List<PortalDeptRoleRel> findByDeptId(String deptId);
}
