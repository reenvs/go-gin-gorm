package model

import "github.com/jinzhu/gorm"

type Advertising struct {
	ID        uint32      `gorm:"primary_key" json:"id"`
	Sign      string      `gorm:"-" json:"-"`
	TimeStamp uint32      `gorm:"-" json:"time_stamp"` // 时间戳
	AdId      string      `json:"ad_id"`
	AdName    string      `json:"ad_name"`
	OrderID   string      `json:"order_id"`
	Materials []*Material `json:"materials"`
}

func (Advertising) TableName() string {
	return "advertising"
}

func initAdvertising(db *gorm.DB) error {
	var err error
	if db.HasTable(&Advertising{}) {
		err = db.AutoMigrate(&Advertising{}).Error
	} else {
		err = db.CreateTable(&Advertising{}).Error
	}
	return err
}

func dropAdvertising(db *gorm.DB) {
	db.DropTableIfExists(&Advertising{})
}
