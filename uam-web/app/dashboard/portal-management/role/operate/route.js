import Ember from 'ember';
var hash = Ember.RSVP.hash;
export default Ember.Route.extend({
    queryParams: {
        id: {
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
            roleId: params.id,
            type: params.type,
            source: params.source,
            selectedUserList: [],
            selectedGroupList: [],
            bakGroupList: [],
            bakUserList: [],
            dataSource: []
        });
    },
    setupController: function(controller, model) {
        controller.set('model', model);
        controller.initOperateRole(model);
    }
});
