package com.idcos.enterprise.portal.convert;

import org.springframework.stereotype.Service;

import com.idcos.enterprise.portal.biz.common.convert.BaseConvertFunction;
import com.idcos.enterprise.portal.dal.entity.PortalResource;
import com.idcos.enterprise.portal.dal.enums.BusinessIdentityEnum;
import com.idcos.enterprise.portal.vo.PortalResourceVO;

/**
 * {controller.converClassName}对象转化类 , 第一次只是生成一个默认的convert
 * <p>第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖</p>
 * <p>这个类是泛型，需要你确定转化的target是什么，ConvertFunction中 Object 需要指定的是待转化为VO的对象class type，这里默认是Object</p>
 *
 * @author yanlv
 * @version PortalPermissionConvert.java, v 1.1 2015-06-09 10:00:37 yanlv Exp $
 */
@Service
public class PortalResourceConvert extends BaseConvertFunction<PortalResource, PortalResourceVO> {
    @Override
    public BusinessIdentityEnum getTabEnum() {
        return BusinessIdentityEnum.PORTAL_RESOURCE;
    }

    @Override
    public PortalResourceVO apply(PortalResource input) {
        PortalResourceVO vo = super.apply(input);
        vo.setTenantId(input.getTenant());
        return vo;
    }
}