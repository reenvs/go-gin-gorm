package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Dsp struct {
	Id        uint32    `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"size:64;unique" json:"dsp_name"`
	DspID     string    `gorm:"size:64;unique" json:"dspId"`    // DSP商家id,由我们计算生成
	DspType   string    `gorm:"size:64;unique" json:"dsp_type"` // 内部码
	Token     string    `gorm:"size:64;unique" json:"-"`        // token
	Orders    []*Order  `json:"orders"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Dsp) TableName() string {
	return "dsp"
}

func initDsp(db *gorm.DB) error {
	var err error
	if db.HasTable(&Dsp{}) {
		err = db.AutoMigrate(&Dsp{}).Error
	} else {
		err = db.CreateTable(&Dsp{}).Error
	}
	return err
}

func dropDsp(db *gorm.DB) {
	db.DropTableIfExists(&Dsp{})
}
