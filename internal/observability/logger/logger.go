package logger



import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultFilePath         = "logs/app.lgo"
	defaultUseLocalTime     = false
	defaultFileMaxSizeInMB  = 100
	defaultFileMaxAgeInDays = 30
	defaultMaxBackups       = 5
	defaultCompress         = true
)

type Config struct {
	FilePath         string        `koanf:"file_path"`
	UseLocalTime     bool          `koanf:"use_local_time"`
	FileMaxSizeInMB  int           `koanf:"file_max_size_mb"`
	FileMaxAgeInDays int           `koanf:"file_max_age_days"`
	MaxBackups       int           `koanf:"max_backups"`
	Compress         bool          `koanf:"compress"`
}

var l *slog.Logger

// init is default logger and Singleton that lets you ensure that a logger has only one instance,
// while providing a global access point to this instance.
func init(){
	fileWriter:=&lumberjack.Logger{
		Filename: defaultFilePath,
		LocalTime: defaultUseLocalTime,
		MaxSize: defaultFileMaxSizeInMB,
		MaxAge: defaultFileMaxAgeInDays,
		MaxBackups: defaultMaxBackups,
		Compress: defaultCompress,
	}

	l=slog.New(
		slog.NewJSONHandler(io.MultiWriter(fileWriter,os.Stdout),&slog.HandlerOptions{}),
	)
}

// L returns the singleton logger instance
func L()*slog.Logger{
	return l
}

func SetDefault(logger *slog.Logger) {
	l = logger
}


// New is constructor logger with special settings.
func New(cfg Config,opt *slog.HandlerOptions)*slog.Logger {
	fileWriter:=&lumberjack.Logger{
		Filename: cfg.FilePath,
		LocalTime: cfg.UseLocalTime,
		MaxSize: cfg.FileMaxSizeInMB,
		MaxAge: cfg.FileMaxAgeInDays,
		MaxBackups: cfg.MaxBackups,
		Compress: cfg.Compress,
	}

	logger:=slog.New(
		slog.NewJSONHandler(io.MultiWriter(fileWriter,os.Stdout),opt),
	)
	return logger
}


func Debug(msg string, args ...any) {
	l.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	l.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	l.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	l.Error(msg, args...)
}
