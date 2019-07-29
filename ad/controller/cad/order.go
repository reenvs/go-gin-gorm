package cad

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"iptv/ad/model"
	"iptv/ad/service"
	"iptv/common/constant"
	"iptv/common/logger"
	"iptv/common/util"
	"net/http"
	"strings"
	"time"
)

/*
	POST /cad/order/save
	生成广告订单
*/
func AdOrderSaveHandler(c *gin.Context) {
	var order, dbOrder model.Order
	var err error

	if err = c.ShouldBind(&order); err != nil {
		err = fmt.Errorf("invalid ad order param")
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	// 校验传入的时间，有效期5分钟
	//if time.Now().Unix() - int64(order.TimeStamp) > 300{
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
	//str := order.AreaCode + order.OrderName + order.DspId + order.DspOrderId + order.AdPositionId +
	//	strconv.Itoa(int(order.TimeStamp)) + service.GetToken(order.DspId,db)
	//if order.Sign != service.GetSign(str) {
	//	err = fmt.Errorf("wrong sign")
	//	logger.Error(err)
	//	c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
	//	return
	//}

	/*
		生成订单ID
		规则：时间(8位) + DspID(4位) + DSP序号(16进制)
	*/
	var dsp model.Dsp
	if err = db.Where("dsp_id = ?", order.DspId).First(&dsp).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	dbOrder.AreaCode = order.AreaCode // todo 根据城市码表转化
	dbOrder.OrderName = order.OrderName
	dbOrder.TotalCount = order.TotalCount
	dbOrder.DspId = order.DspId
	dbOrder.DspName = service.GetDspName(order.DspId,db)
	dbOrder.DspOrderId = order.DspOrderId
	dbOrder.AdPositionId = order.AdPositionId
	dbOrder.OrderDesc = order.OrderDesc
	dbOrder.BeginDate = order.BeginDate
	dbOrder.EndDate = order.EndDate
	dbOrder.BeginTime = order.BeginTime
	dbOrder.EndTime = order.EndTime
	dbOrder.DayControlType = order.DayControlType
	//dbOrder.DayCounts = order.DayCounts	todo
	dbOrder.UnitPrice = order.UnitPrice
	dbOrder.Contact = order.Contact
	dbOrder.Mobile = order.Mobile
	dbOrder.Email = order.Email
	dbOrder.Status = model.AdStatusSubmiting
	dbOrder.Disable = model.AdDisableOnShelf

	// 创建新订单
	if err = db.Create(&dbOrder).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	// 更新订单order_id
	hexID := fmt.Sprintf("%08x", dbOrder.ID)
	now := time.Now().Format("20060102")
	orderID := now + order.DspId + hexID

	if err = db.Model(&dbOrder).Update("order_id", orderID).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": dbOrder})
}

/*
	GET /cad/order/list
	订单列表
*/
func AdOrderListHandler(c *gin.Context) {
	type param struct {
		Offset    uint32 `form:"offset" binding:"exists"`
		Limit     uint32 `form:"limit" binding:"required"`
		OrderName string `form:"order_name"`
		Status    uint32 `form:"status"`
		Disable   uint32 `form:"disable"`
	}

	var p param
	var err error
	if err = c.ShouldBind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	// 参数校验
	if p.Status != 0 && p.Status != model.AdStatusSubmiting && p.Status != model.AdStatusReviewing && p.Status != model.AdStatusReady &&
		p.Status != model.AdStatusRunning && p.Status != model.AdStatusOver {
		err = fmt.Errorf("invaild status %v", p.Status)
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	if p.Disable != model.AdDisableAll && p.Disable != model.AdDisableOnShelf && p.Disable != model.AdDisableOffShelf {
		err = fmt.Errorf("invalid disable：%d", p.Disable)
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	dbExec := db.Order(constant.DefaultColumn + " " + constant.SortDesc + ", id desc")

	if p.OrderName != "" {
		dbExec = dbExec.Where("order_name like ?", "%"+strings.TrimSpace(p.OrderName)+"%")
	}

	if p.Status != 0 {
		dbExec = dbExec.Where("status = ?", p.Status)
	}

	if p.Disable != 0 {
		dbExec = dbExec.Where("disable = ?", p.Disable)
	}

	var count uint32
	if err = dbExec.Model(model.Order{}).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	var orders []*model.Order
	if err = dbExec.Offset(p.Offset).Limit(p.Limit).Find(&orders).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	for _,order := range orders{
		if err = db.Where("order_id = ?", order.OrderID).Find(&order.Ads).Error; err != nil {
			//if err = db.First(&order).Related(&order.Ads, "Ads").Error; err != nil {
			logger.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "total": count, "data": orders})
}

/*
	GET /cad/order/detail
	广告订单详情
*/
func AdOrderDetailHandler(c *gin.Context) {
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

	var order model.Order
	if err = db.First(&order, p.ID).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	if err = db.Where("order_id = ?", order.OrderID).Find(&order.Ads).Error; err != nil {
		//if err = db.First(&order).Related(&order.Ads, "Ads").Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": order})
}

/*
	POST /cad/order/status
	广告订单状态修改
*/
func AdOrderStatusHandler(c *gin.Context) {
	type param struct {
		ID     uint32 `json:"id"`
		Status uint32 `json:"status"`
	}

	var p param
	var err error
	if err = c.ShouldBind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	// 校验状态
	if p.Status != model.AdStatusSubmiting && p.Status != model.AdStatusReviewing && p.Status != model.AdStatusReady &&
		p.Status != model.AdStatusRunning && p.Status != model.AdStatusOver {
		err = fmt.Errorf("invaild status %v", p.Status)
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	if err = db.Model(model.Order{}).Where("id = ?", p.ID).Update("status", p.Status).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}

/*
	POST /cad/order/disable
	广告订单上下架
*/
func AdOrderDisableHandler(c *gin.Context) {
	type param struct {
		ID      uint32 `json:"id" binding:"required"`
		Disable uint32 `json:"disable" binding:"required"`
	}

	var p param
	var err error
	if err = c.ShouldBind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	if p.Disable != model.AdDisableAll && p.Disable != model.AdDisableOnShelf && p.Disable != model.AdDisableOffShelf {
		err = fmt.Errorf("invalid disable：%d", p.Disable)
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var order model.Order
	if err = db.First(&order, p.ID).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	if order.Disable == p.Disable {
		err = fmt.Errorf("cannot be modified to the same state：%d", p.Disable)
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	if err = db.Model(model.Order{}).Where("id = ?", p.ID).Update("disable", p.Disable).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}

/*
	POST /cad/order/delete
	广告订单删除
*/
func AdOrderDeleteHandler(c *gin.Context) {
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
		c.Set(constant.ContextTableName, model.Order{}.TableName())
		c.Set(constant.ContextOperationType, constant.OperationTypeDelete)

		if err != nil {
			tx.Rollback()
			c.Set(constant.ContextError, err)
			return
		}

		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
	}()

	if err = tx.Where("id in (?)", p.ID).Delete(model.Order{}).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}
}
