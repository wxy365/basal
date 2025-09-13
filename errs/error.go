package errs

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/rs/zerolog"
	"github.com/wxy365/basal/i18n"
	"github.com/wxy365/basal/text"
)

type Err struct {
	Code    string `json:"code" form:"code"`
	Message string `json:"message" form:"message"`
	Cause   error  `json:"cause" form:"cause"`
	Status  int    `json:"status" form:"status"` // Equivalent to http status
}

type ErrMarshaller struct {
	Tgt *Err
}

func (e *ErrMarshaller) MarshalZerologObject(event *zerolog.Event) {
	if len(e.Tgt.Code) > 0 {
		event.Str("code", e.Tgt.Code)
	}
	if len(e.Tgt.Message) > 0 {
		event.Str("message", e.Tgt.Message)
	}
	if e.Tgt.Cause != nil {
		if err, ok := e.Tgt.Cause.(*Err); ok {
			event.Object("cause", &ErrMarshaller{err})
		} else {
			event.Str("cause", e.Tgt.Cause.Error())
		}
	}
}

func (e *Err) Error() string {
	b := text.Build("{")
	if len(e.Code) > 0 {
		b.Push("\"code\":\"", e.Code, "\",")
	}
	if len(e.Message) > 0 {
		b.Push("\"message\":\"", e.Message, "\",")
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
	msg := i18n.Message(ctx, key, args...)
	return &Err{
		Code:    key,
		Message: msg,
	}
}

func Wrap(cause error, msg string, args ...any) *Err {
	if len(args) > 0 {
		msg = text.Render(msg, args...)
	}
	return &Err{
		Message: msg,
		Cause:   cause,
	}
}

func I18nWrap(ctx context.Context, key string, cause error, args ...any) *Err {
	msg := i18n.Message(ctx, key, args...)
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
