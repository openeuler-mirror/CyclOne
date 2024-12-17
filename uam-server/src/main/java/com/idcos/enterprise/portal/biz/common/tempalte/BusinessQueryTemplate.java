/**
 * yunjikeji Inc.
 * Copyright (c) 2004-2015 All Rights Reserved.
 */
package com.idcos.enterprise.portal.biz.common.tempalte;

import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.common.biz.BizCoreException;
import com.idcos.common.biz.CommonResultCode;
import com.idcos.enterprise.portal.biz.common.CommonBizException;
import com.idcos.enterprise.portal.biz.common.AuthInfoException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

/**
 * 查询模板处理
 * 基本的业务处理流程
 * <ul>
 * <li>检查请求参数基本信息</li>
 * <li>检查业务参数基本信息</li>
 * <li>处理业务信息，提供事务支持</li>
 * <li>返回处理结果</li>
 * </ul>
 *
 * @author yanlv
 * @version $Id: BusinessQueryTemplate.java, v 0.1 2015年3月27日 下午1:20:08 yanlv Exp $
 */
@Service
public class BusinessQueryTemplate {

    /**
     * 日志
     */
    private static final Logger logger = LoggerFactory.getLogger(BusinessQueryTemplate.class);

    /**
     * @param callback 业务回调实现类
     * @return
     */
    public <T> CommonResult<T> process(final BusinessQueryCallback<T> callback) {

        final CommonResult<T> result = new CommonResult<T>();

        try {

            //对业务参数信息进行检查
            callback.checkParam();

            //这里执行真正的业务逻辑，业务逻辑是在原子事务里面的
            T object = callback.doQuery();

            result.setSuccess(true);
            result.setResultObject(object);

            return result;

        } catch (AuthInfoException e) {

            logger.error("查询用户信息异常:{}", e.getMessage());
            result.setSuccess(false);
            result.setResultCode(e.getResultCode().getCode());
            result.setResultMessage(e.getMessage());

            return result;

        } catch (CommonBizException e) {

            logger.error("执行查询操作异常", e);
            result.setSuccess(false);
            result.setResultCode(e.getResultCode().getCode());
            result.setResultMessage(e.getMessage());

            return result;

        } catch (BizCoreException e) {

            logger.error("执行查询操作异常", e);
            result.setSuccess(false);
            result.setResultCode(CommonResultCode.UNKNOWN_EXCEPTION.getCode());
            result.setResultMessage(e.getMessage());

            return result;

        } catch (Exception e) {

            logger.error("执行查询操作异常", e);
            result.setSuccess(false);
            result.setResultMessage(e.getMessage());
            return result;
        }

    }

}
