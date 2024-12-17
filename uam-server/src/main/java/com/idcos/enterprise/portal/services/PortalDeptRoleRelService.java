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
import com.idcos.enterprise.jooq.tables.PortalDeptRoleRel;

/**
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月06 上午11:14 souakiragen Exp $
 */
@Service
public class PortalDeptRoleRelService {

    @Autowired
    private DSLContext dslContext;

    public List<String> getRoleIdsByDeptIds(List<String> deptIds) {
        PortalDeptRoleRel portalDeptRoleRel = Tables.PORTAL_DEPT_ROLE_REL;
        List<String> ids = dslContext.select(portalDeptRoleRel.ROLE_ID).from(portalDeptRoleRel)
            .where(portalDeptRoleRel.DEPT_ID.in(deptIds)).groupBy(portalDeptRoleRel.ROLE_ID).fetch().into(String.class);
        if (ListUtils.isEmpty(ids)) {
            ids = Lists.newArrayList();
        }
        return ids;
    }
}
