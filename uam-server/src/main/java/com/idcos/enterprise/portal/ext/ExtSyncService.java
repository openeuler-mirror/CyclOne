package com.idcos.enterprise.portal.ext;

/**
 * @Author: Dai
 * @Date: 2018/11/8 10:26 PM
 * @Description: 外部系统同步用户和部门必须继承这个接口
 */
public interface ExtSyncService {
    /**
     * 同步部门和人员
     */
    void syncDeptAndUser();
}
