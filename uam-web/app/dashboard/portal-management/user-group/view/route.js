import Ember from 'ember';
const {get,
    set,
    inject
} = Ember;
export default Ember.Route.extend({
    model: function(param) {
        let userColumns = [{
            "propertyName": "name",
            "title": "用户名称"
        }, {
            "propertyName": "tenantId",
            "title": "租户"
        }, {
            "propertyName": "title",
            "title": "职位"
        }, {
            "propertyName": "deptName",
            "title": "部门名称"
        }];
        let roleColumns = [{
            "propertyName": "code",
            "title": "编码"
        }, {
            "propertyName": "name",
            "title": "名称"
        }, {
            "propertyName": "remark",
            "title": "备注"
        }, {
            "propertyName": "gmtCreate",
            "title": "创建时间"
        }];

        Ember.set(param, 'userColumns', userColumns);
        Ember.set(param, 'roleColumns', roleColumns);
        Ember.set(param, 'indexNumber', 1);
        return param;
    },
    setupController: function(controller, model) {
        controller.set('model', model);
        controller.initGroupView(model);
    }

});