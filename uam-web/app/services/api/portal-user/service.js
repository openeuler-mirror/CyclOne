/*
 * 用户管理
 * 这里的接口对应java中的 PortalUserController提供的api,具体见java中PortalUserController.java的实现
 * 2015-11-07 18:52:53
 */
import Ember from 'ember'
var ajax = Ember.$.ajax

export default Ember.Service.extend({
    /*
     * 分配用户组
     */
    allocateGroup: function(userInfo) {
        var url = '/portal/user/allocate/userGroup'
        //替换{参数占位符}
        // url = url.replace("{id}", id);
        //生成发请求数据对象
        var data = {}
        data.userInfo = userInfo

        //发送ajax请求
        return ajax({
            method: 'POST',
            url: url,
            datatype: 'json',
            data: data
        })
    },

    /*
     * 分配用户组
     */
    allocate: function(userInfo) {
        var url = '/portal/user/'
        //生成发请求数据对象
        var data = {}
        data = userInfo
        var tenant = window.sessionStorage.getItem('tenant')
        data.tenantId = tenant
        //发送ajax请求
        return ajax({
            method: 'POST',
            url: url,
            data: data
        })
    },

    /**
     * 获取某用户的权限信息
     * @param  {[type]} id [description]
     * @return {[type]}    [description]
     */
    queryPermissionsById: function(id) {
        var url = '/portal/user/permissions/{id}'
        //替换{参数占位符}
        url = url.replace('{id}', id)
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
     * 获获取某用户角色信息
     * @param  {[type]} id [description]
     * @return {[type]}    [description]
     */
    queryRolesById: function(id) {
        var url = '/portal/user/roles/{id}'
        //替换{参数占位符}
        url = url.replace('{id}', id)
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
     * 获取某用户用户组信息
     * @param  {[type]} id [description]
     * @return {[type]}    [description]
     */
    queryGroupsById: function(id) {
        var url = '/portal/user/groups/{id}'
        //替换{参数占位符}
        url = url.replace('{id}', id)
        //生成发请求数据对象
        var data = {}

        //发送ajax请求
        return ajax({
            method: 'GET',
            url: url,
            data: data
        })
    },

    /*
     * 查询用户组所关联的角色信息
     */
    queryUserPageByGroupIdAndCnd: function(form) {
        var url = '/portal/user/page/group'

        //发送ajax请求
        return ajax({
            method: 'GET',
            url: url,
            data: form
        })
    },

    /*
     * 查询用户组所关联的角色信息
     */
    queryUserPageByRoleIdAndCnd: function(form) {
        var url = '/portal/user/page/role'

        //发送ajax请求
        return ajax({
            method: 'GET',
            url: url,
            data: form
        })
    },

    /*
     * 获取某用户信息
     */
    queryById: function(id) {
        var url = '/portal/user/{id}'
        //替换{参数占位符}
        url = url.replace('{id}', id)
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
     * 新增用户
     * @param {[type]} data [description]
     */
    addUser: function(data) {
        var url = '/portal/user/add'
        return ajax({
            method: 'POST',
            url: url,
            data: data
        })
    },

    /**
     * 修改用户
     * @param  {[type]} data [description]
     * @return {[type]}      [description]
     */
    updateUser: function(data) {
        var url = '/portal/user'
        return ajax({
            method: 'PUT',
            url: url,
            data: data
        })
    },

    /**
     * 激活用户
     * @param  {[type]} id [description]
     * @return {[type]}      [description]
     */
    enabledUser: function(id) {
        var url = '/portal/user/enabled/{id}'
        //替换{参数占位符}
        url = url.replace('{id}', id)
        var data = {}
        return ajax({
            method: 'PUT',
            url: url,
            data: data
        })
    },

    /**
     * 禁用用户
     * @param  {[type]} id [description]
     * @return {[type]}      [description]
     */
    disabledUser: function(id) {
        var url = '/portal/user/disabled/{id}'
        //替换{参数占位符}
        url = url.replace('{id}', id)
        var data = {}
        return ajax({
            method: 'PUT',
            url: url,
            data: data
        })
    },

    /**
     * 删除用户
     */
    deleteUser: function(id) {
        var url = '/portal/user/{id}'
        //替换{参数占位符}
        url = url.replace('{id}', id)
        var data = {}
        return ajax({
            method: 'DELETE',
            url: url,
            data: data
        })
    },

    /**
     * admin给其他用户发放Token
     * @param  {[type]} data [description]
     * @return {[type]}      [description]
     */
    grantToken: function(loginId,tenantId) {
        var url = '/sso/token/admin?loginId='+loginId+'&tenantId='+tenantId;
        var data = {}
        return ajax({
            method: 'GET',
            url: url,
            data: data
        })
    },

    /**
     * 重置密码
     */
    resetPW: function(id) {
        var url = '/portal/user/resetPW'
        var data = {}
        data.userId = id
        return ajax({
            method: 'POST',
            url: url,
            data: data
        })
    },

    /**
     * 根据excel名，返回上传预览
     */
    importPriview: function(id) {
        var url = '/portal/user/importPriview/fileName'
        var data = {}
        data.fileName = id
        return ajax({
            method: 'GET',
            url: url,
            data: data
        })
    },
    /**
     * 导入预览中的合法数据，并删除用户列表缓存
     */
    saveImportUsers: function(id) {
        var url = '/portal/user/importPriview/fileName'
        var data = {}
        data.fileName = id
        return ajax({
            method: 'POST',
            url: url,
            data: data
        })
    }
})
