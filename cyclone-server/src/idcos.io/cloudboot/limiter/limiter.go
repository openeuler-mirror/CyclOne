package limiter

// Limiter 装机(进入bootos)限流器
type Limiter interface {
	// Route 给目标设备路由到一个合适的令牌桶
	Route(sn string) (Bucket, error)
	// DropBucket 销毁目标令牌桶
	DropBucket(bucketID string) error
}
