

package com.idcos.enterprise.portal.manager.auto;

// auto generated imports

import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.enterprise.portal.form.PortalResourceCreateForm;
import com.idcos.enterprise.portal.form.PortalResourceUpdateForm;

/**
 * 权限资源类型查询接口
 *
 * @author pengganyu
 * @version $Id: PortalResourceQueryManager.java, v 0.1 2016年5月10日 上午9:36:03 pengganyu Exp $
 */
public interface PortalResourceOperateManager {

    /**
     * 保存权限资源类型
     *
     * @param form
     * @return ${method.returnType}
     */
    CommonResult<?> create(PortalResourceCreateForm form);

    /**
     * 修改权限资源类型
     *
     * @param form
     * @param code
     * @return ${method.returnType}
     */
    CommonResult<?> update(String code, PortalResourceUpdateForm form);

    /**
     * 删除权限资源类型
     *
     * @param id 权限资源id
     * @return
     */
    CommonResult<?> delete(String id);

}
