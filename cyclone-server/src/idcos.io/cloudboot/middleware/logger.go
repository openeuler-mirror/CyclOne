package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"idcos.io/cloudboot/logger"
)

// ctxLoggerKey 注入的logger.Logger对应的查询Key
var ctxLoggerKey uint8
// panicFileKey 注入的panic file对应的查询Key
var panicFileKey uint8

// LoggerFromContext 从ctx中获取model.Repo
func LoggerFromContext(ctx context.Context) (logger.Logger, bool) {
	log, ok := ctx.Value(&ctxLoggerKey).(logger.Logger)
	return log, ok
}

// PanicFileFromContext 从ctx中获取model.Repo
func PanicFileFromContext(ctx context.Context) (*os.File, bool) {
	panicFile, ok := ctx.Value(&panicFileKey).(*os.File)
	return panicFile, ok
}

// InjectLogger 注入logger.Logger
func InjectLogger(logger logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), &ctxLoggerKey, logger))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// InjectFile 注入Panic File
func InjectFile(panic *os.File) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), &panicFileKey, panic))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// LogPanic http panic recover
func LogPanic(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		file, _ := PanicFileFromContext(r.Context())

		if _, err := os.Stat(file.Name()); err != nil {
			if os.IsNotExist(err) {
				file, _ = os.OpenFile(file.Name(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			}
		}

		defer func() {
			if rvr := recover(); rvr != nil {
				// append panic info to file
				panicTime := time.Now().Format("2006-01-02 15:04:05")
				if _, err := file.WriteString(fmt.Sprintf("\n%s: %v\n %s", panicTime, rvr, debug.Stack())); err != nil {
					fmt.Printf("write panic file error , %s", err.Error())
				}

				//若出现异常，则进行返回500
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				respMap := map[string]interface{}{
					"status":  "failure",
					"message": fmt.Sprintf("unknown error, error message in %s", file.Name()),
				}
				json.NewEncoder(w).Encode(respMap)
			}
		}()
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
