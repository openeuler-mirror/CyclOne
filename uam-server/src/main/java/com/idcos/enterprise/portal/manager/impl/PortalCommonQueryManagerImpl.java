

package com.idcos.enterprise.portal.manager.impl;

// auto generated imports

import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.cloud.core.common.util.ListUtil;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryCallback;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryTemplate;
import com.idcos.enterprise.portal.biz.common.utils.CurrentUser;
import com.idcos.enterprise.portal.biz.common.CommonBizException;
import com.idcos.enterprise.portal.convert.PortalRoleConvert;
import com.idcos.enterprise.portal.convert.PortalUserGroupConvert;
import com.idcos.enterprise.portal.dal.repository.PortalRoleRepository;
import com.idcos.enterprise.portal.dal.repository.PortalUserGroupRepository;
import com.idcos.enterprise.portal.manager.auto.PortalCommonQueryManager;
import com.idcos.enterprise.portal.vo.PortalRoleVO;
import com.idcos.enterprise.portal.vo.PortalUserGroupVO;

import static com.idcos.enterprise.portal.UamConstant.ADMIN;

/**
 * Manager实现类
 * <p>第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖</p>
 *
 * @author yanlv
 * @version v 1.1 2015-06-09 10:37:31 yanlv Exp $
 */
@Service
public class PortalCommonQueryManagerImpl implements PortalCommonQueryManager {

    @Autowired
    private BusinessQueryTemplate businessQueryTemplate;
    @Autowired
    private PortalUserGroupRepository portalUserGroupRepository;
    @Autowired
    private PortalUserGroupConvert portalUserGroupConvert;
    @Autowired
    private PortalRoleRepository portalRoleRepository;
    @Autowired
    private PortalRoleConvert portalRoleConvert;
    @Autowired
    private CurrentUser currentUser;

    @Override
    public CommonResult<?> getRoleList() {
        if (!ADMIN.equals(currentUser.getUser().getLoginId())) {
            throw new CommonBizException("only admin permits");
        }      
        return businessQueryTemplate.process(new BusinessQueryCallback<List<PortalRoleVO>>() {
            @Override
            public void checkParam() {
            }

            @Override
            public List<PortalRoleVO> doQuery() {
                return ListUtil.transform(portalRoleRepository.queryRoleList(), portalRoleConvert);
            }
        });
    }

    @Override
    public CommonResult<?> getUserGroupList(final String tenantId) {
        if (!ADMIN.equals(currentUser.getUser().getLoginId())) {
            throw new CommonBizException("only admin permits");
        }
        return businessQueryTemplate.process(new BusinessQueryCallback<List<PortalUserGroupVO>>() {
            @Override
            public void checkParam() {
            }

            @Override
            public List<PortalUserGroupVO> doQuery() {
                return ListUtil.transform(portalUserGroupRepository.queryUserGroupList(tenantId),
                        portalUserGroupConvert);
            }
        });
    }

}