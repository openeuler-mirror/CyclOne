/**
 * 杭州云霁科技有限公司 http://www.idcos.com Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.manager.common;

import java.util.*;

import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.google.common.collect.Lists;
import com.idcos.enterprise.portal.dal.entity.*;
import com.idcos.enterprise.portal.dal.repository.*;
import com.idcos.enterprise.portal.services.PortalDeptRoleRelService;
import com.idcos.enterprise.portal.services.PortalGroupRoleRelService;

/**
 * @author Dana
 * @version CommonManager.java, v1 2017/12/10 下午1:35 Dana Exp $$
 */
@Service
public class CommonManager {
    @Autowired
    private PortalUserRepository portalUserRepository;

    @Autowired
    private PortalGroupUserRelRepository portalGroupUserRelRepository;

    @Autowired
    private PortalGroupRoleRelService portalGroupRoleRelService;

    @Autowired
    private PortalDeptRoleRelService portalDeptRoleRelService;

    @Autowired
    private PortalDeptRepository portalDeptRepository;

    /**
     * 根据用户Id获取所有角色(已去重)
     *
     * @param id
     * @return
     */
    public List<String> getRoleIdsById(String id) {
        PortalUser portalUser = portalUserRepository.findOne(id);
        List<String> roleIds = Lists.newArrayList();
        /**
         * 获取用户所属的用户组
         */
        List<PortalGroupUserRel> portalGroupUserRels = portalGroupUserRelRepository.findByUserId(id);
        /**
         * 获取用户组对应的角色
         */
        List<String> groupIds = Lists.newArrayList();
        for (PortalGroupUserRel portalGroupUserRel : portalGroupUserRels) {
            groupIds.add(portalGroupUserRel.getGroupId());
        }
        roleIds.addAll(portalGroupRoleRelService.getRoleIdsByGroupIds(groupIds));
        /**
         * 获取用户对应的部门，包括所有的父部门(admin用户的部门id可能为null)
         */
        if (portalUser.getDeptId() != null) {
            PortalDept portalDept = portalDeptRepository.findOne(portalUser.getDeptId());
            if (portalDept != null) {
                List<String> deptIds = this.getDeptIdsByParentIdAndAllParent(portalDept.getParentId());
                deptIds.add(portalDept.getId());
                /**
                 * 获取部门对应的角色
                 */
                roleIds.addAll(portalDeptRoleRelService.getRoleIdsByDeptIds(deptIds));
            }
        }

        List<String> distinctRoleIds = new ArrayList<>(new HashSet<>(roleIds));

        return distinctRoleIds;
    }

    /**
     * 根据parent级联查询所有父级部门id
     *
     * @param parentId
     * @return
     */
    public List<String> getDeptIdsByParentIdAndAllParent(String parentId) {
        List<String> parentIds = Lists.newArrayList();
        parentIds.add(parentId);
        while (StringUtils.isNotBlank(parentId)) {
            PortalDept portalDept = portalDeptRepository.findByDeptId(parentId);
            String tempParentId = portalDept.getParentId();
            if (StringUtils.isNotBlank(tempParentId)) {
                parentIds.add(tempParentId);
                parentId = tempParentId;
            } else {
                parentId = tempParentId;
            }
        }
        return parentIds;
    }

}