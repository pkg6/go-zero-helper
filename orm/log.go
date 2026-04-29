package orm

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm/logger"
	"time"
)

type Logger struct {
	LogLevel logger.LogLevel
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l

}

func (l *Logger) Info(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel < logger.Info {
		return
	}
	logx.WithContext(ctx).Infof(s, i...)
}

func (l *Logger) Warn(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel < logger.Warn {
		return
	}
	logx.WithContext(ctx).Infof(s, i...)
}

func (l *Logger) Error(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel < logger.Error {
		return
	}
	logx.WithContext(ctx).Infof(s, i...)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	logx.WithContext(ctx).WithDuration(elapsed).Infof("[%.3fms] [rows: %v] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
}
