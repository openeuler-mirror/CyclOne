package uam

import "fmt"

const (
	// Version 库版本号
	Version = "0.2.0"
)

// Client UAM客户端
type Client struct {
	log  Logger
	conf *Config
}

// NewClient 实例化
func NewClient(rootEndpoint, token string, options ...SetOptionFunc) *Client {
	c := Config{
		RootEndpoint: rootEndpoint,
		Token:        token,
	}
	for _, f := range options {
		f(&c)
	}
	c.LoadDefault()
	return &Client{
		log:  defaultLog,
		conf: &c,
	}
}

const (
	success = "success"
)

const (
	tokenHeaderName       = "access-token"
	acceptHeaderName      = "Accept"
	contentTypeHeaderName = "Content-Type"
	appIDParam            = "appId"
	tokenParam            = "token"
)

const (
	jsonMediaType = "application/json"
)

func (cli *Client) acceptJSONHeader() (string, string) {
	return acceptHeaderName, jsonMediaType
}

func (cli *Client) contentTypeJSONHeader() (string, string) {
	return contentTypeHeaderName, jsonMediaType
}

func (cli *Client) tokenHeader() (string, string) {
	return tokenHeaderName, fmt.Sprintf("Bearer %s", cli.conf.Token)
}

// PageLimiter 分页限制
type PageLimiter struct {
	PageNo   int64
	PageSize int64
}

func genDefaultPageLimiter() *PageLimiter {
	return &PageLimiter{
		PageNo:   1,
		PageSize: 10,
	}
}
