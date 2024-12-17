/*
 * 授权管理
 * 这里的接口对应java中的 PortalPermissionController提供的api,具体见java中PortalPermissionController.java的实现
 * 2015-11-05 16:44:09
 */
import Ember from 'ember';
var ajax = Ember.$.ajax;

export default Ember.Service.extend({


    /*
     * 获取授权资源ID
     * 参数form对象信息
     * ---------参数名称:参数说明-------
     * ---authResType:授权资源类型------
     * ---authObjType:授权对象类型------
     * ---authObjId:授权对象ID------
     */
    queryAuthRes: function(form) {
        var url = "/portal/permission/queryAuthResAction";
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
     * 保存授权对象关系
     * 参数form对象信息
     * ---------参数名称:参数说明-------
     * ---authResType:授权资源类型------
     * ---authResId:授权资源ID------
     * ---authObjType:授权对象类型------
     * ---authObjIds:null------
     */
    saveAuthObj: function(form) {
        var url = "/portal/permission/saveAuthObjAction";
        //生成发请求数据对象
        var data = {};
        data = form;

        //发送ajax请求
        return ajax({
            'method': 'POST',
            'url': url,
            'data': data,
        });

    },

    /*
     * 保存授权资源关系
     * 参数form对象信息
     * ---------参数名称:参数说明-------
     * ---authResType:授权资源类型------
     * ---authResIds:null------
     * ---authObjType:授权对象类型------
     * ---authObjId:授权对象ID------
     */
    saveAuthRes: function(permissionInfo) {
        var url = "/portal/permission/saveAuthResAction";
        //生成发请求数据对象
        var data = {};
        data.permissionInfo = permissionInfo;

        //发送ajax请求
        return ajax({
            'method': 'POST',
            'datatypec':'json',
            'url': url,
            'data': data,
        });

    },

    /*
     * 获取授权对象ID
     * 参数form对象信息
     * ---------参数名称:参数说明-------
     * ---authResType:授权资源类型------
     * ---authResId:授权资源ID------
     * ---authObjType:授权对象类型------
     */
    queryAuthObj: function(form) {
        var url = "/portal/permission/queryAuthObjAction";
        //生成发请求数据对象
        var data = {};
        data = form;

        //发送ajax请求
        return ajax({
            'method': 'GET',
            'url': url,
            'data': data,
        });

    }

});
