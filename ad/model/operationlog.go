package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type OperationLog struct {
	Id            uint32    `gorm:"primary_key" json:"id"`
	AppId         uint32    `json:"app_id"`
	Table         string    `gorm:"size:64" json:"table"`
	Type          uint32    `json:"type"`
	Method        string    `gorm:"size:16" json:"method"`
	RequestUrl    string    `gorm:"size:255" json:"request_url"`
	RequestBody   string    `gorm:"type:longtext" json:"request_body"`
	OldValue      string    `gorm:"type:longtext" json:"old_value"`
	NewValue      string    `gorm:"type:longtext" json:"new_value"`
	Error         string    `gorm:"size:255" json:"error"`
	Operator      string    `gorm:"size:32" json:"operator"`
	OperatorId    uint32    `json:"operator_id"`
	OperatorIp    string    `gorm:"size:64" json:"operator_ip"`
	ServerIp      string    `gorm:"size:64" json:"server_ip"`
	ExecutionTime uint32    `json:"execution_time"` // millisecond
	Status        uint32    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

func (OperationLog) TableName() string {
	return "operation_log"
}

const (
	OperationTypeUnknown = 0 //未知类型
	OperationTypeCreate  = 1 //创建
	OperationTypeUpdate  = 2 //更新
	OperationTypeList    = 3 //列表查询
	OperationTypeDetail  = 4 //详情
	OperationTypeDelete  = 5 //删除
	OperationTypeLogin   = 6 //登录
	OperationTypeLogout  = 7 //登出
)

func initOperationLog(db *gorm.DB) error {
	var err error

	if db.HasTable(&OperationLog{}) {
		err = db.AutoMigrate(&OperationLog{}).Error
	} else {
		err = db.CreateTable(&OperationLog{}).Error
	}
	return err
}

func dropOperationLog(db *gorm.DB) {
	db.DropTableIfExists(&OperationLog{})
}
