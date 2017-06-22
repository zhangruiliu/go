package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Middleware httprouter 中间件
type Middleware func(httprouter.Handle) httprouter.Handle

// Chain 中间件串
func Chain(outer Middleware, middlewares ...Middleware) Middleware {
	return func(h httprouter.Handle) httprouter.Handle {
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}
		return outer(h)
	}
}

// WrapHandler 包装http.Handle
func WrapHandler(handle http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		handle.ServeHTTP(w, r)
	}
}

// UnwrapHandler 解包成 http.Handler
func UnwrapHandler(handle httprouter.Handle) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handle(w, r, nil)
	})
}
