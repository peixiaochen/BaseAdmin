package BaseAdmin

import (
	_ "github.com/peixiaochen/BaseAdmin/pkg/config"
	orm "github.com/peixiaochen/BaseAdmin/pkg/database"
	_ "github.com/peixiaochen/BaseAdmin/routers"
)

func init() {
	orm.CloseDB()
}
