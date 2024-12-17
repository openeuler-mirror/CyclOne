import Ember from 'ember'
var hash = Ember.RSVP.hash
export default Ember.Route.extend({
    queryParams: {
        id: {
            refreshModel: true
        }
    },
    model: function(param) {
        if (!param.id) {
            param['id'] = ''
        }
        return hash({
        columns: [
            {
                propertyName: 'loginId',
                title: '登录名',
                template: 'partial/importPriview/loginId'
            },
            {
                propertyName: 'name',
                title: '姓名',
                template: 'partial/importPriview/name'
            },
            {
                propertyName: 'deptFullName',
                title: '具体部门名称',
                template: 'partial/importPriview/deptFullName'
            },
            {
                propertyName: 'title',
                title: '职务',
                template: 'partial/importPriview/title'
            },
            {
                propertyName: 'password',
                title: '密码',
                template: 'partial/importPriview/password'
            },
            {
                propertyName: 'email',
                title: '邮箱',
                template: 'partial/importPriview/email'
            },
            {
                propertyName: 'mobile1',
                title: '移动电话1',
                template: 'partial/importPriview/mobile'
            },
            {
                propertyName: 'officeTel1',
                title: '办公电话1',
                template: 'partial/importPriview/officeTel'
            },
            {
                propertyName: 'sourceType',
                title: '数据来源',
                template: 'partial/importPriview/sourceType'
            },
            {
                propertyName: 'remark',
                title: '备注',
                template: 'partial/importPriview/remark'
            },
            {
                propertyName: 'dataStatus',
                title: '数据状态',
                template: 'partial/importPriview/dataStatus'
            }
        ],
        id: param.id
    })  
        
    },
    setupController: function(controller, model) {
        controller.set('model', model)
        controller.initUser(model)
    }
    
})

