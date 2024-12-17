import Ember from 'ember';

export default Ember.Route.extend({
    model: function(param) {
        let groupColumns = [{ "propertyName": "name", "title": "名称" }, { "propertyName": "remark", "title": "备注" }, { "propertyName": "gmtCreate", "title": "创建时间" }];
        let roleColumns = [{ "propertyName": "code", "title": "编码" }, { "propertyName": "name", "title": "名称" }, { "propertyName": "remark", "title": "备注" }, { "propertyName": "gmtCreate", "title": "创建时间" }];

        Ember.set(param, 'groupColumns', groupColumns);
        Ember.set(param, 'roleColumns', roleColumns);
        Ember.set(param, 'indexNumber', 1);
        return param;
    },
    setupController: function(controller, model) {
        controller.set('model', model);
        controller.initUserView(model);
    }

});
