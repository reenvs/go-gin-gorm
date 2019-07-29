package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Module struct {
	Id          uint32    `gorm:"primary_key" json:"id"`
	ParentId    uint32    `json:"parent_id"`
	Public      bool      `json:"public"`
	Sort        uint32    `json:"sort"`
	Url         string    `gorm:"unique;size:128" json:"url"`
	Name        string    `gorm:"size:32" json:"name"`
	Description string    `gorm:"size:64" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Module) TableName() string {
	return "module"

}

//
// 2018-10-09 XGF
// 模块再此定义注册了，才能使用access key 远程访问
// 或者通过 /ams/accesskey/save 单独创建 新增模块的 access key
//

func initModule(db *gorm.DB) error {
	var err error
	if db.HasTable(&Module{}) {
		err = db.AutoMigrate(&Module{}).Error
	} else {
		err = db.CreateTable(&Module{}).Error
	}

	if err != nil {
		return err
	}

	type moduleImport struct {
		Url         string          `json:"url:"`
		Name        string          `json:"name"`
		Public      bool            `json:"public"`
		Description string          `json:"description"`
		Children    []*moduleImport `json:"children"`
	}

	var modules = []*moduleImport{
		&moduleImport{
			Url:         "/ams/admin/general",
			Name:        "ams用户操作",
			Description: "ams用户操作",
			Public:      true,
			Children: []*moduleImport{
				&moduleImport{Url: "/ams/admin/login",
					Name:        "ams管理员登录",
					Description: "ams管理员登录",
					Public:      true},
				&moduleImport{Url: "/ams/admin/logout",
					Name:        "ams管理员退出登录",
					Description: "ams管理员退出登录",
					Public:      true},
			},
		}, &moduleImport{
			Url:         "/ams/admin",
			Name:        "ams权限",
			Description: "ams权限",
			Public:      true,
			Children: []*moduleImport{
				&moduleImport{Url: "/ams/admin/checkmodule",
					Name:        "检查管理员对当前模块是否有授权",
					Description: "检查管理员对当前模块是否有授权",
					Public:      true},
				&moduleImport{Url: "/ams/admin/page",
					Name:        "ams获取有权限页面",
					Description: "ams获取有权限页面",
					Public:      true},
				&moduleImport{Url: "/ams/admin/verify",
					Name:        "ams管理员Token验证",
					Description: "ams管理员Token验证",
					Public:      true},
				&moduleImport{Url: "/ams/accesskey/verify",
					Name:        "ams管理员AccessKey验证",
					Description: "ams管理员AccessKey验证",
					Public:      true},
			},
		}, &moduleImport{
			Url:         "/ams/operationlog",
			Name:        "ams日志",
			Description: "ams日志",

			Children: []*moduleImport{
				&moduleImport{Url: "/ams/operationlog/list",
					Name:        "获取管理员操作日志",
					Description: "获取管理员操作日志"},
				&moduleImport{Url: "/ams/operationlog/create",
					Name:        "创建操作日志",
					Description: "创建操作日志",
					Public:      true},
			},
		}, &moduleImport{
			Url:         "/ams/module",
			Name:        "ams模块管理",
			Description: "ams模块管理",

			Children: []*moduleImport{
				&moduleImport{Url: "/ams/module/list",
					Name:        "获取模块列表",
					Description: "获取模块列表"},
				&moduleImport{Url: "/ams/module/save",
					Name:        "保存、发布模块",
					Description: "保存、发布模块"},
				&moduleImport{Url: "/ams/module/delete",
					Name:        "删除模块",
					Description: "删除模块"},
			},
		}, &moduleImport{
			Url:         "/ams/role",
			Name:        "ams角色管理",
			Description: "ams角色管理",

			Children: []*moduleImport{
				&moduleImport{Url: "/ams/role/list",
					Name:        "获取角色列表",
					Description: "获取角色列表"},
				&moduleImport{Url: "/ams/role/save",
					Name:        "保存、发布角色",
					Description: "保存、发布角色"},
				&moduleImport{Url: "/ams/role/delete",
					Name:        "删除角色",
					Description: "删除角色"},
				&moduleImport{Url: "/ams/role/detail",
					Name:        "角色详情",
					Description: "角色详情"},
				&moduleImport{Url: "/ams/role/scope/add",
					Name:        "ams管理员添加角色权限",
					Description: "ams管理员添加角色权限",
					Public:      true},
			},
		}, &moduleImport{
			Url:         "/ams/admin/manage",
			Name:        "ams管理员管理",
			Description: "ams管理员管理",

			Children: []*moduleImport{
				&moduleImport{Url: "/ams/admin/list",
					Name:        "获取管理员列表",
					Description: "获取管理员列表"},
				&moduleImport{Url: "/ams/admin/save",
					Name:        "保存、发布管理员",
					Description: "保存、发布管理员"},
				&moduleImport{Url: "/ams/admin/delete",
					Name:        "删除管理员",
					Description: "删除管理员"},
				&moduleImport{Url: "/ams/admin/detail",
					Name:        "管理员详情",
					Description: "管理员详情"},
			},
		}, &moduleImport{
			Url:         "/ams/accesskey",
			Name:        "ams Token管理",
			Description: "ams Token管理",

			Children: []*moduleImport{
				&moduleImport{Url: "/ams/accesskey/list",
					Name:        "获取token列表",
					Description: "获取token列表"},
				&moduleImport{Url: "/ams/accesskey/save",
					Name:        "保存、发布token",
					Description: "保存、发布token"},
				&moduleImport{Url: "/ams/accesskey/delete",
					Name:        "删除token",
					Description: "删除token"},
				&moduleImport{Url: "/ams/accesskey/detail",
					Name:        "获取token详情",
					Description: "获取token详情"},
			},
		},
	}

	tx := db.Begin()
	for _, module := range modules {
		var dbModule Module
		tx.Model(Module{}).Where("url=?", module.Url).First(&dbModule)
		dbModule.Name = module.Name
		dbModule.Description = module.Description
		dbModule.Url = module.Url
		dbModule.Public = module.Public
		tx.Save(&dbModule)

		for _, subModule := range module.Children {
			var subDbModule Module
			tx.Model(Module{}).Where("url=?", subModule.Url).First(&subDbModule)

			subDbModule.ParentId = dbModule.Id
			subDbModule.Name = subModule.Name
			subDbModule.Description = subModule.Description
			subDbModule.Url = subModule.Url
			subDbModule.Public = subModule.Public
			tx.Save(&subDbModule)
		}
	}
	tx.Commit()
	return err
}

func dropModule(db *gorm.DB) {
	db.DropTableIfExists(&Module{})
}
