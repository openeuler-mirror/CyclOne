package cloudbootserver

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/logger"
	myhttp "idcos.io/cloudboot/utils/http"
)

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func makeDirector(log logger.Logger, conf *config.Config, target *url.URL) func(*http.Request) {
	targetQuery := target.RawQuery
	return func(req *http.Request) {
		log.Debugf("%s %s", req.Method, req.URL.String())
		// 以下代码源于golang官方实现，谨慎修改。
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}

		// 在http header 'X-Forwarded-For'中记录每一级代理IP地址，
		// 并通过英文逗号分隔，最新一级的代理信息会被追加到末尾。
		if v, ok := req.Header[myhttp.XForwardedFor]; !ok || len(v) <= 0 {
			req.Header.Set(myhttp.XForwardedFor, conf.ReverseProxy.IP)
		} else {
			req.Header.Set(myhttp.XForwardedFor, fmt.Sprintf("%s,%s", v[0], conf.ReverseProxy.IP))
		}

		// 在http header 'X-Forwarded-Origin'中记录每一级代理的信息，
		// 并通过英文逗号分隔，最新一级的代理信息会被追加到末尾。
		if conf.ReverseProxy.Origin != "" {
			if v, ok := req.Header[myhttp.XForwardedOrigin]; !ok || len(v) <= 0 {
				req.Header.Set(myhttp.XForwardedOrigin, conf.ReverseProxy.Origin)
			} else {
				req.Header.Set(myhttp.XForwardedOrigin, fmt.Sprintf("%s,%s", v[0], conf.ReverseProxy.Origin))
			}
		}
	}
}

// NewReverseProxyServer 实例化http服务
func NewReverseProxyServer(log logger.Logger, conf *config.Config) (*Server, error) {
	targetURL, err := url.Parse(conf.ReverseProxy.URL)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if conf.ReverseProxy.IP == "" {
		conf.ReverseProxy.IP, _ = loadIP()
	}
	log.Debugf("Origin Node: %s, Origin Node IP: %s", conf.ReverseProxy.Origin, conf.ReverseProxy.IP)

	rproxy := httputil.NewSingleHostReverseProxy(targetURL)
	rproxy.Director = makeDirector(log, conf, targetURL)
	return &Server{
		Conf:    conf,
		Log:     log,
		handler: rproxy,
	}, nil
}

func loadIP() (addr string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for i := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := addrs[i].(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "127.0.0.1", nil
}
