package com.idcos.enterprise.portal.biz.common.convert;

import com.google.common.base.Function;
import com.idcos.cloud.core.common.util.ObjectUtil;
import com.idcos.enterprise.portal.dal.enums.BusinessIdentityEnum;

import java.lang.reflect.ParameterizedType;

/**
 * PO2VO转换类
 *
 * @author jiaohuizhe
 * @version $Id: ConvertFunction.java, v 0.1 2015年5月8日 上午9:35:17 jiaohuizhe Exp $
 */
public abstract class BaseConvertFunction<F, T> implements Function<F, T> {
    private Class<T> entityClass;

    @Override
    @SuppressWarnings("unchecked")
    public T apply(F input) {
        if (entityClass == null) {
            entityClass = (Class<T>) ((ParameterizedType) getClass().getGenericSuperclass())
                    .getActualTypeArguments()[1];
        }
        T t = ObjectUtil.copy(input, entityClass);
        return t;
    }

    public BusinessIdentityEnum getTabEnum() {
        return null;
    }

}
