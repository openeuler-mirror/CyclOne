/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2016 All Rights Reserved.
 */
package com.idcos.enterprise.swagger;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.util.StopWatch;
import springfox.documentation.builders.ParameterBuilder;
import springfox.documentation.builders.PathSelectors;
import springfox.documentation.builders.RequestHandlerSelectors;
import springfox.documentation.schema.ModelRef;
import springfox.documentation.service.ApiInfo;
import springfox.documentation.service.Contact;
import springfox.documentation.service.Parameter;
import springfox.documentation.spi.DocumentationType;
import springfox.documentation.spring.web.plugins.Docket;
import springfox.documentation.swagger2.annotations.EnableSwagger2;

import java.util.Arrays;

/**
 * 启动swagger进行API描述自动生成在线文档。
 *
 * @author aol_aog
 * @version $Id: SwaggerConfiguration.java, v 0.1 2016年4月9日 上午11:07:09 aol_aog Exp $
 */
@Configuration
@EnableSwagger2
public class SwaggerConfiguration {

    private final Logger log = LoggerFactory.getLogger(SwaggerConfiguration.class);
    Parameter parameterAccessToken = new ParameterBuilder().name("access-token")
            .description("Access token identifier").modelRef(new ModelRef("string"))
            .parameterType("header").required(false).build();

    /**
     * 第三方系统的接口API
     *
     * @return
     */
    @Bean
    public Docket thirdPartSysDocket() {
        log.debug("Directory系统接口API,version: v1");
        StopWatch watch = new StopWatch();
        watch.start();
        Docket swaggerSpringMvcPlugin = new Docket(DocumentationType.SWAGGER_2)
                .groupName("Directory_V1_API").apiInfo(apiInfo()).select()
                .apis(RequestHandlerSelectors.basePackage("com.idcos.enterprise"))
                .paths(PathSelectors.any())
                .build();
        swaggerSpringMvcPlugin.globalOperationParameters(Arrays.asList(new Parameter[]{parameterAccessToken}));
        watch.stop();
        log.info("Started {} Swagger2 in {} ms", "Directory系统接口API", watch.getTotalTimeMillis());
        return swaggerSpringMvcPlugin;
    }

    private ApiInfo apiInfo() {
        return new ApiInfo("rbac系统内部模块间API接口说明", "主要描述相关API接口的详细信息，包括输入、输出参数、调用方式等信息。", "1.0", "无",
                new Contact("", "", ""), "", "无");
    }
}
