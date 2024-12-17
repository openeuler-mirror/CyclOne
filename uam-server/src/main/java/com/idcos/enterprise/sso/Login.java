package com.idcos.enterprise.sso;

import javax.servlet.http.HttpServletResponse;

import org.apache.commons.lang3.StringUtils;
import org.hibernate.validator.internal.util.privilegedactions.GetResource;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.*;

import com.alibaba.fastjson.JSONObject;
import com.idcos.cloud.biz.common.permisson.OperateAction;
import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.cloud.web.common.BaseResultVO;
import com.idcos.cloud.web.common.JsonResultUtil;
import com.idcos.enterprise.portal.biz.common.utils.DesCipherUtil;
import com.idcos.enterprise.portal.dal.repository.PortalUserRepository;
import com.idcos.enterprise.portal.form.LoginForm;
import com.idcos.enterprise.portal.form.TokenForm;
import com.idcos.enterprise.portal.manager.auto.AuthManager;
import com.idcos.enterprise.portal.web.GlobalValue;
import com.idcos.enterprise.sso.service.LoginResult;
import com.idcos.enterprise.sso.service.impl.LoginServiceImpl;

import io.swagger.annotations.*;

/**
 * @author GuanBin
 * @version Login.java, v1 2017/11/1 下午8:06 GuanBin Exp $$
 */
@Controller
@OperateAction("F_AUTH_REST")
@Api(tags = "01.登录与token管理接口", description = "Login")
public class Login {
    private static final Logger logger = LoggerFactory.getLogger(Login.class);

    @Autowired
    private AuthManager authManager;

    @Autowired
    private LoginServiceImpl loginService;

    @Autowired
    private GlobalValue globalValue;

    @Autowired
    private PortalUserRepository portalUserRepository;

    /**
     * 登录接口
     *
     * @param form
     * @return
     */
    @RequestMapping(method = RequestMethod.POST, value = "sso/login")
    @ResponseBody
    @OperateAction("O_PORTAL_USER_GET")
    @ApiOperation(value = "登录", notes = "登录。")
    public JSONObject login(@RequestBody LoginForm form, HttpServletResponse response,
        @RequestParam(name = "customer", required = false) String customer) {
        String afterDecrypt = DesCipherUtil.decrypt(form.getPassword(), globalValue.getDecryptKey());
        form.setPassword(afterDecrypt);
        CommonResult<?> result = loginService.login(form);
        LoginResult loginResult = (LoginResult)result.getResultObject();
        String accessToken = null;
        if (loginResult != null) {
            accessToken = authManager.getAuth(form);
        }

        JSONObject json = new JSONObject();
        json.put("status", result.getResultCode());
        json.put("message", result.getResultMessage());
        json.put("content", accessToken);

        if (loginResult != null && !loginResult.isStatus()) {
            json.put("status", loginResult.isStatus());
            json.put("message", loginResult.getMessage());
        }

        loginService.afterLogin(response, customer);
        return json;
    }

    @ApiOperation(value = "是否已经登录判定", notes = "是否已经登录判定")
    @RequestMapping(value = "sso/info", method = RequestMethod.POST, produces = "application/json; charset=UTF-8")
    public @ResponseBody JSONObject userInfo(@CookieValue(value = "LtpaToken2", required = false) String token) {
        JSONObject json = new JSONObject();
        if (StringUtils.isBlank(token)) {
            json.put("status", "fail");
            return json;
        }
        return json;
    }

    /**
     * 登录页面
     *
     * @return
     */
    @RequestMapping(value = "/login", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    public String loginHtml(@RequestParam(name = "customer", required = false) String customer, Model model) {
        model.addAttribute("isMultiTenant", globalValue.getMultiTenant());
        model.addAttribute("customer", customer);
        model.addAttribute("decryptKey",
            globalValue.getDecryptKey() == null ? "com.idcos.enterprise.sso" : globalValue.getDecryptKey());
        if (StringUtils.isBlank(customer) || !getLoginTemplates(customer)) {
            return "default/login";
        } else {
            return customer + "/login";
        }
    }

    /**
     * 发放token
     *
     * @return BaseResultVO
     */
    @OperateAction("O_PORTAL_COMMON_ROLELIST")
    @ApiOperation(value = "生成token", notes = "生成token。")
    @RequestMapping(method = RequestMethod.POST, value = "sso/token")
    @ResponseBody
    public BaseResultVO grantToken(@ModelAttribute TokenForm form) {
        return JsonResultUtil.getResult(
            authManager.grantToken(form.getLoginId(), form.getPassword(), form.getTenantId(), form.getTime()));
    }

    /**
     * 分析token
     *
     * @return BaseResultVO
     */
    @OperateAction("O_PORTAL_COMMON_ROLELIST")
    @ApiOperation(value = "分析token", notes = "分析token。")
    @RequestMapping(method = RequestMethod.POST, value = "sso/token/parse")
    @ResponseBody
    public BaseResultVO parseToken(@RequestBody String token) {
        return JsonResultUtil.getResult(authManager.parseToken(token));
    }

    /**
     * admin用户给每个用户发放token
     *
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.GET, value = "sso/token/admin")
    @ResponseBody
    @OperateAction("O_PORTAL_COMMON_ROLELIST")
    @ApiOperation(value = "admin给每个用户发放永久token", notes = "admin给每个用户发放永久token。")
    public BaseResultVO grantTokenByAdmin(@RequestParam String loginId, @RequestParam String tenantId) {
        return JsonResultUtil.getResult(authManager.grantTokenByAdmin(loginId, tenantId));
    }

    /**
     * 发放永久token（主要为命令行提供）
     *
     * @return BaseResultVO
     */
    @OperateAction("O_PORTAL_COMMON_ROLELIST")
    @ApiOperation(value = "快速生成永久token",
        notes = "使用curl生成永久token的格式：curl -G http://域名:端口/sso/token --data-urlencode 'password=密码' -d 'loginId=登录名&tenantId=租户' ")
    @ApiImplicitParams(value = {
        @ApiImplicitParam(paramType = "query", dataType = "String", name = "loginId", value = "登录id", required = true),
        @ApiImplicitParam(paramType = "query", dataType = "String", name = "password", value = "密码", required = true),
        @ApiImplicitParam(paramType = "query", dataType = "String", name = "tenantId", value = "租户id",
            required = true)})
    @RequestMapping(method = RequestMethod.GET, value = "sso/token")
    @ResponseBody
    public BaseResultVO grantForeverToken(@RequestParam String loginId,
        @RequestParam(name = "password", required = true) String password, @RequestParam String tenantId) {

        return JsonResultUtil.getResult(authManager.grantToken(loginId, password, tenantId, null));
    }

    private Boolean getLoginTemplates(String customer) {
        Boolean bool = false;
        try {
            GetResource.class.getClassLoader().getResource("templates/" + customer).getPath();
            bool = true;
        } catch (Exception e) {
        }
        return bool;
    }
}
