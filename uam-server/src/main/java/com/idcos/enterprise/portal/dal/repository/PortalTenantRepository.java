

package com.idcos.enterprise.portal.dal.repository;

// auto generated imports

import com.idcos.cloud.core.dal.common.jpa.BaseRepository;
import com.idcos.enterprise.portal.dal.entity.PortalTenant;
import org.springframework.data.jpa.repository.Query;

import java.util.List;

/**
 * 自动生成PortalRoleRepository
 * <p>
 * 数据库配置文件自动生成,此文件属于自动生成的,具体可以参考codegen工程
 * Generated by <tt>jap-codgen</tt> on 2017-09-26 09:22:12.
 *
 * @author GuanBin
 * @version PortalRoleRepository.java, v 1.1 2017-09-26 09:22:12 yanlv Exp $
 */
public interface PortalTenantRepository extends BaseRepository<PortalTenant, String> {
    /**
     * 根据名称查找
     *
     * @param name
     * @return
     */
    @Query("select u from PortalTenant u where u.name = ?1 and u.isActive = 'Y'")
    PortalTenant findByName(String name);

    /**
     * 查询所有租户
     *
     * @return
     */
    @Query(nativeQuery = true, value = "select * from PORTAL_TENANT u where u.IS_ACTIVE = 'Y' order by convert(u.DISPLAY_NAME using gbk)")
    List<PortalTenant> getAllPortalTenant();

    /**
     * 根据displayName查找
     *
     * @param displayName
     * @return
     */
    @Query("select u from PortalTenant u where u.displayName = ?1 and u.isActive = 'Y'")
    PortalTenant findByDisplayName(String displayName);

    /**
     * NameNotId
     *
     * @param name
     * @param id
     * @return
     */
    @Query("select u from PortalTenant u where u.name = ?1 and u.id <> ?2 and u.isActive = 'Y'")
    PortalTenant findByNameNotId(String name, String id);

    /**
     * DisplayNameNotId
     *
     * @param displayName
     * @param id
     * @return
     */
    @Query("select u from PortalTenant u where u.displayName = ?1 and u.id <> ?2 and u.isActive = 'Y'")
    PortalTenant findByDisplayNameNotId(String displayName, String id);

    /**
     * 根据id查找
     *
     * @param id
     * @return
     */
    @Query("select u from PortalTenant u where u.id = ?1 and u.isActive = 'Y'")
    PortalTenant findById(String id);
}