import Ember from 'ember'
const { get, set, inject } = Ember

export default Ember.Route.extend({
    queryParams: {
        deptID: {
            refreshModel: true
        }
    },
    permissions: [],
    model: function(param) {
        if (!param.deptID) {
            param['deptID'] = ''
        }
        let ifTenant = window.localStorage.getItem("IS_MULTI_TENANT");
        this.set('permissions', window.sessionStorage.getItem("permissions").split(','))

        return Ember.RSVP.hash({
            columns: [
                {propertyName: 'loginId',title: '登录名'},
                {propertyName: 'name',title: '姓名',template: 'partial/userLookUp'},
                ifTenant === 'true' ? {propertyName: 'tenantId',title: '租户'} : {},
                {propertyName: 'selGroups',title: '所属用户组',template: 'partial/userSelGroups'},
                {propertyName: 'deptName',title: '部门名称'},
                {propertyName: 'sourceType',title: '来源',template: 'partial/userSourcetype'},
                {propertyName: 'lastModifiedTime',title: '最后修改时间'},
                {propertyName: 'status',title: '状态',template: 'partial/userStatus'},
                this.permissions.includes("user.operate")? {title: '操作',template: 'partial/userOper'}: {}
                
            ],
            deptID: param.deptID
        })
        
    },
    setupController: function(controller, model) {
        set(model, 'permissions', this.permissions);
        controller.set('model', model)
        controller.initUser(model)
    }
    
})
