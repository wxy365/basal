package lei

import (
	"errors"
	"github.com/rs/zerolog"
	"github.com/wxy365/basal/text"
	"os"
)

var (
	logger zerolog.Logger
)

func init() {
	var level zerolog.Level
	var err error
	levelValue, ok := os.LookupEnv("BASAL_LOG_LEVEL")
	if !ok {
		level = zerolog.WarnLevel
	} else {
		level, err = zerolog.ParseLevel(levelValue)
		if err != nil {
			panic(Wrap("Cannot parse log level: "+levelValue, err))
		}
	}

	out := os.Stdout
	logPath, ok := os.LookupEnv("BASAL_LOG_PATH")
	if ok {
		out, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(Wrap("Cannot open log file", err))
		}
	}

	logger = zerolog.New(
		out,
	).Level(level).With().Timestamp().Caller().Logger()
}

func Debug(format string, args ...any) {
	log(logger.Debug(), format, args...)
}

func DebugErr(err error) {
	logErr(logger.Debug(), err)
}

func DebugErrF(format string, err error, args ...any) {
	logErrF(logger.Debug(), format, err, args...)
}

func Info(format string, args ...any) {
	log(logger.Info(), format, args...)
}

func InfoErr(err error) {
	logErr(logger.Info(), err)
}

func InfoErrF(format string, err error, args ...any) {
	logErrF(logger.Info(), format, err, args...)
}

func Warn(format string, args ...any) {
	log(logger.Warn(), format, args...)
}

func WarnErr(err error) {
	logErr(logger.Warn(), err)
}

func WarnErrF(format string, err error, args ...any) {
	logErrF(logger.Warn(), format, err, args...)
}

func Error(format string, args ...any) {
	log(logger.Error(), format, args...)
}

func ErrorErr(err error) {
	logErr(logger.Error(), err)
}

func ErrorErrF(format string, err error, args ...any) {
	logErrF(logger.Error(), format, err, args...)
}

func Fatal(format string, args ...any) {
	log(logger.Fatal(), format, args...)
}

func FatalErr(err error) {
	logErr(logger.Fatal(), err)
}

func FatalErrF(format string, err error, args ...any) {
	logErrF(logger.Fatal(), format, err, args...)
}

func Panic(format string, args ...any) {
	log(logger.Panic(), format, args...)
}

func PanicErr(err error) {
	logErr(logger.Panic(), err)
}

func PanicErrF(format string, err error, args ...any) {
	logErrF(logger.Panic(), format, err, args...)
}

func log(event *zerolog.Event, format string, args ...any) {
	msg := text.Render(format, args...)
	event.CallerSkipFrame(2).Msg(msg)
}

func logErr(event *zerolog.Event, err error) {
	event.CallerSkipFrame(2)
	var e *Err
	if errors.As(err, &e) {
		if e.Cause != nil {
			appendCause(event, e.Cause)
		}
		event.Msg(e.Message)
		return
	}
	event.Msg(err.Error())
}

func logErrF(event *zerolog.Event, format string, err error, args ...any) {
	if err != nil {
		appendCause(event, err)
	}
	event.CallerSkipFrame(2).Msg(text.Render(format, args...))
}

func appendCause(event *zerolog.Event, err error) *zerolog.Event {
	var e *Err
	if errors.As(err, &e) {
		event.Object("cause", &ErrMarshaller{e})
	} else {
		event.Str("cause", err.Error())
	}
	return event
}
