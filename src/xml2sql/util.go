package xml2sql

import (
	"bytes"
	"errors"
	"fmt"
	"methods"
	"reflect"
	"regexp"
	"strings"

	//"strconv"
	"dbdefault"
	"text/template"
	"time"
	"timex"
)

const NIL = "<nil>"

// const NULL = "<NULL>"

var def = timex.Time(time.Time{})

//获取 # # token
var reg = regexp.MustCompile(`#[\s\S]+?#`)

//匹配空元素
var regnil, _ = regexp.Compile("<nil>")

//空元素不再 <nil>号中
var regnil1, _ = regexp.Compile(`"[\s\S]+?<nil>[\s\S]+?"`)

func GenerateSql(sql string, data interface{}) (string, error) {
	tokens := GetToken(sql)
	mtok, err := ReplaceToke(tokens, data)
	if err != nil {
		return "", err
	}
	return MakeSql(sql, mtok), nil
}

//获取标识
func GetToken(sql string) []string {
	return reg.FindAllString(sql, -1)
}

//替换
func ReplaceToke(tokens []string, data interface{}) (map[string]string, error) {
	m := make(map[string]string)
	buf := bytes.NewBuffer(nil)
	for _, v := range tokens {
		temp := template.New("")

		temp.Funcs(funcMap)
		temp, err := temp.Parse(v)
		if err != nil {
			return m, err
		}
		err = temp.Execute(buf, data)
		if err != nil {
			return m, err
		}
		m[v] = string(buf.Bytes())
		buf.Reset()
	}
	return m, nil
}

//替换
func MakeSql(sql string, m map[string]string) string {
	for k, v := range m {
		if regnil.FindString(v) != "" && regnil1.FindString(v) == "" {
			sql = strings.Replace(sql, k, "", -1)
		} else {
			v = strings.Trim(v, "#")
			sql = strings.Replace(sql, k, v, -1)
		}
	}
	sql = strings.Replace(sql, "\r", " ", -1)
	sql = strings.Replace(sql, "\n", " ", -1)
	//tracebug.Trace(sql)
	//fmt.Println(sql)
	return sql
}

var funcMap template.FuncMap

func init() {
	funcMap = template.FuncMap{
		//	"notEmpty":   NotEmpty,
		"empty": Empty,
		//"allowEmpty": AllowEmpty,
		"like":    Like,
		"inOfStr": InOfStr,
		//"null":       null,
		"leftlike": LeftLike,
		"oraDate":  oracleDate,
		//"testeq" : testeq,
		"raw":              raw,
		"test":             test,
		"trim00":           trim00,
		"isempty":          isEmpty,
		"split":            split,
		"escape":           escape,
	}
}

func escape(data interface{}) (interface{}, error) {
	s, ok := data.(string)
	if !ok {
		return "", nil
	}
	return strings.Replace(s, "'", "''", -1), nil

}
func split(data interface{}) (interface{}, error) {
	s, ok := data.(string)
	if !ok {
		return nil, errors.New(`not valid string`)
	}
	s = strings.Trim(s, ",")
	if dbdefault.IsStringNullOrEmpty(s) {
		return []string{}, nil
	}
	items := strings.Split(s, ",")
	return items, nil

}

func isEmpty(a string) bool {
	return dbdefault.IsStringNullOrEmpty(a)
}

func trim00(a string) string {
	str := strings.TrimRight(a, "0")
	if len(str)%2 == 1 {
		str += "0"
	}
	return str
}

func test(ok bool, str ...string) string {
	if len(str) < 1 {
		return NIL
	}
	if ok {
		return str[0]
	}
	if len(str) > 1 {
		return str[1]
	}
	return NIL
}

func raw(arg interface{}) (s string, err error) {
	if arg == nil {
		return NIL, nil
	}
	value := reflect.ValueOf(arg)
	if value.Type().Kind() == reflect.Ptr && value.IsNil() {
		return NIL, nil
	}
	for {
		if value.Type().Kind() == reflect.Ptr {
			value = value.Elem()
		} else {
			break
		}
	}
	val := value.Interface()
	switch val.(type) {
	case string:
		v := val.(string)
		if dbdefault.IsStringNullOrEmpty(v) {
			return NIL, nil
		}
		return v, nil

	default:
		return "", fmt.Errorf("not support type %T", val)
	}
}


func oracleDate(arg interface{}) (s string, err error) {
	if arg == nil {
		return NIL, nil
	}
	value := reflect.ValueOf(arg)
	if value.Type().Kind() == reflect.Ptr && value.IsNil() {
		return NIL, nil
	}

	for {
		if value.Type().Kind() == reflect.Ptr {
			value = value.Elem()
		} else {
			break
		}
	}
	vI, ok := value.Interface().(timex.Time)
	if !ok {
		vT, ok := value.Interface().(time.Time)
		if !ok {
			return NIL, nil
		}
		vI = timex.Time(vT)
	}
	//to_date('{{empty.Accept_Starttime}}
	//','
	//yyyy - mm - dd
	//hh24:mi:ss')#
	if dbdefault.IsTimeNull(vI) {
		return NIL, nil
	}
	return `to_date('` + vI.ToString() + `','yyyy-mm-dd hh24:mi:ss')`, nil

}

var timeNames = []string{
	"Time",
}

func isTime(name string) bool {
	for _, v := range timeNames {
		if v == name {
			return true
		}
	}
	return false
}

//InOfStr 字句
func InOfStr(arg interface{}) (string, error) {
	s, ok := arg.(string)
	if !ok {
		return NIL, nil
	}
	if dbdefault.IsStringNullOrEmpty(s) {
		return NIL, nil
	}

	sarr := strings.Split(s, ",")

	return methods.MakeInClause(sarr), nil
}

//Like 函数
func Like(arg interface{}) (string, error) {
	if arg == nil {
		return NIL, nil
	}

	value := reflect.ValueOf(arg)

	if value.Type().Kind() == reflect.Ptr && value.IsNil() {
		return NIL, nil
	}

	for {
		if value.Type().Kind() == reflect.Ptr {
			value = value.Elem()
		} else {
			break
		}
	}
	val := value.Interface()
	switch val.(type) {
	case string:
		v := val.(string)
		if dbdefault.IsStringNull(v) {
			return NIL, nil
		}

		if v == "" {
			return NIL, nil
		}

		v = strings.Replace(v, `'`, `''`, -1)
		return `like '%` + v + `%'`, nil
	default:
		return "", fmt.Errorf("not support type %T", val)
	}
}

func LeftLike(arg interface{}) (string, error) {
	if arg == nil {
		return NIL, nil
	}

	value := reflect.ValueOf(arg)

	if value.Type().Kind() == reflect.Ptr && value.IsNil() {
		return NIL, nil
	}

	for {
		if value.Type().Kind() == reflect.Ptr {
			value = value.Elem()
		} else {
			break
		}
	}
	val := value.Interface()
	switch val.(type) {
	case string:
		v := val.(string)
		if dbdefault.IsStringNull(v) {
			return NIL, nil
		}

		if v == "" {
			return NIL, nil
		}

		v = strings.Replace(v, `'`, `''`, -1)
		return `like '` + v + `%'`, nil
	default:
		return "", fmt.Errorf("not support type %T", val)
	}
}

//Empty 去除空的
func Empty(arg interface{}) (interface{}, error) {
	if arg == nil {
		return NIL, nil
	}
	value := reflect.ValueOf(arg)
	if value.Type().Kind() == reflect.Ptr && value.IsNil() {
		return NIL, nil
	}
	for {
		if value.Type().Kind() == reflect.Ptr {
			value = value.Elem()
		} else {
			break
		}
	}
	val := value.Interface()
	switch val.(type) {
	case []byte:
		return NIL, nil
	case int64:
		v := val.(int64)
		if dbdefault.IsIntNull(v) {
			return NIL, nil
		}
		return v, nil
	case string:
		v := val.(string)
		if dbdefault.IsStringNullOrEmpty(v) {
			return NIL, nil
		}

		return EscapeQuotes(v), nil
	case float64:
		v := val.(float64)
		if dbdefault.IsFloat64Null(v) {
			return NIL, nil
		}

		return val.(float64), nil
	case timex.Time:
		v := val.(timex.Time)
		if dbdefault.IsTimeNull(v) {
			return NIL, nil
		}
		if dbdefault.SYSDATE.Equal(v) {
			return timex.Time(time.Now()).ToString(), nil
		}
		return v.ToString(), nil
	default:
		return nil, fmt.Errorf("not support type %T", val)
	}
}

//AllowEmpty 允许为空
func AllowEmpty(arg interface{}) (interface{}, error) {
	if arg == nil {
		return nil, errors.New("nil error")
	}
	value := reflect.ValueOf(arg)
	if value.Type().Kind() == reflect.Ptr && value.IsNil() {
		return NIL, nil
	}
	for {
		if value.Type().Kind() == reflect.Ptr {
			value = value.Addr()
		} else {
			break
		}
	}

	val := value.Interface()
	switch val.(type) {
	case []byte:
		return NIL, nil
	case int64:
		v := val.(int64)
		return v, nil
	case string:
		v := val.(string)
		return EscapeQuotes(v), nil
	case float64:
		return val.(float64), nil
	case timex.Time:
		v := val.(timex.Time)
		return v.ToString(), nil
	default:
		return nil, fmt.Errorf("not support type %T", val)
	}
}
