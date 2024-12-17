package collector

import (
	"strings"
	"sync"
)

// pool 采集器池
var pool = make(map[string]Collector)
var mux sync.Mutex

// Register 注册RAID阵列卡硬件及其处理Collector实例
func Register(name string, collector Collector) {
	mux.Lock()
	defer mux.Unlock()
	if collector == nil {
		panic("raid: Register Collector is nil")
	}
	name = strings.ToUpper(name)
	if _, dup := pool[name]; dup {
		panic("raid: Register called twice for Collector " + name)
	}
	pool[name] = collector
}

// Registered 返回已注册的RAID硬件
func Registered() (items []string) {
	mux.Lock()
	defer mux.Unlock()

	for key := range pool {
		items = append(items, key)
	}
	return items
}

// SelectCollector 根据RAID硬件名称获取相应的Collector
func SelectCollector(name string) Collector {
	mux.Lock()
	defer mux.Unlock()

	name = strings.ToUpper(name)
	for key := range pool {
		if key == name {
			return pool[key]
		}
	}
	return nil
}
