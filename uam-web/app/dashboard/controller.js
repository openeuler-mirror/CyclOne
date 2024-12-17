import Ember from 'ember';

export default Ember.Controller.extend({
    activeMenu: "",
    actions: {
        linkToAction: function(item) {
            this.transitionToRoute('dashboard.portalManagement.' + item);
            Ember.set(this, 'activeMenu', item);
        }
    }
});
