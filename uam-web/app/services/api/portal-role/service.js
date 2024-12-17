/*
 * 角色管理
 * 这里的接口对应java中的 PortalRoleController提供的api,具体见java中PortalRoleController.java的实现
 * 2015-10-30 15:14:19
 */
import Ember from 'ember'
var ajax = Ember.$.ajax

export default Ember.Service.extend({
    /*
     * 分配用户组
     */
    allocateGroup: function(id, selGroups) {
        var url = '/portal/role/{id}/allocate/group'
        //替换{参数占位符}
        url = url.replace('{id}', id)
        //生成发请求数据对象
        var data = {}
        data.selGroups = selGroups

        //发送ajax请求
        return ajax({
            method: 'POST',
            url: url,
            data: data
        })
    },

    /*
     * 分配用户
     */
    allocateUser: function(id, selUsers) {
        var url = '/portal/role/{id}/allocate/users'
        //替换{参数占位符}
        url = url.replace('{id}', id)
        //生成发请求数据对象
        var data = {}
        data.selUsers = selUsers

        //发送ajax请求
        return ajax({
            method: 'POST',
            url: url,
            data: data
        })
    },

    /*
     * 获取某个角色信息
     */
    queryById: function(id) {
        var url = '/portal/role/{id}'
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
     * 获取某个角色权限信息
     */
    queryPermissionsById: function(id) {
        var url = '/portal/role/permissions/{id}'
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
     * 获取某个角色用户信息
     */
    queryUsersById: function(id) {
        var url = '/portal/role/users/{id}'
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
     *  获取某个角色用户组信息
     */
    queryGroupsById: function(id) {
        var url = '/portal/role/groups/{id}'
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
    queryRolePageByGroupIdAndCnd: function(form) {
        var url = '/portal/role/page/group'

        //发送ajax请求
        return ajax({
            method: 'GET',
            url: url,
            data: form
        })
    },

    /*
     * 查询用户所关联的角色信息
     */
    queryRolePageByUserIdAndCnd: function(form) {
        var url = '/portal/role/page/user'

        //发送ajax请求
        return ajax({
            method: 'GET',
            url: url,
            data: form
        })
    },

    /*
     * 删除角色信息
     */
    delete: function(id) {
        var url = '/portal/role/{id}'
        //替换{参数占位符}
        url = url.replace('{id}', id)
        //生成发请求数据对象
        var data = {}

        //发送ajax请求
        return ajax({
            method: 'DELETE',
            url: url,
            data: data
        })
    },

    /*
     * 更新角色信息
     * 参数form对象信息
     * ---------参数名称:参数说明-------
     * ---name:名称------
     * ---code:编码------
     * ---remark:备注------
     * ---editType:编辑类型------
     */
    update: function(id, form) {
        var url = '/portal/role/{id}'
        //替换{参数占位符}
        url = url.replace('{id}', id)
        //生成发请求数据对象
        var data = {}
        data = form
        data['_method'] = 'PUT'

        //发送ajax请求
        return ajax({
            method: 'POST',
            url: url,
            data: data
        })
    },

    /*
     * 查询角色信息
     * 参数form对象信息
     * ---------参数名称:参数说明-------
     * ---name:名称------
     * ---code:编码------
     * ---remark:备注------
     */
    queryByPage: function(offset, limit, form) {
        var url = '/portal/role/page/{offset}/{limit}'
        //替换{参数占位符}
        url = url.replace('{offset}', offset)
        //替换{参数占位符}
        url = url.replace('{limit}', limit)
        //生成发请求数据对象
        var data = {}
        data = form

        //发送ajax请求
        return ajax({
            method: 'GET',
            url: url,
            data: data
        })
    },

    /*
     * 创建角色信息
     * 参数form对象信息
     * ---------参数名称:参数说明-------
     * ---name:名称------
     * ---code:编码------
     * ---remark:备注------
     * ---editType:编辑类型------
     */
    create: function(form) {
        var url = '/portal/role/'
        //生成发请求数据对象
        var data = {}
        data = form

        //发送ajax请求
        return ajax({
            method: 'POST',
            url: url,
            data: data
        })
    }
})
