package httpcontext

import (
	"g"
	"io/ioutil"
	"net/http"
	"uuid"

	"bytes"

	"encoding/json"

	"sync"

	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

// HttpContext 表示一个请求上下文
type HttpContext struct {
	W http.ResponseWriter
	R *http.Request
	// RouterData 表示路由参数
	RouterData httprouter.Params
}

// Init 表示初始化请求上下文
func (ctx *HttpContext) Init(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	ctx.W = w
	ctx.R = r
	ctx.RouterData = params
}

// WriteString writestring
func (ctx *HttpContext) WriteString(str string) (int, error) {
	return ctx.W.Write([]byte(str))
}

// Write Write
func (ctx *HttpContext) Write(buf []byte) (int, error) {
	return ctx.W.Write(buf)
}

// CheckWrite CheckWrite
func (ctx *HttpContext) CheckWrite(err error) {
	if err == nil {
		return
	}
	ret := make(map[string]interface{})
	m := make(map[string]interface{})
	m["message"] = err.Error()
	ret["error"] = m
	ret["status"] = "failure"
	buf, _ := json.Marshal(ret)
	ctx.Write(buf)
	panic(err)
}

// JsonWrite JsonWrite
func (ctx *HttpContext) JsonWrite(data interface{}) {
	ret := make(map[string]interface{})
	ret["status"] = "success"
	ret["result"] = data
	buf, _ := json.Marshal(ret)
	ctx.Write(buf)
}
func (ctx *HttpContext) JsonWritefailure(data interface{}) {
	ret := make(map[string]interface{})
	ret["status"] = "failure"
	ret["result"] = data
	buf, _ := json.Marshal(ret)
	ctx.Write(buf)
}

// Check error check
func (ctx *HttpContext) Check(err error) {
	if err == nil {
		return
	}
	panic(err)
}

// UnJson Json反序列化
func (ctx *HttpContext) UnJson(data interface{}) error {
	getBytes, err := ioutil.ReadAll(ctx.R.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(getBytes, data)
	if err != nil {
		return err
	}

	ctx.R.Body = ioutil.NopCloser(bytes.NewBuffer(getBytes))
	// if body, ok := ctx.R.Body.(*poolBytes); ok {
	// 	body.Reset()
	// 	_, err = body.Write(getBytes)
	// 	if err != nil {
	// 		return err
	// 	}
	// } else {
	// 	if len(a) == 0 {
	// 		a, _ = ioutil.ReadFile(`C:\Users\shenf\Desktop\ask2v3.3.zip`)
	// 	}
	// 	ctx.R.Body, _, _ = newPoolBytes([]byte(a))

	// 	//	ctx.R.Body, _, _ = newPoolBytes(getBytes)
	// }

	return nil
}

var pool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(nil)
	},
}

func (ctx *HttpContext) DB() *gorm.DB {
	return g.DB0
}

func (ctx *HttpContext) Uuid() string {
	return uuid.New()
}

/*
测试了下，内存没有释放
http body 的 close 没有调用 可能是bug
type poolBytes struct {
	*bytes.Buffer
}

func newPoolBytes(by []byte) (ret *poolBytes, n int, err error) {
	ret = &poolBytes{
		pool.Get().(*bytes.Buffer),
	}
	n, err = ret.Buffer.Write(by)
	return
}

func (pb *poolBytes) Read(p []byte) (n int, err error) {
	return pb.Buffer.Read(p)
}

func (pb *poolBytes) Close() error {
	pb.Reset()
	pool.Put(pb.Buffer)
	return nil
}

var pool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(nil)
	},
}
*/
