/**
 * 杭州云霁科技有限公司 http://www.idcos.com Copyright (c) 2015 All Rights Reserved.
 */

package com.idcos.enterprise.portal.manager.impl;

import java.util.*;

import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.google.common.collect.Lists;
import com.idcos.cloud.biz.common.check.CommonParamtersChecker;
import com.idcos.cloud.biz.common.check.FormChecker;
import com.idcos.cloud.core.common.BaseLinkVO;
import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.cloud.core.common.util.ListUtil;
import com.idcos.cloud.core.dal.common.page.PageUtils;
import com.idcos.cloud.core.dal.common.page.Pagination;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryCallback;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryTemplate;
import com.idcos.enterprise.portal.biz.common.utils.*;
import com.idcos.enterprise.portal.convert.*;
import com.idcos.enterprise.portal.dal.entity.*;
import com.idcos.enterprise.portal.dal.enums.*;
import com.idcos.enterprise.portal.dal.repository.*;
import com.idcos.enterprise.portal.form.PortalQueryByPageForm;
import com.idcos.enterprise.portal.form.PortalUserGroupQueryByPageForm;
import com.idcos.enterprise.portal.manager.auto.PortalUserGroupQueryManager;
import com.idcos.enterprise.portal.services.PortalUserGroupService;
import com.idcos.enterprise.portal.vo.PortalUserGroupVO;
import com.idcos.enterprise.portal.vo.PortalUserVO;
import com.idcos.enterprise.portal.biz.common.CommonBizException;
import com.idcos.enterprise.portal.biz.common.utils.CurrentUser;
import static com.idcos.enterprise.portal.UamConstant.ADMIN;

/**
 * Manager实现类
 * <p>
 * 第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖
 * </p>
 *
 * @author yanlv
 * @version v 1.1 2015-06-09 09:11:34 yanlv Exp $
 */
@Service
public class PortalUserGroupQueryManagerImpl implements PortalUserGroupQueryManager {
    private static final Logger LOGGER = LoggerFactory.getLogger(PortalUserGroupQueryManagerImpl.class);

    @Autowired
    private PortalUserGroupConvert func;

    @Autowired
    private PortalUserGroupConvert portalUserGroupConvert;

    @Autowired
    private PortalSysDictConvert portalSysDictConvert;

    @Autowired
    private PortalUserGroupInfoConvert funcInfo;

    @Autowired
    private BusinessQueryTemplate businessQueryTemplate;

    @Autowired
    private CrudUtilService crudUtilService;

    @Autowired
    private PortalRoleRepository portalRoleRepository;

    @Autowired
    private PortalPermissionRepository portalPermissionRepository;

    @Autowired
    private PortalGroupRoleRelRepository portalGroupRoleRelRepository;

    @Autowired
    private PortalGroupUserRelRepository portalGroupUserRelRepository;

    @Autowired
    private PortalRelQueryRepository portalRelQueryRepository;

    @Autowired
    private PortalUserRepository portalUserRepository;

    @Autowired
    private PortalResourceRepository portalResourceRepository;

    @Autowired
    private PortalUserGroupRepository portalUserGroupRepository;

    @Autowired
    private CurrentUser currentUser;

    @Autowired
    private PortalSysDictRepository portalSysDictRepository;

    @Autowired
    private PortalUserGroupService portalUserGroupService;

    @Override
    public CommonResult<?> queryByPage(final String offset, final String limit,
        final PortalUserGroupQueryByPageForm form, final String cnd) {
        if (!ADMIN.equals(currentUser.getUser().getLoginId())) {
            throw new CommonBizException("only admin permits");
        }
        return businessQueryTemplate.process(new BusinessQueryCallback<Pagination<PortalUserGroupVO>>() {
            @Override
            public void checkParam() {}

            @Override
            public Pagination<PortalUserGroupVO> doQuery() {
                // 如果页面没传tenantId则从token里取
                if (StringUtils.isBlank(form.getTenantId())) {
                    form.setTenantId(currentUser.getUser().getTenantId());
                }
                return portalUserGroupService.getPageList(Integer.parseInt(offset), Integer.parseInt(limit), form, cnd);
            }
        });
    }

    @Override
    public CommonResult<?> queryById(final String id) {
        return crudUtilService.findOneAndIsActive(id, funcInfo);
    }

    @Override
    public CommonResult<?> queryRolesById(final String groupId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(groupId);
            }

            @Override
            public Object doQuery() {

                // 声明变量
                BaseLinkVO vo = new BaseLinkVO();

                List<PortalRole> roleList = Lists.newArrayList();

                // 获取角色，获取权限
                for (PortalGroupRoleRel rel : portalGroupRoleRelRepository.findByGroupId(groupId)) {
                    PortalRole role = portalRoleRepository.findOne(rel.getRoleId());
                    if (IsActiveEnum.HAS_ACTIVE.getCode().equals(role.getIsActive())) {
                        roleList.add(role);
                    }
                }

                vo.putLinkObj("ROLE", PortalFilterUtil.filterRole(roleList));

                return vo;
            }
        });
    }

    @Override
    public CommonResult<?> queryUsersById(final String groupId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(groupId);
            }

            @Override
            public Object doQuery() {
                BaseLinkVO vo = new BaseLinkVO();
                List<PortalUser> userList = Lists.newArrayList();

                for (PortalGroupUserRel rel : portalGroupUserRelRepository.findByGroupId(groupId)) {
                    PortalUser user = portalUserRepository.findOne(rel.getUserId());

                    if (user != null) {
                        if (IsActiveEnum.HAS_ACTIVE.getCode().equals(user.getIsActive())) {
                            userList.add(user);
                        }
                    }
                }

                vo.putLinkObj("USER", PortalFilterUtil.filterUser(userList));
                return vo;
            }
        });
    }

    @Override
    public CommonResult<?> queryUsersByGroupId(final String groupId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(groupId);
            }

            @Override
            public Object doQuery() {
                List<PortalUserVO> userVOS = portalUserGroupService.listUsersByGroupId(groupId);
                return userVOS;
            }
        });
    }

    @Override
    public CommonResult<?> queryUserByGroupNameAndTenant(final String groupName, final String tenantId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(groupName);
                CommonParamtersChecker.checkNotBlank(tenantId);
            }

            @Override
            public Object doQuery() {
                List<PortalUserVO> userVOS =
                    portalUserGroupService.listUsersByGroupNameAndTenantId(groupName, tenantId);
                return userVOS;
            }
        });
    }

    @Override
    public CommonResult<?> queryPermissionsById(final String groupId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(groupId);
            }

            @Override
            public Object doQuery() {

                // 声明变量
                BaseLinkVO vo = new BaseLinkVO();

                List<PortalPermission> perList = Lists.newArrayList();

                // 获取权限
                for (PortalGroupRoleRel rel : portalGroupRoleRelRepository.findByGroupId(groupId)) {

                    perList.addAll(portalPermissionRepository
                        .queryAuthResByAuthObjTypeAndAuthObjId(AuthObjTypeEnum.ROLE.getCode(), rel.getRoleId()));
                }

                Map<?, List<PortalPermission>> perGroups = ListUtil.groupBy(perList, "authResType");
                List<Map<String, Object>> resultList = Lists.newArrayList();

                for (Object key : perGroups.keySet()) {
                    Map<String, Object> tmp = new HashMap<>(3);

                    PortalResource res = portalResourceRepository.findByCodeAndIsActive((String)key, "Y");

                    if (res == null) {
                        continue;
                    }
                    // 资源类型
                    tmp.put("resName", res.getName());
                    // 系统名称
                    tmp.put("systemName", res.getAppId());
                    // 权限数据
                    tmp.put("permissions", PortalFilterUtil.filterPermission(perGroups.get(key)));

                    resultList.add(tmp);
                }

                vo.putLinkObj("PERMISSION", resultList);

                return vo;
            }
        });
    }

    @Override
    public CommonResult<?> queryGroupPageByUserIdAndCnd(final PortalQueryByPageForm form) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Pagination<PortalUserGroupVO>>() {
            @Override
            public void checkParam() {
                //FormChecker.check(form); TODO：jar包待下线
            }

            @Override
            public Pagination<PortalUserGroupVO> doQuery() {

                List<PortalUserGroup> totalGroupList = PortalFilterUtil.filterGroup(
                    portalUserGroupRepository.findGroupPageByUserIdAndCnd(form.getId(), "%" + form.getCnd() + "%"));

                List<PortalUserGroup> groupList =
                    PortalFilterUtil.filterGroup(portalRelQueryRepository.findGroupPageByUserIdAndCnd(form.getId(),
                        "%" + form.getCnd() + "%", form.getPageNo(), form.getPageSize()));

                List<PortalUserGroupVO> voList = ListUtil.transform(groupList, portalUserGroupConvert);

                return PageUtils.toPagination(voList, form.getPageNo(), form.getPageSize(), totalGroupList.size());

            }
        });
    }

    @Override
    public CommonResult<?> queryGroupPageByRoleIdAndCnd(final PortalQueryByPageForm form) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Pagination<PortalUserGroupVO>>() {
            @Override
            public void checkParam() {
                //FormChecker.check(form); TODO：jar包待下线
            }

            @Override
            public Pagination<PortalUserGroupVO> doQuery() {

                List<PortalUserGroup> totalGroupList = PortalFilterUtil.filterGroup(
                    portalUserGroupRepository.findGroupPageByRoleIdAndCnd(form.getId(), "%" + form.getCnd() + "%"));

                List<PortalUserGroup> groupList =
                    PortalFilterUtil.filterGroup(portalRelQueryRepository.findGroupPageByRoleIdAndCnd(form.getId(),
                        "%" + form.getCnd() + "%", form.getPageNo(), form.getPageSize()));

                List<PortalUserGroupVO> voList = ListUtil.transform(groupList, portalUserGroupConvert);

                return PageUtils.toPagination(voList, form.getPageNo(), form.getPageSize(), totalGroupList.size());

            }
        });
    }

    @Override
    public CommonResult<?> getAccountGroupTree(final String tenantId, final String treeStyle) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(tenantId);
                CommonParamtersChecker.checkNotBlank(treeStyle);
            }

            @Override
            public Object doQuery() {
                List<PortalUserGroup> portalUserGroups = portalUserGroupRepository.findAllUserGroupByTenantId(tenantId);
                List<PortalUserGroupVO> userGroupVOList = convert2Dto(portalUserGroups);
                if (TreeStyleEnum.Z_TREE.getCode().equals(treeStyle)) {
                    return userGroupVOList;
                } else if (TreeStyleEnum.IO_TREE.getCode().equals(treeStyle)) {
                    return userGroupVOList;
                } else if (TreeStyleEnum.ALL.getCode().equals(treeStyle)) {
                    return new ArrayList<>();
                }
                return null;
            }
        });
    }

    @Override
    public CommonResult<?> getAccountGroupListByPermissionCode(final String code) {
        return businessQueryTemplate.process(new BusinessQueryCallback<List<PortalUserGroup>>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(code);
            }

            @Override
            public List<PortalUserGroup> doQuery() {
                List<String> userGroupIds = portalUserGroupRepository.getUserGroupIdsByPermissionCode(code);
                List<PortalUserGroup> portalUserGroups = new ArrayList<>();
                for (Iterator iterator = userGroupIds.iterator(); iterator.hasNext();) {
                    String userGroupId = (String)iterator.next();
                    PortalUserGroup userGroup =
                        portalUserGroupRepository.findUserGroupByUserGroupIdAndIsActive(userGroupId, "Y");
                    if (userGroup == null) {
                        continue;
                    }
                    portalUserGroups.add(userGroup);
                }
                return portalUserGroups;
            }
        });
    }

    @Override
    public CommonResult<?> getUserGroupType(final String tenantId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<List<PortalSysDict>>() {
            @Override
            public void checkParam() {

            }

            @Override
            public List<PortalSysDict> doQuery() {
                List<PortalSysDict> userGroupType = portalSysDictRepository.getUserGroupTypeAndValue(tenantId);
                return userGroupType;
            }
        });
    }

    @Override
    public CommonResult<?> queryGroupByType(final String type) {
        return businessQueryTemplate.process(new BusinessQueryCallback<List<PortalUserGroupVO>>() {
            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(type);
            }

            @Override
            public List<PortalUserGroupVO> doQuery() {
                List<PortalUserGroup> userGroup = portalUserGroupRepository.findUserGroupByType(type);
                return convert2Dto(userGroup);
            }
        });
    }

    @Override
    public CommonResult<?> queryGroupByTypeAndTenantId(final String type, final String tenantId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<List<PortalUserGroupVO>>() {
            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(type);
            }

            @Override
            public List<PortalUserGroupVO> doQuery() {
                List<PortalUserGroup> userGroup =
                    portalUserGroupRepository.findUserGroupByTypeAndTenantId(type, tenantId);
                return convert2Dto(userGroup);
            }
        });
    }

    @Override
    public CommonResult<?> queryGroupByTypeAndUserId(final String type, final String userId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<List<PortalUserGroupVO>>() {
            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(type);
                CommonParamtersChecker.checkNotBlank(userId);
            }

            @Override
            public List<PortalUserGroupVO> doQuery() {
                List<PortalUserGroupVO> userGroup = portalUserGroupService.queryGroupByTypeAndUserId(type, userId);
                return userGroup;
            }
        });
    }

    @Override
    public CommonResult<?> queryGroupByLoginIdAndTenantId(final String loginId, final String tenantId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<List<PortalUserGroupVO>>() {
            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(loginId);
                CommonParamtersChecker.checkNotBlank(tenantId);
            }

            @Override
            public List<PortalUserGroupVO> doQuery() {
                PortalUser portalUser = portalUserRepository.findPortalUserById(tenantId, loginId);
                List<PortalUserGroupVO> userGroup = portalUserGroupService.listUserGroupByUserId(portalUser.getId());
                return userGroup;
            }
        });
    }

    /**
     * 转化为PortalUserGroupVO
     *
     * @param portalUserGroups
     * @return
     */
    private List<PortalUserGroupVO> convert2Dto(List<PortalUserGroup> portalUserGroups) {
        List<PortalUserGroupVO> dtos = new ArrayList<>();
        for (PortalUserGroup userGroup : portalUserGroups) {
            PortalUserGroupVO userGroupVO = new PortalUserGroupVO();
            userGroupVO.setId(userGroup.getId());
            userGroupVO.setName(userGroup.getName());
            userGroupVO.setIsActive(userGroup.getIsActive());
            userGroupVO.setTenantId(userGroup.getTenant());
            userGroupVO.setRemark(userGroup.getRemark());
            userGroupVO.setType(userGroup.getType());
            dtos.add(userGroupVO);
        }
        return dtos;
    }

}