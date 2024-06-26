package scheduler

import (
	"context"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
	"yeji-bot/bot/openapi"
	"yeji-bot/bot/qbot"
	"yeji-bot/pkg/kit"
)

type GroupAtMessageHandlerFunc func(ctx context.Context, api *openapi.Openapi, msg *qbot.WSGroupAtMessageData) (err error)

type GroupAtMessageHandlerMiddleware func(handler GroupAtMessageHandlerFunc) GroupAtMessageHandlerFunc

func chain(handlers ...GroupAtMessageHandlerMiddleware) GroupAtMessageHandlerMiddleware {
	return func(handler GroupAtMessageHandlerFunc) GroupAtMessageHandlerFunc {
		for i := len(handlers) - 1; i >= 0; i-- {
			handler = handlers[i](handler)
		}
		return handler
	}
}

type CmdHandler struct {
	cmd         string
	middlewares []GroupAtMessageHandlerMiddleware
	handler     GroupAtMessageHandlerFunc
}

func (ch *CmdHandler) Middlewares(middlewares ...GroupAtMessageHandlerMiddleware) {
	ch.middlewares = append([]GroupAtMessageHandlerMiddleware{}, middlewares...)
}

type GroupAtMessageHandlerScheduler struct {
	globalMiddlewares []GroupAtMessageHandlerMiddleware
	defaultHandler    *CmdHandler
	cmdHandlers       []*CmdHandler
}

func NewGroupAtMessageHandlerScheduler() *GroupAtMessageHandlerScheduler {
	s := &GroupAtMessageHandlerScheduler{
		globalMiddlewares: make([]GroupAtMessageHandlerMiddleware, 0),
		defaultHandler: &CmdHandler{
			cmd:         "",
			middlewares: make([]GroupAtMessageHandlerMiddleware, 0),
			handler: func(ctx context.Context, api *openapi.Openapi, msg *qbot.WSGroupAtMessageData) (err error) {
				return nil
			},
		},
		cmdHandlers: make([]*CmdHandler, 0),
	}
	return s
}

func (s *GroupAtMessageHandlerScheduler) GlobalMiddlewares(middlewares ...GroupAtMessageHandlerMiddleware) {
	s.globalMiddlewares = middlewares
}

func (s *GroupAtMessageHandlerScheduler) DefaultHandler(handler GroupAtMessageHandlerFunc) *CmdHandler {
	var defaultCmdHandler = &CmdHandler{
		cmd:         "",
		middlewares: make([]GroupAtMessageHandlerMiddleware, 0),
		handler:     handler,
	}
	s.defaultHandler = defaultCmdHandler
	return defaultCmdHandler
}

func (s *GroupAtMessageHandlerScheduler) Register(cmd string, handler GroupAtMessageHandlerFunc) *CmdHandler {
	cmdHandler := &CmdHandler{
		cmd:         cmd,
		middlewares: make([]GroupAtMessageHandlerMiddleware, 0),
		handler:     handler,
	}
	s.cmdHandlers = append(s.cmdHandlers, cmdHandler)
	return cmdHandler
}

func (s *GroupAtMessageHandlerScheduler) match(cmd string) *CmdHandler {
	for _, cmdHandler := range s.cmdHandlers {
		if strings.HasPrefix(cmd, cmdHandler.cmd) {
			return cmdHandler
		}
	}
	return s.defaultHandler
}

func (s *GroupAtMessageHandlerScheduler) Handler() qbot.GroupAtMessageHandler {
	return func(api *openapi.Openapi, event *qbot.WSPayload, msg *qbot.WSGroupAtMessageData) (err error) {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					buf := make([]byte, 4096)
					buf = buf[:runtime.Stack(buf, false)]
					logrus.Errorf("[GroupAtMessageHandlerScheduler.Handler] panic. err: %v, event: %+v, msg: %+v, stack: %s", e, *event, *msg, string(buf))
				}
			}()
			msg.Content = strings.TrimSpace(msg.Content)
			cmdHandler := s.match(msg.Content)
			msg.CmdName = cmdHandler.cmd
			msg.Content = strings.TrimPrefix(msg.Content, cmdHandler.cmd)
			msg.Content = strings.TrimSpace(msg.Content)
			ms := make([]GroupAtMessageHandlerMiddleware, 0, len(s.globalMiddlewares)+len(cmdHandler.middlewares))
			ms = append(append(ms, s.globalMiddlewares...), cmdHandler.middlewares...)
			chainHandler := chain(ms...)(cmdHandler.handler)
			ctx := context.Background()
			ctx = kit.SetSeq(ctx, 0)
			err = chainHandler(ctx, api, msg)
			if err != nil {
				logrus.Infof("[GroupAtMessageHandlerScheduler.Handler] do chainHandler() fail. err: %v", err)
			}
		}()

		return nil
	}

}
