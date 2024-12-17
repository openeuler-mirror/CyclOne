/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.web.auto.excelfilehandler;

import com.google.common.cache.Cache;
import com.google.common.cache.CacheBuilder;
import com.idcos.enterprise.portal.vo.PortalUserImportVO;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.concurrent.TimeUnit;

/**
 * 为了生成excel导入预览，缓存导入时生成的List
 *
 * @author Dana
 * @version UserExcelCache.java, v1 2017/12/27 上午9:37 Dana Exp $$
 */
public class UserExcelManager {
    private static final Logger logger = LoggerFactory.getLogger(UserExcelManager.class);

    /**
     * 为提高性能，缓存userImportList，此缓存会在30分钟无操作后失效(key:excel名称，value:List<PortalUserImportVO>)。
     */

    private Cache<String, List<PortalUserImportVO>> userImportListCache;

    private static UserExcelManager instance;

    /**
     * @param timeout 缓存超时时间，单位：秒。
     */
    private void init(int timeout) {
        logger.info("----Account Timeout in minutes: {}", timeout);
        this.userImportListCache = CacheBuilder.newBuilder().expireAfterAccess(timeout, TimeUnit.MINUTES).build();
    }

    /**
     * 得到excel管理器实例。
     *
     * @return
     */
    public static UserExcelManager getInstance() {
        synchronized (UserExcelManager.class) {
            if (instance == null) {
                instance = new UserExcelManager();
                //缓存时间设置为30分钟
                instance.init(30);
            }
        }
        return instance;
    }

    /**
     * 删除缓存中的UserImportList
     *
     * @param fileName excel文件名
     */
    public void removeUserImportList(String fileName) {
        this.userImportListCache.invalidate(fileName);
    }

    /**
     * 向缓存中添加一个UserImportList
     *
     * @param fileName
     * @param userImportList
     */
    public void addUserImportList(String fileName, List<PortalUserImportVO> userImportList) {
        this.userImportListCache.put(fileName, userImportList);
    }

    /**
     * 得到已缓存的UserImportList，返回当前系统UserImportList。
     *
     * @param fileName excel文件名
     * @return 返回UserImportList，如果要查询的UserImportList没有在缓存中存在，则返回null.
     */
    public List<PortalUserImportVO> getUserImportList(final String fileName) {
        List<PortalUserImportVO> userImportList = null;
        try {
            userImportList = this.userImportListCache.getIfPresent(fileName);
        } catch (Exception e) {
            logger.error("Load loggined account from cache failed,cause by：" + e.getMessage(), e);
        }
        return userImportList;
    }
}