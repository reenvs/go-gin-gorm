package model

import (
	"github.com/jinzhu/gorm"
	"iptv/common/logger"
)

func InitAmsModel(db *gorm.DB) error {
	var err error

	//解决服务端配置没有配置字符编码情况下，使用程序自动建表时编码不是utf-8问题
	db = db.Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4")

	err = initModule(db)
	if err != nil {
		logger.Fatal("Init db role failed, ", err)
		return err
	}

	err = initRole(db)
	if err != nil {
		logger.Fatal("Init db role failed, ", err)
		return err
	}

	err = initAdmin(db)
	if err != nil {
		logger.Fatal("Init db admin failed, ", err)
		return err
	}

	err = initOperationLog(db)
	if err != nil {
		logger.Fatal("Init db operation_log failed, ", err)
		return err
	}

	err = initToken(db)
	if err != nil {
		logger.Fatal("Init db token failed, ", err)
		return err
	}

	return err
}

func ReBuildAmsModel(db *gorm.DB)  {
	dropOperationLog(db)
	dropToken(db)
	dropAdmin(db)
	dropRole(db)
}


func InitCadModel(db *gorm.DB) error {
	var err error

	//解决服务端配置没有配置字符编码情况下，使用程序自动建表时编码不是utf-8问题
	db = db.Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4")

	err = initCityCode(db)
	if err != nil {
		logger.Fatal("Init db advertising failed, ", err)
		return err
	}

	err = initMaterial(db)
	if err != nil {
		logger.Fatal("Init db material failed, ", err)
		return err
	}

	err = initAdvertising(db)
	if err != nil {
		logger.Fatal("Init db advertising failed, ", err)
		return err
	}

	err = initOrder(db)
	if err != nil {
		logger.Fatal("Init db order failed, ", err)
		return err
	}

	err = initDsp(db)
	if err != nil {
		logger.Fatal("Init db dsp failed, ", err)
		return err
	}
	return err
}


func ReBuildModel(db *gorm.DB)  {
	dropDsp(db)
	dropOrder(db)
	dropAdvertising(db)
	dropMaterial(db)
	dropCityCode(db)
}