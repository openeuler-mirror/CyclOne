

package com.idcos.enterprise.portal.manager.auto;

// auto generated imports

import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.enterprise.portal.form.ModifyPasswordForm;
import com.idcos.enterprise.portal.form.PortalUserAddForm;
import com.idcos.enterprise.portal.form.PortalUserAllocateGroupForm;
import com.idcos.enterprise.portal.form.PortalUserUpdateForm;

/**
 * manager层controller相关的接口自动生成，此文件属于自动生成的，请勿直接修改,具体可以参考codegen工程
 * Generated by <tt>controller-codegen</tt> on 2015-11-07 16:44:40.
 *
 * @author jiaohuizhe
 * @version PortalUserOperateManager.java, v 1.1 2015-11-07 16:44:40 jiaohuizhe Exp $
 */
public interface PortalUserOperateManager {

    /**
     * 添加用户
     *
     * @param form
     * @return BaseResultVO
     */
    CommonResult<?> add(PortalUserAddForm form);

    /**
     * 修改用户
     *
     * @param form
     * @return
     */
    CommonResult<?> update(PortalUserUpdateForm form);

    /**
     * 激活用户
     *
     * @param id 用户id
     * @return
     */
    CommonResult<?> enabled(String id);

    /**
     * 禁用用户
     *
     * @param id 用户id
     * @return
     */
    CommonResult<?> disabled(String id);

    /**
     * 删除用户
     *
     * @param id 用户id
     * @return
     */
    CommonResult<?> delete(String id);

    /**
     * 分配用户组
     *
     * @param form
     * @return BaseResultVO
     */
    CommonResult<?> allocateGroup(PortalUserAllocateGroupForm form);

    /**
     * 修改用户密码
     *
     * @param form
     * @return BaseResultVO
     */
    CommonResult<?> modifyPassword(ModifyPasswordForm form);

    /**
     * 重置用户密码
     *
     * @param userId
     * @return BaseResultVO
     */
    CommonResult<?> resetPassword(String userId);

    /**
     * 禁用token
     *
     * @param tokenId
     * @return BaseResultVO
     */
    CommonResult<?> forbiddenToken(String tokenId);

    /**
     * 查询token
     *
     * @param tenantId
     * @param loginId
     * @return BaseResultVO
     */
    CommonResult<?> listTokenByTenantIdAndLoginId(String tenantId, String loginId);

}
