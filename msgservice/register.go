package msgservice

import (
	"yeji-bot/bot/scheduler"
	"yeji-bot/middleware"
)

func Register(s *scheduler.GroupAtMessageHandlerScheduler) {
	s.Register("/帮助", IntroHelp)
	s.Register("/更新", UpdateData).Middlewares(middleware.MustAdmin())

	s.Register("/厨师", QueryChef).Middlewares(middleware.Helper(chefHelperStr()))
	s.Register("/菜谱", QueryRecipe).Middlewares(middleware.Helper(recipeHelperStr()))
	s.Register("/厨具", QueryEquip).Middlewares(middleware.Helper(equipHelperStr()))
	s.Register("/食材", QueryMaterial).Middlewares(middleware.Helper(materialHelperStr()))
	s.Register("/贵客", QueryGuest).Middlewares(middleware.Helper(guestHelperStr()))
	s.Register("/符文", QueryAntique).Middlewares(middleware.Helper(antiqueHelperStr()))
	s.Register("/任务", QueryQuest).Middlewares(middleware.Helper(questHelperStr()))
	s.Register("/抽签", QueryTarot)
}
