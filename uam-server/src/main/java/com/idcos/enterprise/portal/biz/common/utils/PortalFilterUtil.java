package com.idcos.enterprise.portal.biz.common.utils;

import java.util.List;

import com.google.common.collect.Lists;
import com.idcos.cloud.core.common.util.ListUtil;
import com.idcos.enterprise.portal.dal.entity.*;

/**
 * @author
 * @version PortalFilterUtil.java, v1 2017/11/18 下午6:22 Exp $$
 */
public class PortalFilterUtil {

    /**
     * 过滤权限 不同的角色可能分配了相同的权限信息
     *
     * @param list
     * @return
     */
    public static List<PortalPermission> filterPermission(List<PortalPermission> list) {
        List<PortalPermission> newList = Lists.newArrayList();
        for (int i = list.size() - 1; i >= 0; i--) {
            if (ListUtil.findOne(newList, "authResId", list.get(i).getAuthResId()) == null) {
                newList.add(list.get(i));
            }

        }
        return newList;
    }

    /**
     * 过滤用户 不同的用户组可能分配相同的用户
     *
     * @param list
     * @return
     */
    public static List<PortalUser> filterUser(List<PortalUser> list) {
        List<PortalUser> newList = Lists.newArrayList();

        for (int i = list.size() - 1; i >= 0; i--) {
            if (ListUtil.findOne(newList, "id", list.get(i).getId()) == null) {
                newList.add(list.get(i));
            }
        }
        return newList;
    }

    /**
     * 过滤角色 不同的用户可能属于不同的角色
     *
     * @param list
     * @return
     */
    public static List<PortalRole> filterRole(List<PortalRole> list) {
        List<PortalRole> newList = Lists.newArrayList();

        for (int i = list.size() - 1; i >= 0; i--) {
            if (ListUtil.findOne(newList, "id", list.get(i).getId()) == null) {
                newList.add(list.get(i));
            }
        }

        return newList;
    }

    /**
     * 过滤用户组
     *
     * @param list
     * @return
     */
    public static List<PortalUserGroup> filterGroup(List<PortalUserGroup> list) {
        List<PortalUserGroup> newList = Lists.newArrayList();

        for (int i = list.size() - 1; i >= 0; i--) {
            if (ListUtil.findOne(newList, "id", list.get(i).getId()) == null) {
                newList.add(list.get(i));
            }
        }

        return newList;
    }
}
