package route

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/voidint/binding"

	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)

// GenPXE4CentOS6UEFI 为CentOS6.x（UEFI）生成PXE文件、修改dhcp配置文件、重启dhcp服务。
func GenPXE4CentOS6UEFI(w http.ResponseWriter, r *http.Request) {
	log, _ := middleware.LoggerFromContext(r.Context())
	repo, _ := middleware.RepoFromContext(r.Context())

	reqData := service.GenPXE4CentOS6UEFIReq{
		SN: chi.URLParam(r, "sn"),
	}
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	filename, err := service.GenPXE4CentOS6UEFI(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"filename": filename,
	}))
}

// GenPXE 为目标设备生成PXE文件
func GenPXE(w http.ResponseWriter, r *http.Request) {
	log, _ := middleware.LoggerFromContext(r.Context())
	repo, _ := middleware.RepoFromContext(r.Context())

	reqData := service.GenPXEReq{
		SN: chi.URLParam(r, "sn"),
	}
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	filename, err := service.GenPXE(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"filename": filename,
	}))
}

// GetPXE 返回目标设备PXE
//From iPXE:
//:netboot
//chain http://osinstall/api/cloudboot/v1/devices/${serial}/pxe?arch=${buildarch} 
func GetPXE(w http.ResponseWriter, r *http.Request) {
	log, _ := middleware.LoggerFromContext(r.Context())
	repo, _ := middleware.RepoFromContext(r.Context())

	reqData := service.GetPXEReq{
		SN: chi.URLParam(r, "sn"),
	}

	reqData.SN, _ = url.QueryUnescape(reqData.SN)

	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	pxe, err := service.GetPXE(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.Text(w, http.StatusOK, []byte(pxe))
}
