/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.manager.impl;

import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.cloud.core.common.util.ListUtil;
import com.idcos.enterprise.portal.biz.common.CommonBizException;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessCallback;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessContext;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessTemplate;
import com.idcos.enterprise.portal.biz.common.utils.Base64Util;
import com.idcos.enterprise.portal.biz.common.utils.CurrentUser;
import com.idcos.enterprise.portal.biz.common.utils.PasswordUtil;
import com.idcos.enterprise.portal.dal.entity.PortalDept;
import com.idcos.enterprise.portal.dal.entity.PortalTenant;
import com.idcos.enterprise.portal.dal.entity.PortalUser;
import com.idcos.enterprise.portal.dal.entity.PortalUserGroup;
import com.idcos.enterprise.portal.dal.enums.IsActiveEnum;
import com.idcos.enterprise.portal.dal.enums.PortalUserStatusEnum;
import com.idcos.enterprise.portal.dal.repository.PortalDeptRepository;
import com.idcos.enterprise.portal.dal.repository.PortalTenantRepository;
import com.idcos.enterprise.portal.dal.repository.PortalUserGroupRepository;
import com.idcos.enterprise.portal.dal.repository.PortalUserRepository;
import com.idcos.enterprise.portal.form.PortalTenantAddForm;
import com.idcos.enterprise.portal.form.PortalTenantUpdateForm;
import com.idcos.enterprise.portal.manager.auto.PortalTenantOperateManager;
import com.idcos.enterprise.portal.services.PortalTenantService;
import com.idcos.enterprise.portal.services.PortalUserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.List;

/**
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月07 上午9:42 souakiragen Exp $
 */
@Service
public class PortalTenantOperateManagerImpl implements PortalTenantOperateManager {

    @Autowired
    private BusinessProcessTemplate businessProcessTemplate;

    @Autowired
    private PortalTenantRepository portalTenantRepository;

    @Autowired
    private CurrentUser currentUser;

    @Autowired
    private PortalTenantService portalTenantService;

    @Autowired
    private PortalUserRepository portalUserRepository;

    @Autowired
    private PortalDeptRepository portalDeptRepository;

    @Autowired
    private PortalUserGroupRepository portalUserGroupRepository;

    @Override
    public CommonResult<?> add(final PortalTenantAddForm form) {
        return businessProcessTemplate.process(new BusinessProcessCallback<Object>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
                PortalTenant portalTenant = portalTenantRepository.findByName(form.getTenantId());
                if (portalTenant != null) {
                    throw new RuntimeException("租户id已存在");
                }
                portalTenant = portalTenantRepository.findByDisplayName(form.getTenantName());
                if (portalTenant != null) {
                    throw new RuntimeException("租户名称已存在");
                }
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {

            }

            @Override
            public Object doBusiness(BusinessProcessContext context) {
                PortalTenant portalTenant = new PortalTenant();
                portalTenant.setName(form.getTenantId());
                portalTenant.setDisplayName(form.getTenantName());
                portalTenant.setGmtCreate(new Date());
                portalTenant.setGmtModified(portalTenant.getGmtCreate());
                portalTenant.setIsActive(IsActiveEnum.HAS_ACTIVE.getCode());
                portalTenantRepository.save(portalTenant);
                createAdmin(portalTenant.getName());
                return null;
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
    }

    /**
     * 新建租户时创建admin用户
     *
     * @param tenantId
     */
    private void createAdmin(String tenantId) {
        PortalUser user = new PortalUser();
        user.setName("管理员");
        user.setLoginId("admin");
        user.setTenantId(tenantId);
        user.setStatus(PortalUserStatusEnum.ENABLED.getCode());
        user.setIsActive(IsActiveEnum.HAS_ACTIVE.getCode());
        user.setCreateUser(currentUser.getUser().getUserId());
        user.setCreateTime(new Date());
        user.setLastModifiedTime(new Date());
        portalUserRepository.save(user);
        PortalUser portalUser = portalUserRepository.findPortalUserById(tenantId, user.getLoginId());
        //密码密码后保存到数据库.
        //由于加密需要用户的ID，所以需要保存用户后再更新密码。
        //加密密码
        try {
            //salt随机产生，生成密码时使用的密码是用户的ID.
            byte[] salt = PasswordUtil.getSalt();
            String encriptPW = PasswordUtil.encrypt("zaq1@WSX", portalUser.getId(), salt);

            String saltStr = Base64Util.encode(salt);
            user.setPassword(encriptPW);
            user.setSalt(saltStr);
            portalUserRepository.save(user);
        } catch (Exception e) {
            throw new RuntimeException("系统错误，请联系管理员。");
        }
    }

    @Override
    public CommonResult<?> update(final PortalTenantUpdateForm form) {
        final PortalTenant oldPortalTenant = portalTenantRepository.findOne(form.getId());
        return businessProcessTemplate.process(new BusinessProcessCallback<Object>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
                if (oldPortalTenant == null) {
                    throw new RuntimeException("该租户不存在");
                }
                PortalTenant portalTenant = portalTenantRepository.findByNameNotId(form.getTenantId(), form.getId());
                if (portalTenant != null) {
                    throw new RuntimeException("租户id已存在");
                }
                portalTenant = portalTenantRepository.findByDisplayNameNotId(form.getTenantName(), form.getId());
                if (portalTenant != null) {
                    throw new RuntimeException("租户名称已存在");
                }
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {

            }

            @Override
            public Object doBusiness(BusinessProcessContext context) {
                PortalTenant portalTenant = new PortalTenant();
                portalTenant.setId(oldPortalTenant.getId());
                portalTenant.setName(form.getTenantId());
                portalTenant.setDisplayName(form.getTenantName());
                portalTenant.setGmtCreate(oldPortalTenant.getGmtCreate());
                portalTenant.setGmtModified(new Date());
                portalTenant.setIsActive(oldPortalTenant.getIsActive());
                portalTenantRepository.save(portalTenant);
                return null;
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
    }

    @Override
    public CommonResult<?> delete(final String id) {
        return businessProcessTemplate.process(new BusinessProcessCallback<Object>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {

            }

            @Override
            public Object doBusiness(BusinessProcessContext context) {
                PortalTenant portalTenant = portalTenantRepository.findById(id);
                //删除租户
                portalTenantService.updateIsActive(id, IsActiveEnum.NO_ACTIVE.getCode());
                //删除用户
                List<PortalUser> portalUsers = portalUserRepository.findByTenantId(portalTenant.getName());
                ListUtil.update(portalUsers, "isActive", IsActiveEnum.HAS_ACTIVE.getCode(),
                        IsActiveEnum.NO_ACTIVE.getCode());
                portalUserRepository.save(portalUsers);
                //删除用户组
                List<PortalUserGroup> userGroups = portalUserGroupRepository
                        .findAllUserGroupByTenantId(portalTenant.getName());
                ListUtil.update(userGroups, "isActive", IsActiveEnum.HAS_ACTIVE.getCode(),
                        IsActiveEnum.NO_ACTIVE.getCode());
                portalUserGroupRepository.save(userGroups);
                //删除部门
                List<PortalDept> portalDepts = portalDeptRepository.findAllDeptByTenantId(portalTenant.getName());
                ListUtil.update(portalDepts, "status", "1", "0");
                portalDeptRepository.save(portalDepts);
                return null;
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
    }
}
