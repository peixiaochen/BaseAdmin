package models

import (
	"github.com/peixiaochen/BaseAdmin/pkg/database"
)

type AdminRoleModel struct {
	Mysql
	//此处就是用空格
	Name          string `json:"name" gorm:"unique;<-:create"`
	Description   string `json:"description" gorm:"<-:update"`
	Status        uint8  `json:"status"`
	MenusIds      []uint `json:"menus" gorm:"-"`
	PermissionIds []uint `json:"permissions" gorm:"-"`
}

func (AdminRoleModel) TableName() string {
	return "admin_role"
}

func (a *AdminRoleModel) Insert() (Id uint64, err error) {
	// 开始事务
	tx := database.Db.Begin()
	if res := tx.Create(&a); res.Error != nil {
		// 遇到错误时回滚事务
		err = res.Error
		tx.Rollback()
		return
	}
	//添加用户与角色的中间表
	for _, v := range a.MenusIds {
		if res := tx.Create(&AdminRoleMenuModel{RoleId: a.ID, MenuId: v}); res.Error != nil {
			// 遇到错误时回滚事务
			err = res.Error
			tx.Rollback()
			return
		}
	}
	//添加用户与角色的中间表
	for _, v := range a.PermissionIds {
		if res := tx.Create(&AdminRolePermissionModel{RoleId: a.ID, PermissionId: v}); res.Error != nil {
			// 遇到错误时回滚事务
			err = res.Error
			tx.Rollback()
			return
		}
	}
	Id = uint64(a.ID)
	// 否则，提交事务
	tx.Commit()
	return
}

func (a *AdminRoleModel) Update(Id uint) (Rows int64, err error) {
	// 开始事务
	tx := database.Db.Begin()
	res := tx.Model(&a).Where("id = ?", Id).Updates(&a)
	if res.Error != nil {
		// 遇到错误时回滚事务
		err = res.Error
		tx.Rollback()
		return
	}
	Rows = res.RowsAffected
	//更新角色与菜单的中间表
	if res := tx.Delete(AdminRoleMenuModel{}, "role_id = ?", Id); res.Error != nil {
		// 遇到错误时回滚事务
		err = res.Error
		tx.Rollback()
		return
	}
	for _, v := range a.MenusIds {
		if res := tx.Create(&AdminRoleMenuModel{RoleId: Id, MenuId: v}); res.Error != nil {
			// 遇到错误时回滚事务
			err = res.Error
			tx.Rollback()
			return
		}
	}
	//更新角色与权限的中间表
	if res := tx.Delete(AdminRolePermissionModel{}, "role_id = ?", Id); res.Error != nil {
		// 遇到错误时回滚事务
		err = res.Error
		tx.Rollback()
		return
	}
	for _, v := range a.PermissionIds {
		if res := tx.Create(&AdminRolePermissionModel{RoleId: Id, PermissionId: v}); res.Error != nil {
			// 遇到错误时回滚事务
			err = res.Error
			tx.Rollback()
			return
		}
	}
	// 否则，提交事务
	tx.Commit()
	return
}

func (a *AdminRoleModel) Delete(Ids string) (Rows int64, err error) {
	res := database.Db.Where("id in (?)", Ids).Delete(&a)
	if res.Error != nil {
		err = res.Error
		return
	}
	Rows = res.RowsAffected
	return
}

type AdminRoleList struct {
	Id          uint                       `json:"id" `
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Status      uint8                      `json:"status"`
	CreatedAt   JsonTime                   `json:"created_at"`
	UpdatedAt   JsonTime                   `json:"updated_at"`
	Menus       []*AdminRoleMenuList       `json:"menus"`
	Permissions []*AdminRolePermissionList `json:"permissions"`
}

func (a *AdminRoleModel) GetAll(RequestData map[string]interface{}, Offset uint64, Limit uint8) (Data interface{}, err error) {

	var (
		count                    uint64
		adminRoleList            []*AdminRoleList
		AdminRoleMenuModel       AdminRoleMenuModel
		AdminRolePermissionModel AdminRolePermissionModel
	)
	model := database.Db.Model(&a)
	if RequestData["Keywords"] != "%%" {
		model = model.Where("name like ? or description like ?", RequestData["Keywords"], RequestData["Keywords"])
	}
	if RequestData["StartTime"] != "" {
		model = model.Where("created_at >= ?", RequestData["StartTime"])
	}
	if RequestData["EndTime"] != "" {
		model = model.Where("created_at <= ?", RequestData["EndTime"])
	}
	model.Count(&count)
	model = model.Order("id desc").Offset(Offset).Limit(Limit)
	if model.Error != nil {
		err = model.Error
		return
	}
	rows, err := model.Rows()
	if err != nil {
		err = model.Error
		return
	}
	defer rows.Close()
	for rows.Next() {
		var row AdminRoleList
		_ = database.Db.ScanRows(rows, &row)
		row.Menus, _ = AdminRoleMenuModel.GetRoleMenuList(row.Id)
		row.Permissions, _ = AdminRolePermissionModel.GetRolePermissionList(row.Id)
		adminRoleList = append(adminRoleList, &row)
	}

	Data = map[string]interface{}{"count": count, "data_list": adminRoleList}
	return
}

func (a *AdminRoleModel) GetOne(Id uint) (Data interface{}, err error) {
	panic("implement me")
}
