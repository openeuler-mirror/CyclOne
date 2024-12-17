

package com.idcos.enterprise.portal.web.manual;

// auto generated imports

import io.swagger.annotations.Api;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;

import com.idcos.cloud.biz.common.permisson.OperateAction;
import com.idcos.enterprise.portal.biz.common.PortalResponse;
import com.idcos.enterprise.portal.manager.auto.PortalRestfulService;

/**
 * @author
 * @version PortalFilterUtil.java, v1 2017/11/18 下午6:22  Exp $$
 */
@Controller
@OperateAction("F_PORTAL_USER")
@RequestMapping(value = "/api/v1")
@Api(tags = "12./api/v1接口", description = "PortalRestfulController")
public class PortalRestfulController {
    //========== manager ==========

    @Autowired
    private PortalRestfulService portalRestfulService;

    /**
     * 根据权限资源类型获取权限资源
     *
     * @param resType
     * @return
     */
    @RequestMapping(method = RequestMethod.GET, value = "/res")
    @ResponseBody
    @OperateAction("O_PORTAL_USER_GET")
    public PortalResponse queryResByResType(@RequestParam("resType") String resType, @RequestParam(name = "tenantId", required = false) String tenantId) {
        return portalRestfulService.queryResource(resType, tenantId);
    }

    /**
     * queryUserIdsByGroupId
     *
     * @param groupId
     * @return
     */
    @RequestMapping(method = RequestMethod.GET, value = "/{groupId}/userIds")
    @ResponseBody
    @OperateAction("O_PORTAL_USER_GET")
    public PortalResponse queryUserIdsById(@PathVariable String groupId) {
        return portalRestfulService.queryUserIdsByGroupId(groupId);
    }

    /**
     * queryUserIdsByGroupName
     *
     * @param groupName
     * @return
     */
    @RequestMapping(method = RequestMethod.GET, value = "/userIds")
    @ResponseBody
    @OperateAction("O_PORTAL_USER_GET")
    public PortalResponse queryUserIdsByGroupName(@RequestParam("groupName") String groupName) {
        return portalRestfulService.queryUserIdsByGroupName(groupName);
    }

    /**
     * 根据查询条件查询用户信息
     *
     * @param tenantId
     * @param cnd
     * @param pageNo
     * @param pageSize
     * @return
     */
    @RequestMapping(method = RequestMethod.GET, value = "/users/cnd")
    @ResponseBody
    @OperateAction("O_PORTAL_USER_GET")
    public PortalResponse queryUsersByCnd(@RequestParam(name = "tenantId", required = false) String tenantId,
                                          @RequestParam("cnd") String cnd, @RequestParam("pageNo") String pageNo,
                                          @RequestParam("pageSize") String pageSize) {
        return portalRestfulService.queryUserByCnd(tenantId, cnd, pageNo, pageSize);
    }
}
