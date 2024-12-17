/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.services;

import com.google.common.collect.Lists;
import com.idcos.cloud.core.dal.common.page.Pagination;
import com.idcos.enterprise.jooq.Tables;
import com.idcos.enterprise.jooq.tables.*;
import com.idcos.enterprise.jooq.tables.records.PortalUserRecord;
import com.idcos.enterprise.portal.dal.enums.IsActiveEnum;
import com.idcos.enterprise.portal.form.PortalQueryByPageForm;
import com.idcos.enterprise.portal.form.PortalUserQueryPageListForm;
import com.idcos.enterprise.portal.vo.AuthInfoVO;
import com.idcos.enterprise.portal.vo.PortalPermissionVO;
import com.idcos.enterprise.portal.vo.PortalRoleVO;
import com.idcos.enterprise.portal.vo.PortalUserVO;
import org.apache.commons.lang3.StringUtils;
import org.jooq.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.thymeleaf.util.ListUtils;

import java.util.*;

import static com.idcos.enterprise.jooq.Tables.*;

/**
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月02 下午8:30 souakiragen Exp $
 */
@Service
public class PortalUserService {
    @Autowired
    private DSLContext dslContext;

    public Pagination<PortalUserVO> getPageList(String teanntId, String deptId, String cnd, int pageNo, int pageSize) {
        PortalUser user = PORTAL_USER;
        PortalDept dept = Tables.PORTAL_DEPT;
        PortalUserGroup userGroup = Tables.PORTAL_USER_GROUP;
        PortalGroupUserRel guRel = Tables.PORTAL_GROUP_USER_REL;

        List<Field<?>> fields = Lists.newArrayList(user.fields());
        Field deptNameField = dept.DISPLAY_NAME.as("DEPT_NAME");
        fields.add(deptNameField);
        SelectJoinStep<Record> select = dslContext.select(fields).from(user).leftJoin(dept)
                .on(user.DEPT_ID.eq(dept.ID));
        Condition condition = user.TENANT_ID.eq(teanntId);
        if (StringUtils.isNotBlank(deptId)) {
            condition = condition.and(user.DEPT_ID.eq(deptId));
        }
        if (StringUtils.isNotBlank(cnd)) {
            condition = condition
                    .and(user.NAME.like("%" + cnd + "%").or(user.LOGIN_ID.like("%" + cnd + "%")));
        }
        condition = condition.and(user.IS_ACTIVE.eq(IsActiveEnum.HAS_ACTIVE.getCode()));
        //默认按LOGIN_ID降序展示
        select.where(condition).orderBy(user.LOGIN_ID.asc());
        int count = dslContext.fetchCount(select);

        if (count == 0) {
            return new Pagination<>(pageNo, pageSize, 0, Lists.<PortalUserVO>newArrayList());
        }
        List<PortalUserVO> portalUserVOS = select.limit((pageNo - 1) * pageSize, pageSize)
                .fetch().into(PortalUserVO.class);

        //查询每个用户所属的用户组集合，放入PortalUserVO中
        for (PortalUserVO portalUserVO : portalUserVOS) {
            List<String> selGroupList = dslContext.selectDistinct(userGroup.NAME).from(userGroup).leftJoin(guRel)
                    .on(userGroup.ID.eq(guRel.GROUP_ID)).where(userGroup.IS_ACTIVE.eq(IsActiveEnum.HAS_ACTIVE.getCode()))
                    .and(guRel.USER_ID.eq(portalUserVO.getId())).fetch().into(String.class);
            String[] selGroups = new String[selGroupList.size()];
            selGroupList.toArray(selGroups);
            portalUserVO.setSelGroups(selGroups);
        }

        return new Pagination<>(pageNo, pageSize, count, portalUserVOS);
    }

    public void updateIsActionById(String isAction, String id) {
        PortalUser portalUser = PORTAL_USER;
        UpdateSetMoreStep<PortalUserRecord> set = dslContext.update(portalUser).set(portalUser.IS_ACTIVE, isAction);
        Condition condition = portalUser.ID.eq(id);
        set.where(condition);
        set.execute();
    }

    public List<com.idcos.enterprise.portal.dal.entity.PortalPermission> getPermissonByAuthObjTypeAndAuthObjId(String authObjType,
                                                                                                               List<String> authObjIds) {
        PortalPermission portalPermission = Tables.PORTAL_PERMISSION;
        /**
         * 查询所有角色拥有的权限并去重
         */
        PortalPermission distinctAs = portalPermission.as("distinctAs");
        SelectJoinStep<Record1<String>> idSelect = dslContext.select(distinctAs.ID.max()).from(distinctAs);
        Condition condition = distinctAs.AUTH_OBJ_TYPE.eq(authObjType);
        condition = condition.and(distinctAs.AUTH_OBJ_ID.in(authObjIds));
        idSelect.where(condition).groupBy(distinctAs.AUTH_RES_ID);
        PortalPermission prAs = portalPermission.as("prAs");
        SelectJoinStep<Record> select = dslContext.select(prAs.fields()).from(prAs);
        select.where(prAs.ID.in(idSelect));
        List<com.idcos.enterprise.portal.dal.entity.PortalPermission> portalPermissionList = select.fetch()
                .into(com.idcos.enterprise.portal.dal.entity.PortalPermission.class);
        if (ListUtils.isEmpty(portalPermissionList)) {
            return Lists.newArrayList();
        }
        return portalPermissionList;
    }

    public List<com.idcos.enterprise.portal.dal.entity.PortalPermission> queryPermissionsByRoleIdsAndAppId(List<String> authObjIds,
                                                                                                           String appId) {
        PortalPermission portalPermission = Tables.PORTAL_PERMISSION;
        /**
         * 查询所有角色拥有的权限并去重
         */
        PortalPermission distinctAs = portalPermission.as("distinctAs");
        SelectJoinStep<Record1<String>> idSelect = dslContext.select(distinctAs.ID.max()).from(distinctAs);
        Condition condition = distinctAs.APP_ID.eq(appId);
        condition = condition.and(distinctAs.AUTH_OBJ_ID.in(authObjIds));
        idSelect.where(condition).groupBy(distinctAs.AUTH_RES_ID);
        PortalPermission prAs = portalPermission.as("prAs");
        SelectJoinStep<Record> select = dslContext.select(prAs.fields()).from(prAs);
        select.where(prAs.ID.in(idSelect));
        List<com.idcos.enterprise.portal.dal.entity.PortalPermission> portalPermissionList = select.fetch()
                .into(com.idcos.enterprise.portal.dal.entity.PortalPermission.class);
        if (ListUtils.isEmpty(portalPermissionList)) {
            return Lists.newArrayList();
        }
        return portalPermissionList;
    }

    /**
     * 查询所有用户信息（包含deptName）
     *
     * @return List<PortalUserVO>
     */
    public List<PortalUserVO> listUserVO() {
        ArrayList<Field<?>> fields = Lists.newArrayList(PORTAL_USER.fields());
        fields.add(PORTAL_DEPT.DISPLAY_NAME.as("DEPT_NAME"));
        return dslContext.select(fields).from(PORTAL_USER).leftJoin(PORTAL_DEPT)
                .on(PORTAL_USER.DEPT_ID.eq(PORTAL_DEPT.ID)).where(PORTAL_USER.IS_ACTIVE.eq("Y"))
                .fetchInto(PortalUserVO.class);
    }

    /**
     * 查询所有用户信息（包含deptName）
     *
     * @param tenantId
     * @return List<PortalUserVO>
     */
    public List<PortalUserVO> listUserVOByTenantId(String tenantId) {
        ArrayList<Field<?>> fields = Lists.newArrayList(PORTAL_USER.fields());
        fields.add(PORTAL_DEPT.DISPLAY_NAME.as("DEPT_NAME"));
        return dslContext.select(fields).from(PORTAL_USER).leftJoin(PORTAL_DEPT)
                .on(PORTAL_USER.DEPT_ID.eq(PORTAL_DEPT.ID)).where(PORTAL_USER.TENANT_ID.eq(tenantId))
                .and(PORTAL_USER.IS_ACTIVE.eq("Y")).fetchInto(PortalUserVO.class);
    }

    /**
     * 根据用户id用户信息（包含deptName）
     *
     * @param userId
     * @return AuthInfoVO
     */
    public AuthInfoVO findOne(String userId) {
        ArrayList<Field<?>> fields = Lists.newArrayList(PORTAL_USER.fields());
        fields.add(PORTAL_DEPT.DISPLAY_NAME.as("DEPT_NAME"));
        fields.add(PORTAL_TENANT.DISPLAY_NAME.as("TENANT_NAME"));
        AuthInfoVO authInfoVO = dslContext.select(fields).from(PORTAL_USER).leftJoin(PORTAL_DEPT)
                .on(PORTAL_USER.DEPT_ID.eq(PORTAL_DEPT.ID)).leftJoin(PORTAL_TENANT)
                .on(PORTAL_USER.TENANT_ID.eq(PORTAL_TENANT.NAME)).where(PORTAL_USER.ID.eq(userId))
                .fetchOneInto(AuthInfoVO.class);
        List<String> userGroupIds = dslContext.select(PORTAL_USER_GROUP.ID).from(PORTAL_USER_GROUP)
                .leftJoin(PORTAL_GROUP_USER_REL).on(PORTAL_USER_GROUP.ID.eq(PORTAL_GROUP_USER_REL.GROUP_ID))
                .where(PORTAL_USER_GROUP.IS_ACTIVE.eq("Y")).and(PORTAL_GROUP_USER_REL.USER_ID.eq(userId)).fetch()
                .into(String.class);
        authInfoVO.setUserGroups(userGroupIds);
        return authInfoVO;
    }

    /**
     * 根据角色id查询所有用户（包含deptName）
     *
     * @param roleId
     * @return List<String>
     */
    public List<PortalUserVO> listUserByRoleId(String roleId) {
        ArrayList<Field<?>> fields = Lists.newArrayList(PORTAL_USER.fields());
        fields.add(PORTAL_DEPT.DISPLAY_NAME.as("DEPT_NAME"));
        return dslContext.selectDistinct(fields).from(PORTAL_USER).leftJoin(PORTAL_DEPT)
                .on(PORTAL_USER.DEPT_ID.eq(PORTAL_DEPT.ID)).leftJoin(PORTAL_DEPT_ROLE_REL)
                .on(PORTAL_USER.DEPT_ID.eq(PORTAL_DEPT_ROLE_REL.DEPT_ID)).leftJoin(PORTAL_GROUP_USER_REL)
                .on(PORTAL_USER.ID.eq(PORTAL_GROUP_USER_REL.USER_ID)).leftJoin(PORTAL_GROUP_ROLE_REL)
                .on(PORTAL_GROUP_USER_REL.GROUP_ID.eq(PORTAL_GROUP_ROLE_REL.GROUP_ID))
                .where(PORTAL_GROUP_ROLE_REL.ROLE_ID.eq(roleId)).or(PORTAL_DEPT_ROLE_REL.ROLE_ID.eq(roleId)).fetch()
                .into(PortalUserVO.class);
    }

    /**
     * 根据角色id分页查询用户
     *
     * @param form
     * @return List<String>
     */
    public Pagination<PortalUserVO> listUserByRoleForm(PortalQueryByPageForm form) {
        ArrayList<Field<?>> fields = Lists.newArrayList(PORTAL_USER.fields());
        fields.add(PORTAL_DEPT.DISPLAY_NAME.as("DEPT_NAME"));
        SelectConditionStep<Record> select = dslContext.selectDistinct(fields).from(PORTAL_USER).leftJoin(PORTAL_DEPT)
                .on(PORTAL_USER.DEPT_ID.eq(PORTAL_DEPT.ID)).leftJoin(PORTAL_DEPT_ROLE_REL)
                .on(PORTAL_USER.DEPT_ID.eq(PORTAL_DEPT_ROLE_REL.DEPT_ID)).leftJoin(PORTAL_GROUP_USER_REL)
                .on(PORTAL_USER.ID.eq(PORTAL_GROUP_USER_REL.USER_ID)).leftJoin(PORTAL_GROUP_ROLE_REL)
                .on(PORTAL_GROUP_USER_REL.GROUP_ID.eq(PORTAL_GROUP_ROLE_REL.GROUP_ID))
                .where(PORTAL_GROUP_ROLE_REL.ROLE_ID.eq(form.getId())).or(PORTAL_DEPT_ROLE_REL.ROLE_ID.eq(form.getId()));

        int count = dslContext.fetchCount(select);

        if (count == 0) {
            return new Pagination<>(form.getPageNo(), form.getPageSize(), 0, Lists.<PortalUserVO>newArrayList());
        }
        List<PortalUserVO> portalUserVOS = select.limit((form.getPageNo() - 1) * form.getPageSize(), form.getPageSize())
                .fetch().into(PortalUserVO.class);
        return new Pagination<>(form.getPageNo(), form.getPageSize(), count, portalUserVOS);
    }

    /**
     * 根据用户id分页查询角色信息，支持模糊查询
     *
     * @param form
     * @return List<String>
     */
    public Pagination<PortalRoleVO> listRoleByUserForm(PortalQueryByPageForm form) {
        ArrayList<Field<?>> fields = Lists.newArrayList(PORTAL_ROLE.fields());
        SelectConditionStep<Record> select = dslContext.selectDistinct(fields).from(PORTAL_ROLE)
                .leftJoin(PORTAL_GROUP_ROLE_REL).on(PORTAL_ROLE.ID.eq(PORTAL_GROUP_ROLE_REL.ROLE_ID))
                .leftJoin(PORTAL_GROUP_USER_REL).on(PORTAL_GROUP_ROLE_REL.GROUP_ID.eq(PORTAL_GROUP_USER_REL.GROUP_ID))
                .leftJoin(PORTAL_DEPT_ROLE_REL).on(PORTAL_ROLE.ID.eq(PORTAL_DEPT_ROLE_REL.ROLE_ID)).leftJoin(PORTAL_USER)
                .on(PORTAL_USER.DEPT_ID.eq(PORTAL_DEPT_ROLE_REL.DEPT_ID))
                .where(PORTAL_ROLE.IS_ACTIVE.eq(IsActiveEnum.HAS_ACTIVE.getCode())
                        .and(PORTAL_GROUP_USER_REL.USER_ID.eq(form.getId())).or(PORTAL_USER.ID.eq(form.getId())));

        if (StringUtils.isNotBlank(form.getCnd())) {
            Condition condition = PORTAL_ROLE.NAME.like(form.getCnd()).or(PORTAL_ROLE.CODE.like(form.getCnd()))
                    .or(PORTAL_ROLE.REMARK.like(form.getCnd()));
            select.and(condition);
        }
        int count = dslContext.fetchCount(select);

        if (count == 0) {
            return new Pagination<>(form.getPageNo(), form.getPageSize(), 0, Lists.<PortalRoleVO>newArrayList());
        }
        List<PortalRoleVO> portalRoleVOS = select.limit((form.getPageNo() - 1) * form.getPageSize(), form.getPageSize())
                .fetch().into(PortalRoleVO.class);
        return new Pagination<>(form.getPageNo(), form.getPageSize(), count, portalRoleVOS);
    }

    /**
     * 根据appId和角色id的List(无重复值)获取权限List（无重复值）
     *
     * @param appId,roleIds
     * @return List<String>
     */
    public List<PortalPermissionVO> getPermissionsByroleIds(String appId, List<String> roleIds) {
        ArrayList<Field<?>> fields = Lists.newArrayList(PORTAL_PERMISSION.fields());
        List<PortalPermissionVO> portalPermissions = dslContext.select(fields).from(PORTAL_PERMISSION)
                .where(PORTAL_PERMISSION.AUTH_OBJ_ID.in(roleIds).and(PORTAL_PERMISSION.APP_ID.eq(appId))).fetch()
                .into(PortalPermissionVO.class);
        return portalPermissions;
    }

}
