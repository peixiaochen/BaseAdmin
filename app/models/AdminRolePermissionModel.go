package models

import (
	"github.com/peixiaochen/BaseAdmin/pkg/database"
	"time"
)

type AdminRolePermissionModel struct {
	RoleId       uint      `json:"role_id"`
	PermissionId uint      `json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (AdminRolePermissionModel) TableName() string {
	return "admin_role_permission"
}

func (a *AdminRolePermissionModel) Insert(err error) {
	if res := database.Db.Create(&a); res.Error != nil {
		err = res.Error
	}
	return
}

type AdminRolePermissionList struct {
	PermissionId uint   `json:"permission_id" `
	Name         string `json:"name"`
}

func (a AdminRolePermissionModel) GetRolePermissionList(RoleId uint) (Data []*AdminRolePermissionList, err error) {
	err = database.Db.Table("admin_role_permission as rp").
		Where("rp.role_id = ?", RoleId).
		Select("rp.permission_id,p.name").
		Joins("left join admin_permission as p on p.id = rp.permission_id").
		Scan(&Data).Error
	return
}

type AdminPermissionRoleList struct {
	RoleId uint   `json:"role_id" `
	Name   string `json:"name"`
}

func (a AdminRolePermissionModel) GetPermissionRoleList(PermissionId uint) (Data []*AdminPermissionRoleList, err error) {
	err = database.Db.Table("admin_role_permission as rp").
		Where("rp.permission_id = ?", PermissionId).
		Select("rp.role_id,r.name").
		Joins("left join admin_role as r on r.id = rp.role_id").
		Scan(&Data).Error
	return
}
