package qbot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"net/http"
	"reflect"
	"runtime"
	"time"
	"yeji-bot/bot/openapi"
)

type Client struct {
	conn       *websocket.Conn
	httpClient *http.Client
	session    *Session
	api        *openapi.Openapi

	eventHandlers *EventHandlers

	messageQueue    chan *WSPayload
	closeErrChan    chan error
	heartBeatTicker *time.Ticker
}

func New(appID uint64, token string, clientSecret string, intent Intent, sandbox bool) (*Client, error) {
	api, err := openapi.New(appID, clientSecret, sandbox)
	if err != nil {
		return nil, errors.Wrap(err, "new openapi fail")
	}
	apResp, err := api.WS(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "get ws ap fail")
	}
	session := &Session{
		ID:  "",
		URL: apResp.URL,
		Token: Token{
			AppID:    appID,
			BotToken: token,
		},
		Intent:  intent,
		LastSeq: 0,
		User: LoginUser{
			ID:       "",
			Username: "",
		},
	}
	return &Client{
		httpClient: &http.Client{},
		session:    session,
		api:        api,

		eventHandlers:   &EventHandlers{},
		messageQueue:    make(chan *WSPayload, 1000),
		closeErrChan:    make(chan error, 10),
		heartBeatTicker: time.NewTicker(10 * time.Second),
	}, nil
}

func (c *Client) Start() (err error) {
	reconnectCnt := 0
	for reconnectCnt < 5 {
		if reconnectCnt != 0 {
			time.Sleep(3 * time.Second) // 3秒内最多尝试重连一次
		}
		reconnectCnt++
		err = c.connect()
		if err != nil {
			logrus.Errorf("try connect fail. err: %v", err)
			continue
		}
		err = c.Identify()
		if err != nil {
			logrus.Errorf("try identify fail. err: %v", err)
			continue
		}
		reconnectCnt = 0
		err = c.listening()
		if err != nil {
			logrus.Errorf("keep listening fail. err: %v", err)
			continue
		}
	}
	return errors.New("try reconnect fail too much count")
}

func (c *Client) connect() (err error) {
	if c.session.URL == "" {
		return errors.New("invalid ws ap url")
	}
	c.conn, _, err = websocket.DefaultDialer.Dial(c.session.URL, nil)
	if err != nil {
		return errors.Wrap(err, "dial err")
	}
	logrus.Infof("websocket dial success. url: %s, intent: %d", c.session.URL, c.session.Intent)
	return nil
}

func (c *Client) listening() error {
	go c.recvMessage()

	for {
		select {
		case <-c.heartBeatTicker.C:
			logrus.Infof("heart beat. session: %+v", *c.session)
			err := c.Write(&WSPayload{
				WSPayloadBase: WSPayloadBase{
					OPCode: OPCodeHeartbeat,
				},
				Data: c.session.LastSeq,
			})
			if err != nil {
				logrus.Errorf("send heart beat fail. err: %v", err)
				return errors.New("send heart beat fail")
			}
		case err := <-c.closeErrChan:
			logrus.Errorf("close err chan. err: %v", err)
			return err
		}
	}
}

// ParseData 解析数据
func ParseData(message []byte, target interface{}) error {
	data := gjson.Get(string(message), "d")
	return json.Unmarshal([]byte(data.String()), target)
}

func (c *Client) recvMessage() {
	defer func() {
		if e := recover(); e != nil {
			buf := make([]byte, 2048)
			buf = buf[:runtime.Stack(buf, false)]
			logrus.Errorf("[recvMessage] panic. session: %+v. err: %v, stack: %s", *c.session, e, buf)
			c.closeErrChan <- fmt.Errorf("recvMessage panic: %v", e)
			return
		}
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			logrus.Errorf("read message fail. err:%v", err)
			c.closeErrChan <- err
			return
		}
		payload := &WSPayload{}
		err = json.Unmarshal(message, payload)
		if err != nil {
			logrus.Errorf("json unmarshal payload fail. err: %v, raw message: %s", err, string(message))
			continue
		}
		payload.RawMessage = message
		if payload.Seq > 0 {
			c.session.LastSeq = payload.Seq
		}
		logrus.Infof("recv data. op: %d %s. raw message: %s", payload.OPCode, OPMeans(payload.OPCode), string(message))

		go func() {
			defer func() {
				if e := recover(); e != nil {
					buf := make([]byte, 2048)
					buf = buf[:runtime.Stack(buf, false)]
					logrus.Errorf("[recvMessage.go] panic. session: %+v. payload: %+v. err: %v, stack: %s", *c.session, payload, e, buf)
					c.closeErrChan <- fmt.Errorf("parseAndHandle payload panic: %v", e)
					return
				}
			}()
			err = c.parseAndHandle(payload)
			if err != nil {
				logrus.Errorf("parseAndHandle fail. err: %v", err)
			}
		}()
	}

}

func (c *Client) Write(message *WSPayload) error {
	b, err := json.Marshal(message)
	if err != nil {
		return errors.Wrapf(err, "json marshal message fail. message: %+v", *message)
	}
	err = c.conn.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return err
	}
	return nil
}

// Identify 进行鉴权
func (c *Client) Identify() error {
	payload := &WSPayload{
		WSPayloadBase: WSPayloadBase{
			OPCode: OPCodeIdentity,
		},
		Data: WSIdentityData{
			Token:   c.session.Token.String(),
			Intents: c.session.Intent,
			Shard:   []uint32{0, 1}, // 不分片
		},
	}
	return c.Write(payload)
}

func (c *Client) RegisterHandlerScheduler(schedulers ...interface{}) {
	for _, s := range schedulers {
		switch s.(type) {
		case IDefaultHandlerScheduler:
			c.eventHandlers.DefaultHandler = s.(IDefaultHandlerScheduler).Handler()
			logrus.Info("register DefaultHandlerScheduler success")
		case IReadyHandlerScheduler:
			c.eventHandlers.ReadyHandler = s.(IReadyHandlerScheduler).Handler()
			logrus.Info("register ReadyHandlerScheduler success")
		case IGroupAtMessageHandlerScheduler:
			logrus.Info("register GroupAtMessageHandlerScheduler success")
			c.eventHandlers.GroupAtMessageHandler = s.(IGroupAtMessageHandlerScheduler).Handler()
		default:
			logrus.Warnf("register unknown handler schcheduler :%s", reflect.TypeOf(s).Name())
		}
	}
}

type eventParseFunc func(*Client, *WSPayload) error

var eventParseFuncMap = map[OPCode]map[EventType]eventParseFunc{
	OPCodeHello: {
		EventTypeNone: (*Client).helloHandler,
	},
	OPCodeReconnect: {
		EventTypeNone: (*Client).reconnectHandler,
	},
	OPCodeInvalidSession: {
		EventTypeNone: (*Client).invalidSessionHandler,
	},
	OPCodeHeartbeatAck: {
		EventTypeNone: (*Client).heartbeatAckHandler,
	},
	OPCodeDispatchEvent: {
		EventReady:                (*Client).readyHandler,
		EventGroupAtMessageCreate: (*Client).groupAtMessageHandler,
	},
}

func (c *Client) parseAndHandle(payload *WSPayload) error {
	handler, has := eventParseFuncMap[payload.OPCode][payload.Type]
	if has {
		return handler(c, payload)
	}
	if c.eventHandlers.DefaultHandler != nil {
		return c.eventHandlers.DefaultHandler(c.api, payload)
	}
	return nil
}

func (c *Client) helloHandler(payload *WSPayload) error {
	helloData := &WSHelloData{}
	err := ParseData(payload.RawMessage, helloData)
	if err != nil {
		return errors.Wrap(err, "parse hello data fail")
	}
	c.heartBeatTicker.Reset(time.Duration(helloData.HeartbeatInterval) * time.Millisecond)
	return nil
}

func (c *Client) reconnectHandler(payload *WSPayload) error {
	// 连接时间太长会收到这个，需要重新连接
	// 直接 close，交给上层执行重连
	c.closeErrChan <- errors.New("need reconnect")
	return nil
}

func (c *Client) invalidSessionHandler(payload *WSPayload) error {
	c.closeErrChan <- errors.New("invalid session")
	return nil
}

func (c *Client) heartbeatAckHandler(payload *WSPayload) error {
	logrus.Infof("heartbeat ack. message: %s", payload.RawMessage)
	return nil
}

func (c *Client) readyHandler(payload *WSPayload) error {
	readyData := &WSReadyData{}
	err := ParseData(payload.RawMessage, readyData)
	if err != nil {
		return errors.Wrap(err, "parse ready data fail")
	}
	c.session.ID = readyData.SessionID
	c.session.User.ID = readyData.User.ID
	c.session.User.Username = readyData.User.Username

	if c.eventHandlers.ReadyHandler != nil {
		return c.eventHandlers.ReadyHandler(c.api, payload, readyData)
	}
	return nil
}

func (c *Client) groupAtMessageHandler(payload *WSPayload) error {
	groupAtMessageData := &WSGroupAtMessageData{}
	err := ParseData(payload.RawMessage, groupAtMessageData)
	if err != nil {
		return errors.Wrap(err, "parse group at message data fail")
	}
	groupAtMessageData.WSPayloadBase = payload.WSPayloadBase
	if c.eventHandlers.GroupAtMessageHandler != nil {
		return c.eventHandlers.GroupAtMessageHandler(c.api, payload, groupAtMessageData)
	}
	return nil
}
