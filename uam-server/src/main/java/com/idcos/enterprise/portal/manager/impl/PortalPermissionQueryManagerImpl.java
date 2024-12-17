

package com.idcos.enterprise.portal.manager.impl;

// auto generated imports

import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.idcos.cloud.biz.common.check.FormChecker;
import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.cloud.core.common.util.ListUtil;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryCallback;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryTemplate;
import com.idcos.enterprise.portal.convert.PortalPermissionConvert;
import com.idcos.enterprise.portal.dal.repository.PortalPermissionRepository;
import com.idcos.enterprise.portal.form.PortalPermissionQueryAuthObjForm;
import com.idcos.enterprise.portal.form.PortalPermissionQueryAuthResForm;
import com.idcos.enterprise.portal.manager.auto.PortalPermissionQueryManager;
import com.idcos.enterprise.portal.vo.PortalPermissionVO;
import com.idcos.enterprise.portal.biz.common.CommonBizException;
import com.idcos.enterprise.portal.biz.common.utils.CurrentUser;
import static com.idcos.enterprise.portal.UamConstant.ADMIN;

/**
 * Manager实现类
 * <p>第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖</p>
 *
 * @author yanlv
 * @version v 1.1 2015-06-09 10:00:37 yanlv Exp $
 */
@Service
public class PortalPermissionQueryManagerImpl implements PortalPermissionQueryManager {
    @Autowired
    private BusinessQueryTemplate businessQueryTemplate;
    @Autowired
    private PortalPermissionConvert portalPermissionConvert;
    @Autowired
    private PortalPermissionRepository portalPermissionRepository;

    @Autowired
    private CurrentUser currentUser; 

    @Override
    public CommonResult<?> queryAuthObj(final PortalPermissionQueryAuthObjForm form) {
        if (!ADMIN.equals(currentUser.getUser().getLoginId())) {
            throw new CommonBizException("only admin permits");
        }
        return businessQueryTemplate.process(new BusinessQueryCallback<List<PortalPermissionVO>>() {

            @Override
            public void checkParam() {
                //FormChecker.check(form); TODO：jar包待下线
            }

            @Override
            public List<PortalPermissionVO> doQuery() {
                return ListUtil.transform(
                        portalPermissionRepository.queryAuthObj(form.getAuthObjType(),
                                form.getAuthResType(), form.getAuthResId()), portalPermissionConvert);
            }
        });
    }

    @Override
    public CommonResult<?> queryAuthRes(final PortalPermissionQueryAuthResForm form) {
        if (!ADMIN.equals(currentUser.getUser().getLoginId())) {
            throw new CommonBizException("only admin permits");
        }
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {
            @Override
            public void checkParam() {
                //FormChecker.check(form); TODO：jar包待下线
            }

            @Override
            public Object doQuery() {
                return ListUtil.transform(
                        portalPermissionRepository.queryAuthRes(form.getAuthObjType(),
                                form.getAuthResType(), form.getAuthObjId()), portalPermissionConvert);
            }
        });
    }

}