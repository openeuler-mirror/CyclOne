import Ember from 'ember';
const {set, get, inject } = Ember;
export default Ember.Controller.extend({
    resSrv: inject.service('api/portal-resource/service'),


    initPerRes: function(model) {
        var resId = model.id,
            resSrv = get(this, 'resSrv'),
            self = this;

        //若id为空，说明为新增，初始化所用的用户组和角色信息
        if (Ember.isBlank(resId)) {
            set(this, 'resForm', {
                name: "",
                code: "",
                dataSource: "",
                remark: "",
            });
        } else {
            resSrv.queryByCode(resId).then(function(data) {
                if (data.status === "success") {
                    set(self, 'resForm', data.item);
                } else {
                    swal(data.message);
                }
            });
        }
    },

    actions: {
        /**
         * 返回链接
         * @return {[type]} [description]
         */
        backPerResAction: function() {
            this.transitionToRoute('dashboard.portalManagement.resource.index');
        },

        /**
         * 保存用户操作
         * @param  {[type]} resForm [description]
         * @return {[type]}          [description]
         */
        saveResAction: function() {
            var resSrv = get(this, 'resSrv'),
                resForm = get(this, 'resForm'),
                self = this,
                resId = resForm.id;

            if (Ember.isBlank(resId)) {
                resSrv.create(resForm).then(function(data) {
                    if (data.status === "success") {
                        swal("新建权限资源" + resForm.name + "成功！");
                        self.send('backPerResAction');
                    } else {
                        swal(data.message);
                    }
                });
            } else {
                resSrv.update(resId, resForm).then(function(data) {
                    if (data.status === "success") {
                        swal("修改权限资源" + resForm.name + "成功！");
                        self.send('backPerResAction');
                    } else {
                        swal(dta.message);
                    }
                });
            }
        }
    }
});
