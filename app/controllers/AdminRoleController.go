package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/peixiaochen/BaseAdmin/app/Requests"
	"github.com/peixiaochen/BaseAdmin/app/models"
	"github.com/peixiaochen/BaseAdmin/pkg/context"
)

func AdminRoleIndex(c *gin.Context) {
	var (
		Response Response
		Request  *Requests.AdminRoleIndexRequest
		Model    models.Model
	)
	if err := c.ShouldBindJSON(&Request); err != nil {
		Response.Msg = err.Error()
		Response.Code = context.CodeClientError
		Response.ServerJson(c)
		return
	}
	Model = &models.AdminRoleModel{}
	if Data, err := Model.GetAll(map[string]interface{}{
		"Keywords":  "%" + Request.Data.Keywords + "%",
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

func AdminRoleDelete(c *gin.Context) {
	var (
		Request  *Requests.AdminRoleDeleteRequest
		Response Response
		Model    models.Model
	)
	if err := c.ShouldBindJSON(&Request); err != nil {
		Response.Msg = err.Error()
		Response.Code = context.CodeClientError
		Response.ServerJson(c)
		return
	}
	Model = &models.AdminRoleModel{}
	if _, err := Model.Delete(Request.Data.Ids); err != nil {
		Response.Code = context.CodeServerError
		Response.Msg = err.Error()
	}
	Response.ServerJson(c)
	return
}

func AdminRoleSave(c *gin.Context) {
	var (
		Request  *Requests.AdminRoleSaveRequest
		Response Response
		Model    models.Model
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

	Model = &models.AdminRoleModel{
		Name:          Request.Data.Name,
		Description:   Request.Data.Description,
		Status:        Request.Data.Status,
		MenusIds:      Request.Data.Menus,
		PermissionIds: Request.Data.Permissions,
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
func AdminRoleStore(c *gin.Context) {
	AdminRoleSave(c)
}
func AdminRoleUpdate(c *gin.Context) {
	AdminRoleSave(c)
}
