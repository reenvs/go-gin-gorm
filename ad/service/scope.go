package service

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"iptv/ad/model"
	"iptv/common/util"
)

//合并管理员所属角色下的全部scope
func MergeAllRoleScope(admin *model.Admin, db *gorm.DB) (string, error) {
	if err := db.Model(admin).Related(&admin.Roles, "Roles").Error; err != nil {
		return "", err
	}
	var scope model.RoleScope
	for _, role := range admin.Roles {
		if len(role.Scope.([]byte)) <= 0 {
			continue
		}
		var tmpScope model.RoleScope
		//logger.Debugf("role.Scope:%v %v", *role, role.Scope)
		if err := json.Unmarshal(role.Scope.([]byte), &tmpScope); err != nil {
			return "", err
		}
		for _, app_id := range tmpScope.AppIds {
			if !util.ContainsInt(scope.AppIds, app_id) {
				scope.AppIds = append(scope.AppIds, app_id)
			}
		}
	}
	if len(scope.AppIds) == 0 {
		return "", nil
	}
	scopeByte, _ := json.Marshal(scope)
	return string(scopeByte), nil
}
