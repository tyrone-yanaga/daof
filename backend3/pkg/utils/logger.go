// internal/utils/logger.go
package utils

import (
    "context"
    "os"
    "runtime"
    "time"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var log *zap.Logger

func InitLogger(env string) {
    config := zap.NewProductionConfig()
    
    if env == "development" {
        config = zap.NewDevelopmentConfig()
        config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    }
    
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

    var err error
    log, err = config.Build(zap.AddCallerSkip(1))
    if err != nil {
        panic(err)
    }
}

func Logger() *zap.Logger {
    return log
}

func LogError(ctx context.Context, msg string, err error, fields ...zap.Field) {
    if pc, file, line, ok := runtime.Caller(1); ok {
        f := runtime.FuncForPC(pc)
        fields = append(fields,
            zap.String("function", f.Name()),
            zap.String("file", file),
            zap.Int("line", line),
        )
    }
    
    fields = append(fields, zap.Error(err))
    log.Error(msg, fields...)
}

func LogInfo(ctx context.Context, msg string, fields ...zap.Field) {
    log.Info(msg, fields...)
}