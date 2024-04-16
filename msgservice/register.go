package msgservice

import "yeji-bot/bot/scheduler"

func Register(s *scheduler.GroupAtMessageHandlerScheduler) {
	s.Register("/帮助", introHelp)
}
