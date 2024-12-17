package webank

import (
	"fmt"
	"time"

	"sync"

	"idcos.io/cloudboot/limiter"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils"
)

type myBucket struct {
	log      logger.Logger
	repo     model.Repo
	id       string
	capacity int // 容量
	ch       chan limiter.Token
	mutex    sync.Mutex //加个锁
}

// newBucket 生成指定容量的令牌桶
func newBucket(log logger.Logger, repo model.Repo, bucketid string, capacity int, genToken limiter.GenTokenFunc) (*myBucket, error) {
	if capacity <= 0 {
		panic("capacity must be greater than 0")
	}

	ch := make(chan limiter.Token, capacity)
	items := make([]*model.DHCPTokenBucket, 0, capacity)

	for i := 0; i < capacity; i++ {
		token := genToken(bucketid)
		log.Debugf("new token:%s into bucket:%s", token, bucketid)
		items = append(items, &model.DHCPTokenBucket{
			Token:  string(token),
			Bucket: bucketid,
		})
		ch <- token
	}
	if err := repo.OverwriteTokenBuckets(items...); err != nil {
		close(ch)
		log.Errorf("Overwrite token bucket records error: %s", err.Error())
		return nil, err
	}
	return &myBucket{
		log:      log,
		repo:     repo,
		id:       bucketid,
		capacity: capacity,
		ch:       ch,
	}, nil
}

// reloadBucket 根据令牌桶名从数据库中加载令牌桶及其令牌
func reloadBucket(log logger.Logger, repo model.Repo, bucket string, capacity int) (*myBucket, error) {
	log.Infof("Start reloading the token bucket: %s", bucket)

	items, err := repo.GetTokenBuckets(&model.DHCPTokenBucket{
		Bucket: bucket,
	})
	if err != nil {
		return nil, err
	}

	ch := make(chan limiter.Token, capacity)

	if len(items) > capacity { // 缩容了
		log.Infof("bucket capacity become smaller,len(%d),limit(%d)", len(items), capacity)
		items = items[:capacity]

	} else if len(items) < capacity { // 扩容了
		log.Infof("bucket capacity become greater,len(%d),limit(%d)", len(items), capacity)
		for len(items) < capacity {
			items = append(items, &model.DHCPTokenBucket{
				Token:  string(genToken(bucket)),
				Bucket: bucket,
			})
		}
	}

	if len(items) != capacity {
		panic("unreachable")
	}

	for i := range items {
		sn := items[i].SN
		if sn != nil && *sn != "" {
			log.Infof("The token(%s) has been assigned to the device(%s)", items[i].Token, *sn)
			continue
		}
		ch <- limiter.Token(items[i].Token) // 使用新生成的令牌将令牌桶填满
		log.Infof("Put the token in the token bucket: %s -> %s", items[i].Token, bucket)
	}

	if err := repo.OverwriteTokenBuckets(items...); err != nil {
		return nil, err
	}

	return &myBucket{
		log:      log,
		repo:     repo,
		id:       bucket,
		capacity: capacity,
		ch:       ch,
	}, nil
}

// ID 返回令牌桶ID(TOR名称)
func (bucket *myBucket) ID() string {
	return bucket.id
}

// Capacity 返回令牌桶容量
func (bucket *myBucket) Capacity() int {
	return bucket.capacity
}

// Acquire 尝试从桶中获取令牌。
// 若桶中无可用令牌，则阻塞并等待。
// 若在超时时间内还没有可用的令牌，则返回ErrAcquireTokenTimeout错误。
func (bucket *myBucket) Acquire(sn string, timeout time.Duration) (limiter.Token, error) {
	bucket.log.Infof("The device(%s) starts requesting a token", sn)

	// 若设备已经被分配了令牌，则直接返回之前分发的令牌。
	if t, _ := bucket.repo.GetTokenBySN(sn); t != "" {
		bucket.log.Infof("The device(%s) reuses the existing token: %s", sn)
		return limiter.Token(t), nil
	}

	bucket.log.Infof("The remaining token in the this bucket is %d", len(bucket.ch))

	select {
	case token := <-bucket.ch:
		bucket.log.Infof("The device(%s) gets a token: %s", sn, token)

		// 将设备和令牌的对应关系写入数据库
		if _, err := bucket.repo.BindSNByTokenBucket(sn, string(token), bucket.id); err != nil {
			bucket.log.Warnf("Binding error: %s", err.Error())
			bucket.ch <- token
			return "", err
		}
		return token, nil

	case <-time.After(timeout):
		bucket.log.Warnf("The device(%s) request token timeout", sn)
		return "", limiter.ErrAcquireTokenTimeout
	}
}

// Return 归还令牌。
// 若令牌归还多次，或者归还的是无效的令牌，则返回ErrInvalidOrReturnedToken错误。
func (bucket *myBucket) Return(sn string, token limiter.Token) error {
	bucket.log.Infof("The device(%s) starts to return a token: %s", sn, token)

	// 校验令牌有效性并在数据库中解绑令牌和设备
	affected, err := bucket.repo.UnbindSNByTokenBucket(sn, string(token), bucket.id)
	if err != nil {
		bucket.log.Warnf("Unbinding error: %s", err.Error())
		return err
	}
	if affected == 0 {
		bucket.log.Error("Invalid or returned token, affected == 0")
		return limiter.ErrInvalidOrReturnedToken
	}
	select {
	case bucket.ch <- token:
	default:
		bucket.log.Errorf("Return token:%s fail, bucket:%s is full", token, bucket.ID())
		return limiter.ErrBucketFull
	}

	bucket.log.Infof("The device(%s) has returned a token: %s", sn, token)
	return nil
}

// Drop 销毁当前令牌桶
func (bucket *myBucket) Drop() error {
	bucket.log.Info("Start destroying the token bucket: %s", bucket.id)

	if _, err := bucket.repo.RemoveTokenBuckets(bucket.id); err != nil {
		return err
	}
	close(bucket.ch)

	bucket.log.Info("Token bucket(%s) has been destroyed", bucket.id)
	return nil
}

// 生成token。token格式为"${bucket}:${UUID}"
func genToken(bucket string) limiter.Token {
	return limiter.Token(fmt.Sprintf("%s:%s", bucket, utils.UUID()))
}
