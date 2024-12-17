package com.idcos.enterprise.portal.dal.repository;

import java.util.List;

import javax.persistence.EntityManager;
import javax.persistence.PersistenceContext;
import javax.persistence.TypedQuery;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Repository;
import org.springframework.transaction.annotation.Transactional;

import com.idcos.enterprise.portal.dal.entity.PortalRole;
import com.idcos.enterprise.portal.dal.entity.PortalUser;
import com.idcos.enterprise.portal.dal.entity.PortalUserGroup;

/**
 * 权限资源管理关联表分页查询实现
 *
 * @author pengganyu
 * @version $Id: PortalRelQueryRepository.java, v 0.1 2016年6月2日 上午10:23:01 pengganyu Exp $
 */
@Repository
@Transactional(readOnly = true, rollbackFor = {Exception.class})
public class PortalRelQueryRepository {

    private static final Logger Logger = LoggerFactory.getLogger(PortalRelQueryRepository.class);

    @PersistenceContext
    private EntityManager em;

    /**
     * 查询用户组关联的用户信息
     *
     * @param groupId  用户组id
     * @param cnd      用户查询条件
     * @param pageNo   页号
     * @param pageSize 页大小
     * @return
     */
    public List<PortalUser> findUserPageByGroupIdAndCnd(String groupId, String cnd, int pageNo,
                                                        int pageSize) {

        String hql = "select u from PortalUser u,PortalGroupUserRel r where u.id=r.userId and r.groupId =?1 and u.isActive='Y' and (u.name like ?2 or u.loginId like ?2)";

        Logger.info("查询用户组关联的用户信息：" + hql);
        TypedQuery<PortalUser> query = em.createQuery(hql, PortalUser.class);

        query.setParameter(1, groupId);
        query.setParameter(2, cnd);
        query.setFirstResult((pageNo - 1) * pageSize);
        query.setMaxResults(pageSize);

        return query.getResultList();
    }

    /**
     * 查询角色所关联的用户信息
     *
     * @param roleId   用户组id
     * @param cnd      用户查询条件
     * @param pageNo   页号
     * @param pageSize 页大小
     * @return
     */
    public List<PortalUser> findUserPageByRoleIdAndCnd(String roleId, String cnd, int pageNo,
                                                       int pageSize) {

        String hql = "select u from PortalUser u,PortalGroupUserRel r,PortalGroupRoleRel g where u.id=r.userId and g.groupId = r.groupId and g.roleId =?1 and u.isActive='Y' and (u.name like ?2 or u.loginId like ?2)";

        Logger.info("查询角色所关联的用户信息：" + hql);
        TypedQuery<PortalUser> query = em.createQuery(hql, PortalUser.class);

        query.setParameter(1, roleId);
        query.setParameter(2, cnd);
        query.setFirstResult((pageNo - 1) * pageSize);
        query.setMaxResults(pageSize);

        return query.getResultList();
    }

    /**
     * 查询用户所关联的用户组信息
     *
     * @param userId   用户组id
     * @param cnd      用户组查询条件
     * @param pageNo   页号
     * @param pageSize 页大小
     * @return
     */
    public List<PortalUserGroup> findGroupPageByUserIdAndCnd(String userId, String cnd, int pageNo,
                                                             int pageSize) {

        String hql = "select u from PortalUserGroup u,PortalGroupUserRel r where u.id=r.groupId and r.userId =?1 and u.isActive='Y' and (u.name like ?2 or u.remark like ?2)";

        Logger.info("查询用户所关联的用户组信息：" + hql);
        TypedQuery<PortalUserGroup> query = em.createQuery(hql, PortalUserGroup.class);

        query.setParameter(1, userId);
        query.setParameter(2, cnd);
        query.setFirstResult((pageNo - 1) * pageSize);
        query.setMaxResults(pageSize);

        return query.getResultList();
    }

    /**
     * 查询角色所关联的用户组信息
     *
     * @param roleId   用户组id
     * @param cnd      用户组查询条件
     * @param pageNo   页号
     * @param pageSize 页大小
     * @return
     */
    public List<PortalUserGroup> findGroupPageByRoleIdAndCnd(String roleId, String cnd, int pageNo,
                                                             int pageSize) {

        String hql = "select u from PortalUserGroup u,PortalGroupRoleRel r where u.id=r.groupId and r.roleId =?1 and u.isActive='Y' and (u.name like ?2 or u.remark like ?2)";

        Logger.info("查询角色所关联的用户组信息：" + hql);
        TypedQuery<PortalUserGroup> query = em.createQuery(hql, PortalUserGroup.class);

        query.setParameter(1, roleId);
        query.setParameter(2, cnd);
        query.setFirstResult((pageNo - 1) * pageSize);
        query.setMaxResults(pageSize);

        return query.getResultList();
    }

    /**
     * 查询用户组所关联的角色信息
     *
     * @param groupId  用户组id
     * @param cnd      用户组查询条件
     * @param pageNo   页号
     * @param pageSize 页大小
     * @return
     */
    public List<PortalRole> findRolePageByGroupId(String groupId, String cnd, int pageNo, int pageSize) {

        String hql = "select distinct u from PortalRole u,PortalGroupRoleRel r where u.id=r.roleId and r.groupId =?1 and u.isActive='Y' and (u.name like ?2 or u.code like ?2 or u.remark like ?2)";

        Logger.info("查询用户组所关联的角色信息：" + hql);
        TypedQuery<PortalRole> query = em.createQuery(hql, PortalRole.class);

        query.setParameter(1, groupId);
        query.setParameter(2, cnd);
        query.setFirstResult((pageNo - 1) * pageSize);
        query.setMaxResults(pageSize);

        return query.getResultList();
    }

    /**
     * 查询用户所关联的角色信息
     *
     * @param userId   用户组id
     * @param cnd      用户组查询条件
     * @param pageNo   页号
     * @param pageSize 页大小
     * @return
     */
    public List<PortalRole> findRolePageByUserId(String userId, String cnd, int pageNo, int pageSize) {

        String hql = "select u from PortalRole u,PortalGroupUserRel r,PortalGroupRoleRel g where u.id=g.roleId and g.groupId = r.groupId and r.userId =?1 and u.isActive='Y' and (u.name like ?2 or u.code like ?2 or u.remark like ?2)";

        Logger.info("查询用户所关联的角色信息：" + hql);
        TypedQuery<PortalRole> query = em.createQuery(hql, PortalRole.class);

        query.setParameter(1, userId);
        query.setParameter(2, cnd);
        query.setFirstResult((pageNo - 1) * pageSize);
        query.setMaxResults(pageSize);

        return query.getResultList();
    }

}