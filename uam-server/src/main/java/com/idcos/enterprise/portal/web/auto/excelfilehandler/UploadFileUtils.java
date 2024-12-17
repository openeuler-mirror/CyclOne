/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.web.auto.excelfilehandler;

import org.springframework.web.multipart.MultipartFile;
import org.springframework.web.multipart.MultipartHttpServletRequest;
import org.springframework.web.multipart.commons.CommonsMultipartResolver;

import javax.servlet.http.HttpServletRequest;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

/**
 * 文件上传工具类
 *
 * @author Dana
 * @version UploadFileUtils.java, v1 2017/12/10 下午6:04 Dana Exp $$
 */
public class UploadFileUtils {

    /**
     *
     */
    private UploadFileUtils() {
    }

    /**
     * 根据上传请求得到其中上传的文件。
     *
     * @param request 上传的请求对象。
     * @return 返回上传文件列表。
     */
    public static MultipartFile[] getUploadFile(HttpServletRequest request) {
        //创建一个通用的多部分解析器
        CommonsMultipartResolver multipartResolver = new CommonsMultipartResolver(
                request.getSession().getServletContext());

        List<MultipartFile> fileList = new ArrayList<MultipartFile>();
        //判断 request 是否有文件上传,即多部分请求
        if (multipartResolver.isMultipart(request)) {
            //转换成多部分request
            MultipartHttpServletRequest multiRequest = (MultipartHttpServletRequest) request;
            //取得request中的所有文件名
            Iterator<String> iter = multiRequest.getFileNames();
            while (iter.hasNext()) {

                //取得上传文件
                MultipartFile file = multiRequest.getFile(iter.next());
                if (file != null) {
                    fileList.add(file);
                }

            }
        }
        return fileList.toArray(new MultipartFile[0]);
    }

    /**
     * 返回上传的文件中第一个文件。
     *
     * @param request
     * @return
     */
    public static MultipartFile getFirstUploadFile(HttpServletRequest request) {
        MultipartFile[] files = getUploadFile(request);

        if (files.length > 0) {
            return files[0];
        }
        return null;
    }
}