package cfg

import (
	"github.com/wxy365/basal/types"
)

var defaultCfg Cfg

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
