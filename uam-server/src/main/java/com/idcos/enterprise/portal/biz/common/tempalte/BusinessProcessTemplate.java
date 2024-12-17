/**
 * yunjikeji Inc.
 * Copyright (c) 2004-2015 All Rights Reserved.
 */
package com.idcos.enterprise.portal.biz.common.tempalte;

import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.common.biz.BizCoreException;
import com.idcos.common.biz.CommonResultCode;
import com.idcos.enterprise.portal.biz.common.CommonBizException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.TransactionStatus;
import org.springframework.transaction.support.TransactionCallback;
import org.springframework.transaction.support.TransactionTemplate;

/**
 * 业务处理模板，这里封装了基本的业务处理处理流程
 * 基本的业务处理流程
 * <ul>
 * <li>检查请求参数基本信息</li>
 * <li>检查业务参数基本信息</li>
 * <li>处理业务信息，提供事务支持</li>
 * <li>返回处理结果</li>
 * </ul>
 *
 * @author yanlv
 * @version $Id: BusinessProcessTemplate.java, v 0.1 2015年3月18日 下午7:18:25 yanlv Exp $
 */
@Service
public class BusinessProcessTemplate {
    /**
     * 日志
     */
    private static final Logger LOGGER = LoggerFactory.getLogger(BusinessProcessTemplate.class);

    /**
     * 事务模板信息
     */
    @Autowired
    private TransactionTemplate transactionTemplate;

    /**
     * 业务实现类，返回一个默认的CommonResult
     *
     * @param callback 业务回调实现类,
     * @return
     */

    public <T> CommonResult<T> process(final BusinessProcessCallback<T> callback) {

        final CommonResult<T> result = new CommonResult<T>();

        final BusinessProcessContext context = new BusinessProcessContext();

        try {

            T object = transactionTemplate.execute(new TransactionCallback<T>() {

                @Override
                public T doInTransaction(TransactionStatus status) {

                    //对业务参数信息进行检查
                    callback.checkParam(context);

                    //对业务相关的类型进行check,比如唯一性，状态等等
                    callback.checkBusinessInfo(context);

                    //这里执行真正的业务逻辑，业务逻辑是在原子事务里面的
                    return callback.doBusiness(context);
                }

            });

            result.setSuccess(true);
            result.setResultObject(object);

        } catch (CommonBizException e) {

            LOGGER.error("执行业务操作异常", e);
            result.setSuccess(false);
            result.setResultCode(e.getResultCode().getCode());
            result.setResultMessage(e.getMessage());

            callback.exceptionProcess(e, context);

        } catch (BizCoreException e) {

            LOGGER.error("执行业务操作异常", e);
            result.setSuccess(false);
            result.setResultCode(CommonResultCode.UNKNOWN_EXCEPTION.getCode());
            result.setResultMessage(e.getMessage());

        } catch (Exception e) {
            LOGGER.error("执行业务操作异常", e);
            result.setSuccess(false);
            result.setResultMessage(e.getMessage());

        }

        return result;

    }
    //
    //    /**
    //     * 业务实现类，参数指定要传入的result对象信息
    //     * @param callback 业务回调实现类
    //     * @return
    //     */
    //
    //    public synchronized <T> void processWithResult(final BusinessProcessCallbackWithResult<T> callback,
    //                                                   final CommonResult<T> result) {
    //
    //        final BusinessProcessContext context = new BusinessProcessContext();
    //
    //        try {
    //
    //            T object = transactionTemplate.execute(new TransactionCallback<T>() {
    //
    //                @Override
    //                public T doInTransaction(TransactionStatus status) {
    //
    //                    //对业务参数信息进行检查
    //                    callback.checkParam(context, result);
    //
    //                    //对业务相关的类型进行check,比如唯一性，状态等等
    //                    callback.checkBusinessInfo(context, result);
    //
    //                    //这里执行真正的业务逻辑，业务逻辑是在原子事务里面的
    //                    return callback.doBusiness(context, result);
    //                }
    //
    //            });
    //
    //            result.setSuccess(true);
    //            result.setResultObject(object);
    //
    //        } catch (CommonBizException e) {
    //
    //            LOGGER.error("执行业务操作异常", e);
    //            result.setSuccess(false);
    //            result.setResultCode(e.getResultCode());
    //            result.setResultMessage(e.getMessage());
    //
    //            callback.exceptionProcess(e, context, result);
    //
    //        } catch (Exception e) {
    //            LOGGER.error("执行业务操作异常", e);
    //            result.setSuccess(false);
    //            result.setResultMessage(e.getMessage());
    //
    //        }
    //
    //    }

}
