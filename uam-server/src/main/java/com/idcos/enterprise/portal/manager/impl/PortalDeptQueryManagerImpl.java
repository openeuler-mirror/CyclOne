/**
 * 杭州云霁科技有限公司 http://www.idcos.com Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.manager.impl;

import java.util.ArrayList;
import java.util.List;

import org.apache.commons.lang.StringUtils;
import org.springframework.beans.BeanUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.google.common.collect.Lists;
import com.idcos.cloud.biz.common.check.CommonParamtersChecker;
import com.idcos.cloud.core.common.BaseLinkVO;
import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryCallback;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryTemplate;
import com.idcos.enterprise.portal.biz.common.utils.PortalFilterUtil;
import com.idcos.enterprise.portal.dal.entity.*;
import com.idcos.enterprise.portal.dal.enums.IsActiveEnum;
import com.idcos.enterprise.portal.dal.enums.TreeStyleEnum;
import com.idcos.enterprise.portal.dal.repository.*;
import com.idcos.enterprise.portal.export.util.DeptTreeUtils;
import com.idcos.enterprise.portal.manager.auto.PortalDeptQueryManager;
import com.idcos.enterprise.portal.services.PortalRoleService;
import com.idcos.enterprise.portal.vo.PortalDeptTreeVO;
import com.idcos.enterprise.portal.vo.PortalDeptVO;
import com.idcos.enterprise.portal.biz.common.CommonBizException;
import com.idcos.enterprise.portal.biz.common.utils.CurrentUser;
import static com.idcos.enterprise.portal.UamConstant.ADMIN;

/**
 * @author Dana
 * @version PortalDeptQueryManagerImpl.java, v1 2017/9/26 下午5:33 Dana Exp $$
 */
@Service
public class PortalDeptQueryManagerImpl implements PortalDeptQueryManager {
    @Autowired
    private BusinessQueryTemplate businessQueryTemplate;

    @Autowired
    private PortalDeptRepository portalDeptRepository;

    @Autowired
    private PortalDeptRoleRelRepository portalDeptRoleRelRepository;

    @Autowired
    private PortalRoleService portalRoleService;

    @Autowired
    private CurrentUser currentUser;  

    @Override
    public CommonResult<?> getAllDept(final String tenantId) {
        if (!ADMIN.equals(currentUser.getUser().getLoginId())) {
            throw new CommonBizException("only admin permits");
        }
        return businessQueryTemplate.process(new BusinessQueryCallback<List<PortalDeptVO>>() {

            @Override
            public void checkParam() {}

            @Override
            public List<PortalDeptVO> doQuery() {
                List<PortalDept> portalDepts;
                if (tenantId == null) {
                    portalDepts = portalDeptRepository.findAll();
                } else {
                    portalDepts = portalDeptRepository.findAllDeptByTenantId(tenantId);
                }
                List<PortalDeptVO> deptVOS = new ArrayList<>();
                for (PortalDept portalDept : portalDepts) {
                    PortalDeptVO portalDeptVO = new PortalDeptVO();
                    portalDeptVO.setId(portalDept.getId());
                    portalDeptVO.setTenantId(portalDept.getTenantId());
                    portalDeptVO.setManagerId(portalDept.getManagerId());
                    portalDeptVO.setDisplayName(portalDept.getDisplayName());
                    portalDeptVO.setStatus(portalDept.getStatus());
                    portalDeptVO.setCode(portalDept.getCode());
                    portalDeptVO.setRemark(portalDept.getRemark());
                    portalDeptVO.setParentId(portalDept.getParentId());
                    deptVOS.add(portalDeptVO);
                }

                return deptVOS;
            }
        });
    }

    @Override
    public CommonResult<?> getDeptsTree(final String tenantId, final String treeStyle) {
        if (!ADMIN.equals(currentUser.getUser().getLoginId())) {
            throw new CommonBizException("only admin permits");
        }
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(tenantId);
                CommonParamtersChecker.checkNotBlank(treeStyle);
            }

            @Override
            public Object doQuery() {
                List<PortalDept> portalDepts = portalDeptRepository.findAllDeptByTenantId(tenantId);
                List<PortalDeptTreeVO> deptTreeList = convert2Dto(portalDepts);
                if (TreeStyleEnum.Z_TREE.getCode().equals(treeStyle)) {
                    return deptTreeList;
                } else if (TreeStyleEnum.IO_TREE.getCode().equals(treeStyle)) {
                    Object tree = DeptTreeUtils.getTree(portalDepts, DeptTreeUtils.TreeStyle.IO_TREE);
                    return tree;
                } else if (TreeStyleEnum.ALL.getCode().equals(treeStyle)) {
                    return new ArrayList<>();
                }
                return null;
            }
        });
    }

    @Override
    public CommonResult<?> getDeptByDeptId(final String tenantId, final String deptId) {
        if (!ADMIN.equals(currentUser.getUser().getLoginId())) {
            throw new CommonBizException("only admin permits");
        }
        return businessQueryTemplate.process(new BusinessQueryCallback<PortalDeptVO>() {

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(tenantId);
                CommonParamtersChecker.checkNotBlank(deptId);
            }

            @Override
            public PortalDeptVO doQuery() {
                PortalDept portalDept = portalDeptRepository.findByDeptId(deptId);
                PortalDeptVO portalDeptVO = new PortalDeptVO();
                portalDeptVO.setId(portalDept.getId());
                portalDeptVO.setTenantId(portalDept.getTenantId());
                portalDeptVO.setManagerId(portalDept.getManagerId());
                portalDeptVO.setDisplayName(portalDept.getDisplayName());
                portalDeptVO.setStatus(portalDept.getStatus());
                portalDeptVO.setSourceType(portalDept.getSourceType());
                portalDeptVO.setCode(portalDept.getCode());
                portalDeptVO.setRemark(portalDept.getRemark());
                portalDeptVO.setParentId(portalDept.getParentId());
                return portalDeptVO;
            }
        });
    }

    @Override
    public CommonResult<?> getDeptById(final String id) {
        if (!ADMIN.equals(currentUser.getUser().getLoginId())) {
            throw new CommonBizException("only admin permits");
        }
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {
            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(id);
            }

            @Override
            public Object doQuery() {
                PortalDept portalDept = portalDeptRepository.findOne(id);
                PortalDeptVO portalDeptVO = new PortalDeptVO();
                BeanUtils.copyProperties(portalDept, portalDeptVO);
                return portalDeptVO;
            }
        });
    }

    @Override
    public CommonResult<?> getRolesById(final String id) {
        if (!ADMIN.equals(currentUser.getUser().getLoginId())) {
            throw new CommonBizException("only admin permits");
        }
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {
            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(id);
            }

            @Override
            public Object doQuery() {
                // 声明变量
                BaseLinkVO vo = new BaseLinkVO();

                List<PortalRole> roleList = Lists.newArrayList();

                List<String> roleIds = Lists.newArrayList();

                // 获取角色，获取权限
                for (PortalDeptRoleRel rel : portalDeptRoleRelRepository.findByDeptId(id)) {
                    roleIds.add(rel.getRoleId());
                }
                List<PortalRole> portalRoles = portalRoleService.getListByIds(roleIds);
                for (PortalRole portalRole : portalRoles) {
                    if (IsActiveEnum.HAS_ACTIVE.getCode().equals(portalRole.getIsActive())) {
                        roleList.add(portalRole);
                    }
                }
                vo.putLinkObj("ROLE", PortalFilterUtil.filterRole(roleList));

                return vo;
            }
        });
    }

    /**
     * 转化为PortalDeptTreeVO
     *
     * @param deptEntityList
     * @return
     */
    private List<PortalDeptTreeVO> convert2Dto(List<PortalDept> deptEntityList) {
        List<PortalDeptTreeVO> dtos = new ArrayList<>();
        for (PortalDept dept : deptEntityList) {
            PortalDeptTreeVO deptTreeVO = new PortalDeptTreeVO();
            deptTreeVO.setId(dept.getId());
            deptTreeVO.setCode(dept.getCode());
            if (StringUtils.isBlank(dept.getDisplayName())) {
                deptTreeVO.setName(dept.getCode());
            } else {
                deptTreeVO.setName(dept.getDisplayName());
            }
            deptTreeVO.setPid(dept.getParentId());
            deptTreeVO.setSourceType(dept.getSourceType());
            deptTreeVO.setRemark(dept.getRemark());
            deptTreeVO.setTenantId(dept.getTenantId());
            deptTreeVO.setStatus(dept.getStatus());
            deptTreeVO.setManagerId(dept.getManagerId());
            dtos.add(deptTreeVO);
        }
        return dtos;
    }

}