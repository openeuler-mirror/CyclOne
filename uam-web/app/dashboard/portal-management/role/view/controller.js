import Ember from 'ember';
const {get,
    set,
    inject
} = Ember;
const pageSize = 10;
export default Ember.Controller.extend({
    roleSrv: inject.service('api/portal-role/service'),
    groupSrv: inject.service('api/portal-user-group/service'),
    userSrv: inject.service('api/portal-user/service'),
    roleSteps: [{
        title: "角色信息",
        status: "process",
        active: true,
        index: 0,
        icon: "rolemanager"
    }, {
        title: "用户组信息",
        status: "wait",
        index: 1,
        icon: "usergroup"
    }, {
        title: "用户信息",
        status: "wait",
        index: 2,
        icon: "usermanager"
    }, {
        title: "权限信息",
        status: "wait",
        index: 3,
        icon: "permissiondistribution",
        last: true
    }],
    loading: false,

    initRoleView: function(model) {
        let groupId = model.id;
        if (Ember.isBlank(groupId)) {
            swal("角色id不能为空，请确认！");
            return;
        }
        this.send('stepAction', 0);
        this.queryById(groupId);
    },

    queryUserPageByRoleIdAndCnd: function(cnd, pageNo) {
        let userSrv = get(this, 'userSrv'),
            model = get(this, 'model');

        let queryForm = {
            id: model.id,
            cnd: cnd,
            pageNo: pageNo,
            pageSize: pageSize
        };

        set(this, 'loading', true);
        userSrv.queryUserPageByRoleIdAndCnd(queryForm).then(res => {
            set(this, 'loading', false);
            if ("success" === res.status) {
                set(model, 'users', res.list);
                set(model, 'usermeta', res.meta);
                set(model, 'indexNumber', (pageNo - 1) * pageSize + 1);
            } else {
                swal("查询角色关联用户信息失败", res.message);
            }
        });
    },

    queryPermissionsById: function(id) {
        let roleSrv = get(this, 'roleSrv'),
            model = get(this, 'model');

        roleSrv.queryPermissionsById(id).then(res => {
            if ("success" === res.status) {
                set(model, 'permissions', res.item.links.PERMISSION.sortBy("systemName"));
            } else {
                swal("查询角色关联的权限信息失败，", res.message);
            }
        });
    },

    queryGroupPageByRoleIdAndCnd: function(cnd, pageNo) {
        let groupSrv = get(this, 'groupSrv'),
            model = get(this, 'model');

        let queryForm = {
            id: model.id,
            cnd: cnd,
            pageNo: pageNo,
            pageSize: pageSize
        };
        set(this, 'loading', true);
        groupSrv.queryGroupPageByRoleIdAndCnd(queryForm).then(res => {
            set(this, 'loading', false);
            if ("success" === res.status) {
                set(model, 'groups', res.list);
                set(model, 'groupmeta', res.meta);
            } else {
                swal("查询角色关联的用户组信息失败，", res.message);
            }
        });
    },

    queryById: function(id) {
        let roleSrv = get(this, 'roleSrv'),
            self = this;

        roleSrv.queryById(id).then(res => {
            if ("success" === res.status) {
                set(self, 'roleForm', res.item);
            } else {
                swal("查询角色信息失败，", res.message);
            }
        });
    },

    actions: {
        /**
         * 点击步骤操作
         * @param  {[type]} index [description]
         * @return {[type]}       [description]
         */
        stepAction: function(index) {
            let roleSteps = get(this, 'roleSteps'),
                model = get(this, 'model');

            let currentStep = roleSteps.findBy("index", index);

            roleSteps.forEach(item => {
                if (item.index !== index) {
                    set(item, 'status', 'wait');
                    return;
                }
                set(currentStep, 'status', 'process');
            });
            set(this, 'currentTab', currentStep.title);


            if (index === 0) {
                this.queryById(model.id);
            }

            if (index === 1) {
                this.queryGroupPageByRoleIdAndCnd("", 1);
            }

            if (index === 2) {
                this.queryUserPageByRoleIdAndCnd("", 1);
            }

            if (index === 3) {
                this.queryPermissionsById(model.id);
            }
        },

        userPageClick: function(pageNo) {
            this.queryUserPageByRoleIdAndCnd("", pageNo);

        },
        groupPageClick: function(pageNo) {
            this.queryGroupPageByRoleIdAndCnd("", pageNo);
        }
    }


});