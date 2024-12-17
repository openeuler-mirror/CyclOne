/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.services;

import com.idcos.enterprise.jooq.Tables;
import com.idcos.enterprise.jooq.tables.PortalRole;
import org.jooq.DSLContext;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

/**
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月07 下午9:10 souakiragen Exp $
 */
@Service
public class PortalRoleService {

    @Autowired
    private DSLContext dslContext;

    public List<com.idcos.enterprise.portal.dal.entity.PortalRole> getListByIds(List<String> ids) {
        PortalRole portalRole = Tables.PORTAL_ROLE;
        List<com.idcos.enterprise.portal.dal.entity.PortalRole> portalRoles = dslContext.select()
                .from(portalRole).where(portalRole.ID.in(ids)).fetch()
                .into(com.idcos.enterprise.portal.dal.entity.PortalRole.class);
        return portalRoles;
    }
}
