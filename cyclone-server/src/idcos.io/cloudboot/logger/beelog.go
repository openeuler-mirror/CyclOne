package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/astaxie/beego/logs"

	"idcos.io/cloudboot/config"
	myfilepath "idcos.io/cloudboot/utils/filepath"
)

// BeeLogger beego log实现
type BeeLogger struct {
	fileLog    *logs.BeeLogger
	consoleLog *logs.BeeLogger
}

func selectLevel(level string) int {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return 7
	case "warn":
		return 4
	case "error":
		return 3
	default:
		return 6 // default level: info
	}
}

// NewBeeLogger 创建BeeLogger实例
func NewBeeLogger(conf *config.Logger) *BeeLogger {
	var beeLog BeeLogger

	if conf.FilePerm == "" {
		conf.FilePerm = "0644"
	}

	// 开启文件日志
	if filename := strings.TrimSpace(conf.LogFile); filename != "" {
		filename, _ = myfilepath.Rel2Abs(filename)
		if err := os.MkdirAll(path.Dir(filename), os.ModePerm); err != nil {
			fmt.Printf("MkdirAll err: %s\n", err)
		}

		fLog := logs.NewLogger(1000)
		fLog.EnableFuncCallDepth(true) // 输出文件名和行号
		fLog.SetLogFuncCallDepth(3)
		logData, _ := json.Marshal(map[string]interface{}{
			"filename": filename,
			"perm":     conf.FilePerm,
			"level":    selectLevel(conf.Level),
			"rotate":   conf.RotateEnabled,
			"daily":    conf.RotateEnabled,
		})
		if err := fLog.SetLogger(logs.AdapterFile, string(logData)); err != nil {
			fmt.Fprintf(os.Stderr, "SetLogger err: %s\n", err)
		}
		beeLog.fileLog = fLog
	}

	// 开启控制台日志
	if conf.ConsoleEnabled {
		consoleLog := logs.NewLogger(1000)
		consoleLog.EnableFuncCallDepth(true)
		consoleLog.SetLogFuncCallDepth(3)
		logData, _ := json.Marshal(map[string]interface{}{
			"level": selectLevel(conf.Level),
		})
		if err := consoleLog.SetLogger(logs.AdapterConsole, string(logData)); err != nil {
			fmt.Fprintf(os.Stderr, "SetLogger err: %s\n", err)
		}
		beeLog.consoleLog = consoleLog
	}

	return &beeLog
}

// SetField 暂不支持的方法
func (log *BeeLogger) SetField(name string, value interface{}) {
	panic("not supported")
}

// Debug logs a debug message. If last parameter is a map[string]string, it's content
// is added as fields to the message.
func (log *BeeLogger) Debug(v ...interface{}) {
	if log.fileLog != nil {
		log.fileLog.Debug("%v", v...)
	}
	if log.consoleLog != nil {
		log.consoleLog.Debug("%v", v...)
	}
}

// Debugf logs a debug message with format. If last parameter is a map[string]string,
// it's content is added as fields to the message.
func (log *BeeLogger) Debugf(format string, v ...interface{}) {
	if log.fileLog != nil {
		log.fileLog.Debug(format, v...)
	}
	if log.consoleLog != nil {
		log.consoleLog.Debug(format, v...)
	}
}

// Info logs a info message. If last parameter is a map[string]string, it's content
// is added as fields to the message.
func (log *BeeLogger) Info(v ...interface{}) {
	if log.fileLog != nil {
		log.fileLog.Info("%v", v...)
	}
	if log.consoleLog != nil {
		log.consoleLog.Info("%v", v...)
	}
}

// Infof logs a info message with format. If last parameter is a map[string]string,
// it's content is added as fields to the message.
func (log *BeeLogger) Infof(format string, v ...interface{}) {
	if log.fileLog != nil {
		log.fileLog.Info(format, v...)
	}
	if log.consoleLog != nil {
		log.consoleLog.Info(format, v...)
	}
}

// Warn logs a warning message. If last parameter is a map[string]string, it's content
// is added as fields to the message.
func (log *BeeLogger) Warn(v ...interface{}) {
	if log.fileLog != nil {
		log.fileLog.Warn("%v", v...)
	}
	if log.consoleLog != nil {
		log.consoleLog.Warn("%v", v...)
	}
}

// Warnf logs a warning message with format. If last parameter is a map[string]string,
// it's content is added as fields to the message.
func (log *BeeLogger) Warnf(format string, v ...interface{}) {
	if log.fileLog != nil {
		log.fileLog.Warn(format, v...)
	}
	if log.consoleLog != nil {
		log.consoleLog.Warn(format, v...)
	}
}

// Error logs an error message. If last parameter is a map[string]string, it's content
// is added as fields to the message.
func (log *BeeLogger) Error(v ...interface{}) {
	if log.fileLog != nil {
		log.fileLog.Error("%v", v...)
	}
	if log.consoleLog != nil {
		log.consoleLog.Error("%v", v...)
	}
}

// Errorf logs an error message with format. If last parameter is a map[string]string,
// it's content is added as fields to the message.
func (log *BeeLogger) Errorf(format string, v ...interface{}) {
	if log.fileLog != nil {
		log.fileLog.Error(format, v...)
	}
	if log.consoleLog != nil {
		log.consoleLog.Error(format, v...)
	}
}

// Print 向目标打印输出日志
func (log *BeeLogger) Print(v ...interface{}) {
	var format string
	for range v {
		format += "%v "
	}
	if log.fileLog != nil {
		log.fileLog.Debug(format, v...)
	}
	if log.consoleLog != nil {
		log.consoleLog.Debug(format, v...)
	}
}
