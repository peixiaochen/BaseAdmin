package models

import (
	"github.com/gookit/goutil/jsonutil"
	"github.com/jinzhu/gorm"
	"github.com/peixiaochen/BaseAdmin/pkg/database"
	"gorm.io/datatypes"
	"time"
)

type AdminUserModel struct {
	gorm.Model
	//此处就是用空格
	Username string      `json:"username" gorm:"unique;<-:create"`
	Password string      `json:"password" gorm:"<-:update"`
	Name     string      `json:"name"`
	Avatar   string      `json:"avatar"`
	Status   uint8       `json:"status"`
	RolesIds []uint      `json:"roles" gorm:"-"`
	Extra    interface{} `json:"extra"`
}

func (AdminUserModel) TableName() string {
	return "admin_user"
}
func (a *AdminUserModel) Insert() (Id uint64, err error) {
	// 开始事务
	tx := database.Db.Begin()
	extra, _ := jsonutil.Encode(a.Extra)
	a.Extra = string(extra)
	if res := tx.Create(&a); res.Error != nil {
		// 遇到错误时回滚事务
		err = res.Error
		tx.Rollback()
		return
	}
	//添加用户与角色的中间表
	for _, v := range a.RolesIds {
		if res := tx.Create(&AdminRoleUserModel{RoleId: v, UserId: a.ID}); res.Error != nil {
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

func (a *AdminUserModel) Update(Id uint) (Rows int64, err error) {
	// 开始事务
	tx := database.Db.Begin()
	extra, _ := jsonutil.Encode(a.Extra)
	a.Extra = string(extra)
	res := tx.Model(&a).Where("id = ?", Id).Updates(&a)
	if res.Error != nil {
		// 遇到错误时回滚事务
		err = res.Error
		tx.Rollback()
		return
	}
	Rows = res.RowsAffected
	//添加用户与角色的中间表
	if res := tx.Delete(AdminRoleUserModel{}, "user_id = ?", Id); res.Error != nil {
		// 遇到错误时回滚事务
		err = res.Error
		tx.Rollback()
		return
	}
	for _, v := range a.RolesIds {
		if res := tx.Create(&AdminRoleUserModel{RoleId: v, UserId: Id}); res.Error != nil {
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

func (a *AdminUserModel) MyInfoUpdate(Id uint) (Rows int64, err error) {
	// 开始事务
	tx := database.Db.Begin()
	extra, _ := jsonutil.Encode(a.Extra)
	a.Extra = string(extra)
	res := tx.Model(&a).Where("id = ?", Id).Updates(&a)
	if res.Error != nil {
		// 遇到错误时回滚事务
		err = res.Error
		tx.Rollback()
		return
	}
	Rows = res.RowsAffected
	tx.Commit()
	return
}

func (a *AdminUserModel) Delete(Ids string) (Rows int64, err error) {
	res := database.Db.Where("id in (?)", Ids).Delete(&a)
	if res.Error != nil {
		err = res.Error
		return
	}
	Rows = res.RowsAffected
	return
}

func (a *AdminUserModel) SetStatus(Ids string, Status uint8) (Rows int64, err error) {
	if Status == 0 || Status == 1 {
		res := database.Db.Model(&a).Where("id IN (?)", Ids).Update("status", Status)
		if res.Error != nil {
			err = res.Error
			return
		}
		Rows = res.RowsAffected
	}
	return
}

type AdminUserList struct {
	Id          uint                 `json:"id" `
	Username    string               `json:"username"`
	Name        string               `json:"name"`
	Avatar      string               `json:"avatar"`
	Status      uint8                `json:"status"`
	Extra       datatypes.JSON       `json:"extra"`
	CreatedAt   JsonTime             `json:"created_at"`
	UpdatedAt   JsonTime             `json:"updated_at"`
	Roles       []*AdminRoleUserList `json:"roles"`
	LastLoginIp string               `json:"last_login_ip"`
}

func (a *AdminUserModel) GetAll(RequestData map[string]interface{}, Offset uint64, Limit uint8) (Data interface{}, err error) {

	var (
		count                  uint64
		adminUserList          []*AdminUserList
		AdminRoleUserModel     AdminRoleUserModel
		AdminOperationLogModel AdminOperationLogModel
	)
	model := database.Db.Model(&a).Where("id>1")
	if RequestData["Keywords"] != "%%" {
		model = model.Where("username like ? or name like ?", RequestData["Keywords"], RequestData["Keywords"])
	}
	if RequestData["RoleId"].(uint64) > 0 {
		model = model.Joins("left join admin_role_user  on admin_role_user.user_id = admin_user.id").Where("admin_role_user.role_id = ?", RequestData["RoleId"])
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
		var row AdminUserList
		_ = database.Db.ScanRows(rows, &row)
		row.LastLoginIp = AdminOperationLogModel.GetLastLoginIp(row.Id)
		row.Roles, _ = AdminRoleUserModel.GetUserRoleList(row.Id)
		adminUserList = append(adminUserList, &row)
	}

	Data = map[string]interface{}{"count": count, "data_list": adminUserList}
	return
}

func (a *AdminUserModel) GetOne(Id uint) (Data interface{}, err error) {
	var (
		adminUserList      AdminUserList
		AdminRoleUserModel AdminRoleUserModel
	)
	database.Db.Model(&a).Where("id = ?", Id).Scan(&adminUserList)
	adminUserList.Roles, _ = AdminRoleUserModel.GetUserRoleList(Id)
	Data = adminUserList
	return
}
func (a *AdminUserModel) GetIdByName(Username string) (Id uint, err error) {
	if err = database.Db.Where("username = ?", Username).First(&a).Error; err != nil {
		return
	}
	Id = a.ID
	return
}
func (a AdminUserModel) GetUserMenuIds(UserId uint) (MenuIds []uint) {
	database.Db.Table("admin_role_menu as a_rm").
		Select("a_rm.menu_id").
		Joins("left join admin_role_user as a_ru on a_ru.role_id = a_rm.role_id").
		Where("a_ru.user_id = ?", UserId).Pluck("a_rm.menu_id", &MenuIds)
	return
}
