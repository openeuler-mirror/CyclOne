/**
 * yunjikeji Inc.
 * Copyright (c) 2004-2015 All Rights Reserved.
 */
package com.idcos.enterprise.portal.biz.common;

import com.idcos.common.biz.BizCoreException;
import org.springframework.util.Assert;

/**
 * idcos管理平台运行时异常，作为业务处理的异常基类
 *
 * @author yanlv
 * @version $Id: idcosBizException.java, v 0.1 2015-1-29 下午7:47:55 yanlv Exp $
 */
public class CommonBizException extends BizCoreException {

    /**
     * serialVersionUID
     */
    private static final long serialVersionUID = 1382559230799468711L;

    /**
     * 结果码处理
     */
    private ResultCode resultCode;

    /**
     * 异常描述对象，可以返回到json前台的信息
     */
    private Object exceptionDescription;

    /**
     * 构造函数
     *
     * @param resultCode
     * @param message
     */
    public CommonBizException(ResultCode resultCode, String message) {

        super(message);

        Assert.notNull(resultCode, message);

        this.resultCode = resultCode;
    }

    /**
     * 构造函数
     *
     * @param resultCode
     * @param resultCode
     */
    public CommonBizException(ResultCode resultCode) {

        super(resultCode.getDescription());

        Assert.notNull(resultCode, resultCode.getDescription());

        this.resultCode = resultCode;
    }

    /**
     * @param message
     */
    public CommonBizException(String message) {

        super(message);

        this.resultCode = ResultCode.UNKNOWN_EXCEPTION;
    }

    /**
     * @param message
     */
    public CommonBizException(String message, Object exceptionDescription) {

        super(message);

        this.resultCode = ResultCode.UNKNOWN_EXCEPTION;

        this.exceptionDescription = exceptionDescription;
    }

    /**
     * @param resultCode
     * @param message
     * @param e
     */
    public CommonBizException(ResultCode resultCode, String message, Throwable e) {

        super(message, e);

        if (resultCode == null) {
            this.resultCode = ResultCode.UNKNOWN_EXCEPTION;
        } else {
            this.resultCode = resultCode;
        }
    }

    /**
     * @param message
     * @param e
     */
    public CommonBizException(String message, Throwable e) {

        super(message, e);

        this.resultCode = ResultCode.UNKNOWN_EXCEPTION;

    }

    public ResultCode getResultCode() {
        return resultCode;
    }

    public Object getExceptionDescription() {
        return exceptionDescription;
    }

}
