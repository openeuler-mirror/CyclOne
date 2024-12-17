package middleware

import (
	"sync"

	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
)

// 定义一个全局的map结构来存储机房id到部署IP的映射
type MapDistribute struct {
	MDistribute map[uint][]string
	sync.Mutex
}

var MapDistributeNode MapDistribute

//去数据库查询所有的机房，然后和conf中的分布式节点比较，初始化map数据
func InitDistributeNode(conf *config.Config, log logger.Logger, repo model.Repo) {
	MapDistributeNode.MDistribute = make(map[uint][]string, 0)
	for _, confdn := range conf.DistributedNodes {
		sr, err := repo.GetServerRoomByName(confdn.NodeName)
		if err != nil {
			log.Errorf("init distribute node:%s fail,%v", confdn.NodeName, err)
			continue
		} else if sr != nil {
			MapDistributeNode.Lock()
			MapDistributeNode.MDistribute[sr.ID] = confdn.NodeIP
			MapDistributeNode.Unlock()
		}
	}
}
