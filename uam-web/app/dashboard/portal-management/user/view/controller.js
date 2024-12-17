import Ember from 'ember';
const {get,
    set,
    inject
} = Ember;
const pageSize = 10;
export default Ember.Controller.extend({
    userSrv: inject.service('api/portal-user/service'),
    apiSrv: inject.service('api/portal-api/service'),
    groupSrv: inject.service('api/portal-user-group/service'),
    roleSrv: inject.service('api/portal-role/service'),
    userSteps: [{
        title: "用户信息",
        status: "process",
        active: true,
        index: 0,
        icon: "usermanager"
    }, {
        title: "用户组信息",
        status: "wait",
        index: 1,
        icon: "usergroup"
    }, {
        title: "角色信息",
        status: "wait",
        index: 2,
        icon: "rolemanager"
    }, {
        title: "权限信息",
        status: "wait",
        index: 3,
        icon: "permissiondistribution",
        last: true
    }],
    loading: false,

    initUserView: function(model) {
        let userId = model.id;
        if (Ember.isBlank(userId)) {
            swal("用户id不能为空，请确认！");
            return;
        }
        this.send('stepAction', 0);
        this.queryUserById(userId);
    },

    queryUserById: function(id) {
        let apiSrv = get(this, 'apiSrv'),
            self = this;

        apiSrv.queryUserById(id).then(res => {
            if ("success" === res.status) {
                set(self, 'userForm', res.item);
            } else {
                swal("查询用户信息失败", res.message);
            }
        });
    },

    queryPermissionsById: function(id) {
        let userSrv = get(this, 'userSrv'),
            model = get(this, 'model');

        userSrv.queryPermissionsById(id).then(res => {
            if ("success" === res.status) {
                set(model, 'permissions', res.content.links.PERMISSION.sortBy("systemName"));
            } else {
                swal("查询用户关联的权限信息失败，", res.message);
            }
        });
    },

    queryRolePageByUserIdAndCnd: function(cnd, pageNo) {
        let roleSrv = get(this, 'roleSrv'),
            model = get(this, 'model');

        let queryForm = {
            id: model.id,
            cnd: cnd,
            pageNo: pageNo,
            pageSize: pageSize
        };

        set(this, 'loading', true);
        roleSrv.queryRolePageByUserIdAndCnd(queryForm).then(res => {
            set(this, 'loading', false);
            if ("success" === res.status) {
                set(model, 'roles', res.list);
                set(model, 'rolemeta', res.meta);
                set(model, 'indexNumber', (pageNo - 1) * pageSize + 1);
            } else {
                swal("查询用户关联的角色信息失败，", res.message);
            }
        });
    },

    queryGroupPageByUserIdAndCnd: function(cnd, pageNo) {
        let groupSrv = get(this, 'groupSrv'),
            model = get(this, 'model');

        let queryForm = {
            id: model.id,
            cnd: cnd,
            pageNo: pageNo,
            pageSize: pageSize
        };

        set(this, 'loading', true);
        groupSrv.queryGroupPageByUserIdAndCnd(queryForm).then(res => {
            set(this, 'loading', false);
            if ("success" === res.status) {
                set(model, 'groups', res.list);
                set(model, 'groupmeta', res.meta);
                set(model, 'indexNumber', (pageNo - 1) * pageSize + 1);
            } else {
                swal("查询用户关联的用户组失败，", res.message);
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
            let userSteps = get(this, 'userSteps'),
                model = get(this, 'model');
            let currentStep = userSteps.findBy("index", index);

            userSteps.forEach(item => {
                if (item.index !== index) {
                    set(item, 'status', 'wait');
                    return;
                }
                set(currentStep, 'status', 'process');
            });
            set(this, 'currentTab', currentStep.title);

            if (index === 0) {
                this.queryUserById(model.id);
            }

            if (index === 1) {
                this.queryGroupPageByUserIdAndCnd("", 1);
            }

            if (index === 2) {
                this.queryRolePageByUserIdAndCnd("", 1);
            }

            if (index === 3) {
                this.queryPermissionsById(model.id);
            }
        },

        groupPageClick: function(pageNo) {
            this.queryGroupPageByUserIdAndCnd("", pageNo);
        },
        rolePageClick: function(pageNo) {
            this.queryRolePageByUserIdAndCnd("", pageNo);
        }
    }
});