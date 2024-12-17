package com.idcos.enterprise.portal.ext;

import org.springframework.boot.autoconfigure.condition.ConditionalOnExpression;
import org.springframework.stereotype.Service;

/**
 * @Author: Dai
 * @Date: 2018/11/20 9:06 PM
 * @Description:
 */
@Service
@ConditionalOnExpression("'${ext.user.sync.type}' != 'LDAP' and '${ext.user.sync.type}' != 'OA'")
public class NoExtSynvServiceImpl implements ExtSyncService {

    @Override
    public void syncDeptAndUser() {
        throw new RuntimeException("同步功能尚未开通");
    }
}
