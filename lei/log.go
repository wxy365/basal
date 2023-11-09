package lei

import (
	"github.com/wxy365/basal/text"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	log           *zap.Logger
	logTraceLevel zapcore.Level
)

func init() {
	encoding, ok := os.LookupEnv("SPROUT_LOG_ENCODING")
	if !ok {
		encoding = "json"
	}

	parseLogLevel := func(rawLevel string) zapcore.Level {
		switch rawLevel {
		case "debug":
			return zapcore.DebugLevel
		case "info":
			return zapcore.InfoLevel
		case "warn":
			return zapcore.WarnLevel
		case "error":
			return zapcore.ErrorLevel
		case "panic":
			return zapcore.PanicLevel
		case "fatal":
			return zapcore.FatalLevel
		default:
			return zapcore.ErrorLevel
		}
	}

	level, ok := os.LookupEnv("SPROUT_LOG_LEVEL")
	if !ok {
		level = "warn"
	}
	zapLevel := parseLogLevel(level)

	traceLevel, ok := os.LookupEnv("SPROUT_LOG_STACK_TRACE_LEVEL")
	if !ok {
		traceLevel = "error"
	}
	logTraceLevel = parseLogLevel(traceLevel)

	logPath, ok := os.LookupEnv("SPROUT_LOG_PATH")
	if !ok {
		logPath = "stdout"
	}

	cfg := zap.Config{
		Encoding:         encoding,
		Level:            zap.NewAtomicLevelAt(zapLevel),
		OutputPaths:      []string{logPath},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,
		},
	}
	log, _ = cfg.Build()
}

func Debug(format string, args ...any) {
	msg := text.Render(format, args...)
	if logTraceLevel <= zapcore.DebugLevel {
		log.Debug(msg, zap.StackSkip("stack", 1))
	} else {
		log.Debug(msg)
	}
}

func DebugErr(err error) {
	if e, ok := err.(*Err); ok {
		DebugErrF(e.Message, e.Cause)
	} else {
		log.Debug(err.Error())
	}
}

func DebugErrF(format string, err error, args ...any) {
	if err == nil {
		Debug(format, args...)
		return
	}
	msg := text.Render(format, args...)
	fields := []zap.Field{getCauseField(err)}
	if logTraceLevel <= zapcore.DebugLevel {
		fields = append(fields, zap.StackSkip("stack", 1))
	}
	log.Debug(msg, fields...)
}

func Info(format string, args ...any) {
	msg := text.Render(format, args...)
	if logTraceLevel <= zapcore.InfoLevel {
		log.Info(msg, zap.StackSkip("stack", 1))
	} else {
		log.Info(msg)
	}
}

func InfoErr(err error) {
	if e, ok := err.(*Err); ok {
		InfoErrF(e.Message, e.Cause)
	} else {
		log.Info(err.Error())
	}
}

func InfoErrF(format string, err error, args ...any) {
	if err == nil {
		Info(format, args...)
		return
	}
	msg := text.Render(format, args...)
	fields := []zap.Field{getCauseField(err)}
	if logTraceLevel <= zapcore.InfoLevel {
		fields = append(fields, zap.StackSkip("stack", 1))
	}
	log.Info(msg, fields...)
}

func Warn(format string, args ...any) {
	msg := text.Render(format, args...)
	if logTraceLevel <= zapcore.WarnLevel {
		log.Warn(msg, zap.StackSkip("stack", 1))
	} else {
		log.Warn(msg)
	}
}

func WarnErr(err error) {
	if e, ok := err.(*Err); ok {
		WarnErrF(e.Message, e.Cause)
	} else {
		log.Warn(err.Error())
	}
}

func WarnErrF(format string, err error, args ...any) {
	if err == nil {
		Warn(format, args...)
		return
	}
	msg := text.Render(format, args...)
	fields := []zap.Field{getCauseField(err)}
	if logTraceLevel <= zapcore.WarnLevel {
		fields = append(fields, zap.StackSkip("stack", 1))
	}
	log.Warn(msg, fields...)
}

func Error(format string, args ...any) {
	msg := text.Render(format, args...)
	if logTraceLevel <= zapcore.ErrorLevel {
		log.Error(msg, zap.StackSkip("stack", 1))
	} else {
		log.Error(msg)
	}
}

func ErrorErr(err error) {
	if e, ok := err.(*Err); ok {
		ErrorErrF(e.Message, e.Cause)
	} else {
		log.Error(err.Error())
	}
}

func ErrorErrF(format string, err error, args ...any) {
	if err == nil {
		Error(format, args...)
		return
	}
	msg := text.Render(format, args...)
	fields := []zap.Field{getCauseField(err)}
	if logTraceLevel <= zapcore.ErrorLevel {
		fields = append(fields, zap.StackSkip("stack", 1))
	}
	log.Error(msg, fields...)
}

func Panic(format string, args ...any) {
	msg := text.Render(format, args...)
	if logTraceLevel <= zapcore.PanicLevel {
		log.Panic(msg, zap.StackSkip("stack", 1))
	} else {
		log.Panic(msg)
	}
}

func PanicErr(err error) {
	if e, ok := err.(*Err); ok {
		PanicErrF(e.Message, e.Cause)
	} else {
		log.Panic(err.Error())
	}
}

func PanicErrF(format string, err error, args ...any) {
	if err == nil {
		Panic(format, args...)
		return
	}
	msg := text.Render(format, args...)
	fields := []zap.Field{getCauseField(err)}
	if logTraceLevel <= zapcore.PanicLevel {
		fields = append(fields, zap.StackSkip("stack", 1))
	}
	log.Panic(msg, fields...)
}

func getCauseField(err error) zap.Field {
	if e, ok := err.(*Err); ok {
		return zap.Object("cause", &ErrMarshaller{e})
	} else {
		return zap.String("cause", err.Error())
	}
}
