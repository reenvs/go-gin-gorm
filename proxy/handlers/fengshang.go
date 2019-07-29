package handlers

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"github.com/valyala/fasthttp"
	"net/http"
)

// 风尚广告代理handler
type FengShang struct {
}

func (fs FengShang) Handler(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	url := `http://59.110.10.59:8081/api/1.0/shaanxi_iptv`

	request := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		// 用完需要释放资源
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(request)
	}()

	// 默认是application/x-www-form-urlencoded
	request.Header.SetContentType("application/json")
	request.Header.SetMethod("POST")

	request.SetRequestURI(url)

	requestBody := []byte(`{
    "deviceInfo": {
        "deviceId": "58749C242C7F97EFA7542B4EB9652E24",
        "ip": "101.255.12.13",
        "brand": "letv",
        "model": "letv2.0",
        "is4K": false,
        "osVersion": "android6.0",
        "screenWidth": 3840,
        "screenHeight": 2160,
        "carrier": "CMCC",
        "deviceMac": "00:9e:c8:29:f0:73",
        "ssid": "wirelessnet",
        "bssid": "00:9e:c8:29:f0:73",
        "city": "1156441500"
    },
    "userInfo": {
        "uid": "186000023334",
        "profile": {
            "gender": "male",
            "age": 28,
            "childrenNum": 0
        },
        "tags": [
            "宝马",
            "肯德基"
        ]
    },
    "appInfo": {
        "platform": "medianame",
        "platformType": "3",
        "version": "2.0"
    },
    "impRequests": [
        {
            "impId": "5d41402abc4b2a76b9719d911017c592",
            "plotId": "3",
            "dealId": "200009",
            "plotType": "front",
            "path": "点播/电影/内地",
            "channel": "电影",
            "tags": [
                "内地",
                "抗日"
            ],
            "title": "铁血战狼",
            "function": "点播",
            "width": 1280,
            "height": 720,
            "minDuration": 15,
            "maxDuration": 30
        }
    ],
    "context": {
        "ts": 1472607351
    }
}`)
	request.SetBody(requestBody)

	if err := fasthttp.Do(request, resp); err != nil {
		fmt.Println("请求失败:", err.Error())
		return nil,nil
	}

	b := resp.Body()

	fmt.Println("result:\r\n", string(b))
	return nil, nil
}

