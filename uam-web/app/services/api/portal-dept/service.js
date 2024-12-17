import Ember from 'ember'
var ajax = Ember.$.ajax

export default Ember.Service.extend({
	/**
	 * 新增部门
	 */
	addDept: function(data) {
		var url = '/portal/dept'

		//发送ajax请求
		return ajax({
			method: 'POST',
			url: url,
			data: data
		})
	},

	/**
	 * 获取部门的信息
	 */
	queryDeptById: function(id) {
		var url = `/portal/dept/${id}`
		var data = {}

		//发送ajax请求
		return ajax({
			method: 'GET',
			url: url,
			data: data
		})
	},

	/**
	 * 修改部门
	 */
	updateDept: function(data) {
		var url = '/portal/dept'

		//发送ajax请求
		return ajax({
			method: 'PUT',
			url: url,
			data: data
		})
	},

	deleteDept: function(id) {
		var url = `/portal/dept/${id}`
		var data = {}

		//发送ajax请求
		return ajax({
			method: 'DELETE',
			url: url,
			data: data
		})
	},

	/**
     * 查询部门所关联的角色信息
     */
	queryRolesByDeptId: function(id) {
		var url = `/portal/dept/roles/${id}`

		var data = {}

		//发送ajax请求
		return ajax({
			method: 'GET',
			url: url,
			data: data
		})
	},

	/**
	 * 保存部门与角色的绑定关系
	 */
	saveDeptToRole: function(id, roleIds) {
		var url = `/portal/dept/roles/${id}`

		var data = {}
		data.roleIds = roleIds.join(',')

		//发送ajax请求
		return ajax({
			method: 'POST',
			url: url,
			data: data
		})
	}
})
