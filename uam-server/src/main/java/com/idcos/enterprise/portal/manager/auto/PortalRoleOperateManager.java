

package com.idcos.enterprise.portal.manager.auto;

// auto generated imports

import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.enterprise.portal.form.PortalRoleCreateForm;
import com.idcos.enterprise.portal.form.PortalRoleUpdateForm;

/**
 * manager层controller相关的接口自动生成，此文件属于自动生成的，请勿直接修改,具体可以参考codegen工程
 * Generated by <tt>controller-codegen</tt> on 2015-10-30 15:00:49.
 *
 * @author jiaohuizhe
 * @version PortalRoleOperateManager.java, v 1.1 2015-10-30 15:00:49 jiaohuizhe Exp $
 */
public interface PortalRoleOperateManager {


    /**
     * 分配用户组
     *
     * @param id
     * @param selGroups
     * @return BaseResultVO
     */
    public CommonResult<?> allocateGroup(String id, String selGroups);


    /**
     * 删除角色信息
     *
     * @param id
     * @return BaseResultVO
     */
    public CommonResult<?> delete(String id);


    /**
     * 更新角色信息
     *
     * @param id
     * @param form
     * @return BaseResultVO
     */
    public CommonResult<?> update(String id, PortalRoleUpdateForm form);


    /**
     * 创建角色信息
     *
     * @param form
     * @return BaseResultVO
     */
    public CommonResult<?> create(PortalRoleCreateForm form);

}
