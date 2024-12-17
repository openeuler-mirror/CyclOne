/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.services;

import com.idcos.enterprise.jooq.Tables;
import com.idcos.enterprise.jooq.tables.PortalDept;
import com.idcos.enterprise.jooq.tables.records.PortalDeptRecord;
import com.idcos.enterprise.portal.dal.enums.StatusEnum;
import org.jooq.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

/**
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月04 下午2:12 souakiragen Exp $
 */
@Service
public class PortalDeptService {

    @Autowired
    private DSLContext dslContext;

    public boolean checkHasByTenantIdAndParentIdAndDisplayName(String tenantId, String parentId, String displayName) {
        PortalDept portalDept = Tables.PORTAL_DEPT;
        SelectJoinStep<Record> select = dslContext.select().from(portalDept);
        Condition condition = portalDept.TENANT_ID.eq(tenantId);
        condition = condition.and(portalDept.PARENT_ID.eq(parentId));
        condition = condition.and(portalDept.STATUS.eq(StatusEnum.Y.getCode()));
        condition = condition.and(portalDept.DISPLAY_NAME.eq(displayName));
        select.where(condition);
        int count = dslContext.fetchCount(select);
        return count > 0;
    }

    public boolean checkHasByTenantIdAndParentIdAndDisplayNameAndNotId(String tenantId, String parentId,
                                                                       String displayName, String id) {
        PortalDept portalDept = Tables.PORTAL_DEPT;
        SelectJoinStep<Record> select = dslContext.select().from(portalDept);
        Condition condition = portalDept.TENANT_ID.eq(tenantId);
        condition = condition.and(portalDept.PARENT_ID.eq(parentId));
        condition = condition.and(portalDept.DISPLAY_NAME.eq(displayName));
        condition = condition.and(portalDept.STATUS.eq(StatusEnum.Y.getCode()));
        condition = condition.and(portalDept.ID.ne(id));
        select.where(condition);
        int count = dslContext.fetchCount(select);
        return count > 0;
    }

    public void updatePortalDeptStatus(String status, String id) {
        PortalDept portalDept = Tables.PORTAL_DEPT;
        UpdateSetMoreStep<PortalDeptRecord> set = dslContext.update(portalDept).set(portalDept.STATUS, status);
        Condition condition = portalDept.ID.eq(id);
        set.where(condition);
        set.execute();
    }

}
