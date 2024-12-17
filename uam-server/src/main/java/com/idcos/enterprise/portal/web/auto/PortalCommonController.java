

package com.idcos.enterprise.portal.web.auto;

// auto generated imports

import io.swagger.annotations.Api;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;

import com.idcos.cloud.biz.common.permisson.OperateAction;
import com.idcos.cloud.web.common.BaseResultVO;
import com.idcos.cloud.web.common.JsonResultUtil;
import com.idcos.enterprise.portal.manager.auto.PortalCommonQueryManager;

/**
 * 统一门户公共查询管理
 * web层controller相关的接口自动生成，此文件属于自动生成的，请勿直接修改,具体可以参考codegen工程
 * Generated by <tt>controller-codegen</tt> on 2015-08-21 10:22:49.
 *
 * @author yanlv
 * @version PortalCommonController.java, v 1.1 2015-08-21 10:22:49 yanlv Exp $
 */

@Controller
@OperateAction("F_PORTAL_COMMON")
@RequestMapping(value = "/portal/common")
@Api(tags = "10.统一门户公共查询管理接口", description = "PortalCommonController")
public class PortalCommonController {

    //========== manager ==========

    @Autowired
    private PortalCommonQueryManager portalCommonQueryManager;

    /**
     * 获取所有角色信息
     *
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.GET, value = "/roleListAction")
    @ResponseBody
    @OperateAction("O_PORTAL_COMMON_ROLELIST")
    public BaseResultVO getRoleList() {
        return JsonResultUtil.getResult(portalCommonQueryManager.getRoleList());

    }

    /**
     * 获取所有用户组信息
     *
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.GET, value = "/userGroupListAction")
    @ResponseBody
    @OperateAction("O_PORTAL_COMMON_USERGROUPLIST")
    public BaseResultVO getUserGroupList(@RequestParam("tenantId") String tenantId) {
        return JsonResultUtil.getResult(portalCommonQueryManager.getUserGroupList(tenantId));

    }

}
