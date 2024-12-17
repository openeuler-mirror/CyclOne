/**
 * yunjikeji Inc.
 * Copyright (c) 2004-2015 All Rights Reserved.
 */
package com.idcos.enterprise.portal.biz.common;

/**
 * 业务处理结果码
 *
 * @author yanlv
 * @version $Id: ResultCode.java, v 0.1 2015-1-29 下午7:50:32 yanlv Exp $
 */
public enum ResultCode {
    /**
     * 处理失败
     */
    FAILE("FAILE", "处理失败"),
    /**
     * 处理成功
     */
    SUCCESS("success", "处理成功"),

    PARAM_NULL("PARAM_NULL", "参数为空"),

    PARAM_ERROR("PARAM_ERROR", "参数错误"),

    PAGE_OFFSET_LIMIT_PARAM_ERROR("PAGE_OFFSET_LIMIT_PARAM_ERROR", "翻页查询条件参数错误"),

    QUERY_RESULT_IS_NULL("QUERY_RESULT_IS_NULL", "查询结果为空"),

    DEVICE_NOT_EXIST_FIND_BY_SN("DEVICE_NOT_EXIST_FIND_BY_SN", "根据设备序列号无法获取对应设备信息"),

    DEVICE_SN_REPEAT("DEVICE_SN_REPEAT", "设备序列号重复"),

    UNKNOWN_EXCEPTION("UNKNOWN_EXCEPTION", "未知异常"),

    USER_NOT_EXIST("USER_NOT_EXIST", "用户信息不存在"),

    PASSWORD_NOT_CORRECT("PASSWORD_NOT_CORRECT", "用户密码错误"),

    USER_AUTH_EXCEPTION("USER_AUTH_EXCEPTION", "用户验证错误"),

    UNAUTHORIZED("UNAUTHORIZED", "未授权"),

    NO_AUTH("NO_AUTH", "未认证"),

    AUTH_FAIL("AUTH_FAIL", "认证失败"),

    INSTRUCTION_EXE_TIMEOUT("INSTRUCTION_EXE_TIMEOUT", "指令执行超时"),

    INSTRUCTION_EXE_WAY_NOT_SUPPORT("INSTRUCTION_EXE_WAY_NOT_SUPPORT", "指令执行方式不支持"),

    STATUS_ERROR("STATUS_ERROR", "状态异常");

    /**
     * 结果代码
     */
    private String code;

    /**
     * 结果描述
     */
    private String description;

    /**
     * @param code
     * @param description
     */
    private ResultCode(String code, String description) {

        this.code = code;
        this.description = description;
    }

    /**
     * Getter method for property <tt>code</tt>.
     *
     * @return property value of code
     */
    public String getCode() {
        return code;
    }

    /**
     * Getter method for property <tt>description</tt>.
     *
     * @return property value of description
     */
    public String getDescription() {
        return description;
    }

}
