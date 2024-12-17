/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.manager.impl;

import com.google.common.collect.Lists;
import com.idcos.cloud.biz.common.check.CommonParamtersChecker;
import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.enterprise.portal.biz.common.CommonBizException;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessCallback;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessContext;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessTemplate;
import com.idcos.enterprise.portal.dal.entity.PortalDept;
import com.idcos.enterprise.portal.dal.entity.PortalDeptRoleRel;
import com.idcos.enterprise.portal.dal.entity.PortalUser;
import com.idcos.enterprise.portal.dal.enums.SourceTypeEnum;
import com.idcos.enterprise.portal.dal.enums.StatusEnum;
import com.idcos.enterprise.portal.dal.repository.PortalDeptRepository;
import com.idcos.enterprise.portal.dal.repository.PortalDeptRoleRelRepository;
import com.idcos.enterprise.portal.dal.repository.PortalUserRepository;
import com.idcos.enterprise.portal.form.PortalDeptAddForm;
import com.idcos.enterprise.portal.form.PortalDeptUpdateForm;
import com.idcos.enterprise.portal.manager.auto.PortalDeptOperateManager;
import com.idcos.enterprise.portal.services.PortalDeptService;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.BeanUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.thymeleaf.util.ListUtils;

import java.util.Date;
import java.util.List;

/**
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月04 上午11:44 souakiragen Exp $
 */
@Service
public class PortalDeptOperateManagerImpl implements PortalDeptOperateManager {

    @Autowired
    private BusinessProcessTemplate businessProcessTemplate;

    @Autowired
    private PortalDeptRepository portalDeptRepository;

    @Autowired
    private PortalDeptService portalDeptService;

    @Autowired
    private PortalDeptRoleRelRepository portalDeptRoleRelRepository;

    @Autowired
    private PortalUserRepository portalUserRepository;

    @Override
    public CommonResult<?> add(final PortalDeptAddForm form) {
        return businessProcessTemplate.process(new BusinessProcessCallback<Object>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
                if (StringUtils.isNotBlank(form.getParentId())) {
                    PortalDept portalDept = portalDeptRepository.findByDeptId(form.getParentId());
                    if (portalDept == null) {
                        throw new RuntimeException("父级部门不存在!");
                    } else if (!SourceTypeEnum.NATIVE.getCode().equals(portalDept.getSourceType())) {
                        throw new RuntimeException("数据源为" + portalDept.getSourceType() + "的部门不允许添加子部门!");
                    }
                }
                if (portalDeptService.checkHasByTenantIdAndParentIdAndDisplayName(form.getTenantId(),
                        form.getParentId(), form.getDisplayName())) {
                    throw new RuntimeException("该部门已存在");
                }
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {

            }

            @Override
            public Object doBusiness(BusinessProcessContext context) {
                PortalDept portalDept = new PortalDept();
                BeanUtils.copyProperties(form, portalDept);
                portalDept.setStatus(StatusEnum.Y.getCode());
                portalDept.setGmtCreate(new Date());
                portalDept.setGmtModified(new Date());
                PortalDept dBdept = portalDeptRepository.save(portalDept);
                return dBdept.getId();
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
    }

    @Override
    public CommonResult<?> update(final PortalDeptUpdateForm form) {
        return businessProcessTemplate.process(new BusinessProcessCallback<Object>() {
            @Override
            public void checkParam(BusinessProcessContext context) {

            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {
                if (StringUtils.isNotBlank(form.getParentId())) {
                    PortalDept parentDept = portalDeptRepository.findByDeptId(form.getParentId());
                    if (parentDept == null) {
                        throw new RuntimeException("父级部门不存在!");
                    } else if (!SourceTypeEnum.NATIVE.getCode().equals(parentDept.getSourceType())) {
                        throw new RuntimeException("数据源为" + parentDept.getSourceType() + "的部门不允许添加子部门!");
                    }
                }
                PortalDept portalDept = portalDeptRepository.findByDeptId(form.getId());
                if (!SourceTypeEnum.NATIVE.getCode().equals(portalDept.getSourceType())) {
                    throw new RuntimeException("数据源为" + portalDept.getSourceType() + "的部门不允许修改!");
                }
                if (portalDeptService.checkHasByTenantIdAndParentIdAndDisplayNameAndNotId(form.getTenantId(),
                        form.getParentId(), form.getDisplayName(), form.getId())) {
                    throw new RuntimeException("该部门已存在");
                }
                if (form.getParentId() != null && form.getParentId().equals(form.getId())) {
                    throw new RuntimeException("不能将自己设置成自己的父部门");
                }
            }

            @Override
            public Object doBusiness(BusinessProcessContext context) {
                PortalDept oldPortalDept = portalDeptRepository.findByDeptId(form.getId());
                PortalDept portalDept = new PortalDept();
                BeanUtils.copyProperties(form, portalDept);
                portalDept.setGmtCreate(oldPortalDept.getGmtCreate());
                portalDept.setGmtModified(new Date());
                portalDept.setStatus(oldPortalDept.getStatus());
                return portalDeptRepository.save(portalDept);
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
                PortalDept portalDept = portalDeptRepository.findByDeptId(id);
                if (portalDept == null) {
                    throw new RuntimeException("要删除的部门不存在");
                } else if (!SourceTypeEnum.NATIVE.getCode().equals(portalDept.getSourceType())) {
                    throw new RuntimeException("数据源为" + portalDept.getSourceType() + "的部门不允许删除!");
                }
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {
                //检查部门下是否存在用户
                List<PortalUser> portalUsers = portalUserRepository.findByDeptId(id);
                if (!ListUtils.isEmpty(portalUsers)) {
                    throw new RuntimeException("部门下存在用户，不能删除");
                }
            }

            @Override
            public Object doBusiness(BusinessProcessContext context) {
                portalDeptService.updatePortalDeptStatus(StatusEnum.N.getCode(), id);
                portalDeptRoleRelRepository.deleteByDeptId(id);
                return "删除成功";
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
    }

    @Override
    public CommonResult<?> assignRole(final String id, final String roleIdsStr) {
        return businessProcessTemplate.process(new BusinessProcessCallback<Object>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
                CommonParamtersChecker.checkNotBlank(id);
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {

            }

            @Override
            public Object doBusiness(BusinessProcessContext context) {

                portalDeptRoleRelRepository.deleteByDeptId(id);
                if (StringUtils.isNotBlank(roleIdsStr)) {
                    List<String> roleIds = Lists.newArrayList(roleIdsStr.split(","));
                    List<PortalDeptRoleRel> portalDeptRoleRels = Lists.newArrayList();
                    for (String roleId : roleIds) {
                        PortalDeptRoleRel portalDeptRoleRel = new PortalDeptRoleRel();
                        portalDeptRoleRel.setDeptId(id);
                        portalDeptRoleRel.setRoleId(roleId);
                        portalDeptRoleRels.add(portalDeptRoleRel);
                    }
                    portalDeptRoleRelRepository.save(portalDeptRoleRels);
                }
                return null;
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
    }

}
