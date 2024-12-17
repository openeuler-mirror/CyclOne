package com.idcos.enterprise.sso.service;

import java.util.Map;

/**
 * @author xizhao
 */
public class LoginResult {
    /**
     * true: login success， false： login failed
     */
    private boolean status;
    private String message;
    private Map<String, Object> otherData;

    public boolean isStatus() {
        return status;
    }

    public void setStatus(boolean status) {
        this.status = status;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public Map<String, Object> getOtherData() {
        return otherData;
    }

    public void setOtherData(Map<String, Object> otherData) {
        this.otherData = otherData;
    }

}
