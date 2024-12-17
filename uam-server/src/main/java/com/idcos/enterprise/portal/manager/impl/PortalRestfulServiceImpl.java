/**
 * 杭州云霁科技有限公司 http://www.idcos.com Copyright (c) 2015 All Rights Reserved.
 */

package com.idcos.enterprise.portal.manager.impl;

import static com.idcos.enterprise.portal.UamConstant.INTERROGATION;

import java.util.*;

import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

// auto generated imports

import com.alibaba.fastjson.JSON;
import com.google.common.collect.Lists;
import com.idcos.cloud.core.common.util.ListUtil;
import com.idcos.cloud.core.dal.common.page.Pagination;
import com.idcos.enterprise.portal.biz.common.PortalResponse;
import com.idcos.enterprise.portal.biz.common.ResultCode;
import com.idcos.enterprise.portal.biz.common.utils.CurrentUser;
import com.idcos.enterprise.portal.dal.entity.*;
import com.idcos.enterprise.portal.dal.enums.IsActiveEnum;
import com.idcos.enterprise.portal.dal.repository.*;
import com.idcos.enterprise.portal.manager.auto.PortalRestfulService;
import com.idcos.enterprise.portal.services.PortalUserService;
import com.idcos.enterprise.portal.vo.PortalUserVO;
import com.idcos.enterprise.portal.web.GlobalValue;

/**
 * Manager实现类
 * <p>
 * 第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖
 * </p>
 *
 * @author pengganyu
 * @version v 1.1 2015-06-09 09:26:24 pengganyu Exp $
 */
@Service
public class PortalRestfulServiceImpl implements PortalRestfulService {

    /**
     * Logger类
     */
    private static final Logger Logger = LoggerFactory.getLogger(PortalRestfulServiceImpl.class);

    @Autowired
    private PortalGroupUserRelRepository portalGroupUserRelRepository;

    @Autowired
    private PortalUserGroupRepository portalUserGroupRepository;

    @Autowired
    private PortalGroupRoleRelRepository portalGroupRoleRelRepository;

    @Autowired
    private PortalPermissionRepository portalPermissionRepository;

    @Autowired
    private PortalResourceRepository portalResourceRepository;

    @Autowired
    private PortalUserService portalUserService;

    @Autowired
    private CurrentUser currentUser;

    @Autowired
    private GlobalValue globalValue;

    /**
     * 页面不传tenantId的话。租户信息从token中获取
     *
     * @param tenantId
     * @return
     */
    private String getTenant(String tenantId) {
        String tenant;
        // 若页面不传的话
        if (StringUtils.isBlank(tenantId)) {
            tenant = currentUser.getUser().getTenantId();
        } else {
            tenant = tenantId;
        }
        return tenant;
    }

    /**
     * 根据用户组id查询用户ID列表
     *
     * @return BaseResultVO
     */
    @Override
    public PortalResponse queryUserIdsByGroupId(String groupId) {
        PortalResponse response = new PortalResponse();

        if (StringUtils.isBlank(groupId)) {
            response.setMessage("用户组Id为空");
            response.setStatus(ResultCode.FAILE.getCode());
            return response;
        }

        List<String> userIds = Lists.newArrayList();

        for (PortalGroupUserRel rel : portalGroupUserRelRepository.findByGroupId(groupId)) {
            userIds.add(rel.getUserId());
        }

        response.setMessage("查询成功");
        response.setStatus(ResultCode.SUCCESS.getCode());
        response.setContent(userIds);

        return response;
    }

    @Override
    public PortalResponse queryUserIdsByGroupName(String groupName) {
        PortalResponse rsp = new PortalResponse();

        if (StringUtils.isBlank(groupName)) {
            rsp.setStatus(ResultCode.FAILE.getCode());
            rsp.setMessage("用户组名称不能为空");
            return rsp;
        }

        List<PortalUserGroup> groups = portalUserGroupRepository.findByName(groupName);

        if (groups.size() == 0) {
            rsp.setStatus(ResultCode.SUCCESS.getCode());
            rsp.setMessage("查询成功");
            rsp.setContent(Lists.newArrayList());
            return rsp;
        }

        if (groups.size() > 1) {
            rsp.setStatus(ResultCode.FAILE.getCode());
            rsp.setMessage("用户组名称不唯一");
            return rsp;
        }
        return queryUserIdsByGroupId(groups.get(0).getId());
    }

    @Override
    public PortalResponse queryUserByCnd(String tenantId, String cnd, String pageNo, String pageSize) {
        Logger.info("===============================根据查询条件查询用户信息，查询条件为：" + cnd);
        Pagination<PortalUserVO> userPagination =
            portalUserService.getPageList(tenantId, null, cnd, Integer.parseInt(pageNo), Integer.parseInt(pageSize));

        PortalResponse rsp = new PortalResponse();

        if (userPagination.getTotalCount() == 0) {
            rsp.setStatus(ResultCode.FAILE.getCode());
            rsp.setMessage("查询用户列表为空");
            return rsp;
        }
        rsp.setStatus(ResultCode.SUCCESS.getCode());
        rsp.setMessage("查询成功");
        rsp.setContent(userPagination);
        return rsp;
    }

    /**
     * 根据用户id查询用户信息
     *
     * @see com.idcos.enterprise.portal.manager.auto.PortalRestfulService#queryAuthority(java.lang.String)
     */
    @SuppressWarnings("unchecked")
    @Override
    public PortalResponse queryAuthority(final String userId) {
        PortalResponse rsp = new PortalResponse();

        if (StringUtils.isBlank(userId)) {
            rsp.setStatus(ResultCode.FAILE.getCode());
            rsp.setMessage("用户id为空");
        }

        Logger.info("===============================开始查询用户信息");
        Logger.info("===============================入参：" + userId);

        Map<String, Object> returnMap = new HashMap<>(3);

        Set<String> roleList = new HashSet<>();
        Set<String> groupList = new HashSet<>();

        Map<String, Set<String>> perMap = new HashMap<>(9);

        // 获取用户组信息
        for (PortalGroupUserRel rel : portalGroupUserRelRepository.findByUserId(userId)) {
            if (StringUtils.isNotBlank(rel.getGroupId())) {
                groupList.add(rel.getGroupId());
            }
        }

        // 根据用户组id查询角色
        for (String groupId : groupList) {
            for (PortalGroupRoleRel rel : portalGroupRoleRelRepository.findByGroupId(groupId)) {
                if (StringUtils.isNotBlank(rel.getRoleId())) {
                    roleList.add(rel.getRoleId());
                }
            }
        }

        // 获取权限信息
        if (roleList.size() > 0) {
            Map<?, List<PortalPermission>> tmpMap = ListUtil.groupBy(
                portalPermissionRepository.queryAuthIdsByAuthObjs(globalValue.getAppId(), roleList), "authResType");
            for (Object key : tmpMap.keySet()) {
                perMap.put((String)key, new HashSet<String>(ListUtil.filter(tmpMap.get(key), "authResId")));
            }
        }

        returnMap.put("roles", roleList);
        returnMap.put("permissions", perMap);
        returnMap.put("groups", groupList);

        rsp.setStatus(ResultCode.SUCCESS.getCode());
        rsp.setMessage("查询成功");
        rsp.setContent(JSON.toJSON(returnMap));

        Logger.info("===============================查询结束。");
        return rsp;
    }

    @Override
    public PortalResponse queryResource(String code, String tenantId) {
        PortalResponse rsp = new PortalResponse();

        // 校验权限资源配置信息
        PortalResource res = portalResourceRepository.findByCodeAndIsActive(code, IsActiveEnum.HAS_ACTIVE.getCode());

        if (res == null) {
            rsp.setStatus(ResultCode.FAILE.getCode());
            rsp.setMessage("未查询到有效的权限资源配置信息");
            return rsp;
        }
        if (StringUtils.isBlank(res.getUrl())) {
            rsp.setStatus(ResultCode.FAILE.getCode());
            rsp.setMessage("权限资源" + res.getName() + "URL配置为空");
            return rsp;
        }

        StringBuffer url = new StringBuffer(res.getUrl());
        if (res.getUrl().contains(INTERROGATION)) {
            url.append("&");
        } else {
            url.append(INTERROGATION);
        }
        url.append("tenantId=" + tenantId);

        Logger.info("===============================开始查询" + res.getName() + "权限资源信息,URL地址为：" + url);

        rsp = new RestTemplate().getForObject(StringUtils.trim(url.toString()), PortalResponse.class);

        Logger.info("===============================查询结束，返回信息:" + rsp.getMessage());

        return rsp;
    }

}
