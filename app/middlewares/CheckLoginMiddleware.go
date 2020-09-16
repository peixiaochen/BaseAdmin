package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/peixiaochen/BaseAdmin/pkg/context"
)

func CheckLoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		AdminUserId := session.Get("AdminUserId")
		if AdminUserId == nil {
			var Response context.Response
			Response.Code = context.CodeClientNoLogin
			Response.ServerJson(c)
			c.Abort()
		} else {
			c.Set("AdminUserId", AdminUserId)
			c.Next()
		}
	}
}
