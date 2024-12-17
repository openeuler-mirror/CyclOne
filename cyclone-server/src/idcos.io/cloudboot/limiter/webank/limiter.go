package webank

import (
	"strings"
	"sync"

	"idcos.io/cloudboot/limiter"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
)

// myLimiter 限流器
type myLimiter struct {
	log          logger.Logger
	repo         model.Repo
	limit        int                  // 每个令牌桶中令牌数量阈值
	buckets      map[string]*myBucket // tor名称与令牌桶键值对
	bucketsMutex *sync.Mutex          // TODO 改造成读写锁
}

// NewLimiter 返回限流器实例
func NewLimiter(log logger.Logger, repo model.Repo, limit int) (limiter.Limiter, error) {
	if limit <= 0 {
		panic("invalid 'limit' parameter")
	}
	lim := &myLimiter{
		log:          log,
		repo:         repo,
		limit:        limit,
		buckets:      make(map[string]*myBucket),
		bucketsMutex: new(sync.Mutex),
	}
	if err := lim.init(); err != nil {
		return nil, err
	}
	return lim, nil
}

// init 初始化限流器实例
func (lim *myLimiter) init() error {
	lim.log.Infof("The limiter starts to initialize")

	// 1、从数据库读取已有令牌桶
	buckets, err := lim.repo.GetBuckets()
	if err != nil {
		lim.log.Errorf("Fetch buckets from db error: %s", err.Error())
		return err
	}

	// 2、重建限流器内的令牌桶及令牌
	for i := range buckets {
		bucket := strings.TrimSpace(buckets[i])
		if lim.buckets[bucket], err = reloadBucket(lim.log, lim.repo, bucket, lim.limit); err != nil {
			lim.log.Errorf("Reload bucket(%s) error: %s", bucket, err.Error())
			return err
		}
		// TODO 探活（ping设备的dhcp ip）
	}
	lim.log.Infof("The limiter is initialized")
	return nil
}

// Route 给目标设备路由到一个合适的令牌桶
func (lim *myLimiter) Route(sn string) (limiter.Bucket, error) {
	lim.log.Infof("The device(%s) starts looking for token bucket", sn)
	defer lim.log.Infof("The device(%s) has found token bucket", sn)

	tor, _ := lim.repo.GetTORBySN(sn)
	if tor == "" {
		lim.log.Errorf("The device(%s) failed to find TOR", sn)
		return nil, limiter.ErrBucketNotFound
	}
	lim.log.Infof("The device(%s) has found TOR: %s", sn, tor)

	lim.bucketsMutex.Lock()
	defer lim.bucketsMutex.Unlock()

	var err error

	bucket := lim.buckets[tor]
	if bucket == nil {
		lim.log.Infof("The device(%s) failed to get token bucket and create a token bucket", sn)
		bucket, err = newBucket(lim.log, lim.repo, tor, lim.limit, genToken)
		if err != nil {
			return nil, err
		}
		lim.buckets[tor] = bucket
	}
	return bucket, nil
}

// DropBucket 删除目标令牌桶
func (lim *myLimiter) DropBucket(bucketID string) error {
	lim.log.Infof("Start destroying the token bucket: %s", bucketID)

	lim.bucketsMutex.Lock()
	defer lim.bucketsMutex.Unlock()

	b := lim.buckets[bucketID]
	if b == nil {
		lim.log.Error("Token bucket not found")
		return limiter.ErrBucketNotFound
	}
	if err := b.Drop(); err != nil {
		return err
	}
	delete(lim.buckets, bucketID)

	lim.log.Infof("The token bucket has been destroyed: %s", bucketID)
	return nil
}
