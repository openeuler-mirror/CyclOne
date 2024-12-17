package cloudbootserver

import (
	"net/http"
	"os"

	"idcos.io/cloudboot/limiter"
	"idcos.io/cloudboot/limiter/webank"

	"github.com/go-chi/chi"
	"idcos.io/cloudboot/config"
	jmysql "idcos.io/cloudboot/job/mysql"
	"idcos.io/cloudboot/logger"
	mw "idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/model/mysqlrepo"
)

// Server API server
type Server struct {
	Conf    *config.Config
	Log     logger.Logger
	Repo    model.Repo
	handler http.Handler
}

// NewServer 实例化http服务
func NewServer(log logger.Logger, conf *config.Config) (*Server, error) {
	repo, err := mysqlrepo.NewRepo(conf, log)
	if err != nil {
		return nil, err
	}

	mw.InitDistributeNode(conf, log, repo)

	jmgr := jmysql.NewJobManager(log, repo, conf)
	if err = jmgr.Rebuild(); err != nil {
		return nil, err
	}

	panicFile, err := os.OpenFile(conf.Logger.PanicLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	if err = mw.InitAuthorizationAPIs(repo); err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(mw.InjectConfig(conf))
	r.Use(mw.InjectLogger(log))
	r.Use(mw.InjectRepo(repo))
	r.Use(mw.InjectFile(panicFile))
	r.Use(mw.InjectJobManager(jmgr))
	r.Use(mw.LogPanic)
	r.Use(mw.Authenticator) //用户认证中间件
	r.Use(mw.Authorization) //API授权中间件
	r.Use(mw.OperateInterceptor)

	if conf.DHCPLimiter.Enable {
		limiter.GlobalLimiter, err = webank.NewLimiter(log, repo, conf.DHCPLimiter.Limit)
		if err != nil {
			return nil, err
		}
		r.Use(mw.InjectDHCPLimiter(limiter.GlobalLimiter))
	}

	registerRoutes(r)

	return &Server{
		Conf:    conf,
		Log:     log,
		Repo:    repo,
		handler: r,
	}, nil
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.handler.ServeHTTP(w, r)
}
