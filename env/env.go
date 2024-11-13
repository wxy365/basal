package env

import (
	"encoding/json"
	"github.com/wxy365/basal/lei"
	"github.com/wxy365/basal/types"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

func GetBool(k string, def ...bool) (bool, error) {
	v, found := lookup(k)
	if found {
		r, err := strconv.ParseBool(v)
		if err != nil {
			return false, lei.New("The value of environment variable '{0}': '{1}' cannot be parsed as bool", k, v).WithCode(ErrCodeEnvBadType)
		}
		return r, nil
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return false, errEnvAbsent(k)
}

func GetInt[T types.IntUnion](k string, def ...T) (T, error) {
	v, found := lookup(k)
	var t T
	if found {
		size := unsafe.Sizeof(t)
		r, err := strconv.ParseInt(v, 10, int(size*8))
		if err != nil {
			return t, errEnvNotNumber(k, v, err)
		}
		return T(r), nil
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return t, errEnvAbsent(k)
}

func GetUint[T types.UintUnion](k string, def ...T) (T, error) {
	v, found := lookup(k)
	var t T
	if found {
		size := unsafe.Sizeof(t)
		r, err := strconv.ParseUint(v, 10, int(size*8))
		if err != nil {
			return t, errEnvNotNumber(k, v, err)
		}
		t = T(r)
		return t, nil
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return t, errEnvAbsent(k)
}

func GetFloat[T types.FloatUnion](k string, def ...T) (T, error) {
	v, found := lookup(k)
	var t T
	if found {
		size := unsafe.Sizeof(t)
		r, err := strconv.ParseFloat(v, int(size*8))
		if err != nil {
			return t, errEnvNotNumber(k, v, err)
		}
		return T(r), err
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return t, errEnvAbsent(k)
}

func GetStr(k string, def ...string) (string, error) {
	v, found := lookup(k)
	if found {
		return v, nil
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return "", errEnvAbsent(k)
}

func GetObj[T any](k string, def ...T) (T, error) {
	v, found := lookup(k)
	var t T
	if found {
		err := json.Unmarshal([]byte(v), &t)
		if err != nil {
			return t, lei.Wrap("Unable to deserialize object from environment variable {0}", err, k)
		}
		return t, nil
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return t, errEnvAbsent(k)
}

func lookup(k string) (string, bool) {
	v, found := os.LookupEnv(k)
	if found {
		return v, true
	}
	if strings.Contains(k, ".") {
		k = strings.ToUpper(strings.ReplaceAll(k, ".", "_"))
		v, found = os.LookupEnv(k)
		if found {
			return v, true
		}
	}
	return v, false
}

func errEnvAbsent(k string) error {
	return lei.New("Environment variable '{0}' is not found", k).WithCode(ErrCodeEnvMissing)
}

func errEnvNotNumber(k string, v string, cause error) error {
	return lei.Wrap("The value of environment variable '{0}': '{1}' cannot be parsed as number", cause, k, v).WithCode(ErrCodeEnvBadType)
}

const (
	ErrCodeEnvMissing = "ENV_MISSING"
	ErrCodeEnvBadType = "ENV_BAD_TYPE"
)
