package hardware

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"idcos.io/cloudboot/logger"
)

// Base 硬件基础结构体
type Base struct {
	log   logger.Logger
	debug bool
}

// SetDebug 设置是否开启debug。若开启debug，会将关键日志信息写入console。
func (base *Base) SetDebug(debug bool) {
	base.debug = debug
}

// SetLog 更换日志实现。默认情况下内部无日志实现。
func (base *Base) SetLog(log logger.Logger) {
	base.log = log
}

// GetLog 返回日志
func (base *Base) GetLog() logger.Logger {
	return base.log
}

// Sleep 休眠2s
func (base *Base) Sleep() {
	time.Sleep(2 * time.Second)
}

// Exec 本地命令执行并返回stdout内容
func (base *Base) Exec(cmd string, args ...string) (output []byte, err error) {
	if base.log != nil {
		base.log.Debugf("==> %s %s", cmd, strings.Join(args, " "))
		if base.debug {
			fmt.Printf("==> %s %s\n", cmd, strings.Join(args, " "))
		}
	}

	outputInfo := exec.Command(cmd, args...)
	output, err = outputInfo.CombinedOutput() //返回标准输出和标准错误
	if output != nil && base.log != nil {
		base.log.Debugf("\n--------------------stdout/stderr begin--------------------\n%s\n--------------------stdout/stderr end--------------------", string(output))
		if base.debug {
			fmt.Printf("--------------------stdout/stderr begin--------------------\n%s\n--------------------stdout/stderr end--------------------\n", string(output))
		}
	}

	if err != nil && base.log != nil {
		base.log.Debugf(err.Error())
		err = errors.New(string(output))
	}

	return output, err
}

const (
	shell = "/bin/bash"
)

// ExecByShell 通过本地shell运行命令
func (base *Base) ExecByShell(cmd string, args ...string) (output []byte, err error) {
	if base.log != nil {
		base.log.Debugf("==> %s %s", cmd, strings.Join(args, " "))
		if base.debug {
			fmt.Printf("==> %s %s\n", cmd, strings.Join(args, " "))
		}
	}

	scriptFile, err := genTempScript([]byte(fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))))
	if err != nil {
		if base.log != nil {
			base.log.Error(err)
		}
		return nil, err
	}
	defer os.Remove(scriptFile)

	outputInfo := exec.Command(shell, scriptFile)
	output, err = outputInfo.CombinedOutput() //返回标准输出和标准错误
	if output != nil && base.log != nil {
		base.log.Debugf("\n--------------------stdout/stderr begin--------------------\n%s\n--------------------stdout/stderr end--------------------", string(output))
		if base.debug {
			fmt.Printf("--------------------stdout/stderr begin--------------------\n%s\n--------------------stdout/stderr end--------------------\n", string(output))
		}
	}

	if err != nil && base.log != nil {
		base.log.Debugf(err.Error())
		err = errors.New(string(output))
	}

	return output, err
}

// genTempScript 在系统临时目录生成可执行脚本文件
func genTempScript(content []byte) (scriptFile string, err error) {
	scriptFile = filepath.Join(os.TempDir(), fmt.Sprintf("%d.sh", time.Now().UnixNano()))
	if err = ioutil.WriteFile(scriptFile, content, 0744); err != nil {
		return "", err
	}
	return scriptFile, nil
}

// ExecWithPipe 完成带管道符的命令输出 比如命令[ls | grep Helloworld],
// 需要将exec.Command("ls")和exec.Command("grep", "Helloworld")的结果分别作为参数传入
func (base *Base) ExecWithPipe(cmds ...*exec.Cmd) ([]byte, error) {
	// At least one command
	if len(cmds) < 1 {
		return nil, nil
	}

	var stdout bytes.Buffer
	//var stderr bytes.Buffer
	var err error
	maxIndex := len(cmds) - 1
	cmds[maxIndex].Stdout = &stdout
	cmds[maxIndex].Stderr = &stdout

	//for i, cmd := range cmds[:maxIndex] {
	for i := 0; i < maxIndex; i++ {
		cmds[i+1].Stdin, err = cmds[i].StdoutPipe()
		if err != nil {
			return nil, err
		}
	}
	// Start each command
	for _, cmd := range cmds {
		err := cmd.Start()
		if err != nil {
			return stdout.Bytes(), err
		}
	}
	// Wait for each command to complete
	for _, cmd := range cmds {
		err := cmd.Wait()
		if err != nil {
			return stdout.Bytes(), err
		}
	}
	return stdout.Bytes(), nil
}

// ExecWithStdinPipe 本地命令执行并读取stdin内容后返回stdout内容。
func (base *Base) ExecWithStdinPipe(inputs []string, cmd string, args ...string) (output []byte, err error) {
	if base.log != nil {
		base.log.Debugf("==> %s %s", cmd, strings.Join(args, " "))
		if base.debug {
			fmt.Printf("==> %s %s\n", cmd, strings.Join(args, " "))
		}
	}

	command := exec.Command(cmd, args...)
	stdin, err := command.StdinPipe()
	if err != nil {
		return nil, err
	}

	go func() {
		defer stdin.Close()
		for i := range inputs {
			fmt.Fprintln(stdin, inputs[i])
			if i < len(inputs)-1 {
				time.Sleep(time.Second)
			}
		}
	}()

	output, err = command.Output()
	if output != nil && base.log != nil {
		base.log.Debugf("\n--------------------stdout begin--------------------\n%s\n--------------------stdout end--------------------", string(output))
		if base.debug {
			fmt.Printf("--------------------stdout begin--------------------\n%s\n--------------------stdout end--------------------\n", string(output))
		}
	}
	if err != nil && base.log != nil {
		base.log.Errorf("\n--------------------stderr begin--------------------\n%s\n--------------------stderr end--------------------", err.Error())
		if base.debug {
			fmt.Fprintf(os.Stderr, "--------------------stderr begin--------------------\n%s\n--------------------stderr end--------------------\n", err.Error())
		}
		return nil, err
	}
	return output, err
}

const (
	// MatchedYES 实际值与预期值相符
	MatchedYES = "yes"
	// MatchedNO 实际值与预期值不符
	MatchedNO = "no"
	// MatchedUnknown 实际值与预期值是否相符未知
	MatchedUnknown = "unknown"
)

// CheckingItem 硬件配置实施后置检查结果
type CheckingItem struct {
	Title    string `json:"title"`    // 检查项名称
	Expected string `json:"expected"` // 预期值
	Actual   string `json:"actual"`   // 实际值
	Matched  string `json:"matched"`  // 是否匹配
	Error    string `json:"error"`    // 检查过程出错信息
}

// CheckingResult 硬件配置实施后置检查结果集合
type CheckingResult struct {
	RAIDItems []CheckingItem `json:"raid"`
	OOBItems  []CheckingItem `json:"oob"`
	BIOSItems []CheckingItem `json:"bios"`
	FWItems   []CheckingItem `json:"firmware"`
	//...继续扩展
}

// MatchFunc 对比预期值与实际值是否匹配的匹配器
type MatchFunc func(expected, actual string) bool

// EqualMatch 相等匹配器
func EqualMatch(expected, actual string) bool {
	return expected == actual
}

// EqualIgnoreCaseMatch 大小写不敏感相等匹配器
func EqualIgnoreCaseMatch(expected, actual string) bool {
	return strings.ToLower(expected) == strings.ToLower(actual)
}

// ContainsMatch 内容包含匹配器
func ContainsMatch(expected, actual string) bool {
	return strings.Contains(actual, expected)
}

// CheckingHelper 硬件配置实施检查帮手。
// 检查帮手默认使用相等匹配模式。
type CheckingHelper struct {
	item    *CheckingItem
	matcher MatchFunc
}

// NewCheckingHelper 实例化检查帮手。默认使用相等匹配器。
func NewCheckingHelper(title, expected, actual string) *CheckingHelper {
	return &CheckingHelper{
		item: &CheckingItem{
			Title:    title,
			Expected: expected,
			Actual:   actual,
		},
		matcher: EqualMatch,
	}
}

// Matcher 重置匹配器
func (h *CheckingHelper) Matcher(matcher MatchFunc) *CheckingHelper {
	h.matcher = matcher
	return h
}

// Do 执行检查并返回检查项结果。
func (h *CheckingHelper) Do() *CheckingItem {
	if h.matcher(h.item.Expected, h.item.Actual) {
		h.item.Matched = MatchedYES
	} else {
		h.item.Matched = MatchedNO
	}
	return h.item
}
