import Ember from 'ember';
import AuthenticatedRouteMixin from 'clouduam-web/mixins/authenticated-route-mixin';
export default Ember.Route.extend(AuthenticatedRouteMixin, {
	model: function(params) {
		var authInfo = this.get("session").authInfo();
		
		var user = authInfo.user;
		var tenantId = user.tenantId;
		let promises = {};
		var tenant = window.sessionStorage.getItem("tenant");
		if (tenant != null && tenant != '' && tenant != 'default') {
			tenantId = tenant;
		}
		user.tenant = tenantId;
		promises.user = user;
		promises.permissions =window.sessionStorage.getItem("permissions").split(',');
		return Ember.RSVP.hash(promises).then(results => results);
	},

	setupController: function(controller, model) {
		let ifTenant= window.localStorage.getItem("IS_MULTI_TENANT");
        model.ifTenant = (ifTenant=='true');

		controller.set('model', model);

	}
});