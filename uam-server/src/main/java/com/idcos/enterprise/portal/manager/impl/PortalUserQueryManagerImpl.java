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
import com.idcos.cloud.core.common.util.StringUtil;
import com.idcos.cloud.core.dal.common.page.PageUtils;
import com.idcos.cloud.core.dal.common.page.Pagination;
import com.idcos.enterprise.portal.biz.common.*;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryCallback;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryTemplate;
import com.idcos.enterprise.portal.biz.common.utils.*;
import com.idcos.enterprise.portal.convert.PortalSysDictConvert;
import com.idcos.enterprise.portal.dal.entity.*;
import com.idcos.enterprise.portal.dal.enums.*;
import com.idcos.enterprise.portal.dal.repository.*;
import com.idcos.enterprise.portal.form.PortalQueryByPageForm;
import com.idcos.enterprise.portal.form.PortalUserQueryPageListForm;
import com.idcos.enterprise.portal.manager.auto.PortalUserQueryManager;
import com.idcos.enterprise.portal.manager.common.CommonManager;
import com.idcos.enterprise.portal.services.PortalUserService;
import com.idcos.enterprise.portal.vo.*;
import com.idcos.enterprise.portal.web.GlobalValue;
import com.idcos.enterprise.portal.biz.common.CommonBizException;
import com.idcos.enterprise.portal.biz.common.utils.CurrentUser;
import static com.idcos.enterprise.portal.UamConstant.ADMIN;
import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;

/**
 * Manager实现类
 * <p>
 * 第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖
 * </p>
 *
 * @author yanlv
 * @version v 1.1 2015-06-09 09:26:24 yanlv Exp $
 */
@Service
public class PortalUserQueryManagerImpl implements PortalUserQueryManager {
    private static final Logger LOGGER = LoggerFactory.getLogger(PortalUserQueryManagerImpl.class);

    @Autowired
    private PortalGroupUserRelRepository portalGroupUserRelRepository;

    @Autowired
    private BusinessQueryTemplate businessQueryTemplate;

    @Autowired
    private PortalUserGroupRepository portalUserGroupRepository;

    @Autowired
    private PortalRoleRepository portalRoleRepository;

    @Autowired
    private PortalUserRepository portalUserRepository;

    @Autowired
    private PortalRelQueryRepository portalRelQueryRepository;

    @Autowired
    private PortalGroupRoleRelRepository portalGroupRoleRelRepository;

    @Autowired
    private PortalResourceRepository portalResourceRepository;

    @Autowired
    private PortalUserService portalUserService;

    @Autowired
    private PortalDeptRepository portalDeptRepository;

    @Autowired
    private PortalTokenRepository portalTokenRepository;

    @Autowired
    private PortalSysDictRepository portalSysDictRepository;

    @Autowired
    private PortalSysDictConvert portalSysDictConvert;

    @Autowired
    private GlobalValue globalValue;

    @Autowired
    private CommonManager commonManager;

    @Autowired
    private CurrentUser currentUser; 

    @Override
    public CommonResult<?> queryPageList(final PortalUserQueryPageListForm form) {
        if (!ADMIN.equals(currentUser.getUser().getLoginId())) {
            throw new CommonBizException("only admin permits");
        }
        return businessQueryTemplate.process(new BusinessQueryCallback<Pagination<PortalUserVO>>() {
            @Override
            public void checkParam() {
                if (form.getPageNo() <= 0) {
                    form.setPageNo(1);
                }
                if (form.getPageSize() <= 0) {
                    form.setPageSize(10);
                }
            }

            @Override
            public Pagination<PortalUserVO> doQuery() {
                return portalUserService.getPageList(form.getTenantId(), form.getDeptId(), form.getName(),
                    form.getPageNo(), form.getPageSize());
            }
        });
    }

    @Override
    public CommonResult<?> queryById(final String userId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<PortalUserVO>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(userId);
            }

            @Override
            public PortalUserVO doQuery() {
                PortalUser portalUser = portalUserRepository.findOne(userId);
                PortalUserVO userVo = convertUserToUserVO(portalUser);

                {
                    List<PortalGroupUserRel> list = portalGroupUserRelRepository.findByUserId(userId);
                    String[] selGroups = new String[list.size()];

                    for (int i = 0; i < list.size(); i++) {
                        selGroups[i] = list.get(i).getGroupId();
                    }
                    userVo.setSelGroups(selGroups);
                }

                return userVo;
            }
        });
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

                for (PortalGroupUserRel rel : portalGroupUserRelRepository.findByUserId(id)) {
                    PortalUserGroup group = portalUserGroupRepository.findOne(rel.getGroupId());
                    if (group != null) {
                        if (IsActiveEnum.HAS_ACTIVE.getCode().equals(group.getIsActive())) {
                            groupList.add(group);
                        }
                    }
                }

                vo.putLinkObj("GROUP", groupList);

                return vo;
            }
        });
    }

    @Override
    public CommonResult<?> queryRolesById(final String id) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(id);
            }

            @Override
            public Object doQuery() {
                BaseLinkVO vo = new BaseLinkVO();
                List<PortalRole> roleList = Lists.newArrayList();

                for (PortalGroupUserRel userRel : portalGroupUserRelRepository.findByUserId(id)) {

                    for (PortalGroupRoleRel roleRel : portalGroupRoleRelRepository
                        .findByGroupId(userRel.getGroupId())) {
                        PortalRole role = portalRoleRepository.findOne(roleRel.getRoleId());

                        if (IsActiveEnum.HAS_ACTIVE.getCode().equals(role.getIsActive())) {
                            roleList.add(role);
                        }
                    }
                }
                vo.putLinkObj("ROLE", PortalFilterUtil.filterRole(roleList));

                return vo;
            }
        });
    }

    @Override
    public CommonResult<?> queryPermissionsById(final String id, final String appId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(id);
            }

            @Override
            public Object doQuery() {
                BaseLinkVO vo = new BaseLinkVO();
                List<PortalPermission> perList = Lists.newArrayList();

                List<String> roleIds = commonManager.getRoleIdsById(id);
                perList =
                    portalUserService.getPermissonByAuthObjTypeAndAuthObjId(AuthObjTypeEnum.ROLE.getCode(), roleIds);

                // 根据appId进行过滤,得到某个系统的权限
                List<PortalPermission> portalPermissions = Lists.newArrayList();
                if (StringUtil.isNotBlank(appId)) {
                    String app = appId.toUpperCase();
                    for (PortalPermission portalPermission : perList) {
                        if (portalPermission.getAuthResType().contains(app)) {
                            portalPermissions.add(portalPermission);
                        }
                    }
                }

                Map<?, List<PortalPermission>> perGroups;
                if (portalPermissions.isEmpty() || portalPermissions == null) {
                    perGroups = ListUtil.groupBy(perList, "authResType");
                } else {
                    perGroups = ListUtil.groupBy(portalPermissions, "authResType");
                }

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
    public CommonResult<?> queryPermissionsByUserIdAndAppId(final String userId, final String appId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(userId);
                CommonParamtersChecker.checkNotBlank(appId);
            }

            @Override
            public Object doQuery() {
                List<String> roleIds = commonManager.getRoleIdsById(userId);
                List<PortalPermission> perminssionList =
                    portalUserService.queryPermissionsByRoleIdsAndAppId(roleIds, appId);

                Map<?, List<PortalPermission>> perGroups = new HashMap<>(2);
                if (perminssionList != null && perminssionList.size() != 0) {
                    perGroups = ListUtil.groupBy(perminssionList, "authResType");
                }

                List<Map<String, Object>> resultList = Lists.newArrayList();
                for (Object key : perGroups.keySet()) {
                    Map<String, Object> tmp = new HashMap<>(2);
                    PortalResource res = portalResourceRepository.findByCodeAndIsActive((String)key, "Y");
                    if (res == null) {
                        continue;
                    }
                    // 权限数据
                    tmp.put("permissions", PortalFilterUtil.filterPermission(perGroups.get(key)));
                    resultList.add(tmp);
                }

                return resultList;
            }
        });
    }

    @Override
    public CommonResult<?> queryUserPageByGroupIdAndCnd(final PortalQueryByPageForm form) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Pagination<PortalUserVO>>() {
            @Override
            public void checkParam() {
                //FormChecker.check(form); TODO：jar包待下线
            }

            @Override
            public Pagination<PortalUserVO> doQuery() {

                List<PortalUser> totalUserList = PortalFilterUtil.filterUser(
                    portalUserRepository.findUserPageByGroupIdAndCnd(form.getId(), "%" + form.getCnd() + "%"));

                List<PortalUser> userList =
                    PortalFilterUtil.filterUser(portalRelQueryRepository.findUserPageByGroupIdAndCnd(form.getId(),
                        "%" + form.getCnd() + "%", form.getPageNo(), form.getPageSize()));

                List<PortalUserVO> voList = new ArrayList<>();
                Iterator<PortalUser> portalUserIterator = userList.iterator();
                while (portalUserIterator.hasNext()) {
                    PortalUser portalUser = portalUserIterator.next();
                    PortalUserVO userVo = convertUserToUserVO(portalUser);
                    voList.add(userVo);
                }
                return PageUtils.toPagination(voList, form.getPageNo(), form.getPageSize(), totalUserList.size());

            }
        });
    }

    @Override
    public CommonResult<?> queryUserPageByRoleIdAndCnd(final PortalQueryByPageForm form) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Pagination<PortalUserVO>>() {
            @Override
            public void checkParam() {
                //FormChecker.check(form); TODO：jar包待下线
            }

            @Override
            public Pagination<PortalUserVO> doQuery() {
                return portalUserService.listUserByRoleForm(form);
            }
        });

    }

    @Override
    public CommonResult<?> getAccountCount(final String tenantId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Long>() {

            @Override
            public void checkParam() {}

            @Override
            public Long doQuery() {
                Long accountCount;
                if (StringUtil.isNotBlank(tenantId)) {
                    accountCount = portalUserRepository.getAccountCount(tenantId);
                } else {
                    accountCount = portalUserRepository.getAccountCount();
                }
                return accountCount;
            }
        });
    }

    @Override
    public CommonResult<?> getAllAcount() {
        return businessQueryTemplate.process(new BusinessQueryCallback<List<PortalUserVO>>() {

            @Override
            public void checkParam() {}

            @Override
            public List<PortalUserVO> doQuery() {
                return portalUserService.listUserVO();
            }
        });
    }

    @Override
    public CommonResult<?> getAllAcount(final String tenantId) {
        if (!ADMIN.equals(currentUser.getUser().getLoginId())) {
            throw new CommonBizException("only admin permits");
        }
        return businessQueryTemplate.process(new BusinessQueryCallback<List<PortalUserVO>>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(tenantId);
            }

            @Override
            public List<PortalUserVO> doQuery() {
                return portalUserService.listUserVOByTenantId(tenantId);
            }
        });
    }

    @Override
    public CommonResult<?> getAcountByDeptId(final String tenantId, final String deptId, final String recurse) {
        return businessQueryTemplate.process(new BusinessQueryCallback<List<PortalUserVO>>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(tenantId);
                CommonParamtersChecker.checkNotBlank(deptId);
            }

            @Override
            public List<PortalUserVO> doQuery() {
                String[] deptIds = deptId.split(",");
                List<PortalUserVO> userVOs = new ArrayList<>();
                List<String> idList = new ArrayList<>();

                if (TrueFalse.TRUE.getCode().equals(recurse)) {
                    // 获取所有的子部门和主部门的id集合
                    for (String deptId : deptIds) {
                        idList.add(deptId);
                        for (PortalDept dept : getChildrenList(deptId)) {
                            idList.add(dept.getId());
                        }
                    }
                } else {
                    for (String deptId : deptIds) {
                        idList.add(deptId);
                    }
                }

                List<PortalUser> portalUsers = portalUserRepository.findPortalUserByDeptIdIn(tenantId, idList);

                for (PortalUser portalUser : portalUsers) {
                    PortalUserVO userVO = convertUserToUserVO(portalUser);
                    userVOs.add(userVO);
                }

                return userVOs;
            }
        });
    }

    /**
     * 获取本部门下及子部门id的集合
     */
    private List<PortalDept> getChildrenList(String deptId) {
        List<PortalDept> allChildrenList = new LinkedList<>();
        // 先得到直接的下级部门
        List<PortalDept> depts = portalDeptRepository.findByParentId(deptId);
        for (PortalDept dept : depts) {
            allChildrenList.add(dept);
            // 递归查询,子部门下的子部门
            List<PortalDept> tmp = getChildrenList(dept.getId().toString());
            for (PortalDept ciClass1 : tmp) {
                allChildrenList.add(ciClass1);
            }
        }
        return allChildrenList;

    }

    @Override
    public CommonResult<?> getByAccountNo(final String tenantId, final String accountNo) {
        return businessQueryTemplate.process(new BusinessQueryCallback<PortalUserVO>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(tenantId);
                CommonParamtersChecker.checkNotBlank(accountNo);
            }

            @Override
            public PortalUserVO doQuery() {
                PortalUser portalUser = portalUserRepository.findPortalUserById(tenantId, accountNo);
                PortalUserVO userVO = convertUserToUserVO(portalUser);
                return userVO;
            }
        });
    }

    @Override
    public CommonResult<?> getUserById(final String id) {
        return businessQueryTemplate.process(new BusinessQueryCallback<PortalUserVO>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(id);
            }

            @Override
            public PortalUserVO doQuery() {
                PortalUser portalUser = portalUserRepository.findOne(id);
                PortalUserVO userVO = convertUserToUserVO(portalUser);
                return userVO;
            }
        });
    }

    @Override
    public CommonResult<?> getByAccountNos(final String tenantId, final List<String> accountNoList) {
        return businessQueryTemplate.process(new BusinessQueryCallback<List<PortalUserVO>>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(tenantId);
            }

            @Override
            public List<PortalUserVO> doQuery() {
                List<PortalUser> portalUsers = portalUserRepository.findUserByIds(tenantId, accountNoList);
                List<PortalUserVO> userVOs = new ArrayList<>();
                for (PortalUser portalUser : portalUsers) {
                    PortalUserVO userVO = convertUserToUserVO(portalUser);
                    userVOs.add(userVO);
                }
                return userVOs;
            }
        });
    }

    @Override
    public CommonResult<?> getLastLoginTime() {
        return businessQueryTemplate.process(new BusinessQueryCallback<Long>() {

            @Override
            public void checkParam() {}

            @Override
            public Long doQuery() {
                Date lastLoginTime = portalUserRepository.getLastLoginTime();
                return lastLoginTime.getTime();
            }
        });

    }

    private PortalUserVO convertUserToUserVO(PortalUser portalUser) {
        PortalUserVO userVO = new PortalUserVO();
        userVO.setId(portalUser.getId());
        userVO.setName(portalUser.getName());
        userVO.setLoginId(portalUser.getLoginId());
        userVO.setDeptId(portalUser.getDeptId());
        userVO.setTenantId(portalUser.getTenantId());
        userVO.setMobile1(portalUser.getMobile1());
        userVO.setMobile2(portalUser.getMobile2());
        userVO.setEmployeeType(portalUser.getEmployeeType());
        userVO.setOfficeTel1(portalUser.getOfficeTel1());
        userVO.setOfficeTel2(portalUser.getOfficeTel2());
        userVO.setStatus(portalUser.getStatus());
        userVO.setWeixin(portalUser.getWeixin());
        userVO.setEmail(portalUser.getEmail());
        userVO.setRemark(portalUser.getRemark());
        userVO.setTitle(portalUser.getTitle());
        userVO.setRtx(portalUser.getRtx());
        PortalDept portalDept = portalDeptRepository.findByDeptId(portalUser.getDeptId());
        if (portalDept != null) {
            userVO.setDeptName(portalDept.getDisplayName());
        }
        return userVO;
    }

    @Override
    public CommonResult<?> authInfo(final String id, final String appId, final String token) {
        return businessQueryTemplate.process(new BusinessQueryCallback<AuthInfoVO>() {
            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(appId);
                CommonParamtersChecker.checkNotBlank(token);
            }

            @Override
            public AuthInfoVO doQuery() {
                String userId = id;
                if (StringUtils.isBlank(userId)) {
                    Claims claims =
                        Jwts.parser().setSigningKey(globalValue.getSecretKey()).parseClaimsJws(token).getBody();

                    userId = (String)claims.get("userId");
                }
                checkToken(userId, token);

                return getAuthInfoByUserIdAndAppId(userId, appId);
            }
        });
    }

    private void checkToken(String userId, String token) {
        PortalUser portalUser = portalUserRepository.findById(userId);
        if (portalUser == null) {
            throw new AuthInfoException(ResultCode.STATUS_ERROR, "userId为" + userId + "的用户不存在");
        }

        PortalToken portalToken = portalTokenRepository.queryTokenByNameAndTokenCrc(token, CrcUtil.crc(token));
        if (portalToken == null) {
            throw new AuthInfoException(ResultCode.STATUS_ERROR,
                "用户TenantId：" + portalUser.getTenantId() + "，LoginId：" + portalUser.getLoginId() + "不存在值相同的token。");
        }
        if (!IsActiveEnum.HAS_ACTIVE.getCode().equals(portalToken.getIsActive())) {
            throw new AuthInfoException(ResultCode.AUTH_FAIL, "token被禁用！！！");
        }
    }

    @Override
    public CommonResult<?> getAuthInfoByLoginIdAndPw(final String loginId, final String password, final String tenantId,
        final String appId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<AuthInfoVO>() {
            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(loginId);
                CommonParamtersChecker.checkNotBlank(password);
            }

            @Override
            public AuthInfoVO doQuery() {
                // 验证用户、租户和密码是否正确
                PortalUser portalUser = portalUserRepository.findPortalUserById(tenantId, loginId);
                if (portalUser == null) {
                    throw new CommonBizException(ResultCode.USER_NOT_EXIST, "帐号不存在，请检查！！！");
                }

                // 验证密码是否匹配
                String plainText = "";
                try {
                    String cipherText = portalUser.getPassword();
                    byte[] salt = Base64Util.decode(portalUser.getSalt());
                    // 解密时的解密用的密码是生成时的密码， 也就是用户的ID。
                    plainText = PasswordUtil.decrypt(cipherText, portalUser.getId(), salt);

                } catch (Exception ex) {
                    LOGGER.error(String.format("解密用户【" + portalUser.getLoginId() + "】的登录密码时出错"));
                    throw new CommonBizException(ResultCode.PASSWORD_NOT_CORRECT,
                        String.format("解密用户【" + portalUser.getLoginId() + "】的登录密码时出错"));

                }

                if (!plainText.equals(password)) {
                    throw new CommonBizException(ResultCode.PASSWORD_NOT_CORRECT, "用户" + loginId + " 的密码校验错误！！！");
                }
                return getAuthInfoByUserIdAndAppId(portalUser.getId(), appId);
            }
        });
    }

    private AuthInfoVO getAuthInfoByUserIdAndAppId(String userId, String appId) {
        AuthInfoVO authInfoVO = portalUserService.findOne(userId);
        List<String> roleIds = commonManager.getRoleIdsById(userId);
        authInfoVO.setRoleIds(roleIds);

        List<PortalPermissionVO> permissions = portalUserService.getPermissionsByroleIds(appId, roleIds);
        Map<String, List<String>> permissionMap = new HashMap<>(4);

        for (PortalPermissionVO permission : permissions) {
            String authResType = permission.getAuthResType();
            if (!permissionMap.containsKey(authResType)) {
                List<String> valueList = new ArrayList<>();
                permissionMap.put(authResType, valueList);
            }
            if (!permissionMap.get(authResType).contains(permission.getAuthResId())) {
                permissionMap.get(authResType).add(permission.getAuthResId());
            }
        }

        authInfoVO.setPermissions(permissionMap);
        return authInfoVO;
    }

    @Override
    public CommonResult<?> getUserStatus() {
        return businessQueryTemplate.process(new BusinessQueryCallback<List<PortalSysDictVO>>() {
            @Override
            public void checkParam() {}

            @Override
            public List<PortalSysDictVO> doQuery() {
                List<PortalSysDict> userStatusList = portalSysDictRepository.getUserStatus();
                if (userStatusList == null || userStatusList.size() == 0) {
                    return Collections.emptyList();
                }
                return ListUtil.transform(userStatusList, portalSysDictConvert);
            }
        });
    }
}