package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/peixiaochen/BaseAdmin/app/Requests"
	"github.com/peixiaochen/BaseAdmin/app/models"
	"github.com/peixiaochen/BaseAdmin/pkg/config"
	"github.com/peixiaochen/BaseAdmin/pkg/context"
	"github.com/peixiaochen/BaseAdmin/pkg/lib"
)

func AdminUserIndex(c *gin.Context) {
	var (
		Response Response
		Request  *Requests.AdminUserIndexRequest
		Model    models.Model
	)
	if err := c.ShouldBindJSON(&Request); err != nil {
		Response.Msg = err.Error()
		Response.Code = context.CodeClientError
		Response.ServerJson(c)
		return
	}
	Model = &models.AdminUserModel{}
	if Data, err := Model.GetAll(map[string]interface{}{
		"Keywords":  "%" + Request.Data.Keywords + "%",
		"RoleId":    Request.Data.RoleId,
		"StartTime": Request.Data.StartTime,
		"EndTime":   Request.Data.EndTime,
	}, (Request.Data.Page-1)*uint64(Request.Data.Limit), Request.Data.Limit); err != nil {
		Response.Code = context.CodeServerError
		Response.Msg = err.Error()
	} else {
		Response.Data = Data
	}
	Response.ServerJson(c)
	return
}

func AdminUserDelete(c *gin.Context) {
	var (
		Request  *Requests.AdminUserDeleteRequest
		Response Response
		Model    models.Model
	)
	if err := c.ShouldBindJSON(&Request); err != nil {
		Response.Msg = err.Error()
		Response.Code = context.CodeClientError
		Response.ServerJson(c)
		return
	}
	Model = &models.AdminUserModel{}
	if _, err := Model.Delete(Request.Data.Ids); err != nil {
		Response.Code = context.CodeServerError
		Response.Msg = err.Error()
	}
	Response.ServerJson(c)
	return
}

func AdminUserStatus(c *gin.Context) {
	var (
		Response Response
		Request  *Requests.AdminUserStatusRequest
		Model    models.AdminUserModel
	)
	if err := c.ShouldBindJSON(&Request); err != nil {
		Response.Msg = err.Error()
		Response.Code = context.CodeClientError
		Response.ServerJson(c)
		return
	}
	if _, err := Model.SetStatus(Request.Data.Ids, Request.Data.Status); err != nil {
		Response.Code = context.CodeServerError
		Response.Msg = err.Error()
	}
	Response.ServerJson(c)
	return
}

func AdminUserSave(c *gin.Context) {
	var (
		Request  *Requests.AdminUserSaveRequest
		Response Response
		Model    models.Model
	)
	if err := c.ShouldBindJSON(&Request); err != nil {
		Response.Msg = err.Error()
		Response.Code = context.CodeClientError
		Response.ServerJson(c)
		return
	}
	if res, msg := ValidateData(Request.Data, Request.Data.Extra); res == false {
		Response.Msg = msg
		Response.Code = context.CodeClientError
		Response.ServerJson(c)
		return
	}
	if Request.Data.Password != "" {
		Request.Data.Password = lib.SubstrMd5(Request.Data.Password, config.PasswordSetting.PasswordStart, config.PasswordSetting.PasswordLength)
	}
	var status uint8 = 1
	if Request.Data.Id > 0 {
		status = Request.Data.Status
	}
	Model = &models.AdminUserModel{
		Username: Request.Data.UserName,
		Password: Request.Data.Password,
		Name:     Request.Data.Name,
		Avatar:   Request.Data.Avatar,
		RolesIds: Request.Data.Roles,
		Status:   status,
		Extra: map[string]interface{}{
			"phone": Request.Data.Extra.Phone,
			"email": Request.Data.Extra.Email,
			"sex":   Request.Data.Extra.Sex,
		},
	}
	if Request.Data.Id == 0 {
		if _, err := Model.Insert(); err != nil {
			Response.Code = context.CodeServerError
			Response.Msg = err.Error()
		}
	} else {
		if _, err := Model.Update(Request.Data.Id); err != nil {
			Response.Code = context.CodeServerError
			Response.Msg = err.Error()
		}
	}
	Response.ServerJson(c)
	return
}

func AdminUserStore(c *gin.Context) {
	AdminUserSave(c)
}

func AdminUserUpdate(c *gin.Context) {
	AdminUserSave(c)
}

func LoginOut(c *gin.Context) {
	var (
		Response Response
	)
	session := sessions.Default(c)
	session.Delete("AdminUserId")
	_ = session.Save()
	Response.Msg = "user logout success"
	Response.ServerJson(c)
}

func Login(c *gin.Context) {
	var (
		Request  *Requests.AdminUserLoginRequest
		Response Response
		Model    models.AdminUserModel
	)
	if err := c.ShouldBindJSON(&Request); err != nil {
		Response.Msg = err.Error()
		Response.Code = context.CodeClientError
		Response.ServerJson(c)
		return
	}
	if res, msg := ValidateData(Request.Data); res == false {
		Response.Msg = msg
		Response.Code = context.CodeClientError
		Response.ServerJson(c)
		return
	}
	session := sessions.Default(c)
	id, err := Model.GetIdByName(Request.Data.UserName)
	if err != nil {
		Response.Msg = err.Error()
		Response.Code = context.CodeServerError
		Response.ServerJson(c)
		return
	}
	session.Set("AdminUserId", id)
	session.Save()
	Response.ServerJson(c)
}

func MyInfoUpdate(c *gin.Context) {
	var (
		Request  *Requests.AdminUserSaveRequest
		Response Response
	)
	adminUserId, _ := c.Get("AdminUserId")
	if err := c.ShouldBindJSON(&Request); err != nil {
		Response.Msg = err.Error()
		Response.Code = context.CodeClientError
		Response.ServerJson(c)
		return
	}
	if res, msg := ValidateData(Request.Data, Request.Data.Extra); res == false {
		Response.Msg = msg
		Response.Code = context.CodeClientError
		Response.ServerJson(c)
		return
	}
	if Request.Data.Password != "" {
		Request.Data.Password = lib.SubstrMd5(Request.Data.Password, config.PasswordSetting.PasswordStart, config.PasswordSetting.PasswordLength)
	}
	Model := &models.AdminUserModel{
		Username: Request.Data.UserName,
		Password: Request.Data.Password,
		Name:     Request.Data.Name,
		Avatar:   Request.Data.Avatar,
		RolesIds: Request.Data.Roles,
		Extra: map[string]interface{}{
			"phone": Request.Data.Extra.Phone,
			"email": Request.Data.Extra.Email,
			"sex":   Request.Data.Extra.Sex,
		},
	}
	if _, err := Model.MyInfoUpdate(adminUserId.(uint)); err != nil {
		Response.Code = context.CodeServerError
		Response.Msg = err.Error()
	}
	Response.ServerJson(c)
	return
	Response.ServerJson(c)
}
