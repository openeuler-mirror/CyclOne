package limiter

import (
	"errors"
	"time"
)

var (
	// ErrInvalidOrReturnedToken 无效的或已经归还的令牌
	ErrInvalidOrReturnedToken = errors.New("invalid or returned token")
	// ErrAcquireTokenTimeout 等待获取令牌超时错误
	ErrAcquireTokenTimeout = errors.New("acquire token timeout")
	// ErrBucketNotFound 令牌桶未发现
	ErrBucketNotFound = errors.New("bucket not found")
	// ErrBucketFull 令牌桶满了
	ErrBucketFull = errors.New("bucket is full")
	// 1 定义全局变量
	// 2 server.go 初始化
	// 3 auto_deploy.go 调用
	GlobalLimiter Limiter
)

// Token 令牌(唯一的不重复字符串)
type Token string

// GenTokenFunc 令牌生成函数
type GenTokenFunc func(str string) Token

// Bucket 令牌桶
type Bucket interface {
	// ID 返回令牌桶ID
	ID() string
	// Capacity 返回令牌桶容量
	Capacity() int
	// Acquire 尝试从桶中获取令牌。
	// 若桶中无可用令牌，则阻塞并等待。
	// 若在超时时间内还没有可用的令牌，则返回ErrAcquireTokenTimeout错误。
	Acquire(sn string, timeout time.Duration) (Token, error)
	// Return 归还令牌。
	// 若令牌归还多次，或者归还的是无效的令牌，则返回ErrInvalidOrReturnedToken错误。
	Return(sn string, token Token) error
	// Drop 销毁当前令牌桶
	Drop() error
}
