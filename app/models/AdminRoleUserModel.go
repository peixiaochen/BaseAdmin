package models

import (
	"github.com/peixiaochen/BaseAdmin/pkg/database"
	"time"
)

type AdminRoleUserModel struct {
	RoleId    uint      `json:"role_id"`
	UserId    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (AdminRoleUserModel) TableName() string {
	return "admin_role_user"
}

func (a *AdminRoleUserModel) Insert(err error) {
	if res := database.Db.Create(&a); res.Error != nil {
		err = res.Error
	}
	return
}

type AdminRoleUserList struct {
	RoleId uint   `json:"role_id" `
	Name   string `json:"name"`
}

func (a AdminRoleUserModel) GetUserRoleList(UserId uint) (Data []*AdminRoleUserList, err error) {
	err = database.Db.Table("admin_role_user as ru").
		Where("ru.user_id = ?", UserId).
		Select("ru.role_id,r.name").
		Joins("left join admin_role as r on r.id = ru.role_id").
		Scan(&Data).Error
	if err != nil {
		Data = []*AdminRoleUserList{}
	}
	return
}
func (a AdminRoleUserModel) GetRoleIds(UserId uint) (RoleIds []uint, err error) {
	err = database.Db.Model(&a).Where("user_id = ?", UserId).Pluck("role_id", &RoleIds).Error
	return
}
