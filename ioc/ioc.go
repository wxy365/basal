package ioc

import (
	"reflect"
	"sync"

	"github.com/wxy365/basal/ds/slices"
	"github.com/wxy365/basal/log"
)

var defaultIOC = new(IOC)

type IOC struct {
	beans            []*bean
	uncompletedBeans []*bean
	mutex            sync.Mutex
}

type bean struct {
	name  string
	value any
}

func RegisterTo[T any](ioc *IOC, t T, name ...string) {
	ioc.mutex.Lock()
	defer ioc.mutex.Unlock()
	otyp := reflect.TypeOf(t)
	oval := reflect.ValueOf(t)
	if otyp.Kind() != reflect.Pointer {
		log.Panic("Objects registered into the IOC container must be pointers to structs.")
	}
	typ := otyp.Elem()
	val := oval.Elem()
	if typ.Kind() != reflect.Struct {
		log.Panic("Objects registered into the IOC container must be pointers to structs.")
	}
	var autowireCount int
	var autowiredCount int
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		beanName, ok := field.Tag.Lookup("autowired")
		if ok {
			autowireCount++
			var depend *bean
			for _, b := range ioc.beans {
				if beanName != "" && b.name != beanName {
					continue
				}
				bTyp := reflect.TypeOf(b.value)
				if field.Type == bTyp {
					depend = b
					break
				}
				if field.Type.Kind() == reflect.Interface && bTyp.Implements(field.Type) {
					depend = b
				}
			}
			if depend != nil {
				fieldValue := val.Field(i)
				fieldValue.Set(reflect.ValueOf(depend.value))
				autowiredCount++
			}
		}
	}
	newBean := &bean{
		name: func() string {
			if len(name) > 0 {
				return name[0]
			} else {
				return typ.Name()
			}
		}(),
		value: t,
	}
	ioc.beans = append(ioc.beans, newBean)

	for i := 0; i < len(ioc.uncompletedBeans); {
		b := ioc.uncompletedBeans[i]
		bType := reflect.TypeOf(b.value).Elem()
		bValue := reflect.ValueOf(b.value).Elem()
		var restAutowireCount int
		var restAutowiredCount int
		for i := 0; i < bType.NumField(); i++ {
			fieldValue := bValue.Field(i)
			if !fieldValue.IsZero() {
				continue
			}
			field := bType.Field(i)
			if depName, ok := field.Tag.Lookup("autowired"); ok {
				restAutowireCount++
				if depName != "" && depName != newBean.name {
					continue
				}
				if field.Type != otyp && !(field.Type.Kind() == reflect.Interface && otyp.Implements(field.Type)) {
					continue
				}
				fieldValue.Set(reflect.ValueOf(newBean.value))
				restAutowiredCount++
			}
		}
		if restAutowireCount == restAutowiredCount {
			ioc.uncompletedBeans = slices.Del(ioc.uncompletedBeans, i)
		} else {
			i++
		}
	}

	if autowiredCount < autowireCount {
		ioc.uncompletedBeans = append(ioc.uncompletedBeans, newBean)
	}
}

func GetFrom[T any](ioc *IOC) (T, bool) {
	var t T
	targetType := reflect.TypeOf(t)
	for _, b := range ioc.beans {
		if bType := reflect.TypeOf(b.value); bType == targetType || bType.Implements(targetType) {
			return b.value.(T), true
		}
	}
	return t, false
}

func GetByNameFrom[T any](ioc *IOC, name string) (T, bool) {
	var t T
	targetType := reflect.TypeOf(t)
	for _, b := range ioc.beans {
		if bType := reflect.TypeOf(b.value); bType != targetType && !bType.Implements(targetType) {
			continue
		}
		if b.name == name {
			return b.value.(T), true
		}
	}
	return t, false
}

func Register[T any](t T, name ...string) {
	RegisterTo[T](defaultIOC, t, name...)
}

func Get[T any]() (T, bool) {
	return GetFrom[T](defaultIOC)
}

func GetByName[T any](name string) (T, bool) {
	return GetByNameFrom[T](defaultIOC, name)
}
