

package com.idcos.enterprise.portal.manager.auto;

// auto generated imports

import com.idcos.cloud.core.common.biz.CommonResult;

/**
 * 权限资源类型查询接口
 *
 * @author pengganyu
 * @version $Id: PortalResourceQueryManager.java, v 0.1 2016年5月10日 上午9:36:03 pengganyu Exp $
 */
public interface PortalResourceQueryManager {

    /**
     * 根据权限资源类型编码和有效标志查询权限资源类型信息
     *
     * @param code 权限资源类型编码
     * @return PortalResource 权限资源类型
     */
    public CommonResult<?> queryByCode(String code);

    /**
     * 查询所有的权限资源类型
     *
     * @return
     */
    public CommonResult<?> queryAll();
}