import Ember from 'ember';
const {set,
    get,
    inject,
    RSVP
} = Ember;

export default Ember.Controller.extend({
    tenantChange: function() {
        var tenant = this.get("model.tenant");
        var oldTenant = window.sessionStorage.getItem("tenant");
        if (tenant !== oldTenant) {
            window.sessionStorage.setItem("tenant", tenant);
            window.location.reload();
        }
        // set(model, 'selPerLength', model.selPers.length);
    }.observes('model.tenant'),
    actions: {
        menuChangeAct: function(url) {
            this.transitionToRoute(url);
        },

        logout: function() {
            if (window.sessionStorage.getItem('tenant')) {
                window.sessionStorage.removeItem("tenant");
            }
            this.set('curUser', null);
            document.cookie = "LtpaToken2="  + '=0;expires=' + new Date(0).toUTCString()+"; path=/"+"; domain=.gf.com.cn";
            document.cookie = "access-token=" + '=; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
            location.href = "/";
        }
    },

});