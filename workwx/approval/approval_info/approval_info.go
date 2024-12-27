package approval_info

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goDeploy/config"
	p "goDeploy/public"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	ApprovalInfoApi = "https://qyapi.weixin.qq.com/cgi-bin/oa/getapprovalinfo"
	maxRetries      = 3
	timeout         = 5 * time.Second
)

// ApprovalInfoRequest 定义请求结构体
type ApprovalInfoRequest struct {
	Starttime string `json:"starttime"`
	Endtime   string `json:"endtime"`
	NewCursor string `json:"new_cursor"`
	Size      int    `json:"size"`
	Filters   []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"filters"`
}

// ApprovalInfoResponse 定义响应结构体
type ApprovalInfoResponse struct {
	Errcode  int      `json:"errcode"`
	Errmsg   string   `json:"errmsg"`
	SpNoList []string `json:"sp_no_list"`
}

// getApprovalInfo 发送请求并获取审批信息
func approvalInfo(url string, request *ApprovalInfoRequest) ([]string, error) {
	var resp ApprovalInfoResponse

	// JSON编码请求体
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	// 创建HTTP客户端，并设置超时
	client := &http.Client{
		Timeout: timeout,
	}

	// 定义重试的请求函数
	retryFunc := func() (*http.Response, error) {
		return client.Post(url, config.ContentType, bytes.NewBuffer(requestBody))
	}

	// 执行重试逻辑
	respBody, err := p.WithRetry(retryFunc, maxRetries, 2*time.Second)
	if err != nil {
		return nil, err
	}
	defer respBody.Body.Close() // 确保关闭响应体

	// 读取响应体
	body, err := ioutil.ReadAll(respBody.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应体
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	// 校验状态码
	if resp.Errcode != 0 {
		fmt.Printf("Error: %s\n", resp.Errmsg)
		return nil, fmt.Errorf("API error: %s", resp.Errmsg)
	}

	// 成功获取sp_no_list
	return resp.SpNoList, nil
}

func Get(startTimestamp, endTimestamp, templateId, accessToken string) ([]string, error) {
	// 构建请求URL和请求体
	url := ApprovalInfoApi + "?access_token=" + accessToken
	request := &ApprovalInfoRequest{
		Starttime: startTimestamp,
		Endtime:   endTimestamp,
		NewCursor: "",
		Size:      10,
		Filters: []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{
			{Key: "template_id", Value: templateId},
			{Key: "sp_status", Value: "1"},
		},
	}

	// 获取审批信息
	spNoList, err := approvalInfo(url, request)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return spNoList, err
}
