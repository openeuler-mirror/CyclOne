/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.form;

import java.util.List;

/**
 * @author Dana
 * @version UserListForm.java, v1 2017/12/7 下午1:55 Dana Exp $$
 */
public class UserListForm {
    /**
     * 用户ID.
     */
    private List<String> loginIdList;

    public List<String> getLoginIdList() {
        return loginIdList;
    }

    public void setLoginIdList(List<String> loginIdList) {
        this.loginIdList = loginIdList;
    }
}