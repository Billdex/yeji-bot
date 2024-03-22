package qbot

// OPCode websocket op 码
type OPCode int

// WS OPCode
// https://bot.q.qq.com/wiki/develop/api/gateway/opcode.html
const (
	OPCodeDispatchEvent OPCode = iota
	OPCodeHeartbeat            // 客户端心跳
	OPCodeIdentity             // 客户端发起验证
	_                          // Presence Update
	_                          // Voice State Update
	_
	OPCodeResume // 重连
	OPCodeReconnect
	_                    // Request Guild Members
	OPCodeInvalidSession // identity 失败下发的错误信息
	OPCodeHello          // 连接成功后下发的第一条消息
	OPCodeHeartbeatAck   // 服务端确认客户端上报的心跳
	HTTPCallbackAck
)

// opMeans op 对应的含义字符串标识
var opMeans = map[OPCode]string{
	OPCodeDispatchEvent:  "Event",
	OPCodeHeartbeat:      "Heartbeat",
	OPCodeIdentity:       "Identity",
	OPCodeResume:         "Resume",
	OPCodeReconnect:      "Reconnect",
	OPCodeInvalidSession: "InvalidSession",
	OPCodeHello:          "Hello",
	OPCodeHeartbeatAck:   "HeartbeatAck",
}

// OPMeans 返回 op 含义
func OPMeans(op OPCode) string {
	means, ok := opMeans[op]
	if !ok {
		means = "unknown"
	}
	return means
}
