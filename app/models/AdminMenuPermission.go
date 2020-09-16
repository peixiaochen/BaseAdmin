package models

import (
	"github.com/peixiaochen/BaseAdmin/pkg/database"
)

type AdminMenuPermissionModel struct {
	MenuId       uint     `json:"menu_id"`
	PermissionId uint     `json:"permission_id"`
	CreatedAt    JsonTime `json:"created_at"`
	UpdatedAt    JsonTime `json:"updated_at"`
}

func (AdminMenuPermissionModel) TableName() string {
	return "admin_menu_permission"
}

func (a *AdminMenuPermissionModel) Insert(err error) {
	if res := database.Db.Create(&a); res.Error != nil {
		err = res.Error
	}
	return
}

type AdminMenuPermissionList struct {
	PermissionId uint   `json:"permission_id" `
	Name         string `json:"name"`
}

func (a AdminMenuPermissionModel) GetMenuPermissionList(MenuId uint) (Data []*AdminMenuPermissionList, err error) {
	err = database.Db.Table("admin_menu_permission as mp").
		Where("mp.menu_id = ?", MenuId).
		Select("mp.permission_id,p.name").
		Joins("left join admin_permission as p on p.id = mp.permission_id").
		Scan(&Data).Error
	return
}
