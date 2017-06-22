package register

import (
	"reflect"
	"service/api"
	"service/card"
	"web/httpcontext"
	"web/middleware"

	"g"

	"github.com/julienschmidt/httprouter"
)

// Example example
type Example struct {
	httpcontext.HttpContext
}

// Get method
func (e *Example) Get() {
	e.WriteString("Hello")
}

// RegisterTypeHandle 注册类型handle
func RegisterTypeHandle(r *httprouter.Router) {
	if g.Configure.Mode == "local" {
		r.GET("/api/:"+httpcontext.InvokeMethodName, middleware.DefaultMiddle(MakeTypeHttpRouter(reflect.TypeOf(card.Card{}))))
		r.POST("/api/:"+httpcontext.InvokeMethodName, middleware.DefaultMiddle(MakeTypeHttpRouter(reflect.TypeOf(card.Card{}))))
	} else if g.Configure.Mode == "public" {
		r.GET("/api/:"+httpcontext.InvokeMethodName, middleware.DefaultMiddle(MakeTypeHttpRouter(reflect.TypeOf(api.Api{}))))
		r.POST("/api/:"+httpcontext.InvokeMethodName, middleware.DefaultMiddle(MakeTypeHttpRouter(reflect.TypeOf(api.Api{}))))
	}
}
