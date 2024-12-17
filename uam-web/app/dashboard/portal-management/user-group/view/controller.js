import Ember from 'ember';
const {get,
    set,
    inject
} = Ember;
const pageSize = 10;
export default Ember.Controller.extend({
    groupSrv: inject.service('api/portal-user-group/service'),
    roleSrv: inject.service('api/portal-role/service'),
    apiSrv: inject.service('api/portal-api/service'),
    userSrv: inject.service('api/portal-user/service'),
    userGroupSteps: [{
        title: "用户组信息",
        status: "process",
        active: true,
        index: 0,
        icon: "usergroup"
    }, {
        title: "角色信息",
        status: "wait",
        index: 1,
        icon: "rolemanager"
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

    initGroupView: function(model) {
        let groupId = model.id;
        if (Ember.isBlank(groupId)) {
            swal("用户组id不能为空，请确认！");
            return;
        }

        let initmeta = {
            pages: 0,
            offset: 1
        };
        set(model, 'usermeta', initmeta);
        this.send('stepAction', 0);
        this.queryById(groupId);
    },

    queryUserPageByGroupIdAndCnd: function(cnd, pageNo) {
        let userSrv = get(this, 'userSrv'),
            model = get(this, 'model');

        let queryForm = {
            id: model.id,
            cnd: cnd,
            pageNo: pageNo,
            pageSize: pageSize
        };

        set(this, 'loading', true);
        userSrv.queryUserPageByGroupIdAndCnd(queryForm).then(res => {
            set(this, 'loading', false);
            if ("success" === res.status) {
                set(model, 'users', res.list);
                set(model, 'usermeta', res.meta);
                set(model, 'indexNumber', (pageNo - 1) * pageSize + 1);
            } else {
                swal("查询用户信息失败", res.message);
            }
        });
    },

    queryPermissionsById: function(id) {
        let groupSrv = get(this, 'groupSrv'),
            model = get(this, 'model');

        groupSrv.queryPermissionsById(id).then(res => {
            if ("success" === res.status) {
                set(model, 'permissions', res.item.links.PERMISSION.sortBy("systemName"));
            } else {
                swal("查询用户组关联的权限信息失败，", res.message);
            }
        });
    },

    queryRolePageByGroupIdAndCnd: function(cnd, pageNo) {
        let roleSrv = get(this, 'roleSrv'),
            model = get(this, 'model');

        let queryForm = {
            id: model.id,
            cnd: cnd,
            pageNo: pageNo,
            pageSize: pageSize
        };

        set(this, 'loading', true);
        roleSrv.queryRolePageByGroupIdAndCnd(queryForm).then(res => {
            set(this, 'loading', false);
            if ("success" === res.status) {
                set(model, 'roles', res.list);
                set(model, 'rolemeta', res.meta);
                set(model, 'indexNumber', (pageNo - 1) * pageSize + 1);
            } else {
                swal("查询用户组关联的角色信息失败，", res.message);
            }
        });
    },

    queryById: function(id) {
        let groupSrv = get(this, 'groupSrv'),
            self = this;

        groupSrv.queryById(id).then(res => {
            if ("success" === res.status) {
                set(self, 'groupForm', res.item);
            } else {
                swal("查询用户组关联的用户组失败，", res.message);
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
            let userGroupSteps = get(this, 'userGroupSteps'),
                model = get(this, 'model');

            let currentStep = userGroupSteps.findBy("index", index);

            userGroupSteps.forEach(item => {
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
                this.queryRolePageByGroupIdAndCnd("", 1);
            }

            if (index === 2) {
                this.queryUserPageByGroupIdAndCnd("", 1);
            }

            if (index === 3) {
                this.queryPermissionsById(model.id);
            }
        },

        userPageClick: function(pageNo) {
            this.queryUserPageByGroupIdAndCnd("", pageNo);
        },
        rolePageClick: function(pageNo) {
            this.queryRolePageByGroupIdAndCnd("", pageNo);
        }
    }


});