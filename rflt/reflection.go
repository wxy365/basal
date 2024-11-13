package rflt

import (
	"encoding/json"
	"fmt"
	"github.com/wxy365/basal/lei"
	"reflect"
)

func SetFieldValue[T any](target any, fieldName string, fieldValue T) error {
	v, err := shouldBePointerToStruct(target)
	if err != nil {
		return err
	}
	fv := v.FieldByName(fieldName)
	newVal := reflect.ValueOf(fieldValue)
	if newVal.IsZero() {
		fv.SetZero()
		return nil
	}
	if fv.Kind() == reflect.Pointer && newVal.Kind() != reflect.Pointer {
		newVal = reflect.ValueOf(&fieldValue)
	} else if fv.Kind() != reflect.Pointer && newVal.Kind() == reflect.Pointer {
		newVal = newVal.Elem()
	}
	fv.Set(newVal)
	return nil
}

func SetFieldValueAny(target any, name string, fieldValue any) error {
	v, err := shouldBePointerToStruct(target)
	if err != nil {
		return err
	}
	fv := v.FieldByName(name)
	fv.Set(reflect.ValueOf(fieldValue))
	return nil
}

func SetFieldIValue[T any](target any, i int, fieldValue T) error {
	v, err := shouldBePointerToStruct(target)
	if err != nil {
		return err
	}
	fv := v.Field(i)
	newVal := reflect.ValueOf(fieldValue)
	if newVal.IsZero() {
		fv.SetZero()
		return nil
	}
	if fv.Kind() == reflect.Pointer && newVal.Kind() != reflect.Pointer {
		newVal = reflect.ValueOf(&fieldValue)
	} else if fv.Kind() != reflect.Pointer && newVal.Kind() == reflect.Pointer {
		newVal = newVal.Elem()
	}
	fv.Set(newVal)
	return nil
}

func SetFieldIValueAny(target any, i int, fieldValue any) error {
	v, err := shouldBePointerToStruct(target)
	if err != nil {
		return err
	}
	fv := v.Field(i)
	fv.Set(reflect.ValueOf(fieldValue))
	return nil
}

func UnmarshalValue(target reflect.Value, from string) error {
	k := target.Kind()
	if k == reflect.String {
		from = "\"" + from + "\""
		target = target.Addr()
	} else if k == reflect.Pointer {
		inner := target.Elem()
		if inner.Kind() == reflect.String {
			from = "\"" + from + "\""
		}
	} else {
		target = target.Addr()
	}

	return json.Unmarshal([]byte(from), target.Interface())
}

func ValueToString(from reflect.Value) (string, error) {
	for from.Kind() == reflect.Pointer {
		from = from.Elem()
	}
	switch from.Kind() {
	case reflect.Bool:
		if from.Bool() {
			return "true", nil
		} else {
			return "false", nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", from.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fmt.Sprintf("%d", from.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", from.Float()), nil
	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%g", from.Complex()), nil
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
		raw, err := json.Marshal(from.Interface())
		return string(raw), err
	default:
		if from.CanInterface() {
			return fmt.Sprintf("%+v", from.Interface()), nil
		} else {
			return "", lei.New("Cannot convert value of {0} kind to string", from.Kind().String())
		}
	}
}

func shouldBePointerToStruct(target any) (reflect.Value, error) {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Pointer {
		return v, lei.New("The target must be a pointer to struct")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return v, lei.New("The target must be a pointer to struct")
	}
	return v, nil
}
