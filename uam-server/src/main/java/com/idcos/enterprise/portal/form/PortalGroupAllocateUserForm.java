

package com.idcos.enterprise.portal.form;

import org.hibernate.validator.constraints.NotBlank;

// auto generated imports
import com.idcos.cloud.biz.common.BaseForm;

/**
 * 表单对象PortalUserUpdateForm
 * <p>form由代码自动生成框架自动生成，不可进行编辑</p>
 *
 * @author
 * @version PortalUserUpdateForm.java, v 1.1 2015-11-07 16:44:40  Exp $
 */

public class PortalGroupAllocateUserForm extends BaseForm {

    @NotBlank(message = "用户id不能为空")
    private String id;

    private String nameCN;

    private String nameEN;

    private String deptId;

    private String deptName;

    private String title;

    private String mobile;

    private String wexin;

    private String fax;

    private String email;

    private String tel;

    private String remark;

    /**
     * 操作类型  I:新增、D:删除
     */
    @NotBlank(message = "操作类型不能为空")
    private String operType;

    //========== getters and setters ==========

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getDeptId() {
        return deptId;
    }

    public void setDeptId(String deptId) {
        this.deptId = deptId;
    }

    public String getDeptName() {
        return deptName;
    }

    public void setDeptName(String deptName) {
        this.deptName = deptName;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getMobile() {
        return mobile;
    }

    public void setMobile(String mobile) {
        this.mobile = mobile;
    }

    public String getNameCN() {
        return nameCN;
    }

    public void setNameCN(String nameCN) {
        this.nameCN = nameCN;
    }

    public String getNameEN() {
        return nameEN;
    }

    public void setNameEN(String nameEN) {
        this.nameEN = nameEN;
    }

    public String getWexin() {
        return wexin;
    }

    public void setWexin(String wexin) {
        this.wexin = wexin;
    }

    public String getFax() {
        return fax;
    }

    public void setFax(String fax) {
        this.fax = fax;
    }

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }

    public String getTel() {
        return tel;
    }

    public void setTel(String tel) {
        this.tel = tel;
    }

    public String getRemark() {
        return remark;
    }

    public void setRemark(String remark) {
        this.remark = remark;
    }

    public String getOperType() {
        return operType;
    }

    public void setOperType(String operType) {
        this.operType = operType;
    }

}