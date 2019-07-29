package ams

import (
	"encoding/json"
	"iptv/ad/model"
	"iptv/common/constant"
	"iptv/common/logger"
	"iptv/common/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

/*
	GET /role/list
	获取角色列表
*/
func RoleListHandler(c *gin.Context) {
	type param struct {
		Name string `form:"name"`
	}
	var p param
	var err error

	if err = c.ShouldBind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	var roles []*model.Role
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	db = db.Select("role.*,(select count(1) from admin_role where role_id=role.id) as admin_count")
	if p.Name != "" {
		db = db.Where("name like ?", "%"+strings.TrimSpace(p.Name)+"%")
	}

	if err = db.Find(&roles).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	for _, role := range roles {
		role.Transcope()
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": roles})
}

/*
   GET /role/detail
   获取角色详情
   @Author: Zhangli
*/
func RoleDetailHandler(c *gin.Context) {
	type param struct {
		Id string `form:"id" binding:"required"`
	}

	var p param
	var err error
	if err = c.ShouldBind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var role model.Role
	if err = db.First(&role, p.Id).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}
	role.Transcope()
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": role})
}

/*
	POST /role/save
	保存、发布角色
*/
func RoleSaveHandler(c *gin.Context) {
	type param struct {
		model.Role
		Scope struct {
			AppIds []uint32 `json:"app_ids"`
		} `json:"scope"`
	}
	var role param
	var oldValue, dbRole model.Role
	var err error
	var operationType int

	if err = c.Bind(&role); err != nil {
		logger.Error(err)
		return
	}

	var isUpdate bool
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	tx := db.Begin()
	defer func() {
		if isUpdate {
			operationType = model.OperationTypeUpdate
			c.Set(constant.ContextOldValue, &oldValue)
		} else {
			operationType = model.OperationTypeCreate
		}
		c.Set(constant.ContextRequestBody, &role)
		c.Set(constant.ContextTableName, model.Role{}.TableName())
		c.Set(constant.ContextOperationType, operationType)

		if err != nil {
			tx.Rollback()
			c.Set(constant.ContextError, err.Error())
			return
		}
		tx.Commit()
		c.Set(constant.ContextNewValue, &dbRole)
		c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": &dbRole})
	}()

	if role.Id == 0 {
		isUpdate = false
	} else {
		isUpdate = true
		dbRole.Id = role.Id
		if err = tx.First(&dbRole).Error; err != nil {
			logger.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
			return
		}
		oldValue = dbRole
	}

	dbRole.Name = role.Name
	dbRole.Description = role.Description
	dbRole.Scope, _ = json.Marshal(role.Scope)

	if isUpdate {
		if err = tx.Save(&dbRole).Error; err != nil {
			logger.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
			return
		}
	} else {
		if err = tx.Create(&dbRole).Error; err != nil {
			logger.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
			return
		}
	}

	dbRole.Transcope()
}

/*
	POST /role/delete
	删除角色
*/
func RoleDeleteHandler(c *gin.Context) {
	type param struct {
		Id []uint32 `json:"id" binding:"required"`
	}

	var p param
	var err error
	if err = c.Bind(&p); err != nil {
		logger.Error("Invalid Role param ", err)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	tx := db.Begin()
	defer func() {
		c.Set(constant.ContextRequestBody, &p)
		c.Set(constant.ContextTableName, model.Role{}.TableName())
		c.Set(constant.ContextOperationType, model.OperationTypeDelete)

		if err != nil {
			tx.Rollback()
			c.Set(constant.ContextError, err.Error())
			return
		}
		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
	}()

	// delete all roles
	if err = tx.Where("id in (?)", p.Id).Delete(&model.Role{}).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	// remove all role_page association attached to this page
	if err = tx.Exec("delete from role_page where role_id in (?)", p.Id).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}
}

/*
	POST /role/scope/add
	为admin的角色新增app权限
*/
func RoleScopeAdd(c *gin.Context) {
	type param struct {
		AppId   uint32 `json:"app_id" binding:"required"`
		AdminId uint32 `json:"admin_id" binding:"required"`
	}

	var p param
	var err error
	if err = c.Bind(&p); err != nil {
		logger.Error("Invalid scope param ", err)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var admin model.Admin
	admin.Id = p.AdminId
	if err = db.Model(admin).Related(&admin.Roles, "Roles").Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	for _, role := range admin.Roles {
		role.Transcope()

		scopes := role.Scope.(model.RoleScope)
		scopes.AppIds = append(scopes.AppIds, p.AppId)

		data, _ := json.Marshal(scopes)

		if err = db.Model(role).Update("scope", data).Error; err != nil {
			logger.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}
