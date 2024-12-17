/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.web;

/**
 * @author xizhao.dai
 * @version RestResponse.java, v1 2017/9/18 下午6:28 xizhao.dai Exp $$
 */
public class RestResponse {

    /**
     * 返回调用操作是否有错误,代表状态码，成功success，失败fail
     * <p>
     * 此处理的错误主要是业务层面的错误信息。
     * </p>
     */
    private String status;

    /**
     * 返回信息，不管操作是否成功，都可以返回相应的信息。
     */
    private String message;

    /**
     * 返回值，不管操作是否成功都可以返回相应的数据。
     */
    private Object content;

    /**
     * 状态代码，用不同的数字标示不同的返回状态
     */
    private String statusCode = ResponseStatusCode.RSP_STATUS_NORMAL;

    public RestResponse() {
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public Object getContent() {
        return content;
    }

    public void setContent(Object content) {
        this.content = content;
    }

    public String getStatusCode() {
        return statusCode;
    }

    public void setStatusCode(String statusCode) {
        this.statusCode = statusCode;
    }

}
