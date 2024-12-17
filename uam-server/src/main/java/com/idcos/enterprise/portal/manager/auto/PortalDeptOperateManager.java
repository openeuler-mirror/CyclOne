/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.manager.auto;

import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.enterprise.portal.form.PortalDeptAddForm;
import com.idcos.enterprise.portal.form.PortalDeptUpdateForm;

/**
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月04 上午11:43 souakiragen Exp $
 */
public interface PortalDeptOperateManager {

    /**
     * 新增部门
     *
     * @param form
     * @return
     */
    CommonResult<?> add(PortalDeptAddForm form);

    /**
     * 修改部门
     *
     * @param form
     * @return
     */
    CommonResult<?> update(PortalDeptUpdateForm form);

    /**
     * 删除部门
     *
     * @param id
     * @return
     */
    CommonResult<?> delete(String id);

    /**
     * 分配角色
     *
     * @param id
     * @param roleIds
     * @return
     */
    CommonResult<?> assignRole(String id, String roleIds);
}
