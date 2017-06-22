package middleware

import (
	"log"
	"net/http"
	"time"

	"io"

	"os"

	"github.com/julienschmidt/httprouter"
)

var file io.Writer

// SetLoggger 设置日志
func SetLoggger(writer io.Writer) {
	file = writer
}

// LoggerMiddleware middle
func LoggerMiddleware(h httprouter.Handle) httprouter.Handle {
	const logTmpl = "%+v | %3d |  %v | %s | %s | %s \n"
	if file == nil {
		file = os.Stdout
	}
	l := log.New(file, "[Golang]", 0)
	return func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
		start := time.Now()
		h(rw, r, params)
		status := 0
		nrw, ok := rw.(*WrapedResponseWriter)
		if ok {
			status = nrw.Status
		}
		l.Printf(logTmpl, start, status, time.Since(start), r.Host, r.Method, r.URL.Path)
	}
}
