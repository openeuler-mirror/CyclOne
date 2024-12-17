import Ember from 'ember'
const { get, set, inject } = Ember
const pageSize = 10
const rootNode = {
    id: '',
    title: '所有部门',
    sourceType: ''
}
export default Ember.Controller.extend({
    apiSrv: inject.service('api/portal-api/service'),
    userSrv: inject.service('api/portal-user/service'),
    deptSrv: inject.service('api/portal-dept/service'),
    commonSrv: inject.service('api/portal-common/service'),
    queryParams: ['deptID'],
    openPersonModalWindow: false,
    openPersonModalWindowTitle: '用户信息',
    openDeptModalWindow: false,
    openDeptRoleWindow: false,
    editUser: {},
    editDept: {},
    allRoles: [],
    selectedRoles: [],
    openResetPasswordWindow: false,
    openNewPasswordContent: '',
    openGrantTokenWindow: false,
    openGrantTokenContent: '',
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
     * 初始化用户操作
     */
    initUser: function(model) {
        var self = this
        var meta = {
            offset: 1
        }

        set(model, 'indexNumber', 1)
        self.queryDeptTree()
        self.queryUserByDeptId(model.deptID, meta.offset)
    },

    queryDeptTree: function() {
        var apiSrv = get(this, 'apiSrv')
        var self = this
        var tenant = window.sessionStorage.getItem('tenant')
        var model = get(this, 'model')
        apiSrv.queryDepartment(tenant).then(res => {
            if (res.status === 'success') {
                //Add root dept
                res.content.insertAt(0, rootNode)
                set(model, 'deptTree', res.content)
            } else {
                swal(res.message)
            }
        })
    },

    queryUserByCnd: function(cnd, pageNo) {
        let apiSrv = get(this, 'apiSrv'),
            model = get(this, 'model')
        let deptId = model.deptID
        set(this, 'loading', true)
        this.queryUserByDeptId(deptId, pageNo, cnd)
    },

    /**
     * 监听查询条件
     * @return {[type]} [description]
     */
    queryUsersByCn: function() {
        this.queryUserByCnd(this.filterString, 1)
    }.observes('filterString'),

    /**
     * 分页查询操作
     */
    queryUserByDeptId: function(deptId, pageNo, cnd) {
        var apiSrv = get(this, 'apiSrv'),
            model = get(this, 'model')
        set(model, 'paginationShow', false)
        set(this, 'loading', true)
        var tenant = window.sessionStorage.getItem('tenant')
        apiSrv.queryUserByDeptId(deptId, pageNo, 10, tenant, cnd).then(res => {
            this.set('loading', false)
            if (res.status === 'success') {

                res.list.forEach(item=>{
                    Ember.set(item,'firstValue',item.selGroups[0]);
                    set(item, 'groupNumbers', item.selGroups.length);
                });
                set(model, 'users', res.list);
                set(model, 'meta', res.meta);
                if (res.meta.total > 0) {
                    set(model, 'paginationShow', true)
                }
            } else {
                swal(res.message)
            }
        })
    },
    //////// 将角色列表根据指定的key，进行转换
    parseTransferList: function(list) {
        list.forEach(item => {
            set(item, 'label', item.name)
            set(item, 'value', item.id)
        })
        return list
    },

    /**
     * [actions description]
     * @type {Object}
     */
    actions: {
        /** 点击部门树，查询用户列表信息 */
        onSelectedFolderNode: function(param) {
            var model = get(this, 'model'),
                pageNo = 1

            set(model, 'deptID', param.id)
            set(model, 'indexNumber', (pageNo - 1) * pageSize + 1)
            this.set('loading', true)
            this.queryUserByDeptId(param.id, pageNo)
        },

        lookUpAction: function(record) {
            this.transitionToRoute(
                'dashboard.portalManagement.user.view',
                record.id
            )
        },
        /**
         * 页面操作
         * @param  {[type]} pageNo [description]
         * @return {[type]}        [description]
         */
        pageClick: function(pageNo) {
            var model = get(this, 'model')
            set(model.users, 'offset', pageNo)
            set(model, 'indexNumber', (pageNo - 1) * pageSize + 1)
            if (this.filterString) {
                this.queryUserByCnd(this.filterString, pageNo)
            } else {
                this.queryUserByDeptId(model.deptID, pageNo)
            }
        },

        /**
         * 分配用户组操作，跳转页面
         */
        editUserGroupAction: function(actionType, group) {
            var id = '',
                model = get(this, 'model')
            if (!Ember.isBlank(group)) {
                id = Ember.isBlank(group.id) ? '' : group.id
            }
            var param = {
                deptID: model.deptID,
                type: actionType
            }
            this.transitionToRoute(
                'dashboard.portalManagement.user.operate',
                id,
                {
                    queryParams: param
                }
            )
        },

        /**
         * 打开用户编辑窗口
         */
        viewPersonModalWindow: function(type, user) {
            var userSrv = get(this, 'userSrv')
            var model = get(this, 'model')
            this.set('openPersonModalWindow', true)
            set(model, 'isAdd', false)
            if (type === 'add') {
                set(model, 'isAdd', true)
                this.set('openPersonModalWindowTitle', '新增用户')
                this.set('editUser', { deptId: model.deptID })
            } else if (type === 'edit') {
                this.set('openPersonModalWindowTitle', '编辑用户')
                userSrv.queryById(user.id).then(res => {
                    if (res.status === 'success') {
                        this.set('editUser', res.item)
                    } else {
                        swal(res.message)
                    }
                })
            }
        },

        /**
         * 关闭用户编辑窗口
         */
        closePersonModalWindow: function() {
            this.set('openPersonModalWindow', false)
            this.set('editUser', {})
        },

        /**
         * 新增用户
         */
        saveUser: function() {
            var userSrv = get(this, 'userSrv')
            let tenantId = window.sessionStorage.getItem('tenant')
            let user = get(this, 'editUser')
            var model = get(this, 'model')
            user.tenantId = tenantId
            if (model.isAdd === true) {
                userSrv.addUser(user).then(res => {
                    if (res.status === 'success') {
                        swal('新增用户成功!')
                        this.set('openPersonModalWindow', false)
                        this.set('editUser', {})
                    } else {
                        swal(res.message)
                    }
                    this.queryUserByCnd(this.filterString, 1)
                })
            } else {
                userSrv.updateUser(user).then(res => {
                    if (res.status === 'success') {
                        swal('修改用户成功!')
                        this.set('openPersonModalWindow', false)
                        this.set('editUser', {})
                    } else {
                        swal(res.message)
                    }
                    this.queryUserByCnd(this.filterString, 1)
                })
            }
        },

        /**
         * 导出用户
         */
        downloadUser: function() {
            let userSrv = get(this, 'userSrv')
            let model = get(this, 'model')
            let deptId=model.deptID
            let tenantId = window.sessionStorage.getItem('tenant');
            location.href = '/portal/user/list/download?deptId='+deptId+'&tenantId='+tenantId;
        },

        /**
         * 激活用户
         */
        enabledUser: function(user) {
            var userSrv = get(this, 'userSrv')
            var model = get(this, 'model');
            if (user.status === 'ENABLED') {
                userSrv.disabledUser(user.id).then(res => {
                    if (res.status === 'success') {
                        swal('禁用用户成功!')
                    } else {
                        swal(res.message)
                    }
                    this.queryUserByCnd(this.filterString, 1)
                })
            } else {
                userSrv.enabledUser(user.id).then(res => {
                    if (res.status === 'success') {
                        swal('激活用户成功!')
                    } else {
                        swal(res.message)
                    }
                    this.queryUserByCnd(this.filterString, 1)
                })
            }
        },

        /**
         * 重置密码
         */
        resetPW: function(user) {
            var userSrv = get(this, 'userSrv');
            var self = this
            swal(
                {
                    title: '是否重置此用户密码?',
                    type: 'warning',
                    showCancelButton: true,
                    confirmButtonClass: 'btn-danger',
                    confirmButtonText: '重置密码',
                    cancelButtonText: '取消操作',
                    closeOnConfirm: true,
                    closeOnCancel: true    
                },
                function(isConfirm) {
                    if (isConfirm) {
                        userSrv.resetPW(user.id).then(res => {
                            set(self, 'openResetPasswordWindow', true);
                            set(self, 'openNewPasswordContent', res.item);
                            self.queryUserByCnd(self.filterString, 1);
                        })
                    }
                }
            )
        },

        closeResetPasswordWindow: function() {
            set(this, 'openResetPasswordWindow', false)
            set(this, 'openNewPasswordContent', "")
        },

        /**
         * 删除用户
         */
        deleteUser: function(user) {
            var self = this
            var userSrv = get(this, 'userSrv')
            // userSrv.deleteUser(user.id)

            if (!user && Ember.isBlank(user.id)) {
                swal('用户信息为空，请确认')
                return
            }
            swal(
                {
                    title: '是否删除此用户?',
                    type: 'warning',
                    showCancelButton: true,
                    confirmButtonClass: 'btn-danger',
                    cancelButtonText: '取消',
                    confirmButtonText: '删除',
                    closeOnConfirm: false
                },
                function(isConfirm) {
                    if (isConfirm) {
                        userSrv.deleteUser(user.id).then(res => {
                            if ('success' === res.status) {
                                swal('删除用户成功!')
                                self.queryUserByCnd(this.filterString, 1)
                            } else {
                                swal(res.message)
                            }
                        })
                    }
                }
            )
        },

        /**
         * admin给其他用户发放Token
         */
        grantToken: function(user) {
            var self = this
            var userSrv = get(this, 'userSrv')
            let tenantId = window.sessionStorage.getItem('tenant')
            let loginId=user.loginId;
            userSrv.grantToken(loginId,tenantId).then(res => {
                if (res.status === 'success') {
                    set(self, 'openGrantTokenWindow', true);
                    set(self, 'openGrantTokenContent', res.item);
                    self.queryUserByCnd(self.filterString, 1);
                    // swal(res.item)
                } else {
                    swal(res.message)
                }
            })
        },

        closeGrantTokenWindow: function() {
            set(this, 'openGrantTokenWindow', false)
            set(this, 'openGrantTokenContent', "")
        },

        viewDeptModalWindow: function(type) {
            var deptSrv = get(this, 'deptSrv')
            var model = get(this, 'model')
            set(model, 'isAddDept', false)
            if (type === 'add') {
                this.set('openDeptModalWindow', true)
                set(model, 'isAddDept', true)
                this.set('openDeptModalWindowTitle', '新增部门')
            } else if (type === 'edit') {
                if (Ember.isBlank(model.deptID)) {
                    swal('请选择一个部门！')
                    return
                }
                deptSrv.queryDeptById(model.deptID).then(res => {
                    if (res.status === 'success') {
                        if(res.item.sourceType !== 'native'){
                            swal('数据源为'+res.item.sourceType+'的部门不允许修改！')
                            return
                        }
                        this.set('editDept', res.item)
                        this.set('openDeptModalWindow', true)
                        this.set('openDeptModalWindowTitle', '部门信息')
                    } else {
                        swal(res.message)
                    }
                })
            }
        },

        closeDeptModalWindow: function() {
            this.set('openDeptModalWindow', false)
            this.set('openDeptModalWindowTitle', '')
            this.set('editDept', {})
        },

        saveDept: function() {
            var model = get(this, 'model')
            var dept = this.editDept
            var deptSrv = get(this, 'deptSrv')
            dept.tenantId = window.sessionStorage.getItem('tenant')
            if (model.isAddDept === true) {
                deptSrv.addDept(dept).then(res => {
                    if (res.status === 'success') {
                        swal('新增部门成功!')
                    } else {
                        swal(res.message)
                    }
                    this.queryDeptTree()
                })
            } else {
                deptSrv.updateDept(dept).then(res => {
                    if (res.status === 'success') {
                        swal('修改部门成功!')
                    } else {
                        swal(res.message)
                    }
                    this.queryDeptTree()
                })
            }
            this.set('openDeptModalWindow', false)
            this.set('editDept', {})
        },

        /**
         * 打开部门绑定角色的窗口
         */
        viewDeptRoleModalWindow: function() {
            var self = this
            var model = get(this, 'model')
            if (Ember.isBlank(model.deptID)) {
                swal('请选择一个部门！')
                return
            }
            this.set('openDeptRoleWindow', true)
            var commonSrv = get(this, 'commonSrv')
            var deptSrv = get(this, 'deptSrv')
            commonSrv.getRoleList().then(res => {
                if (res.status === 'success') {
                    this.set('allRoles', self.parseTransferList(res.list))
                } else {
                    swal(res.message)
                }
            })
            deptSrv.queryRolesByDeptId(model.deptID).then(res => {
                if (res.status === 'success') {
                    let roleIds = []
                    for (let item of res.item.links.ROLE) {
                        roleIds.push(item.id)
                    }
                    this.set('selectedRoles', roleIds)
                } else {
                    swal(res.message)
                }
            })
        },

        /**
         * 关门部门绑定角色的窗口
         */
        closeDeptRoleModalWindow: function() {
            this.set('openDeptRoleWindow', false)
            this.set('selectedRoles', [])
        },

        /**
         * 保存部门与角色的绑定
         */
        saveDeptToRole: function() {
            var model = get(this, 'model')
            var deptSrv = get(this, 'deptSrv')
            deptSrv
                .saveDeptToRole(model.deptID, this.selectedRoles)
                .then(res => {
                    if (res.status === 'success') {
                        swal('绑定部门角色成功!')
                    } else {
                        swal(res.message)
                    }
                })
            this.set('openDeptRoleWindow', false)
            this.set('selectedRoles', [])
        },

        /**
         * 删除部门
         */
        deleteDept: function() {
            var self = this
            var model = get(this, 'model')
            if (Ember.isBlank(model.deptID)) {
                swal('请选择一个部门！')
                return
            }
            var deptSrv = get(this, 'deptSrv')
            swal(
                {
                    title: '是否删除此部门?',
                    type: 'warning',
                    showCancelButton: true,
                    confirmButtonClass: 'btn-danger',
                    cancelButtonText: '取消',
                    confirmButtonText: '删除',
                    closeOnConfirm: false
                },
                function(isConfirm) {
                    if (isConfirm) {
                        deptSrv.deleteDept(model.deptID).then(res => {
                            if ('success' === res.status) {
                                swal('删除部门成功!')
                                set(model, 'deptID', '')
                                self.queryDeptTree()
                            } else {
                                swal(res.message)
                            }
                        })
                    }
                }
            )
        }
    }
})
