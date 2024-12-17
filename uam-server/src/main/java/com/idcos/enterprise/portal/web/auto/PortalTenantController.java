/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.web.auto;

import com.idcos.cloud.biz.common.permisson.OperateAction;
import com.idcos.cloud.web.common.BaseResultVO;
import com.idcos.cloud.web.common.JsonResultUtil;
import com.idcos.enterprise.portal.form.PortalTenantAddForm;
import com.idcos.enterprise.portal.form.PortalTenantQueryPageListForm;
import com.idcos.enterprise.portal.form.PortalTenantUpdateForm;
import com.idcos.enterprise.portal.manager.auto.PortalTenantManager;
import com.idcos.enterprise.portal.manager.auto.PortalTenantOperateManager;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.*;

/**
 * 租户管理
 *
 * @author souakiragen
 * @version $Id: , v 0.1 2017年11月07 上午9:09 souakiragen Exp $
 */
@Controller
@OperateAction("F_PORTAL_TENANT")
@RequestMapping(value = "/portal/tenant")
@Api(tags = "06.租户管理的接口", description = "PortalTenantController")
public class PortalTenantController {

    @Autowired
    private PortalTenantOperateManager portalTenantOperateManager;

    @Autowired
    private PortalTenantManager portalTenantManager;

    @RequestMapping(method = RequestMethod.POST, value = "")
    @ResponseBody
    @OperateAction("ADD")
    @ApiOperation(value = "新增租户", notes = "新增租户")
    public BaseResultVO add(@ModelAttribute PortalTenantAddForm form) {
        return JsonResultUtil.getResult(portalTenantOperateManager.add(form));
    }

    @RequestMapping(method = RequestMethod.GET, value = "/{id}")
    @ResponseBody
    @OperateAction("QUERY")
    @ApiOperation(value = "获取租户信息", notes = "获取租户信息")
    public BaseResultVO query(@PathVariable String id) {
        return JsonResultUtil.getResult(portalTenantManager.getTenantInfo(id));
    }

    @RequestMapping(method = RequestMethod.GET, value = "/pageList/{pageNo}/{pageSize}")
    @ResponseBody
    @OperateAction("QUERY")
    @ApiOperation(value = "获取分页数据", notes = "获取分页数据")
    public BaseResultVO getPageList(@PathVariable int pageNo, @PathVariable int pageSize,
                                    @ModelAttribute PortalTenantQueryPageListForm form) {
        return JsonResultUtil.getResult(portalTenantManager.getPageList(pageNo, pageSize, form));
    }

    @RequestMapping(method = RequestMethod.PUT, value = "")
    @ResponseBody
    @OperateAction("UPDATE")
    @ApiOperation(value = "更新租户", notes = "更新租户")
    public BaseResultVO update(@ModelAttribute PortalTenantUpdateForm form) {
        return JsonResultUtil.getResult(portalTenantOperateManager.update(form));
    }

    @RequestMapping(method = RequestMethod.DELETE, value = "/{id}")
    @ResponseBody
    @OperateAction("DELETE")
    @ApiOperation(value = "删除租户", notes = "删除租户")
    public BaseResultVO delete(@PathVariable String id) {
        return JsonResultUtil.getResult(portalTenantOperateManager.delete(id));
    }
}
