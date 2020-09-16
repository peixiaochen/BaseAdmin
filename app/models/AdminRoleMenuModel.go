package models

import (
	"github.com/peixiaochen/BaseAdmin/pkg/database"
	"time"
)

type AdminRoleMenuModel struct {
	RoleId    uint     `json:"role_id"`
	MenuId    uint     `json:"menu_id"`
	CreatedAt JsonTime `json:"created_at"`
	UpdatedAt JsonTime `json:"updated_at"`
}

func (AdminRoleMenuModel) TableName() string {
	return "admin_role_menu"
}

func (a *AdminRoleMenuModel) Insert(err error) {
	if res := database.Db.Create(&a); res.Error != nil {
		err = res.Error
	}
	return
}

type AdminRoleMenuList struct {
	MenuId uint   `json:"menu_id" `
	Name   string `json:"title"`
}

func (a AdminRoleMenuModel) GetRoleMenuList(RoleId uint) (Data []*AdminRoleMenuList, err error) {
	err = database.Db.Table("admin_role_menu as rm").
		Where("rm.role_id = ?", RoleId).
		Select("rm.menu_id,m.title").
		Joins("left join admin_menu as m on m.id = rm.menu_id").
		Scan(&Data).Error
	return
}
