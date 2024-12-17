package http

import (
	"net/http"
	"strings"
)

const (
	// XForwardedOrigin 自定义http header，用于记录经过的各级代理信息
	XForwardedOrigin = "X-Forwarded-Origin"
	// XForwardedFor 自定义http header，用于记录经过的各级代理IP
	XForwardedFor = "X-Forwarded-For"
)

// ExtractOriginNodeWithDefault 从http request中提取源代理节点(第一级代理节点)信息，若值为空，则返回默认值。
func ExtractOriginNodeWithDefault(r *http.Request, def string) string {
	if v := ExtractOriginNode(r); v != "" {
		return v
	}
	return def
}

// ExtractOriginNode 从http request中提取源代理节点(第一级代理节点)信息
func ExtractOriginNode(r *http.Request) string {
	return strings.Split(r.Header.Get(XForwardedOrigin), ",")[0]
}

// ExtractOriginNodeIP 从http request中提取源代理节点(第一级代理节点)IP地址
func ExtractOriginNodeIP(r *http.Request) string {
	v := r.Header.Get(XForwardedFor)
	if idx := strings.Index(v, ","); idx > 0 {
		return strings.TrimSpace(v[:idx])
	}
	return v
}
