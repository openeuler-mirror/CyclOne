import Ember from 'ember'
const { set, get, inject, RSVP } = Ember
export default Ember.Controller.extend({
    queryParams: ['deptID', 'type', 'source'],
    portalUserSrv: inject.service('api/portal-user/service'),
    portalGroupSrv: inject.service('api/portal-user-group/service'),
    portalCommonSrv: inject.service('api/portal-common/service'),
    apiSrv: inject.service('api/portal-api/service'),
    userId: '',

    unSelectedRoleList: [],
    selectedRoleList: [],
    unSelectedGroupList: [],
    selectedGroupList: [],

    /////////// 转换数组，将数据转换为
    parseTransferList: function(list) {
        list.forEach(item => {
            set(item, 'label', item.name)
            set(item, 'value', item.id)
        })
        return list
    },

    /**
     * 初始化用户信息
     * @param  {[type]} model [description]
     * @return {[type]}       [description]
     */
    initOperateUser: function(model) {
        this.queryUserGroupType();
        var userId = model.id,
            portalUserSrv = get(this, 'portalUserSrv'),
            groupTypeSrv = get(this, 'apiSrv'),
            self = this
        var tenantId = window.sessionStorage.getItem('tenant')

        //若id为空，说明为新增，初始化所用的用户组和角色信息
        RSVP.hash({
            allGroups: groupTypeSrv.getUserGroupByTypeAndTenantId("default",tenantId).then(res => {
                return res.list
            }),
            userData: portalUserSrv.queryById(userId).then(res => {
                return res.item
            })
        }).then(res => {
            var allGroups = res.allGroups,
                userData = res.userData
            /////////  已经选择的用户组信息 //////////
            var selGroups = []

            if (userData) {
                selGroups = userData.selGroups

                allGroups.forEach(item => {
                    if (selGroups.contains(item.id)) {
                        set(item, 'chose', item.id)
                    }
                })

            }
            set(model, 'backGroupList', Ember.copy(selGroups))
            set(model, 'selGroups', selGroups)
            set(model, 'groupSource', self.parseTransferList(allGroups))
            set(self, 'model', model)
        })
    },

    /**
     * 比较两个list，将删除，新增标志与id号进行拼接
     * @param  {[type]} backList     [description]
     * @param  {[type]} selectedList [description]
     * @return {[type]}              [description]
     */
    compareList: function(backList, selectedList) {
        var idListUnion = [] //id号与操作的拼接数据 id:I id:D

        //将selectedList与backList对比
        //若两个list都存在相同的id，从backList当中移除，只保留未选择的
        //若selectedList当中的id在backList当中不存在，则标记为I
        selectedList.forEach(item => {
            if (backList.contains(item)) {
                backList.removeAt(backList.indexOf(item))
            } else {
                idListUnion.addObject(item + ':I')
            }
        })
        //将backList里面所有的id都标记为D
        backList.forEach(item => {
            idListUnion.addObject(item + ':D')
        })
        return idListUnion.join()
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
        const self = this
        groupTypeSrv.getUserGroupByTypeAndTenantId(model.groupType,tenant).then(resp => {
            if (resp.status === 'success') {
                let userGroups = resp.list;
                Ember.set(model, 'groupSource', self.parseTransferList(userGroups));
            } else {
                swal(resp.message)
            }
        });
    }.observes('model.groupType'),

    /**
     * 事件集合
     * @type {Object}
     */
    actions: {
        linkAction: function(targetRoute) {
            var param = {
                id: '',
                type: 'ALL',
                source: 'user.operate'
            }
            this.transitionToRoute(
                'dashboard.portalManagement.' + targetRoute + '.operate',
                {
                    queryParams: param
                }
            )
        },
        /**
         * tab标签切换
         * @param  {[type]} selectTab [description]
         * @return {[type]}           [description]
         */
        selectTabAction: function(selectTab) {
            set(this, 'currentTab', selectTab)
        },

        /**
         * 返回链接
         * @return {[type]} [description]
         */
        backUserAction: function() {
            var source = get(this, 'source'),
                model = get(this, 'model')
            if (Ember.isBlank(source)) {
                this.transitionToRoute('dashboard.portalManagement.user', {
                    queryParams: {
                        deptID: model.deptID
                    }
                })
            } else {
                this.transitionToRoute('dashboard.portalManagement.' + source)
            }
        },

        /**
         * 保存用户组操作
         * @return {[type]} [description]
         */
        saveGroupAction: function() {
            var model = get(this, 'model'),
                userId = get(model, 'id'),
                selectedGroupList = get(model, 'selGroups'),
                backGroupList = get(model, 'backGroupList'),
                portalUserSrv = get(this, 'portalUserSrv'),
                apiSrv = get(this, 'apiSrv')

            var selGroups = this.compareList(backGroupList, selectedGroupList)

            let userForm = {
                id: userId,
                tenantId: window.sessionStorage.getItem('tenant'),
                selGroups: selGroups
            }

            portalUserSrv.allocate(userForm).then(function(data) {
                if ('success' === data.status) {
                    swal('保存成功')
                    set(model, 'backGroupList', Ember.copy(selectedGroupList))
                } else {
                    swal('保存失败', data.message)
                }
            })
        }
    }
})
