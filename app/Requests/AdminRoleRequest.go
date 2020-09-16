package Requests

import (
	"github.com/gookit/validate"
	"github.com/peixiaochen/BaseAdmin/pkg/database"
)

/**
角色保存
*/
type AdminRoleSaveRequest struct {
	Localtime uint64                    `json:"localtime" binding:"required"`
	Data      *AdminRoleSaveRequestData `json:"data" binding:"required"`
}
type AdminRoleSaveRequestData struct {
	Id          uint   `json:"id" validate:"uint|CheckId"`
	Name        string `json:"name" validate:"required|Unique"`
	Description string `json:"description" validate:"string|minLen:2"`
	Status      uint8  `json:"status" validate:"uint8"`
	Menus       []uint `json:"menus" validate:"isInts|CheckMenusCount"`
	Permissions []uint `json:"permissions" validate:"isInts|CheckPermissionsCount"`
}

// Messages 您可以自定义验证器错误消息
func (a AdminRoleSaveRequestData) Messages() map[string]string {
	return validate.MS{
		"Id.CheckId":                        "{field}当前不存在，请检查后重新提交",
		"Name.Unique":                       "{field}已经存在，请勿重复使用",
		"Menus.CheckMenusCount":             "{field}菜单数量有误，请检查后重新提交",
		"Permissions.CheckPermissionsCount": "{field}权限数量有误，请检查后重新提交",
	}
}

// Translates 你可以自定义字段翻译
func (a AdminRoleSaveRequestData) Translates() map[string]string {
	return validate.MS{
		"Name":        "菜单名",
		"Description": "描述",
		"Menus":       "菜单",
		"Permissions": "权限",
	}
}

//自定义检查用户名是否唯一
func (a AdminRoleSaveRequestData) CheckId(val uint) bool {
	var count uint64
	err := database.Db.
		Table("admin_role").
		Where("id = ?", val).
		Count(&count).Error
	if err != nil || count == 0 {
		return false
	}
	return true
}

//自定义检查用户名是否唯一
func (a AdminRoleSaveRequestData) Unique(val string) bool {
	var count uint64
	db := database.Db.
		Table("admin_role").
		Select("id").
		Where("name = ?", val)
	if a.Id > 0 {
		db = db.Not("id = ?", a.Id)
	}
	err := db.Count(&count).Error
	if err != nil || count > 0 {
		return false
	}
	return true
}
func (a AdminRoleSaveRequestData) CheckMenusCount(val []uint) bool {
	var count int
	err := database.Db.
		Table("admin_menu").
		Select("id").
		Where("id in (?)", val).
		Count(&count).Error
	if err != nil || count != len(val) {
		return false
	}
	return true
}
func (a AdminRoleSaveRequestData) CheckPermissionsCount(val []uint) bool {
	var count int
	err := database.Db.
		Table("admin_permission").
		Select("id").
		Where("id in (?)", val).
		Count(&count).Error
	if err != nil || count != len(val) {
		return false
	}
	return true
}

//角色删除
type AdminRoleDeleteRequest struct {
	Localtime uint64                      `json:"localtime" binding:"required"`
	Data      *AdminRoleDeleteRequestData `json:"data" binding:"required"`
}
type AdminRoleDeleteRequestData struct {
	Ids string `json:"ids" validate:"string"`
}

type AdminRoleIndexRequest struct {
	Localtime uint64                     `json:"localtime" binding:"required"`
	Data      *AdminRoleIndexRequestData `json:"data" binding:"required"`
}
type AdminRoleIndexRequestData struct {
	Page      uint64 `json:"page" binding:"required,gte=1"`
	Limit     uint8  `json:"limit" binding:"required,gte=1"`
	Keywords  string `json:"keywords" binding:"min:2"`
	StartTime string `json:"start_time" binding:"required" time_format:"2020-09-01 10:22:22"`
	EndTime   string `json:"end_time" binding:"required" time_format:"2020-09-01 10:22:22"`
}
