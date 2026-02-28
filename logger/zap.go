package logger

import (
	"context"
	"errors"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger
	config        *Config
}

func NewZapLogger(cfg *Config) (Logger, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	zl := &zapLogger{
		config: cfg,
	}

	encoder := zl.buildZapEncoder()
	level := zl.mapZapLogLevel()

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)

	var opts []zap.Option

	if cfg.EnableCaller {
		opts = append(opts, zap.AddCaller())
		opts = append(opts, zap.AddCallerSkip(1))
	}

	if cfg.EnableTrace {
		opts = append(opts, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	log := zap.New(core, opts...)

	zl.logger = log
	zl.sugaredLogger = log.Sugar()

	return zl, nil
}

func (z *zapLogger) buildZapEncoder() zapcore.Encoder {
	var encoderConfig zapcore.EncoderConfig
	var encoder zapcore.Encoder

	if z.config.Mode == ModeProduction {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeName = zapcore.FullNameEncoder

		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
		encoderConfig.ConsoleSeparator = "|"
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeName = zapcore.FullNameEncoder

		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	return encoder
}

func (z *zapLogger) mapZapLogLevel() zapcore.Level {
	switch z.config.Level {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func (z *zapLogger) convertFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	return zapFields
}

func (z *zapLogger) convertFieldsWithContext(ctx context.Context, fields []Field) []zap.Field {
	zapFields := z.convertFields(fields)

	if traceID := ctx.Value("trace_id"); traceID != nil {
		zapFields = append(zapFields, zap.Any("trace_id", traceID))
	}
	if requestID := ctx.Value("request_id"); requestID != nil {
		zapFields = append(zapFields, zap.Any("request_id", requestID))
	}
	if userID := ctx.Value("span_id"); userID != nil {
		zapFields = append(zapFields, zap.Any("span_id", userID))
	}

	return zapFields
}

func (z *zapLogger) Sync() error {
	if z.sugaredLogger != nil {
		return z.sugaredLogger.Sync()
	}

	return nil
}

func String(key, val string) Field {
	return Field{Key: key, Value: val}
}

func Int(key string, val int) Field {
	return Field{Key: key, Value: val}
}

func Int64(key string, val int64) Field {
	return Field{Key: key, Value: val}
}

func Float64(key string, val float64) Field {
	return Field{Key: key, Value: val}
}

func Bool(key string, val bool) Field {
	return Field{Key: key, Value: val}
}

func Any(key string, val interface{}) Field {
	return Field{Key: key, Value: val}
}

func Error(err error) Field {
	return Field{Key: "error", Value: err}
}

func (z *zapLogger) Debug(args ...interface{}) {
	z.sugaredLogger.Debug(args...)
}

func (z *zapLogger) Info(args ...interface{}) {
	z.sugaredLogger.Info(args...)
}

func (z *zapLogger) Warn(args ...interface{}) {
	z.sugaredLogger.Warn(args...)
}

func (z *zapLogger) Error(args ...interface{}) {
	z.sugaredLogger.Error(args...)
}

func (z *zapLogger) Fatal(args ...interface{}) {
	z.sugaredLogger.Fatal(args...)
}

func (z *zapLogger) Debugf(str string, args ...interface{}) {
	z.sugaredLogger.Debugf(str, args...)
}

func (z *zapLogger) Infof(str string, args ...interface{}) {
	z.sugaredLogger.Infof(str, args...)
}

func (z *zapLogger) Warnf(str string, args ...interface{}) {
	z.sugaredLogger.Warnf(str, args...)
}

func (z *zapLogger) Errorf(str string, args ...interface{}) {
	z.sugaredLogger.Errorf(str, args...)
}

func (z *zapLogger) Fatalf(str string, args ...interface{}) {
	z.sugaredLogger.Fatalf(str, args...)
}

func (z *zapLogger) Debugw(msg string, fields ...Field) {
	z.logger.Debug(msg, z.convertFields(fields)...)
}

func (z *zapLogger) Infow(msg string, fields ...Field) {
	z.logger.Info(msg, z.convertFields(fields)...)
}

func (z *zapLogger) Warnw(msg string, fields ...Field) {
	z.logger.Warn(msg, z.convertFields(fields)...)
}

func (z *zapLogger) Errorw(msg string, fields ...Field) {
	z.logger.Error(msg, z.convertFields(fields)...)
}

func (z *zapLogger) Fatalw(msg string, fields ...Field) {
	z.logger.Fatal(msg, z.convertFields(fields)...)
}

func (z *zapLogger) DebugCtx(ctx context.Context, msg string, fields ...Field) {
	z.logger.Debug(msg, z.convertFieldsWithContext(ctx, fields)...)
}

func (z *zapLogger) InfoCtx(ctx context.Context, msg string, fields ...Field) {
	z.logger.Info(msg, z.convertFieldsWithContext(ctx, fields)...)
}

func (z *zapLogger) WarnCtx(ctx context.Context, msg string, fields ...Field) {
	z.logger.Warn(msg, z.convertFieldsWithContext(ctx, fields)...)
}

func (z *zapLogger) ErrorCtx(ctx context.Context, msg string, fields ...Field) {
	z.logger.Error(msg, z.convertFieldsWithContext(ctx, fields)...)
}

func (z *zapLogger) FatalCtx(ctx context.Context, msg string, fields ...Field) {
	z.logger.Fatal(msg, z.convertFieldsWithContext(ctx, fields)...)
}

func (z *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{
		logger: z.logger.With(z.convertFields(fields)...),
	}
}
