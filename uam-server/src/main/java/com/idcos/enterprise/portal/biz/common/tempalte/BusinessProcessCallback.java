/**
 * yunjikeji Inc.
 * Copyright (c) 2004-2015 All Rights Reserved.
 */
package com.idcos.enterprise.portal.biz.common.tempalte;

import com.idcos.enterprise.portal.biz.common.CommonBizException;

/**
 * 业务处理回调接口,负责处理具体的业务信息
 *
 * @author yanlv
 * @version $Id: BusinessProcessCallback.java, v 0.1 2015年3月18日 下午7:19:29 yanlv Exp $
 */
public interface BusinessProcessCallback<T> {

    /**
     * 检查参数信息，如果失败，抛出<code>CommonBizException</code>异常
     *
     * @param context
     */
    void checkParam(BusinessProcessContext context);

    /**
     * 检查参数信息，如果失败，抛出<code>CommonBizException</code>异常
     *
     * @param context
     */
    void checkBusinessInfo(BusinessProcessContext context);

    /**
     * 做相关的业务处理，提供事务相关的支持,如果失败，抛出<code>CommonBizException</code>异常
     *
     * @param context
     * @return T
     */
    T doBusiness(BusinessProcessContext context);

    /**
     * 发生异常的时候处理方式
     *
     * @param exception 异常信息
     * @param context   上下文信息
     */
    void exceptionProcess(CommonBizException exception, BusinessProcessContext context);

}
