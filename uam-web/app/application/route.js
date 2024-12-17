import Ember from 'ember';
const {set,
    get,
    inject,
    RSVP
} = Ember;
export default Ember.Route.extend({
    apiSrv: inject.service('api/portal-api/service'),
    session: inject.service('auth-service/service'),
    beforeModel: function() {
        if (window.location.hash === '') {
            this.transitionTo('dashboard.portalManagement.user');
        }
    },
    model: function() {
        let promises = {};
        let tenants = this.get('apiSrv').queryTenantList().then(results => results.content);
        let user = this.get("session").authInfo().user;

        promises.tenants = tenants;
        promises.user = user;

        this.get("apiSrv").isMultiTenant().then(resp => {
            window.localStorage.setItem("IS_MULTI_TENANT", resp.content==null?'true':resp.content);
        });

        return Ember.RSVP.hash(promises).then(results => results);

    },
    setupController: function(controller, model) {
        var tenant = window.sessionStorage.getItem("tenant");
        //storage里没有取到租户信息，则从当前用户的登录信息里获取
        var authInfo = this.get("session").authInfo();
        model.authInfo = authInfo;
        if (Ember.isEmpty(tenant)) {
            if (authInfo.authenticated && authInfo.user) {
                tenant = authInfo.user.tenantId;
            }
        }
        model.tenant = tenant;
        controller.set("model", model);

        let ifTenant= window.localStorage.getItem("IS_MULTI_TENANT");
        model.ifTenant = (ifTenant=='true');
    },
    actions: {
        loading(transition) {
            //displayLoadingSpinner();
            this.router.one('didTransition', function() {
                // hideLoadingSpinner();
                // console.log('loading done');
            });
            // substate implementation when returning `true`
            return true;
        },
        error(error, transition) {
            // Ember.onerror(error);
            if (error && error.status === 400) {
                return this.transitionTo('error');
            }
            return true;
        }
    }
});