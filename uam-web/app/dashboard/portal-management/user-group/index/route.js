import Ember from 'ember';
const {set,
    get,
    inject,
    RSVP
} = Ember;
export default Ember.Route.extend({
    permissions: [],
    model: function() {
        this.set('permissions', window.sessionStorage.getItem("permissions").split(','))
		var columns = [{ "propertyName": "name", "title": "名称", "template": "partial/lookUp" }, 
                    { "propertyName": "remark", "title": "描述" }, 
                    { "propertyName": "type", "title": "类型" , "template": "partial/groupType"}, 
                    { "propertyName": "selRoles", "title": "拥有角色", "template": "partial/groupSelRoles" }, 
                    { "propertyName": "gmtModified", "title": "最后修改时间" },
                    this.permissions.includes("userGroup.operate")? { "title": "操作", "template": "partial/groupOper" }: {}
                    ];
        return { "columns": columns };
    },
    setupController: function(controller, model) {
        set(model, 'permissions', this.permissions);
        controller.initUserGroup(model);
        controller.set('model', model);
    }
});
