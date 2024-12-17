/**
 * yunjikeji Inc.
 * Copyright (c) 2004-2015 All Rights Reserved.
 */
package com.idcos.enterprise.portal.biz.common.tempalte;

/**
 * 查询回调处理
 *
 * @author yanlv
 * @version $Id: BusinessQueryCallback.java, v 0.1 2015年3月27日 下午1:21:21 yanlv Exp $
 */
public interface BusinessQueryCallback<T> {

    /**
     * 检查参数信息，是否满足查询条件
     */
    void checkParam();

    /**
     * 进行查询相关的业务处理
     *
     * @return T
     */
    T doQuery();

}
