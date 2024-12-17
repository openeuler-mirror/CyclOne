

package com.idcos.enterprise.portal.manager.auto;

// auto generated imports

import com.idcos.enterprise.portal.biz.common.PortalResponse;

/**
 * 第三方API service类
 *
 * @author pengganyu
 * @version $Id: PortalResfulService.java, v 0.1 2016年5月11日 下午5:03:50 pengganyu Exp $
 */
public interface PortalRestfulService {

    /**
     * 根据用户id获取用户的用户组，角色，权限信息
     *
     * @param userId 用户id
     * @return BaseResultVO
     */
    PortalResponse queryAuthority(String userId);

    /**
     * 根据查询条件查询用户信息
     *
     * @param cnd      查询条件
     * @param pageNo   页号
     * @param pageSize 页大小
     * @param tenantId
     * @return
     */
    PortalResponse queryUserByCnd(String tenantId, String cnd, String pageNo, String pageSize);

    /**
     * 根据用户组id查询用户ID列表
     *
     * @param groupId 用户组id
     * @return BaseResultVO
     */
    PortalResponse queryUserIdsByGroupId(String groupId);

    /**
     * 根据用户组名称查询用户ID列表
     *
     * @param groupName 用户组id
     * @return BaseResultVO
     */
    PortalResponse queryUserIdsByGroupName(String groupName);

    /**
     * 获取权限资源信息
     *
     * @param resType  权限资源类型
     * @param tenantId 目标租户code
     * @return
     */
    PortalResponse queryResource(String resType, String tenantId);

}
