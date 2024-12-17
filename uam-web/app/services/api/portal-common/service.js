/*
* 统一门户公共查询管理
* 这里的接口对应java中的 PortalCommonController提供的api,具体见java中PortalCommonController.java的实现
* 2015-10-27 14:18:30
*/
import Ember from 'ember';
var ajax = Ember.$.ajax;

export default Ember.Service.extend({


    /*
    * 获取所有网段信息
    */
    queryNetSegList : function() {
        var url = "/portal/common/netSegListAction";     
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
    * 获取所有用户组信息
    */
    getUserGroupList: function(tenantId) {
        var url = "/portal/common/userGroupListAction";     
		//生成发请求数据对象
		var data = {};
        data.tenantId=tenantId;
		//发送ajax请求
        return ajax({
			'method': 'GET',
			'url': url,
			'data': data,
        });
        
    }, 

    /*
    * 获取功能树
    */
    queryFuncTree : function(showFuncOper) {
        var url = "/portal/common/funcTreeAction";     
		//生成发请求数据对象
		var data = {};
		data.showFuncOper = showFuncOper;

		//发送ajax请求
        return ajax({
			'method': 'GET',
			'url': url,
			'data': data,
        });
        
    }, 

    /*
    * 获取所有角色信息
    */
    getRoleList : function() {
        var url = "/portal/common/roleListAction";     
		//生成发请求数据对象
		var data = {};

		//发送ajax请求
        return ajax({
			'method': 'GET',
			'url': url,
			'data': data,
        });
        
    }

});