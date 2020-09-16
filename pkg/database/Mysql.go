package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/peixiaochen/BaseAdmin/pkg/config"
	"log"
	"time"
)

var Db *gorm.DB

//type Mysql struct {
//	ID        uint64    `gorm:"primary_key" json:"id"`
//	CreatedAt time.Time `json:"created_at"`
//	UpdatedAt time.Time `json:"updated_at"`
//}

func init() {
	var err error
	Db, err = gorm.Open(config.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DatabaseSetting.User,
		config.DatabaseSetting.Password,
		config.DatabaseSetting.Host,
		config.DatabaseSetting.Name))

	if err != nil {
		log.Println(err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return config.DatabaseSetting.TablePrefix + defaultTableName
	}

	Db.SingularTable(true)
	Db.LogMode(true)
	//空闲
	Db.DB().SetMaxIdleConns(50)
	//打开
	Db.DB().SetMaxOpenConns(100)
	//超时
	Db.DB().SetConnMaxLifetime(time.Second * 30)
}

func CloseDB() {
	defer Db.Close()
}
