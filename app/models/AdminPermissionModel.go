package models

import (
	"github.com/peixiaochen/BaseAdmin/pkg/database"
)

type AdminPermissionModel struct {
	Mysql
	//此处就是用空格
	Name        string `json:"name" gorm:"unique;<-:create"`
	Description string `json:"description" gorm:"<-:update"`
	HttpMethod  string `json:"http_method"`
	HttpPath    string `json:"http_path"`
}

func (AdminPermissionModel) TableName() string {
	return "admin_permission"
}
func (a *AdminPermissionModel) Insert() (Id uint64, err error) {
	// 开始事务
	tx := database.Db.Begin()
	if res := tx.Create(&a); res.Error != nil {
		// 遇到错误时回滚事务
		err = res.Error
		tx.Rollback()
		return
	}
	//添加权限与角色的中间表
	if res := tx.Create(&AdminRolePermissionModel{RoleId: 1, PermissionId: a.ID}); res.Error != nil {
		// 遇到错误时回滚事务
		err = res.Error
		tx.Rollback()
		return
	}
	Id = uint64(a.ID)
	// 否则，提交事务
	tx.Commit()
	return
}

func (a *AdminPermissionModel) Update(Id uint) (Rows int64, err error) {
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
	// 否则，提交事务
	tx.Commit()
	return
}

func (a *AdminPermissionModel) Delete(Ids string) (Rows int64, err error) {
	res := database.Db.Where("id in (?)", Ids).Delete(&a)
	if res.Error != nil {
		err = res.Error
		return
	}
	Rows = res.RowsAffected
	return
}

type AdminPermissionList struct {
	Id          uint                       `json:"id" `
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	HttpMethod  string                     `json:"http_method"`
	HttpPath    string                     `json:"http_path"`
	CreatedAt   JsonTime                   `json:"created_at"`
	UpdatedAt   JsonTime                   `json:"updated_at"`
	Roles       []*AdminPermissionRoleList `json:"roles"`
}

func (a *AdminPermissionModel) GetAll(RequestData map[string]interface{}, Offset uint64, Limit uint8) (Data interface{}, err error) {
	var (
		count                    uint64
		adminPermissionList      []*AdminPermissionList
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
		var row AdminPermissionList
		_ = database.Db.ScanRows(rows, &row)
		row.Roles, _ = AdminRolePermissionModel.GetPermissionRoleList(row.Id)
		adminPermissionList = append(adminPermissionList, &row)
	}

	Data = map[string]interface{}{"count": count, "data_list": adminPermissionList}
	return
}

func (a *AdminPermissionModel) GetOne(Id uint) (Data interface{}, err error) {
	panic("implement me")
}

func (a *AdminPermissionModel) GetPermissionCheck(path string, roleIds []uint) (adminPermissionList AdminPermissionList, count uint64, err error) {
	var (
		permissionIds    []uint
		rolePermissionId []uint
		MenuPermissionId []uint
	)
	database.Db.Table("admin_role_permission").Where("role_id in (?)", roleIds).Pluck("permission_id", &rolePermissionId)
	database.Db.Table("admin_menu_permission").Where("menu_id in (select permission_id from admin_role_menu where role_id in (?))", roleIds).Pluck("permission_id", &MenuPermissionId)
	permissionIds = append(permissionIds, rolePermissionId...)
	permissionIds = append(permissionIds, MenuPermissionId...)
	err = database.Db.Model(&a).
		Where("http_path like ?", "%"+path+"%").
		Where("id in (?)", permissionIds).
		Limit(1).Scan(&adminPermissionList).Count(&count).Error
	return
}
