

package com.idcos.enterprise.portal.form;

// auto generated imports

import com.idcos.cloud.biz.common.BaseForm;
import com.idcos.cloud.biz.common.check.anonations.DateFormat;
import com.idcos.cloud.biz.common.check.anonations.In;
import org.hibernate.validator.constraints.NotBlank;


/**
 * 表单对象PortalCommandShellApiReportForm
 * <p>form由代码自动生成框架自动生成，不可进行编辑</p>
 *
 * @author
 * @version PortalCommandShellApiReportForm.java, v 1.1 2015-10-12 17:04:52  Exp $
 */

public class PortalCommandShellApiReportForm extends BaseForm {

    //========== properties ==========
    /**
     * 主机名
     */
    @NotBlank(message = "主机名不能为空")
    private String host;

    /**
     * 执行用户
     */
    @NotBlank(message = "执行用户不能为空")
    private String user;

    /**
     * 执行的具体命令
     */
    @NotBlank(message = "执行的具体命令不能为空")
    private String command;

    /**
     * 命令执行开始时间
     */
    @NotBlank(message = "命令执行开始时间不能为空")
    @DateFormat(format = "yyyy-MM-dd hh:mm:ss", message = "命令执行开始时间必须为'yyyy-MM-dd hh:mm:ss'格式")
    private String startAt;

    /**
     * 命令执行结束时间
     */
    @NotBlank(message = "命令执行结束时间不能为空")
    @DateFormat(format = "yyyy-MM-dd hh:mm:ss", message = "命令执行结束时间必须为'yyyy-MM-dd hh:mm:ss'格式")
    private String endAt;

    /**
     * 执行状态
     */
    @NotBlank(message = "执行状态不能为空")
    @In(value = {"success", "fail"}, message = "执行状态的值必须在{'success', 'fail'}之中")
    private String status;


    //========== getters and setters ==========

    public String getHost() {
        return host;
    }

    public void setHost(String host) {
        this.host = host;
    }


    public String getUser() {
        return user;
    }

    public void setUser(String user) {
        this.user = user;
    }


    public String getCommand() {
        return command;
    }

    public void setCommand(String command) {
        this.command = command;
    }


    public String getStartAt() {
        return startAt;
    }

    public void setStartAt(String startAt) {
        this.startAt = startAt;
    }


    public String getEndAt() {
        return endAt;
    }

    public void setEndAt(String endAt) {
        this.endAt = endAt;
    }


    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }


}