import Ember from 'ember';
const {get,
    set,
    inject,
    RSVP
} = Ember;
const pageSize = 20;
const rootNode = {
    id: "",
    title: "所有部门"
};
export default Ember.Controller.extend({
    queryParams: ["id", "type", "source"],
    portalUserGroupSrv: inject.service('api/portal-user-group/service'),
    portalCommonSrv: inject.service('api/portal-common/service'),
    apiSrv: inject.service('api/portal-api/service'),

    userGroupSteps: [{
        title: "用户组信息",
        status: "process",
        active: true,
        index: 1,
        icon: "usergroup"
    }, {
        title: "角色清单",
        status: "wait",
        index: 2,
        icon: "rolemanager"
    }, {
        title: "用户清单",
        status: "wait",
        index: 3,
        icon: "usermanager",
        last: true
    }],
    currentTab: "用户组信息",
    /**
     * 比较两个list，将删除，新增标志与id号进行拼接
     * @param  {[type]} backList     [description]
     * @param  {[type]} selectedList [description]
     * @return {[type]}              [description]
     */
    compareList: function(backList, selectedList) {
        var idListUnion = []; //id号与操作的拼接数据 id:I id:D

        //将selectedList与backList对比
        //若两个list都存在相同的id，从backList当中移除，只保留未选择的
        //若selectedList当中的id在backList当中不存在，则标记为I
        for (var i = selectedList.length - 1; i >= 0; i--) {

            var selectedId = "";
            if (typeof(selectedList[i]) === "string") {
                selectedId = selectedList[i];
            } else {
                selectedId = selectedList[i].id;
            }

            if (backList.contains(selectedId)) {
                backList.removeAt(backList.indexOf(selectedId));
            } else {
                idListUnion.addObject(selectedId + ":I");
            }
        }

        //将backList里面所有的id都标记为D
        for (var j = backList.length - 1; j >= 0; j--) {
            idListUnion.addObject(backList[j] + ":D");
        }
        return idListUnion.join();
    },

    /**
     * 监听用户table的选择操作
     * @return {[type]} [description]
     */
    observeCheckbox: function() {
        let model = get(this, 'model');
        if (model.users) {
            model.users.forEach(item => {
                if (item.checked && !model.selUsers.findBy("id", item.id)) {
                    model.selUsers.addObject(item);
                }

                if (!item.checked && model.selUsers.findBy("id", item.id)) {
                    model.selUsers.removeObject(model.selUsers.findBy("id", item.id));
                }
            });
        }

        if (model.selUsers) {
            set(model, 'selUsersLength', model.selUsers.length);
        }
    }.observes('model.users.@each.checked'),

    /**
     * 监听用户查询条件
     * @return {[type]} [description]
     */
    observeQueryParam: function() {
        let model = get(this, 'model'),
            apiSrv = get(this, 'apiSrv');

        /////// 第一次加载页面是为undfined
        /// 若model.queryParam=""， 会查询两次，根据deptID查询一次，然后监听查询条件再查询一次
        if (model.queryParam === undefined) {
            return;
        }
        var tenantId = window.sessionStorage.getItem("tenant");
        apiSrv.queryUserByCnd(tenantId, model.queryParam, 1, pageSize).then(res => {
            if ("success" === res.status) {
                let userData = res.content.list;

                /////// 对已经选择的用户添加checked
                /////// model.selUsers：已经选择的用户列表
                userData.forEach(item => {
                    if (model.selUsers.findBy("id", item.id)) {
                        set(item, 'checked', true);
                    }
                });
                set(model, 'users', userData);
                set(model.meta, 'totalPage', res.content.totalPage);
                set(model.meta, 'page', res.content.page);
            } else {
                // swal("查询失败，", res.message);
            }
        });

    }.observes('model.queryParam'),


    //////// 将角色列表根据指定的key，进行转换
    parseTransferList: function(list) {
        list.forEach((item) => {
            set(item, 'label', item.name);
            set(item, 'value', item.id);
        });
        return list;
    },
    /**
     * 查询用户组所有类型
     * @return {[type]} [description]
     */
    queryUserGroupType: function(){
        let groupTypeSrv = get(this, 'apiSrv');
        let tenant=window.sessionStorage.getItem("tenant");
        let model = get(this, 'model')
        groupTypeSrv.getUserGroupType(tenant).then(resp => {
            if (resp.status === 'success') {
                set(model, 'groupTypeTree', resp.list);
            } else {
                swal(resp.message)
            }
        });
    },
    /**
     * 初始化用户组信息
     * @return {[type]} [description]
     */
    initUserGroup: function(model) {
        this.queryUserGroupType();
        var portalUserGroupSrv = get(this, 'portalUserGroupSrv'),
            portalCommonSrv = get(this, 'portalCommonSrv'),
            apiSrv = get(this, 'apiSrv'),
            userGroupId = model.userGroupId,
            type = model.type,
            userGroupAllTabs = get(this, 'userGroupAllTabs'),
            self = this;

        set(model, 'meta', {
            page: 1,
            limit: pageSize
        });
        var tenant = window.sessionStorage.getItem("tenant");
        //判断是否为修改，为空，则为新增
        if (Ember.isBlank(userGroupId)) {
            set(this, 'currentTab', "用户组信息");
            set(this, 'userGroupTabs', userGroupAllTabs);
            set(model, 'userGroupId', "");

            set(this, 'userGroupForm', {
                name: "",
                remark: ""
            });

            RSVP.hash({
                "allRoles": portalCommonSrv.getRoleList().then(function(data) {
                    return data.list;
                }),
                "allDept": apiSrv.queryDepartment(tenant).then(res => {
                    return res.content;
                })
            }).then(function(data) {
                const allDept = data.allDept;
                allDept.insertAt(0, rootNode); ///// 添加所有部门的节点
                set(model, 'dataSource', self.parseTransferList(data.allRoles));
                set(model, 'selectedRoleList', []);
                set(model, 'bakRoleList', []);
                set(model, 'bakUserList', []);
                set(model, 'selUsers', []);
                set(model, 'deptTree', allDept);
                set(self, 'model', model);
                self.queryUserByDeptId('', 1, tenant);
            });
        } else {

            RSVP.hash({
                "userGroup": portalUserGroupSrv.queryById(userGroupId).then(res => {
                    return res.item;
                }),
                "allRoles": portalCommonSrv.getRoleList().then((data) => {
                    return data.list;
                }),
                "allDept": apiSrv.queryDepartment(tenant).then(res => {
                    return res.content;
                }),
                "selUsersInfo": portalUserGroupSrv.queryUsersById(userGroupId).then(res => {
                    return res.item;
                })
            }).then((res) => {
                //组装用户组与角色的对应关系
                var allRoles = res.allRoles,
                    userGroup = res.userGroup,
                    db_selRoles = userGroup.selRoles, // 数据库已经选择的角色信息
                    db_selUsers = userGroup.selUsers, // 数据库已经选择的用户信息
                    allDept = res.allDept,
                    selUsersInfo = res.selUsersInfo.links.USER;

                //      组装穿梭框数据---已经选择的角色及未选择的角色     //
                allRoles.forEach((item) => {
                    if (db_selRoles.contains(item.id)) {
                        set(item, 'chose', item.id);
                    }
                });

                allDept.insertAt(0, rootNode);
                set(model, 'selectedRoleList', db_selRoles);
                set(model, 'dataSource', self.parseTransferList(allRoles));
                set(model, 'selUsers', selUsersInfo);
                set(model, 'bakRoleList', Ember.copy(db_selRoles)); ///// 备份数据库已经选择的角色信息
                set(model, 'bakUserList', Ember.copy(db_selUsers)); ///// 备份数据库已经选择的用户信息
                set(model, 'deptTree', allDept);
                set(model, 'deptID', '');
                set(model, 'queryParam', '');
                set(self, 'model', model);
                set(self, 'userGroupForm', res.userGroup);
                self.queryUserByDeptId('', 1, tenant);
            });
        }

        //为ALL，全部显示，此时控制上一步，下一步按钮的显示
        if (type === "ALL" || type === "INFO") {
            this.send('stepAction', 1);
        }
        if (type === "ROLE") {
            this.send('stepAction', 2);
        }
        if (type === "USER") {
            this.send('stepAction', 3);
        }
    },

    /**
     * 根据部门id查询用户信息，并关联已经选择的用户与未选择的用户
     * 根据部门id,加载用户信息
     * 若用户组id为空，则不加载关联信息
     * @param  {[type]} deptId   [description]
     * @param  {[type]} pageNo   [description]
     * @param  {[type]} pageSize [description]
     * @return {[type]}          [description]
     */
    queryUserByDeptId: function(deptId, pageNo) {
        var apiSrv = get(this, 'apiSrv'),
            model = get(this, 'model'),
            portalUserGroupSrv = get(this, 'portalUserGroupSrv'),
            self = this;
        var tenant = window.sessionStorage.getItem("tenant");
        apiSrv.queryUserByDeptId(deptId, pageNo, pageSize, tenant).then(res => {
            if (res.status === "success") {
                let userData = res.list;

                set(model.meta, 'totalPage', res.meta.pages);
                set(model.meta, 'page', res.meta.offset);
                set(model, 'indexNumber', (pageNo - 1) * pageSize + 1);

                ////// 用户在选择了其他部门或者其他table页的用户信息也需要添加checked标志
                userData.forEach(item => {
                    if (model.selUsers.findBy("id", item.id)) {
                        set(item, 'checked', true);
                    }
                });

                set(model, 'users', userData);

                //////  用户切换部门的时候，需要将全选标记去掉
                set(model, 'AllChecked', false);

                ///// 若用户组id为空，不加载db当中已经选择的用户信息
                if (Ember.isBlank(model.userGroupId)) {
                    return;
                }

                portalUserGroupSrv.queryById(model.userGroupId).then(res => {
                    if ("success" === res.status) {
                        let groupData = res.item;
                        let db_selUsers = groupData.selUsers;

                        model.users.forEach(item => {
                            if (db_selUsers.contains(item.id)) {
                                set(item, 'checked', true);
                            }
                        });
                        set(model, 'bakUserList', Ember.copy(db_selUsers));
                    } else {
                        self.transitionToRoute('error', '查询用户组关联用户信息异常');
                    }
                });

            } else {
                self.transitionToRoute('error', '查询部门人员异常');
            }
        });
    },

    /**
     * 事件集合
     * @type {Object}
     */
    actions: {

        /**
         * 全选操作
         * @return {[type]} [description]
         */
        checkBoxClick: function() {
            let model = get(this, 'model'),
                checked = !model.AllChecked;

            model.users.forEach(item => {
                set(item, 'checked', checked);
            });
            /////// 控制全选的样式
            set(model, 'AllChecked', checked);
        },

        /**
         * 页面操作
         * @param  {[type]} pageNo [description]
         * @return {[type]}        [description]
         */
        pageClick: function(pageNo) {
            let model = get(this, 'model'),
                apiSrv = get(this, 'apiSrv'),
                self = this;
            var tenant = window.sessionStorage.getItem("tenant");

            //////////  根据部门id分页查询用户列表信息
            apiSrv.queryUserByDeptId(model.deptID, pageNo, pageSize, tenant).then(res => {
                if (res.status === "success") {
                    let userData = res.list;
                    //set(model.meta, 'totalPage', res.meta.total);
                    //set(model.meta, 'page', res.meta.pages);

                    ////// 用户在选择了其他部门或者其他table页的用户信息也需要添加checked标志
                    userData.forEach(item => {
                        if (model.selUsers.findBy("id", item.id)) {
                            set(item, 'checked', true);
                        }
                    });
                    set(model, 'users', userData);
                    set(model, 'AllChecked', false); //翻页的时候将全选标志标记为false
                } else {
                    self.transitionToRoute('error', '查询部门人员异常');
                }
            });


        },
        /**
         * 链接事件
         * @param  {[type]} targetRoute [description]
         * @return {[type]}             [description]
         */
        linkAction: function(targetRoute) {
            var param = {
                id: "",
                type: "ALL",
                source: "userGroup.operate",
            };
            this.transitionToRoute('dashboard.portalManagement.' + targetRoute + '.operate', {
                queryParams: param
            });
        },
        /**
         * tab标签切换
         * @param  {[type]} selectTab [description]
         * @return {[type]}           [description]
         */
        selectTabAction: function(step) {
            let userGroupSteps = get(this, 'userGroupSteps');

            userGroupSteps.forEach(item => {
                if (item.index !== step.index) {
                    set(item, 'status', 'wait');
                    return;
                }
                set(step, 'status', 'process');
            });

            set(this, 'currentTab', step.title);
        },

        /**
         * 点击部门节点，根据部门id分页查询用户列表信息
         * @param  {[type]} node [description]
         * @return {[type]}      [description]
         */
        onSelectedFolderNode: function(param) {
            var model = get(this, 'model'),
                pageNo = 1;
            var tenant = window.sessionStorage.getItem("tenant");
            set(model, 'deptID', param.id);
            this.queryUserByDeptId(param.id, pageNo, tenant);
        },

        /**
         * 点击步骤操作
         * @param  {[type]} index [description]
         * @return {[type]}       [description]
         */
        stepAction: function(index) {
            let currentStep = get(this, 'userGroupSteps').findBy("index", index);
            this.send('selectTabAction', currentStep);
        },

        /**
         * 从已经选择的用户列表当中，删除某一用户
         * @param  {[type]} user [description]
         * @return {[type]}      [description]
         */
        removeUserAction: function(user) {
            let model = get(this, 'model'),
                selUsers = model.selUsers;

            if (!selUsers) {
                selUsers = [];
                return;
            }
            selUsers.removeObject(user);

            ///// 删除用户后，取消当前用户列表此用户的checked状态
            model.users.forEach(item => {
                if (item.id === user.id) {
                    set(item, 'checked', false);
                }
            });
        },

        /**
         * 保存/更新用户组信息接口
         * @param  {[type]} userGroupForm [description]
         * @return {[type]}               [description]
         */
        saveUesrGroupAction: function() {
            var portalUserGroupSrv = get(this, 'portalUserGroupSrv'),
                userGroupForm = get(this, 'userGroupForm'),
                model = get(this, 'model'),
                userGroupId = get(model, 'userGroupId'),
                name = userGroupForm.name,
                type = userGroupForm.type,
                remark = userGroupForm.remark,
                self = this;
            if (Ember.isBlank(userGroupId)) {
                portalUserGroupSrv.create(name, type,remark).then(function(data) {
                    if (data.status === "success") {
                        swal(name + "用户组保存成功！");
                        set(model, 'userGroupId', data.item.id);
                        self.send('stepAction', 2);
                    } else {
                        swal(data.message);
                    }

                });

            } else {
                portalUserGroupSrv.update(userGroupId, name, type,remark).then(function(data) {
                    if (data.status === "success") {
                        swal(name + "用户组修改成功！");
                        self.send('stepAction', 2);
                    } else {
                        swal(data.message);
                    }

                });
            }
        },

        /**
         * 保存用户组与角色的关联关系
         * @return {[type]} [description]
         */
        saveGroupRoleAction: function() {
            var portalUserGroupSrv = get(this, 'portalUserGroupSrv'),
                model = get(this, 'model'),
                userGroupId = get(model, 'userGroupId'),
                selRoles = get(model, 'selectedRoleList'),
                bakRoleList = get(model, 'bakRoleList');

            if (Ember.isBlank(userGroupId)) {
                swal("用户组信息不能为空");
                return;
            }

            ////////// 将备份的角色与现选择的角色进行比对，区分新增/删除
            ////////// 若为新增，id:I，若为删除，id:D
            var roleIdListStr = this.compareList(bakRoleList, selRoles);

            portalUserGroupSrv.allocateRole(userGroupId, roleIdListStr).then(function(data) {
                if (data.status === "success" && !Ember.isBlank(userGroupId)) {
                    swal("保存成功");
                    set(model, 'bakRoleList', Ember.copy(selRoles));
                } else {
                    swal(data.message);
                }
            });
        },

        /**
         * 保存用户组与用户的关联关系
         * @return {[type]} [description]
         */
        saveGroupUserAction: function() {
            var portalUserGroupSrv = get(this, 'portalUserGroupSrv'),
                model = get(this, 'model'),
                userGroupId = get(model, 'userGroupId'),
                bakUserList = get(model, 'bakUserList');

            if (Ember.isBlank(userGroupId)) {
                swal("用户组信息不能为空");
                return;
            }

            let selUsers = Ember.copy(model.selUsers);

            // 比较用户，区分删除，新增状态
            selUsers.forEach(item => {
                set(item, 'operType', ""); //防止注入operType
                if (bakUserList.contains(item.id)) {
                    bakUserList.removeAt(bakUserList.indexOf(item.id));
                } else {
                    set(item, 'operType', "I");
                }
            });

            bakUserList.forEach(item => {
                let userInfo = {
                    "id": item,
                    "operType": "D"
                };
                selUsers.pushObject(userInfo);
            });

            portalUserGroupSrv.allocateUser(userGroupId, JSON.stringify(selUsers)).then(function(data) {
                if (data.status === "success") {
                    swal("保存成功！");
                    //////// 保存成功后，刷新备份的用户id列表
                    set(model, 'bakUserList', model.selUsers.mapBy("id"));
                } else {
                    swal(data.message);
                }

            });
        }
    }
});