package openapi

import (
	"context"
	"encoding/json"
	"net/http"
)

type WSReq struct {
}

func (W *WSReq) Method() string {
	return http.MethodGet
}

func (W *WSReq) URI() string {
	return gatewayURI
}

func (W *WSReq) PathParams() map[string]string {
	return map[string]string{}
}

func (W *WSReq) Marshal() ([]byte, error) {
	return json.Marshal(*W)
}

type WSResp struct {
	URL string `json:"url"`
}

func (W *WSResp) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, W)
}

// WS 获取 websockets 接入点信息
func (o *Openapi) WS(ctx context.Context) (*WSResp, error) {
	ws := new(WSResp)
	err := o.request(ctx, &WSReq{}, ws, nil)
	return ws, err
}
