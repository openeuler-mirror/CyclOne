import Ember from 'ember';
import ajaxUploadFile from 'clouduam-web/components/common/io-file-upload/ajax-file-upload';

export default Ember.Controller.extend({
    actions: {
        saveAction: function () {
            var self = this;
            var files = $('.file-uploader')[0].files;
            if (Ember.isEmpty(files) || files.length <= 0) {
                swal("操作失败，请选择要上传的文件！");
                return;
            }
            var tenantId = window.sessionStorage.getItem("tenant");
            ajaxUploadFile({
                url: '/portal/user/upload_users?tenantId='+tenantId,
                fileElement: $('.file-uploader')[0],
                onStart() {
                },
                onSuccess(response) {
                    var data = $.parseJSON(response);
                    if (data.status === "success") {
                        // swal("导入用户成功！");
                        self.transitionToRoute('dashboard.portalManagement.user.importPriview', data.item);
                    } else {
                        swal("操作失败！");
                        return;
                    }
                },
                onError(err) {
                    swal("操作失败：" + err);
                    return;
                }
            });
        },
    }
});