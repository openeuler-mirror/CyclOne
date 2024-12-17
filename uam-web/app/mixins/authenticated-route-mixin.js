import Ember from 'ember';

const {
    service
} = Ember.inject;

export default Ember.Mixin.create({
    session: service('auth-service/service'),
    /**
     * [beforeModel description]
     * @param  {[type]} transition [description]
     * @return {[type]}            [description]
     */
    beforeModel(transition) {
        var authInfo = this.get("session").authInfo();
        var user = authInfo.user;

        if (!authInfo.authenticated) {
            transition.abort();
            this.set('session.attemptedTransition', transition);
            Ember.run.later(() => {
                location = authInfo.login;
            }, 10);
        } else {
            // 如果没有login的权限，不允许进入
            let permissions = window.sessionStorage.getItem("permissions").split(',');
            if (authInfo.user.loginName !== "admin" && !permissions.includes("login")) {
                this.transitionTo('error', "用户无登陆权限，请确认！");
                Ember.run.later(() => {
                    document.cookie = "access-token=" + '=; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
                    location.href = "/";
                }, 3000);
            }
            return this._super(...arguments);
        }
        return null;
    }
});