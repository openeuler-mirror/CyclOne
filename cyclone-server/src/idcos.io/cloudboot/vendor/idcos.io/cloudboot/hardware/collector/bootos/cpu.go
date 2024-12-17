package bootos

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"

	"idcos.io/cloudboot/hardware/collector"
)

// CPU 采集并返回当前设备的CPU信息
func (c *bootosC) CPU() (*collector.CPU, error) {
	var cpu collector.CPU
	var err error

	cpu.Physicals, err = c.physicalCPUs()
	if err != nil {
		return nil, err
	}
	cpu.TotalPhysicals = len(cpu.Physicals)
	for i := range cpu.Physicals {
		cpu.TotalCores += cpu.Physicals[i].Cores
	}
	cpu.TotalLogicals, err = c.countLogicalCPUs()
	if err != nil {
		return nil, err
	}
	cpu.Threads, err = c.cpuThreads()
	if err != nil {
		return nil, err
	}
	return &cpu, nil
}

// physicalCPUs 返回物理CPU列表
func (c *bootosC) physicalCPUs() (items []collector.PhysicalCPU, err error) {
	output, err := c.ExecByShell("dmidecode", "-t", "processor")
	if err != nil {
		if log := c.Base.GetLog(); log != nil {
			log.Error(err)
		}
		return nil, err
	}

	procs := bytes.Split(output, []byte{'\n', '\n'})
	for i := range procs {
		if !bytes.Contains(procs[i], []byte("Processor Information")) {
			continue
		}
		cpu, err := c.parsePhysicalCPU(procs[i])
		if err != nil {
			return nil, err
		}
		if cpu == nil {
			continue
		}
		items = append(items, *cpu)
	}
	return items, nil
}

// parsePhysicalCPU 返回单个物理CPU及其对应的physical id
func (c *bootosC) parsePhysicalCPU(out []byte) (pcpu *collector.PhysicalCPU, err error) {
	var cpu collector.PhysicalCPU
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.Contains(line, ":") {
			continue
		}

		if strings.HasPrefix(line, "Core Count") {
			cpu.Cores, _ = strconv.Atoi(c.extractValue(line, colonSeparator))

		} else if strings.HasPrefix(line, "Version") {
			cpu.ModelName = c.extractValue(line, colonSeparator)

		} else if strings.HasPrefix(line, "Max Speed") {
			cpu.ClockSpeed = c.extractValue(line, colonSeparator)
		}
	}
	if cpu.ModelName == "" {
		return nil, scanner.Err()
	}
	return &cpu, scanner.Err()
}

// countLogicalCPUs 返回逻辑CPU个数
func (c *bootosC) countLogicalCPUs() (int, error) {
	// cat /proc/cpuinfo| grep "processor"| wc -l
	output, err := c.Base.ExecByShell("cat", "/proc/cpuinfo", "|", "grep", `"processor"`, "|", "wc", "-l")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(output)))
}

// cpuThreads 返回CPU超线程数
func (c *bootosC) cpuThreads() (int, error) {
	// lscpu | grep "Thread(s) per core:"
	output, err := c.Base.ExecByShell("lscpu", "|", "grep", `"Thread(s) per core:"`)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(string(output), "Thread(s) per core:")))
}
