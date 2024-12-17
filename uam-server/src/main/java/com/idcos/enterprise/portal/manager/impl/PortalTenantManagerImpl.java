

package com.idcos.enterprise.portal.manager.impl;

// auto generated imports

import com.idcos.cloud.biz.common.check.CommonParamtersChecker;
import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.cloud.core.common.util.DateUtil;
import com.idcos.cloud.core.common.util.ListUtil;
import com.idcos.enterprise.portal.biz.common.tempalte.*;
import com.idcos.enterprise.portal.convert.PortalTenantConvert;
import com.idcos.enterprise.portal.dal.entity.PortalTenant;
import com.idcos.enterprise.portal.dal.repository.PortalTenantRepository;
import com.idcos.enterprise.portal.form.PortalTenantQueryPageListForm;
import com.idcos.enterprise.portal.manager.auto.PortalTenantManager;
import com.idcos.enterprise.portal.services.PortalTenantService;
import com.idcos.enterprise.portal.vo.PortalTenantVO;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

/**
 * Manager实现类
 * <p>第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖</p>
 *
 * @author yanlv
 * @version v 1.1 2015-06-09 09:26:24 yanlv Exp $
 */
@Service
public class PortalTenantManagerImpl implements PortalTenantManager {
    @Autowired
    private BusinessQueryTemplate businessQueryTemplate;

    @Autowired
    private PortalTenantRepository portalTenantRepository;

    @Autowired
    private PortalTenantConvert portalTenantConvert;

    @Autowired
    private PortalTenantService portalTenantService;

    @Override
    public CommonResult<?> getAllTenant() {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {

            @Override
            public Object doQuery() {

                List<PortalTenant> portalTenantList = portalTenantRepository.getAllPortalTenant();
                return ListUtil.transform(portalTenantList, portalTenantConvert);
            }

            @Override
            public void checkParam() {
            }
        });
    }

    @Override
    public CommonResult<?> getTenantByTenantId(final String tenantId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {

            @Override
            public Object doQuery() {

                PortalTenant portalTenant = portalTenantRepository.findByName(tenantId);
                return portalTenantConvert.apply(portalTenant);
            }

            @Override
            public void checkParam() {
            }
        });
    }

    @Override
    public CommonResult<?> getTenantInfo(final String tenantId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {
            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(tenantId);
            }

            @Override
            public Object doQuery() {
                PortalTenant portalTenant = portalTenantRepository.findOne(tenantId);
                if (portalTenant == null) {
                    throw new RuntimeException("未找到改租户");
                }
                return po2Vo(portalTenant);
            }
        });
    }

    @Override
    public CommonResult<?> getPageList(final int pageNo, final int pageSize, final PortalTenantQueryPageListForm form) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {
            @Override
            public void checkParam() {

            }

            @Override
            public Object doQuery() {
                int pageNoCheck = pageNo;
                int pageSizeCheck = pageSize;
                if (pageNoCheck == 0) {
                    pageNoCheck = 10;
                }
                if (pageSizeCheck == 0) {
                    pageSizeCheck = 10;
                }
                return portalTenantService.queryPageList(pageNoCheck, pageSizeCheck, form);
            }
        });
    }

    private PortalTenantVO po2Vo(PortalTenant portalTenant) {
        PortalTenantVO portalTenantVO = new PortalTenantVO();
        portalTenantVO.setTenantId(portalTenant.getName());
        portalTenantVO.setTenantName(portalTenant.getDisplayName());
        portalTenantVO.setId(portalTenant.getId());
        portalTenantVO.setGmtCreate(portalTenant.getGmtCreate());
        portalTenantVO.setGmtModified(portalTenant.getGmtModified());
        return portalTenantVO;
    }
}