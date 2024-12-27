package user

import (
	"encoding/json"
	"fmt"
	p "goDeploy/public"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	maxRetries = 3
	timeout    = 5 * time.Second
	retryDelay = 2 * time.Second
	baseURL    = "https://qyapi.weixin.qq.com/cgi-bin/department/simplelist"
)

// ListResponse  定义部门列表响应结构体
type ListResponse struct {
	Errcode  int    `json:"errcode"`
	Errmsg   string `json:"errmsg"`
	DeptList []Info `json:"department_id"`
}

// Info  定义部门信息结构体
type Info struct {
	ID       int `json:"id"`
	ParentID int `json:"parentid"`
	Order    int `json:"order"`
}

// getDepartmentList 发送请求并获取子部门ID列表
func getDepartmentList(accessToken string, departmentID int) (*ListResponse, error) {
	// 创建HTTP客户端，并设置超时
	client := &http.Client{
		Timeout: timeout,
	}

	// 定义重试的请求函数
	retryFunc := func() (*http.Response, error) {
		url := fmt.Sprintf("%s?access_token=%s&id=%d", baseURL, accessToken, departmentID)
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
	var deptListResp ListResponse
	err = json.Unmarshal(body, &deptListResp)
	if err != nil {
		return nil, err
	}

	// 校验状态码
	if deptListResp.Errcode != 0 {
		return nil, fmt.Errorf("API error: %s", deptListResp.Errmsg)
	}

	// 成功获取子部门ID列表
	return &deptListResp, nil
}

func Ids(accessToken string) []int {
	// 调用接口凭证和部门ID
	departmentID := 0 // 可以设置为0或不填，以获取全量组织架构

	// 获取子部门ID列表
	deptListResp, err := getDepartmentList(accessToken, departmentID)
	if err != nil {
		log.Fatalf("Get user list failed: %v\n", err)
	}

	var ids []int
	for _, info := range deptListResp.DeptList {
		ids = append(ids, info.ID)
	}
	return ids
}
