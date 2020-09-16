package routers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/peixiaochen/BaseAdmin/app/controllers"
	"github.com/peixiaochen/BaseAdmin/app/middlewares"
	"github.com/peixiaochen/BaseAdmin/pkg/config"
	"time"
)

func init() {
	// 创建一个默认的路由引擎
	r := gin.Default()
	store, _ := redis.NewStore(int(60*time.Minute), "tcp", config.RedisSetting.Host, config.RedisSetting.Password, []byte("secret"))
	r.Use(sessions.Sessions("GinBaseAdminSessionId", store))
	v1 := r.Group("/v1")
	{
		v1.POST("/user/login", controllers.Login)
		// 默认使用了2个中间件Logger(), Recovery()
		admin := v1.Group("/", middlewares.CheckLoginMiddleware(), middlewares.CheckPermissionMiddleware()) //需要中间件验证是否登录
		{
			admin.POST("user/logout", controllers.LoginOut)
			admin.POST("welcome", controllers.Welcome)
			admin.POST("my_info/update", controllers.MyInfoUpdate)
			adminMenu := admin.Group("/admin_menu/")
			{
				adminMenu.POST("store", controllers.AdminMenuStore)
				adminMenu.POST("update", controllers.AdminMenuUpdate)
				adminMenu.POST("index", controllers.AdminMenuIndex)
				adminMenu.POST("sort", controllers.AdminMenuSort)
				adminMenu.GET("detail/:Id", controllers.AdminMenuDetail)
			}
			adminUser := admin.Group("/admin_user/")
			{
				adminUser.POST("store", controllers.AdminUserStore)
				adminUser.POST("update", controllers.AdminUserUpdate)
				adminUser.POST("delete", controllers.AdminUserDelete)
				adminUser.POST("status", controllers.AdminUserStatus)
				adminUser.POST("index", controllers.AdminUserIndex)
			}
			adminRole := admin.Group("/admin_role/")
			{
				adminRole.POST("store", controllers.AdminRoleStore)
				adminRole.POST("update", controllers.AdminRoleUpdate)
				adminRole.POST("delete", controllers.AdminRoleDelete)
				adminRole.POST("index", controllers.AdminRoleIndex)
			}
			adminPermission := admin.Group("/admin_permission/")
			{
				adminPermission.POST("store", controllers.AdminPermissionStore)
				adminPermission.POST("update", controllers.AdminPermissionUpdate)
				adminPermission.POST("delete", controllers.AdminPermissionDelete)
				adminPermission.POST("index", controllers.AdminPermissionIndex)
			}
			adminOperation := admin.Group("/admin_operation_log/")
			{
				adminOperation.POST("index", controllers.AdminOperationLogIndex)
			}
		}
	}
	r.Run(fmt.Sprintf(":%d", config.ServerSetting.HttpPort))
}
