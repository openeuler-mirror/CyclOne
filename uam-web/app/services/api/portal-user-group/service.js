/*
 * 用户组管理
 * 这里的接口对应java中的 PortalUserGroupController提供的api,具体见java中PortalUserGroupController.java的实现
 * 2015-10-28 14:30:49
 */
import Ember from 'ember';
var ajax = Ember.$.ajax;

export default Ember.Service.extend({


    /*
     * 分配角色信息
     */
    allocateRole: function(id, selectRoles) {
        var url = "/portal/usergroup/{id}/allocate/role";
        //替换{参数占位符}
        url = url.replace("{id}", id);
        //生成发请求数据对象
        var data = {};
        data.selectRoles = selectRoles;

        //发送ajax请求
        return ajax({
            'method': 'POST',
            'url': url,
            'data': data,
        });

    },

    /*
     * 分配用户信息
     */
    allocateUser: function(id, selectUsers) {
        var url = "/portal/usergroup/{id}/allocate/user";
        //替换{参数占位符}
        url = url.replace("{id}", id);
        //生成发请求数据对象
        var data = {};
        data.selectUsers = selectUsers;

        //发送ajax请求
        return ajax({
            'method': 'POST',
            'url': url,
            'data': data,
        });

    },

    /*
     * 获取某用户组角色信息
     */
    queryRolesById: function(id) {
        var url = "/portal/usergroup/roles/{id}";
        //替换{参数占位符}
        url = url.replace("{id}", id);
        //生成发请求数据对象
        var data = {};

        //发送ajax请求
        return ajax({
            'method': 'GET',
            'url': url,
            'data': data,
        });

    },

    /*
     * 获取某用户组用户列表
     */
    queryUsersById: function(id) {
        var url = "/portal/usergroup/users/{id}";
        //替换{参数占位符}
        url = url.replace("{id}", id);
        //生成发请求数据对象
        var data = {};

        //发送ajax请求
        return ajax({
            'method': 'GET',
            'url': url,
            'data': data,
        });

    },

    /*
     * 查询用户所关联的用户组信息
     */
    queryGroupPageByUserIdAndCnd: function(form) {
        var url = "/portal/usergroup/page/user";

        //发送ajax请求
        return ajax({
            'method': 'GET',
            'url': url,
            'data': form
        });

    },

    /*
     * 查询角色所关联的用户组信息
     */
    queryGroupPageByRoleIdAndCnd: function(form) {
        var url = "/portal/usergroup/page/role";

        //发送ajax请求
        return ajax({
            'method': 'GET',
            'url': url,
            'data': form
        });

    },

    /*
     * 获取某用户组权限信息
     */
    queryPermissionsById: function(id) {
        var url = "/portal/usergroup/permissions/{id}";
        //替换{参数占位符}
        url = url.replace("{id}", id);
        //生成发请求数据对象
        var data = {};

        //发送ajax请求
        return ajax({
            'method': 'GET',
            'url': url,
            'data': data,
        });

    },

    /*
     * 获取某用户组信息
     */
    queryById: function(id) {
        var url = "/portal/usergroup/{id}";
        //替换{参数占位符}
        url = url.replace("{id}", id);
        //生成发请求数据对象
        var data = {};

        //发送ajax请求
        return ajax({
            'method': 'GET',
            'url': url,
            'data': data,
        });

    },

    /*
     * 删除用户组
     */
    delete: function(id) {
        var url = "/portal/usergroup/{id}";
        //替换{参数占位符}
        url = url.replace("{id}", id);
        //生成发请求数据对象
        var data = {};

        //发送ajax请求
        return ajax({
            'method': 'DELETE',
            'url': url,
            'data': data,
        });

    },

    /*
     * 更新用户组
     */
    update: function(id, name, type,remark) {
        var url = "/portal/usergroup/{id}";
        //替换{参数占位符}
        url = url.replace("{id}", id);
        //生成发请求数据对象
        var data = {};
        data.name = name;
        data.type = type;
        data.remark = remark;
        data['_method'] = "PUT";

        //发送ajax请求
        return ajax({
            'method': 'POST',
            'url': url,
            'data': data,
        });

    },

    /*
     * 查询用户组信息
     * 参数form对象信息
     * ---------参数名称:参数说明-------
     * ---name:名称------
     * ---remark:备注------
     * ---selectRoles:所属角色------
     * ---selectUsers:包含用户------
     */
    queryByPage: function(offset, limit, form, cnd) {
        var url = "/portal/usergroup/page/{offset}/{limit}?cnd="+cnd;
        //替换{参数占位符}
        url = url.replace("{offset}", offset);
        //替换{参数占位符}
        url = url.replace("{limit}", limit);

        var tenantId = window.sessionStorage.getItem("tenant");
        form.tenantId = tenantId;
        //生成发请求数据对象
        var data = {};
        data = form;

        //发送ajax请求
        return ajax({
            'method': 'GET',
            'url': url,
            'data': data,
        });

    },

    /*
     * 创建用户组信息
     */
    create: function(name, type,remark) {
        var url = "/portal/usergroup/";
        //生成发请求数据对象
        var data = {};

        var tenantId = window.sessionStorage.getItem("tenant");
        data.tenantId = tenantId;
        data.name = name;
        data.type = type;
        data.remark = remark;

        //发送ajax请求
        return ajax({
            'method': 'POST',
            'url': url,
            'data': data,
        });

    }

});