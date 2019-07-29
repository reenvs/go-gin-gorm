package middleware

import (
	"fmt"
	"iptv/common/constant"
	"iptv/common/logger"
	"iptv/common/util"
	"net/http"
	"net/http/httputil"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

func CIBNRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				logger.Errorf("[Recovery] panic recovered:\n%s\n%s\n", string(httprequest), err)
				logger.Error(string(debug.Stack()))

				c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(fmt.Errorf("%v", err)))
			}
		}()
		c.Next()
	}
}

func CIBNAPiVersion() gin.HandlerFunc {
	return func(c *gin.Context) {
		uriPath := c.Request.RequestURI
		if len(uriPath) == 0 {
			uriPath = c.Request.URL.Path
		}
		ap := strings.Split(uriPath, "/")

		c.Set(constant.ContextApiVersion, ap[2])

		c.Next()
	}
}
