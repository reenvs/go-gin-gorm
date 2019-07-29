package cad

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"iptv/ad/model"
	"iptv/ad/service"
	"iptv/common/constant"
	"iptv/common/logger"
	"iptv/common/util"
	"net/http"
	"strconv"
	"time"
)

/*
	GET /cad/dsp/list
	广告商列表
*/
func DspListHandler(c *gin.Context) {
	type param struct {
		Offset uint32 `form:"offset" binding:"exists"`
		Limit  uint32 `form:"limit" binding:"required"`
	}

	var p param
	var err error
	if err = c.ShouldBind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var dsps []*model.Dsp
	var count uint32

	if err = db.Model(model.Dsp{}).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	if err = db.Offset(p.Offset).Limit(p.Limit).Find(&dsps).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "total": count, "data": dsps})
}

/*
	GET /cad/dsp/detail
	广告商详情
*/
func DspDetailHandler(c *gin.Context) {
	type param struct {
		ID uint32 `form:"id"  binding:"required"`
	}

	var p param
	var err error
	if err = c.ShouldBind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var dsp model.Dsp
	if err = db.First(&dsp, p.ID).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	if err = db.Where("dsp_id = ?", dsp.DspID).Find(&dsp.Orders).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": dsp})
}

/*
	GET /cad/dsp/type
	广告商分类列表
*/
func DspTypeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": constant.Success})
}

/*
	POST /cad/dsp/save
	广告商保存
*/
func DspSaveHandler(c *gin.Context) {
	type param struct {
		ID      uint32 `json:"id"`
		Name    string `json:"name"`
		DspType string `json:"dsp_type"`
	}

	var p param
	var err error
	if err = c.ShouldBind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var dsp model.Dsp
	// 计算广告商编号
	now := time.Now().Unix()
	nowStr := strconv.Itoa(int(now))
	dspId := service.GetDspID(p.Name, nowStr, db)
	dsp.Name = p.Name
	dsp.DspID = dspId
	dsp.DspType = p.DspType

	if p.ID != 0 {
		dsp.Id = p.ID
		if err = db.Save(&dsp).Error; err != nil {
			logger.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"err_code": constant.Success,"data":dsp})
		return
	}

	// 新建
	str := p.Name + dspId + strconv.Itoa(int(time.Now().UnixNano()))
	dsp.Token = service.GetSign(str)

	if err = db.Create(&dsp).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": dsp})
}

/*
	POST /cad/dsp/delete
	广告商删除
*/
func DspDeleteHandler(c *gin.Context) {
	type param struct {
		ID []uint32 `json:"id"`
	}

	var p param
	var err error
	if err = c.ShouldBind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	tx := db.Begin()
	defer func() {
		c.Set(constant.ContextRequestBody, &p)
		c.Set(constant.ContextTableName, model.Dsp{}.TableName())
		c.Set(constant.ContextOperationType, constant.OperationTypeDelete)

		if err != nil {
			tx.Rollback()
			c.Set(constant.ContextError, err)
			return
		}

		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
	}()

	if err = tx.Where("id in (?)", p.ID).Delete(model.Dsp{}).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}
}
