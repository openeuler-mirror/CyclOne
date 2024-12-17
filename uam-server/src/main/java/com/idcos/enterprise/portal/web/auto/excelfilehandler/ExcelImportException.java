/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.web.auto.excelfilehandler;

/**
 * @author Dana
 * @version ServiceException.java, v1 2017/12/10 下午6:14 Dana Exp $$
 */
public class ExcelImportException extends Exception {

    /**
     * serialUID.
     */
    private static final long serialVersionUID = -7902772542663897926L;

    /**
     *
     */
    public ExcelImportException() {
    }

    /**
     * @param message
     */
    public ExcelImportException(String message) {
        super(message);
    }

    /**
     * @param cause
     */
    public ExcelImportException(Throwable cause) {
        super(cause);
    }

    /**
     * @param message
     * @param cause
     */
    public ExcelImportException(String message, Throwable cause) {
        super(message, cause);
    }

    /**
     * @param message
     * @param cause
     * @param enableSuppression
     * @param writableStackTrace
     */
    public ExcelImportException(String message, Throwable cause, boolean enableSuppression, boolean writableStackTrace) {
        super(message, cause, enableSuppression, writableStackTrace);
    }
}