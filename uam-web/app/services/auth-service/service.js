import Ember from 'ember';

export default Ember.Service.extend({
    authInfo: function() {
        var model = {
            authenticated: false,
        };

        if (this.get("curUser")) {
            model.authenticated = true;
            model.user = this.get("curUser");
            return model;
        }

        var token = null;
        if (document.cookie !== undefined) {
            document.cookie.split('; ').forEach(function(data) {
                if (data.indexOf("access-token=") === 0) {
                    token = data.substring(13);
                }
            });
        }

        var url = "/auth";
        var self = this;
        let callBackUrl=window.location.protocol+"//"+window.location.host + window.location.pathname + "#/authenticated?loginHash=" + window.location.hash.replace("#/","");
        let customer="default";
        window.sessionStorage.removeItem("permissions");
        Ember.$.ajax({
            'method': 'GET',
            'url': url,
            headers: {
                "access-token": "Bearer " + token
            },
            async: false,
            success: function(data) {
                if (data.status === "success") {
                    model.user = data.content;
                    model.authenticated = true;
                    self.set("curUser", data.content);
                    window.sessionStorage.setItem("permissions", data.content.authorizationInfo.permissions.CLOUDUAM_MENU);
                } else if (data.status === "AUTH_FAILED") {
                    model.login = "/login?authCallbackUrl=" + encodeURIComponent(callBackUrl)+"&customer="+customer;
                } else {
                    swal(data.message);
                }
            },
            error: function(resp) {
                if(resp.status === 401){
                    model.login = resp.responseJSON.content.ssoWebUrl + encodeURIComponent(callBackUrl)+"&customer="+customer;
                }
                // swal("系统异常");
            }
        });

        return model;
    },
});