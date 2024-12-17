package com.idcos.enterprise.portal.biz.common.utils;

/**
 * @Author: Dai
 * @Date: 2018/10/24 11:44 PM
 * @Description:
 */
public class SourceTypeUtil {
    /**
     * 数据库存储的sourceType转换成页面显示的sourceType
     *
     * @param sourceType
     * @return
     */
    public static String getSourceType(String sourceType) {
        switch (sourceType) {
            case "native":
                return "UAM";
            case "webank":
                return "webank";
            default:
                return "UAM";
        }
    }

    /**
     * 页面显示的sourceType转换成数据库存储的sourceType
     *
     * @param sourceType
     * @return
     */
    public static String setSourceType(String sourceType) {
        switch (sourceType) {
            case "UAM":
                return "native";
            case "LDAP":
            case "ldap":
                return "ldap";
            case "webank":
                return "webank";
            default:
                return "native";
        }
    }
}
