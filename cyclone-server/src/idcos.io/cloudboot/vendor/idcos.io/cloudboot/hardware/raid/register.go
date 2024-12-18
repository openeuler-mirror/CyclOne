package raid

import (
	"strings"
	"sync"
)

// workerPool 处理RAID工作的工人池
var workerPool = make(map[string]Worker)
var mux sync.Mutex

// Register 注册RAID阵列卡硬件及其处理worker实例
func Register(name string, worker Worker) {
	mux.Lock()
	defer mux.Unlock()
	if worker == nil {
		panic("raid: Register worker is nil")
	}
	name = strings.ToUpper(name)
	if _, dup := workerPool[name]; dup {
		panic("raid: Register called twice for worker " + name)
	}
	workerPool[name] = worker
}

// Registered 返回已注册的RAID硬件
func Registered() (items []string) {
	mux.Lock()
	defer mux.Unlock()

	for key := range workerPool {
		items = append(items, key)
	}
	return items
}

// SelectWorker 根据RAID硬件名称获取相应的Worker
func SelectWorker(name string) Worker {
	mux.Lock()
	defer mux.Unlock()

	name = strings.ToUpper(name)
	for key := range workerPool {
		if key == name {
			return workerPool[key]
		}
	}
	return nil
}
