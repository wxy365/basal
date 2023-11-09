package lei

import (
	"context"
	"github.com/wxy365/basal/text"
	"go.uber.org/zap/zapcore"
	"strings"
)

type Err struct {
	Code    string
	Message string
	Cause   error
	Status  int
}

type ErrMarshaller struct {
	tgt *Err
}

func (e *ErrMarshaller) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	if len(e.tgt.Code) > 0 {
		encoder.AddString("code", e.tgt.Code)
	}
	if len(e.tgt.Message) > 0 {
		encoder.AddString("message", e.tgt.Message)
	}
	if e.tgt.Cause != nil {
		if err, ok := e.tgt.Cause.(*Err); ok {
			err1 := encoder.AddObject("cause", &ErrMarshaller{err})
			if err1 != nil {
				return err1
			}
		} else {
			encoder.AddString("cause", e.tgt.Cause.Error())
		}
	}
	return nil
}

func (e *Err) Error() string {
	var b strings.Builder
	b.WriteString("{")
	if len(e.Code) > 0 {
		b.WriteString("\"code\":\"" + e.Code + "\",")
	}
	if len(e.Message) > 0 {
		b.WriteString("\"msg\":\"" + e.Message + "\",")
	}
	if e.Cause != nil {
		if err, yes := e.Cause.(*Err); yes {
			b.WriteString("\"cause\":" + string(err.Error()))
		} else {
			b.WriteString("\"cause\":\"" + e.Cause.Error() + "\"")
		}
	}
	return strings.Trim(b.String(), ",") + "}"
}

func (e *Err) Unwrap() error {
	return e.Cause
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
