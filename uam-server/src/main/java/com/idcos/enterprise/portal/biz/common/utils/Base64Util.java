package com.idcos.enterprise.portal.biz.common.utils;

import sun.misc.BASE64Decoder;
import sun.misc.BASE64Encoder;

import java.io.IOException;

/**
 * base64编码解码工具类
 *
 * @author Dana
 * @version Base64Util.java, v1 2017/11/30 下午11:38 Dana Exp $$
 */
public class Base64Util {
    private static BASE64Decoder base64Decoder = new BASE64Decoder();

    private static BASE64Encoder base64Encoder = new BASE64Encoder();

    public static final String encode(byte[] b) {
        return base64Encoder.encode(b);
    }

    public static final byte[] decode(String str) throws IOException {
        return base64Decoder.decodeBuffer(str);
    }
}