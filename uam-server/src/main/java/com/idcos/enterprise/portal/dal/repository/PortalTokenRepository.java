/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.dal.repository;

import com.idcos.cloud.core.dal.common.jpa.BaseRepository;
import com.idcos.enterprise.portal.dal.entity.PortalToken;
import org.springframework.data.jpa.repository.Query;

import java.util.List;

/**
 * @author Dana
 * @version PortalTokenRepository.java, v1 2017/11/30 上午1:18 Dana Exp $$
 */
public interface PortalTokenRepository extends BaseRepository<PortalToken, String> {
    /**
     * 根据tenantId和loginId获取token列表
     * 过期时间大于当前时间，代表token未过期
     *
     * @param tenantId
     * @param loginId
     * @return
     */
    @Query("select t from PortalToken t where t.tenantId = ?1 and t.loginId = ?2 and t.expireTime > current_timestamp")
    List<PortalToken> listTokenByTenantIdAndLoginId(String tenantId, String loginId);

    /**
     * 根据tokenId获取token
     *
     * @param tokenId
     * @return
     */
    @Query("select t from PortalToken t where t.id = ?1")
    PortalToken queryTokenByTokenId(String tokenId);

    /**
     * 根据token名和token名的hash值获取token
     *
     * @param name
     * @param tokenCrc
     * @return
     */
    @Query("select t from PortalToken t where t.name = ?1 and t.tokenCrc=?2")
    PortalToken queryTokenByNameAndTokenCrc(String name, long tokenCrc);

}
