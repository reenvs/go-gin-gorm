package cad

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"iptv/ad/model"
	"iptv/common/constant"
	"iptv/common/logger"
	"iptv/common/util"
	"net/http"
)

/*
	POST /cad/advertising/code
	加载审核通过code码
*/
func AdvertisingCodeHandler(c *gin.Context) {
	type param struct {
		ID uint32 `json:"id" binding:"required"`
		//OrderID string `json:"order_id" binding:"required"`
		Code string `json:"code"  binding:"required"`
	}

	var ps []*param
	var err error
	if err = c.ShouldBind(&ps); err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var orderId string
	for _, p := range ps {
		var material model.Material
		if err = db.First(&material, p.ID).Update("code", p.Code).Error; err != nil {
			logger.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
			return
		}
		orderId = material.OrderID
	}

	// 修改订单为就绪状态
	if err = db.Model(model.Order{}).Where("order_id = ?", orderId).Update("status", model.AdStatusReady).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}

/*
	POST /cad/advertising/save
	创意保存
*/
func AdvertisingSaveHandler(c *gin.Context) {
	var ad model.Advertising
	//var ad,dbad model.Advertising
	var err error

	if err = c.ShouldBind(&ad); err != nil {
		logger.Error("Invalid ad review param", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	// 参数校验
	// 校验传入的时间，有效期5分钟
	//if time.Now().Unix()-int64(ad.TimeStamp) > 300 {
	//	err = fmt.Errorf("wrong sign")
	//	logger.Error(err)
	//	c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
	//	return
	//}

	/*
		签名规则：
		1. 拼接字符串：area_code + order_name + dsp_id + dsp_order_id + ad_position_id + time_stamp + token
		2. 对字符串进行sha1加密,得到40位字符串签名密钥
	*/

	//str := ad.AdId + ad.AdName + ad.OrderID + strconv.Itoa(int(ad.TimeStamp)) + service.GetToken(service.GetDspNameByOrderId(ad.OrderID, db), db)
	//if ad.Sign != service.GetSign(str) {
	//	err = fmt.Errorf("wrong sign")
	//	logger.Error(err)
	//	c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
	//	return
	//}

	// 目前只支持新建，不支持修改
	for _,material := range ad.Materials{
		if err = db.Create(&material).Error; err != nil {
			logger.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
			return
		}
	}

	if err = db.Create(&ad).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	// 修改创意订单状态为审核中
	if err = db.Model(model.Order{}).Where("order_id = ?", ad.OrderID).Update("status", model.AdStatusReviewing).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": ad})
}

/*
	GET /cad/advertising/list
	创意列表
*/
func AdvertisingListHandler(c *gin.Context) {
	type param struct {
		Offset  uint32 `form:"offset" binding:"exists"`
		Limit   uint32 `form:"limit" binding:"required"`
		OrderID string `form:"order_id"`
	}

	var p param
	var err error
	if err = c.ShouldBind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	dbExec := db.Order(constant.DefaultColumn + " " + constant.SortDesc + ", id desc")

	if p.OrderID != "" {
		dbExec = dbExec.Where("order_id = ?", p.OrderID)
	}

	var count uint32
	if err = dbExec.Model(model.Material{}).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	var ads []*model.Material
	if err = dbExec.Offset(p.Offset).Limit(p.Limit).Find(&ads).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "total": count, "data": ads})
}

/*
	GET /cad/advertising/status
	根据状态展示创意
*/
func AdvertisingStatusHandler(c *gin.Context) {
	type param struct {
		OrderId string `json:"order_id" binding:"required"`
		Status  uint32 `json:"status" binding:"required"`
	}

	var p param
	var err error
	if err = c.ShouldBind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	// 校验状态
	if p.Status != model.CodeStatusReviewing && p.Status != model.CodeStatusFailed && p.Status != model.CodeStatusReady {
		err = fmt.Errorf("invaild code status %v", p.Status)
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	dbExec := db
	if p.Status == model.CodeStatusReviewing {
		dbExec = dbExec.Where("code = ''")
	} else if p.Status == model.CodeStatusFailed {
		dbExec = dbExec.Where("code = '未通过'")
	} else {
		dbExec = dbExec.Where("code != '' AND code != '未通过'")
	}

	var ads []*model.Material
	if err = dbExec.Where("order_id = ?", p.OrderId).Find(&ads).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}

/*
	GET /cad/advertising/detail
	创意详情
*/
func AdvertisingDetailHandler(c *gin.Context) {
	type param struct {
		ID uint32 `form:"id"`
	}

	var p param
	var err error
	if err = c.ShouldBind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var material model.Material
	material.ID = p.ID
	if err = db.First(&material).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": material})
}

/*
	POST /cad/advertising/delete
	创意删除
*/
func AdvertisingDeleteHandler(c *gin.Context) {
	type param struct {
		ID []uint32 `json:"id" binding:"required"`
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
		c.Set(constant.ContextTableName, model.Advertising{}.TableName())
		c.Set(constant.ContextOperationType, constant.OperationTypeDelete)

		if err != nil {
			tx.Rollback()
			c.Set(constant.ContextError, err)
			return
		}

		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
	}()

	if err = tx.Where("id in (?)", p.ID).Delete(model.Material{}).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}
}
