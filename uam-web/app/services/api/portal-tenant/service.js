import Ember from 'ember'
var ajax = Ember.$.ajax

export default Ember.Service.extend({
	/**
	 * 分页查询用户组
	 */
	queryByPage: function(offset, limit, form) {
		var url = `/portal/tenant/pageList/${offset}/${limit}`

		//发送ajax请求
		return ajax({
			method: 'GET',
			url: url,
			data: form
		})
	},
	/**
	 * 新增租户
	 */
	addUser: function(data) {
		var url = `/portal/tenant`

		//发送ajax请求
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
		var url = `/portal/tenant`

		//发送ajax请求
		return ajax({
			method: 'PUT',
			url: url,
			data: data
		})
	},
	/**
	 * 获取租户信息
	 * @param  {[type]} id [description]
	 * @return {[type]}    [description]
	 */
	queryTenantById: function(id) {
		var url = `/portal/tenant/${id}`

		//发送ajax请求
		return ajax({
			method: 'GET',
			url: url
		})
	},
	deleteTenant: function(id) {
		var url = `/portal/tenant/${id}`

		//发送ajax请求
		return ajax({
			method: 'DELETE',
			url: url
		})
	}
})
