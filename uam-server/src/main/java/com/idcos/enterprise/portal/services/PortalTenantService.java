/**
 * 杭州云霁科技有限公司 http://www.idcos.com Copyright (c) 2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.services;

import java.util.List;

import org.apache.commons.lang3.StringUtils;
import org.jooq.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.google.common.collect.Lists;
import com.idcos.cloud.core.dal.common.page.Pagination;
import com.idcos.enterprise.jooq.Tables;
import com.idcos.enterprise.jooq.tables.PortalTenant;
import com.idcos.enterprise.portal.dal.enums.IsActiveEnum;
import com.idcos.enterprise.portal.form.PortalTenantQueryPageListForm;
import com.idcos.enterprise.portal.vo.PortalTenantVO;

/**
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月07 上午10:11 souakiragen Exp $
 */
@Service
public class PortalTenantService {

    @Autowired
    private DSLContext dslContext;

    public Pagination queryPageList(int pageNo, int pageSize, PortalTenantQueryPageListForm form) {
        PortalTenant portalTenant = Tables.PORTAL_TENANT;
        SelectJoinStep<Record3<String, String, String>> select = dslContext
            .select(portalTenant.ID, portalTenant.NAME.as("TENANT_ID"), portalTenant.DISPLAY_NAME.as("TENANT_NAME"))
            .from(portalTenant);
        Condition condition = portalTenant.IS_ACTIVE.eq(IsActiveEnum.HAS_ACTIVE.getCode());
        if (StringUtils.isNotBlank(form.getTenantId())) {
            condition = condition.and(portalTenant.NAME.like("%" + form.getTenantId() + "%"));
        }
        if (StringUtils.isNotBlank(form.getTenantName())) {
            condition = condition.and(portalTenant.DISPLAY_NAME.like("%" + form.getTenantName() + "%"));
        }
        select.where(condition);
        int total = dslContext.fetchCount(select);
        if (total == 0) {
            return new Pagination(pageNo, pageSize, 0, Lists.newArrayList());
        }
        List<PortalTenantVO> portalTenantVOS = select.orderBy(portalTenant.GMT_CREATE.desc())
            .limit((pageNo - 1) * pageSize, pageSize).fetch().into(PortalTenantVO.class);
        return new Pagination(pageNo, pageSize, total, portalTenantVOS);
    }

    public void updateIsActive(String id, String isActive) {
        PortalTenant portalTenant = Tables.PORTAL_TENANT;
        dslContext.update(portalTenant).set(portalTenant.IS_ACTIVE, isActive).where(portalTenant.ID.eq(id)).execute();
    }
}
