package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

func InitLogger() error {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	var err error
	globalLogger, err = config.Build(zap.AddCaller())
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(globalLogger)
	return nil
}

func Get() *zap.Logger {
	if globalLogger == nil {
		panic("logger not initialized - call InitLogger first")
	}
	return globalLogger
}

func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}
