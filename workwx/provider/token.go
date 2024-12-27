package provider

import (
	"encoding/json"
	"fmt"
	"goDeploy/config"
	p "goDeploy/public"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	maxRetries = 3
	timeout    = 5 * time.Second
	retryDelay = 2 * time.Second
	baseURL    = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
)

// TokenResponse 定义获取access_token的响应结构体
type TokenResponse struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// getAccessToken 发送请求并获取access_token
func getAccessToken(corpid, corpsecret string) (*TokenResponse, error) {
	// 创建HTTP客户端，并设置超时
	client := &http.Client{
		Timeout: timeout,
	}

	// 定义重试的请求函数
	retryFunc := func() (*http.Response, error) {
		url := fmt.Sprintf("%s?corpid=%s&corpsecret=%s", baseURL, corpid, corpsecret)
		return client.Get(url)
	}

	// 执行重试逻辑
	resp, err := p.WithRetry(retryFunc, maxRetries, retryDelay)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // 确保关闭响应体

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应体
	var tokenResp TokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return nil, err
	}

	// 校验状态码
	if tokenResp.Errcode != 0 {
		return nil, fmt.Errorf("API error: %s", tokenResp.Errmsg)
	}

	// 成功获取access_token
	return &tokenResp, nil
}

func Token() (*string, error) {
	// 获取access_token
	tokenResp, err := getAccessToken(config.Corpid, config.ProviderSecret)
	if err != nil {
		return nil, err
	}

	return &tokenResp.AccessToken, nil
}
