package com.idcos.enterprise.portal.convert;

import com.idcos.enterprise.portal.biz.common.convert.BaseConvertFunction;
import com.idcos.enterprise.portal.dal.entity.PortalToken;
import com.idcos.enterprise.portal.vo.PortalTokenVO;
import org.springframework.stereotype.Service;

/**
 * @author Dana
 * @version PortalTokenConvert.java, v1 2017/11/30 下午5:09 Dana Exp $$
 */
@Service
public class PortalTokenConvert extends BaseConvertFunction<PortalToken, PortalTokenVO> {
    @Override
    public PortalTokenVO apply(PortalToken input) {
        PortalTokenVO vo = super.apply(input);
        return vo;
    }

}