package com.idcos.enterprise.portal.biz.common;

import com.idcos.common.biz.BizCoreException;
import org.springframework.util.Assert;

/**
 * @author Xizhao.Dai
 * @version CommonBizException2.java, v1 2018/5/28 下午9:12 Xizhao.Dai Exp $$
 */
public class AuthInfoException extends BizCoreException {

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
    public AuthInfoException(ResultCode resultCode, String message) {

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
    public AuthInfoException(ResultCode resultCode) {

        super(resultCode.getDescription());

        Assert.notNull(resultCode, resultCode.getDescription());

        this.resultCode = resultCode;
    }

    /**
     * @param message
     */
    public AuthInfoException(String message) {

        super(message);

        this.resultCode = ResultCode.UNKNOWN_EXCEPTION;
    }

    /**
     * @param message
     */
    public AuthInfoException(String message, Object exceptionDescription) {

        super(message);

        this.resultCode = ResultCode.UNKNOWN_EXCEPTION;

        this.exceptionDescription = exceptionDescription;
    }

    /**
     * @param resultCode
     * @param message
     * @param e
     */
    public AuthInfoException(ResultCode resultCode, String message, Throwable e) {

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
    public AuthInfoException(String message, Throwable e) {

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