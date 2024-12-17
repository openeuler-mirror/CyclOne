package com.idcos.enterprise.portal.biz.common;

/**
 * 抽象配置项管理系统的返回信息对象。
 *
 * @author pengganyu
 * @version $Id: PortalResponse.java, v 0.1 2016年5月10日 下午4:45:34 pengganyu Exp $
 */
public class PortalResponse {

    /**
     * 返回调用操作是否有错误,代表状态码，成功success，失败error
     * <p>
     * 此处理的错误主要是业务层面的错误信息。
     * </p>
     */
    private String status;

    /**
     * 返回信息，不管操作是否成功，都可以返回相应的信息。
     */
    private String message;

    /**
     * 返回值，不管操作是否成功都可以返回相应的数据。
     */
    private Object content;

    public PortalResponse() {
        content = "";
    }

    /**
     * Getter method for property <tt>error</tt>.
     *
     * @return property value of error
     */
    public String getStatus() {
        return status;
    }

    /**
     * Setter method for property <tt>error</tt>.
     *
     * @param error value to be assigned to property error
     */
    public void setStatus(String status) {
        this.status = status;
    }

    /**
     * Getter method for property <tt>message</tt>.
     *
     * @return property value of message
     */
    public String getMessage() {
        return message;
    }

    /**
     * Setter method for property <tt>message</tt>.
     *
     * @param message value to be assigned to property message
     */
    public void setMessage(String message) {
        this.message = message;
    }

    /**
     * Getter method for property <tt>data</tt>.
     *
     * @return property value of data
     */
    public Object getContent() {
        return content;
    }

    /**
     * Setter method for property <tt>data</tt>.
     *
     * @param data value to be assigned to property data
     */
    public void setContent(Object content) {
        this.content = content;
    }

}
