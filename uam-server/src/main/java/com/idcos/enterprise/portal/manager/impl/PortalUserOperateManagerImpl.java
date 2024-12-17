

package com.idcos.enterprise.portal.manager.impl;

// auto generated imports

import com.idcos.cloud.biz.common.check.CommonParamtersChecker;
import com.idcos.cloud.biz.common.check.FormChecker;
import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.cloud.core.common.util.ListUtil;
import com.idcos.cloud.core.common.util.StringUtil;
import com.idcos.enterprise.portal.biz.common.CommonBizException;
import com.idcos.enterprise.portal.biz.common.ResultCode;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessCallback;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessContext;
import com.idcos.enterprise.portal.biz.common.tempalte.BusinessProcessTemplate;
import com.idcos.enterprise.portal.biz.common.utils.Base64Util;
import com.idcos.enterprise.portal.biz.common.utils.CrudUtilService;
import com.idcos.enterprise.portal.biz.common.utils.CurrentUser;
import com.idcos.enterprise.portal.biz.common.utils.PasswordUtil;
import com.idcos.enterprise.portal.convert.PortalTokenConvert;
import com.idcos.enterprise.portal.dal.entity.*;
import com.idcos.enterprise.portal.dal.enums.PortalUserStatusEnum;
import com.idcos.enterprise.portal.dal.enums.IsActiveEnum;
import com.idcos.enterprise.portal.dal.repository.*;
import com.idcos.enterprise.portal.form.ModifyPasswordForm;
import com.idcos.enterprise.portal.form.PortalUserAddForm;
import com.idcos.enterprise.portal.form.PortalUserAllocateGroupForm;
import com.idcos.enterprise.portal.form.PortalUserUpdateForm;
import com.idcos.enterprise.portal.manager.auto.PortalUserOperateManager;
import com.idcos.enterprise.portal.services.PortalUserService;
import com.idcos.enterprise.portal.vo.PortalTokenVO;
import org.apache.commons.lang3.StringUtils;
import org.apache.poi.ss.formula.functions.T;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.BeanUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.thymeleaf.util.ListUtils;

import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.util.Date;
import java.util.List;

/**
 * Manager实现类
 * <p>第一次由自动生成代码工具初始化，后续可以编辑，再次生成的时候不会进行覆盖</p>
 *
 * @author yanlv
 * @version v 1.1 2015-06-09 09:26:24 yanlv Exp $
 */
@Service
public class PortalUserOperateManagerImpl implements PortalUserOperateManager {
    private static final Logger LOGGER = LoggerFactory.getLogger(PortalUserOperateManagerImpl.class);

    @Autowired
    private PortalGroupUserRelRepository portalGroupUserRelRepository;

    @Autowired
    private PortalRoleRepository portalRoleRepository;

    @Autowired
    private PortalUserRepository portalUserRepository;

    @Autowired
    private BusinessProcessTemplate businessProcessTemplate;

    @Autowired
    private CurrentUser currentUser;

    @Autowired
    private PortalDeptRepository portalDeptRepository;

    @Autowired
    private PortalTokenRepository portalTokenRepository;

    @Autowired
    private PortalUserService portalUserService;

    @Autowired
    private PortalTokenConvert portalTenantConvert;

    @Override
    public CommonResult<?> add(final PortalUserAddForm form) {
        final PortalUser user = new PortalUser();
        return businessProcessTemplate.process(new BusinessProcessCallback<Object>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
                //检查密码规则是否合法
                PasswordUtil.checkPassword(form.getPassword());

                if (!form.getPassword().equals(form.getConfirmPassword())) {
                    throw new RuntimeException("两次输入密码不一致，添加失败！");
                }

            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {
                //如果所属部门不空，则检查所属部门。
                if (StringUtils.isNotBlank(form.getDeptId())) {
                    PortalDept dept = portalDeptRepository.findByDeptId(form.getDeptId());
                    if (dept == null) {
                        throw new RuntimeException(String.format("用户所属部门(id: %s)没有找到。", form.getDeptId()));
                    }
                }
                //检测用户登录名是否已存在
                PortalUser tempUser = portalUserRepository.findPortalUserById(form.getTenantId(), form.getLoginId());
                if (tempUser != null) {
                    throw new RuntimeException(String.format("用户登录名(%s)已经存在。", form.getLoginId()));
                }
            }

            @Override
            public Object doBusiness(BusinessProcessContext context) {
                BeanUtils.copyProperties(form, user);
                user.setIsActive(IsActiveEnum.HAS_ACTIVE.getCode());
                user.setCreateUser(currentUser.getUser().getUserId());
                user.setCreateTime(new Date());
                user.setLastModifiedTime(new Date());
                user.setStatus(PortalUserStatusEnum.ENABLED.getCode());
                portalUserRepository.save(user);
                PortalUser portalUser = portalUserRepository.findPortalUserById(form.getTenantId(), form.getLoginId());

                //密码密码后保存到数据库.
                //由于加密需要用户的ID，所以需要保存用户后再更新密码。
                //加密密码
                try {
                    //salt随机产生，生成密码时使用的密码是用户的ID.
                    byte[] salt = PasswordUtil.getSalt();
                    String encriptPW = PasswordUtil.encrypt(user.getPassword(), portalUser.getId(), salt);

                    String saltStr = Base64Util.encode(salt);
                    user.setPassword(encriptPW);
                    user.setSalt(saltStr);
                    portalUserRepository.save(user);
                } catch (Exception e) {
                    throw new RuntimeException("系统错误，请联系管理员。");
                }
                return portalUser.getId();
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
    }

    @Override
    public CommonResult<?> update(final PortalUserUpdateForm form) {
        final PortalUser oldPortalUser = portalUserRepository.findById(form.getId());
        return businessProcessTemplate.process(new BusinessProcessCallback<Object>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
                //如果所属部门不空，则检查所属部门。
                if (StringUtils.isNotBlank(form.getDeptId())) {
                    PortalDept dept = portalDeptRepository.findByDeptId(form.getDeptId());
                    if (dept == null) {
                        throw new RuntimeException(String.format("用户所属部门(id: %s)没有找到。", form.getDeptId()));
                    }
                }
                if (oldPortalUser == null) {
                    throw new RuntimeException("该用户不存在");
                }
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {

            }

            @Override
            public Object doBusiness(BusinessProcessContext context) {
                PortalUser portalUser = new PortalUser();
                BeanUtils.copyProperties(form, portalUser);
                portalUser.setId(oldPortalUser.getId());
                portalUser.setIsActive(oldPortalUser.getIsActive());
                portalUser.setStatus(oldPortalUser.getStatus());
                portalUser.setPassword(oldPortalUser.getPassword());
                portalUser.setSalt(oldPortalUser.getSalt());
                portalUser.setLastModifiedTime(new Date());
                portalUserRepository.save(portalUser);
                return portalUser;
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
    }

    @Override
    public CommonResult<?> enabled(final String id) {
        return businessProcessTemplate.process(new BusinessProcessCallback<Object>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
                PortalUser portalUser = portalUserRepository.findOne(id);
                if (portalUser == null) {
                    throw new RuntimeException("激活失败,此用户不存在!");
                }
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {

            }

            @Override
            public Object doBusiness(BusinessProcessContext context) {
                PortalUser portalUser = portalUserRepository.findOne(id);
                portalUser.setStatus(PortalUserStatusEnum.ENABLED.getCode());
                portalUserRepository.save(portalUser);
                return null;
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
    }

    @Override
    public CommonResult<?> disabled(final String id) {
        return businessProcessTemplate.process(new BusinessProcessCallback<Object>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
                PortalUser portalUser = portalUserRepository.findOne(id);
                if (portalUser == null) {
                    throw new RuntimeException("禁用失败,此用户不存在!");
                }
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {

            }

            @Override
            public Object doBusiness(BusinessProcessContext context) {
                PortalUser portalUser = portalUserRepository.findOne(id);
                portalUser.setStatus(PortalUserStatusEnum.DISABLED.getCode());
                portalUserRepository.save(portalUser);
                return null;
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
    }

    @Override
    public CommonResult<?> delete(final String id) {
        return businessProcessTemplate.process(new BusinessProcessCallback<Object>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
                PortalUser portalUser = portalUserRepository.findOne(id);
                if (portalUser == null) {
                    throw new RuntimeException("删除失败,此用户不存在!");
                }
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {

            }

            @Override
            public Object doBusiness(BusinessProcessContext context) {
                portalUserService.updateIsActionById(IsActiveEnum.NO_ACTIVE.getCode(), id);
                /**
                 * 删除用户与工作组的关系
                 */
                portalGroupUserRelRepository.deleteByUserId(id);
                return null;
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
    }

    @Override
    public CommonResult<?> allocateGroup(final PortalUserAllocateGroupForm form) {
        CommonResult<T> result = businessProcessTemplate.process(new BusinessProcessCallback<T>() {

            @Override
            public void checkParam(BusinessProcessContext context) {
                //FormChecker.check(form); TODO：jar包待下线
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {

            }

            @Override
            public T doBusiness(BusinessProcessContext context) {
                String selGroups = form.getSelGroups();

                //   保存用户组与用户的关联关系
                if (StringUtils.isNoneBlank(selGroups)) {
                    String[] groupIdList = StringUtils.split(selGroups, ",");
                    for (String str : groupIdList) {
                        //根据 :  来截取id和操作类型
                        String groupId = str.split(":")[0];
                        String operType = str.split(":")[1];

                        //D代表删除
                        if ("D".equals(operType)) {
                            portalGroupUserRelRepository.deleteByGroupIdAndUserId(groupId, form.getId());
                        }

                        //I代表新增
                        if ("I".equals(operType)) {
                            PortalGroupUserRel relPo = new PortalGroupUserRel();
                            relPo.setUserId(form.getId());
                            relPo.setGroupId(groupId);
                            if (StringUtils.isNotBlank(form.getTenantId()) && !"null".equals(form.getTenantId())) {
                                relPo.setTenant(form.getTenantId());
                            } else {
                                relPo.setTenant(currentUser.getUser().getTenantId());
                            }
                            portalGroupUserRelRepository.save(relPo);
                        }
                    }
                }
                return null;
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {
            }
        });
        return result;
    }


    @Override
    public CommonResult<?> forbiddenToken(final String tokenId) {
        return businessProcessTemplate.process(new BusinessProcessCallback<Object>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
                CommonParamtersChecker.checkNotBlank(tokenId);
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {
            }

            @Override
            public Object doBusiness(BusinessProcessContext context) {
                PortalToken portalToken = portalTokenRepository.queryTokenByTokenId(tokenId);
                if (portalToken == null) {
                    throw new CommonBizException(ResultCode.STATUS_ERROR, "token不存在！！！");
                }
                portalToken.setIsActive("N");
                portalTokenRepository.save(portalToken);
                return "success";
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
    }

    @Override
    public CommonResult<?> listTokenByTenantIdAndLoginId(final String tenantId, final String loginId) {
        return businessProcessTemplate.process(new BusinessProcessCallback<List<PortalTokenVO>>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
                CommonParamtersChecker.checkNotBlank(tenantId);
                CommonParamtersChecker.checkNotBlank(loginId);
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {
            }

            @Override
            public List<PortalTokenVO> doBusiness(BusinessProcessContext context) {
                List<PortalToken> portalTokens = portalTokenRepository.listTokenByTenantIdAndLoginId(tenantId, loginId);
                return ListUtil.transform(portalTokens, portalTenantConvert);
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
    }

    @Override
    public CommonResult<?> modifyPassword(final ModifyPasswordForm form) {
        return businessProcessTemplate.process(new BusinessProcessCallback<Object>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {
            }

            @Override
            public Object doBusiness(BusinessProcessContext context) {
                PortalUser portalUser = portalUserRepository.findById(form.getUserId());
                if (portalUser == null) {
                    throw new CommonBizException("该用户不存在");
                }
                //检查新密码规则是否合法
                PasswordUtil.checkPassword(form.getNewPassword());
                //检查老的密码是否与传过来的密码一样。
                if (!StringUtil.isEmpty(portalUser.getPassword())) {
                    //解密用户密码。
                    try {
                        String cipherText = portalUser.getPassword();
                        byte[] salt = Base64Util.decode(portalUser.getSalt());
                        //解密时的解密用的密码是生成时的密码，也就是用户的ID。
                        String plainText = PasswordUtil.decrypt(cipherText, form.getUserId(), salt);
                        if (!plainText.equals(form.getOldPassword())) {
                            throw new CommonBizException(
                                    String.format("用户:%s输入的旧密码与系统中存在的密码不相同，不能修改。", portalUser.getLoginId()));
                        }

                    } catch (UnsupportedEncodingException e) {
                        LOGGER.error(String.format("解密用户:%s密码的登录密码时出错。", portalUser.getLoginId()), e);
                        throw new CommonBizException(String.format("解密用户:%s密码的登录密码时出错。", portalUser.getLoginId()));
                    } catch (IOException e) {
                        throw new CommonBizException("系统错误，请联系管理员。");
                    }
                }
                //加密密码
                try {
                    //salt随机产生，密码是用户的ID.
                    byte[] salt = PasswordUtil.getSalt();
                    String encriptPW = PasswordUtil.encrypt(form.getNewPassword(), form.getUserId(), salt);
                    //修改
                    portalUser.setPassword(encriptPW);
                    portalUser.setSalt(Base64Util.encode(salt));
                    portalUserRepository.save(portalUser);
                } catch (Exception e) {
                    throw new CommonBizException("系统错误，请联系管理员。");
                }
                return "新密码为" + form.getNewPassword();
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
    }

    @Override
    public CommonResult<?> resetPassword(final String userId) {
        return businessProcessTemplate.process(new BusinessProcessCallback<Object>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
                CommonParamtersChecker.checkNotBlank(userId);
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {

            }

            @Override
            public Object doBusiness(BusinessProcessContext context) {
                PortalUser portalUser = portalUserRepository.findById(userId);
                if (portalUser == null) {
                    throw new CommonBizException("该用户不存在");
                }
                //随机生成一个8位的密码用于用户登录。
                String newPassword = PasswordUtil.getStringRandom(8);
                //加密
                String encriptPW;
                try {
                    encriptPW = PasswordUtil.encrypt(newPassword, userId, Base64Util.decode(portalUser.getSalt()));
                } catch (UnsupportedEncodingException e) {
                    LOGGER.error(String.format("加密用户:(%s,%s)的登录密码时出错。", portalUser.getLoginId(), portalUser.getName()),
                            e);
                    throw new CommonBizException(
                            String.format("加密用户:(%s,%s)的登录密码时出错。", portalUser.getLoginId(), portalUser.getName()));
                } catch (IOException e) {
                    throw new CommonBizException("系统错误，请联系管理员。");
                }
                portalUser.setPassword(encriptPW);
                portalUser.setLastModifiedTime(new Date());
                portalUserRepository.save(portalUser);
                return newPassword;
            }

            @Override
            public void exceptionProcess(CommonBizException exception, BusinessProcessContext context) {

            }
        });
    }

}
