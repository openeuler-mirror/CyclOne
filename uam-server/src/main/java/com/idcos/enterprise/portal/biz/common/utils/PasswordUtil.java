package com.idcos.enterprise.portal.biz.common.utils;

import com.idcos.enterprise.portal.biz.common.CommonBizException;
import com.idcos.enterprise.portal.dal.enums.ArrayNumEnum;

import javax.crypto.Cipher;
import javax.crypto.SecretKey;
import javax.crypto.SecretKeyFactory;
import javax.crypto.spec.PBEKeySpec;
import javax.crypto.spec.PBEParameterSpec;
import java.security.Key;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Random;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * @author
 * @version PortalFilterUtil.java, v1 2017/11/18 下午6:22  Exp $$
 */
public class PasswordUtil {

    /**
     * JAVA6支持以下任意一种算法 PBEWITHMD5ANDDES PBEWITHMD5ANDTRIPLEDES
     * PBEWITHSHAANDDESEDE PBEWITHSHA1ANDRC2_40 PBKDF2WITHHMACSHA1
     * */

    /**
     * 定义使用的算法为:PBEWITHMD5andDES算法
     */
    public static final String ALGORITHM = "PBEWithMD5AndDES";

    /**
     * 定义迭代次数为1000次
     */
    private static final int ITERATIONCOUNT = 1000;

    /**
     * 获取加密算法中使用的盐值,解密中使用的盐值必须与加密中使用的相同才能完成操作. 盐长度必须为8字节
     *
     * @return byte[] 盐值
     */
    public static byte[] getSalt() throws Exception {
        // 实例化安全随机数
        //        SecureRandom random = new SecureRandom();
        // 产出盐
        return randomSalt(8);
    }

    /**
     * SecureRandom有坑，这里自定义一个假的盐
     *
     * @param length
     * @return
     */
    private static byte[] randomSalt(int length) {
        byte[] letter = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789".getBytes();
        int letterLength = letter.length;

        Random random = new Random();
        random.setSeed(System.currentTimeMillis());

        byte[] randomValue = new byte[length];
        for (int i = 0; i < length; i++) {
            int index = random.nextInt(length);
            randomValue[i] = letter[index];
        }
        return randomValue;
    }

    /**
     * 根据PBE密码生成一把密钥
     *
     * @param password 生成密钥时所使用的密码
     * @return Key PBE算法密钥
     */
    private static Key getPBEKey(String password) {
        // 实例化使用的算法
        SecretKeyFactory keyFactory;
        SecretKey secretKey = null;
        try {
            keyFactory = SecretKeyFactory.getInstance(ALGORITHM);
            // 设置PBE密钥参数
            PBEKeySpec keySpec = new PBEKeySpec(password.toCharArray());
            // 生成密钥
            secretKey = keyFactory.generateSecret(keySpec);
        } catch (Exception e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        }

        return secretKey;
    }

    /**
     * 加密明文字符串
     *
     * @param plaintext 待加密的明文字符串
     * @param password  生成密钥时所使用的密码
     * @param salt      盐值
     * @return 加密后的密文字符串
     * @throws Exception
     */
    public static String encrypt(String plaintext, String password, byte[] salt) {

        Key key = getPBEKey(password);
        byte[] encipheredData = null;
        PBEParameterSpec parameterSpec = new PBEParameterSpec(salt, ITERATIONCOUNT);
        try {
            Cipher cipher = Cipher.getInstance(ALGORITHM);

            cipher.init(Cipher.ENCRYPT_MODE, key, parameterSpec);

            encipheredData = cipher.doFinal(plaintext.getBytes());
        } catch (Exception e) {
        }
        return bytesToHexString(encipheredData);
    }

    /**
     * 解密密文字符串
     *
     * @param ciphertext 待解密的密文字符串
     * @param password   生成密钥时所使用的密码(如需解密,该参数需要与加密时使用的一致)
     * @param salt       盐值(如需解密,该参数需要与加密时使用的一致)
     * @return 解密后的明文字符串
     * @throws Exception
     */
    public static String decrypt(String ciphertext, String password, byte[] salt) {

        Key key = getPBEKey(password);
        byte[] passDec = null;
        PBEParameterSpec parameterSpec = new PBEParameterSpec(salt, ITERATIONCOUNT);
        try {
            Cipher cipher = Cipher.getInstance(ALGORITHM);

            cipher.init(Cipher.DECRYPT_MODE, key, parameterSpec);

            passDec = cipher.doFinal(hexStringToBytes(ciphertext));
        } catch (Exception e) {
            // TODO: handle exception
            e.printStackTrace();
        }
        return new String(passDec);
    }

    /**
     * 将字节数组转换为十六进制字符串
     *
     * @param src 字节数组
     * @return
     */
    public static String bytesToHexString(byte[] src) {
        StringBuilder stringBuilder = new StringBuilder("");
        if (src == null || src.length <= 0) {
            return null;
        }
        for (int i = 0; i < src.length; i++) {
            int v = src[i] & 0xFF;
            String hv = Integer.toHexString(v);
            if (hv.length() < 2) {
                stringBuilder.append(0);
            }
            stringBuilder.append(hv);
        }
        return stringBuilder.toString();
    }

    /**
     * 将十六进制字符串转换为字节数组
     *
     * @param hexString 十六进制字符串
     * @return
     */
    public static byte[] hexStringToBytes(String hexString) {
        if (hexString == null || "".equals(hexString)) {
            return null;
        }
        hexString = hexString.toUpperCase();
        int length = hexString.length() / 2;
        char[] hexChars = hexString.toCharArray();
        byte[] d = new byte[length];
        for (int i = 0; i < length; i++) {
            int pos = i * 2;
            d[i] = (byte) (charToByte(hexChars[pos]) << 4 | charToByte(hexChars[pos + 1]));
        }
        return d;
    }

    private static byte charToByte(char c) {
        return (byte) "0123456789ABCDEF".indexOf(c);
    }

    /**
     * public static void main(String[] args) {
     * <p>
     * System.out.print(-1 % 2 == 0);
     * String str = "r00t!@#";
     * String password = "";
     * <p>
     * System.out.println("明文:" + str);
     * System.out.println("密码:" + password);
     * <p>
     * try {
     * byte[] salt = PasswordUtil.getSalt();
     * String ciphertext = PasswordUtil.encrypt(str, password, salt);
     * System.out.println("密文:" + ciphertext);
     * String plaintext = PasswordUtil.decrypt(ciphertext, password, salt);
     * System.out.println("1明文:" + plaintext);
     * String tmp = new String(salt,"ISO-8859-1");
     * System.out.println(tmp);
     * <p>
     * salt = tmp.getBytes("ISO-8859-1");
     * plaintext = PasswordUtil.decrypt(ciphertext, password, salt);
     * System.out.println("2明文:" + plaintext);
     * } catch (Exception e) {
     * e.printStackTrace();
     * }
     * }
     */
    private static final Pattern PATTERN1 = Pattern.compile("[a-z]+");

    private static final Pattern PATTERN2 = Pattern.compile("[A-Z]+");

    private static final Pattern PATTERN3 = Pattern.compile("[\\d]+");

    /**
     * ~!@#$%^&*()_-+=|\{}[]"';:/?.>,<
     */
    private static final Pattern PATTERN4 = Pattern
            .compile("[~!@#\\$%\\^&\\*\\(\\)_\\-\\+=\\|\\\\\\{\\}\\[\\]\\\"';:\\/\\?\\.>,<]+");

    private static final Integer PASSWORD_MAX_LENGTH = 16;

    private static final Integer PASSWORD_MIN_LENGTH = 8;

    public static final void checkPassword(String password) {
        // 密码强度校验
        // 要求:用户设定密码时需包含大小写字母，数字，特殊字符中的3种，长度为8-16个字符
        if (password.length() < PASSWORD_MIN_LENGTH || password.length() > PASSWORD_MAX_LENGTH) {
            throw new CommonBizException("密码长度必须为8-16位，且包含大小写字母、数字、英文特殊字符!");
        }
        Matcher m1 = PATTERN1.matcher(password);
        if (!m1.find()) {
            throw new CommonBizException("密码长度必须为8-16位，且包含大小写字母、数字、英文特殊字符!");
        }
        Matcher m2 = PATTERN2.matcher(password);
        if (!m2.find()) {
            throw new CommonBizException("密码长度必须为8-16位，且包含大小写字母、数字、英文特殊字符!");
        }
        Matcher m3 = PATTERN3.matcher(password);
        if (!m3.find()) {
            throw new CommonBizException("密码长度必须为8-16位，且包含大小写字母、数字、英文特殊字符!");
        }
        Matcher m4 = PATTERN4.matcher(password);
        if (!m4.find()) {
            throw new CommonBizException("密码长度必须为8-16位，且包含大小写字母、数字、英文特殊字符!");
        }
    }

    /**
     * 重置密码生成随机字符串，必须包含大小写字母、数字和特殊字符
     *
     * @param length
     * @return
     */
    public static String getStringRandom(int length) {
        if (length < ArrayNumEnum.EIGHT.getCode()) {
            throw new CommonBizException("密码长度至少为8!");
        }
        String num = "23456789";
        String letter = "abcdefghjkmnpqrstuvwxyz";
        String specialChar = "~!@#$%^&*";
        Random random = new Random();
        StringBuffer sb = new StringBuffer();
        sb.append(specialChar.charAt(random.nextInt(specialChar.length())));
        for (int i = 0; i < ArrayNumEnum.TWO.getCode(); i++) {
            sb.append(num.charAt(random.nextInt(num.length())));
        }
        for (int i = 0; i < ArrayNumEnum.TWO.getCode(); i++) {
            sb.append(letter.toUpperCase().charAt(random.nextInt(letter.toUpperCase().length())));
        }
        for (int i = 0; i < length - ArrayNumEnum.FIVE.getCode(); i++) {
            sb.append(letter.charAt(random.nextInt(letter.length())));
        }
        char[] c = sb.toString().toCharArray();
        List<Character> lst = new ArrayList<>();
        for (int i = 0; i < c.length; i++) {
            lst.add(c[i]);
        }
        Collections.shuffle(lst);
        String resultStr = "";
        for (int i = 0; i < lst.size(); i++) {
            resultStr += lst.get(i);
        }
        return resultStr;
    }

    /**
     * 对传入的密码加密处理
     *
     * @param id
     * @param password
     * @param salt
     * @return
     */
    public static String encryptPassword(String id, String password, String salt) {
        String encriptPW;
        try {
            byte[] saltBytes = Base64Util.decode(salt);
            encriptPW = encrypt(password, id, saltBytes);
        } catch (Exception e) {
            throw new RuntimeException("密码加密错误:", e);
        }
        return encriptPW;
    }

}
