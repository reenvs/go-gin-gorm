// +build debug

package middleware

import "github.com/gin-gonic/gin"

// 本地调试不校验授权
// only production environnement enables signature function
func SignatureVerifyHandler(isProductionEnv bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			c.Next()
		}()
	}
}
