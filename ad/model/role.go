package model

import (
	"encoding/json"
	"iptv/common/logger"
	"time"

	"github.com/jinzhu/gorm"
)

type Role struct {
	Id          uint32      `gorm:"primary_key" json:"id"`
	Name        string      `gorm:"size:64" json:"name"`
	Description string      `gorm:"size:255" json:"description"`
	AdminCount  uint32      `gorm:"-" json:"admin_count"`
	Scope       interface{} `gorm:"type:longtext" json:"scope"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func (Role) TableName() string {
	return "role"
}

type RoleScope struct {
	AppIds []uint32 `json:"app_ids"`
}

func (r *Role) Transcope() *Role {
	var scope RoleScope
	if len(r.Scope.([]byte)) > 0 {
		json.Unmarshal(r.Scope.([]byte), &scope)
	}
	r.Scope = scope
	return r
}

func initRole(db *gorm.DB) error {
	var err error
	if db.HasTable(&Role{}) {
		err = db.AutoMigrate(&Role{}).Error
	} else {
		err = db.CreateTable(&Role{}).Error
	}

	role := Role{}
	role.Name = "超级管理员"

	exists := false
	if err = db.Model(Role{}).Where("name=?", role.Name).First(&role).Error; err == nil {
		exists = true
	}

	role.Description = "超级管理员"
	emptyScope, _ := json.Marshal(RoleScope{})
	role.Scope = string(emptyScope)
	if !exists {
		if err = db.Model(Role{}).Create(&role).Error; err != nil {
			logger.Error(err)
			return err
		}
	}
	return err
}

func dropRole(db *gorm.DB) {
	db.DropTableIfExists(&Role{})
}
