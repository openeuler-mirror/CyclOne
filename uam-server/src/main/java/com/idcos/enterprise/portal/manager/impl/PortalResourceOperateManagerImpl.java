

package com.idcos.enterprise.portal.manager.impl;

// auto generated imports

import com.idcos.cloud.core.common.util.ListUtil;
import com.idcos.enterprise.portal.biz.common.utils.CurrentUser;
import com.idcos.enterprise.portal.dal.entity.PortalPermission;
import com.idcos.enterprise.portal.dal.repository.PortalPermissionRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.idcos.cloud.biz.common.check.CommonParamtersChecker;
import com.idcos.cloud.biz.common.check.FormChecker;
import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.cloud.core.common.util.ObjectUtil;
import com.idcos.enterprise.portal.biz.common.CommonBizException;
import com.idcos.enterprise.portal.biz.common.ResultCode;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessCallbackAdator;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessContext;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessTemplate;
import com.idcos.enterprise.portal.biz.common.utils.CrudUtilService;
import com.idcos.enterprise.portal.convert.PortalResourceConvert;
import com.idcos.enterprise.portal.dal.entity.PortalResource;
import com.idcos.enterprise.portal.dal.enums.IsActiveEnum;
import com.idcos.enterprise.portal.dal.repository.PortalResourceRepository;
import com.idcos.enterprise.portal.form.PortalResourceCreateForm;
import com.idcos.enterprise.portal.form.PortalResourceUpdateForm;
import com.idcos.enterprise.portal.manager.auto.PortalResourceOperateManager;

import java.util.ArrayList;
import java.util.List;

/**
 * PortalResourceOperateManagerImpl
 *
 * @author pengganyu
 * @version $Id: PortalResourceOperateManagerImpl.java, v 0.1 2016年5月10日 上午9:50:27 pengganyu Exp $
 */
@Service
public class PortalResourceOperateManagerImpl implements PortalResourceOperateManager {
    @Autowired
    private PortalResourceRepository portalResourceRepository;
    @Autowired
    private PortalPermissionRepository portalPermissionRepository;
    @Autowired
    private BusinessProcessTemplate businessProcessTemplate;
    @Autowired
    private CrudUtilService crudUtilService;
    @Autowired
    private PortalResourceConvert portalResourceConvert;
    @Autowired
    private CurrentUser currentUser;

    @Override
    public CommonResult<?> create(final PortalResourceCreateForm form) {
        return businessProcessTemplate.process(new BusinessProcessCallbackAdator<Object>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
                //FormChecker.check(form); TODO：jar包待下线

                if (portalResourceRepository.findByCodeAndIsActive(form.getCode(),
                        IsActiveEnum.HAS_ACTIVE.getCode()) != null) {
                    throw new CommonBizException("不能保存已经存在的权限资源类型" + form.getCode());
                }
            }

            @Override
            public Object doBusiness(BusinessProcessContext context) {
                PortalResource res = new PortalResource();

                ObjectUtil.copyProperties(form, res);
                res.setIsActive("Y");
                res.setTenant(currentUser.getUser().getTenantId());
                portalResourceRepository.save(res);

                return portalResourceConvert.apply(res);
            }
        });
    }

    @Override
    public CommonResult<?> update(final String id, final PortalResourceUpdateForm form) {
        return businessProcessTemplate.process(new BusinessProcessCallbackAdator<Object>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
                //FormChecker.check(form); TODO：jar包待下线
                CommonParamtersChecker.checkNotBlank(id);

            }

            @Override
            public Object doBusiness(BusinessProcessContext context) {

                PortalResource resource = portalResourceRepository.findOne(id);

                if (resource == null) {
                    throw new CommonBizException(ResultCode.QUERY_RESULT_IS_NULL, "未查询到可更新的权限资源信息");
                }

                changePermissionAuthResId(resource, form);


                ObjectUtil.copyProperties(form, resource);

                resource.setTenant(currentUser.getUser().getTenantId());
                portalResourceRepository.save(resource);

                return portalResourceConvert.apply(resource);
            }
        });
    }

    /**
     * 需要修改permission表的权限资源编码，防止脏数据产生
     *
     * @param resource
     * @param form
     */
    private void changePermissionAuthResId(PortalResource resource, PortalResourceUpdateForm form) {
        List<PortalPermission> portalPermissions = new ArrayList<>();
        if (!form.getCode().equals(resource.getCode())) {
            //如果form中的code和数据库不相同，说明权限资源编码发生变化，需要修改permission表的权限资源编码，防止脏数据产生
            portalPermissions = portalPermissionRepository.queryAuthObjByAppIdAndAuthResType(resource.getAppId(), resource.getCode());
        }
        if (portalPermissions.size() > 0) {
            ListUtil.update(portalPermissions, "authResType", resource.getCode(), form.getCode());
            portalPermissionRepository.save(portalPermissions);
        }
    }

    @Override
    public CommonResult<?> delete(final String id) {
        return crudUtilService.delete(id, portalResourceConvert);
    }

}
