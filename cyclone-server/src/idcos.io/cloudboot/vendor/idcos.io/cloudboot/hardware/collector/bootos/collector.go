package bootos

import (
	"strings"

	"idcos.io/cloudboot/hardware"
	"idcos.io/cloudboot/hardware/collector"
)

const (
	// collectorName 采集器名称
	collectorName = collector.DefaultCollector
)

func init() {
	// 注册默认的采集器实现
	collector.Register(collectorName, new(bootosC))
}

// bootosC BOOTOS中运行的采集器实现
type bootosC struct {
	hardware.Base
}

const (
	// colonSeparator 分隔符
	colonSeparator = ":"
	// eqSeprator 等号分隔符
	eqSeprator = "="
)

// extractValue 截取kv对中v的内容。假设，kv内容为"name : voidint"，那么将返回"voidint"。
func (c *bootosC) extractValue(kv, sep string) (value string) {
	if !strings.Contains(kv, sep) {
		return kv
	}
	return strings.TrimSpace(strings.SplitN(kv, sep, 2)[1])
}
