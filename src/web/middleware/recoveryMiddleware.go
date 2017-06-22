package middleware

import (
	"net/http"
	"os"

	"encoding/json"
	"runtime"

	"fmt"

	"github.com/julienschmidt/httprouter"
)

// RecoveryMiddleware 日志
func RecoveryMiddleware(h httprouter.Handle) httprouter.Handle {
	if file == nil {
		file = os.Stdout
	}
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		defer func() {
			if rec := recover(); rec != nil {
				var p [1024]byte
				n := runtime.Stack(p[:1023], false)
				p[n] = '\n'
				wr, ok := w.(*WrapedResponseWriter)
				fmt.Fprintf(file, "%v\n", rec)
				file.Write(p[:n+1])
				if ok {
					if r.Method == "GET" {
						if wr.Written {
							return
						}
						if wr.Status != 0 {
							return
						}
						wr.WriteHeader(500)
						wr.Write(p[:n+1])
					} else {
						if wr.Written {
							return
						}
						m := make(map[string]interface{})
						wr.WriteHeader(500)
						mm := make(map[string]interface{})
						mm["message"] = fmt.Sprintf("%v", r)
						m["error"] = mm
						buf, _ := json.Marshal(m)
						wr.Write(buf)
					}
				}

			}
		}()
		h(w, r, params)
	}
}
