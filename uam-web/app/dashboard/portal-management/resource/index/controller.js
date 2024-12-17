import Ember from 'ember';
const {set, get, inject } = Ember;
export default Ember.Controller.extend({
    portalResSrv: inject.service('api/portal-resource/service'),

    /**
     * 初始化resInfo主页面
     * @param  {[type]} currentPage [description]
     * @return {[type]}             [description]
     */
    initResource: function(model) {
        var portalResSrv = get(this, 'portalResSrv');
        this.set('loading', true);
        portalResSrv.queryAll().then((res) => {
            this.set('loading', false);
            set(model, 'data', res.list);
        });
    },
    /**
     * [loading description]
     * @type {Boolean}
     */
    loading: false,
    /**
     * [filterString description]
     * @type {String}
     */
    filterString: '',

    /**
     * 事件集合
     * @type {Object}
     */
    actions: {

        showModelAction: function(resInfo) {
            let id = "";
            if (resInfo) {
                id = resInfo.code;
            }
            this.transitionToRoute('dashboard.portalManagement.resource.operate', id);
        },

        lookUpAction: function(resInfo) {
            if (Ember.isBlank(resInfo.code)) {
                swal("权限资源编码为空，请确认");
                return;
            }
            this.transitionToRoute('dashboard.portalManagement.resource.view', resInfo.code);
        },

        /**
         * 对角色进行操作包括新增，删除，分配用户，分配用户组等
         * @param  {[type]} actionType [description]
         * @param  {[type]} resInfo       [description]
         * @return {[type]}            [description]
         */
        saveResAction: function() {
            var model = get(this, 'model'),
                resForm = get(this, 'resForm'),
                portalResSrv = get(this, 'portalResSrv'),
                self = this;

            if (resForm.id) {
                portalResSrv.update(resForm.id, resForm).then((res) => {
                    if (res.status === "success") {
                        swal("修改成功!");
                        self.send('toggleModal', 'AddShowing');
                        self.initResource(model);
                    } else {
                        swal(res.message);
                    }
                });
            } else {
                portalResSrv.create(resForm).then((res) => {
                    if (res.status === "success") {
                        swal("保存成功!");
                        self.send('toggleModal', 'AddShowing');
                        self.initResource(model);
                    } else {
                        swal(res.message);
                    }
                });
            }
        },

        /**
         * 删除角色信息，同时刷新角色列表
         * @param  {[type]} resInfo [description]
         * @return {[type]}      [description]
         */
        deleteResAction: function(resInfo) {
            var id = resInfo.id,
                model = get(this, 'model'),
                self = this,
                portalResSrv = get(this, 'portalResSrv');

            if (Ember.isBlank(id)) {
                swal("权限资源编码不能为空，请确认！");
                return;
            }
            swal({
                title: "是否删除此权限资源?",
                type: "warning",
                showCancelButton: true,
                confirmButtonClass: "btn-danger",
                cancelButtonText: "取消",
                confirmButtonText: "删除",
                closeOnConfirm: false
            }, function(isConfirm) {
                if (isConfirm) {
                    portalResSrv.delete(id).then(function(data) {
                        if ("success" === data.status) {
                            swal.close();
                            self.initResource(model);
                        } else {
                            swal(data.message);
                        }
                    });
                }
            });
        }
    }
});
