import Ember from 'ember';
const { get, set, inject } = Ember
export default Ember.Controller.extend({
	userSrv: inject.service('api/portal-user/service'),
	/**
     * [loading description]
     * @type {Boolean}
     */
    loading: false,
    clickToSave: true,

    /**
     * 初始化导入用户预览操作
     */
    initUser: function(model) {
        var self = this
        self.importPriview(model.id);
    },

    importPriview: function(id) {
        var userSrv = get(this, 'userSrv')
        var self = this
        var model = get(this, 'model')
        userSrv.importPriview(id).then(res => {

            if (res.status === 'success') {
				set(model, 'users', res.list);

            } else {
                swal(res.message)
            }
        })
    },

actions: {
        saveAction: function() {
            var userSrv = get(this, 'userSrv')
        	var self = this
        	var model = get(this, 'model')
            if(get(this, 'clickToSave')===true){
                set(this, 'clickToSave', false);
                userSrv.saveImportUsers(model.id).then(res => {
                    if (res.status === 'success') {
                        swal("导入用户成功！");
                        set(this, 'clickToSave', true);
                        self.transitionToRoute('dashboard.portalManagement.user.index');
                    } else {
                        swal("导入用户失败！");
                        set(this, 'clickToSave', true);
                    }
                })
            }
        	
        }

    }
});