package cfg

import (
	"github.com/wxy365/basal/env"
	"github.com/wxy365/basal/types"
)

var defaultCfg Cfg

func init() {
	var cfg Cfg
	jsonPath, err := env.GetStr("BASAL_CFG_PATH", "./cfg/cfg.json")
	if err == nil {
		cfg, _ = ParseJsonFile(jsonPath)
	}

	jsonDir, err := env.GetStr("BASAL_CFG_DIR", "./cfg/")
	if err == nil {
		cfg1, _ := ParseJsonDir(jsonDir)
		if len(cfg1) > 0 {
			if len(cfg) == 0 {
				cfg = cfg1
			} else {
				cfg.Merge(cfg1)
			}
		}
	}
}

func SetDefault(cfg Cfg) {
	defaultCfg = cfg
}

func SetDefaultFromFile(jsonPath string) (Cfg, error) {
	var err error
	defaultCfg, err = ParseJsonFile(jsonPath)
	return defaultCfg, err
}

func SetDefaultFromStr(jsonStr string) (Cfg, error) {
	var err error
	defaultCfg, err = ParseJsonString(jsonStr)
	return defaultCfg, err
}

func SetDefaultFromDir(jsonDir string) (Cfg, error) {
	var err error
	defaultCfg, err = ParseJsonDir(jsonDir)
	return defaultCfg, err
}

func HasDefault() bool {
	return len(defaultCfg) > 0
}

func GetBoolDef(key string, def ...bool) (bool, error) {
	return GetBool(defaultCfg, key, def...)
}

func GetIntDef[T types.IntUnion](key string, def ...T) (T, error) {
	return GetInt(defaultCfg, key, def...)
}

func GetUintDef[T types.UintUnion](key string, def ...T) (T, error) {
	return GetUint(defaultCfg, key, def...)
}

func GetFloatDef[T types.FloatUnion](key string, def ...T) (T, error) {
	return GetFloat(defaultCfg, key, def...)
}

func GetStrDef(key string, def ...string) (string, error) {
	return GetStr(defaultCfg, key, def...)
}

func GetObjDef[T any](key string, def ...T) (T, error) {
	return GetObj(defaultCfg, key, def...)
}
