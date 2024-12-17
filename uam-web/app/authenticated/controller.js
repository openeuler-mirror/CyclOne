import Ember from 'ember';

export default Ember.Controller.extend({
    queryParams: ["token", "redirectTo"],

    initController: function() {
        var model = this.get("model");
        var token = model.token;
        var redirectTo = model.redirectTo;

        var expireDate = new Date();
        expireDate.setTime(expireDate.getTime() + (18 * 60 * 60 * 1000));

        document.cookie = "access-token=" + token + ";expires=" + expireDate.toUTCString();
        if(redirectTo == null || redirectTo == "") {
            location.href = "/";
        } else {
            location.href = redirectTo;
        }
    }
});
