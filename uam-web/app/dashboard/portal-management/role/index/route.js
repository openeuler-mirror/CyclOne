import Ember from 'ember';

export default Ember.Route.extend({
    model: function() {
        var columns = [{ "propertyName": "name", "title": "名称","template":"partial/lookUp" }, { "propertyName": "code", "title": "编码" }, { "title": "操作", "template": "partial/roleOper" }];
        return { 'columns': columns };
    },
    setupController: function(controller, model) {
        controller.initRole(model);
        controller.set('model', model);
    }
});
