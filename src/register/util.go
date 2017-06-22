package register

import (
	"net/http"
	"reflect"

	"web/httpcontext"

	"github.com/julienschmidt/httprouter"
)

// MakeTypeHttpRouter 根据类型构造httprouter.Handle
func MakeTypeHttpRouter(tp reflect.Type) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, m httprouter.Params) {
		if len(m) < 1 {
			panic(" not valid reflect.Type call")
		}

		methodName := m.ByName(httpcontext.InvokeMethodName)

		value := reflect.New(tp)
		initMethod := value.MethodByName(httpcontext.InitMethod)
		if initMethod.IsValid() {
			initMethod.Call([]reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r), reflect.ValueOf(m)})
		} else {
			panic(" no " + httpcontext.InitMethod + " method " + " of type " + tp.Name())
		}

		callmethod := value.MethodByName(methodName)
		if callmethod.IsValid() {
			callmethod.Call([]reflect.Value{})
		} else {
			panic(" no  method " + methodName + " of type " + tp.Name())
		}

	}
}
