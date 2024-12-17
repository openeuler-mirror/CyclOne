import Ember from 'ember';
const {set, get, inject, RSVP } = Ember;
export default Ember.Controller.extend({
    queryParams: ["id", "type", "source"],
    portalRoleSrv: inject.service('api/portal-role/service'),
    portalCommonSrv: inject.service('api/portal-common/service'),
    apiSrv: inject.service('api/portal-api/service'),
    resSrv: inject.service('api/portal-resource/service'),
    perSrv: inject.service('api/portal-permission/service'),
    roleId: "",
    roleForm: {
        code: "",
        name: "",
        remark: ""
    },
    unSelectedUserList: [],
    selectedUserList: [],
    unSelectedGroupList: [],
    selectedGroupList: [],
    selectedAuthResList: [],

    roleAllTabs: [{
        title: "角色信息",
        status: "process",
        active: true,
        index: 1,
        icon: "rolemanager"
    }, {
        title: "用户组信息",
        status: "wait",
        index: 2,
        icon: "usergroup"
    }, {
        title: "权限资源信息",
        status: "wait",
        index: 3,
        icon: "permissionmanager",
        last: true
    }],

    loading: false,

    //备份角色关联信息，用于提交数据进行比对，判断新增/删除
    bakUserList: [], //角色的关联用户信息备份，
    bakGroupList: [], //角色关联用户组信息备份
    bakAuthResList: [], //角色关联的功能信息备份

    // 监听系统名称
    selectedAppId: function() {
        var model = get(this, 'model');
        if (model.resData) {
            var resTypeData = model.resData.filterBy('appId', model.appId);
            set(model, 'resTypeData', resTypeData);
            var resType = resTypeData.mapBy('code').uniq();
            if (resType.length > 0) {
                set(model, 'resType', resType[0]);
            }else{
                set(model, 'resType', '');
            }
        }
    }.observes('model.appId'),

    // 穿梭框转换
    parseTransferList: function(list) {
        list.forEach((item) => {
            set(item, 'label', item.name);
            set(item, 'value', item.id);
        });
        return list;
    },
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
        selectedList.forEach((item) => {
            var selectedId = "";
            if (typeof(item) === "string") {
                selectedId = item;
            } else {
                selectedId = item.id;
            }

            if (backList.contains(selectedId)) {
                backList.removeAt(backList.indexOf(selectedId));
            } else {
                idListUnion.addObject(selectedId + ":I");
            }
        });

        backList.forEach((item) => {
            idListUnion.addObject(item + ":D");
        });
        return idListUnion.join();
    },
    /**
     * 遍历树，加载已经选择的节点信息
     * @param  {[type]} children [树的节点list]
     * @param  {[type]} perData  [已经选择的权限信息]
     */
    iteatorTree: function(children, perData) {
        let self = this,
            model = get(this, 'model');
        children.forEach(function(item) {
            if (perData.contains(item.id)) {
                set(item, 'selected', true);
                if (model.selRes.indexOf(item) < 0) {
                    model.selRes.pushObject(item);
                }
            }
            if (item.children.length > 0) {
                self.iteatorTree(item.children, perData);
            }
        });
    },

    observeSelRes: function() {
        let model = get(this, 'model');

        set(model, 'selPers', []);
        if (model.selRes) {
            model.selRes.forEach(item => {
                if (item.authResId) {
                    if (!model.selPers.findBy("id", item.authResId)) {
                        model.selPers.pushObject({ id: item.authResId, title: item.authResName });
                    }
                } else {
                    if (item.id && !model.selPers.findBy("id", item.id)) {
                        model.selPers.pushObject({ id: item.id, title: item.title });
                    }
                }
            });
        }
        set(model, 'selPersLength', model.selPers.length);
    }.observes('model.selRes.length'),

    /**
     * 操作角色初始化
     * 分配权限初始化
     * 分配用户组和用户以及修改角色基本信息初始化
     * @param  {[type]} model [description]
     * @return {[type]}       [description]
     */
    initOperateRole: function(model) {
        this.queryUserGroupType();
        var roleId = model.roleId,
            type = model.type,
            selectedGroupList = get(model, 'selectedGroupList'),
            portalRoleSrv = get(this, 'portalRoleSrv'),
            bakGroupList = get(model, 'bakGroupList'),
            bakUserList = get(model, 'bakUserList'),
            portalCommonSrv = get(this, 'portalCommonSrv'),
            resSrv = get(this, 'resSrv'),
            groupTypeSrv = get(this, 'apiSrv'),
            self = this;
        var tenantId = window.sessionStorage.getItem("tenant");
        //若id为空，说明为新增，初始化所有的用户组和角色信息
        if (Ember.isBlank(model.roleId)) {
            set(model, 'roleId', "");

            RSVP.hash({
                "allGroups": groupTypeSrv.getUserGroupByTypeAndTenantId("default",tenantId).then(res => {
                    return res.list
                }),
                "allRes": resSrv.queryAll().then(function(res) {
                    return res.list;
                })
            }).then(function(res) {
                set(model, 'dataSource', self.parseTransferList(res.allGroups));
                set(model, 'resData', res.allRes);
                var appIdList = res.allRes.mapBy("appId").uniq();
                set(model, 'appIdList', appIdList);
                set(model, 'appId', appIdList[0]);
                set(self, 'roleForm', {
                    code: "",
                    name: "",
                    remark: ""
                });
                set(self, 'model', model);
            });
            this.send('stepAction', 1);

        } else {
            /**
             * 若为分配权限，对权限信息进行处理
             * 1. 获取所有已经选择的权限信息对应关系
             * 2. 对功能列表进行初始化，并对已经选择的权限信息进行添加checkBox，checked等对象
             * 3. 备份功能信息，用于保存时进行比对，判断新增/删除
             * @param  {[type]} "PERMISSION" [description]
             * @return {[type]}              [description]
             */

            RSVP.hash({
                "roleItem": portalRoleSrv.queryById(roleId).then(function(res) {
                    return res.item;
                }),
                "allGroups": groupTypeSrv.getUserGroupByTypeAndTenantId("default",tenantId).then(res => {
                    return res.list
                }),
                "allRes": resSrv.queryAll().then(function(res) {
                    return res.list;
                }),
            }).then(function(res) {
                var roleItem = res.roleItem,
                    allGroups = res.allGroups,
                    selGroupIds = roleItem.selGroups;

                var appIdList = res.allRes.mapBy("appId").uniq();

                bakGroupList.pushObjects(selGroupIds);

                //处理已经选择的用户组和未选择的用户组
                allGroups.forEach((item) => {
                    if (selGroupIds.contains(item.id)) {
                        set(item, 'chosen', item.id);
                    }
                });

                set(model, 'resData', res.allRes);
                set(model, 'dataSource', self.parseTransferList(allGroups));
                selectedGroupList.pushObjects(selGroupIds);
                set(self, 'roleForm', roleItem);
                set(model, 'appIdList', appIdList);
                set(model, 'appId', appIdList[0]);
                set(model, 'selPers', []);
                set(self, 'model', model);
            });

            //控制标签页显示
            if ("INFO" === type || "ALL" === type) {
                this.send('stepAction', 1);
            }

            if ("GROUP" === type) {
                this.send('stepAction', 2);
            }

            if ("PERMISSION" === type) {
                this.send('stepAction', 3);
            }
        }
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
     * 监听用户组类型选择操作
     * @return {[type]} [description]
     */
    observeUserGroupType: function() {
        let groupTypeSrv = get(this, 'apiSrv');
        let tenant=window.sessionStorage.getItem("tenant");
        let model = get(this, 'model');
        let self = this
        groupTypeSrv.getUserGroupByTypeAndTenantId(model.groupType,tenant).then(resp => {
            if (resp.status === 'success') {
                let userGroups = resp.list;
                Ember.set(model, 'dataSource', self.parseTransferList(userGroups));
            } else {
                swal(resp.message)
            }
        });
    }.observes('model.groupType'),

    /**
     * 监听分配对象类型，功能id
     * 查询现有的功能与对象的关联信息
     * @return {[type]} [description]
     */
    changeTreeData: function() {
        var model = get(this, 'model'),
            apiSrv = get(this, 'apiSrv'),
            perSrv = get(this, 'perSrv'),
            self = this;

        if (Ember.isBlank(model.appId)) {
            return;
        }

        if (Ember.isBlank(model.resType)) {
            return;
        }


        set(this, 'loading', true);
        set(model, 'selRes', []);

        var tenantId = window.sessionStorage.getItem('tenant');
        apiSrv.queryResByResType(model.resType,tenantId).then((res) => {
            if (res.status === "success") {
                this.set('loading', false);
                let treeData = res.content;
                set(treeData, 'expanded', true);

                set(model, 'resTree', treeData);

                if (Ember.isBlank(model.roleId)) {
                    return;
                }

                var queryResForm = {
                    authResType: model.resType,
                    authObjType: "ROLE",
                    authObjId: model.roleId
                };

                perSrv.queryAuthRes(queryResForm).then((res) => {
                    if (res.status === "success") {
                        let perData = res.list.mapBy("authResId");

                        self.iteatorTree(model.resTree.children, perData);
                        set(model, 'backRes', Ember.copy(perData));

                    } else {
                        self.transitionToRoute(res.message);
                    }
                });

            } else {
                self.transitionToRoute(res.message);
            }
        });

    }.observes('model.resType', 'model.roleId'),

    /**
     * 将树平铺为List返回
     * @param  {[type]} children [树的节点list]
     */
    treeToList: function(children) {
        let list = [];
        children.forEach(item => { 
            if (list.indexOf(item) < 0) {
                list.push(item);
            }
            if (item.children && item.children.length > 0) {
                list=list.concat(this.treeToList(item.children));
            }
        });
        return list;
    },

    /**
     * 事件集合
     * @type {Object}
     */
    actions: {
        linkAction: function(targetRoute) {
            var param = {
                id: "",
                type: "ALL",
                source: "role.operate"
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
            get(this, 'roleAllTabs').forEach(item => {
                if (item.index === step.index) {
                    set(item, 'status', 'process');
                    return;
                }
                set(item, 'status', 'wait');
            });
            set(this, 'currentTab', step.title);
        },

        /**
         * 点击步骤操作
         * @param  {[type]} index [description]
         * @return {[type]}       [description]
         */
        stepAction: function(index) {
            let roleAllTabs = get(this, 'roleAllTabs');
            let currentStep = roleAllTabs.findBy("index", index);
            this.send('selectTabAction', currentStep);
        },

        /**
         * 删除已经选择的权限信息
         * @param  {[type]} per [description]
         * @return {[type]}     [description]
         */
        removePermission: function(per) {
            let model = get(this, 'model');
            model.selPers.removeObject(per);
            let node = model.selRes.findBy('id', per.id);
            set(node, 'selected', false);
            model.selRes.removeObject(node);
        },

        /**
         * 返回链接
         * @return {[type]} [description]
         */
        backRoleAction: function() {
            var source = get(this, 'source');
            if (Ember.isBlank(source)) {
                this.transitionToRoute('dashboard.portalManagement.role.index');
            } else {
                if (source === "user.operate") {
                    this.transitionToRoute('dashboard.portalManagement.user');
                } else {
                    this.transitionToRoute('dashboard.portalManagement.' + source);
                }
                set(this, 'source', "");
            }
        },

        /**
         * 保存角色操作
         * @param  {[type]} roleForm [description]
         * @return {[type]}          [description]
         */
        saveRoleAction: function() {
            var portalRoleSrv = get(this, 'portalRoleSrv'),
                model = get(this, 'model'),
                roleForm = get(this, 'roleForm'),
                self = this;


            if (Ember.isBlank(model.roleId)) {
                portalRoleSrv.create(roleForm).then(function(data) {
                    if (data.status === "success") {
                        swal(roleForm.name + "角色保存成功！");
                        self.send('stepAction', 2);
                        set(model, 'roleId', data.item.id);
                    } else {
                        swal(data.message);
                    }
                });
            } else {
                portalRoleSrv.update(model.roleId, roleForm).then(function(data) {
                    if (data.status === "success") {
                        swal(roleForm.name + "角色修改成功！");
                    } else {
                        swal(data.message);
                    }
                });
            }
        },
        /**
         * 保存用户组与角色关联操作
         * @return {[type]} [description]
         */
        saveGroupRoleAction: function() {
            var model = get(this, 'model'),
                selectedGroupList = get(model, 'selectedGroupList'),
                portalRoleSrv = get(this, 'portalRoleSrv'),
                bakGroupList = get(this, 'bakGroupList'),
                roleId = get(model, 'roleId'),
                self = this;

            var selGroups = this.compareList(bakGroupList, selectedGroupList);

            if (Ember.isBlank(roleId)) {
                swal("角色信息不能为空，请确认");
                return;
            }

            portalRoleSrv.allocateGroup(roleId, selGroups).then(function(data) {
                if (data.status === "success") {
                    swal("保存成功！");
                    set(self, 'bakGroupList', Ember.copy(selectedGroupList));
                    // self.send('stepAction', 3);
                } else {
                    swal(data.message);
                }
            });
        },

        /**
         * 保存权限信息
         * @return {[type]} [description]
         */
        savePermissionAction: function() {
            let model = get(this, 'model'),
                perSrv = get(this, 'perSrv'),
                self = this;

            if (Ember.isBlank(model.roleId)) {
                swal("角色信息为空，请确认");
                return;
            }

            let selPers = Ember.copy(model.selPers),
                resTree = get(model,'resTree'),
                operList = self.treeToList(resTree.children);

            let commitData = []; 
            operList.forEach(item => {
                let perInfo = { authResType: model.resType, authResId: item.id, authResName: item.title, authObjType: "ROLE", authObjId: model.roleId, operType: item.selected ? 'I' : 'D' };
                commitData.pushObject(perInfo);
            });

            perSrv.saveAuthRes(JSON.stringify(commitData)).then(function(data) {
                if ("success" === data.status) {
                    swal("保存成功！");
                    set(model, 'backRes', model.selPers.mapBy("id"));
                } else {
                    swal(data.message);
                }

            });
        },
    }
});
