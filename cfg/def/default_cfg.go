package def

import (
	"github.com/wxy365/basal/cfg"
	"github.com/wxy365/basal/env"
	"github.com/wxy365/basal/types"
)

var defaultCfg cfg.Cfg

func init() {
	var config cfg.Cfg

	jsonStr, err := env.GetStr("BASAL_CFG")
	if err == nil && jsonStr != "" {
		config, _ = cfg.ParseJsonString(jsonStr)
	}

	jsonPath, err := env.GetStr("BASAL_CFG_PATH", "./config/config.json")
	if err == nil {
		cfg1, _ := cfg.ParseJsonFile(jsonPath)
		if len(config) == 0 {
			config = cfg1
		} else {
			config.Merge(cfg1)
		}
	}

	jsonDir, err := env.GetStr("BASAL_CFG_DIR", "./config/")
	if err == nil {
		cfg1, _ := cfg.ParseJsonDir(jsonDir)
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

func SetDefault(cfg cfg.Cfg) {
	defaultCfg = cfg
}

func SetDefaultFromFile(jsonPath string) (cfg.Cfg, error) {
	var err error
	defaultCfg, err = cfg.ParseJsonFile(jsonPath)
	return defaultCfg, err
}

func SetDefaultFromStr(jsonStr string) (cfg.Cfg, error) {
	var err error
	defaultCfg, err = cfg.ParseJsonString(jsonStr)
	return defaultCfg, err
}

func SetDefaultFromDir(jsonDir string) (cfg.Cfg, error) {
	var err error
	defaultCfg, err = cfg.ParseJsonDir(jsonDir)
	return defaultCfg, err
}

func HasDefault() bool {
	return len(defaultCfg) > 0
}

func GetBool(key string, def ...bool) (bool, error) {
	return cfg.GetBool(defaultCfg, key, def...)
}

func GetInt[T types.IntUnion](key string, def ...T) (T, error) {
	return cfg.GetInt(defaultCfg, key, def...)
}

func GetUint[T types.UintUnion](key string, def ...T) (T, error) {
	return cfg.GetUint(defaultCfg, key, def...)
}

func GetFloat[T types.FloatUnion](key string, def ...T) (T, error) {
	return cfg.GetFloat(defaultCfg, key, def...)
}

func GetStr(key string, def ...string) (string, error) {
	return cfg.GetStr(defaultCfg, key, def...)
}

func GetObj[T any](key string, def ...T) (T, error) {
	return cfg.GetObj(defaultCfg, key, def...)
}
