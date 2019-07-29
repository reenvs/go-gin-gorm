package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Token struct {
	Id        uint32     `gorm:"primary_key" json:"-"`
	AdminId   uint32     `gorm:"size:64" json:"-"`
	Token     string     `gorm:"size:128;unique" json:"token"`
	ExpiredAt *time.Time `json:"-"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
}

func (Token) TableName() string {
	return "token"
}

func initToken(db *gorm.DB) error {
	var err error
	if db.HasTable(&Token{}) {
		err = db.AutoMigrate(&Token{}).Error
	} else {
		err = db.CreateTable(&Token{}).Error
	}
	return err
}

func dropToken(db *gorm.DB) {
	db.DropTableIfExists(&Token{})
}
