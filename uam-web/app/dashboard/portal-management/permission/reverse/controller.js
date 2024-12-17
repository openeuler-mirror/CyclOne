import Ember from 'ember';
const {set,
    get,
    inject,
    RSVP
} = Ember;
const operSelect = ['选择权限', '选择角色'];
export default Ember.Controller.extend({
    perSrv: inject.service('api/portal-permission/service'),
    commSrv: inject.service('api/portal-common/service'),
    apiSrv: inject.service('api/portal-api/service'),
    resSrv: inject.service('api/portal-resource/service'),

    loading: false,
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

            set(model, 'appIdList', appIdList);
            set(model, 'resData', data.resData);
            set(model, 'operSelect', operSelect);
            set(model, 'operId', operSelect[1]);
            set(model, 'roleData', self.parseRoleList(data.allRols));
            set(model, 'appId', appIdList[0]);
            set(model, 'perId', "");
            set(model, 'selPers', []); //// 已经选择的数据集合，obj数组
            set(model, 'selRes', []); //// 勾选的数据集合
            set(model, 'bakRes', []); //// 备份的角色权限，只有id，用于selRes的比对
            set(self, 'model', model);
        });
    },

    /**
     * 监听分配对象类型，功能id
     * 查询现有的功能与对象的关联信息
     * @return {[type]} [description]
     */
    observePerTree: function() {
        var model = get(this, 'model'),
            apiSrv = get(this, 'apiSrv'),
            self = this;

        if (Ember.isBlank(model.appId)) {
            return;
        }

        if (Ember.isBlank(model.resType)) {
            return;
        }

        set(this, 'loading', true);

        var tenantId = window.sessionStorage.getItem('tenant');
        apiSrv.queryResByResType(model.resType,tenantId).then((res) => {
            if (res.status === "success") {
                this.set('loading', false);
                //////// 获取权限树数据，并展开一级
                let treeData = res.content;
                treeData.expanded = true;
                set(model, 'resTree', treeData);
            } else {
                self.transitionToRoute(res.message);
            }
        });
    }.observes('model.resType'),


    /**
     * 监听权限id，查询角色权限的关系数据信息
     * @return {[type]} [description]
     */
    obserPerId: function() {
        let model = get(this, 'model'),
            perId = model.perId,
            perSrv = get(this, 'perSrv'),
            self = this;

        if (Ember.isBlank(perId)) {
            return;
        }

        let queryForm = {
            authResId: perId,
            authResType: model.resType,
            authObjType: "ROLE"
        };

        perSrv.queryAuthObj(queryForm).then(function(res) {
            if ("success" === res.status) {
                let perData = res.list.mapBy("authObjId");
                ///////  添加备份信息
                set(model, 'bakRes', perData);
                // 修改roleData的状态，添加selected=true
                self.iteratorTree(model.roleData, perData);

            } else {
                swal("查询角色权限数据信息失败", res.message);
            }
        });
    }.observes('model.perId'),

    /**
     * 解析角色列表，将角色列表数据结构转换为树形结构
     * @param  {[type]} roleList [description]
     * @return {[type]}          [description]
     */
    parseRoleList: function(roleList) {
        roleList.forEach(item => {
            set(item, 'title', item.name);
        });
        return roleList;
    },

    /**
     * 监听操作类型
     * @return {[type]} [description]
     */
    obserOper: function() {
        let model = get(this, 'model');

        if (model.operId === operSelect[1]) {
            return;
        }
        this.transitionToRoute('dashboard.portalManagement.permission');
    }.observes('model.operId'),

    /**
     * 监听操作系统名称
     * @return {[type]} [description]
     */
    selectedAppId: function() {
        var model = get(this, 'model');
        var resTypeData = model.resData.filterBy('appId', model.appId);
        set(model, 'resTypeData', resTypeData);

        var resType = resTypeData.mapBy('code').uniq();
        if (resType.length > 0) {
            set(model, 'resType', resType[0]);
        }else{
            set(model, 'resType', '');
        }

    }.observes('model.appId'),


    iteratorTree: function(roleData, perData) {
        let model = get(this, 'model');

        roleData.forEach(item => {
            if (perData.contains(item.id)) {
                set(item, 'selected', true);
                if (!model.selPers.findBy("id", item.id)) {
                    model.selRes.pushObject(item);
                }
            } else {
                set(item, 'selected', false);
            }
        });
    },

    /**
     * 监听选择树节点，选择后添加至已经选择的列表当中
     * @return {[type]} [description]
     */
    observeSelRes: function() {
        let model = get(this, 'model');

        set(model, 'selPers', []);
        if (model.selRes) {
            model.selRes.forEach(item => {
                if (item.id && !model.selPers.findBy("id", item.id)) {
                    model.selPers.pushObject({
                        id: item.id,
                        title: item.title
                    });
                }
            });
        }

        set(model, 'selPersLength', model.selPers.length);
    }.observes('model.selRes.length'),


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
            let node = model.selRes.findBy('id', per.id);

            set(node, 'selected', false);
            model.selRes.removeObject(node);
        },

        selPerAction: function(param) {
            let model = get(this, 'model');
            if (!param.id) {
                swal("权限id为空，请确认");
                return;
            }
            set(model, 'perId', param.id);
            set(model, 'perName', param.title);
        },


        selRoleAction: function(param) {
            let model = get(this, 'model');
            if (!param.id) {
                swal("权限id为空，请确认");
                return;
            }
            set(model, 'perId', param.id);
            set(model, 'perName', param.title);
        },

        savePermissionAction: function() {
            var model = get(this, 'model'),
                bakRes = get(model, 'bakRes'),
                perSrv = get(this, 'perSrv');

            if (Ember.isBlank(model.perId)) {
                swal("当前权限信息为空，请确认");
                return;
            }

            var selPers = Ember.copy(model.selPers);

            selPers.forEach(item => {
                if (bakRes.contains(item.id)) {
                    bakRes.removeAt(bakRes.indexOf(item.id));
                } else {
                    set(item, 'operType', "I");
                }
            });

            bakRes.forEach(item => {
                let perInfo = {
                    id: item,
                    operType: "D"
                };
                selPers.pushObject(perInfo);
            });


            let commitData = [];
            selPers.forEach(item => {
                if (item.id) {
                    let perInfo = {
                        authResType: model.resType,
                        authResId: model.perId,
                        authResName: model.perName,
                        authObjType: "ROLE",
                        authObjId: item.id,
                        operType: item.operType
                    };
                    commitData.pushObject(perInfo);
                }
            });

            perSrv.saveAuthRes(JSON.stringify(commitData)).then(function(data) {
                if ("success" === data.status) {
                    swal("保存成功！");
                    set(model, 'bakRes', model.selPers.mapBy("id"));
                } else {
                    swal(data.message);
                }

            });
        },
    }
});
