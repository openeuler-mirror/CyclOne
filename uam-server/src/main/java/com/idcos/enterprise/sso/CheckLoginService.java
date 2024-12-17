package com.idcos.enterprise.sso;

import com.idcos.enterprise.portal.biz.common.CommonBizException;
import com.idcos.enterprise.portal.biz.common.utils.PasswordUtil;
import com.idcos.enterprise.portal.dal.entity.PortalTenant;
import com.idcos.enterprise.portal.dal.entity.PortalUser;
import com.idcos.enterprise.portal.dal.enums.PortalUserStatusEnum;
import com.idcos.enterprise.portal.dal.repository.PortalTenantRepository;
import com.idcos.enterprise.portal.dal.repository.PortalUserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

/**
 * @Author: Dai
 * @Date: 2018/10/8 下午7:33
 * @Description:
 */
@Service
public class CheckLoginService {
    @Autowired
    private PortalUserRepository portalUserRepository;

    @Autowired
    private PortalTenantRepository portalTenantRepository;

    public PortalUser checkUser(String loginId, String tenantId) {
        //验证用户、租户和密码是否正确
        PortalUser portalUser = portalUserRepository.findPortalUserById(tenantId, loginId);
        if (portalUser == null || !portalUser.getStatus().equals(PortalUserStatusEnum.ENABLED.getCode())) {
            throw new CommonBizException("您输入的用户或密码不正确，请检查。");
        }
        return portalUser;
    }

    public void checkPassword(PortalUser portalUser, String password) {
        //验证密码是否匹配
        String pw = PasswordUtil.encryptPassword(portalUser.getId(), password, portalUser.getSalt());
        if (!pw.equals(portalUser.getPassword())) {
            throw new CommonBizException("您输入的用户或密码不正确，请检查。");
        }
    }

    public PortalTenant checkTenant(String tenantId) {
        //验证租户
        PortalTenant portalTenant = portalTenantRepository.findByName(tenantId);
        if (portalTenant == null) {
            throw new CommonBizException("租户【" + tenantId + "】不存在，请检查。");
        }
        return portalTenant;
    }
}
