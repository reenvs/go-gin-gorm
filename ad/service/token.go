package service

import (
	"fmt"
	"iptv/ad/model"
	"iptv/common/logger"
	"iptv/common/util"
	"time"

	"github.com/jinzhu/gorm"
)

// this method allows to get the admin by token, if token is valid, then continue it automatically.
func GetAdminFromToken(tokenStr string, db *gorm.DB) (*model.Token, *model.Admin) {
	if tokenStr == "" {
		return nil, nil
	}

	var token model.Token
	if err := db.Where("token = ? AND expired_at > ?", tokenStr, time.Now()).First(&token).Error; err != nil {
		logger.Error(err)
		return nil, nil
	}

	// extend expired_at automatically 20 minutes
	expiredAt := time.Now().Add(time.Minute * 20)
	token.ExpiredAt = &expiredAt
	if err := db.Save(&token).Error; err != nil {
		logger.Error(err)
	}

	var admin model.Admin
	admin.Id = token.AdminId
	if err := db.First(&admin).Error; err != nil {
		logger.Error(err)
		return &token, nil
	}
	return &token, &admin
}

func GenerateTokenForAdmin(admin *model.Admin, ip string, db *gorm.DB) string {
	var token model.Token
	token.AdminId = admin.Id
	token.Token = util.Md5Hash(fmt.Sprintf("%d%s%s", admin.Id, ip, time.Now()))
	expiredAt := time.Now().Add(time.Minute * 20)
	token.ExpiredAt = &expiredAt
	if err := db.Create(&token).Error; err != nil {
		logger.Error(err)
		return ""
	}
	return token.Token
}
