package models

import (
	"github.com/gookit/goutil/jsonutil"
	"github.com/peixiaochen/BaseAdmin/pkg/database"
	"gorm.io/datatypes"
	"time"
)

type AdminOperationLogModel struct {
	Id        uint        `json:"id"`
	UserId    uint        `json:"user_id"`
	Path      string      `json:"path"`
	Method    string      `json:"method"`
	Ip        string      `json:"ip"`
	Input     interface{} `json:"input"`
	CreatedAt JsonTime    `json:"created_at"`
	UpdatedAt JsonTime    `json:"updated_at"`
}

func (AdminOperationLogModel) TableName() string {
	return "admin_operation_log"
}
func (a *AdminOperationLogModel) Insert() (Id uint64, err error) {
	// 开始事务
	tx := database.Db.Begin()
	extra, _ := jsonutil.Encode(a.Input)
	a.Input = extra
	if res := tx.Create(&a); res.Error != nil {
		// 遇到错误时回滚事务
		err = res.Error
		tx.Rollback()
		return
	}
	Id = uint64(a.Id)
	// 否则，提交事务
	tx.Commit()
	return
}
func (a *AdminOperationLogModel) Update(Id uint) (Rows int64, err error) {
	panic("implement me")
}

func (a *AdminOperationLogModel) Delete(Ids string) (Rows int64, err error) {
	panic("implement me")
}

type AdminOperationLogList struct {
	Id        uint           `json:"id" `
	Name      string         `json:"name"`
	Ip        string         `json:"ip"`
	Method    string         `json:"method"`
	Path      string         `json:"path"`
	UserId    uint           `json:"user_id"`
	Username  string         `json:"username"`
	Input     datatypes.JSON `json:"input"`
	CreatedAt JsonTime       `json:"created_at"`
	UpdatedAt JsonTime       `json:"updated_at"`
}

func (a *AdminOperationLogModel) GetAll(RequestData map[string]interface{}, Offset uint64, Limit uint8) (Data interface{}, err error) {
	var (
		count                 uint64
		adminOperationLogList []*AdminOperationLogList
	)
	model := database.Db.Model(&a)
	if RequestData["Keywords"] != "%%" {
		model = model.
			Where("admin_user.username like ? or admin_user.name like ? or admin_operation_log.path like ? or admin_operation_log.ip like ?", RequestData["Keywords"], RequestData["Keywords"], RequestData["Keywords"], RequestData["Keywords"])
	}
	if RequestData["StartTime"] != "" {
		model = model.Where("created_at >= ?", RequestData["StartTime"])
	}
	if RequestData["EndTime"] != "" {
		model = model.Where("created_at <= ?", RequestData["EndTime"])
	}
	model.Count(&count)
	model = model.Select("admin_operation_log.*,admin_user.name,admin_user.username").Joins("left join admin_user on admin_user.id = admin_operation_log.user_id").Order("id desc").Offset(Offset).Limit(Limit)
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
		var row AdminOperationLogList
		_ = database.Db.ScanRows(rows, &row)
		adminOperationLogList = append(adminOperationLogList, &row)
	}

	Data = map[string]interface{}{"count": count, "data_list": adminOperationLogList}
	return
}

func (a *AdminOperationLogModel) GetOne(Id uint) (Data interface{}, err error) {
	panic("implement me")
}

func (a *AdminOperationLogModel) GetLastLoginIp(UserId uint) (Ip string) {
	database.Db.Where("user_id = ?", UserId).Order("id desc").First(&a)
	Ip = a.Ip
	return
}
