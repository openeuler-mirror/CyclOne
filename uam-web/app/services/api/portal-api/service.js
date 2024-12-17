/*
 * 权限资源管理
 * 这里的接口对应java中的 PortalUserGroupController提供的api,具体见java中PortalUserGroupController.java的实现
 * 2015-10-28 14:30:49
 */
import Ember from 'ember'
var ajax = Ember.$.ajax

export default Ember.Service.extend({
    /**
     * 查询部门信息
     * @param  {[type]} code [description]
     * @return {[type]}      [description]
     */
    queryDepartment: function(tenantId) {
        var url = '/rbac/api/dept/tree'
        //生成发请求数据对象
        var data = {}
        data.tenantId = tenantId
        data.treeStyle = 'ioTree'
        //发送ajax请求
        return ajax({
            method: 'GET',
            url: url,
            data: data
        })
    },

    /**
     * 根据deptID获取权限资源数据
     * @param  {[type]} resType [description]
     * @return {[type]}         [description]
     */
    queryResByResType: function(resType,tenantId) {
        var url = '/api/v1/res'
        //生成发请求数据对象
        var data = {}
        data.resType = resType;
        data.tenantId = tenantId;

        //发送ajax请求
        return ajax({
            method: 'GET',
            url: url,
            data: data
        })
    },

    /**
     * 根据部门id查询人员信息
     * @param  {[type]} deptID [description]
     * @return {[type]}        [description]
     */
    queryUserByDeptId: function(deptId, pageNo, pageSize, tenant, name) {
        var url = '/portal/user/pageList'
        //生成发请求数据对象
        var data = {}
        data.deptId = deptId
        data.pageNo = pageNo
        data.pageSize = pageSize
        data.tenantId = tenant
        data.name = name

        //发送ajax请求
        return ajax({
            method: 'GET',
            url: url,
            data: data
        })
    },

    /**
     * 根据部门id查询人员信息
     * @param  {[type]} deptID [description]
     * @return {[type]}        [description]
     */
    queryUserByCnd: function(tenantId, cnd, pageNo, pageSize) {
        var url = '/api/v1/users/cnd'
        //生成发请求数据对象
        var data = {}
        data.tenantId = tenantId
        data.cnd = cnd
        data.pageNo = pageNo
        data.pageSize = pageSize

        //发送ajax请求
        return ajax({
            method: 'GET',
            url: url,
            data: data
        })
    },

    /**
     * 根据部门id查询人员信息
     * @param  {[type]} deptID [description]
     * @return {[type]}        [description]
     */
    queryUserById: function(userId) {
        var url = '/portal/user/' + userId
        //生成发请求数据对象
        var data = {}

        //发送ajax请求
        return ajax({
            method: 'GET',
            url: url,
            data: data
        })
    },

    /**
     * 查询部门信息
     * @param  {[type]} code [description]
     * @return {[type]}      [description]
     */
    queryTenantList: function() {
        var url = '/rbac/api/tenant/getAll'
        //生成发请求数据对象
        var data = {}

        //发送ajax请求
        return ajax({
            method: 'GET',
            url: url,
            data: data
        })
    },

    /**
     * 查询是否是多租户
     * @return {[type]}      [description]
     */
    isMultiTenant: function() {
        var url = '/auth/isMultiTenant'
        //生成发请求数据对象
        var data = {}

        //发送ajax请求
        return ajax({
            method: 'GET',
            url: url,
            data: data
        })
    },

    /**
     * 根据类型和租户id查询用户组集合
     * @return {[type]}      [description]
     */
    getUserGroupByTypeAndTenantId: function(type,tenantId) {
        var url = '/portal/usergroup/group/type/tenantId'
        //生成发请求数据对象
        var data = {}
        data.type = type
        data.tenantId = tenantId

        //发送ajax请求
        return ajax({
            method: 'GET',
            url: url,
            data: data
        })
    },

    /**
     * 查询字典的用户组类型
     * @return {[type]}      [description]
     */
    getUserGroupType: function(tenantId) {
        var url = '/portal/usergroup/userGroupType'
        //生成发请求数据对象
        var data = {}
        data.tenantId = tenantId

        //发送ajax请求
        return ajax({
            method: 'GET',
            url: url,
            data: data
        })
    },

    /**
     * 查询字典的用户状态
     * @return {[type]}      [description]
     */
    getUserStatus: function() {
        var url = '/portal/user/status'
        //生成发请求数据对象
        var data = {}

        //发送ajax请求
        return ajax({
            method: 'GET',
            url: url,
            data: data
        })
    }
})
