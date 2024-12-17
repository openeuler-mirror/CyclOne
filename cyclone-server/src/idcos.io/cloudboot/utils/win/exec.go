package win

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/utils"
)

// ExecScript windows系统下，执行scriptFile标识的脚本文件，并将脚本执行的标准输出写入stdout，标准错误输出写入stderr。
func ExecScript(scriptFile string, stdout, stderr *bytes.Buffer) (err error) {
	cmd := exec.Command("cmd", "/c", scriptFile)
	if stdout != nil {
		cmd.Stdout = stdout
	}
	if stderr != nil {
		cmd.Stderr = stderr
	}
	return cmd.Run()
}

// ExecScriptCombinedOutput windows系统下，执行scriptFile标识的脚本文件，并将脚本执行的标准输出和标准错误输出都通过字节切片output返回。
func ExecScriptCombinedOutput(scriptFile string) (output []byte, err error) {
	return exec.Command("cmd", "/c", scriptFile).CombinedOutput()
}

// Exec windows系统下，执行命令字符串cmdAndArgs，并将命令执行的标准输出写入stdout，标准错误输出写入stderr。
func Exec(log logger.Logger, cmdAndArgs string, stdout, stderr *bytes.Buffer) (err error) {
	scriptFile, err := GenTempScript([]byte(cmdAndArgs))
	if err != nil {
		if log != nil {
			log.Error(err)
		}
		return err
	}
	defer os.Remove(scriptFile)

	if log != nil {
		log.Infof("Exec script %s", scriptFile)
	}
	if err = ExecScript(scriptFile, stdout, stderr); err != nil {
		return err
	}
	return nil
}

// ExecOutput windows系统下，执行命令字符串cmdAndArgs，并将命令执行的标准输出和标准错误输出都通过字节切片output返回。
func ExecOutput(log logger.Logger, cmdAndArgs string) (output []byte, err error) {
	scriptFile, err := GenTempScript([]byte(cmdAndArgs))
	if err != nil {
		if log != nil {
			log.Error(err)
		}
		return nil, err
	}
	defer os.Remove(scriptFile)

	if log != nil {
		log.Infof("Exec script %s", scriptFile)
	}
	if output, err = exec.Command("cmd", "/c", scriptFile).Output(); err != nil {
		if log != nil {
			log.Error(err)
		}
	}
	return output, err
}

// ExecOutputWithLog 与ExecOutput函数功能相同，但是会输出更丰富的日志。
func ExecOutputWithLog(log logger.Logger, cmdAndArgsUTF8 string) (outUTF8 []byte, err error) {
	output, err := ExecOutput(log, cmdAndArgsUTF8)
	if err != nil {
		if log != nil {
			log.Errorf("Exec %q err: %s\noutput:\n%s", cmdAndArgsUTF8, err, string(output))
		}
		return nil, err
	}
	outputUTF8 := utils.GBK2UTF8(string(output))
	if log != nil {
		log.Infof("%s ==>\n%s", cmdAndArgsUTF8, outputUTF8)
	}
	return []byte(outputUTF8), nil
}

// ExecOneByOneWithLog 将脚本文件内容按照分隔符拆分成若干条语句并逐条执行
func ExecOneByOneWithLog(log logger.Logger, scriptBody, sep string, continueOnError bool) {
	scripts := strings.Split(scriptBody, sep)
	for i := range scripts {
		scripts[i] = strings.TrimSpace(scripts[i])
		if scripts[i] == "" {
			continue
		}
		_, err := ExecOutputWithLog(log, scripts[i])
		if err == nil {
			continue
		}
		if !continueOnError {
			log.Warnf("An error occurred and terminated")
			break
		}
	}
}

// GenTempScript 在系统临时目录生成bat脚本文件
func GenTempScript(content []byte) (scriptFile string, err error) {
	scriptFile = filepath.Join(os.TempDir(), fmt.Sprintf("%s.bat", utils.UUID()))
	if err = ioutil.WriteFile(scriptFile, content, 0744); err != nil {
		return "", err
	}
	return scriptFile, nil
}
