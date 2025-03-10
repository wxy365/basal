package cfg

import (
	"errors"
	"github.com/wxy365/basal/lei"
)

func cfgAbsent(key string) error {
	return lei.New("The configuration item '{0}' cannot be found", key).WithCode(ErrCodeCfgMissing)
}

const (
	ErrCodeCfgNotInitialized = "CFG_NOT_INITIALIZED"
	ErrCodeCfgMissing        = "CFG_MISSING"
	ErrCodeCfgBadType        = "CFG_BAD_TYPE"
)

var (
	ErrCfgNotInitialized = lei.New("configuration not initialized").WithCode(ErrCodeCfgNotInitialized)
)

func IsCfgMissingErr(err error) bool {
	var e *lei.Err
	if !errors.As(err, &e) {
		return false
	}
	return e.Code == ErrCodeCfgMissing
}
