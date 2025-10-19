package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormZapLogger 实现了 gorm.io/gorm/logger.Interface 接口
type GormZapLogger struct {
	ZapLogger *zap.Logger
	LogLevel  logger.LogLevel
}

// NewGormZapLogger 创建一个新的 GORM 日志记录器实例
func NewGormZapLogger(zapLogger *zap.Logger) *GormZapLogger {
	return &GormZapLogger{
		ZapLogger: zapLogger,
		LogLevel:  logger.Info, // 默认 GORM 日志级别
	}
}

func (l *GormZapLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *GormZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.ZapLogger.Info(fmt.Sprintf(msg, data...))
	}
}

func (l *GormZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.ZapLogger.Warn(fmt.Sprintf(msg, data...))
	}
}

func (l *GormZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.ZapLogger.Error(fmt.Sprintf(msg, data...))
	}
}

func (l *GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []zap.Field{
		zap.Duration("elapsed", elapsed),
		zap.String("sql", sql),
		zap.Int64("rows", rows),
	}

	switch {
	case err != nil && l.LogLevel >= logger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		l.ZapLogger.Error("gorm_trace", append(fields, zap.Error(err))...)
	case elapsed > 200*time.Millisecond && l.LogLevel >= logger.Warn: // 慢查询阈值
		l.ZapLogger.Warn("gorm_trace_slow_query", append(fields, zap.Duration("threshold", 200*time.Millisecond))...)
	case l.LogLevel >= logger.Info:
		l.ZapLogger.Debug("gorm_trace", fields...)
	}
}
