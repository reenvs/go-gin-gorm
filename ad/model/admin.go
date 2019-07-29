package model

import (
	"iptv/common/logger"
	"iptv/common/util"
	"time"

	"github.com/jinzhu/gorm"
)

type Admin struct {
	Id        uint32        `gorm:"primary_key" json:"id"`
	Username  string        `gorm:"size:64;unique" json:"username"`
	Realname  string        `gorm:"size:64" json:"realname"`
	Mobile    NullString `gorm:"size:16;unique" json:"mobile"`
	Email     NullString `gorm:"size:64;unique" json:"email"`
	Password  string        `gorm:"size:255" json:"password"`
	Roles     []*Role    `gorm:"many2many:admin_role" json:"roles"`
	LoginAt   *time.Time    `json:"login_at"`
	Locked    bool          `json:"locked"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	WebUIType uint32        `gorm:"-" json:"-"`
}

func (Admin) TableName() string {
	return "admin"
}

func initAdmin(db *gorm.DB) error {
	var err error
	if db.HasTable(&Admin{}) {
		err = db.AutoMigrate(&Admin{}).Error
	} else {
		err = db.CreateTable(&Admin{}).Error
	}

	var admin Admin
	admin.Username = "superadmin"
	if err = db.Where("username=?", admin.Username).First(&admin).Error; err != nil {
		admin.Password = util.SHA512("stars@admin")
		db.Create(&admin)
	}

	var roles []*Role
	if err = db.Where("name='超级管理员'").Find(&roles).Error; err != nil {
		logger.Error(err)
		return err
	}

	if err = db.Model(&admin).Association("Roles").Replace(roles).Error; err != nil {
		logger.Error(err)
		return err
	}
	return err
}

func dropAdmin(db *gorm.DB) {
	db.DropTableIfExists(&Admin{})
}
