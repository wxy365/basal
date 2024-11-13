package cfg

import (
	"encoding/json"
	"errors"
	"github.com/wxy365/basal/ds/maps"
	"github.com/wxy365/basal/env"
	"github.com/wxy365/basal/lei"
	"github.com/wxy365/basal/types"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func ParseJsonFile(jsonPath string) (Cfg, error) {
	file, err := os.Open(jsonPath)
	if err != nil {
		return nil, lei.Wrap("error in opening configuration json file '{0}'", err, jsonPath)
	}
	jsonBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, lei.Wrap("failed to read configuration file '{0}'", err, jsonPath)
	}
	return ParseJsonString(string(jsonBytes))
}

func ParseJsonString(jsonStr string) (Cfg, error) {
	c := make(Cfg)
	err := json.Unmarshal([]byte(jsonStr), &c)
	if err != nil {
		return nil, lei.Wrap("failed to parse the configuration content", err)
	}
	return c, nil
}

func ParseJsonDir(jsonDir string) (Cfg, error) {
	c := make(Cfg)
	err := filepath.WalkDir(jsonDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".json") {
			return nil
		}
		cfg, err := ParseJsonFile(path)
		if err != nil {
			return err
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

func (c Cfg) GetMap(key string, def ...map[string]any) (map[string]any, error) {
	return GetObj(c, key, def...)
}

func (c Cfg) GetObj(key string, t *any, def ...any) error {
	m, err := GetObj[map[string]any](c, key)
	if err != nil {
		var er *lei.Err
		if errors.As(err, &er) && er.Code == ErrCodeCfgMissing {
			if len(def) > 0 {
				*t = def[0]
				return nil
			}
			return err
		}
	}
	maps.ToObj(m, t)
	return nil
}

func GetBool(c Cfg, key string, def ...bool) (bool, error) {
	val, err := env.GetBool(key)
	if err == nil {
		return val, nil
	}
	return lookupInCfg(c, key, def...)
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
	return lookupInCfg(c, key, def...)
}

func GetObj[T any](c Cfg, key string, def ...T) (T, error) {
	val, err := env.GetObj[T](key)
	if err == nil {
		return val, nil
	}
	return lookupObjInCfg(c, key, def...)
}

func lookupInCfg[T any](m map[string]any, key string, def ...T) (T, error) {
	var t T
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
		if b, ok := res.(T); ok {
			return b, nil
		} else {
			return t, lei.New("The value of key '{0}': '{1}' cannot be parsed as {2}", key, res, reflect.TypeOf(t).Name()).WithCode(ErrCodeCfgBadType)
		}
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return t, cfgAbsent(key)
}

func lookupNumberInCfg[T types.BasicNumberUnion](m map[string]any, key string, def ...T) (T, error) {
	var t T
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
		if b, ok := res.(T); ok {
			return b, nil
		} else if f, ok := res.(float64); ok {
			// json number is always parsed as float64 in unmarshalling
			// as we want int eg., convert it in unsafe way
			return T(f), nil
		} else {
			return t, lei.New("The value of key '{0}': '{1}' cannot be parsed as {2}", key, res, reflect.TypeOf(t).Name()).WithCode(ErrCodeCfgBadType)
		}
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return t, cfgAbsent(key)
}

func lookupObjInCfg[T any](m map[string]any, key string, def ...T) (T, error) {
	var t T
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
		if b, ok := res.(T); ok {
			return b, nil
		} else if m, ok := res.(map[string]any); ok {
			tt := new(T)
			maps.ToObj[T](m, tt)
			return *tt, nil
		} else {
			return t, lei.New("The value of key '{0}': '{1}' cannot be parsed as {2}", key, res, reflect.TypeOf(t).Name()).WithCode(ErrCodeCfgBadType)
		}
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return t, cfgAbsent(key)
}

func cfgAbsent(key string) error {
	return lei.New("The configuration item '{0}' cannot be found", key).WithCode(ErrCodeCfgMissing)
}

const (
	ErrCodeCfgMissing = "CFG_MISSING"
	ErrCodeCfgBadType = "CFG_BAD_TYPE"
)
