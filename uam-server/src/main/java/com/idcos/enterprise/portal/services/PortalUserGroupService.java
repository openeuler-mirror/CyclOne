/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.services;

import com.google.common.collect.Lists;
import com.idcos.cloud.core.dal.common.page.Pagination;
import com.idcos.enterprise.jooq.Tables;
import com.idcos.enterprise.jooq.tables.*;
import com.idcos.enterprise.jooq.tables.records.PortalUserGroupRecord;
import com.idcos.enterprise.portal.dal.enums.IsActiveEnum;
import com.idcos.enterprise.portal.form.PortalUserGroupQueryByPageForm;
import com.idcos.enterprise.portal.vo.PortalUserGroupVO;
import com.idcos.enterprise.portal.vo.PortalUserVO;
import org.apache.commons.lang3.StringUtils;
import org.jooq.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;

import static com.idcos.enterprise.jooq.Tables.PORTAL_GROUP_USER_REL;
import static com.idcos.enterprise.jooq.Tables.PORTAL_USER_GROUP;

/**
 * @author Dana
 * @version PortalUserGroupService.java, v1 2017/11/27 上午11:23 Dana Exp $$
 */
@Service
public class PortalUserGroupService {
    @Autowired
    private DSLContext dslContext;

    /**
     * 修改用户组的IsActive属性
     *
     * @param isAction
     * @param id
     */
    public void updateIsActiveById(String isAction, String id) {
        PortalUserGroup userGroup = PORTAL_USER_GROUP;
        UpdateSetMoreStep<PortalUserGroupRecord> set = dslContext.update(userGroup).set(userGroup.IS_ACTIVE, isAction);
        Condition condition = userGroup.ID.eq(id);
        set.where(condition);
        set.execute();
    }

    /**
     * 根据用户id和用户组type查询用户组集合
     *
     * @param type,id
     * @return List<PortalUserGroupVO>
     */
    public List<PortalUserGroupVO> queryGroupByTypeAndUserId(String type, String userId) {
        ArrayList<Field<?>> fields = Lists.newArrayList(PORTAL_USER_GROUP.fields());
        return dslContext.select(fields).from(PORTAL_USER_GROUP).leftJoin(PORTAL_GROUP_USER_REL)
                .on(PORTAL_USER_GROUP.ID.eq(PORTAL_GROUP_USER_REL.GROUP_ID)).where(PORTAL_GROUP_USER_REL.USER_ID.eq(userId))
                .and(PORTAL_USER_GROUP.TYPE.eq(type)).fetch().into(PortalUserGroupVO.class);
    }

    /**
     * 根据用户id获取工作组的list
     *
     * @param userId
     * @return List<PortalUserGroupVO>
     */
    public List<PortalUserGroupVO> listUserGroupByUserId(String userId) {
        ArrayList<Field<?>> fields = Lists.newArrayList(PORTAL_USER_GROUP.fields());
        return dslContext.select(fields).from(PORTAL_USER_GROUP).leftJoin(PORTAL_GROUP_USER_REL)
                .on(PORTAL_USER_GROUP.ID.eq(PORTAL_GROUP_USER_REL.GROUP_ID)).where(PORTAL_GROUP_USER_REL.USER_ID.eq(userId))
                .fetch().into(PortalUserGroupVO.class);
    }

    /**
     * 根据用户组id获取用户的list
     *
     * @param groupId
     * @return List<PortalUserVO>
     */
    public List<PortalUserVO> listUsersByGroupId(String groupId) {
        PortalUser user = Tables.PORTAL_USER;
        PortalGroupUserRel guRel = Tables.PORTAL_GROUP_USER_REL;
        ArrayList<Field<?>> fields = Lists.newArrayList(user.fields());
        return dslContext.select(fields).from(user).leftJoin(guRel).on(user.ID.eq(guRel.USER_ID))
                .where(guRel.GROUP_ID.eq(groupId)).and(user.IS_ACTIVE.eq("Y")).fetch().into(PortalUserVO.class);
    }

    /**
     * 根据用户组名称和租户Id获取用户的list
     *
     * @param groupName
     * @param tenantId
     * @return List<PortalUserVO>
     */
    public List<PortalUserVO> listUsersByGroupNameAndTenantId(String groupName, String tenantId) {
        PortalUser user = Tables.PORTAL_USER;
        PortalUserGroup userGroup = Tables.PORTAL_USER_GROUP;
        PortalGroupUserRel guRel = Tables.PORTAL_GROUP_USER_REL;
        ArrayList<Field<?>> fields = Lists.newArrayList(user.fields());
        return dslContext.select(fields).from(user).leftJoin(guRel).on(user.ID.eq(guRel.USER_ID)).leftJoin(userGroup)
                .on(guRel.GROUP_ID.eq(userGroup.ID)).where(userGroup.NAME.eq(groupName)).and(userGroup.TENANT.eq(tenantId))
                .and(userGroup.IS_ACTIVE.eq("Y")).and(user.IS_ACTIVE.eq("Y")).fetch().into(PortalUserVO.class);
    }

    /**
     * 分页查询用户组（带每个组拥有的角色列表）
     *
     * @param pageNo
     * @param pageSize
     * @param form
     * @return
     */
    public Pagination<PortalUserGroupVO> getPageList(Integer pageNo, Integer pageSize,
                                                     PortalUserGroupQueryByPageForm form, String cnd) {
        PortalUserGroup userGroup = Tables.PORTAL_USER_GROUP;
        PortalGroupRoleRel grRel = Tables.PORTAL_GROUP_ROLE_REL;
        PortalRole role = Tables.PORTAL_ROLE;
        List<Field<?>> fields = Lists.newArrayList(userGroup.fields());
        SelectJoinStep<Record> select = dslContext.select(fields).from(userGroup);
        Condition condition = userGroup.TENANT.eq(form.getTenantId())
                .and(userGroup.IS_ACTIVE.eq(IsActiveEnum.HAS_ACTIVE.getCode()));
        if (StringUtils.isNotBlank(cnd)) {
            condition = condition.and(userGroup.NAME.like("%" + cnd + "%").or(userGroup.REMARK.like("%" + cnd + "%")));
        }
        //按最后更新时间降序展示
        select.where(condition).orderBy(userGroup.GMT_MODIFIED.desc());
        int count = dslContext.fetchCount(select);
        if (count == 0) {
            return new Pagination<>(pageNo, pageSize, 0, Lists.<PortalUserGroupVO>newArrayList());
        }
        List<PortalUserGroupVO> portalUserGroupVOS = select.limit((pageNo - 1) * pageSize, pageSize).fetch()
                .into(PortalUserGroupVO.class);

        //查询每个用户所属的用户组集合，放入PortalUserVO中
        for (PortalUserGroupVO portalUserGroupVO : portalUserGroupVOS) {
            List<String> selRoleList = dslContext.selectDistinct(role.NAME).from(role).leftJoin(grRel)
                    .on(role.ID.eq(grRel.ROLE_ID)).where(role.IS_ACTIVE.eq(IsActiveEnum.HAS_ACTIVE.getCode()))
                    .and(grRel.GROUP_ID.eq(portalUserGroupVO.getId())).fetch().into(String.class);
            String[] selRoles = new String[selRoleList.size()];
            selRoleList.toArray(selRoles);
            portalUserGroupVO.setSelRoles(selRoles);
        }
        return new Pagination<>(pageNo, pageSize, count, portalUserGroupVOS);
    }

}