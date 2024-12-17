package httpserver

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"idcos.io/cloudboot/logger"
	mw "idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/hwserver/config"
)

// Server HTTP server
type Server struct {
	conf    *config.Configuration
	log     logger.Logger
	handler http.Handler
}

// NewServer 返回HTTP server实例
func NewServer(conf *config.Configuration, log logger.Logger) (srv *Server, err error) {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(mw.InjectLogger(log))
	r.Use(mw.InjectHWConfig(conf))

	registerHandlers(r)

	return &Server{
		conf:    conf,
		log:     log,
		handler: r,
	}, nil
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.handler.ServeHTTP(w, r)
}
