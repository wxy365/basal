package text

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func Render(format string, args ...any) string {
	bd := strings.Builder{}
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
					bd.WriteString(format[origin:lidx])
					bd.WriteString(fmt.Sprintf("%+v", args[num]))
					origin = i + 1
					lidx = -1
				}
			}
		}
	}
	if origin == 0 {
		return format
	} else {
		bd.WriteString(format[origin:])
		return bd.String()
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
