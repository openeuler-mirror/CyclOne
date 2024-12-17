/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.web;

/**
 * 返回状态码常量
 *
 * @author xizhao.dai
 * @version ResponseStatusCode.java, v1 2017/9/18 下午6:28 xizhao.dai Exp $$
 */
public class ResponseStatusCode {

    /**
     * 返回状态码 －－ 正常状态
     */
    public final static String RSP_STATUS_NORMAL = "1000";

    /**
     * 返回状态码 －－ 登录超时
     */
    public final static String RSP_STATUS_LOGINOUT = "7001";

    /**
     * 返回状态码 －－ 非法请求token.
     */
    public final static String RSP_STATUS_ILLEGAL_TOKEN = "7002";

}
