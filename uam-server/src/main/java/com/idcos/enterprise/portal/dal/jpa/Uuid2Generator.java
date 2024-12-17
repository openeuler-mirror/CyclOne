/*
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2018 All Rights Reserved.
 */
package com.idcos.enterprise.portal.dal.jpa;

import com.idcos.cloud.core.common.BaseVO;
import com.idcos.cloud.core.common.util.FieldUtil;
import org.hibernate.HibernateException;
import org.hibernate.engine.spi.SessionImplementor;
import org.hibernate.id.UUIDGenerator;

import java.io.Serializable;

/**
 * 自定义UUID生成器
 *
 * @author sevenlin
 * @version Uuid2Generator.java, v1 2018/9/27 14:42 sevenlin Exp $$
 */
public class Uuid2Generator extends UUIDGenerator {

    @Override
    public Serializable generate(SessionImplementor session, Object object) throws HibernateException {
        // 如果存在自定义ID则使用自定义ID
        if (object instanceof BaseVO) {
            Object id = FieldUtil.readField(object, "id");
            if (id != null) {
                return (Serializable) id;
            }
        }
        return super.generate(session, object);
    }
}
