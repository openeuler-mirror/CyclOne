/**
 * 杭州云霁科技有限公司 http://www.idcos.com Copyright (c) 2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.services;

import java.util.List;

import org.jooq.DSLContext;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.thymeleaf.util.ListUtils;

import com.google.common.collect.Lists;
import com.idcos.enterprise.jooq.Tables;
import com.idcos.enterprise.jooq.tables.PortalGroupRoleRel;

/**
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月04 下午5:04 souakiragen Exp $
 */
@Service
public class PortalGroupRoleRelService {

    @Autowired
    private DSLContext dslContext;

    public List<String> getRoleIdsByGroupIds(List<String> ids) {
        PortalGroupRoleRel portalGroupRoleRel = Tables.PORTAL_GROUP_ROLE_REL;
        List<String> roleIds = dslContext.select(portalGroupRoleRel.ROLE_ID).from(portalGroupRoleRel)
            .where(portalGroupRoleRel.GROUP_ID.in(ids)).groupBy(portalGroupRoleRel.ROLE_ID).fetch().into(String.class);
        if (ListUtils.isEmpty(roleIds)) {
            return Lists.newArrayList();
        }
        return roleIds;
    }
}
