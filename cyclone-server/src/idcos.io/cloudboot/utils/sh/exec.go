package sh

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"idcos.io/cloudboot/logger"
)

// ExecScriptWithBuf *nix系统下，执行scriptFile标识的脚本文件，并将脚本执行的标准输出写入stdout，标准错误输出写入stderr。
func ExecScriptWithBuf(scriptFile string, stdout, stderr *bytes.Buffer) (err error) {
	cmd := exec.Command("/bin/bash", scriptFile)
	if stdout != nil {
		cmd.Stdout = stdout
	}
	if stderr != nil {
		cmd.Stderr = stderr
	}
	return cmd.Run()
}

// ExecScriptCombinedOutput *nix系统下，执行scriptFile标识的脚本文件，并将脚本执行的标准输出和标准错误输出都通过字节切片output返回。
func ExecScriptCombinedOutput(scriptFile string) (output []byte, err error) {
	return exec.Command("/bin/bash", scriptFile).CombinedOutput()
}

// ExecWithBuf *nix系统下，执行命令字符串cmdAndArgs，并将命令执行的标准输出写入stdout，标准错误输出写入stderr。
func ExecWithBuf(log logger.Logger, cmdAndArgs string, stdout, stderr *bytes.Buffer) (err error) {
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
	if err = ExecScriptWithBuf(scriptFile, stdout, stderr); err != nil {
		return err
	}
	return nil
}

// ExecOutput 执行命令字符串cmdAndArgs，并将命令执行的标准输出和标准错误输出都通过字节切片output返回。
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
	if output, err = exec.Command("/bin/bash", scriptFile).Output(); err != nil {
		if log != nil {
			log.Error(err)
		}
	}
	return output, err
}

// ExecDesensitizeOutputWithLog 脱敏日志输出
func ExecDesensitizeOutputWithLog(log logger.Logger, cmdAndArgsUTF8, keyWord, newWord string) (outUTF8 []byte, err error) {

	scriptFile, err := GenTempScript([]byte(cmdAndArgsUTF8))

	defer os.Remove(scriptFile)

	if err != nil {
		if log != nil {
			log.Error(err)
		}
		return nil, err
	}
	output, err := ExecScriptCombinedOutput(scriptFile)

	desensitizedCmd := strings.Replace(cmdAndArgsUTF8, keyWord, newWord, -1)
	outDesensitized := strings.Replace(string(output), keyWord, newWord, -1)
	if err != nil {
		if log != nil {
			log.Errorf("Exec %q err: %s\noutput:\n%s", desensitizedCmd, err, outDesensitized)
		}
		return output, err
	}
	if log != nil {
		log.Infof("%s ==>\n%s", desensitizedCmd, outDesensitized)
	}
	return output, nil
}

func CmdDesensitization(cmd string) string {
	strs := strings.Split(cmd, " ")
	for i, v := range strs {
		if strings.EqualFold(strings.ToUpper(v), "-P") {
			cmd = strings.Replace(cmd, strs[i+1], "******", -1)
		}
		if strings.EqualFold(strings.ToLower(v), "password") && strings.Contains(cmd, "user set password") {
			cmd = strings.Replace(cmd, strs[i+2], "******", -1)
		}
	}
	return cmd
}

// ExecOutputWithLog 与ExecOutput函数功能相同，但是会输出更丰富的日志。
func ExecOutputWithLog(log logger.Logger, cmdAndArgsUTF8 string) (outUTF8 []byte, err error) {
	output, err := ExecOutput(log, cmdAndArgsUTF8)

	desensitizationCmd := CmdDesensitization(cmdAndArgsUTF8)
	desensitizationOp := CmdDesensitization(string(output))

	if err != nil {
		if log != nil {
			log.Errorf("Exec %q err: %s\noutput:\n%s", desensitizationCmd, err, desensitizationOp)
		}
		return output, err
	}
	if log != nil {
		log.Infof("[CMD-EXECUTED]==>\n%s", desensitizationCmd)
		//log.Debugf("[CMD-OUTPUT] ==>\n%s", desensitizationOp)
	}
	return output, nil
}

// GenTempScript 在系统临时目录生成可执行脚本文件
func GenTempScript(content []byte) (scriptFile string, err error) {
	scriptFile = filepath.Join(os.TempDir(), fmt.Sprintf("%s.sh", genUUID()))
	if err = ioutil.WriteFile(scriptFile, content, 0744); err != nil {
		return "", err
	}
	return scriptFile, nil
}

// genUUID 返回UUID字符串
func genUUID() string {
	return fmt.Sprintf("%d%d", time.Now().UnixNano(), time.Now().Nanosecond())
}
