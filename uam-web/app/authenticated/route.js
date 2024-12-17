import Ember from 'ember';

export default Ember.Route.extend({

    queryParams: {
        token: {
            refreshModel: true
        },
        redirectTo: {
            refreshModel: true
        }
    },

    model: function(params) {
        return params;
    },

    setupController: function(controller, model) {
        controller.set('model', model);
        controller.initController();
    }
});