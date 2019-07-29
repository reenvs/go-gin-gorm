package main

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"iptv/proxy/processes"
	"log"
	"net/http"
	"strings"
)

func main() {
	// config:=flag.String("conf","./conf/config.yaml","配置文件路径")
	// flag.Parse()
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		fmt.Printf("当前请求URL：%s\n", req.URL)
		// TODO：后续增加广告商时这里应当过滤请求
		for host, proc := range processes.HandlerChan {
			// if strings.Contains(req.URL.Host, host) {
			// 	return proc.Handler(req *http.Request, ctx *goproxy.ProxyCtx)
			// }
			if strings.Contains(req.URL.Host, "haijd.tech") {
				fmt.Println(host)
				return proc.Handler(req, ctx)
			}
		}
		return req, nil
	})
	fmt.Println("服务器启动")
	proxy.Verbose = false
	log.Fatal(http.ListenAndServe(":8888", proxy))
}

func orPanic(err error) {
	if err != nil {
		panic(err)
	}
}
