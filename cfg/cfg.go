package cfg

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/wxy365/basal/ds/maps"
	"github.com/wxy365/basal/env"
	"github.com/wxy365/basal/errs"
	"github.com/wxy365/basal/fn"
	"github.com/wxy365/basal/types"
	"gopkg.in/yaml.v3"
)

func ParseJsonFile(jsonPath string) (Cfg, error) {
	file, err := os.Open(jsonPath)
	if err != nil {
		return nil, errs.Wrap(err, "error in opening configuration json file [{0}]", jsonPath)
	}
	jsonBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, errs.Wrap(err, "failed to read configuration file [{0}]", jsonPath)
	}
	return ParseJsonString(string(jsonBytes))
}

func ParseYamlFile(yamlPath string) (Cfg, error) {
	file, err := os.Open(yamlPath)
	if err != nil {
		return nil, errs.Wrap(err, "error in opening configuration yaml file [{0}]", yamlPath)
	}
	yamlBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, errs.Wrap(err, "failed to read configuration file [{0}]", yamlPath)
	}
	return ParseYamlString(string(yamlBytes))
}

func ParseJsonString(jsonStr string) (Cfg, error) {
	c := make(Cfg)
	err := json.Unmarshal([]byte(jsonStr), &c)
	if err != nil {
		return nil, errs.Wrap(err, "failed to parse the configuration content")
	}
	return c, nil
}

func ParseYamlString(yamlStr string) (Cfg, error) {
	c := make(Cfg)
	err := yaml.Unmarshal([]byte(yamlStr), &c)
	if err != nil {
		return nil, errs.Wrap(err, "failed to parse the configuration content")
	}
	return c, nil
}

func LoadFromDir(dir string) (Cfg, error) {
	c := make(Cfg)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		var cfg Cfg
		var err error
		if strings.HasSuffix(d.Name(), ".json") {
			cfg, err = ParseJsonFile(path)
		} else if strings.HasSuffix(d.Name(), ".yaml") || strings.HasSuffix(d.Name(), ".yml") {
			cfg, err = ParseYamlFile(path)
		} else {
			return nil
		}
		if err != nil {
			return errs.Wrap(err, "failed to load configuration file [{0}]", path)
		}

		c.Merge(cfg)
		return nil
	})
	return c, err
}

type Cfg map[string]any

func (c Cfg) Merge(another Cfg) Cfg {
	return maps.Merge(c, another)
}

func (c Cfg) GetBool(key string, def ...bool) (bool, error) {
	return GetBool(c, key, def...)
}

func (c Cfg) GetInt(key string, def ...int) (int, error) {
	return GetInt(c, key, def...)
}

func (c Cfg) GetInt8(key string, def ...int8) (int8, error) {
	return GetInt(c, key, def...)
}

func (c Cfg) GetInt16(key string, def ...int16) (int16, error) {
	return GetInt(c, key, def...)
}

func (c Cfg) GetInt32(key string, def ...int32) (int32, error) {
	return GetInt(c, key, def...)
}

func (c Cfg) GetInt64(key string, def ...int64) (int64, error) {
	return GetInt(c, key, def...)
}

func (c Cfg) GetUint(key string, def ...uint) (uint, error) {
	return GetUint(c, key, def...)
}

func (c Cfg) GetUint8(key string, def ...uint8) (uint8, error) {
	return GetUint(c, key, def...)
}

func (c Cfg) GetUint16(key string, def ...uint16) (uint16, error) {
	return GetUint(c, key, def...)
}

func (c Cfg) GetUint32(key string, def ...uint32) (uint32, error) {
	return GetUint(c, key, def...)
}

func (c Cfg) GetUint64(key string, def ...uint64) (uint64, error) {
	return GetUint(c, key, def...)
}

func (c Cfg) GetFloat32(key string, def ...float32) (float32, error) {
	return GetFloat(c, key, def...)
}

func (c Cfg) GetFloat64(key string, def ...float64) (float64, error) {
	return GetFloat(c, key, def...)
}

func (c Cfg) GetStr(key string, def ...string) (string, error) {
	return GetStr(c, key, def...)
}

func (c Cfg) Clear() {
	clear(c)
}

func GetBool(c Cfg, key string, def ...bool) (bool, error) {
	val, err := env.GetBool(key)
	if err == nil {
		return val, nil
	}
	return lookupInCfg(c, key, func(a any) (bool, bool) {
		b, ok := a.(bool)
		return b, ok
	}, def...)
}

func GetInt[T types.IntUnion](c Cfg, key string, def ...T) (T, error) {
	val, err := env.GetInt[T](key)
	if err == nil {
		return val, nil
	}
	return lookupNumberInCfg(c, key, def...)
}

func GetUint[T types.UintUnion](c Cfg, key string, def ...T) (T, error) {
	val, err := env.GetUint[T](key)
	if err == nil {
		return val, nil
	}
	return lookupNumberInCfg(c, key, def...)
}

func GetFloat[T types.FloatUnion](c Cfg, key string, def ...T) (T, error) {
	val, err := env.GetFloat[T](key)
	if err == nil {
		return val, nil
	}
	return lookupNumberInCfg(c, key, def...)
}

func GetStr(c Cfg, key string, def ...string) (string, error) {
	val, err := env.GetStr(key)
	if err == nil {
		return val, nil
	}
	return lookupInCfg(c, key,
		func(a any) (string, bool) {
			return fmt.Sprintf("%v", a), true
		},
		def...,
	)
}

func GetObj[T any](c Cfg, key string, def ...T) (T, error) {
	val, err := env.GetObj[T](key)
	if err == nil {
		return val, nil
	}
	return lookupObjInCfg(c, key, def...)
}

func Clear(c Cfg) {
	clear(c)
}

func lookupInCfg[T any](m map[string]any, key string, converter fn.TryFunction[any, T], def ...T) (T, error) {
	var t T
	if len(m) == 0 {
		return t, ErrCfgNotInitialized
	}
	parts := strings.Split(key, ".")
	for i := 0; i < len(parts)-1; i++ {
		if v, found := m[parts[i]]; found {
			var ok bool
			if m, ok = v.(map[string]any); !ok {
				return t, cfgAbsent(key)
			}
		} else {
			return t, cfgAbsent(key)
		}
	}

	if res, found := m[parts[len(parts)-1]]; found {
		if t1, ok := res.(T); ok {
			return t1, nil
		} else if t, ok = converter(res); ok {
			return t, nil
		} else {
			return t, cfgBadType(key, res, reflect.TypeOf(t).Name())
		}
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return t, cfgAbsent(key)
}

func lookupNumberInCfg[T types.BasicNumberUnion](m map[string]any, key string, def ...T) (T, error) {
	return lookupInCfg(
		m,
		key,
		func(a any) (T, bool) {
			if f, ok := a.(float64); ok {
				return T(f), true
			}
			var t T
			return t, false
		},
		def...,
	)
}

func lookupObjInCfg[T any](m map[string]any, key string, def ...T) (T, error) {
	return lookupInCfg(
		m,
		key,
		func(a any) (T, bool) {
			switch mp := a.(type) {
			case map[string]any:
				newT := new(T)
				maps.ToObj(mp, newT)
				return *newT, true
			case []any:
				newT := new(T)
				maps.ToObj(mp, newT)
				return *newT, true
			}
			var t T
			return t, false
		},
		def...,
	)
}
