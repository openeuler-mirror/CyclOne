/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2015-2017 All Rights Reserved.
 */
package com.idcos.enterprise.portal.export;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.*;

import com.idcos.enterprise.portal.export.util.CommonRestResult;
import com.idcos.enterprise.portal.export.util.CommonRestResultUtil;
import com.idcos.enterprise.portal.form.*;
import com.idcos.enterprise.portal.manager.auto.*;

import io.swagger.annotations.Api;

import java.util.HashMap;

/**
 * rbac提供给外部系统的restful接口
 *
 * @author Dana
 * @version AuthRestController.java, v1 2017/9/26 下午2:34 Dana Exp $$
 */
@Controller
@RequestMapping(value = "/rbac/api")
@Api(tags = "14.uam-client访问的接口", description = "AuthRestController")
public class AuthRestController {
    private static final Logger logger = LoggerFactory.getLogger(AuthRestController.class);

    @Autowired
    private PortalUserQueryManager portalUserQueryManager;

    @Autowired
    private PortalUserOperateManager portalUserOperateManager;

    @Autowired
    private PortalDeptQueryManager portalDeptQueryManager;

    @Autowired
    private PortalTenantManager portalTenantManager;

    @Autowired
    private PortalRoleQueryManager portalRoleQueryManager;

    @Autowired
    private PortalUserGroupQueryManager portalUserGroupQueryManager;

    @Autowired
    private PortalUserGroupOperateManager portalUserGroupOperateManager;

    /**
     * 根据userId和appId查询该用户所有信息
     *
     * @return
     */
    @RequestMapping(value = "/authInfo", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult authInfo(@RequestParam(name = "userId", required = false) String userId, @RequestParam String appId,
                                     @RequestParam String token) {
        return CommonRestResultUtil.getResult(portalUserQueryManager.authInfo(userId, appId, token));
    }

    /**
     * 根据loginId和password查询该用户所有信息
     *
     * @return
     */
    @RequestMapping(value = "/userGroup/loginId", method = RequestMethod.POST, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult getAuthInfoByLoginIdAndPw(@RequestBody AuthInfoQueryForm form) {
        return CommonRestResultUtil.getResult(portalUserQueryManager.getAuthInfoByLoginIdAndPw(form.getLoginId(),
                form.getPassword(), form.getTenantId(), form.getAppId()));
    }

    /**
     * 根据租户id查询所有的帐号信息
     *
     * @return
     */
    @RequestMapping(value = "/account/getAllAcountByTenantId", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult getAllAcountByTenantId(@RequestParam String tenantId) {
        return CommonRestResultUtil.getResult(portalUserQueryManager.getAllAcount(tenantId));
    }

    /**
     * 查询所有的帐号信息
     *
     * @return
     */
    @RequestMapping(value = "/account/getAllAcount", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult getAllAcount() {
        return CommonRestResultUtil.getResult(portalUserQueryManager.getAllAcount());
    }

    /**
     * 查询租户下系统的帐号总数，不传租户则查询所有的账号总数
     *
     * @return
     */
    @RequestMapping(value = "/account/count", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult getAccountCount(@RequestParam(name = "tenantId", required = false) String tenantId) {
        return CommonRestResultUtil.getResult(portalUserQueryManager.getAccountCount(tenantId));
    }

    /**
     * 查询部门下的所有的员工(及子部门下的帐号)
     *
     * @return
     */
    @RequestMapping(value = "/account/getByDeptId", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult getByDeptId(@RequestParam String tenantId, @RequestParam String deptId,
                                        @RequestParam(name = "recurse", required = false) String recurse) {
        return CommonRestResultUtil.getResult(portalUserQueryManager.getAcountByDeptId(tenantId, deptId, recurse));
    }

    /**
     * 查询最后登录时间
     *
     * @return
     */
    @RequestMapping(value = "/account/getLastLoginTime", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult getLastLoginTime() {
        return CommonRestResultUtil.getResult(portalUserQueryManager.getLastLoginTime());
    }

    /**
     * 根据登录id获取用户信息。
     *
     * @param accountNo 账号的登录名.
     * @return
     */
    @RequestMapping(value = "/account/getByAccountNo", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult getByAccountNo(@RequestParam String tenantId, @RequestParam String accountNo) {
        return CommonRestResultUtil.getResult(portalUserQueryManager.getByAccountNo(tenantId, accountNo));
    }

    /**
     * 根据用户id获取用户信息。
     *
     * @param id 账号的登录名.
     * @return
     */
    @RequestMapping(value = "/account/id", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult getUserById(@RequestParam String id) {
        return CommonRestResultUtil.getResult(portalUserQueryManager.getUserById(id));
    }

    /**
     * 根据账号登录ids查询账号信息。
     * 目前res只用到了accountNo和Email两个属性，所以后面可以考虑去掉这个方法
     *
     * @param form 账号的登录名.
     * @return
     */
    @RequestMapping(value = "/account/getByAccountNos", method = RequestMethod.POST, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult getByAccountNos(@RequestParam String tenantId, @RequestBody UserListForm form) {
        return CommonRestResultUtil.getResult(portalUserQueryManager.getByAccountNos(tenantId, form.getLoginIdList()));
    }

    /**
     * 分页模糊查询用户列表
     *
     * @param form
     * @return
     */
    @RequestMapping(method = RequestMethod.GET, value = "/pageList", produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult queryPagelist(@ModelAttribute PortalUserQueryPageListForm form) {
        return CommonRestResultUtil.getResult(portalUserQueryManager.queryPageList(form));
    }

    /**
     * 获取所有部门信息
     *
     * @return
     */
    @RequestMapping(value = "/dept/allDept", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult getAllDept(@RequestParam(name = "tenantId", required = false) String tenantId) {
        return CommonRestResultUtil.getResult(portalDeptQueryManager.getAllDept(tenantId));
    }

    /**
     * 获取所有部门的Tree型的树
     *
     * @return
     */
    @RequestMapping(value = "/dept/tree", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult getDeptsTree(@RequestParam String tenantId, @RequestParam String treeStyle) {
        return CommonRestResultUtil.getResult(portalDeptQueryManager.getDeptsTree(tenantId, treeStyle));
    }

    /**
     * 获取所有用户组的Tree型的树
     *
     * @return
     */
    @RequestMapping(value = "/accountGroup/tree", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult getAccountGroupTree(@RequestParam String tenantId, @RequestParam String treeStyle,
                                                @RequestParam String isOpen) {
        return CommonRestResultUtil.getResult(portalUserGroupQueryManager.getAccountGroupTree(tenantId, treeStyle));
    }

    /**
     * 根据权限code查询工作组列表
     *
     * @return
     */
    @RequestMapping(value = "/accountGroup/getListByCode", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult getAccountGroupListByPermissionCode(@RequestParam String code) {
        return CommonRestResultUtil.getResult(portalUserGroupQueryManager.getAccountGroupListByPermissionCode(code));
    }

    /**
     * 根据部门id获取部门信息
     *
     * @return
     */
    @RequestMapping(value = "/dept/get", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult getDeptByDeptId(@RequestParam String tenantId, @RequestParam String deptId) {
        return CommonRestResultUtil.getResult(portalDeptQueryManager.getDeptByDeptId(tenantId, deptId));
    }

    /**
     * 根据部门id获取部门角色
     *
     * @return
     */
    @RequestMapping(value = "/account/getRoles", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult getRolesByDeptId(@RequestParam String tenantId, @RequestParam String deptId) {
        return CommonRestResultUtil.getResult(null);
    }

    /**
     * 获取所有的租户信息
     *
     * @return
     */
    @RequestMapping(value = "/tenant/getAll", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult getAllTenant() {
        return CommonRestResultUtil.getResult(portalTenantManager.getAllTenant());
    }

    /**
     * 通过租户Id获取租户信息
     *
     * @return
     */
    @RequestMapping(value = "/tenant/getByTenantId", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult getTenantByTenantId(@RequestParam String tenantId) {
        return CommonRestResultUtil.getResult(portalTenantManager.getTenantByTenantId(tenantId));
    }

    /**
     * 根据tenantId获取所有角色信息
     *
     * @return
     */
    @RequestMapping(value = "/roles", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult getAllRolesByTenantId(@RequestParam(name = "tenantId", required = false) String tenantId) {
        return CommonRestResultUtil.getResult(portalRoleQueryManager.queryByTenantId(tenantId));
    }

    /**
     * 根据tenantId和用户id获取所有角色信息
     *
     * @return
     */
    @RequestMapping(value = "/roles/accountId", method = RequestMethod.GET, produces = "application/json; charset=UTF-8")
    @ResponseBody
    public CommonRestResult getAllRolesByAccountId(@RequestParam String accountId, @RequestParam String tenantId) {
        return CommonRestResultUtil.getResult(portalRoleQueryManager.queryByAccountNoAndTenantId(accountId, tenantId));
    }

    /**
     * 获取某用户权限信息
     *
     * @param userId,appId
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.GET, value = "/permissions")
    @ResponseBody
    public CommonRestResult queryPermissionsByUserIdAndAppId(@RequestParam String userId,
                                                             @RequestParam(name = "appId") String appId) {
        return CommonRestResultUtil.getResult(portalUserQueryManager.queryPermissionsByUserIdAndAppId(userId, appId));
    }

    /**
     * 获取某用户用户组信息
     *
     * @param id
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.GET, value = "/groups/{id}")
    @ResponseBody
    public CommonRestResult queryGroupsById(@PathVariable String id) {
        return CommonRestResultUtil.getResult(portalUserQueryManager.queryGroupsById(id));
    }

    /**
     * 根据用户组type查询用户组集合
     *
     * @param type
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.GET, value = "/group/type")
    @ResponseBody
    public CommonRestResult queryGroupByType(@RequestParam String type) {
        return CommonRestResultUtil.getResult(portalUserGroupQueryManager.queryGroupByType(type));
    }

    /**
     * 根据用户组type和租户id查询用户组集合
     *
     * @param type
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.GET, value = "/group/type/tenantId")
    @ResponseBody
    public CommonRestResult queryGroupByTypeAndTenantId(@RequestParam String type, @RequestParam String tenantId) {
        return CommonRestResultUtil.getResult(portalUserGroupQueryManager.queryGroupByTypeAndTenantId(type, tenantId));
    }

    /**
     * 根据用户id和用户组type查询用户组集合
     *
     * @param type,userId
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.GET, value = "/group/type/userId")
    @ResponseBody
    public CommonRestResult queryGroupByTypeAndUserId(@RequestParam String type, @RequestParam String userId) {
        return CommonRestResultUtil.getResult(portalUserGroupQueryManager.queryGroupByTypeAndUserId(type, userId));
    }

    /**
     * 根据loginId和tenantId查询用户组集合
     *
     * @param loginId,tenantId
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.GET, value = "/group/loginId")
    @ResponseBody
    public CommonRestResult queryGroupByLoginIdAndTenantId(@RequestParam String loginId,
                                                           @RequestParam String tenantId) {
        return CommonRestResultUtil
                .getResult(portalUserGroupQueryManager.queryGroupByLoginIdAndTenantId(loginId, tenantId));
    }

    /**
     * 修改用户密码
     *
     * @param form
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.POST, value = "/modifyPW")
    @ResponseBody
    public CommonRestResult modifyPassword(@RequestBody ModifyPasswordForm form) {
        return CommonRestResultUtil.getResult(portalUserOperateManager.modifyPassword(form));
    }

    /**
     * 重置用户密码
     *
     * @param userId
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.POST, value = "/resetPW")
    @ResponseBody
    public CommonRestResult resetPassword(@RequestParam String userId) {
        return CommonRestResultUtil.getResult(portalUserOperateManager.resetPassword(userId));
    }

    /**
     * 对外提供的新增用户
     *
     * @param tenantId
     * @param map
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.POST, value = "/account/add")
    @ResponseBody
    public CommonRestResult addAccount(@RequestParam String tenantId, @RequestBody HashMap<String, String> map) {
        PortalUserAddForm userAddForm = new PortalUserAddForm();
        userAddForm.setTenantId(tenantId);
        userAddForm.setName(map.get("name"));
        userAddForm.setLoginId(map.get("accountNo"));
        userAddForm.setMobile1(map.get("mobile"));
        userAddForm.setDeptId(map.get("deptID"));
        userAddForm.setRemark(map.get("remark"));
        userAddForm.setEmail(map.get("email"));
        userAddForm.setOfficeTel1(map.get("officeTel"));
        userAddForm.setPassword(map.get("password"));
        userAddForm.setConfirmPassword(map.get("password"));
        return CommonRestResultUtil.getResult(portalUserOperateManager.add(userAddForm));
    }

    /**
     * 对外提供的修改用户
     *
     * @param tenantId
     * @param map
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.POST, value = "/account/update")
    @ResponseBody
    public CommonRestResult updateAccount(@RequestParam String tenantId, @RequestParam String id, @RequestBody HashMap<String, String> map) {
        PortalUserUpdateForm userUpdateForm = new PortalUserUpdateForm();
        userUpdateForm.setId(id);
        userUpdateForm.setTenantId(tenantId);
        userUpdateForm.setName(map.get("name"));
        userUpdateForm.setLoginId(map.get("accountNo"));
        userUpdateForm.setMobile1(map.get("mobile"));
        userUpdateForm.setDeptId(map.get("deptID"));
        userUpdateForm.setRemark(map.get("remark"));
        userUpdateForm.setEmail(map.get("email"));
        userUpdateForm.setOfficeTel1(map.get("officeTel"));
        return CommonRestResultUtil.getResult(portalUserOperateManager.update(userUpdateForm));
    }

    /**
     * 获取某用户组下所有用户
     *
     * @param groupId
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.GET, value = "/groupId")
    @ResponseBody
    public CommonRestResult queryUserByGroupId(@RequestParam String groupId) {
        return CommonRestResultUtil.getResult(portalUserGroupQueryManager.queryUsersByGroupId(groupId));
    }

    /**
     * 根据用户组名称和租户Id获取某用户组下所有用户
     *
     * @param groupName
     * @param tenantId
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.GET, value = "/users")
    @ResponseBody
    public CommonRestResult queryUserByGroupNameAndTenant(@RequestParam String groupName,
                                                          @RequestParam String tenantId) {
        return CommonRestResultUtil
                .getResult(portalUserGroupQueryManager.queryUserByGroupNameAndTenant(groupName, tenantId));
    }

    /**
     * 删除用户组
     *
     * @param groupId
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.POST, value = "/group/delete")
    @ResponseBody
    public CommonRestResult deleteUserGroup(@RequestParam String groupId) {
        return CommonRestResultUtil.getResult(portalUserGroupOperateManager.delete(groupId));
    }

    /**
     * 新建用户组(包括用户)
     * 用户组名称必传，用户组类型传空值代表default类型，传其他值都是去查找字典表，字典表没有对应数据，则新建字典
     *
     * @param form
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.POST, value = "/group/create")
    @ResponseBody
    public CommonRestResult createUserGroupAndUser(@RequestBody PortalGroupAndUserCreateForm form) {
        return CommonRestResultUtil.getResult(portalUserGroupOperateManager.createUserGroupAndUser(form));
    }

    /**
     * 修改用户组(包括用户)
     * 用户组id和用户组名称必传，用户组类型传空值代表default类型，传其他值都是去查找字典表，字典表没有对应数据，则新建字典
     *
     * @param form
     * @return BaseResultVO
     */
    @RequestMapping(method = RequestMethod.POST, value = "/group/update")
    @ResponseBody
    public CommonRestResult updateUserGroupAndUser(@RequestBody PortalGroupAndUserUpdateForm form) {
        return CommonRestResultUtil.getResult(portalUserGroupOperateManager.updateUserGroupAndUser(form));
    }

}