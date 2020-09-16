package Requests

import (
	"github.com/gookit/validate"
	"github.com/peixiaochen/BaseAdmin/pkg/config"
	"github.com/peixiaochen/BaseAdmin/pkg/database"
	"github.com/peixiaochen/BaseAdmin/pkg/lib"
)

/**
用户登录
*/
type AdminUserLoginRequest struct {
	Localtime uint64                     `json:"localtime" binding:"required"`
	Data      *AdminUserLoginDataRequest `json:"data" binding:"required"`
}
type AdminUserLoginDataRequest struct {
	UserName string `json:"username" validate:"required|string|CheckExists"`
	Password string `json:"password" validate:"required|string|CheckPasswordIsEq"`
}

// Messages 您可以自定义验证器错误消息
func (a AdminUserLoginDataRequest) Messages() map[string]string {
	return validate.MS{
		"UserName.CheckExists":       "当前 {field} 不存在或已禁止登录，请联系管理员",
		"Password.CheckPasswordIsEq": "{field}不正确，请检查后重新输入",
	}
}

// Translates 你可以自定义字段翻译
func (a AdminUserLoginDataRequest) Translates() map[string]string {
	return validate.MS{
		"UserName": "用户名",
		"Password": "密码",
	}
}

//自定义检查用户名是否存在
func (a AdminUserLoginDataRequest) CheckExists(val string) bool {
	var count uint64
	err := database.Db.
		Table("admin_user").Select("id").
		Where("username = ?", val).
		Where("status = ?", 1).
		Where("deleted_at is null").
		Count(&count).Error
	if err != nil || count == 0 {
		return false
	}
	return true
}

//自定义检查密码是否相同
func (a AdminUserLoginDataRequest) CheckPasswordIsEq(val string) bool {
	var count uint64
	err := database.Db.
		Table("admin_user").Select("id").
		Where("username = ?", a.UserName).
		Where("password = ?", lib.SubstrMd5(val, config.PasswordSetting.PasswordStart, config.PasswordSetting.PasswordLength)).
		Where("status = ?", 1).
		Where("deleted_at is null").
		Count(&count).Error
	if err != nil || count == 0 {
		return false
	}
	return true
}

/**
用户保存
*/
type AdminUserSaveRequest struct {
	Localtime uint64                    `json:"localtime" binding:"required"`
	Data      *AdminUserSaveRequestData `json:"data" binding:"required"`
}
type AdminUserSaveRequestData struct {
	Id                   uint                           `json:"id" validate:"uint|CheckId"`
	Name                 string                         `json:"name" validate:"required|string"`
	UserName             string                         `json:"username" validate:"required|string|Unique|minLen:2"`
	Password             string                         `json:"password" validate:"required_without:Id|minLen:6"`
	PasswordConfirmation string                         `json:"password_confirmation" validate:"required_without:Id|eq_field:Password"`
	Avatar               string                         `json:"avatar" validate:"string"`
	Extra                *AdminUserSaveRequestDataExtra `json:"extra" validate:"required"`
	Roles                []uint                         `json:"roles" validate:"isInts|CheckCount"`
}

// Messages 您可以自定义验证器错误消息
func (a AdminUserSaveRequestData) Messages() map[string]string {
	return validate.MS{
		"UserName.Unique":  "{field}已经存在，请勿重复使用",
		"Roles.CheckCount": "{field}角色数量有误，请检查后重新提交",
	}
}

// Translates 你可以自定义字段翻译
func (a AdminUserSaveRequestData) Translates() map[string]string {
	return validate.MS{
		"UserName": "用户名",
		"Roles":    "角色ID",
	}
}

//自定义检查用户名是否唯一
func (a AdminUserSaveRequestData) CheckId(val uint) bool {
	var count uint64
	err := database.Db.
		Table("admin_user").
		Where("id = ?", val).
		Count(&count).Error
	if err != nil || count == 0 {
		return false
	}
	return true
}

//自定义检查用户名是否唯一
func (a AdminUserSaveRequestData) Unique(val string) bool {
	var count uint64
	db := database.Db.
		Table("admin_user").
		Select("id").
		Where("username = ?", val)
	if a.Id > 0 {
		db = db.Not("id = ?", a.Id)
	}
	err := db.Count(&count).Error
	if err != nil || count > 0 {
		return false
	}
	return true
}
func (a AdminUserSaveRequestData) CheckCount(val []uint) bool {
	var count int
	err := database.Db.
		Table("admin_role").
		Select("id").
		Where("id in (?)", val).
		Count(&count).Error
	if err != nil || count != len(val) {
		return false
	}
	return true
}

type AdminUserSaveRequestDataExtra struct {
	Email string `json:"email" validate:"string|email"`
	Phone string `json:"phone" validate:"string|isCnMobile"`
	Sex   int    `json:"sex" validate:"int|enum:0,1,2"`
}

// Messages 您可以自定义验证器错误消息
func (a AdminUserSaveRequestDataExtra) Messages() map[string]string {
	return validate.MS{}
}

// Translates 你可以自定义字段翻译
func (a AdminUserSaveRequestDataExtra) Translates() map[string]string {
	return validate.MS{
		"Email": "邮箱",
		"Phone": "手机号码",
		"Sex":   "性别",
	}
}

type AdminUserDeleteRequest struct {
	Localtime uint64                      `json:"localtime" binding:"required"`
	Data      *AdminUserDeleteRequestData `json:"data" binding:"required"`
}
type AdminUserDeleteRequestData struct {
	Ids string `json:"ids" validate:"string"`
}
type AdminUserIndexRequest struct {
	Localtime uint64                     `json:"localtime" binding:"required"`
	Data      *AdminUserIndexRequestData `json:"data" binding:"required"`
}
type AdminUserIndexRequestData struct {
	Page      uint64 `json:"page" binding:"required,gte=1"`
	Limit     uint8  `json:"limit" binding:"required,gte=1"`
	Keywords  string `json:"keywords" binding:"min:2"`
	RoleId    uint64 `json:"role_id" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

type AdminUserStatusRequest struct {
	Localtime uint64                      `json:"localtime" binding:"required"`
	Data      *AdminUserStatusRequestData `json:"data" binding:"required"`
}
type AdminUserStatusRequestData struct {
	Ids    string `json:"ids" binding:"required"`
	Status uint8  `json:"status" binding:"required,in:0,1"`
}
