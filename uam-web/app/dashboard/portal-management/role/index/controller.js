import Ember from 'ember';
const {set, get, inject } = Ember;
const pageSize = 10;
export default Ember.Controller.extend({
    portalRoleSrv: inject.service('api/portal-role/service'),
    queryForm: {
        name: "",
    },

    /**
     * 查询form里面name的监听事件
     */
    queryByName: function() {
        this.initRole();
    }.observes('queryForm.name'),

    queryRoleByPage: function(offset, queryForm) {
        var portalRoleSrv = get(this, 'portalRoleSrv'),
            model = get(this, 'model');

        set(model, 'paginationShow', false);
        this.set('loading', true);
        portalRoleSrv.queryByPage(offset, pageSize, queryForm).then((data) => {
            this.set('loading', false);
            set(model, 'data', data.list);
            set(model, 'meta', data.meta);
            if (data.meta.pages > 1) {
                set(model, 'paginationShow', true);
            }
        });

    },

    /**
     * 初始化role主页面
     * @param  {[type]} currentPage [description]
     * @return {[type]}             [description]
     */
    initRole: function(model) {
        set(model, 'indexNumber', 1);
        set(model, 'pageNo', 1);
        set(model, 'paginationShow', false);
        set(this, 'model', model);
        this.queryRoleByPage(1, this.queryForm);
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
        pageClick: function(pageNo) {
            var model = get(this, 'model');
            set(model, 'indexNumber', (pageNo - 1) * model.meta.limit + 1);
            set(model.meta, 'offset', pageNo);
            this.queryRoleByPage(pageNo, this.queryForm);
        },
        /**
         * 查询角色信息
         * @return {[type]} [description]
         */
        queryRoleDataAction: function() {
            var model = get(this, 'model');
            this.initRole(model);
        },

        /**
         * 对角色进行操作包括新增，删除，分配用户，分配用户组等
         * @param  {[type]} actionType [description]
         * @param  {[type]} role       [description]
         * @return {[type]}            [description]
         */
        operateRoleAction: function(actionType, role) {
            var id = "";
            if (!Ember.isBlank(role)) {
                id = Ember.isBlank(role.id) ? "" : role.id;
            }

            var param = {
                id: id,
                type: actionType
            };

            this.transitionToRoute('dashboard.portalManagement.role.operate', {
                queryParams: param
            });
        },

        lookUpAction: function(record) {
            if (Ember.isBlank(record.id)) {
                swal("角色id为空，请确认");
                return;
            }

            this.transitionToRoute('dashboard.portalManagement.role.view', record.id);
        },

        /**
         * 删除角色信息，同时刷新角色列表
         * @param  {[type]} role [description]
         * @return {[type]}      [description]
         */
        deleteRoleAction: function(role) {
            var roleId = role.id,
                model = get(this, 'model'),
                self = this,
                portalRoleSrv = get(this, 'portalRoleSrv');

            if (Ember.isBlank(roleId)) {
                swal("角色信息为空，请确认");
                return;
            }
            swal({
                title: "是否删除此角色?",
                type: "warning",
                showCancelButton: true,
                confirmButtonClass: "btn-danger",
                cancelButtonText: "取消",
                confirmButtonText: "删除",
                closeOnConfirm: false
            }, function(isConfirm) {
                if (isConfirm) {
                    portalRoleSrv.delete(roleId).then(function(data) {
                        swal.close();
                        if ("success" === data.status) {
                            self.queryRoleByPage(model.meta.offset, self.queryForm);
                        }
                    });
                }
            });
        }
    }
});
