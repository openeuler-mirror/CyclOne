/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2018 All Rights Reserved.
 */
package com.idcos.enterprise.portal.web.permission;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;

import io.swagger.annotations.Api;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.util.ResourceUtils;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;
import org.yaml.snakeyaml.Yaml;

import com.idcos.common.service.vo.CommonRestResult;

/**
 * @author Xizhao.Dai
 * @version PermissionController.java, v1 2018/7/19 下午5:40 Xizhao.Dai Exp $$
 */
@RestController
@RequestMapping("/uam/permission")
@Api(tags = "13.CloudUam获取自身权限的接口", description = "PermissionController")
public class PermissionController {
    private static final Logger logger = LoggerFactory.getLogger(PermissionController.class);

    @RequestMapping(value = "/menuCodes", method = RequestMethod.GET)
    public CommonRestResult<PermissionNode> menuCodesResource() {
        Yaml yaml = new Yaml();

        PermissionNode permissionNode = null;
        try {
            File file = ResourceUtils.getFile("classpath:menu.yaml");
            permissionNode = yaml.loadAs(new FileInputStream(file), PermissionNode.class);
        } catch (FileNotFoundException e) {
            logger.error("menu resource file menu.yaml not found", e);
        }

        return new CommonRestResult<>(permissionNode);
    }
}