package main

import (
	"flag"
	"fmt"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"iptv/ad/config"
	"iptv/ad/controller/ams"
	"iptv/ad/controller/cad"
	"iptv/ad/model"
	"iptv/common/constant"
	"iptv/common/logger"
	"iptv/common/middleware"
	"log"
)

func main() {
	// 加载配置文件
	configPath := flag.String("conf", "./ad/config/config.json", "Config file path")
	listenPort := flag.Int("port", 26880, "listen port")
	flag.Parse()

	err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("Config Failed!!!!", err)
		return
	}

	logger.SetLevel(config.GetLoggerLevel())

	// 连接数据库，做初始化
	amsdb, err := gorm.Open(config.GetDBName(), config.GetDBSourceAms())
	if err != nil {
		logger.Fatal("Open ams db Failed!!!!", err)
		return
	}
	if config.IsOrmLogEnabled() {
		amsdb.LogMode(true)
	} else {
		amsdb.LogMode(false)
	}

	model.InitAmsModel(amsdb)

	caddb, err := gorm.Open(config.GetDBName(), config.GetDBSourceCad())
	if err != nil {
		logger.Fatal("Open cad db Failed!!!!", err)
		return
	}
	if config.IsOrmLogEnabled() {
		caddb.LogMode(true)
	} else {
		caddb.LogMode(false)
	}

	model.InitCadModel(caddb)

	r := gin.New()

	if config.IsProductionEnv() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	if config.IsPprofEnabled() {
		ginpprof.Wrap(r)
	}
	if config.IsPrometheusEnabled() {
		prometheusMiddleware := middleware.NewPrometheus("gin")
		prometheusMiddleware.Use(r)
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	//r.Use(middleware.CIBNRecovery())

	r.Use(middleware.CorsAllowHandler)
	r.OPTIONS("*f", func(c *gin.Context) {})

	dbMiddleware := middleware.GetDbPrepareHandler(config.GetDBName(), config.GetDBSourceAms(), config.IsOrmLogEnabled(), constant.ContextDb)
	adminMiddleware := middleware.AdminVerifyHandler(config.GetBindAddr())
	logMiddleware := middleware.OperationLogHandler(config.GetBindAddr(), "11e87e97-4a1f-1dba-ac15-00163e08201c")
	//logMiddleware := admw.OperationLogHandler()

	// 以下是登陆验证部分
	api := r.Group("ams")

	api.Use(dbMiddleware, logMiddleware)
	{
		api.POST("/admin/login", ams.AdminLoginHandler)
		api.POST("/admin/ftp/login", ams.AdminFtpLoginHandler)
	}

	apiAuth := api.Use(adminMiddleware)
	{
		apiAuth.POST("/admin/logout", ams.AdminLogoutHandler)
		apiAuth.POST("/role/scope/add", ams.RoleScopeAdd)
		apiAuth.GET("/operationlog/list", ams.OperationLogListHandler)
		apiAuth.POST("/operationlog/create", ams.OperationLogCreateHandler)
	}

	apiPermission := apiAuth.Use(middleware.PermissionVerifyHandler)
	{
		apiPermission.GET("/role/list", ams.RoleListHandler)
		apiPermission.GET("/role/detail", ams.RoleDetailHandler)
		apiPermission.POST("/role/save", ams.RoleSaveHandler)
		apiPermission.POST("/role/delete", ams.RoleDeleteHandler)

		apiPermission.GET("/admin/list", ams.AdminListHandler)
		apiPermission.GET("/admin/detail", ams.AdminDetailHandler)
		apiPermission.POST("/admin/save", ams.AdminSaveHandler)
		apiPermission.POST("/admin/delete", ams.AdminDeleteHandler)
	}

	signApi := r.Group("ams")
	signApi.Use(dbMiddleware, middleware.ModuleSignatureVerifyHandler)
	{
		//signApi.POST("/accesskey/verify", controller.AccessKeyVerifyHandler)
		signApi.POST("/admin/verify", ams.AdminVerifyHandler)
	}

	// 以下是计算广告接口
	dbCreateMiddleware := middleware.GetDbPrepareHandler(config.GetDBName(),
		config.GetDBSourceCad(),
		config.IsOrmLogEnabled(),
	constant.ContextDb)
	caddbMiddleware := middleware.GetDbPrepareHandler(config.GetDBName(), config.GetDBSourceCad(), config.IsOrmLogEnabled(), constant.ContextDb)
	//dbMiddleware := middleware.GetDbHandler(db, constant.ContextDb)
	//amsBindAddr := config.GetBindAddr()
	//if err != nil {
	//	logger.Error(err)
	//}
	//logMiddleware := middleware.OperationLogHandler(amsBindAddr, amsAccessKey)
	//adminMiddleware := middleware.AdminVerifyHandler(amsBindAddr)

	adminScopeApi := r.Group("cad")
	adminScopeApi.Use(dbCreateMiddleware, caddbMiddleware, logMiddleware, adminMiddleware)
	adminScopeApi.Use(middleware.CorsAllowHandler) // todo scope验证
	//initAdminChkScoapApi(adminScopeApi)

	// 广告商
	adminScopeApi.GET("dsp/list", cad.DspListHandler)
	adminScopeApi.GET("dsp/detail", cad.DspDetailHandler)
	adminScopeApi.GET("dsp/type", cad.DspTypeHandler)
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
	adminScopeApi.GET("advertising/status", cad.AdvertisingStatusHandler)
	adminScopeApi.GET("advertising/detail", cad.AdvertisingDetailHandler)
	adminScopeApi.POST("advertising/delete", cad.AdvertisingDeleteHandler)

	r.Run(fmt.Sprintf(":%d", *listenPort))
}
