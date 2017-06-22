package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// WrapedResponseWriter 包装的ResponseWriter
type WrapedResponseWriter struct {
	http.ResponseWriter
	Written bool
	Status  int
}

// WriteHeader Header
func (r *WrapedResponseWriter) WriteHeader(code int) {
	if r.Status == 0 {
		r.Written = true
		r.Status = code
		r.ResponseWriter.WriteHeader(code)
	}
}

// Write written
func (r *WrapedResponseWriter) Write(b []byte) (int, error) {
	if !r.Written {
		r.Written = true
		r.WriteHeader(http.StatusOK)
	}
	return r.ResponseWriter.Write(b)
}

// WrapResponseMiddleware 包装 http.ResponseWriter
func WrapResponseMiddleware(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ww := &WrapedResponseWriter{
			w,
			false,
			0,
		}
		h(ww, r, params)
	}
}
