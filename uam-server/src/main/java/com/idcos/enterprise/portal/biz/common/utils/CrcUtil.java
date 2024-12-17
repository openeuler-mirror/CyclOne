package com.idcos.enterprise.portal.biz.common.utils;

import java.util.zip.CRC32;

/**
 * @Author: Dai
 * @Date: 2018/11/15 12:03 PM
 * @Description: crc加密工具类
 */
public class CrcUtil {
    public static final long crc(String value) {
        CRC32 crc32 = new CRC32();
        crc32.update(value.getBytes());
        return crc32.getValue();
    }
}
