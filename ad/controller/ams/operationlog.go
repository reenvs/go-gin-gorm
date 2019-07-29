package ams

import (
	"errors"
	"iptv/ad/model"
	"iptv/common/constant"
	"iptv/common/logger"
	"iptv/common/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

/*
	POST /ams/operationlog/list
	获取管理员操作日志
*/
func OperationLogListHandler(c *gin.Context) {
	type param struct {
		Offset  int    `form:"offset" binding:"exists"`
		Limit   int    `form:"limit" binding:"required"`
		Order   string `form:"order"`
		Sort    string `form:"sort"`
		Module  string `form:"module"`
		AdminId uint32 `form:"admin_id"`
	}

	var p param
	var err error
	if err = c.Bind(&p); err != nil {
		logger.Error("Invalid OperationLogListHandler param ", err)
		return
	}

	if p.Offset < 0 ||
		p.Limit <= 0 ||
		p.Limit > constant.MaxPageSize {
		err = errors.New("bad request")
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	// set id as default order
	if p.Order == "" {
		p.Order = constant.DefaultColumn
	}

	sort := constant.SortDesc
	if p.Sort == constant.SortAsc {
		sort = constant.SortAsc
	}

	var logs []*model.OperationLog
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	dbExec := db

	if p.Module != "" {
		dbExec = dbExec.Where("module = ?", p.Module)
	}

	if p.AdminId != 0 {
		dbExec = dbExec.Where("operator_id = ?", p.AdminId)
	}

	dbCount := db
	if err = dbExec.Offset(p.Offset).Limit(p.Limit).Order(p.Order + " " + sort).Find(&logs).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	var count uint32
	if err = dbCount.Model(&model.OperationLog{}).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": logs})
}

/*
	POST /ams/operationlog/create
	创建操作日志
*/
func OperationLogCreateHandler(c *gin.Context) {
	var opLog model.OperationLog
	var err error
	if err = c.Bind(&opLog); err != nil {
		logger.Debug("Invalid OperationLogSaveHandler param ", err)
		return
	}

	// collect all useful information and create a operation_log.
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

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": &opLog})
}
