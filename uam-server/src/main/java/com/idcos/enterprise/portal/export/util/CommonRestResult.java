/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.export.util;

import java.util.List;

/**
 * @author Dana
 * @version CommonRestResult.java, v1 2017/11/27 下午7:30 Dana Exp $$
 */
public class CommonRestResult<T> {
    /**
     * 返回结果
     */
    private T content;
    /**
     * 业务的处理结果
     */
    private String status = "success";
    /**
     *
     */
    private String message;

    private String resultCode;
    /**
     * 更多错误
     */
    private List<?> errors;

    /**
     * 默认构造函数
     */
    public CommonRestResult() {

    }

    public CommonRestResult(final T content) {
        this();
        this.content = content;
    }

    public CommonRestResult(final String status, final String message) {
        this(null);
        this.status = status;
        this.message = message;
    }

    public CommonRestResult(final String status, final String message, final String resultCode) {
        this(null);
        this.status = status;
        this.message = message;
        this.resultCode = resultCode;
    }

    public T getContent() {
        return content;
    }

    public void setContent(T content) {
        this.content = content;
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

    public boolean success() {
        return "success".equals(status);
    }

    public List<?> getErrors() {
        return errors;
    }

    public void setErrors(List<?> errors) {
        this.errors = errors;
    }

    public String getResultCode() {
        return resultCode;
    }

    public void setResultCode(String resultCode) {
        this.resultCode = resultCode;
    }
}