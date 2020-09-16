package config

import (
	"gopkg.in/ini.v1"
	"log"
	"time"
)

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

type Password struct {
	PasswordStart  int
	PasswordLength int
}

var PasswordSetting = &Password{}

var cfg *ini.File

// Setup initialize the configuration instance
func init() {
	var err error
	cfg, err = ini.Load("config/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'config/app.ini': %v", err)
	}

	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)
	mapTo("password", PasswordSetting)

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
