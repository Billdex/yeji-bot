package openapi

import (
	"context"
	"encoding/json"
	"net/http"
)

type FileType int

const (
	FileTypeImage FileType = 1
	FileTypeVideo FileType = 2
	FileTypeVoice FileType = 3
	FileTypeFile  FileType = 4
)

// UploadGroupFileReq 上传群文件
type UploadGroupFileReq struct {
	FileType   FileType `json:"file_type"`    // 文件类型
	URL        string   `json:"url"`          // 文件 url
	SrvSendMsg bool     `json:"srv_send_msg"` // 是否直接发送到目标端，为 true 则占用主动发送频次
}

func (u *UploadGroupFileReq) Method() string {
	return http.MethodPost
}

func (u *UploadGroupFileReq) URI() string {
	return uploadGroupFileURI
}

func (u *UploadGroupFileReq) Marshal() ([]byte, error) {
	return json.Marshal(*u)
}

type UploadGroupFileResp struct {
	FileUUID string `json:"file_uuid"`
	FileInfo string `json:"file_info"`
	TTL      int    `json:"ttl"`
	Id       string `json:"id"`
}

func (u *UploadGroupFileResp) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, u)
}

func (o *Openapi) UploadGroupFile(ctx context.Context, groupOpenid string, req *UploadGroupFileReq) (*UploadGroupFileResp, error) {
	resp := new(UploadGroupFileResp)
	err := o.request(ctx, req, resp, map[string]string{
		"group_openid": groupOpenid,
	})
	return resp, err
}
