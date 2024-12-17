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
            columns: [],
            userGroupId: params.id,
            type: params.type,
            source: params.source,
            dataSource: [],
            AllChecked: false
        });
    },
    setupController: function(controller, model) {
        let ifTenant= window.localStorage.getItem("IS_MULTI_TENANT");
        model.ifTenant = (ifTenant=='true');
        controller.set('model', model);
        controller.initUserGroup(model);
    }
});
