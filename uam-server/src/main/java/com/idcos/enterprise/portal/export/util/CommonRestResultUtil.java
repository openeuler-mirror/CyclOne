/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.export.util;

import com.idcos.cloud.core.common.biz.CommonResult;

/**
 * @author Dana
 * @version CommonRestResultUtil.java, v1 2017/10/14 下午5:22 Dana Exp $$
 */
public class CommonRestResultUtil {
    /**
     * 根据CommonResult返回结果信息
     *
     * @param result
     * @return
     */
    @SuppressWarnings("unchecked")
    public static final CommonRestResult getResult(CommonResult<?> result) {
        CommonRestResult commonRestResult = new CommonRestResult(result.isSuccess() ? "success" : "fail",
                result.getResultMessage());
        commonRestResult.setContent(result.getResultObject());
        commonRestResult.setResultCode(result.getResultCode());
        return commonRestResult;
    }
}