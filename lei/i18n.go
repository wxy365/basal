package lei

import (
	"context"
	cmn "github.com/wxy365/basal/text"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
	"sync"
)

var (
	bd         = catalog.NewBuilder()
	printerMap = sync.Map{}
)

func AddMessages(locale string, messages map[string]string) error {
	langTag, err := language.Parse(locale)
	if err != nil {
		return Wrap("Failed to parse language {0}", err, locale)
	}
	for k, v := range messages {
		err := bd.SetString(langTag, k, v)
		if err != nil {
			return Wrap("Failed to set message {{0}:{1}}", err, k, v)
		}
	}
	printer := message.NewPrinter(langTag, message.Catalog(bd))
	printerMap.Store(langTag, printer)
	return nil
}

func MustAddMessages(locale string, messages map[string]string) {
	err := AddMessages(locale, messages)
	if err != nil {
		panic(err)
	}
}

type locale int

var localeKey locale

func Message(ctx context.Context, key string, args ...any) string {
	locales := ctx.Value(localeKey).(string)
	languages := bd.Languages()
	matcher := language.NewMatcher(languages)
	tag, _ := language.MatchStrings(matcher, locales)
	printer, ok := printerMap.Load(tag)
	if ok {
		msg := printer.(*message.Printer).Sprintf(key)
		if len(args) > 0 {
			msg = cmn.Render(msg, args...)
		}
		return msg
	}
	return key
}

func WithLocale(ctx context.Context, locales string) context.Context {
	return context.WithValue(ctx, localeKey, locales)
}
