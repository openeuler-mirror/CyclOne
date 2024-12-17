/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.web.auto;

import com.idcos.cloud.biz.common.permisson.OperateAction;
import com.idcos.cloud.web.common.BaseResultVO;
import com.idcos.cloud.web.common.JsonResultUtil;
import com.idcos.enterprise.portal.form.PortalDeptAddForm;
import com.idcos.enterprise.portal.form.PortalDeptUpdateForm;
import com.idcos.enterprise.portal.manager.auto.PortalDeptOperateManager;
import com.idcos.enterprise.portal.manager.auto.PortalDeptQueryManager;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.*;

/**
 * 部门管理
 *
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月03 上午10:33 souakiragen Exp $
 */
@Controller
@OperateAction("F_PORTAL_DEPT")
@RequestMapping(value = "/portal/dept")
@Api(tags = "05.部门管理的接口", description = "PortalDeptController")
public class PortalDeptController {

    @Autowired
    private PortalDeptOperateManager portalDeptOperateManager;

    @Autowired
    private PortalDeptQueryManager portalDeptQueryManager;

    @RequestMapping(method = RequestMethod.POST, value = "")
    @ResponseBody
    @OperateAction("ADD")
    @ApiOperation(value = "新增部门", notes = "新增部门")
    public BaseResultVO add(@ModelAttribute PortalDeptAddForm form) {
        return JsonResultUtil.getResult(portalDeptOperateManager.add(form));
    }

    @RequestMapping(method = RequestMethod.PUT, value = "")
    @ResponseBody
    @OperateAction("update")
    @ApiOperation(value = "修改部门", notes = "修改部门")
    public BaseResultVO update(@ModelAttribute PortalDeptUpdateForm form) {
        return JsonResultUtil.getResult(portalDeptOperateManager.update(form));
    }

    @RequestMapping(method = RequestMethod.DELETE, value = "/{id}")
    @ResponseBody
    @OperateAction("delete")
    @ApiOperation(value = "删除部门", notes = "删除部门")
    public BaseResultVO delete(@PathVariable String id) {
        return JsonResultUtil.getResult(portalDeptOperateManager.delete(id));
    }

    @RequestMapping(method = RequestMethod.POST, value = "/roles/{id}", produces = "application/json;charset=utf-8")
    @ResponseBody
    @OperateAction("assign")
    @ApiOperation(value = "分配角色", notes = "为部门分配角色")
    public BaseResultVO assignRole(@PathVariable String id, @RequestParam String roleIds) {
        return JsonResultUtil.getResult(portalDeptOperateManager.assignRole(id, roleIds));
    }

    @RequestMapping(method = RequestMethod.GET, value = "/roles/{id}")
    @ResponseBody
    @OperateAction("query")
    @ApiOperation(value = "获取部门对应的角色", notes = "获取部门对应的角色")
    public BaseResultVO queryRolesById(@PathVariable String id) {
        return JsonResultUtil.getResult(portalDeptQueryManager.getRolesById(id));
    }

    @RequestMapping(method = RequestMethod.GET, value = "/{id}")
    @ResponseBody
    @OperateAction("query")
    @ApiOperation(value = "获取部门信息", notes = "获取部门信息")
    public BaseResultVO queryDeptById(@PathVariable String id) {
        return JsonResultUtil.getResult(portalDeptQueryManager.getDeptById(id));
    }

}
