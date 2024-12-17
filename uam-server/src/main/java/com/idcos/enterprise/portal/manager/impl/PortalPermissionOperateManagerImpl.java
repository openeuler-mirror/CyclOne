

package com.idcos.enterprise.portal.manager.impl;

// auto generated imports

import com.alibaba.fastjson.JSON;
import com.idcos.cloud.biz.common.check.FormChecker;
import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.cloud.core.common.util.ListUtil;
import com.idcos.cloud.core.common.util.ObjectUtil;
import com.idcos.enterprise.portal.biz.common.CommonBizException;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessCallbackAdator;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessContext;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessTemplate;
import com.idcos.enterprise.portal.biz.common.utils.CurrentUser;
import com.idcos.enterprise.portal.dal.entity.PortalPermission;
import com.idcos.enterprise.portal.dal.entity.PortalResource;
import com.idcos.enterprise.portal.dal.entity.PortalRole;
import com.idcos.enterprise.portal.dal.enums.AuthObjTypeEnum;
import com.idcos.enterprise.portal.dal.enums.IsActiveEnum;
import com.idcos.enterprise.portal.dal.repository.PortalPermissionRepository;
import com.idcos.enterprise.portal.dal.repository.PortalResourceRepository;
import com.idcos.enterprise.portal.dal.repository.PortalRoleRepository;
import com.idcos.enterprise.portal.form.PortalPermissionSaveAuthObjForm;
import com.idcos.enterprise.portal.form.PortalPermissionSaveAuthResForm;
import com.idcos.enterprise.portal.manager.auto.PortalPermissionOperateManager;
import org.apache.commons.lang3.StringUtils;
import org.apache.poi.ss.formula.functions.T;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import static com.idcos.enterprise.portal.UamConstant.ADMIN;
import java.util.List;

/**
 * Manager实现类
 * <p>第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖</p>
 *
 * @author yanlv
 * @version v 1.1 2015-06-09 10:00:37 yanlv Exp $
 */
@Service
public class PortalPermissionOperateManagerImpl implements PortalPermissionOperateManager {
    @Autowired
    private PortalPermissionRepository portalPermissionRepository;

    @Autowired
    private BusinessProcessTemplate businessProcessTemplate;

    @Autowired
    private CurrentUser currentUser;

    @Autowired
    private PortalResourceRepository portalResourceRepository;

    @Autowired
    private PortalRoleRepository portalRoleRepository;

    @Override
    public CommonResult<?> saveAuthObj(final PortalPermissionSaveAuthObjForm form) {
        return businessProcessTemplate.process(new BusinessProcessCallbackAdator<T>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
                //FormChecker.check(form); TODO：jar包待下线
            }

            @Override
            public T doBusiness(BusinessProcessContext context) {

                //取出authObjId信息，组装需要保存和删除的数据信息
                String authObjIds = form.getAuthObjIds();
                if (StringUtils.isNotBlank(authObjIds)) {

                    String[] authObjArr = authObjIds.split(",");

                    for (String str : authObjArr) {
                        String authObjId = str.split(":")[0];
                        String operType = str.split(":")[1];

                        if ("I".equals(operType)) {
                            PortalPermission perPo = new PortalPermission();

                            perPo.setAuthObjId(authObjId);
                            perPo.setAuthObjType(StringUtils.upperCase(form.getAuthObjType()));
                            perPo.setAuthResId(form.getAuthResId());
                            perPo.setAuthResType(StringUtils.upperCase(form.getAuthResType()));

                            //获取资源系统的id并把AppId放到PortalPermissions中的AppId
                            perPo.setAppId(getAppId(form.getAuthResType()));

                            perPo.setAuthResName(StringUtils.isBlank(form.getAuthResName()) ? form.getAuthResId()
                                    : form.getAuthResName());
                            perPo.setTenant(currentUser.getUser().getTenantId());
                            portalPermissionRepository.save(perPo);
                        }

                        if ("D".equals(operType)) {
                            portalPermissionRepository.deleteByAuthResIdAndAuthObjId(form.getAuthResId(), authObjId);
                        }
                    }
                }

                return null;
            }
        });
    }

    /**
     * 获取资源系统的id
     *
     * @param code
     * @return
     */
    private String getAppId(String code) {
        PortalResource portalResource = portalResourceRepository.findByCodeAndIsActive(code,
                IsActiveEnum.HAS_ACTIVE.getCode());
        if (portalResource == null) {
            throw new CommonBizException("未查询到要保存的权限资源信息");
        }
        return portalResource.getAppId();
    }

    @Override
    public CommonResult<?> saveAuthRes(final String permissionInfo) {
        if (!ADMIN.equals(currentUser.getUser().getLoginId())) {
            throw new CommonBizException("非admin用户不允许给其他用户分配权限！！！");
        }
        return businessProcessTemplate.process(new BusinessProcessCallbackAdator<T>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
            }

            @Override
            public T doBusiness(BusinessProcessContext context) {
                List<PortalPermissionSaveAuthResForm> permissionList = JSON.parseArray(permissionInfo,
                        PortalPermissionSaveAuthResForm.class);
                savePermissionList(permissionList, currentUser.getUser().getTenantId());

                return null;
            }
        });
    }

    @Override
    public CommonResult<?> assignByRoleCode(final String tenantId, final String roleCode,
                                            final List<PortalPermissionSaveAuthResForm> permissionList) {
        return businessProcessTemplate.process(new BusinessProcessCallbackAdator<T>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
            }

            @Override
            public T doBusiness(BusinessProcessContext context) {
                List<PortalRole> roles = portalRoleRepository.findByCodeAndIsActive(roleCode, "Y");
                if (roles == null || roles.size() == 0) {
                    throw new RuntimeException("角色Code不存在，无法绑定权限！");
                }
                String roleId = roles.get(0).getId();

                for (PortalPermissionSaveAuthResForm form : permissionList) {
                    form.setAuthObjType(AuthObjTypeEnum.ROLE.getCode());
                    form.setAuthObjId(roleId);
                }

                savePermissionList(permissionList, tenantId);

                return null;
            }
        });
    }

    /**
     * 增加或删除权限与角色的绑定关系
     *
     * @param permissionList
     * @param tenantId
     */
    private void savePermissionList(List<PortalPermissionSaveAuthResForm> permissionList, String tenantId) {
        //这里要防止重复插入
        if (permissionList == null || permissionList.size() == 0) {
            return;
        }
        List<PortalPermission> dbPermissions = portalPermissionRepository.queryAuthRes("ROLE",
                permissionList.get(0).getAuthResType(), permissionList.get(0).getAuthObjId());
        for (PortalPermissionSaveAuthResForm form : permissionList) {
            //FormChecker.check(form); TODO：jar包待下线
            if (dbPermissions != null && dbPermissions.size() > 0) {
                //数据库有数据
                PortalPermission existPermission = ListUtil.findOne(dbPermissions, "authResId", form.getAuthResId());
                if (existPermission != null) {
                    //数据库有列表数据
                    if ("D".equals(form.getOperType())) {
                        portalPermissionRepository.deleteByAuthResIdAndAuthObjId(form.getAuthResId(),
                                form.getAuthObjId());
                    }
                } else {
                    //数据库没有列表数据
                    if ("I".equals(form.getOperType())) {
                        PortalPermission permission = new PortalPermission();
                        ObjectUtil.copyProperties(form, permission);
                        permission.setTenant(tenantId);
                        permission.setAppId(getAppId(form.getAuthResType()));
                        portalPermissionRepository.save(permission);
                    }
                }
            } else if (dbPermissions == null || dbPermissions.size() <= 0) {
                //数据库没有数据
                if ("I".equals(form.getOperType())) {
                    PortalPermission permission = new PortalPermission();
                    ObjectUtil.copyProperties(form, permission);
                    permission.setTenant(tenantId);
                    permission.setAppId(getAppId(form.getAuthResType()));
                    portalPermissionRepository.save(permission);
                }
            }

        }
    }
}
