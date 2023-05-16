package logger

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzlogrus "github.com/hertz-contrib/logger/logrus"
	"github.com/sirupsen/logrus"
	"github.com/three-body/hertz-scaffold/config"
)

func Init() {
	logger := hertzlogrus.NewLogger(hertzlogrus.WithHook(&RequestIdHook{}))
	hlog.SetLogger(logger)
	hlog.SetLevel(getLevel(config.GetConf().Logger.Level))
}

type RequestIdHook struct{}

func (h *RequestIdHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *RequestIdHook) Fire(e *logrus.Entry) error {
	ctx := e.Context
	if ctx == nil {
		return nil
	}
	value := ctx.Value("X-Request-ID")
	if value != nil {
		e.Data["log_id"] = value
	}
	return nil
}

func getLevel(level string) hlog.Level {
	switch level {
	case "trace":
		return hlog.LevelTrace
	case "debug":
		return hlog.LevelDebug
	case "info":
		return hlog.LevelInfo
	case "warn":
		return hlog.LevelWarn
	case "error":
		return hlog.LevelError
	case "fatal":
		return hlog.LevelFatal
	default:
		return hlog.LevelTrace
	}
}
