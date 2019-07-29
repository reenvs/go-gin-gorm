package ams

import (
	"errors"
	"iptv/ad/model"
	"iptv/ad/service"
	"iptv/common/logger"
	"runtime/debug"

	"github.com/jinzhu/gorm"
)

type data struct {
	*model.Admin
	Scope string `json:"scope"`
}

func GetAdminFromToken(token string, db *gorm.DB) (*data, error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("panic: ", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	var err error
	_, admin := service.GetAdminFromToken(token, db)
	//logger.Debugf("Admin info:%v", *admin)
	if admin == nil {
		err = errors.New("invalid admin token")
		logger.Error("GetAdminFromToken ", err)
		return nil, err
	}

	if err := db.Model(admin).Related(&admin.Roles, "Roles").Error; err != nil {
		logger.Error("GetAdminFromToken ", err)
		return nil, err
	}

	scope, err := service.MergeAllRoleScope(admin, db)
	if err != nil {
		logger.Error("GetAdminFromToken MergeAllRoleScope ", err)
		return nil, err
	}

	for _, role := range admin.Roles {
		role.Transcope()
	}
	return &data{admin, scope}, nil
}
