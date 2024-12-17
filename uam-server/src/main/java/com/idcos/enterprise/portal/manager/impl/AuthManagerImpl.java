

package com.idcos.enterprise.portal.manager.impl;

// auto generated imports

import com.idcos.cloud.biz.common.check.CommonParamtersChecker;
import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.cloud.core.common.util.DateUtil;
import com.idcos.cloud.core.common.util.StringUtil;
import com.idcos.cloud.core.common.util.UUIDUtil;
import com.idcos.enterprise.portal.biz.common.CommonBizException;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryCallback;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessQueryTemplate;
import com.idcos.enterprise.portal.biz.common.utils.Base64Util;
import com.idcos.enterprise.portal.biz.common.utils.CrcUtil;
import com.idcos.enterprise.portal.biz.common.utils.CurrentUser;
import com.idcos.enterprise.portal.dal.entity.PortalTenant;
import com.idcos.enterprise.portal.dal.entity.PortalToken;
import com.idcos.enterprise.portal.dal.entity.PortalUser;
import com.idcos.enterprise.portal.dal.repository.PortalTokenRepository;
import com.idcos.enterprise.portal.dal.repository.PortalUserRepository;
import com.idcos.enterprise.portal.form.LoginForm;
import com.idcos.enterprise.portal.manager.auto.AuthManager;
import com.idcos.enterprise.portal.vo.JwtTokenVO;
import com.idcos.enterprise.portal.web.GlobalValue;
import com.idcos.enterprise.sso.CheckLoginService;
import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.SignatureAlgorithm;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.io.UnsupportedEncodingException;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.Iterator;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import static com.idcos.cloud.core.common.util.DateUtil.simpleFormat;
import static com.idcos.enterprise.portal.UamConstant.ADMIN;
import static com.idcos.enterprise.portal.UamConstant.USER_ID;

/**
 * Manager实现类
 * <p>第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖</p>
 *
 * @author yanlv
 * @version v 1.1 2015-06-09 09:26:24 yanlv Exp $
 */
@Service
public class AuthManagerImpl implements AuthManager {
    @Autowired
    private BusinessQueryTemplate businessQueryTemplate;

    @Autowired
    private GlobalValue globalValue;

    @Autowired
    private CheckLoginService checkLoginService;

    @Autowired
    private PortalUserRepository portalUserRepository;

    @Autowired
    private PortalTokenRepository portalTokenRepository;

    @Autowired
    private CurrentUser currentUser;

    private static final Logger logger = LoggerFactory.getLogger(AuthManagerImpl.class);

    private static final String NEVER_EXPIRE_TIME = "2037-10-01 15:00";

    @Override
    public String getAuth(LoginForm loginForm) {
        PortalUser portalUser = checkLoginService.checkUser(loginForm.getLoginId(), loginForm.getTenantId());
        PortalTenant portalTenant = checkLoginService.checkTenant(portalUser.getTenantId());
        return "token=" + buildToken(portalUser, portalTenant,
                System.currentTimeMillis() / 1000L + Integer.valueOf(globalValue.getAccessTimeout()));
    }

    @Override
    public CommonResult<?> grantToken(final String loginId, final String password, final String tenantId,
                                      final String time) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {

            @Override
            public Object doQuery() {
                PortalTenant portalTenant = checkLoginService.checkTenant(tenantId);
                PortalUser portalUser = checkLoginService.checkUser(loginId, tenantId);
                checkLoginService.checkPassword(portalUser, password);
                //返回token(生成永久token或者传入给定过期时间token)；
                String expireTime = StringUtil.isBlank(time) ? NEVER_EXPIRE_TIME : time;
                return buildToken(portalUser, portalTenant, DateUtil.strToSimpleFormat(expireTime).getTime() / 1000L);
            }

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(loginId);
                CommonParamtersChecker.checkNotBlank(tenantId);
                CommonParamtersChecker.checkNotBlank(password);
            }
        });
    }

    @Override
    public CommonResult<?> grantTokenByAdmin(final String loginId, final String tenantId) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {

            @Override
            public Object doQuery() {
                if (!ADMIN.equals(currentUser.getUser().getLoginId())) {
                    throw new CommonBizException("非admin用户不允许给其他用户发放token！！！");
                }
                PortalTenant portalTenant = checkLoginService.checkTenant(tenantId);
                PortalUser portalUser = checkLoginService.checkUser(loginId, tenantId);
                //返回token(生成永久token)；
                String expireTime = NEVER_EXPIRE_TIME;
                return buildToken(portalUser, portalTenant, DateUtil.strToSimpleFormat(expireTime).getTime() / 1000L);
            }

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(loginId);
                CommonParamtersChecker.checkNotBlank(tenantId);
            }
        });
    }


    /**
     * 生成token
     *
     * @return
     */
    private String buildToken(PortalUser user, PortalTenant portalTenant, Long exp) {
        String token = Jwts.builder().setSubject(user.getName())
                .signWith(SignatureAlgorithm.HS256, globalValue.getSecretKey()).claim("userId", user.getId())
                .claim("name", user.getName()).claim("loginId", user.getLoginId()).claim("loginName", user.getLoginId())
                .claim("tenantId", user.getTenantId()).claim("timeout", Integer.valueOf(globalValue.getAccessTimeout()))
                .claim("exp", exp).claim("creatTime", System.currentTimeMillis())
                .claim("tenantName", portalTenant.getDisplayName()).compact();

        //将token写入数据库
        PortalToken portalToken = new PortalToken();
        portalToken.setId(UUIDUtil.get());
        portalToken.setName(token);
        portalToken.setTokenCrc(CrcUtil.crc(token));
        portalToken.setLoginId(user.getLoginId());
        portalToken.setExpireTime(new Date(exp * 1000));
        portalToken.setGmtCreate(new Date());
        portalToken.setGmtModified(new Date());
        portalToken.setIsActive("Y");
        portalToken.setTenantId(user.getTenantId());
        portalTokenRepository.save(portalToken);
        logger.info("===============用户：" + user.getLoginId() + "@" + user.getTenantId() + "生成token:" + token);
        return token;

    }

    @Override
    public CommonResult<?> parseToken(final String token) {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {

            @Override
            public Object doQuery() {
                Claims claims = Jwts.parser().setSigningKey(globalValue.getSecretKey()).parseClaimsJws(token).getBody();
                if (claims == null || claims.get(USER_ID) == null) {
                    throw new CommonBizException("token中的用户id为空值导致解码失败");
                }

                JwtTokenVO jwtTokenVO = new JwtTokenVO();
                jwtTokenVO.setUserId(claims.get("userId").toString());
                jwtTokenVO.setName(claims.get("name") == null ? "" : claims.get("name").toString());
                jwtTokenVO.setLoginId(claims.get("loginId") == null ? "" : claims.get("loginId").toString());
                jwtTokenVO.setTenantId(claims.get("tenantId") == null ? "" : claims.get("tenantId").toString());
                jwtTokenVO.setTenantName(claims.get("tenantName") == null ? "" : claims.get("tenantName").toString());

                String expireTime;
                try {
                    expireTime = simpleFormat(new Date((Integer) claims.get("exp") * 1000L));
                } catch (Exception e) {
                    throw new RuntimeException("解析token过期时间异常！");
                }
                jwtTokenVO.setExpireTime(expireTime);

                String creatTime;
                try {
                    SimpleDateFormat dateFormat = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
                    creatTime = dateFormat.format(claims.get("creatTime"));
                } catch (Exception e) {
                    throw new RuntimeException("解析token创建时间异常！");
                }
                jwtTokenVO.setCreatTime(creatTime);
                return jwtTokenVO;
            }

            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(token);
            }
        });
    }

    private static final Pattern MODIFY_SALT_PATTERN = Pattern.compile("[A-Za-z0-9=]+");

    @Override
    public CommonResult<?> modifySalt() {
        return businessQueryTemplate.process(new BusinessQueryCallback<Object>() {
            @Override
            public Object doQuery() {
                List<PortalUser> portalUsers = portalUserRepository.findAll();
                List<PortalUser> newPortalUsers = new ArrayList<>();

                Iterator<PortalUser> iterator1 = portalUsers.iterator();
                while (iterator1.hasNext()) {
                    PortalUser portalUser = iterator1.next();
                    Matcher matcher = MODIFY_SALT_PATTERN.matcher(portalUser.getSalt());
                    if (matcher.matches()) {
                        throw new CommonBizException("接口已经执行过，不需要再执行了。");
                    }
                }

                Iterator<PortalUser> iterator2 = portalUsers.iterator();
                while (iterator2.hasNext()) {
                    PortalUser portalUser = iterator2.next();
                    try {
                        portalUser.setSalt(Base64Util.encode(portalUser.getSalt().getBytes("ISO-8859-1")));
                    } catch (UnsupportedEncodingException e) {
                        throw new CommonBizException("系统错误，请联系管理员。");
                    }
                    newPortalUsers.add(portalUser);
                }
                portalUserRepository.save(newPortalUsers);
                return null;
            }

            @Override
            public void checkParam() {

            }
        });
    }
}