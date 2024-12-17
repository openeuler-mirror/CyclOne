

package com.idcos.enterprise.portal.manager.impl;

// auto generated imports

import java.util.List;

import com.idcos.enterprise.portal.biz.common.utils.CurrentUser;
import com.idcos.enterprise.portal.dal.entity.PortalPermission;
import com.idcos.enterprise.portal.dal.repository.PortalPermissionRepository;
import org.apache.commons.collections.CollectionUtils;
import org.apache.commons.lang3.StringUtils;
import org.apache.poi.ss.formula.functions.T;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.idcos.cloud.biz.common.check.CommonParamtersChecker;
import com.idcos.cloud.biz.common.check.FormChecker;
import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.cloud.core.common.util.ObjectUtil;
import com.idcos.enterprise.portal.biz.common.CommonBizException;
import com.idcos.enterprise.portal.biz.common.ResultCode;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessCallback;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessContext;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessTemplate;
import com.idcos.enterprise.portal.biz.common.utils.CrudUtilService;
import com.idcos.enterprise.portal.convert.PortalRoleConvert;
import com.idcos.enterprise.portal.dal.entity.PortalGroupRoleRel;
import com.idcos.enterprise.portal.dal.entity.PortalRole;
import com.idcos.enterprise.portal.dal.enums.IsActiveEnum;
import com.idcos.enterprise.portal.dal.repository.PortalGroupRoleRelRepository;
import com.idcos.enterprise.portal.dal.repository.PortalRoleRepository;
import com.idcos.enterprise.portal.form.PortalRoleCreateForm;
import com.idcos.enterprise.portal.form.PortalRoleUpdateForm;
import com.idcos.enterprise.portal.manager.auto.PortalRoleOperateManager;
import com.idcos.enterprise.portal.vo.PortalRoleVO;

/**
 * Manager实现类
 * <p>第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖</p>
 *
 * @author yanlv
 * @version v 1.1 2015-06-06 10:24:34 yanlv Exp $
 */
@Service
public class PortalRoleOperateManagerImpl implements PortalRoleOperateManager {
    @Autowired
    private PortalRoleConvert portalRoleConvert;
    @Autowired
    private CrudUtilService crudUtilService;
    @Autowired
    private PortalGroupRoleRelRepository portalGroupRoleRelRepository;
    @Autowired
    private PortalRoleRepository portalRoleRepository;
    @Autowired
    private BusinessProcessTemplate businessProcessTemplate;

    @Autowired
    private PortalPermissionRepository portalPermissionRepository;

    @Autowired
    private CurrentUser currentUser;

    @Override
    public CommonResult<?> create(final PortalRoleCreateForm form) {
        CommonResult<PortalRoleVO> result = businessProcessTemplate
                .process(new BusinessProcessCallback<PortalRoleVO>() {

                    @Override
                    public void checkParam(BusinessProcessContext context) {
                        //FormChecker.check(form); TODO：jar包待下线
                    }

                    @Override
                    public void checkBusinessInfo(BusinessProcessContext context) {
                        //编码重复校验
                        List<PortalRole> roleList = portalRoleRepository.findByCodeAndIsActive(
                                form.getCode(), IsActiveEnum.HAS_ACTIVE.getCode());
                        if (CollectionUtils.isNotEmpty(roleList)) {
                            throw new CommonBizException(ResultCode.UNKNOWN_EXCEPTION, "编码"
                                    + form.getCode()
                                    + "已经存在，无法保存");
                        }
                    }

                    @Override
                    public PortalRoleVO doBusiness(BusinessProcessContext context) {
                        PortalRole rolePo = new PortalRole();
                        ObjectUtil.copyProperties(form, rolePo);
                        rolePo.setIsActive(IsActiveEnum.HAS_ACTIVE.getCode());
                        rolePo.setTenant(currentUser.getUser().getTenantId());

                        return portalRoleConvert.apply(portalRoleRepository.save(rolePo));
                    }

                    @Override
                    public void exceptionProcess(CommonBizException exception,
                                                 BusinessProcessContext context) {

                    }

                });
        return result;
    }

    @Override
    public CommonResult<?> delete(final String id) {

        //删除和角色关联的用户组关系信息
        deleteRoleGroupRel(id);
        //删除和角色关联的权限关系信息
        deleteRolePermissionsRel(id);
        return crudUtilService.delete(id, portalRoleConvert);
    }

    private void deleteRoleGroupRel(String id) {
        List<PortalGroupRoleRel> portalGroupRoleRelList = portalGroupRoleRelRepository.findByRoleId(id);
        if (!portalGroupRoleRelList.isEmpty() && portalGroupRoleRelList != null) {
            portalGroupRoleRelRepository.delete(portalGroupRoleRelList);
        }
    }


    private void deleteRolePermissionsRel(String id) {
        List<PortalPermission> portalPermissions = portalPermissionRepository.findByAuthObjId(id);
        if (!portalPermissions.isEmpty() && portalPermissions != null) {
            portalPermissionRepository.delete(portalPermissions);
        }
    }


    @Override
    public CommonResult<?> update(final String id, final PortalRoleUpdateForm form) {
        CommonResult<PortalRoleVO> result = businessProcessTemplate
                .process(new BusinessProcessCallback<PortalRoleVO>() {

                    @Override
                    public void checkParam(BusinessProcessContext context) {
                        //FormChecker.check(form); TODO：jar包待下线
                        CommonParamtersChecker.checkNotBlank(id);
                    }

                    @Override
                    public void checkBusinessInfo(BusinessProcessContext context) {
                        PortalRole rolePo = portalRoleRepository.findOne(id);
                        if (rolePo == null
                                || IsActiveEnum.NO_ACTIVE.getCode().equals(rolePo.getIsActive())) {
                            throw new CommonBizException(ResultCode.QUERY_RESULT_IS_NULL, "未查询到有效的角色信息");
                        }

                        //编码重复校验
                        List<PortalRole> roleList = portalRoleRepository.findByCodeAndIsActive(
                                form.getCode(), IsActiveEnum.HAS_ACTIVE.getCode());

                        for (PortalRole role : roleList) {
                            if (!role.getId().equals(id)) {
                                throw new CommonBizException(ResultCode.UNKNOWN_EXCEPTION,
                                        "编码" + form.getCode() + "已经存在，无法保存");
                            }
                        }

                        context.put("PortalRole", rolePo);

                    }

                    @Override
                    public PortalRoleVO doBusiness(BusinessProcessContext context) {
                        PortalRole rolePo = (PortalRole) context.get("PortalRole");
                        ObjectUtil.copyProperties(form, rolePo);

                        return portalRoleConvert.apply(portalRoleRepository.save(rolePo));
                    }

                    @Override
                    public void exceptionProcess(CommonBizException exception,
                                                 BusinessProcessContext context) {

                    }

                });
        return result;
    }

    @Override
    public CommonResult<?> allocateGroup(final String id, final String selGroups) {
        CommonResult<T> result = businessProcessTemplate.process(new BusinessProcessCallback<T>() {

            @Override
            public void checkParam(BusinessProcessContext context) {
                //FormChecker.check(id);
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {

            }

            @Override
            public T doBusiness(BusinessProcessContext context) {

                if (StringUtils.isNoneBlank(selGroups)) {
                    String[] groupList = selGroups.split(",");
                    for (String str : groupList) {
                        //根据 :  来截取id和操作类型
                        String groupId = str.split(":")[0];
                        String operType = str.split(":")[1];

                        //D代表删除
                        if ("D".equals(operType)) {
                            portalGroupRoleRelRepository.deleteByGroupIdAndRoleId(groupId, id);
                        }

                        //I代表新增
                        if ("I".equals(operType)) {
                            PortalGroupRoleRel relPo = new PortalGroupRoleRel();
                            relPo.setGroupId(groupId);
                            relPo.setRoleId(id);
                            relPo.setTenant(currentUser.getUser().getTenantId());
                            portalGroupRoleRelRepository.save(relPo);
                        }
                    }
                }

                return null;
            }

            @Override
            public void exceptionProcess(CommonBizException exception,
                                         BusinessProcessContext context) {
            }
        });
        return result;
    }
}
