package qbot

import "fmt"

// Token 机器人验证信息，在开放平台后台查询
type Token struct {
	AppID    uint64 // 机器人 id（不是 QQ 号）
	BotToken string // 机器人 ws 用的令牌
}

func (t *Token) String() string {
	return fmt.Sprintf("%d.%s", t.AppID, t.BotToken)
}

type LoginUser struct {
	ID       string
	Username string
}

// Session 连接的 session 结构信息
type Session struct {
	ID      string
	URL     string
	Token   Token
	Intent  Intent
	LastSeq int64
	User    LoginUser
}

func (s Session) String() string {
	return fmt.Sprintf("{ID:%s, URL:%s, Intent:%d, LastSeq:%d, User:{ID:%s, Username:%s}}", s.ID, s.URL, s.Intent, s.LastSeq, s.User.ID, s.User.Username)
}
