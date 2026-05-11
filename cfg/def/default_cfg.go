package def

import (
	"fmt"
	"os"
	"sync"

	"github.com/wxy365/basal/cfg"
	"github.com/wxy365/basal/env"
	"github.com/wxy365/basal/types"
)

var (
	defaultCfg cfg.Cfg
	cfgMu      sync.RWMutex
)

func init() {
	var config cfg.Cfg

	jsonStr, err := env.GetStr("APP_CFG")
	if err == nil && jsonStr != "" {
		config, err = cfg.ParseJsonString(jsonStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARN: failed to parse APP_CFG env var: %v\n", err)
		}
	}

	jsonPath, err := env.GetStr("APP_CFG_PATH", "./config/config.json")
	if err == nil {
		cfg1, err := cfg.ParseJsonFile(jsonPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARN: failed to load config file [%s]: %v\n", jsonPath, err)
		}
		if len(config) == 0 {
			config = cfg1
		} else {
			config.Merge(cfg1)
		}
	}

	jsonDir, err := env.GetStr("APP_CFG_DIR", "./config/")
	if err == nil {
		cfg1, err := cfg.LoadFromDir(jsonDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARN: failed to load config dir [%s]: %v\n", jsonDir, err)
		}
		if len(cfg1) > 0 {
			if len(config) == 0 {
				config = cfg1
			} else {
				config.Merge(cfg1)
			}
		}
	}

	defaultCfg = config
}

func SetDefault(c cfg.Cfg) {
	cfgMu.Lock()
	defaultCfg = c
	cfgMu.Unlock()
}

func SetDefaultFromFile(jsonPath string) (cfg.Cfg, error) {
	cfgMu.Lock()
	defer cfgMu.Unlock()
	var err error
	defaultCfg, err = cfg.ParseJsonFile(jsonPath)
	return defaultCfg, err
}

func SetDefaultFromStr(jsonStr string) (cfg.Cfg, error) {
	cfgMu.Lock()
	defer cfgMu.Unlock()
	var err error
	defaultCfg, err = cfg.ParseJsonString(jsonStr)
	return defaultCfg, err
}

func SetDefaultFromDir(jsonDir string) (cfg.Cfg, error) {
	cfgMu.Lock()
	defer cfgMu.Unlock()
	var err error
	defaultCfg, err = cfg.LoadFromDir(jsonDir)
	return defaultCfg, err
}

func HasDefault() bool {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return len(defaultCfg) > 0
}

func GetBool(key string, def ...bool) (bool, error) {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return cfg.GetBool(defaultCfg, key, def...)
}

func GetInt[T types.IntUnion](key string, def ...T) (T, error) {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return cfg.GetInt(defaultCfg, key, def...)
}

func GetUint[T types.UintUnion](key string, def ...T) (T, error) {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return cfg.GetUint(defaultCfg, key, def...)
}

func GetFloat[T types.FloatUnion](key string, def ...T) (T, error) {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return cfg.GetFloat(defaultCfg, key, def...)
}

func GetStr(key string, def ...string) (string, error) {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return cfg.GetStr(defaultCfg, key, def...)
}

func GetObj[T any](key string, def ...T) (T, error) {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return cfg.GetObj(defaultCfg, key, def...)
}

func Clear() {
	cfgMu.Lock()
	defer cfgMu.Unlock()
	defaultCfg.Clear()
}

func GetCfg() cfg.Cfg {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return defaultCfg
}
