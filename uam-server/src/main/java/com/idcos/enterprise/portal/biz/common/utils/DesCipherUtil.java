package com.idcos.enterprise.portal.biz.common.utils;

import org.bouncycastle.jce.provider.BouncyCastleProvider;

import javax.crypto.*;
import javax.crypto.spec.DESKeySpec;
import java.io.IOException;
import java.nio.charset.Charset;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.security.NoSuchProviderException;
import java.security.Security;
import java.security.spec.InvalidKeySpecException;

/**
 * @author Xizhao.Dai
 * @version DesCipherUtil.java, v1 2018/8/3 上午11:05 Xizhao.Dai Exp $$
 */
public class DesCipherUtil {
    private DesCipherUtil() {
        throw new AssertionError("No DesCipherUtil instances for you!");
    }

    static {
        // add BC provider
        Security.addProvider(new BouncyCastleProvider());
    }

    /**
     * 加密
     *
     * @param encryptText 需要加密的信息
     * @param key         加密密钥
     * @return 加密后Base64编码的字符串
     */
    public static String encrypt(String encryptText, String key) {

        if (encryptText == null || key == null) {
            throw new IllegalArgumentException("encryptText or key must not be null");
        }

        try {
            DESKeySpec desKeySpec = new DESKeySpec(key.getBytes());
            SecretKeyFactory secretKeyFactory = SecretKeyFactory.getInstance("DES");
            SecretKey secretKey = secretKeyFactory.generateSecret(desKeySpec);

            Cipher cipher = Cipher.getInstance("DES/ECB/PKCS7Padding", "BC");
            cipher.init(Cipher.ENCRYPT_MODE, secretKey);
            byte[] bytes = cipher.doFinal(encryptText.getBytes(Charset.forName("UTF-8")));
            return Base64Util.encode(bytes);

        } catch (NoSuchAlgorithmException | InvalidKeyException | InvalidKeySpecException | NoSuchPaddingException
                | BadPaddingException | NoSuchProviderException | IllegalBlockSizeException e) {
            throw new RuntimeException("encrypt failed", e);
        }

    }

    /**
     * 解密
     *
     * @param decryptText 需要解密的信息
     * @param key         解密密钥，经过Base64编码
     * @return 解密后的字符串
     */
    public static String decrypt(String decryptText, String key) {

        if (decryptText == null || key == null) {
            throw new IllegalArgumentException("decryptText or key must not be null");
        }

        try {
            DESKeySpec desKeySpec = new DESKeySpec(key.getBytes());
            SecretKeyFactory secretKeyFactory = SecretKeyFactory.getInstance("DES");
            SecretKey secretKey = secretKeyFactory.generateSecret(desKeySpec);

            Cipher cipher = Cipher.getInstance("DES/ECB/PKCS7Padding", "BC");
            cipher.init(Cipher.DECRYPT_MODE, secretKey);
            byte[] decode = Base64Util.decode(decryptText);
            byte[] bytes = cipher.doFinal(decode);
            return new String(bytes, Charset.forName("UTF-8"));

        } catch (NoSuchAlgorithmException | InvalidKeyException | InvalidKeySpecException | NoSuchPaddingException
                | BadPaddingException | NoSuchProviderException | IllegalBlockSizeException | IOException e) {
            throw new RuntimeException("decrypt failed", e);
        }
    }
}