package msgservice

import (
	"yeji-bot/bot/scheduler"
	"yeji-bot/middleware"
)

func Register(s *scheduler.GroupAtMessageHandlerScheduler) {
	s.Register("/帮助", IntroHelp)
	s.Register("/更新", UpdateData).Middlewares(middleware.MustAdmin())

	s.Register("/厨师", QueryChef)
}
