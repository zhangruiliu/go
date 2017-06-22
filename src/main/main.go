package main

import (
	"fmt"
	"g"
	"jsonconf"
	"net/http"
	"os"
	"path/filepath"
	"register"
	"strings"
	"time"
	"web/middleware"

	"github.com/julienschmidt/httprouter"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	os.Setenv("NLS_LANG", "AMERICAN_AMERICA.AL32UTF8")

	// 配置初始化
	err := jsonconf.ConfFromFile("conf.json", &g.Configure)
	check(err)

	err = os.MkdirAll(g.Configure.LogPath, os.ModePerm)
	check(err)
	// 日志文件
	f, err := os.Create(filepath.Join(g.Configure.LogPath, strings.Replace(time.Now().Format(time.RFC3339Nano), ":", "-", -1)))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	g.Writer = f
	fmt.Fprintf(g.Writer, "%+v\n", g.Configure)

	err = g.Before()
	if err != nil {
		fmt.Fprintf(g.Writer, "%+v\n", err)
		return
	}
	defer func() {
		g.Clean()
	}()

	// 日志
	middleware.SetLoggger(g.Writer)

	err = g.LoadSqlTmpl()
	check(err)

	// 静态服务器
	file := http.FileServer(http.Dir(g.Configure.Static))

	afile := middleware.DefaultMiddle(middleware.WrapHandler(file))
	http.Handle("/", middleware.UnwrapHandler(afile))

	// 路由
	router := httprouter.New()
	register.RegisterTypeHandle(router)
	http.Handle("/api/", router)

	fmt.Fprintf(g.Writer, "%s\n", "Listening on "+g.Configure.Port+"....")

	err = http.ListenAndServe(":"+g.Configure.Port, nil)
	fmt.Println("1111")
	if err != nil {
		fmt.Fprintf(g.Writer, "%+v", err)
	}
}
