package env

import (
	"encoding/json"
	"github.com/wxy365/basal/lei"
	"os"
	"strconv"
	"unsafe"
)

func GetBool(k string, def ...bool) (bool, error) {
	v, found := os.LookupEnv(k)
	if found {
		return strconv.ParseBool(v)
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return false, errEnvAbsent(k)
}

func GetInt[T any](k string, def ...T) (T, error) {
	v, found := os.LookupEnv(k)
	var t T
	if found {
		size := unsafe.Sizeof(t)
		r, e := strconv.ParseInt(v, 10, int(size*8))
		if e != nil {
			return t, errEnvNotNumber(k, e)
		}
		t = *(*T)(unsafe.Pointer(&r))
		return t, nil
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return t, errEnvAbsent(k)
}

func GetUint[T any](k string, def ...T) (T, error) {
	v, found := os.LookupEnv(k)
	var t T
	if found {
		size := unsafe.Sizeof(t)
		r, e := strconv.ParseUint(v, 10, int(size*8))
		if e != nil {
			return t, errEnvNotNumber(k, e)
		}
		t = *(*T)(unsafe.Pointer(&r))
		return t, nil
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return t, errEnvAbsent(k)
}

func GetFloat[T any](k string, def ...T) (T, error) {
	v, found := os.LookupEnv(k)
	var t T
	if found {
		size := unsafe.Sizeof(t)
		r, e := strconv.ParseFloat(v, int(size*8))
		if e != nil {
			return t, errEnvNotNumber(k, e)
		}
		t = *(*T)(unsafe.Pointer(&r))
		return t, e
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return t, errEnvAbsent(k)
}

func GetStr(k string, def ...string) (string, error) {
	v, found := os.LookupEnv(k)
	if found {
		return v, nil
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return "", errEnvAbsent(k)
}

func GetObj[T any](k string, def ...T) (T, error) {
	v, found := os.LookupEnv(k)
	var t T
	if found {
		err := json.Unmarshal([]byte(v), &t)
		if err != nil {
			return t, lei.Wrap("Unable to deserialize object from environment {0}", err, k)
		}
		return t, nil
	}
	if len(def) > 0 {
		return t, nil
	}
	return t, errEnvAbsent(k)
}

func errEnvAbsent(k string) error {
	return lei.New("Environment {0} cannot be found", k)
}

func errEnvNotNumber(k string, cause error) error {
	return lei.Wrap("Environment {0} cannot be parsed as number", cause, k)
}
