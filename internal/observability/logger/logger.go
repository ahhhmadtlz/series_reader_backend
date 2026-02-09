package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Level string

const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
)

type Config struct {
	Level            Level  `koanf:"level"`
	FilePath         string `koanf:"file_path"`
	MaxSizeInMB      int    `koanf:"file_max_size_mb"`
	MaxAgeInDays     int    `koanf:"file_max_age_days"`
	MaxBackups       int    `koanf:"max_backups"`
	Compress         bool   `koanf:"compress"`
}

func Default() Config {
	return Config{
		Level:        InfoLevel,
		FilePath:     "logs/app.log",
		MaxSizeInMB:  100,
		MaxAgeInDays: 30,
		MaxBackups:   5,
		Compress:     true,
	}
}

// Global logger instance
var Log *zap.Logger

// New creates a new logger and sets it as the global logger
func New(cfg Config) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if err := level.Set(string(cfg.Level)); err != nil {
		return nil, err
	}

	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSizeInMB,
		MaxAge:     cfg.MaxAgeInDays,
		MaxBackups: cfg.MaxBackups,
		Compress:   cfg.Compress,
	})

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.NewMultiWriteSyncer(writer, zapcore.AddSync(os.Stdout)),
		level,
	)

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// Set global logger
	Log = zapLogger

	return zapLogger, nil
}

// Global helper functions (always use these)

func Info(msg string, fields ...interface{}) {
	if Log != nil {
		Log.Sugar().Infow(msg, fields...)
	}
}

func Error(msg string, fields ...interface{}) {
	if Log != nil {
		Log.Sugar().Errorw(msg, fields...)
	}
}

func Debug(msg string, fields ...interface{}) {
	if Log != nil {
		Log.Sugar().Debugw(msg, fields...)
	}
}

func Warn(msg string, fields ...interface{}) {
	if Log != nil {
		Log.Sugar().Warnw(msg, fields...)
	}
}

func Fatal(msg string, fields ...interface{}) {
	if Log != nil {
		Log.Sugar().Fatalw(msg, fields...)
	}
}