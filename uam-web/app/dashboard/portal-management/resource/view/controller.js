import Ember from 'ember';
const {set, get, inject } = Ember;
export default Ember.Controller.extend({
    resSrv: inject.service('api/portal-resource/service'),

    initPerRes: function(model) {
        let resSrv = get(this, 'resSrv'),
            self = this;

        if (Ember.isBlank(model.id)) {
            swal("权限资源id为空，请确认");
            return;
        }

        resSrv.queryByCode(model.id).then(res => {
            if ("success" === res.status) {
                set(self, 'resForm', res.item);
                set(model, 'readOnly', true);
                set(self, 'model', model);
            } else {
                swal("查询权限资源信息失败", res.message);
            }
        });
    }
});
