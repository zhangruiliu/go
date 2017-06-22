package methods

import (
	"errors"
	"reflect"
)

var (
	nilErr       = errors.New("from or to parameter can't be nil")
	notPtrErr    = errors.New("from or to not a pointer")
	notStructErr = errors.New("from or to not a pointer to struct")
)

//FieldCopy 字段赋值 不要使用
func FieldCopy(from interface{}, to interface{}) error {
	if from == nil || to == nil {
		return nilErr
	}
	fv := reflect.ValueOf(from)
	tv := reflect.ValueOf(to)
	if fv.Kind() != reflect.Ptr || tv.Kind() != reflect.Ptr {
		return notPtrErr
	}

	fv = fv.Elem()
	tv = tv.Elem()

	if fv.Kind() != reflect.Struct || tv.Kind() != reflect.Struct {
		return notStructErr
	}

	fvtype := fv.Type()
	fieldNum := fvtype.NumField()
	for i := 0; i < fieldNum; i++ {
		fvFieldName := fvtype.Field(i).Name
		fvField := fv.Field(i)
		fvFieldType := fvField.Type()

		tvField := tv.FieldByName(fvFieldName)
		if !tvField.IsValid() {
			continue
		}
		tvFieldType := tvField.Type()
		if tvFieldType.Name() != fvFieldType.Name() {
			continue
		}

		if tvField.CanSet() {
			tvField.Set(fvField)
		}
	}

	return nil
}
