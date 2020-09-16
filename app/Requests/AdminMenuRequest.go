package Requests

import (
	"github.com/gookit/validate"
	"github.com/peixiaochen/BaseAdmin/pkg/database"
)

type AdminMenuSaveRequest struct {
	Localtime uint64                    `json:"localtime" binding:"required"`
	Data      *AdminMenuSaveRequestData `json:"data" binding:"required"`
}
type AdminMenuSaveRequestData struct {
	Id          uint   `json:"id" validate:"uint|CheckId"`
	ParentId    uint   `json:"parent_id" validate:"required|uint|CheckExists"`
	Title       string `json:"title" validate:"required|Unique"`
	Icon        string `json:"icon" validate:"string"`
	Uri         string `json:"uri" validate:"required|string|minLen:2"`
	Permissions []uint `json:"permissions" validate:"isInts|CheckPermissionsCount"`
}

// Messages 您可以自定义验证器错误消息
func (a AdminMenuSaveRequestData) Messages() map[string]string {
	return validate.MS{
		"Id.CheckId":                        "{field}当前不存在，请检查后重新提交",
		"ParentId.Exists":                   "{field}当前父级菜单不存在，请检查后重新提交",
		"Title.Unique":                      "{field}已经存在，请勿重复使用",
		"Permissions.CheckPermissionsCount": "{field}权限数量有误，请检查后重新提交",
	}
}

// Translates 你可以自定义字段翻译
func (a AdminMenuSaveRequestData) Translates() map[string]string {
	return validate.MS{
		"ParentId":    "父级菜单",
		"Title":       "标题",
		"Icon":        "小icon",
		"Uri":         "路径",
		"Permissions": "权限列表",
	}
}

//自定义检查用户名是否唯一
func (a AdminMenuSaveRequestData) CheckId(val uint) bool {
	var count uint64
	err := database.Db.
		Table("admin_menu").
		Where("id = ?", val).
		Count(&count).Error
	if err != nil || count == 0 {
		return false
	}
	return true
}

//检查是否存在标题
func (a AdminMenuSaveRequestData) CheckExists(val uint) bool {
	var count uint8
	err := database.Db.
		Table("admin_menu").Select("id").
		Where("id = ?", val).
		Count(&count).Error
	if err != nil || count == 0 {
		return false
	}
	return true
}

//自定义检查用户名是否唯一
func (a AdminMenuSaveRequestData) Unique(val string) bool {
	var count uint64
	db := database.Db.
		Table("admin_menu").
		Select("id").
		Where("title = ?", val)
	if a.Id > 0 {
		db = db.Not("id = ?", a.Id)
	}
	err := db.Count(&count).Error
	if err != nil || count > 0 {
		return false
	}
	return true
}

func (a AdminMenuSaveRequestData) CheckPermissionsCount(val []uint) bool {
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

type AdminMenuDetailRequest struct {
	Id uint `uri:"Id";binding:"required"`
}

type AdminMenuSortRequest struct {
	Localtime uint64                    `json:"localtime" binding:"required"`
	Data      *AdminMenuSortRequestData `json:"data" binding:"required"`
}
type AdminMenuSortRequestData struct {
	OrderIds []AdminMenuSortRequestDataOrderIds `json:"order_id" validate:"required"`
}
type AdminMenuSortRequestDataOrderIds struct {
	Id       uint                               `json:"id" validate:"required"`
	Children []AdminMenuSortRequestDataOrderIds `json:"children" `
}
