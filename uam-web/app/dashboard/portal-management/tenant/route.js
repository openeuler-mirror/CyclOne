import Ember from 'ember'

export default Ember.Route.extend({
	model: function() {
		var columns = [
			{
				propertyName: 'tenantName',
				title: '名称',
				template: 'partial/tenantLookUp'
			},
			{
				propertyName: 'tenantId',
				title: '租户编码'
			},
			{
				title: '操作',
				template: 'partial/tenantOper'
			}
		]
		return { columns: columns }
	},
	setupController: function(controller, model) {
		controller.initTenant(model)
		controller.set('model', model)
	}
})
