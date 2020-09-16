package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
	"github.com/peixiaochen/BaseAdmin/app/models"
	"github.com/peixiaochen/BaseAdmin/pkg/context"
)

type Response struct {
	context.Response
}

func Welcome(c *gin.Context) {
	var (
		Response       Response
		adminUserModel models.AdminUserModel
	)
	adminUserId, _ := c.Get("AdminUserId")
	myInfo, _ := adminUserModel.GetOne(adminUserId.(uint))
	Response.Data = map[string]interface{}{
		"my_info": myInfo,
	}
	Response.ServerJson(c)
}

func ValidateData(data ...interface{}) (res bool, msg string) {
	for _, i := range data {
		v := validate.Struct(i)
		if !v.Validate() { // validate ok
			msg = v.Errors.One()
			return
		}
	}
	return true, ""
}
