package com.idcos.enterprise.portal.manager.auto;

import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.enterprise.portal.form.PortalTenantQueryPageListForm;

/**
 * @author GuanBin
 * @version PortalTenantManager.java, v1 2017/10/10 上午7:40 GuanBin Exp $$
 */
public interface PortalTenantManager {

    /**
     * 获取租户户信息
     *
     * @return BaseResultVO
     */
    CommonResult<?> getAllTenant();

    /**
     * 获取租户户信息
     *
     * @param tenantId
     * @return BaseResultVO
     */
    CommonResult<?> getTenantByTenantId(String tenantId);

    /**
     * 获取单个租户信息
     *
     * @param id
     * @return
     */
    CommonResult<?> getTenantInfo(String id);

    /**
     * 分页查询租户
     *
     * @param pageNo
     * @param pageSize
     * @param form
     * @return
     */
    CommonResult<?> getPageList(int pageNo, int pageSize, PortalTenantQueryPageListForm form);
}
