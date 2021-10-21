package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

const (
	logFormatLogfmt = "logfmt"
	logFormatJSON   = "json"
)

func main() {
	var logger log.Logger
	logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	r := gin.Default()
	r.POST("/webhook/alert/:notifier")
	// r.POST("/webhook/notify/")
	r.Run()
}

func setupLogger(lvl string, fmt string) (logger log.Logger) {
	var filter level.Option
	switch lvl {
	case "error":
		filter = level.AllowError()
	case "warn":
		filter = level.AllowWarn()
	case "debug":
		filter = level.AllowDebug()
	default:
		filter = level.AllowInfo()
	}
	if fmt == logFormatJSON {
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	} else {
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	}
	logger = level.NewFilter(logger, filter)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	return
}
