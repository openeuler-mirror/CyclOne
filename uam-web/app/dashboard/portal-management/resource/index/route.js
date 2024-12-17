import Ember from 'ember';

export default Ember.Route.extend({
    model: function() {
        var columns = [{
            "propertyName": "name",
            "title": "名称",
            "template": "partial/lookUp"
        }, {
            "propertyName": "code",
            "title": "编码"
        }, {
            "propertyName": "url",
            "title": "URL"
        }, {
            "propertyName": "appId",
            "title": "系统名称"
        }, {
            "title": "操作",
            "template": "partial/resOper"
        }];
        return Ember.RSVP.hash({
            columns: columns,
            indexNumber: 1
        });
    },
    setupController: function(controller, model) {
        controller.initResource(model);
        controller.set('model', model);
    }
});