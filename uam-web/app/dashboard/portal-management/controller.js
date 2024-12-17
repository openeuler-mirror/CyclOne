import Ember from 'ember';
var set = Ember.set;

export default Ember.Controller.extend({
    currentItem: "",

    initPortalMenu: function() {
        // set(this, 'currentItem', "user");
    },

    /**
     * 事件集合
     * @type {Object}
     */
    actions: {
        clickItemAction: function(item) {
            set(this, 'currentItem', item);
            this.transitionToRoute("dashboard.portalManagement." + item);
        }
    }
});
