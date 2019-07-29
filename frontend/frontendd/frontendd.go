package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"iptv/common/middleware"
	"iptv/frontend/config"
	"log"
	"net/http"
)

func main() {
	// 加载配置文件
	configPath := flag.String("conf", "./frontend/config/config.json", "Config file path")
	listenPort := flag.Int("port", 26890, "listen port")
	flag.Parse()

	err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("Config Failed!!!!", err)
		return
	}

	r := gin.New()
	gin.SetMode(gin.DebugMode)

	r.Use(gin.Logger())
	//r.Use(gin.Recovery())
	r.Use(middleware.CIBNRecovery())

	r.LoadHTMLFiles("./frontend/html/index.html")
	//r.LoadHTMLFiles(config.GetImsRoot()+"html/index.html")

	r.NoRoute(func(c *gin.Context) {
		title := config.GetTitle()

		c.HTML(http.StatusOK, "index.html",
			gin.H{
				"title":   title,
				"product": config.GetProduct(),
			})
	})

	r.Static("/static", config.GetImsRoot()+"html/static")
	r.Run(fmt.Sprintf(":%d", *listenPort))
}
