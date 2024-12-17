package com.idcos.enterprise.portal.web.manual;

import com.idcos.cloud.biz.common.permisson.OperateAction;
import com.idcos.common.service.vo.CommonRestResult;
import com.idcos.enterprise.portal.dal.enums.RbacMenuEnum;
import com.idcos.enterprise.portal.vo.Item;
import io.swagger.annotations.Api;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.context.WebApplicationContext;
import org.springframework.web.context.support.WebApplicationContextUtils;
import org.springframework.web.method.HandlerMethod;
import org.springframework.web.servlet.mvc.method.RequestMappingInfo;
import org.springframework.web.servlet.mvc.method.annotation.RequestMappingHandlerMapping;

import javax.servlet.http.HttpServletRequest;
import java.util.Map;

/**
 * 权限资源管理系统对外提供权限资源信息的接口
 *
 * @author pengganyu
 * @version $Id: PortalPermissionResourceService.java, v 0.1 2016年6月2日 下午4:36:41 pengganyu Exp $
 */
@RestController
@RequestMapping("/api/v1/rbac")
@Api(tags = "11.对外提供权限资源信息的接口", description = "PortalPermissionResourceService")
public class PortalPermissionResourceService {
    @RequestMapping(value = "/operateCodes", method = RequestMethod.GET)
    public CommonRestResult<Item> operateCodesResource(HttpServletRequest request) {

        Item root = new Item();
        root.setId("");
        root.setTitle("RBAC操作权限列表");

        WebApplicationContext wc = WebApplicationContextUtils
                .getRequiredWebApplicationContext(request.getSession().getServletContext());
        RequestMappingHandlerMapping handlerMapping = wc
                .getBean(RequestMappingHandlerMapping.class);
        Map<RequestMappingInfo, HandlerMethod> map = handlerMapping.getHandlerMethods();
        for (Map.Entry<RequestMappingInfo, HandlerMethod> entry : map.entrySet()) {
            HandlerMethod handlerMethod = entry.getValue();
            //class level
            OperateAction operateAction = handlerMethod.getBeanType().getAnnotation(
                    OperateAction.class);
            if (operateAction == null) {
                //假设method级别的annotation不能单独存在
                continue;
            }
            Item item = new Item();
            item.setId(operateAction.value());
            item.setTitle(operateAction.value());

            Item classLevel = root.addChild(item);

            OperateAction methodLevelOperateAction = handlerMethod
                    .getMethodAnnotation(OperateAction.class);
            if (methodLevelOperateAction != null) {
                Item methodLevel = new Item();
                methodLevel.setId(methodLevelOperateAction.value());
                methodLevel.setTitle(methodLevelOperateAction.value());
                classLevel.addChild(methodLevel);
            }

        }
        return new CommonRestResult<>(root);

    }

    @RequestMapping(value = "/menuCodes", method = RequestMethod.GET)
    public CommonRestResult<Item> menuCodesResource() {
        Item root = new Item();

        root.setId("");
        root.setTitle("RBAC菜单列表");

        // 遍历枚举，组装菜单数据
        for (RbacMenuEnum rbEnum : RbacMenuEnum.values()) {
            Item childLevel = new Item();

            childLevel.setId(rbEnum.getCode());
            childLevel.setTitle(rbEnum.getName());

            root.addChild(childLevel);
        }

        return new CommonRestResult<>(root);
    }
}
