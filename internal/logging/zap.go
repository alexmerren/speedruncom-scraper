package logger

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func NewDevelopmentLogger(encoding, logLevel string) (*ZapLogger, error) {
	if encoding == "" {
		return nil, errors.New("encoding not found")
	}

	if logLevel == "" {
		return nil, errors.New("log level not found")
	}

	var level zapcore.Level
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		return nil, err
	}

	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Encoding:         encoding,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     "\n",
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		},
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{
		logger: logger.Sugar(),
	}, nil
}

func NewProductionLogger(logLevel string) (*ZapLogger, error) {
	zapConfig := zap.NewProductionConfig()
	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{
		logger: logger.Sugar(),
	}, nil
}

func (z *ZapLogger) Debug(msg interface{}) {
	z.logger.Debug(msg)
}

func (z *ZapLogger) Info(msg interface{}) {
	z.logger.Info(msg)
}

func (z *ZapLogger) Error(msg interface{}) {
	z.logger.Error(msg)
}

func (z *ZapLogger) Fatal(msg interface{}) {
	z.logger.Fatal(msg)
}

func (z *ZapLogger) With(key string, value interface{}) Logger {
	return &ZapLogger{
		z.logger.With(key, value),
	}
}

func (z *ZapLogger) WithError(err error) Logger {
	return &ZapLogger{
		z.logger.With("err", err),
	}
}

func (z *ZapLogger) Sync() error {
	return z.logger.Sync()
}
