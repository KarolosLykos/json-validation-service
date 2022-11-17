package logruslog

import (
	"context"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/KarolosLykos/json-validation-service/internal/config"
	"github.com/KarolosLykos/json-validation-service/internal/logger"
)

type logruslog struct {
	log *logrus.Logger
}

type key string

// Settings holds the settings for the logger.
var Settings = struct {
	ErrorKey key
	CtxKey   key
}{
	ErrorKey: "error",
	CtxKey:   "ctxKey",
}

var defaultLogger = &logrus.Logger{
	Out:          os.Stderr,
	Hooks:        make(logrus.LevelHooks),
	Level:        logrus.DebugLevel,
	ExitFunc:     os.Exit,
	ReportCaller: false,
	Formatter:    &logrus.JSONFormatter{},
}

func New(l *logrus.Logger) logger.Logger {
	return &logruslog{
		log: l,
	}
}

func DefaultLogger(cfg *config.Config) logger.Logger {
	log := defaultLogger
	log.SetFormatter(setFormatter(cfg.Logger.Format))
	log.SetLevel(setLevel(cfg.Logger.Level))

	return New(defaultLogger)
}

func setFormatter(format string) logrus.Formatter {
	switch strings.ToLower(format) {
	case "json":
		return &logrus.JSONFormatter{}
	default:
		return &logrus.TextFormatter{}
	}
}

func setLevel(lvl string) logrus.Level {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return logrus.InfoLevel
	}

	return level
}

func (l *logruslog) Panic(ctx context.Context, err error, messages ...interface{}) {
	le := l.parseMessages(ctx, err)
	le.Panic(messages...)
}

func (l *logruslog) Error(ctx context.Context, err error, messages ...interface{}) {
	le := l.parseMessages(ctx, err)
	le.Error(messages...)
}

func (l *logruslog) Info(ctx context.Context, messages ...interface{}) {
	le := l.parseMessages(ctx, nil)
	le.Info(messages...)
}

func (l *logruslog) Debug(ctx context.Context, messages ...interface{}) {
	le := l.parseMessages(ctx, nil)
	le.Debug(messages...)
}

func (l *logruslog) parseMessages(ctx context.Context, err error) *logrus.Entry {
	if ctx == nil {
		ctx = context.TODO()
	}

	e := l.log.WithFields(logrus.Fields{string(Settings.CtxKey): ctx.Value(Settings.CtxKey)})

	if err != nil {
		e = e.WithField(string(Settings.ErrorKey), err.Error())
	}

	return e
}
