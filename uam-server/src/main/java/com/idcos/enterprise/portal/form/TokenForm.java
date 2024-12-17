

package com.idcos.enterprise.portal.form;

// auto generated imports

import com.idcos.cloud.biz.common.BaseForm;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import org.hibernate.validator.constraints.NotBlank;

/**
 * 表单对象TokenForm
 * <p>第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖</p>
 *
 * @author yanlv
 * @version TokenForm.java, v 1.1 2015-06-09 10:00:37 yanlv Exp $
 */
@ApiModel("生成Token")
public class TokenForm extends BaseForm {

    //========== properties ==========
    /**
     * loginId
     */
    @ApiModelProperty(name = "loginId", value = "员工登录名,eg:admin", notes = "员工的登录名，不能重复", required = true)
    @NotBlank(message = "登录id")
    private String loginId;

    /**
     * password
     */
    @ApiModelProperty(name = "password", value = "密码", notes = "登录密码", required = true)
    @NotBlank(message = "密码")
    private String password;

    /**
     * 租户id
     */
    @ApiModelProperty(name = "tenantId", value = "租户id,eg:default", notes = "租户id", required = true)
    @NotBlank(message = "tenant")
    private String tenantId;

    /**
     * token有效时间长度
     */
    @ApiModelProperty(name = "time", value = "token过期时间,eg:'2018-10-01 08:00'表示该token会在2018年10月1日上午8点过期；这里不填写则会生成系统的永久token", notes = "token过期时间")
    private String time;

    //========== getters and setters ==========

    public String getLoginId() {
        return loginId;
    }

    public void setLoginId(String loginId) {
        this.loginId = loginId;
    }

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public String getTenantId() {
        return tenantId;
    }

    public void setTenantId(String tenantId) {
        this.tenantId = tenantId;
    }

    public String getTime() {
        return time;
    }

    public void setTime(String time) {
        this.time = time;
    }
}