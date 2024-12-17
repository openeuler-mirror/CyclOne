/**
 * 系统字典工具类信息
 *
 * @author yanlv
 * @version $Id: CfgSysDictUtils.java, v 0.1 2015年6月17日 上午11:04:54 yanlv Exp $
 */
package com.idcos.enterprise.portal.biz.common.utils;

import com.idcos.common.component.cache.CacheProvider;
import com.idcos.common.component.cache.CommonBizCache;
import org.springframework.beans.factory.InitializingBean;

import java.util.Collection;

/**
 * 配置缓存工具类
 * @author yanlv
 * @version $Id: CfgCacheUtil.java, v 0.1 2015年6月25日 下午6:30:51 yanlv Exp $
 */
public class CfgCacheUtil implements InitializingBean {


    @Override
    public void afterPropertiesSet() throws Exception {

    }

    /** 缓存提供者 */
    private static CacheProvider cacheProvider;

    /**
     * 默认的构造函数
     * @param cacheProvider
     */
    public CfgCacheUtil(CacheProvider cacheProvider) {
        CfgCacheUtil.cacheProvider = cacheProvider;
    }

    /**
     * 根据key以及缓存class 获取缓存中的数据类型
     * @param cacheEnum 缓存类型
     * @param key 缓存中的key的值
     * @param returnField 需要的返回字段值
     * @return
     */

    public static final Object find(Class<?> cacheClass, String key) {

        CommonBizCache<?> cache = cacheProvider.getCache(cacheClass);

        Object object = cache.get(key);
        return object;

    }

    /**
     * 获取缓存中所有的对象数据
     * @param cacheClass
     * @return
     */
    public static final Collection<?> findAll(Class<?> cacheClass) {
        CommonBizCache<?> cache = cacheProvider.getCache(cacheClass);
        return cache.getAll();

    }

}
