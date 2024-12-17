

package com.idcos.enterprise.portal.web.auto;

// auto generated imports

import io.swagger.annotations.Api;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.ResponseBody;

import com.idcos.cloud.biz.common.permisson.OperateAction;
import com.idcos.cloud.web.common.BaseResultVO;
import com.idcos.cloud.web.common.JsonResultUtil;
import com.idcos.enterprise.portal.form.PortalResourceCreateForm;
import com.idcos.enterprise.portal.form.PortalResourceUpdateForm;
import com.idcos.enterprise.portal.manager.auto.PortalResourceOperateManager;
import com.idcos.enterprise.portal.manager.auto.PortalResourceQueryManager;

/**
 * PortalResourceController
 * 权限资源管理Controller
 *
 * @author pengganyu
 * @version $Id: PortalResourceController.java, v 0.1 2016年5月11日 下午5:11:26 pengganyu Exp $
 */
@Controller
@OperateAction("F_PORTAL_RESOURCE")
@RequestMapping(value = "/portal/res")
@Api(tags = "09.权限资源管理的接口", description = "PortalResourceController")
public class PortalResourceController {

    @Autowired
    private PortalResourceQueryManager portalResourceQueryManager;
    @Autowired
    private PortalResourceOperateManager portalResourceOperateManager;

    /**
     * 根据权限资源编码查询权限信息
     *
     * @param code
     * @return
     */
    @RequestMapping(method = RequestMethod.GET, value = "/{code}")
    @ResponseBody
    @OperateAction("O_PORTAL_USER_GET")
    public BaseResultVO queryByCode(@PathVariable String code) {
        return JsonResultUtil.getResult(portalResourceQueryManager.queryByCode(code));

    }

    /**
     * 查询所有权限信息
     *
     * @return
     */
    @RequestMapping(method = RequestMethod.GET, value = "/all")
    @ResponseBody
    @OperateAction("O_PORTAL_USER_GET")
    public BaseResultVO queryAll() {
        return JsonResultUtil.getResult(portalResourceQueryManager.queryAll());

    }

    /**
     * 删除权限资源
     *
     * @param id
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.DELETE, value = "/{id}")
    @ResponseBody
    @OperateAction("O_PORTAL_USER_DELETE")
    public BaseResultVO delete(@PathVariable String id) {
        return JsonResultUtil.getResult(portalResourceOperateManager.delete(id));
    }

    /**
     * 更新权限资源
     *
     * @param id
     * @param form
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.PUT, value = "/{id}")
    @ResponseBody
    @OperateAction("O_PORTAL_USER_PUT")
    public BaseResultVO update(@PathVariable String id, PortalResourceUpdateForm form) {
        return JsonResultUtil.getResult(portalResourceOperateManager.update(id, form));

    }

    /**
     * 增加权限资源
     *
     * @param form
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.POST, value = "/")
    @ResponseBody
    @OperateAction("O_PORTAL_USER_CREATE")
    public BaseResultVO create(PortalResourceCreateForm form) {
        return JsonResultUtil.getResult(portalResourceOperateManager.create(form));

    }

}
