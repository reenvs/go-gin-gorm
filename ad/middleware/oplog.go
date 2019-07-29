package middleware

import (
	"encoding/json"
	"iptv/ad/model"
	"iptv/common/constant"
	"iptv/common/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func OperationLogHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		//c.Header("Server-Version", constant.Version)
		c.Next()
		var elapsed uint32 // millisecond
		elapsed = uint32(time.Since(start) / time.Millisecond)
		var opLog model.OperationLog
		var err string
		var table string
		var typ uint32

		tab, exists := c.Get(constant.ContextTableName)
		if exists && tab != nil {
			table = tab.(string)
		} else {
			// if table name is not set, then return
			return
		}

		opType, exists := c.Get(constant.ContextOperationType)
		if exists && tab != nil {
			typ = uint32(opType.(int))
		}

		e, exists := c.Get(constant.ContextError)
		if exists && e != nil {
			err = e.(string)
		}

		var adminId uint32
		ad, exists := c.Get(constant.ContextAdminId)
		if exists && ad != nil {
			adminId = ad.(uint32)
		}

		var adminName string
		adName, exists := c.Get(constant.ContextAdminName)
		if exists && adName != nil {
			adminName = adName.(string)
		}

		loginIp := c.ClientIP()
		requestUrl := c.Request.URL.RequestURI()

		requestBody, _ := c.Get(constant.ContextRequestBody)
		oldValue, _ := c.Get(constant.ContextOldValue)
		newValue, _ := c.Get(constant.ContextNewValue)

		opLog.Table = table
		opLog.Method = c.Request.Method
		opLog.RequestUrl = requestUrl
		opLog.Type = typ

		if requestBody != nil {
			b, err := json.Marshal(requestBody)
			if err != nil {
				logger.Error(err)
				return
			}
			opLog.RequestBody = string(b)
		}

		if oldValue != nil {
			b, err := json.Marshal(oldValue)
			if err != nil {
				logger.Error(err)
				return
			}
			opLog.OldValue = string(b)
		}

		if newValue != nil {
			b, err := json.Marshal(newValue)
			if err != nil {
				logger.Error(err)
				return
			}
			opLog.NewValue = string(b)
		}

		opLog.Error = err
		opLog.Operator = adminName
		opLog.OperatorId = adminId
		opLog.OperatorIp = loginIp
		opLog.ExecutionTime = elapsed

		db := c.MustGet(constant.ContextDb).(*gorm.DB)

		var admin *model.Admin
		if a, adminExists := c.Get(constant.ContextAdmin); adminExists == true {
			admin = a.(*model.Admin)
		}

		opLog.ServerIp = c.ClientIP()
		if admin != nil {
			opLog.Operator = admin.Username
			opLog.OperatorId = admin.Id
		}

		if err := db.Create(&opLog).Error; err != nil {
			logger.Error(err)
		}
	}
}
