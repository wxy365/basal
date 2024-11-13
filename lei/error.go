package lei

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"github.com/wxy365/basal/text"
	"reflect"
	"strings"
)

type Err struct {
	Code    string
	Message string
	Cause   error
	Status  int // Equivalent to http status
}

type ErrMarshaller struct {
	tgt *Err
}

func (e *ErrMarshaller) MarshalZerologObject(event *zerolog.Event) {
	if len(e.tgt.Code) > 0 {
		event.Str("code", e.tgt.Code)
	}
	if len(e.tgt.Message) > 0 {
		event.Str("message", e.tgt.Message)
	}
	if e.tgt.Cause != nil {
		if err, ok := e.tgt.Cause.(*Err); ok {
			event.Object("cause", &ErrMarshaller{err})
		} else {
			event.Str("cause", e.tgt.Cause.Error())
		}
	}
}

func (e *Err) Error() string {
	b := text.Build("{")
	if len(e.Code) > 0 {
		b.Push("\"code\":\"", e.Code, "\",")
	}
	if len(e.Message) > 0 {
		b.Push("\"msg\":\"", e.Message, "\",")
	}
	if e.Cause != nil {
		var er *Err
		if errors.As(e.Cause, &er) {
			b.Push("\"cause\":", er.Error())
		} else {
			b.Push("\"cause\":\"", e.Cause.Error(), "\"")
		}
	}
	return strings.Trim(b.String(), ",") + "}"
}

func (e *Err) Unwrap() error {
	return e.Cause
}

func (e *Err) Is(target error) bool {
	return errors.Is(e, target)
}

func IsOk(err error) bool {
	return reflect.ValueOf(err).IsZero()
}

func (e *Err) WithCode(code string) *Err {
	e.Code = code
	return e
}

func (e *Err) WithStatus(status int) *Err {
	e.Status = status
	return e
}

func New(msg string, args ...any) *Err {
	if len(args) > 0 {
		msg = text.Render(msg, args...)
	}
	return &Err{
		Message: msg,
	}
}

func I18nNew(ctx context.Context, key string, args ...any) *Err {
	msg := Message(ctx, key, args...)
	return &Err{
		Code:    key,
		Message: msg,
	}
}

func Wrap(msg string, cause error, args ...any) *Err {
	if len(args) > 0 {
		msg = text.Render(msg, args...)
	}
	return &Err{
		Message: msg,
		Cause:   cause,
	}
}

func I18nWrap(ctx context.Context, key string, cause error, args ...any) *Err {
	msg := Message(ctx, key, args...)
	return &Err{
		Code:    key,
		Message: msg,
		Cause:   cause,
	}
}

func Adapt(err error) *Err {
	if e, ok := err.(*Err); ok {
		return e
	} else {
		return New(err.Error())
	}
}
