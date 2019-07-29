package processes

import (
	"github.com/elazarl/goproxy"
	"iptv/proxy/handlers"
	"net/http"
)

func init(){
	HandlerChan["fengshang"] = new(handlers.FengShang)
}

// 请求处理器链
// key为拦截请求的URL，value为处理器实现
var HandlerChan = make(map[string]Process)

// 操作接口
type Process interface { // 当前需要拦截的URL
	Handler(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) // 请求拦截处理函数
	// OperationRequestHeader(header *http.Header)                                       // 需要操作拦截请求的header
	// OperationResponseHeader(header *http.Header)                                      // 需要操作的响应体的header
	// OperationRequestBody(body *io.ReadCloser)                                         // 需要操作的请求体
	// OperationResponseBody(body *io.ReadCloser)                                        // 需要操作的响应体
}
