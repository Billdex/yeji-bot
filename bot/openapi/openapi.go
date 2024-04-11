package openapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type Openapi struct {
	appID        uint64
	clientSecret string
	sandbox      bool

	httpClient *http.Client
	timeout    time.Duration

	accessToken   string
	tokenExpireAt time.Time
	t             *time.Ticker
}

func New(appID uint64, appSecret string, sandbox bool) (*Openapi, error) {
	openapi := &Openapi{
		appID:        appID,
		clientSecret: appSecret,
		sandbox:      sandbox,
		httpClient:   &http.Client{},
		timeout:      3 * time.Second,
		t:            time.NewTicker(30 * time.Second),
	}

	err := openapi.refreshAccessToken()
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case <-openapi.t.C:
				err := openapi.refreshAccessToken()
				if err != nil {
					logrus.Errorf("refresh access token fail. %v", err)
				}
			}
		}
	}()

	return openapi, nil
}

func (o *Openapi) Close() {
	o.t.Stop()
}

func (o *Openapi) request(ctx context.Context, request Request, response Response, pathParams map[string]string) (err error) {
	if time.Now().Sub(o.tokenExpireAt).Seconds() < 60 {
		err = o.refreshAccessToken()
		if err != nil {
			return err
		}
	}
	byteBody, err := request.Marshal()
	if err != nil {
		return err
	}
	req, err := http.NewRequest(request.Method(), o.getURL(request.URI(), pathParams), bytes.NewReader(byteBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("QQBot %s", o.accessToken))
	req.Header.Set("X-Union-Appid", fmt.Sprintf("%d", o.appID))
	startAt := time.Now()
	resp, err := o.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	if resp.StatusCode > 204 {
		var errResponse ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResponse)
		if err != nil {
			return err
		}
		return fmt.Errorf("%d %s", errResponse.Code, errResponse.Message)
	}
	byteRespBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logrus.Infof("openapi request. url: %s, t: %s, req_body: %s, resp_body: %s", req.URL, time.Now().Sub(startAt).String(), string(byteBody), string(byteRespBody))
	err = response.Unmarshal(byteRespBody)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return nil
}

// refreshAccessToken 刷新 access token
func (o *Openapi) refreshAccessToken() (err error) {
	// 所用的 url/header 格式与其他请求不同，不调用通用 request
	var reqBody = struct {
		AppId        string `json:"appId"`
		ClientSecret string `json:"clientSecret"`
	}{
		AppId:        fmt.Sprintf("%d", o.appID),
		ClientSecret: o.clientSecret,
	}
	byteReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s://%s", scheme, accessTokenURL), bytes.NewReader(byteReqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := o.httpClient.Do(req)
	if err != nil {
		return err
	}
	var respBody struct {
		Code        int64       `json:"code"`
		Message     string      `json:"message"`
		AccessToken string      `json:"access_token"`
		ExpiresIn   json.Number `json:"expires_in"`
	}
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&respBody)
	if err != nil {
		return err
	}
	if respBody.Code != 0 {
		return fmt.Errorf("%d %s", respBody.Code, respBody.Message)
	}
	o.accessToken = respBody.AccessToken
	expiresIn, _ := respBody.ExpiresIn.Int64()
	o.tokenExpireAt = time.Now().Add(time.Duration(expiresIn) * time.Second)

	return nil
}
