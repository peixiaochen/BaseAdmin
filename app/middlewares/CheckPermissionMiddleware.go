package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/gookit/goutil/strutil"
	"github.com/peixiaochen/BaseAdmin/app/models"
	"github.com/peixiaochen/BaseAdmin/pkg/context"
	"strings"
)

func CheckPermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminUserId, _ := c.Get("AdminUserId")
		var (
			path                 = strutil.Substr(c.Request.URL.Path, 4, len(c.Request.URL.Path))
			adminRoleUserModel   models.AdminRoleUserModel
			roleIds              []uint
			adminPermissionModel models.AdminPermissionModel
		)
		paths := strings.Split(path, "/")
		if len(paths) > 2 {
			path = paths[0] + "/" + paths[1]
		}
		roleIds, _ = adminRoleUserModel.GetRoleIds(adminUserId.(uint))
		permissionInfo, count, err := adminPermissionModel.GetPermissionCheck(path, roleIds)
		if err != nil || count <= 0 {
			var Response context.Response
			Response.Code = context.CodeClientPermissionError
			Response.ServerJson(c)
			c.Abort()
		}
		var model models.Model
		model = &models.AdminOperationLogModel{
			UserId: adminUserId.(uint),
			Path:   c.FullPath(),
			Method: c.Request.Method,
			Ip:     c.ClientIP(),
			Input: map[string]interface{}{
				"data":       c.Params,
				"permission": permissionInfo,
			},
		}
		if _, err := model.Insert(); err != nil {
			var Response context.Response
			Response.Code = context.CodeServerError
			Response.Msg = err.Error()
			Response.ServerJson(c)
			c.Abort()
		} else {
			c.Next()
		}
	}
}
