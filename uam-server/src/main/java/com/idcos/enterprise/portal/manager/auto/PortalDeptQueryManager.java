/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.manager.auto;

import com.idcos.cloud.core.common.biz.CommonResult;

/**
 * @author Dana
 * @version PortalDeptQueryManager.java, v1 2017/9/26 下午5:26 Dana Exp $$
 */
public interface PortalDeptQueryManager {
    /**
     * 获取所有部门信息
     *
     * @param tenantId
     * @return BaseResultVO
     */
    CommonResult<?> getAllDept(String tenantId);

    /**
     * 获取所有部门的Tree型的树
     *
     * @param tenantId
     * @param treeStyle
     * @return BaseResultVO
     */
    CommonResult<?> getDeptsTree(String tenantId, String treeStyle);

    /**
     * 根据租户id和部门id获取部门信息
     *
     * @param tenantId
     * @param deptId
     * @return BaseResultVO
     */
    CommonResult<?> getDeptByDeptId(String tenantId, String deptId);

    /**
     * 根据部门id获取部门信息
     *
     * @param id
     * @return
     */
    CommonResult<?> getDeptById(String id);

    /**
     * 根据部门id获取对应的角色列表
     *
     * @param id
     * @return
     */
    CommonResult<?> getRolesById(String id);
}
