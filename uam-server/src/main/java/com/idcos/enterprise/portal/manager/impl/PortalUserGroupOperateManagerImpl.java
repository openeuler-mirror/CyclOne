

package com.idcos.enterprise.portal.manager.impl;

// auto generated imports

import com.alibaba.fastjson.JSON;
import com.idcos.cloud.biz.common.check.CommonParamtersChecker;
import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.cloud.core.common.util.ListUtil;
import com.idcos.enterprise.portal.biz.common.CommonBizException;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessCallback;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessContext;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessTemplate;
import com.idcos.enterprise.portal.biz.common.utils.CrudUtilService;
import com.idcos.enterprise.portal.biz.common.utils.CurrentUser;
import com.idcos.enterprise.portal.convert.PortalUserGroupConvert;
import com.idcos.enterprise.portal.dal.entity.*;
import com.idcos.enterprise.portal.dal.enums.PortalUserGroupTypeEnum;
import com.idcos.enterprise.portal.dal.enums.IsActiveEnum;
import com.idcos.enterprise.portal.dal.repository.*;
import com.idcos.enterprise.portal.form.PortalGroupAllocateUserForm;
import com.idcos.enterprise.portal.form.PortalGroupAndUserCreateForm;
import com.idcos.enterprise.portal.form.PortalGroupAndUserUpdateForm;
import com.idcos.enterprise.portal.manager.auto.PortalUserGroupOperateManager;
import com.idcos.enterprise.portal.services.PortalUserGroupService;
import com.idcos.enterprise.portal.vo.PortalUserGroupVO;
import org.apache.commons.lang3.StringUtils;
import org.apache.poi.ss.formula.functions.T;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Date;
import java.util.List;

/**
 * Manager实现类
 * <p>第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖</p>
 *
 * @author yanlv
 * @version v 1.1 2015-06-09 09:11:34 yanlv Exp $
 */
@Service
public class PortalUserGroupOperateManagerImpl implements PortalUserGroupOperateManager {
    @Autowired
    private BusinessProcessTemplate businessProcessTemplate;

    @Autowired
    private PortalUserGroupConvert portalUserGroupConvert;

    @Autowired
    private CrudUtilService crudUtilService;

    @Autowired
    private PortalUserGroupService portalUserGroupService;

    @Autowired
    private PortalUserRepository portalUserRepository;

    @Autowired
    private PortalUserGroupRepository portalUserGroupRepository;

    @Autowired
    private PortalGroupRoleRelRepository portalGroupRoleRelRepository;

    @Autowired
    private PortalGroupUserRelRepository portalGroupUserRelRepository;

    @Autowired
    private CurrentUser currentUser;

    @Autowired
    private PortalSysDictRepository portalSysDictRepository;

    @Override
    public CommonResult<?> create(final String tenant, final String name, final String type, final String remark) {
        CommonResult<PortalUserGroupVO> result = businessProcessTemplate
                .process(new BusinessProcessCallback<PortalUserGroupVO>() {

                    @Override
                    public void checkParam(BusinessProcessContext context) {
                        CommonParamtersChecker.checkNotBlank(name);
                    }

                    @Override
                    public void checkBusinessInfo(BusinessProcessContext context) {
                        if (portalUserGroupRepository.findByNameAndTenant(name, tenant).size() > 0) {
                            throw new CommonBizException("用户组名称已经存在，无法保存");
                        }
                    }

                    @Override
                    public PortalUserGroupVO doBusiness(BusinessProcessContext context) {
                        PortalUserGroup po = new PortalUserGroup();

                        po.setName(name);
                        po.setRemark(remark);
                        if (StringUtils.isNotBlank(type)) {
                            po.setType(type);
                        } else {
                            po.setType(PortalUserGroupTypeEnum.DEFAULT.getCode());
                        }
                        po.setIsActive(IsActiveEnum.HAS_ACTIVE.getCode());
                        if (StringUtils.isNotBlank(tenant)) {
                            po.setTenant(tenant);
                        } else {
                            po.setTenant(currentUser.getUser().getTenantId());
                        }
                        return portalUserGroupConvert.apply(portalUserGroupRepository.save(po));
                    }

                    @Override
                    public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

                    }
                });
        return result;
    }

    @Override
    public CommonResult<?> delete(final String id) {
        return businessProcessTemplate.process(new BusinessProcessCallback<Object>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
                PortalUserGroup userGroup = portalUserGroupRepository.findOne(id);
                if (userGroup == null) {
                    throw new RuntimeException("删除失败,此用户组不存在!");
                }
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {

            }

            @Override
            public Object doBusiness(BusinessProcessContext context) {
                //删除用户组和人员关系数据
                deleteGroupUser(id);
                //删除用户组和角色关联数据
                deleteGroupRole(id);

                portalUserGroupService.updateIsActiveById(IsActiveEnum.NO_ACTIVE.getCode(), id);
                return null;
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
    }

    /**
     * 删除用户组关联的用户关系数据
     */
    private void deleteGroupUser(String id) {
        List<PortalGroupUserRel> portalGroupUserRelList = portalGroupUserRelRepository.findByGroupId(id);
        if (portalGroupUserRelList != null && !portalGroupUserRelList.isEmpty()) {
            portalGroupUserRelRepository.delete(portalGroupUserRelList);
        }
    }

    /**
     * 删除用户组关联的角色关系数据
     */
    private void deleteGroupRole(String id) {
        List<PortalGroupRoleRel> portalGroupRoleRelList = portalGroupRoleRelRepository.findByGroupId(id);
        if (portalGroupRoleRelList != null && !portalGroupRoleRelList.isEmpty()) {
            portalGroupRoleRelRepository.delete(portalGroupRoleRelList);
        }
    }

    @Override
    public CommonResult<?> update(final String id, final String name, final String type, final String remark) {
        CommonResult<PortalUserGroupVO> result = businessProcessTemplate
                .process(new BusinessProcessCallback<PortalUserGroupVO>() {

                    @Override
                    public void checkParam(BusinessProcessContext context) {
                        CommonParamtersChecker.checkNotBlank(id);
                        CommonParamtersChecker.checkNotBlank(name);
                    }

                    @Override
                    public void checkBusinessInfo(BusinessProcessContext context) {
                        PortalUserGroup po = portalUserGroupRepository.findOne(id);

                        if (po == null || IsActiveEnum.NO_ACTIVE.getCode().equals(po.getIsActive())) {
                            throw new CommonBizException("未查询到有效的用户组信息");
                        }

                        context.put("PortalUserGroup", po);
                    }

                    @Override
                    public PortalUserGroupVO doBusiness(BusinessProcessContext context) {
                        PortalUserGroup po = (PortalUserGroup) context.get("PortalUserGroup");
                        po.setName(name);
                        po.setRemark(remark);
                        if (StringUtils.isNotBlank(type)) {
                            po.setType(type);
                        } else {
                            po.setType(PortalUserGroupTypeEnum.DEFAULT.getCode());
                        }
                        return portalUserGroupConvert.apply(portalUserGroupRepository.save(po));
                    }

                    @Override
                    public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

                    }
                });
        return result;
    }

    @Override
    public CommonResult<?> allocateRole(final String id, final String selectRoles) {
        CommonResult<T> result = businessProcessTemplate.process(new BusinessProcessCallback<T>() {

            @Override
            public void checkParam(BusinessProcessContext context) {
                CommonParamtersChecker.checkNotBlank(id);
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {

            }

            @Override
            public T doBusiness(BusinessProcessContext context) {
                // 查询用户组信息
                PortalUserGroup group = portalUserGroupRepository.findOne(id);
                if (StringUtils.isNoneBlank(selectRoles)) {
                    String[] roleIdList = selectRoles.split(",");
                    for (String role : roleIdList) {
                        String roleId = role.split(":")[0];
                        String operType = role.split(":")[1];

                        if ("I".equals(operType)) {
                            PortalGroupRoleRel po = new PortalGroupRoleRel();
                            po.setRoleId(roleId);
                            po.setGroupId(id);
                            po.setTenant(group.getTenant());
                            portalGroupRoleRelRepository.save(po);
                        }

                        if ("D".equals(operType)) {
                            portalGroupRoleRelRepository.deleteByGroupIdAndRoleId(id, roleId);
                        }

                    }
                }
                return null;
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
        return result;
    }

    @Override
    public CommonResult<?> allocateUser(final String id, final String selectUsers) {
        CommonResult<T> result = businessProcessTemplate.process(new BusinessProcessCallback<T>() {

            @Override
            public void checkParam(BusinessProcessContext context) {
                CommonParamtersChecker.checkNotBlank(id);

            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {

            }

            @Override
            public T doBusiness(BusinessProcessContext context) {
                // 保存用户组与用户的关联关系
                // 查询用户组信息
                PortalUserGroup group = portalUserGroupRepository.findOne(id);
                for (PortalGroupAllocateUserForm form : JSON.parseArray(selectUsers,
                        PortalGroupAllocateUserForm.class)) {
                    if ("I".equals(form.getOperType())) {
                        PortalGroupUserRel po = new PortalGroupUserRel();
                        po.setUserId(form.getId());
                        po.setGroupId(id);
                        po.setTenant(group.getTenant());
                        portalGroupUserRelRepository.save(po);
                    }

                    if ("D".equals(form.getOperType())) {
                        portalGroupUserRelRepository.deleteByGroupIdAndUserId(id, form.getId());
                    }

                }
                return null;
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
        return result;
    }

    @Override
    public CommonResult<?> createUserGroupAndUser(final PortalGroupAndUserCreateForm form) {
        CommonResult<String> result = businessProcessTemplate.process(new BusinessProcessCallback<String>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
                CommonParamtersChecker.checkNotBlank(form.getGroupName());
                CommonParamtersChecker.checkNotBlank(form.getTenantId());
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {
                if (portalUserGroupRepository.findByNameAndTenant(form.getGroupName(), form.getTenantId()).size() > 0) {
                    throw new CommonBizException("用户组名称已经存在，无法创建");
                }
            }

            @Override
            public String doBusiness(BusinessProcessContext context) {
                String type = (form.getGroupType() == null || "".equals(form.getGroupType())) ? "default"
                        : form.getGroupType();
                List<String> userGroupTypes = portalSysDictRepository.getUserGroupType(form.getTenantId());
                if (!userGroupTypes.contains(type)) {
                    PortalSysDict sysDict = new PortalSysDict();
                    sysDict.setTypeCode("userGroupType");
                    sysDict.setCode(type);
                    sysDict.setValue(type);
                    sysDict.setTenantId(form.getTenantId());
                    portalSysDictRepository.save(sysDict);
                }
                PortalUserGroup userGroup = new PortalUserGroup();
                userGroup.setName(form.getGroupName());
                userGroup.setType(type);
                userGroup.setIsActive("Y");
                userGroup.setRemark(form.getGroupRemark());
                userGroup.setTenant(form.getTenantId());
                userGroup.setCreateUser("system");
                userGroup.setGmtCreate(new Date());
                userGroup.setGmtModified(new Date());
                PortalUserGroup portalUserGroup = portalUserGroupRepository.save(userGroup);

                if (form.getLoginIds() != null && !form.getLoginIds().isEmpty()) {
                    List<String> loginIds = Arrays.asList(form.getLoginIds().split(","));
                    List<PortalUser> portalUsers = portalUserRepository.findUserByIds(form.getTenantId(), loginIds);
                    List<PortalGroupUserRel> groupUserRels = new ArrayList<>();
                    for (PortalUser user : portalUsers) {
                        PortalGroupUserRel po = new PortalGroupUserRel();
                        po.setGroupId(portalUserGroup.getId());
                        po.setTenant(form.getTenantId());
                        po.setUserId(user.getId());
                        groupUserRels.add(po);
                    }
                    portalGroupUserRelRepository.save(groupUserRels);
                }

                return portalUserGroup.getId();
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
        return result;
    }

    @Override
    public CommonResult<?> updateUserGroupAndUser(final PortalGroupAndUserUpdateForm form) {
        CommonResult<String> result = businessProcessTemplate.process(new BusinessProcessCallback<String>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
                CommonParamtersChecker.checkNotBlank(form.getGroupId());
                CommonParamtersChecker.checkNotBlank(form.getGroupName());
                CommonParamtersChecker.checkNotBlank(form.getTenantId());
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {
                PortalUserGroup po = portalUserGroupRepository.findOne(form.getGroupId());

                if (po == null || IsActiveEnum.NO_ACTIVE.getCode().equals(po.getIsActive())) {
                    throw new CommonBizException("未查询到有效的用户组信息");
                }

                context.put("PortalUserGroup", po);
            }

            @Override
            public String doBusiness(BusinessProcessContext context) {
                String type = (form.getGroupType() == null || "".equals(form.getGroupType())) ? "default"
                        : form.getGroupType();
                List<String> userGroupTypes = portalSysDictRepository.getUserGroupType(form.getTenantId());
                if (!userGroupTypes.contains(type)) {
                    PortalSysDict sysDict = new PortalSysDict();
                    sysDict.setTypeCode("userGroupType");
                    sysDict.setValue(type);
                    sysDict.setCode(type);
                    sysDict.setTenantId(form.getTenantId());
                    portalSysDictRepository.save(sysDict);
                }
                PortalUserGroup userGroup = new PortalUserGroup();
                userGroup.setId(form.getGroupId());
                userGroup.setName(form.getGroupName());
                userGroup.setType(type);
                userGroup.setIsActive("Y");
                userGroup.setRemark(form.getGroupRemark());
                userGroup.setTenant(form.getTenantId());
                userGroup.setGmtModified(new Date());
                PortalUserGroup portalUserGroup = portalUserGroupRepository.save(userGroup);

                if (form.getLoginIds() == null || form.getLoginIds().isEmpty()) {
                    deleteGroupUser(form.getGroupId());
                } else {
                    //此处逻辑为先全部删除，再添加（这样逻辑简单）
                    deleteGroupUser(form.getGroupId());
                    List<String> loginIds = Arrays.asList(form.getLoginIds().split(","));
                    List<PortalUser> portalUsers = portalUserRepository.findUserByIds(form.getTenantId(), loginIds);
                    List<String> ids = ListUtil.filter(portalUsers, "id");
                    List<PortalGroupUserRel> groupUserRels = new ArrayList<>();
                    for (String id : ids) {
                        PortalGroupUserRel po = new PortalGroupUserRel();
                        po.setGroupId(portalUserGroup.getId());
                        po.setTenant(form.getTenantId());
                        po.setUserId(id);
                        groupUserRels.add(po);
                    }
                    portalGroupUserRelRepository.save(groupUserRels);
                }
                return "success";
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
        return result;
    }
}
