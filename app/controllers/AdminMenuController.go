package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/peixiaochen/BaseAdmin/app/Requests"
	"github.com/peixiaochen/BaseAdmin/app/models"
	"github.com/peixiaochen/BaseAdmin/pkg/context"
)

func AdminMenuDetail(c *gin.Context) {
	var (
		Response Response
		Request  *Requests.AdminMenuDetailRequest
		Model    models.Model
	)
	if err := c.ShouldBindUri(&Request); err != nil {
		Response.Msg = err.Error()
		Response.Code = context.CodeClientError
		Response.ServerJson(c)
		return
	}
	Model = &models.AdminMenuModel{}
	if Data, err := Model.GetOne(Request.Id); err != nil {
		Response.Code = context.CodeServerError
		Response.Msg = err.Error()
	} else {
		Response.Data = Data
	}
	Response.ServerJson(c)
	return
}

func AdminMenuIndex(c *gin.Context) {
	var (
		Response       Response
		AdminUserModel models.AdminUserModel
		AdminMenuModel models.AdminMenuModel
	)
	adminUserId, isGet := c.Get("AdminUserId")
	if isGet {
		//获取用户的所有菜单列表
		if Data, err := AdminMenuModel.GetMenuTree(1, AdminUserModel.GetUserMenuIds(adminUserId.(uint))); err != nil {
			Response.Code = context.CodeServerError
			Response.Msg = err.Error()
		} else {
			Response.Data = Data
		}
	}
	Response.ServerJson(c)
	return
}

func AdminMenuSort(c *gin.Context) {
	var (
		Request        *Requests.AdminMenuSortRequest
		Response       Response
		adminMenuModel models.AdminMenuModel
	)
	if err := c.ShouldBindJSON(&Request); err != nil {
		Response.Msg = err.Error()
		Response.Code = context.CodeClientError
		Response.ServerJson(c)
		return
	}
	if _, err := adminMenuModel.SortAdminMenu(Request.Data.OrderIds, 1, 1); err != nil {
		Response.Code = context.CodeServerError
		Response.Msg = err.Error()
	}
	Response.ServerJson(c)
	return
}
func AdminMenuSave(c *gin.Context) {
	var (
		Request  *Requests.AdminMenuSaveRequest
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

	Model = &models.AdminMenuModel{
		ParentId:      Request.Data.ParentId,
		Order:         0,
		Title:         Request.Data.Title,
		Icon:          Request.Data.Icon,
		Uri:           Request.Data.Uri,
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

func AdminMenuStore(c *gin.Context) {
	AdminMenuSave(c)
}
func AdminMenuUpdate(c *gin.Context) {
	AdminMenuSave(c)
}
