# BaseAdmin
后台框架的基础模板，集成权限、菜单、角色、用户的基本增删改查、登陆与退出功能等

- 支持的版本Go1.14+
- 依赖环境Mysql、Redis
- Go版本支持mod

# 安装
```go
    go get -u github.com/peixiaochen/BaseAdmin
```
# 快速开始

## 新建 *config/app.ini* 目录 
在自己的项目目录下新建 config/app.ini
将下列代码粘贴到自己平级目录中，修改Mysql与Redis的配置即可

```go
# possible values : production, development
app_mode = development

[server]
#debug or release
RunMode = debug
HttpPort = 8080
ReadTimeout = 60
WriteTimeout = 60

[database]
Type = mysql
User = root
Password = root
Host = 127.0.0.1:3306
Name = base_admin
TablePrefix =

[redis]
Host = 127.0.0.1:6379
Password =

[password]
PasswordStart = 2
PasswordLength = 18
```

|Key|描述|配置|
|:-------|:----|:----|
|app_mode|gin自带配置||
|RunMode|gin自带配置||
|HttpPort| 项目运行端口||
|ReadTimeout | 读取超时时间 ||
|WriteTimeout| 写入超时时间||
|Mysql配置|----|----|
|Type|连接方式|mysql|
|User|用户名|right|
|Password|密码|right|
|Host|数据库地址:端口|right|
|Name|数据库名|right|
|TablePrefix|数据库前缀|目前可不填，暂无用|
|Redis配置|----|----|
|Host|redis地址|right|
|Password|连接密码|right|
|后台密码加密规则配置|----|----|
|PasswordStart| 密码开始截取处| |
|PasswordLength|密码开始截取长度| |

## 运行sql文件至mysql

下载的目录中携带一个 **BaseAdmin.sql** 文件，运行后生成mysql的基础数据库结构与基础数据

## 运行demo

```go
package main

import (
	"fmt"
	_ "github.com/peixiaochen/BaseAdmin"
)

func main() {
	fmt.Println("hello world")
	//业务逻辑
}


```

## 使用middleware验证权限


1. 导入`"github.com/peixiaochen/BaseAdmin/app/middlewares"`包
2. 使用 `middlewares.CheckLoginMiddleware(), middlewares.CheckPermissionMiddleware()`
3. 注意：如果不使用这两个中间件验证路由，或者自己书写验证逻辑，要不然后台的页面有翻墙使用的漏洞！！！
