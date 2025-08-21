package cfg

import (
	"errors"

	"github.com/wxy365/basal/errs"
)

func cfgAbsent(key string) error {
	return errs.New("The configuration item [{0}] cannot be found", key).WithCode(ErrCodeCfgMissing)
}

func cfgBadType(key string, value any, expectedType string) error {
	return errs.New("The value of configuration item [{0}]: [{1}] cannot be parsed as [{2}]", key, value, expectedType).WithCode(ErrCodeCfgBadType)
}

const (
	ErrCodeCfgNotInitialized = "CFG_NOT_INITIALIZED"
	ErrCodeCfgMissing        = "CFG_MISSING"
	ErrCodeCfgBadType        = "CFG_BAD_TYPE"
)

var (
	ErrCfgNotInitialized = errs.New("configuration not initialized").WithCode(ErrCodeCfgNotInitialized)
)

func IsCfgMissingErr(err error) bool {
	var e *errs.Err
	if !errors.As(err, &e) {
		return false
	}
	return e.Code == ErrCodeCfgMissing
}
