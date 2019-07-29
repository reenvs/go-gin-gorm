package main

import (
	"github.com/gin-gonic/gin"
	"iptv/ad/controller/cad"
)

func initAdminChkScoapApi(adminScopeApi gin.IRoutes) {
	// 广告商
	adminScopeApi.GET("dsp/list", cad.DspListHandler)
	adminScopeApi.GET("dsp/detail", cad.DspDetailHandler)
	adminScopeApi.POST("dsp/save", cad.DspSaveHandler)
	adminScopeApi.POST("dsp/delete", cad.DspDeleteHandler)
	// 广告订单
	adminScopeApi.POST("order/save", cad.AdOrderSaveHandler)
	adminScopeApi.GET("order/list", cad.AdOrderListHandler)
	adminScopeApi.GET("order/detail", cad.AdOrderDetailHandler)
	adminScopeApi.POST("order/status", cad.AdOrderStatusHandler)
	adminScopeApi.POST("order/disable", cad.AdOrderDisableHandler)
	adminScopeApi.POST("order/delete", cad.AdOrderDeleteHandler)
	// 广告创意
	adminScopeApi.POST("advertising/code", cad.AdvertisingCodeHandler)
	adminScopeApi.POST("advertising/save", cad.AdvertisingSaveHandler)
	adminScopeApi.GET("advertising/list", cad.AdvertisingListHandler)
	adminScopeApi.GET("advertising/detail", cad.AdvertisingDetailHandler)
	adminScopeApi.POST("advertising/delete", cad.AdvertisingDeleteHandler)
}