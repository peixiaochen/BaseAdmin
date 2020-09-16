package models

import (
	"fmt"
	"github.com/peixiaochen/BaseAdmin/app/Requests"
	"github.com/peixiaochen/BaseAdmin/pkg/database"
)

type AdminMenuModel struct {
	ID            uint     `gorm:"primary_key" json:"id"`
	ParentId      uint     `json:"parent_id"`
	Order         int      `json:"order"`
	Title         string   `json:"title"`
	Icon          string   `json:"icon"`
	Uri           string   `json:"uri"`
	CreatedAt     JsonTime `json:"created_at"`
	UpdatedAt     JsonTime `json:"updated_at"`
	PermissionIds []uint   `json:"permissions" gorm:"-"`
}

func (a AdminMenuModel) SortAdminMenu(OrderIds []Requests.AdminMenuSortRequestDataOrderIds, Start uint, ParentId uint) (res bool, err error) {
	if ParentId == 0 {
		ParentId = 1
	}
	for _, OrderId := range OrderIds {
		err = database.Db.Model(&a).Where("id = ?", OrderId.Id).Update(map[string]interface{}{"order": Start, "parent_id": ParentId}).Error
		if err != nil {
			break
		}
		Start++
		_, err = a.SortAdminMenu(OrderId.Children, Start, OrderId.Id)
		if err != nil {
			break
		}
	}
	return
}

func (a *AdminMenuModel) Insert() (Id uint64, err error) {
	// 开始事务
	tx := database.Db.Begin()
	if res := tx.Create(&a); res.Error != nil {
		// 遇到错误时回滚事务
		err = res.Error
		tx.Rollback()
		return
	}
	//添加菜单与权限的中间表
	for _, v := range a.PermissionIds {
		if res := tx.Create(&AdminMenuPermissionModel{PermissionId: v, MenuId: a.ID}); res.Error != nil {
			// 遇到错误时回滚事务
			err = res.Error
			tx.Rollback()
			return
		}
	}
	//添加菜单给予超级管理员
	if res := tx.Create(&AdminRoleMenuModel{MenuId: a.ID, RoleId: 1}); res.Error != nil {
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

func (a *AdminMenuModel) Update(Id uint) (Rows int64, err error) {
	// 开始事务
	tx := database.Db.Begin()
	res := tx.Model(&a).Where("id = ?", Id).Updates(&a)
	if res.Error != nil {
		// 遇到错误时回滚事务
		err = res.Error
		tx.Rollback()
		return
	}
	//更新菜单与权限的中间表
	if res := tx.Delete(AdminMenuPermissionModel{}, "menu_id = ?", Id); res.Error != nil {
		// 遇到错误时回滚事务
		err = res.Error
		tx.Rollback()
		return
	}
	for _, v := range a.PermissionIds {
		fmt.Println(v)
		if res := tx.Create(&AdminMenuPermissionModel{PermissionId: v, MenuId: Id}); res.Error != nil {
			// 遇到错误时回滚事务
			err = res.Error
			tx.Rollback()
			return
		}
	}
	Rows = res.RowsAffected
	// 否则，提交事务
	tx.Commit()
	return
}

func (a *AdminMenuModel) Delete(Ids string) (Rows int64, err error) {
	res := database.Db.Where("id in (?)", Ids).Delete(&a)
	if res.Error != nil {
		err = res.Error
		return
	}
	Rows = res.RowsAffected
	return
}

func (a *AdminMenuModel) GetAll(RequestData map[string]interface{}, Offset uint64, Limit uint8) (Data interface{}, err error) {
	panic("implement me")
}

type AdminMenuDetail struct {
	AdminMenuTree `gorm:"embedded"`
	Permissions   []*AdminMenuPermissionList `json:"permissions"`
}

func (a *AdminMenuModel) GetOne(Id uint) (Data interface{}, err error) {
	var (
		AdminUserModel           AdminUserModel
		adminMenuDetail          AdminMenuDetail
		adminMenuPermissionModel AdminMenuPermissionModel
	)
	database.Db.Model(&a).Where("id = ?", Id).Scan(&adminMenuDetail)
	adminMenuDetail.ChildData, err = a.GetMenuTree(Id, AdminUserModel.GetUserMenuIds(2))
	adminMenuDetail.Permissions, err = adminMenuPermissionModel.GetMenuPermissionList(Id)
	Data = adminMenuDetail
	return
}

func (AdminMenuModel) TableName() string {
	return "admin_menu"
}

type AdminMenuTree struct {
	ID        uint             `gorm:"primary_key" json:"id"`
	ParentId  uint             `json:"parent_id"`
	Order     int              `json:"order"`
	Title     string           `json:"title"`
	Icon      string           `json:"icon"`
	Uri       string           `json:"uri"`
	ChildData []*AdminMenuTree `json:"child_data"`
	CreatedAt JsonTime         `json:"created_at"`
	UpdatedAt JsonTime         `json:"updated_at"`
}

func (a AdminMenuModel) GetMenuTree(ParentId uint, MenuIds []uint) (Data []*AdminMenuTree, err error) {
	var adminMenuTree []*AdminMenuTree
	rows, rowsError := database.Db.Model(&a).Where("parent_id = ?", ParentId).Where("id in (?)", MenuIds).Order("`order` asc,id asc").Rows()
	if rowsError != nil {
		err = rowsError
		return
	}
	defer rows.Close()
	for rows.Next() {
		var row AdminMenuTree
		_ = database.Db.ScanRows(rows, &row)
		row.ChildData, _ = a.GetMenuTree(row.ID, MenuIds)
		adminMenuTree = append(adminMenuTree, &row)
	}

	Data = adminMenuTree
	return
}
