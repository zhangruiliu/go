package methods

import (
	"dbdefault"
	"fmt"
	"path"
	"reflect"
	"strings"
	"sync"
)

func DefaultValueOtherSystem(data interface{}) error {
	if data == nil {
		return fmt.Errorf(`nil can't be a paramater`)
	}
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		return fmt.Errorf(`parameter must be a pointer`)
	}

	refValue := reflect.ValueOf(data)
	if refValue.Type().Kind() == reflect.Ptr || refValue.Type().Kind() == reflect.Interface {
		refValue = refValue.Elem()
	}

	if !refValue.IsValid() {
		return fmt.Errorf(`not a valid value`)
	}
	if refValue.Type().Kind() != reflect.Struct {
		return fmt.Errorf(`pointer must point to a struct `)
	}

	fieldPos := getTypeFieldPos(refValue)

	for _, v := range fieldPos {
		tmpV := refValue.FieldByIndex(v[0])
		switch tmpV.Type().Kind() {
		case reflect.Int64:
			if dbdefault.IsIntNull(tmpV.Int()) {
				tmpV.SetInt(0)
			}
		case reflect.Float64:
			if dbdefault.IsFloat64Null(tmpV.Float()) {
				tmpV.SetFloat(0)
			}
		case reflect.String:
			if dbdefault.IsStringNullOrEmpty(tmpV.String()) {
				tmpV.SetString("")
			}
		}
	}

	return nil
}

func DefaultValue(data interface{}) error {
	//
	if data == nil {
		return fmt.Errorf(`nil can't be a paramater`)
	}
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		return fmt.Errorf(`parameter must be a pointer`)
	}

	refValue := reflect.ValueOf(data)
	if refValue.Type().Kind() == reflect.Ptr || refValue.Type().Kind() == reflect.Interface {
		refValue = refValue.Elem()
	}

	if !refValue.IsValid() {
		return fmt.Errorf(`not a valid value`)
	}
	if refValue.Type().Kind() != reflect.Struct {
		return fmt.Errorf(`pointer must point to a struct `)
	}

	fieldPos := getTypeFieldPos(refValue)

	for _, v := range fieldPos {
		tmpV := refValue.FieldByIndex(v[0])
		switch tmpV.Type().Kind() {
		case reflect.Int64:
			tmpV.SetInt(dbdefault.DbIntNull)
		case reflect.Float64:
			tmpV.SetFloat(dbdefault.DbFloat64Null)
		case reflect.String:
			tmpV.SetString(dbdefault.DbStringNull)
		case reflect.Struct:
			if isTime(tmpV.Type().String()) {
				tmpV.Set(reflect.ValueOf(dbdefault.DbTimeNull))
			}
		}
	}

	return nil
}

var cachedStruct = struct {
	//l   sync.Mutex
	rwl sync.RWMutex
	m   map[string]map[string][][]int
}{
	//l:sync.Mutex{},
	rwl: sync.RWMutex{},
	m:   make(map[string]map[string][][]int),
}


//获取反射结构
func getTypeFieldPos(value reflect.Value) map[string][][]int {
	tp := value.Type()
	pkg := tp.PkgPath()
	name := tp.Name()
	qs := path.Join(pkg, name)

	cachedStruct.rwl.Lock()
	fpos, ok := cachedStruct.m[qs]
	cachedStruct.rwl.Unlock()
	if !ok {
		cachedStruct.rwl.Lock()
		if fpos, ok = cachedStruct.m[qs]; ok {
			cachedStruct.rwl.Unlock()
			return fpos
		}
		m := make(map[string][][]int)
		handle(value, m, []int{})
		cachedStruct.m[qs] = m
		fpos = m
		cachedStruct.rwl.Unlock()
	}

	return fpos
}

var tagName = "sql"
var tagValue = "field"

//生成结构
func handle(value reflect.Value, m map[string][][]int, dep []int) {
	tp := value.Type()
	numField := tp.NumField()

	arr := makeArr(dep)
	for i := 0; i < numField; i++ {
		field := tp.Field(i)
		fieldValue := value.Field(i)

		// if field.Tag.Get(tagName) != tagValue {
		// 	continue
		// }

		if !fieldValue.CanSet() {
			continue
		}

		fieldName := strings.ToLower(field.Name)

		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			handle(fieldValue, m, makeNArr(arr, i))
		}

		switch field.Type.Kind() {
		case reflect.Int64, reflect.Float64, reflect.String:
			m[fieldName] = append(m[fieldName], makeNArr(arr, i))
		case reflect.Struct:
			if isTime(field.Type.Name()) {
				m[fieldName] = append(m[fieldName], makeNArr(arr, i))
			}
		case reflect.Slice, reflect.Array:
			if field.Type.Elem().Kind() == reflect.Uint8 {
				m[fieldName] = append(m[fieldName], makeNArr(arr, i))
			}
		}
	}

	return
}

func makeArr(arr []int) []int {
	ret := make([]int, len(arr))
	copy(ret, arr)
	return ret
}
func makeNArr(arr []int, n int) []int {
	ret := make([]int, len(arr)+1)
	copy(ret, arr)
	ret[len(arr)] = n
	return ret
}

//timeNames 所有的时间类型
var timeNames = []string{
	"Time",
}

//是否是时间类型
func isTime(name string) bool {
	for _, v := range timeNames {
		if v == name {
			return true
		}
	}
	return false
}
