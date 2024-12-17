import Ember from 'ember'
const { set, get, inject } = Ember
const pageSize = 10
export default Ember.Controller.extend({
	tenantSrv: inject.service('api/portal-tenant/service'),
	queryForm: {
		tenantId: null,
		tenantName: null
	},
	filterString: '',
	loading: false,
	editTenant: {},
	infoTenant: {},
	openModalWindow: false,
	openModalWindowTitle: '租户信息',
	openInfoModalWindow: false,
	openInfoModalWindowTitle: '租户信息',
	initTenant: function(model) {
		var meta = {
			offset: 1
		}
		set(model, 'meta', meta)
		set(model, 'indexNumber', 1)
		set(this, 'model', model)
		this.queryTenantByPage(meta.offset, this.queryForm)
	},
	queryTenantByPage: function(offset, queryForm) {
		var tenantSrv = get(this, 'tenantSrv'),
			model = get(this, 'model')
		set(model, 'paginationShow', true)
		set(this, 'loading', true)
		tenantSrv.queryByPage(offset, pageSize, queryForm).then(res => {
			set(this, 'loading', false)
			if ('success' === res.status) {
				set(model, 'data', res.list)
				set(model, 'meta', res.meta)
			} else {
				swal(res.message)
			}
		})
	},
	actions: {
		/**
		 * 分页查询
		 */
		pageClick: function(pageNo) {
			var model = get(this, 'model')

			// 设置页面的number开头
			set(model, 'indexNumber', (pageNo - 1) * pageSize + 1)
			set(model.meta, 'offset', pageNo)
			this.queryTenantByPage(model.meta.offset, this.queryForm)
		},
		/**
		 * 打开租户编辑窗口
		 */
		viewModalWindow: function(type, tenant) {
			set(this, 'openModalWindow', true)
			var model = get(this, 'model')
			if (type === 'add') {
				set(model, 'isAdd', true)
				set(this, 'openModalWindowTitle', '新增租户')
			} else {
				set(this, 'openModalWindowTitle', '修改租户')
				var tenantSrv = get(this, 'tenantSrv')
				tenantSrv.queryTenantById(tenant.id).then(res => {
					if ('success' === res.status) {
						set(this, 'editTenant', res.item)
					} else {
						swal(res.message)
					}
				})
			}
		},
		closeModalWindow: function() {
            set(this, 'openModalWindow', false)
            set(this, 'editTenant', {})
		},
		saveTenant: function() {
			var self = this
			var model = get(this, 'model')
			var tenantSrv = get(this, 'tenantSrv')
			if (model.isAdd) {
				tenantSrv.addUser(this.editTenant).then(res => {
					if ('success' === res.status) {
						swal('新增租户成功')
					} else {
						swal(res.message)
					}
					self.queryTenantByPage(1, this.queryForm)
				})
			} else {
				tenantSrv.updateUser(this.editTenant).then(res => {
					if ('success' === res.status) {
						swal('修改租户成功')
					} else {
						swal(res.message)
					}
					self.queryTenantByPage(1, this.queryForm)
				})
			}
			set(this, 'openModalWindow', false)
			set(this, 'editTenant', {})
		},
		deleteTenant: function(tenant) {
			var self = this
			if (!tenant && Ember.isBlank(tenant.id)) {
				swal('用户信息为空，请确认')
				return
			}
			var tenantSrv = get(this, 'tenantSrv')
			swal(
				{
					title: '是否删除此租户?',
					type: 'warning',
					showCancelButton: true,
					confirmButtonClass: 'btn-danger',
					cancelButtonText: '取消',
					confirmButtonText: '删除',
					closeOnConfirm: false
				},
				function(isConfirm) {
					if (isConfirm) {
						tenantSrv.deleteTenant(tenant.id).then(res => {
							if ('success' === res.status) {
								swal('删除租户成功!')
								self.queryTenantByPage(1, this.queryForm)
							} else {
								swal(res.message)
							}
						})
					}
				}
			)
		},
		viewInfoModalWindow: function(tenant) {
			set(this, 'infoTenant', tenant)
			set(this, 'openInfoModalWindow', true)
		},
		closeInfoModalWindow: function() {
			set(this, 'infoTenant', {})
			set(this, 'openInfoModalWindow', false)
		}
	}
})
