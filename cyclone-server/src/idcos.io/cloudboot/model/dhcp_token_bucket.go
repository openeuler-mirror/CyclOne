package model

import "time"

// DHCPTokenBucket 令牌桶
type DHCPTokenBucket struct {
	Token     string    `gorm:"column:token"`
	Bucket    string    `gorm:"column:bucket"`
	SN        *string   `gorm:"column:sn"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// TableName 指定数据库表名
func (DHCPTokenBucket) TableName() string {
	return "dhcp_token_bucket"
}

// IDHCPTokenBucket 令牌桶操作接口
type IDHCPTokenBucket interface {
	// AddTokenBuckets 新增令牌记录
	AddTokenBuckets(items ...*DHCPTokenBucket) (err error)
	// OverwriteTokenBuckets 覆写令牌桶及其令牌
	OverwriteTokenBuckets(items ...*DHCPTokenBucket) (err error)
	// BindSNByTokenBucket 绑定SN与令牌桶中令牌
	BindSNByTokenBucket(sn, token, bucket string) (affected int64, err error)
	// UnbindSNByTokenBucket 解绑SN与令牌桶中令牌
	UnbindSNByTokenBucket(sn, token, bucket string) (affected int64, err error)
	// GetTokenBuckets 返回满足过滤条件的令牌桶及其令牌
	GetTokenBuckets(cond *DHCPTokenBucket) (items []*DHCPTokenBucket, err error)
	// GetUnbindingTokensByBucket 返回未绑定的SN的令牌列表
	GetUnbindingTokensByBucket(bucket string) (tokens []string, err error)
	// GetBuckets 返回当前所有令牌桶
	GetBuckets() (buckets []string, err error)
	// GetTokenBySN 返回目标设备在令牌桶内的令牌
	GetTokenBySN(sn string) (token string, err error)
	// RemoveTokenBuckets 移除指定名称的令牌桶及其令牌
	RemoveTokenBuckets(buckets ...string) (affected int64, err error)
}
