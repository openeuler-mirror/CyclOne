package com.idcos.enterprise.portal.convert;

import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.idcos.enterprise.portal.biz.common.convert.BaseConvertFunction;
import com.idcos.enterprise.portal.dal.entity.PortalGroupRoleRel;
import com.idcos.enterprise.portal.dal.entity.PortalGroupUserRel;
import com.idcos.enterprise.portal.dal.entity.PortalUserGroup;
import com.idcos.enterprise.portal.dal.enums.BusinessIdentityEnum;
import com.idcos.enterprise.portal.dal.repository.PortalGroupRoleRelRepository;
import com.idcos.enterprise.portal.dal.repository.PortalGroupUserRelRepository;
import com.idcos.enterprise.portal.vo.PortalUserGroupVO;

/**
 * 用户组PO2VO转换类
 *
 * @author jiaohuizhe
 * @version $Id: PortalUserGroupVOConvert.java, v 0.1 2015年5月8日 下午12:33:32 jiaohuizhe Exp $
 */
@Service
public class PortalUserGroupInfoConvert extends BaseConvertFunction<PortalUserGroup, PortalUserGroupVO> {

    @Autowired
    private PortalGroupRoleRelRepository portalGroupRoleRelRepository;
    @Autowired
    private PortalGroupUserRelRepository portalGroupUserRelRepository;

    @Override
    public PortalUserGroupVO apply(PortalUserGroup input) {
        PortalUserGroupVO vo = super.apply(input);
        vo.setTenantId(input.getTenant());
        // 加载角色信息，加载权限信息
        {
            List<PortalGroupRoleRel> list = portalGroupRoleRelRepository.findByGroupId(input
                    .getId());
            String[] roles = new String[list.size()];

            for (int i = 0; i < roles.length; i++) {
                roles[i] = list.get(i).getRoleId();
            }
            vo.setSelRoles(roles);

        }

        // 加载用户信息
        {
            List<PortalGroupUserRel> list = portalGroupUserRelRepository.findByGroupId(input
                    .getId());
            String[] users = new String[list.size()];

            for (int i = 0; i < users.length; i++) {
                users[i] = list.get(i).getUserId();
            }
            vo.setSelUsers(users);
        }
        return vo;
    }

    @Override
    public BusinessIdentityEnum getTabEnum() {
        return BusinessIdentityEnum.PORTAL_USER_GROUP;
    }

}