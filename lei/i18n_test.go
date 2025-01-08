package lei

import (
	"context"
	"embed"
	"testing"
)

//go:embed assets
var fs embed.FS

func TestI18n(t *testing.T) {
	err := AddMessagesFromEmbedFS(fs)
	if err != nil {
		panic(err)
	}
	ctx := WithLocale(context.Background(), "zh-CN")
	if Message(ctx, "one") != "一" {
		t.Fail()
	}
	if Message(ctx, "two") != "二" {
		t.Fail()
	}
	if Message(ctx, "three") != "三" {
		t.Fail()
	}
	if Message(ctx, "four") != "四" {
		t.Fail()
	}
}
