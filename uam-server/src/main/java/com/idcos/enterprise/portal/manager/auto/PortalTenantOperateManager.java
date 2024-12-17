/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.manager.auto;

import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.enterprise.portal.form.PortalTenantAddForm;
import com.idcos.enterprise.portal.form.PortalTenantUpdateForm;

/**
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月07 上午9:32 souakiragen Exp $
 */
public interface PortalTenantOperateManager {

    /**
     * 添加租户
     *
     * @param form
     * @return
     */
    CommonResult<?> add(PortalTenantAddForm form);

    /**
     * 更新租户
     *
     * @param form
     * @return
     */
    CommonResult<?> update(PortalTenantUpdateForm form);

    /**
     * 删除租户
     *
     * @param id
     * @return
     */
    CommonResult<?> delete(String id);
}
