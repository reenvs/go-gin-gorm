package model

import "github.com/jinzhu/gorm"

type Material struct {
	ID                  uint32 `gorm:"primary_key" json:"id"`
	Name                string `json:"name"`
	OrderID             string `json:"order_id"`
	Code                string `json:"code"`
	Type                uint32 `json:"type"`
	ActionType          uint32 `json:"action_type"`
	LandingURL          string `json:"landing_url"`
	PackageName         string `json:"package_name"`
	IsOrigin            uint32 `json:"is_origin"`
	AdOriginContentType uint32 `json:"ad_origin_content_type"`
	AdOriginContent     string `json:"ad_origin_content"`
	EndTime             string `json:"end_time"`
	Title               string `json:"title"`
	Src                 string `json:"src"`
	Duration            uint32 `json:"duration"`
}

const (
	CodeStatusReviewing = 1 // 审核中
	CodeStatusFailed    = 2 // 未通过
	CodeStatusReady     = 3 // 已通过
)

func (Material) TableName() string {
	return "material"
}

func initMaterial(db *gorm.DB) error {
	var err error
	if db.HasTable(&Material{}) {
		err = db.AutoMigrate(&Material{}).Error
	} else {
		err = db.CreateTable(&Material{}).Error
	}
	return err
}

func dropMaterial(db *gorm.DB) {
	db.DropTableIfExists(&Material{})
}
