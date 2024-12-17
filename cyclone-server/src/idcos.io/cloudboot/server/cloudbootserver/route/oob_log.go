package route

import (
	"net/http"

	"github.com/go-chi/chi"

	"idcos.io/cloudboot/job/mysql"
	"idcos.io/cloudboot/middleware"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)

//GetOOBlogBySN 根据sn获取impo传感器、事件数据
func GetOOBlogBySN(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	req := mysql.NewInspectionJob(log, repo, conf, "")
	sn := chi.URLParam(r, "sn")

	output := req.Collect(sn)
	if output.Error != "" {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(output.Error))
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
			"power_status": output,
		}),
	)
}
