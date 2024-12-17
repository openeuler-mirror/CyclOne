/*
 * 权限资源管理
 * 这里的接口对应java中的 PortalUserGroupController提供的api,具体见java中PortalUserGroupController.java的实现
 * 2015-10-28 14:30:49
 */
import Ember from 'ember';
var ajax = Ember.$.ajax;

export default Ember.Service.extend({



    /*
     * 获取某权限资源信息
     */
    queryByCode: function(code) {
        var url = "/portal/res/{code}";
        //替换{参数占位符}
        url = url.replace("{code}", code);
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
     * 删除权限资源
     */
    delete: function(id) {
        var url = "/portal/res/{id}";
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
     * 更新权限资源
     */
    update: function(id, form) {
        var url = "/portal/res/{id}";
        //替换{参数占位符}
        url = url.replace("{id}", id);
        //生成发请求数据对象
        var data = {};
        data = form;
        //发送ajax请求
        return ajax({
            'method': 'PUT',
            'url': url,
            'data': data,
        });

    },

    /*
     * 查询权限资源信息
     * 参数form对象信息
     * ---------参数名称:参数说明-------
     * ---name:名称------
     * ---remark:备注------
     * ---selectRoles:所属角色------
     * ---selectUsers:包含用户------
     */
    queryAll: function() {
        var url = "/portal/res/all";
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
     * 创建权限资源信息
     */
    create: function(form) {
        var url = "/portal/res/";
        //生成发请求数据对象
        var data = form;

        //发送ajax请求
        return ajax({
            'method': 'POST',
            'url': url,
            'data': data,
        });

    }

});
