package log

import "go.uber.org/zap"

type Log struct {
	*zap.SugaredLogger
}

type Logger interface {
	Infof(pattern string, msg ...interface{})
	Errorf(pattern string, msg ...interface{})
	Fatalf(pattern string, msg ...interface{})
}

func NewLogger() (Logger, error) {
	cfg := initLoggerConfig()
	logger, _ := cfg.Build()
	defer logger.Sync()

	return &Log{logger.Sugar()}, nil
}
