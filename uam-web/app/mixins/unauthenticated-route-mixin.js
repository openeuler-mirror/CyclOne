import Ember from 'ember';

const { service } = Ember.inject;

export default Ember.Mixin.create({
    authSrv: service('auth-service/service'),
    /**
     * [beforeModel description]
     * @param  {[type]} transition [description]
     * @return {[type]}            [description]
     */
    beforeModel(transition) {
        var user = this.get("authSrv").getUser();
        if (user) {
            // 如果没有RBAC的权限，不允许进入
            let permissions = authInfo.user.authorizationInfo.permissions;

            if (!permissions.RBAC) {
                this.transitionToRoute('error', "用户没有系统登陆权限，请确认");
                return;
            }
            transition.abort();
            this.transitionTo('dashboard.srvManagement.srvBpAppInfo.bpList');
        } else {
            return this._super(...arguments);
        }
    }
});
