package config

import (
	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
)

type botConfig struct {
	AppID       uint64   `ini:"app_id"`       // 腾讯开放平台后台获取的机器人 app_id
	BotToken    string   `ini:"bot_token"`    // 腾讯开放平台后台获取的 bot_token
	AppSecret   string   `ini:"app_secret"`   // 腾讯开放平台后台获取的 client_secret
	Sandbox     bool     `ini:"sandbox"`      // 是否是沙盒环境
	AdminGroups []string `ini:"admin_groups"` // 管理员可以处理管理功能的群组的 openid 列表
	AdminUsers  []string `ini:"admin_users"`  // 管理员的 openid 列表
}

type dbConfig struct {
	DSN string `ini:"dsn"`
}

var AppConfig struct {
	Bot botConfig `ini:"bot"`
	DB  dbConfig  `ini:"db"`
}

func LoadConfig(path string) error {
	cfg, err := ini.Load(path)
	if err != nil {
		return errors.Wrap(err, "ini.Load fail")
	}
	err = cfg.MapTo(&AppConfig)
	if err != nil {
		return errors.Wrap(err, "cfg.MapTo fail")
	}
	return nil
}
