import Ember from 'ember';
var hash = Ember.RSVP.hash;
export default Ember.Route.extend({
    queryParams: {
        deptID: {
            refreshModel: true
        },
        type: {
            refreshModel: true
        },
        source: {
            refreshModel: true
        }
    },
    model: function(params) {
        return hash({
            id: params.id,
            deptID: params.deptID,
            type: params.type,
            groupSource: [],
            selGroups: []
        });
    },
    setupController: function(controller, model) {
        controller.set('model', model);
        controller.initOperateUser(model);
    }
});
