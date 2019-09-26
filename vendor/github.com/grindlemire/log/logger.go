package log

import (
	"fmt"
	stdlogger "log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Opts are options you can configure the logger with
type Opts struct {
	Level              Level `json:"log_level"            env:"log_level"             default:"INFO"  description:"log level to log at (possible values are debug, info, warn, error, fatal, panic)"`
	MaxLogSize         int   `json:"log_max_size"         env:"log_max_size"          default:"10"    description:"Max size of a log before rolling over"`
	MaxLogBackups      int   `json:"log_max_backups"      env:"log_max_backups"       default:"5"     description:"Max number of backups to keep"`
	CompressBackupLogs bool  `json:"log_compress_backups" env:"log_compress_backups"  default:"false" description:"Whether to compress backups or not"`
	Console            bool  `json:"log_console"          env:"log_console"           default:"true"  description:"Whether to log to the console or not (through stdout)"`
	CallerSkip         int   `json:"log_caller_skip"      env:"log_caller_skip"       default:"1"     description:"How many levels of stack to skip before logging in your application (defaults to 1 for this library)"`
}

// Default is the default config for on the fly use
var Default = Opts{
	Level:              InfoLevel,
	CallerSkip:         1,
	MaxLogSize:         10,
	MaxLogBackups:      5,
	CompressBackupLogs: false,
	Console:            true,
}

// Fields are key values that we will decorate the message with
type Fields map[string]interface{}

// Level is a simple wrapper for log levels so you just import this library
type Level string

// Wrap the log levels so you only have to import this library
const (
	DebugLevel = "DEBUG"
	InfoLevel  = "INFO"
	WarnLevel  = "WARN"
	ErrorLevel = "ERROR"
	FatalLevel = "FATAL"
	PanicLevel = "PANIC"
)

// loaded is a simple internal flag to specify whether we have been initialized or
// if we should fall back on the default go logger
var loaded = false

// rfc3339Encoder encodes the time field to rfc3339 with millisecond precision
func rfc3339Encoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000Z"))
}

// createEncoderConfig encapsulates all the custom configuration changes we are making to the encoder
func createEncoderConfig() (encoderConfig zapcore.EncoderConfig) {
	encoderConfig = zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = rfc3339Encoder
	encoderConfig.CallerKey = "caller"
	encoderConfig.TimeKey = "ts"
	encoderConfig.MessageKey = "msg"
	encoderConfig.LevelKey = "level"
	return encoderConfig
}

// Init initializes the logger a level and any number of log files (including none)
func Init(opts Opts, logPaths ...string) (err error) {
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		l := zapcore.InfoLevel
		l.UnmarshalText([]byte(opts.Level))
		return lvl >= l
	})
	encoderConfig := createEncoderConfig()
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	cores := []zapcore.Core{}

	if opts.Console {
		stdout := zapcore.Lock(os.Stdout)
		cores = append(cores, zapcore.NewCore(encoder, stdout, infoPriority))
	}

	if opts.CallerSkip <= 0 {
		return fmt.Errorf("caller skip must be > 0 but was %d", opts.CallerSkip)
	}

	if len(logPaths) > 0 {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)

		for _, logPath := range logPaths {
			f, err := os.OpenFile(logPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
			if err != nil {
				return err
			}
			f.Close()

			logFile := zapcore.AddSync(&lumberjack.Logger{
				Filename:   logPath,
				MaxSize:    opts.MaxLogSize,
				MaxBackups: opts.MaxLogBackups,
				Compress:   opts.CompressBackupLogs,
			})
			cores = append(cores, zapcore.NewCore(encoder, logFile, infoPriority))
		}
	}

	core := zapcore.NewTee(cores...)
	zap.ReplaceGlobals(zap.New(core, zap.AddCaller(), zap.AddCallerSkip(opts.CallerSkip)))
	loaded = true
	return nil
}

// Info logs info statements
func Info(args ...interface{}) {
	if !loaded {
		stdlogger.Print(args...)
		return
	}
	zap.S().Info(args...)
}

// Infof logs infof statements
func Infof(template string, args ...interface{}) {
	if !loaded {
		stdlogger.Printf(template, args...)
		return
	}

	zap.S().Infof(template, args...)
}

// Infow logs infow statements
func Infow(msg string, fields Fields) {
	if !loaded {
		stdlogger.Print(msg)
		return
	}
	zap.S().Infow(msg, convertToZapFields(fields)...)
}

// Warn logs Warn statements
func Warn(args ...interface{}) {
	if !loaded {
		stdlogger.Print(args...)
		return
	}
	zap.S().Warn(args...)
}

// Warnf logs Warnf statements
func Warnf(template string, args ...interface{}) {
	if !loaded {
		stdlogger.Printf(template, args...)
		return
	}

	zap.S().Warnf(template, args...)
}

// Warnw logs Warnw statements
func Warnw(msg string, fields Fields) {
	if !loaded {
		stdlogger.Print(msg)
		return
	}
	zap.S().Warnw(msg, convertToZapFields(fields)...)
}

// Error logs Error statements
func Error(args ...interface{}) {
	if !loaded {
		stdlogger.Print(args...)
		return
	}
	zap.S().Error(args...)
}

// Errorf logs Errorf statements
func Errorf(template string, args ...interface{}) {
	if !loaded {
		stdlogger.Printf(template, args...)
		return
	}

	zap.S().Errorf(template, args...)
}

// Errorw logs Errorw statements
func Errorw(msg string, fields Fields) {
	if !loaded {
		stdlogger.Print(msg)
		return
	}
	zap.S().Errorw(msg, convertToZapFields(fields)...)
}

// Debug logs Debug statements
func Debug(args ...interface{}) {
	if !loaded {
		stdlogger.Print(args...)
		return
	}
	zap.S().Debug(args...)
}

// Debugf logs Debugf statements
func Debugf(template string, args ...interface{}) {
	if !loaded {
		stdlogger.Printf(template, args...)
		return
	}

	zap.S().Debugf(template, args...)
}

// Debugw logs Debugw statements
func Debugw(msg string, fields Fields) {
	if !loaded {
		stdlogger.Print(msg)
		return
	}
	zap.S().Debugw(msg, convertToZapFields(fields)...)
}

// Fatal logs Fatal statements
func Fatal(args ...interface{}) {
	if !loaded {
		stdlogger.Print(args...)
		os.Exit(1)
		return
	}
	zap.S().Fatal(args...)
}

// Fatalf logs Fatalf statements
func Fatalf(template string, args ...interface{}) {
	if !loaded {
		stdlogger.Printf(template, args...)
		os.Exit(1)
	}

	zap.S().Fatalf(template, args...)
}

// Fatalw logs Fatalw statements
func Fatalw(msg string, fields Fields) {
	if !loaded {
		stdlogger.Print(msg)
		os.Exit(1)
	}
	zap.S().Fatalw(msg, convertToZapFields(fields)...)
}

// Panic logs Panic statements
func Panic(args ...interface{}) {
	if !loaded {
		panic(fmt.Sprint(args...))
	}
	zap.S().Panic(args...)
}

// Panicf logs Panicf statements
func Panicf(template string, args ...interface{}) {
	if !loaded {
		panic(fmt.Sprintf(template, args...))
	}

	zap.S().Panicf(template, args...)
}

// Panicw logs Panicw statements
func Panicw(msg string, fields Fields) {
	if !loaded {
		panic(fmt.Sprint(msg))
	}
	zap.S().Panicw(msg, convertToZapFields(fields)...)
}

func convertToZapFields(f Fields) (output []interface{}) {
	for k, v := range f {
		output = append(output, k, v)
	}
	return output
}
