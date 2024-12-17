package route

import (
	"encoding/json"
	"fmt"
	"idcos.io/cloudboot/model"
	"net/http"

	"idcos.io/cloudboot/middleware"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)

// PermissionCodeTree 查询权限码树
func PermissionCodeTree(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())

	typ := r.FormValue("type")

	if typ != model.MenuPermissionType && typ != model.ButtonPermissionType && typ != model.DataPermissionType && typ != model.APIPermissionType {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(fmt.Sprintf("未知权限类型,%s", typ)))
		return
	}

	setting, err := repo.GetSystemSetting(typ)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	body := myhttp.RespBody{}
	if err := json.Unmarshal([]byte(setting.Value), &body); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, &body)
}
