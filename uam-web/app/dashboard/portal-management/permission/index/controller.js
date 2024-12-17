import Ember from 'ember';
const {set,
    get,
    inject,
    RSVP
} = Ember;
const operSelect = ['选择权限'];
export default Ember.Controller.extend({
    perSrv: inject.service('api/portal-permission/service'),
    commSrv: inject.service('api/portal-common/service'),
    apiSrv: inject.service('api/portal-api/service'),
    resSrv: inject.service('api/portal-resource/service'),

    loading: false,
    multi: [],

    /**
     * 初始化权限信息
     * 初始化树结构
     * @return {[type]} [description]
     */
    initPermission: function(model) {
        var commSrv = get(this, 'commSrv'),
            resSrv = get(this, 'resSrv'),
            self = this;

        RSVP.hash({
            "resData": resSrv.queryAll().then((res) => {
                return res.list;
            }),
            "allRols": commSrv.getRoleList().then((res) => {
                return res.list;
            }),
        }).then((data) => {
            var appIdList = data.resData.mapBy("appId").uniq();
            let roles = self.parseRoleList(data.allRols);
            set(model, 'appIdList', appIdList);
            set(model, 'resData', data.resData);
            set(model, 'operSelect', operSelect);
            set(model, 'operId', operSelect[0]);
            set(model, 'selectedRole', roles[0]);
            set(model, 'roleData', roles);
            set(model, 'selPers', []);
            if (appIdList.length > 0) {
                set(model, 'appId', appIdList[0]);
            }

            if (data.allRols.length > 0) {
                set(model, 'roleId', data.allRols[0].id);
                set(model, 'roleName', data.allRols[0].name);
            }
            
            set(self, 'model', model);
        });
    },

    /**
     * 解析角色列表，将角色列表数据结构转换为树形结构
     * @param  {[type]} roleList [description]
     * @return {[type]}          [description]
     */
    parseRoleList: function(roleList) {
        roleList.forEach((item) => {
            set(item, 'title', item.name);
        });
        return roleList;
    },

    /**
     * 遍历树，加载已经选择的节点信息
     * 为已经选择的节点添加selected=true的属性
     * 同时，
     * @param  {[type]} children [树的节点list]
     * @param  {[type]} perData  [已经选择的权限信息]
     */
    iteatorTree: function(children, perData) {
        let self = this,
        multi = get(this, 'multi');
        children.forEach(function(item) {
            if (perData.contains(item.id)) {
                set(item, 'selected', true);
                if (multi.indexOf(item) < 0) {
                    multi.pushObject(item);
                }
            }
            if (item.children && item.children.length > 0) {
                self.iteatorTree(item.children, perData);
            }
        });
    },

    /**
     * 监听操作类型
     * @return {[type]} [description]
     */
    selectOper: function() {
        let model = get(this, 'model');

        if (model.operId === operSelect[0]) {
            return;
        }
        this.transitionToRoute('dashboard.portalManagement.permission.reverse');
    }.observes('model.operId'),

    /**
     * 监听操作系统名称
     * @return {[type]} [description]
     */
    selectedAppId: function() {
        var model = get(this, 'model');
        if (model.resData) {
            var resTypeData = model.resData.filterBy('appId', model.appId);
            var resType = resTypeData.mapBy('code').uniq();
            set(model, 'resTypeData', resTypeData);
            if (resType.length > 0) {
                set(model, 'resType', resType[0]);
            }else{
                 set(model, 'resType', '');
            }
        }
    }.observes('model.appId'),

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
        // set(model, 'selRes', []);
        set(model, 'selPers', []);

        var tenantId = window.sessionStorage.getItem('tenant');
        apiSrv.queryResByResType(model.resType,tenantId).then((res) => {
            if (res.status === "success") {
                this.set('loading', false);
                let treeData = res.content;
                set(this, 'multi', []);
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
     * 监听选择树节点，选择后添加至已经选择的列表当中
     * @return {[type]} [description]
     */
    observeSelRes: function() {
        let model = get(this, 'model'),
            multi = get(this, 'multi');

        if (multi) {
            set(model, 'selPers', []);
            //遍历获取已经选择的权限列表信息
            multi.forEach(item => {
                if (item.authResId) {
                    if (!model.selPers.findBy("id", item.authResId)) {
                        model.selPers.pushObject({
                            id: item.authResId,
                            title: item.authResName
                        });
                    }
                } else {
                    if (item.id && !model.selPers.findBy("id", item.id)) {
                        model.selPers.pushObject({
                            id: item.id,
                            title: item.title
                        });
                    }
                }
            });
        }
        set(model, 'selPerLength', model.selPers.length);
    }.observes('multi.length'),

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
        selectResType: function(resType) {
            var model = get(this, 'model');
            set(model, 'resType', resType);
        },

        /**
         * 删除已经选择的权限信息
         * @param  {[type]} per [description]
         * @return {[type]}     [description]
         */
        removePermission: function(per) {
            let model = get(this, 'model');
            model.selPers.removeObject(per);
            if (model.selRes.findBy('authResId', per.id)) {
                set(model.selRes.findBy('authResId', per.id), 'selected', false);
            }

            if (model.selRes.findBy('id', per.id)) {
                set(model.selRes.findBy('id', per.id), 'selected', false);
            }
        },

        selRoleAction: function(param) {
            let model = get(this, 'model');
            if (!param.id) {
                swal("角色信息为空，请确认");
                return;
            }
            set(model, 'roleId', param.id);
            set(model, 'roleName', param.title);
        },

        savePermissionAction: function() {
            let model = get(this, 'model'),
                perSrv = get(this, 'perSrv'),
                self = this;

            if (Ember.isBlank(model.roleId)) {
                swal("角色信息为空，请确认");
                return;
            }

            let resTree = get(model,'resTree'),
                operList = self.treeToList(resTree.children);

            let commitData = [];
            operList.forEach(item => {
                if (item.id) {
                    let perInfo = {
                        authResType: model.resType,
                        authResId: item.id,
                        authResName: item.title,
                        authObjType: "ROLE",
                        authObjId: model.roleId,
                        operType: item.selected ? 'I' : 'D'
                    };
                    commitData.pushObject(perInfo);
                }
            });
            perSrv.saveAuthRes(JSON.stringify(commitData)).then(function(data) {
                if ("success" === data.status) {
                    swal("保存成功！");
                } else {
                    swal(data.message);
                }

            });
        },
    }
});