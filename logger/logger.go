package logger

import "context"

type Field struct {
	Key   string
	Value interface{}
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugf(str string, args ...interface{})
	Infof(str string, args ...interface{})
	Warnf(str string, args ...interface{})
	Errorf(str string, args ...interface{})
	Fatalf(str string, args ...interface{})

	Debugw(msg string, fields ...Field)
	Infow(msg string, fields ...Field)
	Warnw(msg string, fields ...Field)
	Errorw(msg string, fields ...Field)
	Fatalw(msg string, fields ...Field)

	DebugCtx(ctx context.Context, msg string, fields ...Field)
	InfoCtx(ctx context.Context, msg string, fields ...Field)
	WarnCtx(ctx context.Context, msg string, fields ...Field)
	ErrorCtx(ctx context.Context, msg string, fields ...Field)
	FatalCtx(ctx context.Context, msg string, fields ...Field)

	With(fields ...Field) Logger

	Sync() error
}
