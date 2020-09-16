package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/peixiaochen/BaseAdmin/app/Requests"
	"github.com/peixiaochen/BaseAdmin/app/models"
	"github.com/peixiaochen/BaseAdmin/pkg/context"
)

func AdminPermissionDelete(c *gin.Context) {
	var (
		Request  *Requests.AdminPermissionDeleteRequest
		Response Response
		Model    models.Model
	)
	if err := c.ShouldBindJSON(&Request); err != nil {
		Response.Msg = err.Error()
		Response.Code = context.CodeClientError
		Response.ServerJson(c)
		return
	}
	Model = &models.AdminPermissionModel{}
	if _, err := Model.Delete(Request.Data.Ids); err != nil {
		Response.Code = context.CodeServerError
		Response.Msg = err.Error()
	}
	Response.ServerJson(c)
	return
}
func AdminPermissionIndex(c *gin.Context) {
	var (
		Response Response
		Request  *Requests.AdminPermissionIndexRequest
		Model    models.Model
	)
	if err := c.ShouldBindJSON(&Request); err != nil {
		Response.Msg = err.Error()
		Response.Code = context.CodeClientError
		Response.ServerJson(c)
		return
	}
	Model = &models.AdminPermissionModel{}
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

func AdminPermissionSave(c *gin.Context) {
	var (
		Request  *Requests.AdminPermissionSaveRequest
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

	Model = &models.AdminPermissionModel{
		Name:        Request.Data.Name,
		Description: Request.Data.Description,
		HttpMethod:  Request.Data.HttpMethod,
		HttpPath:    Request.Data.HttpPath,
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

func AdminPermissionStore(c *gin.Context) {
	AdminPermissionSave(c)
}
func AdminPermissionUpdate(c *gin.Context) {
	AdminPermissionSave(c)
}
