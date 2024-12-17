import Ember from 'ember';

export default Ember.Route.extend({
    model: function() {
        return Ember.RSVP.hash({});
    },
    setupController: function(controller, model) {
        controller.initPermission(model);
    }
});
