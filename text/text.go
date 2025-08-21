package text

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/wxy365/basal/fn"
)

func Render(format string, args ...any) string {
	b := Build()
	lidx := -1
	origin := 0
	for i := 0; i < len(format); i++ {
		c := format[i]
		if c == '{' {
			lidx = i
		} else if c == '}' {
			if lidx > -1 {
				num, err := strconv.ParseInt(format[lidx+1:i], 10, 32)
				if err == nil && int(num) < len(args) && int(num) > -1 {
					b.Push(format[origin:lidx], fmt.Sprintf("%+v", args[num]))
					origin = i + 1
					lidx = -1
				}
			}
		}
	}
	if origin == 0 {
		return format
	} else {
		b.Push(format[origin:])
		return b.String()
	}
}

func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

func Pascal2Snake(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	l := len(r)
	o := make([]rune, 0)
	o = append(o, unicode.ToLower(r[0]))
	for i := 1; i < l-1; i++ {
		if unicode.IsUpper(r[i]) {
			if unicode.IsLower(r[i-1]) || unicode.IsLower(r[i+1]) {
				o = append(o, '_', unicode.ToLower(r[i]))
			} else {
				o = append(o, unicode.ToLower(r[i]))
			}
		} else {
			o = append(o, r[i])
		}
	}
	if unicode.IsUpper(r[l-1]) && l >= 2 && unicode.IsLower(r[l-2]) {
		o = append(o, '_', unicode.ToLower(r[l-1]))
	} else {
		o = append(o, unicode.ToLower(r[l-1]))
	}
	return string(o)
}

func Snake2Pascal(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	l := len(r)
	o := make([]rune, 0)
	if r[0] != '_' {
		o = append(o, unicode.ToUpper(r[0]))
	}
	for i := 1; i < l; i++ {
		if r[i] != '_' {
			if r[i-1] == '_' {
				o = append(o, unicode.ToUpper(r[i]))
			} else {
				o = append(o, r[i])
			}
		}
	}
	return string(o)
}

func Join[T any](items []T, delimiter string, f fn.Function[T, string], preSuffix ...string) string {
	var content string
	switch len(items) {
	case 0:
	case 1:
		content = f(items[0])
	default:
		b := Build(f(items[0]))
		for _, i := range items[1:] {
			b.Push(delimiter, f(i))
		}
		content = b.String()
	}
	if len(preSuffix) > 0 {
		content = preSuffix[0] + content
		if len(preSuffix) > 1 {
			content += preSuffix[1]
		}
	}
	return content
}

type Builder struct {
	builder strings.Builder
}

func (b *Builder) Push(s ...string) *Builder {
	for _, item := range s {
		b.builder.WriteString(item)
	}
	return b
}

func (b *Builder) PushF(format string, args ...any) *Builder {
	b.Push(Render(format, args...))
	return b
}

func (b *Builder) PushInt(i int) *Builder {
	b.builder.WriteString(strconv.FormatInt(int64(i), 10))
	return b
}

func (b *Builder) String() string {
	return b.builder.String()
}

func Build(s ...string) *Builder {
	b := new(Builder)
	b.Push(s...)
	return b
}
