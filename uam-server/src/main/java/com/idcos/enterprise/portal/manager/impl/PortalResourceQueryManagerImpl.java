

package com.idcos.enterprise.portal.manager.impl;

// auto generated imports

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.idcos.cloud.biz.common.check.CommonParamtersChecker;
import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.cloud.core.common.util.ListUtil;
import com.idcos.enterprise.portal.biz.common.CommonBizException;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryCallback;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryTemplate;
import com.idcos.enterprise.portal.convert.PortalResourceConvert;
import com.idcos.enterprise.portal.dal.entity.PortalResource;
import com.idcos.enterprise.portal.dal.enums.IsActiveEnum;
import com.idcos.enterprise.portal.dal.repository.PortalResourceRepository;
import com.idcos.enterprise.portal.manager.auto.PortalResourceQueryManager;
import com.idcos.enterprise.portal.biz.common.utils.CurrentUser;
import static com.idcos.enterprise.portal.UamConstant.ADMIN;
/**
 * PortalResourceOperateManagerImpl
 *
 * @author pengganyu
 * @version $Id: PortalResourceOperateManagerImpl.java, v 0.1 2016年5月10日 上午9:50:27 pengganyu Exp $
 */
@Service
public class PortalResourceQueryManagerImpl implements PortalResourceQueryManager {
    @Autowired
    private PortalResourceRepository portalResourceRepository;
    @Autowired
    private BusinessQueryTemplate businessQueryTemplate;
    @Autowired
    private PortalResourceConvert portalResourceConvert;
    @Autowired
    private CurrentUser currentUser;  

    @Override
    public CommonResult<?> queryByCode(final String code) {
        if (!ADMIN.equals(currentUser.getUser().getLoginId())) {
            throw new CommonBizException("only admin permits");
        }
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {

            @Override
            public Object doQuery() {
                PortalResource res = portalResourceRepository.findByCodeAndIsActive(code, "Y");

                if (res == null) {
                    throw new CommonBizException("未查询到有效的权限资源信息");
                }

                return portalResourceConvert.apply(res);
            }

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(code);
            }
        });
    }

    @Override
    public CommonResult<?> queryAll() {
        if (!ADMIN.equals(currentUser.getUser().getLoginId())) {
            throw new CommonBizException("only admin permits");
        }
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {

            @Override
            public Object doQuery() {

                return ListUtil.transform(
                        portalResourceRepository.findByIsActive(IsActiveEnum.HAS_ACTIVE.getCode()),
                        portalResourceConvert);
            }

            @Override
            public void checkParam() {
            }
        });
    }
}
