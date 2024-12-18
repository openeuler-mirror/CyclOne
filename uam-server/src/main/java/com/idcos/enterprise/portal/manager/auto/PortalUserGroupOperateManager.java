

package com.idcos.enterprise.portal.manager.auto;

// auto generated imports

import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.enterprise.portal.form.PortalGroupAndUserCreateForm;
import com.idcos.enterprise.portal.form.PortalGroupAndUserUpdateForm;

/**
 * manager层controller相关的接口自动生成，此文件属于自动生成的，请勿直接修改,具体可以参考codegen工程
 * Generated by <tt>controller-codegen</tt> on 2015-10-28 14:17:42.
 *
 * @author jiaohuizhe
 * @version PortalUserGroupOperateManager.java, v 1.1 2015-10-28 14:17:42 jiaohuizhe Exp $
 */
public interface PortalUserGroupOperateManager {

    /**
     * 分配角色信息
     *
     * @param id
     * @param selectRoles
     * @return BaseResultVO
     */
    CommonResult<?> allocateRole(String id, String selectRoles);

    /**
     * 分配用户信息
     *
     * @param id
     * @param selectUsers
     * @return BaseResultVO
     */
    CommonResult<?> allocateUser(String id, String selectUsers);

    /**
     * 删除用户组
     *
     * @param id
     * @return BaseResultVO
     */
    CommonResult<?> delete(String id);

    /**
     * 更新用户组
     *
     * @param id
     * @param name
     * @param remark
     * @param type
     * @return BaseResultVO
     */
    CommonResult<?> update(String id, String name, String type, String remark);

    /**
     * 创建用户组信息
     *
     * @param name
     * @param remark
     * @param type
     * @param tenant
     * @return BaseResultVO
     */
    CommonResult<?> create(String tenant, String name, String type, String remark);

    /**
     * 新建用户组(包括用户)
     *
     * @param form
     * @return BaseResultVO
     */
    CommonResult<?> createUserGroupAndUser(PortalGroupAndUserCreateForm form);

    /**
     * 修改用户组(包括用户)
     *
     * @param form
     * @return BaseResultVO
     */
    CommonResult<?> updateUserGroupAndUser(PortalGroupAndUserUpdateForm form);

}
