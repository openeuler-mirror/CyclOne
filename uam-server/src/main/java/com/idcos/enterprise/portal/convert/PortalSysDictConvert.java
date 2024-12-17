package com.idcos.enterprise.portal.convert;

import org.springframework.stereotype.Service;

import com.idcos.enterprise.portal.biz.common.convert.BaseConvertFunction;
import com.idcos.enterprise.portal.dal.entity.PortalSysDict;
import com.idcos.enterprise.portal.vo.PortalSysDictVO;

/**
 * @author Dana
 * @version PortalSysDictConvert.java, v1 2017/12/5 上午9:58 Dana Exp $$
 */
@Service
public class PortalSysDictConvert extends BaseConvertFunction<PortalSysDict, PortalSysDictVO> {
    @Override
    public PortalSysDictVO apply(PortalSysDict input) {
        PortalSysDictVO vo = super.apply(input);
        return vo;
    }
}