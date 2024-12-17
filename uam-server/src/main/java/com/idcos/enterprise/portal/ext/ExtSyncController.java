package com.idcos.enterprise.portal.ext;

import io.swagger.annotations.Api;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.ResponseBody;

/**
 * @Author: Dai
 * @Date: 2018/11/8 10:28 PM
 * @Description:
 */
@Controller
@RequestMapping(value = "/rbac/ext")
@Api(tags = "15.外部同步数据对接接口", description = "ExtSyncController")
public class ExtSyncController {
    @Autowired
    private ExtSyncService extSyncService;


    /**
     * 同步外部系统的部门和人员到UAM
     *
     * @return
     */
    @RequestMapping(value = "/sync", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public void syncDeptAndUser() {
        extSyncService.syncDeptAndUser();
    }
}

