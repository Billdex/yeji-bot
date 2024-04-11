package openapi

import (
	"fmt"
	"net/url"
	"strings"
)

const accessTokenURL = "bots.qq.com/app/getAppAccessToken"

const (
	scheme = "https"

	domain        = "api.sgroup.qq.com"
	sandboxDomain = "sandbox.api.sgroup.qq.com"
)

const (
	gatewayURI    = "/gateway"     // 获取 ws 接入地址
	gatewayBotURI = "/gateway/bot" // 获取 ws 接入地址与分片信息

	uploadGroupFileURI = "/v2/groups/{group_openid}/files" // 上传群聊用的富媒体文件

	postGroupMessageURI = "/v2/groups/{group_openid}/messages" // 发送群聊消息
)

func (o *Openapi) getURL(endpoint string, pathParams map[string]string) string {
	d := domain
	if o.sandbox {
		d = sandboxDomain
	}
	if len(pathParams) > 0 {
		for p, v := range pathParams {
			endpoint = strings.Replace(endpoint, "{"+p+"}", url.PathEscape(v), -1)
		}
	}
	return fmt.Sprintf("%s://%s%s", scheme, d, endpoint)
}
