package lei

import (
	"context"
	"embed"
	"encoding/json"
	cmn "github.com/wxy365/basal/text"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
	"strings"
	"sync"
)

var (
	bd         = catalog.NewBuilder()
	printerMap = sync.Map{}
)

func AddMessagesFromEmbedFS(fs embed.FS) error {
	return readDir(fs, "i18n")
}

func readDir(fs embed.FS, dir string) error {
	entries, err := fs.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		name := entry.Name()
		path := dir + "/" + name
		if entry.IsDir() {
			err = readDir(fs, path)
			if err != nil {
				return err
			}
		} else if strings.HasSuffix(name, ".json") {
			eles := strings.Split(strings.TrimSuffix(name, ".json"), "_")
			var lcl string
			l := len(eles)
			switch len(eles) {
			case 0:
				return New("illegal i18n json file name: {0}", path)
			case 1:
				lcl = eles[0]
			default:
				lcl = eles[l-1]
			}
			content, err := fs.ReadFile(path)
			if err != nil {
				return err
			}
			messages := make(map[string]string)
			err = json.Unmarshal(content, &messages)
			if err != nil {
				return New("illegal content of file {0}", path)
			}
			err = AddMessages(lcl, messages)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

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
	locales, _ := ctx.Value(localeKey).(string)
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
