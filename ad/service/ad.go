package service

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/jinzhu/gorm"
	"iptv/ad/model"
)

func GetDspID(name, now string, db *gorm.DB) string {
	str := name + now
	hash := sha1.New()
	hash.Write([]byte(str))
	bytes := hash.Sum(nil)
	sha1Str := hex.EncodeToString(bytes)
	dspID := sha1Str[:4]

	// 验证dspID是否已存在，存在则重新生成
	var order model.Order
	if err := db.Where("dsp_id = ?", dspID).First(&order).Error; err == nil {
		now += "1"
		dspID = GetDspID(name, now, db)
	}

	return dspID
}

func GetDspNameByOrderId(order_id string,db *gorm.DB) string {
	var order model.Order
	if err := db.Where("order_id = ?",order_id).First(&order).Error;err!=nil{
		return ""
	}
	return order.DspId
}

func GetDspName(dsp_id string,db *gorm.DB) string {
	var dsp model.Dsp
	if err := db.Where("dsp_id = ?",dsp_id).First(&dsp).Error;err!=nil{
		return ""
	}
	return dsp.Name
}

func GetToken(dsp_id string,db *gorm.DB) string {
	var dsp model.Dsp
	if err := db.Where("dsp_id = ?",dsp_id).First(&dsp).Error;err!=nil{
		return ""
	}
	return dsp.Token
}

func GetSign(str string) string {
	hash := sha1.New()
	hash.Write([]byte(str))
	bytes := hash.Sum(nil)
	return hex.EncodeToString(bytes)
}