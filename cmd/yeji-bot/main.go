package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"yeji-bot/bot/qbot"
	"yeji-bot/bot/scheduler"
	"yeji-bot/config"
	"yeji-bot/dao"
	"yeji-bot/msgservice"
)

func main() {
	cfgPath := flag.String("c", "config.ini", "配置文件路径")
	flag.Parse()

	err := config.LoadConfig(*cfgPath)
	if err != nil {
		logrus.Errorf("加载配置文件失败 %v", err)
		return
	}
	logrus.Infof("加载配置文件成功")

	err = dao.InitDao()
	if err != nil {
		logrus.Errorf("初始化数据库连接失败 %v", err)
		return
	}
	logrus.Infof("初始化数据库连接成功")

	msgScheduler := scheduler.NewGroupAtMessageHandlerScheduler()
	msgservice.Register(msgScheduler)

	client, err := qbot.New(
		config.AppConfig.Bot.AppID,
		config.AppConfig.Bot.BotToken,
		config.AppConfig.Bot.AppSecret,
		0|qbot.IntentGroupAtMessage,
		config.AppConfig.Bot.Sandbox,
	)
	if err != nil {
		logrus.Errorf("初始化机器人配置失败 %v", err)
		return
	}
	client.RegisterHandlerScheduler(msgScheduler)
	err = client.Start()
	if err != nil {
		logrus.Errorf("启动机器人失败 %v", err)
		return
	}

	return
}
