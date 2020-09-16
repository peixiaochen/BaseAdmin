package Requests

import (
	"github.com/gookit/validate"
	"github.com/peixiaochen/BaseAdmin/pkg/database"
)

type AdminPermissionSaveRequest struct {
	Localtime uint64                          `json:"localtime" binding:"required"`
	Data      *AdminPermissionSaveRequestData `json:"data" binding:"required"`
}
type AdminPermissionSaveRequestData struct {
	Id          uint   `json:"id" validate:"uint|CheckId"`
	Name        string `json:"name" validate:"required|Unique"`
	Description string `json:"description" validate:"string|minLen:2"`
	HttpMethod  string `json:"http_method" validate:"string|minLen:2"`
	HttpPath    string `json:"http_path" validate:"required|string|minLen:2"`
}

// Messages 您可以自定义验证器错误消息
func (a AdminPermissionSaveRequestData) Messages() map[string]string {
	return validate.MS{
		"Id.CheckId":  "{field}当前不存在，请检查后重新提交",
		"Name.Unique": "{field}已经存在，请勿重复使用",
	}
}

// Translates 你可以自定义字段翻译
func (a AdminPermissionSaveRequestData) Translates() map[string]string {
	return validate.MS{
		"Name":        "菜单名",
		"Description": "描述",
	}
}

//自定义检查用户名是否唯一
func (a AdminPermissionSaveRequestData) CheckId(val uint) bool {
	var count uint64
	err := database.Db.
		Table("admin_permission").
		Where("id = ?", val).
		Count(&count).Error
	if err != nil || count == 0 {
		return false
	}
	return true
}

//自定义检查用户名是否唯一
func (a AdminPermissionSaveRequestData) Unique(val string) bool {
	var count uint64
	db := database.Db.
		Table("admin_permission").
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

//角色删除
type AdminPermissionDeleteRequest struct {
	Localtime uint64                            `json:"localtime" binding:"required"`
	Data      *AdminPermissionDeleteRequestData `json:"data" binding:"required"`
}
type AdminPermissionDeleteRequestData struct {
	Ids string `json:"ids" validate:"string"`
}

//权限列表
type AdminPermissionIndexRequest struct {
	Localtime uint64                           `json:"localtime" binding:"required"`
	Data      *AdminPermissionIndexRequestData `json:"data" binding:"required"`
}
type AdminPermissionIndexRequestData struct {
	Page      uint64 `json:"page" binding:"required,gte=1"`
	Limit     uint8  `json:"limit" binding:"required,gte=1"`
	Keywords  string `json:"keywords" binding:"min:2"`
	StartTime string `json:"start_time" binding:"required" time_format:"2020-09-01 10:22:22"`
	EndTime   string `json:"end_time" binding:"required" time_format:"2020-09-01 10:22:22"`
}
