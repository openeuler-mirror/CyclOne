package uam

import "github.com/astaxie/beego/logs"

// Logger 系统logger接口
type Logger interface {
	// Debug logs a debug message. If last parameter is a map[string]string, it's content
	// is added as fields to the message.
	Debug(v ...interface{})
	// Debug logs a debug message with format. If last parameter is a map[string]string,
	// it's content is added as fields to the message.
	Debugf(format string, v ...interface{})
	// Info logs a info message. If last parameter is a map[string]string, it's content
	// is added as fields to the message.
	Info(v ...interface{})
	// Info logs a info message with format. If last parameter is a map[string]string,
	// it's content is added as fields to the message.
	Infof(format string, v ...interface{})
	// Warn logs a warning message. If last parameter is a map[string]string, it's content
	// is added as fields to the message.
	Warn(v ...interface{})
	// Warn logs a warning message with format. If last parameter is a map[string]string,
	// it's content is added as fields to the message.
	Warnf(format string, v ...interface{})
	// Error logs an error message. If last parameter is a map[string]string, it's content
	// is added as fields to the message.
	Error(v ...interface{})
	// Error logs an error message with format. If last parameter is a map[string]string,
	// it's content is added as fields to the message.
	Errorf(format string, v ...interface{})
}

// defaultLog 默认日志(输出到控制台)
var defaultLog = newBeeLog()

type beeLog struct {
	l *logs.BeeLogger
}

func newBeeLog() Logger {
	l := logs.NewLogger(1000)
	l.SetLogger(logs.AdapterConsole)
	return &beeLog{
		l: l,
	}
}

func (log *beeLog) Debug(v ...interface{}) {
	log.l.Debug("%v", v...)
}

func (log *beeLog) Debugf(format string, v ...interface{}) {
	log.l.Debug(format, v...)
}

func (log *beeLog) Info(v ...interface{}) {
	log.l.Info("%v", v...)
}

func (log *beeLog) Infof(format string, v ...interface{}) {
	log.l.Info(format, v...)
}

func (log *beeLog) Warn(v ...interface{}) {
	log.l.Warn("%v", v...)
}

func (log *beeLog) Warnf(format string, v ...interface{}) {
	log.l.Warn(format, v...)
}

func (log *beeLog) Error(v ...interface{}) {
	log.l.Error("%v", v...)
}

func (log *beeLog) Errorf(format string, v ...interface{}) {
	log.l.Error(format, v...)
}
