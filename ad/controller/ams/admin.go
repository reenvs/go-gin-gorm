package ams

import (
	"errors"
	"iptv/ad/model"
	"iptv/ad/service"
	"iptv/common/constant"
	"iptv/common/logger"
	"iptv/common/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

/*
	POST /admin/login
	管理员登录
	@Author:HYK
	http://localhost:2000/#!./ams/ams-admin.md
*/
func AdminLoginHandler(c *gin.Context) {
	type param struct {
		Account  string `json:"account" binding:"required"`  //username 或 mobile 或 email
		Password string `json:"password" binding:"required"` //登录密码, password, smscode至少需要一项有值
	}
	type adminInfo struct {
		Id       uint32 `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
		Token    string `json:"token"`
	}
	var p param
	var err error
	if err = c.ShouldBind(&p); err != nil {
		logger.Debug("Invalid request param ", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var dbAdmin model.Admin
	var info adminInfo

	// add operation log when handler return
	defer func() {
		// not log the password
		c.Set(constant.ContextRequestBody, &param{Account: p.Account})
		c.Set(constant.ContextTableName, model.Admin{}.TableName())
		c.Set(constant.ContextOperationType, model.OperationTypeLogin)
		c.Set(constant.ContextAdmin, &dbAdmin)
		if err != nil {
			c.Set(constant.ContextError, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": info})
	}()

	// Find the user
	if err = db.Where("(username = ? OR email = ? OR mobile = ?)", p.Account, p.Account, p.Account).First(&dbAdmin).Error; err != nil {
		logger.Error(err)
		c.JSON(http.StatusOK, gin.H{"err_code": constant.AdminNotExists, "err_msg": constant.TranslateErrCode(constant.AdminNotExists)})
		return
	}
	if err = db.Model(&dbAdmin).Related(&dbAdmin.Roles, "Roles").Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	if dbAdmin.Password != util.SHA512(p.Password) {
		err = errors.New("wrong password")
		logger.Error(err)
		c.JSON(http.StatusOK, gin.H{"err_code": constant.WrongUsernamePassword, "err_msg": constant.TranslateErrCode(constant.WrongUsernamePassword)})
		return
	}

	// if admin login successfully, then create an admin token
	token := service.GenerateTokenForAdmin(&dbAdmin, c.Request.RemoteAddr, db)
	info.Id = dbAdmin.Id
	info.Username = dbAdmin.Username
	info.Role = dbAdmin.Roles[0].Name
	info.Token = token
}

/*
	POST /admin/ftp/login
	管理员ftp登录
	@Author:
	http://localhost:2000/#!./ams/ams-admin.md
*/
func AdminFtpLoginHandler(c *gin.Context) {
	type param struct {
		Account  string `json:"account" binding:"required"`  //username 或 mobile 或 email
		Password string `json:"password" binding:"required"` //登录密码, password, smscode至少需要一项有值
	}

	var p param
	var err error
	if err = c.Bind(&p); err != nil {
		logger.Debug("Invalid request param ", err)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var dbAdmin model.Admin

	// add operation log when handler return
	defer func() {
		// not log the password
		c.Set(constant.ContextRequestBody, &param{Account: p.Account})
		c.Set(constant.ContextTableName, model.Admin{}.TableName())
		c.Set(constant.ContextOperationType, model.OperationTypeLogin)
		c.Set(constant.ContextAdmin, &dbAdmin)
		if err != nil {
			c.Set(constant.ContextError, err.Error())
			return
		}
		c.Set(constant.ContextAdminId, dbAdmin.Id)
		c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": dbAdmin})
	}()

	// Find the user
	if err = db.Where("(username = ? OR email = ? OR mobile = ?)", p.Account, p.Account, p.Account).First(&dbAdmin).Error; err != nil {
		logger.Error(err)
		c.JSON(http.StatusOK, gin.H{"err_code": constant.AdminNotExists, "err_msg": constant.TranslateErrCode(constant.AdminNotExists)})
		return
	}

	if dbAdmin.Password != util.SHA512(p.Password) {
		err = errors.New("wrong password")
		logger.Error(err)
		c.JSON(http.StatusOK, gin.H{"err_code": constant.WrongUsernamePassword, "err_msg": constant.TranslateErrCode(constant.WrongUsernamePassword)})
		return
	}
}

/*
	POST /admin/verify
	管理员验证
	@Author:HYK
	http://localhost:2000/#!./ams/ams-admin.md
*/
func AdminVerifyHandler(c *gin.Context) {
	type param struct {
		Token string `json:"token" binding:"required"` // 管理员令牌
		Url   string `json:"url" binding:"required"`   // 访问api的url
	}

	var p param
	var err error
	if err = c.Bind(&p); err != nil {
		logger.Debug("Invalid request param ", err)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	// var admin *model.Admin
	// var token *model.Token

	// add operation log when handler return
	/* Do not log this api
	defer func() {
		// not log the password
		c.Set(constant.ContextRequestBody, &p)
		c.Set(constant.ContextTableName, model.Admin{}.TableName())
		c.Set(constant.ContextOperationType, model.OperationTypeLogin)
		c.Set(constant.ContextAdmin, admin)
		if err != nil {
			c.Set(constant.ContextError, err.Error())
		}
	}() */

	//logger.Debugf("AdminVerifyHandler %v", p)
	data, err := GetAdminFromToken(p.Token, db)
	if err != nil {
		logger.Error(err, p.Token, p.Url)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": data})

	/*
		var module model.Module
		if err = db.Where("url = ?", p.Url).First(&module).Error; err != nil {
			logger.Error(err)
			c.JSON(http.StatusOK, gin.H{"err_code": constant.ApiNotRegisted, "err_msg": constant.TranslateErrCode(constant.ApiNotRegisted)})
			return
		}

		allowed := true //service.IsModuleAllowed(&module, admin, db)
		if !allowed {
			logger.Error(err)
			c.JSON(http.StatusOK, gin.H{"err_code": constant.ModuleAccessDenied, "err_msg": constant.TranslateErrCode(constant.ModuleAccessDenied)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": admin})
	*/
}

/*
	POST /admin/logout
	退出登陆
	http://localhost:2000/#!./ams/ams-admin.md
*/
func AdminLogoutHandler(c *gin.Context) {
	// add operation log when handler return
	defer func() {
		c.Set(constant.ContextTableName, model.Admin{}.TableName())
		c.Set(constant.ContextOperationType, model.OperationTypeLogout)
	}()

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	adminId := c.MustGet(constant.ContextAdminId).(uint32)

	// delete all admin's token
	if err := db.Exec("delete from token where admin_id = ?", adminId).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}

/*
	GET /admin/list
	获取管理员列表
	@Author: HYK
	http://localhost:2000/#!./ims/ims-admin.md
*/
func AdminListHandler(c *gin.Context) {
	type param struct {
		Offset   int    `form:"offset" binding:"exists"`
		Limit    int    `form:"limit" binding:"required"`
		Order    string `form:"order"`
		Sort     string `form:"sort"`
		Username string `form:"username"`
		RoleId   uint32 `form:"role_id"`
	}

	var p param
	var err error
	if err = c.ShouldBind(&p); err != nil {
		logger.Error("Invalid Product param ", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	if p.Offset < 0 ||
		p.Limit <= 0 ||
		p.Limit > constant.MaxPageSize {
		err = errors.New("bad request")
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, util.Err2JsonObj(err))
		return
	}

	// set id as default order
	if p.Order == "" {
		p.Order = "id"
	}

	sort := constant.SortDesc
	if p.Sort == constant.SortAsc {
		sort = constant.SortAsc
	}

	var admins []*model.Admin
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	dbExec := db

	if p.Username != "" {
		dbExec = dbExec.Where("username like ?", "%"+strings.TrimSpace(p.Username)+"%")
	}

	if p.RoleId != 0 {
		dbExec = dbExec.Joins("inner join admin_role ar on admin.id = ar.admin_id")
		dbExec = dbExec.Where("ar.role_id = ?", p.RoleId)
	}

	dbCount := db
	if err = dbExec.Offset(p.Offset).Limit(p.Limit).Order(p.Order + " " + sort).Find(&admins).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	for _, admin := range admins {
		if err = db.Model(&admin).Related(&admin.Roles, "Roles").Error; err != nil {
			logger.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
			return
		}

		for _, role := range admin.Roles {
			role.Transcope()
		}
	}

	var count uint32
	if err = dbCount.Model(&model.Admin{}).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": admins, "total": count})
}

/*
	POST /admin/save
	保存、发布管理员
	@Author: HYK
	http://localhost:2000/#!./ims/ims-admin.md
*/
func AdminSaveHandler(c *gin.Context) {
	var admin, oldValue, dbAdmin model.Admin
	var err error
	var operationType int
	if err = c.Bind(&admin); err != nil {
		logger.Error(err)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var isUpdate bool
	tx := db.Begin()
	defer func() {
		if isUpdate {
			operationType = model.OperationTypeUpdate
			c.Set(constant.ContextOldValue, &oldValue)
		} else {
			operationType = model.OperationTypeCreate
		}
		c.Set(constant.ContextRequestBody, &admin)
		c.Set(constant.ContextTableName, model.Admin{}.TableName())
		c.Set(constant.ContextOperationType, operationType)

		if err != nil {
			tx.Rollback()
			c.Set(constant.ContextError, err.Error())
			return
		}
		tx.Commit()
		c.Set(constant.ContextNewValue, &dbAdmin)
		c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": &dbAdmin})
	}()

	if admin.Id == 0 {
		isUpdate = false
	} else {
		isUpdate = true
		dbAdmin.Id = admin.Id
		if err = tx.First(&dbAdmin).Error; err != nil {
			logger.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
			return
		}
		oldValue = dbAdmin
	}
	// find new roles
	var newRoleIds []uint32
	for _, role := range admin.Roles {
		newRoleIds = append(newRoleIds, role.Id)
	}
	var newRoles []*model.Role
	if err = tx.Where("id in (?)", newRoleIds).Find(&newRoles).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	// update dbmodel
	dbAdmin.Username = admin.Username
	if admin.Email.String != "" {
		dbAdmin.Email.Valid = true
		dbAdmin.Email.String = admin.Email.String
	} else {
		dbAdmin.Email.Valid = false
	}
	if admin.Mobile.String != "" {
		dbAdmin.Mobile.Valid = true
		dbAdmin.Mobile.String = admin.Mobile.String
	} else {
		dbAdmin.Mobile.Valid = false
	}
	dbAdmin.Realname = admin.Realname
	dbAdmin.Locked = admin.Locked
	// if admin's password is not empty, update it
	if admin.Password != "" && admin.Password != "######" {
		dbAdmin.Password = util.SHA512(admin.Password)
	}

	if isUpdate {
		if err = tx.Model(&dbAdmin).Association("Roles").Replace(newRoles).Error; err != nil {
			logger.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
			return
		}
		if err = tx.Save(&dbAdmin).Error; err != nil {
			logger.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
			return
		}
	} else {
		dbAdmin.Roles = newRoles
		if err = tx.Create(&dbAdmin).Error; err != nil {
			logger.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
			return
		}
	}
}

/*
   GET /admin/detail
   获取角色详情
   @Author: Zhangli
*/
func AdminDetailHandler(c *gin.Context) {
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

	var admin model.Admin
	if err = db.First(&admin, p.Id).Related(&admin.Roles, "Roles").Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": &admin})
}

/*
	POST /admin/delete
	删除管理员
	@Author: HYK
	http://localhost:2000/#!./ims/ims-admin.md
*/
func AdminDeleteHandler(c *gin.Context) {
	type param struct {
		Id []uint32 `json:"id" binding:"required"`
	}

	var p param
	var err error
	if err = c.Bind(&p); err != nil {
		logger.Error("Invalid Admin param ", err)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	tx := db.Begin()
	defer func() {
		c.Set(constant.ContextRequestBody, &p)
		c.Set(constant.ContextTableName, model.Admin{}.TableName())
		c.Set(constant.ContextOperationType, model.OperationTypeDelete)

		if err != nil {
			tx.Rollback()
			c.Set(constant.ContextError, err.Error())
			return
		}
		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
	}()

	// delete all admins
	if err = tx.Where("id in (?)", p.Id).Delete(&model.Admin{}).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}

	// remove all role_module association attached to this module
	if err = tx.Exec("delete from admin_role where admin_id in (?)", p.Id).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(err))
		return
	}
}
