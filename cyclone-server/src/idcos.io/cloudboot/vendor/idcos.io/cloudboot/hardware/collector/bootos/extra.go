package bootos

import (
	"encoding/json"
	"os/exec"

	"idcos.io/cloudboot/hardware/collector"
	"idcos.io/cloudboot/utils/sh"
)

// Extra 执行采集脚本并返回采集到的信息。
func (c *bootosC) Extra(scripts [][]byte) *collector.Extra {
	store := make(map[string]interface{})
	for i := range scripts {
		entries, _ := c.execScript(scripts[i])
		for k := range entries {
			store[k] = entries[k]
		}
	}
	e := collector.Extra(store)
	return &e
}

// execScript 将脚本内容写入临时文件，并赋予该临时文件执行权限，再执行。
// 读取脚本执行的JSON输出并反序列化成map，若脚本的标准输出非JSON格式，或者是非kv结构，则返回非nil的error。
func (c *bootosC) execScript(script []byte) (entries map[string]interface{}, err error) {
	filename, err := sh.GenTempScript(script)
	if err != nil {
		if log := c.Base.GetLog(); log != nil {
			log.Error(err)
		}
		return nil, err
	}

	if log := c.Base.GetLog(); log != nil {
		log.Debugf("%s ==>\n%s\n", filename, script)
	}

	// 不能调用某个具体程序(bash、python)去执行该脚本，因为目标脚本的编写语言并不确定。
	// 脚本的编写者应该在shebang中指定执行该脚本的程序，如'#! /usr/bin/env python'。
	cmd := exec.Command(filename)
	output, err := cmd.Output()
	if err != nil {
		if log := c.Base.GetLog(); log != nil {
			log.Errorf("Exec err: %s\noutput:\n%s", err, output)
		}
		return nil, err
	}
	if log := c.Base.GetLog(); log != nil {
		log.Debugf("Exec output ==>\n%s", output)
	}
	return entries, json.Unmarshal(output, &entries)
}
