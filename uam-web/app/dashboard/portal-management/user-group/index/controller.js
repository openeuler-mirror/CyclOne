import Ember from 'ember';
const {set, get, inject } = Ember;
const pageSize = 10;
export default Ember.Controller.extend({
    groupSrv: inject.service('api/portal-user-group/service'),
    //查询form
    queryForm: {
        code: null,
        name: null
    },
    filterString: '',
    queryUserGroupData: function(offset, queryForm, cnd) {
        var groupSrv = get(this, 'groupSrv'),
            model = get(this, 'model');
        set(model, 'paginationShow', true);

        this.set('loading', true);
        groupSrv.queryByPage(offset, pageSize, queryForm, cnd).then((data) => {
            set(model, 'meta', data.meta);
            set(model, 'data', data.list);
            data.list.forEach(item=>{
                    Ember.set(item,'firstValue',item.selRoles[0]);
                    set(item, 'roleNumbers', item.selRoles.length);
            });
            this.set('loading', false);
            if (data.meta.pages > 1) {
                set(model, 'paginationShow', true);
            }
        });
    },

    queryUserGroupByCnd: function(cnd, offset) {
        set(this, 'loading', true)
        this.queryUserGroupData(offset, this.queryForm, cnd)
    },

    /**
     * 监听查询条件
     * @return {[type]} [description]
     */
    queryUserGroupsByCn: function() {
        this.queryUserGroupByCnd(this.filterString, 1);
    }.observes('filterString'),


    initUserGroup: function(model) {
        var meta = {
            offset: 1,
        };
        set(model, 'meta', meta);
        set(model, 'indexNumber', 1);
        set(this, 'model', model);
        this.queryUserGroupData(model.meta.offset, this.queryForm, "");
    },

    // /**
    //  * 监听页面录入名称，查询数据列表信息
    //  * @param  {[type]} queryForm [description]
    //  * @return {[type]}           [description]
    //  */
    // queryByName: function(queryForm) {
    //     var model = get(this, 'model');
    //     this.queryUserGroupData(model.meta.offset, model.meta.pageSize, this.queryForm);
    // }.observes('queryForm.name'),
    loading: false,
    actions: {
        /**
         * 分页查询操作
         * @param  {[type]} pageNo [description]
         * @return {[type]}        [description]
         */
        pageClick: function(pageNo) {
            var model = get(this, 'model');

            // 设置页面的number开头
            set(model, 'indexNumber', (pageNo - 1) * pageSize + 1);
            set(model.meta, 'offset', pageNo);
            if (this.filterString) {
                this.queryUserGroupByCnd(this.filterString, pageNo);
            } else {
                this.queryUserGroupData(model.meta.offset, this.queryForm, "");
            }
        },

        editUserGroupAction: function(actionType, group) {
            var id = "";
            if (!Ember.isBlank(group)) {
                id = Ember.isBlank(group.id) ? "" : group.id;
            }
            var param = {
                id: id,
                type: actionType
            };
            this.transitionToRoute('dashboard.portalManagement.userGroup.operate', {
                queryParams: param
            });
        },

        lookUpAction: function(record) {

            if (Ember.isBlank(record.id)) {
                swal("用户组id不能为空，请确认");
                return;
            }

            this.transitionToRoute('dashboard.portalManagement.userGroup.view', record.id);
        },

        deleteUserGroupAction: function(group) {
            var userGroupId = group.id,
                model = get(this, 'model'),
                self = this,
                groupSrv = get(this, 'groupSrv');
            if (Ember.isBlank(userGroupId)) {
                swal("用户组信息为空，请确认");
                return;
            }
            swal({
                title: "是否删除此用户组?",
                type: "warning",
                showCancelButton: true,
                confirmButtonClass: "btn-danger",
                cancelButtonText: "取消",
                confirmButtonText: "删除",
                closeOnConfirm: false
            }, function(isConfirm) {
                if (isConfirm) {
                    groupSrv.delete(userGroupId).then(function(data) {
                        swal.close();
                        if ("success" === data.status) {
                            self.queryUserGroupByCnd(self.filterString, 1);
                            // self.queryUserGroupData(model.meta.offset, self.queryForm);
                        }
                    });
                }
            });

        }
    }
});
