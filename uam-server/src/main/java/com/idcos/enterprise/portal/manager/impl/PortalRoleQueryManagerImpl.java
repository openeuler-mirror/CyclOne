/**
 * 杭州云霁科技有限公司 http://www.idcos.com Copyright (c) 2015 All Rights Reserved.
 */

package com.idcos.enterprise.portal.manager.impl;

import java.util.*;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.google.common.collect.Lists;
import com.idcos.cloud.biz.common.check.CommonParamtersChecker;
import com.idcos.cloud.biz.common.check.FormChecker;
import com.idcos.cloud.core.common.BaseLinkVO;
import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.cloud.core.common.util.ListUtil;
import com.idcos.cloud.core.dal.common.page.*;
import com.idcos.enterprise.portal.biz.common.CommonBizException;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryCallback;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryTemplate;
import com.idcos.enterprise.portal.biz.common.utils.CrudUtilService;
import com.idcos.enterprise.portal.biz.common.utils.PortalFilterUtil;
import com.idcos.enterprise.portal.convert.PortalRoleConvert;
import com.idcos.enterprise.portal.convert.PortalRoleInfoConvert;
import com.idcos.enterprise.portal.dal.entity.*;
import com.idcos.enterprise.portal.dal.enums.AuthObjTypeEnum;
import com.idcos.enterprise.portal.dal.enums.IsActiveEnum;
import com.idcos.enterprise.portal.dal.repository.*;
import com.idcos.enterprise.portal.form.PortalQueryByPageForm;
import com.idcos.enterprise.portal.form.PortalRoleQueryByPageForm;
import com.idcos.enterprise.portal.manager.auto.PortalRoleQueryManager;
import com.idcos.enterprise.portal.services.PortalUserService;
import com.idcos.enterprise.portal.vo.PortalRoleVO;
import com.idcos.enterprise.portal.biz.common.utils.CurrentUser;
import static com.idcos.enterprise.portal.UamConstant.ADMIN;
/**
 * Manager实现类
 * <p>
 * 第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖
 * </p>
 *
 * @author yanlv
 * @version v 1.1 2015-06-06 10:24:34 yanlv Exp $
 */
@Service
public class PortalRoleQueryManagerImpl implements PortalRoleQueryManager {
    @Autowired
    private PortalRoleConvert func;

    @Autowired
    private PortalRoleInfoConvert funcInfo;

    @Autowired
    private CrudUtilService crudUtilService;

    @Autowired
    private BusinessQueryTemplate businessQueryTemplate;

    @Autowired
    private PortalUserGroupRepository portalUserGroupRepository;

    @Autowired
    private PortalRoleRepository portalRoleRepository;

    @Autowired
    private PortalGroupRoleRelRepository portalGroupRoleRelRepository;

    @Autowired
    private PortalGroupUserRelRepository portalGroupUserRelRepository;

    @Autowired
    private PortalPermissionRepository portalPermissionRepository;

    @Autowired
    private PortalUserRepository portalUserRepository;

    @Autowired
    private PortalRelQueryRepository portalRelQueryRepository;

    @Autowired
    private PortalResourceRepository portalResourceRepository;

    @Autowired
    private PortalUserService portalUserService;

    @Autowired
    private CurrentUser currentUser;  

    @Override
    public CommonResult<?> queryByPage(final String offset, final String limit, final PortalRoleQueryByPageForm form) {
        if (!ADMIN.equals(currentUser.getUser().getLoginId())) {
            throw new CommonBizException("only admin permits");
        }
        return crudUtilService.query(PortalRole.class, form, new PageForm(offset, limit, null, null), func, true);
    }

    @Override
    public CommonResult<?> queryById(final String id) {
        return crudUtilService.findOneAndIsActive(id, funcInfo);
    }

    @Override
    public CommonResult<?> queryGroupsById(final String id) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(id);
            }

            @Override
            public Object doQuery() {
                BaseLinkVO vo = new BaseLinkVO();
                List<PortalUserGroup> groupList = Lists.newArrayList();

                for (PortalGroupRoleRel rel : portalGroupRoleRelRepository.findByRoleId(id)) {
                    PortalUserGroup group = portalUserGroupRepository.findOne(rel.getGroupId());
                    if (IsActiveEnum.HAS_ACTIVE.getCode().equals(group.getIsActive())) {
                        groupList.add(group);
                    }
                }

                vo.putLinkObj("GROUP", groupList);
                return vo;
            }
        });
    }

    @Override
    public CommonResult<?> queryPermissionsById(final String id) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(id);
            }

            @Override
            public Object doQuery() {
                BaseLinkVO vo = new BaseLinkVO();
                List<PortalPermission> perList = Lists.newArrayList();

                perList.addAll(portalPermissionRepository
                    .queryAuthResByAuthObjTypeAndAuthObjId(AuthObjTypeEnum.ROLE.getCode(), id));

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
    public CommonResult<?> queryUsersById(final String id) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {
            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(id);
            }

            @Override
            public Object doQuery() {
                BaseLinkVO vo = new BaseLinkVO();
                vo.putLinkObj("USER", portalUserService.listUserByRoleId(id));
                return vo;
            }
        });
    }

    @Override
    public CommonResult<?> queryRolePageByUserIdAndCnd(final PortalQueryByPageForm form) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Pagination<PortalRoleVO>>() {
            @Override
            public void checkParam() {
                //FormChecker.check(form); TODO：jar包待下线
            }

            @Override
            public Pagination<PortalRoleVO> doQuery() {
                Pagination<PortalRoleVO> portalRoleVOPagination = portalUserService.listRoleByUserForm(form);
                return portalRoleVOPagination;
            }
        });
    }

    @Override
    public CommonResult<?> queryRolePageByGroupIdAndCnd(final PortalQueryByPageForm form) {
        if (!ADMIN.equals(currentUser.getUser().getLoginId())) {
            throw new CommonBizException("only admin permits");
        }
        return businessQueryTemplate.process(new BusinessQueryCallback<Pagination<PortalRoleVO>>() {
            @Override
            public void checkParam() {
                //FormChecker.check(form); TODO：jar包待下线
            }

            @Override
            public Pagination<PortalRoleVO> doQuery() {
                List<PortalRole> totalRoleList =
                    portalRoleRepository.findRolePageByGroupIdAndCnd(form.getId(), "%" + form.getCnd() + "%");

                List<PortalRole> roleList = portalRelQueryRepository.findRolePageByGroupId(form.getId(),
                    "%" + form.getCnd() + "%", form.getPageNo(), form.getPageSize());

                List<PortalRoleVO> voList = ListUtil.transform(roleList, func);

                return PageUtils.toPagination(voList, form.getPageNo(), form.getPageSize(), totalRoleList.size());

            }
        });
    }

    @Override
    public CommonResult<?> queryByTenantId(final String tenantId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<List<PortalRoleVO>>() {
            @Override
            public void checkParam() {}

            @Override
            public List<PortalRoleVO> doQuery() {
                List<PortalRole> portalRoles = portalRoleRepository.findByIsActive("Y");
                return ListUtil.transform(portalRoles, func);
            }
        });
    }

    @Override
    public CommonResult<?> queryByAccountNoAndTenantId(final String accountNo, final String tenantId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<List<PortalRoleVO>>() {
            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(accountNo);
                CommonParamtersChecker.checkNotBlank(tenantId);
            }

            @Override
            public List<PortalRoleVO> doQuery() {

                // 查询用户是否存在
                PortalUser portalUser = portalUserRepository.findPortalUserByAccountId(tenantId, accountNo);

                if (portalUser == null) {

                    throw new CommonBizException("用户信息不存在，请检查");
                }

                List<PortalRole> portalRoles = new ArrayList<>();

                for (PortalGroupUserRel userRel : portalGroupUserRelRepository.findByUserId(portalUser.getId())) {

                    for (PortalGroupRoleRel roleRel : portalGroupRoleRelRepository
                        .findByGroupId(userRel.getGroupId())) {

                        PortalRole portalRole = portalRoleRepository.findOne(roleRel.getRoleId());
                        if (portalRole != null) {
                            portalRoles.add(portalRole);
                        }

                    }
                }

                return ListUtil.transform(portalRoles, func);

            }
        });
    }
}