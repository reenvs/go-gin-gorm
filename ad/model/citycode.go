package model

import "github.com/jinzhu/gorm"

type CityCode struct {
	ID    uint32 `json:"id"`
	Code  string `json:"code"`
	City  string `json:"city"`
	Scope string `json:"scope"`
}

func (CityCode) TableName() string {
	return "city_code"
}

func initCityCode(db *gorm.DB) error {
	var err error
	if db.HasTable(&CityCode{}) {
		err = db.AutoMigrate(&CityCode{}).Error
	} else {
		err = db.CreateTable(&CityCode{}).Error
	}
	return err
}

func dropCityCode(db *gorm.DB) {
	db.DropTableIfExists(&CityCode{})
}