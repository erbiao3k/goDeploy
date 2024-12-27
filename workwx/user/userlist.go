package user

import (
	"encoding/json"
	"fmt"
	p "goDeploy/public"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	userListUrl = "https://qyapi.weixin.qq.com/cgi-bin/user/list"
)

// UserListResponse 定义用户列表响应结构体
type UserListResponse struct {
	Errcode  int       `json:"errcode"`
	Errmsg   string    `json:"errmsg"`
	UserList []Details `json:"userlist"`
}

// Details 定义用户详细信息结构体
type Details struct {
	UserID           string          `json:"userid"`
	Name             string          `json:"name"`
	Department       []int           `json:"user"`
	Order            []int           `json:"order"`
	Position         string          `json:"position"`
	Mobile           string          `json:"mobile"`
	Gender           int             `json:"gender"`
	Email            string          `json:"email"`
	BizMail          string          `json:"biz_mail"`
	IsLeaderInDept   []int           `json:"is_leader_in_dept"`
	DirectLeader     []string        `json:"direct_leader"`
	Avatar           string          `json:"avatar"`
	ThumbAvatar      string          `json:"thumb_avatar"`
	Telephone        string          `json:"telephone"`
	Alias            string          `json:"alias"`
	Status           int             `json:"run"`
	Address          string          `json:"address"`
	EnglishName      string          `json:"english_name"`
	OpenUserid       string          `json:"open_userid"`
	MainDepartment   int             `json:"main_department"`
	Extattr          Extattr         `json:"extattr"`
	QRCode           string          `json:"qr_code"`
	ExternalPosition string          `json:"external_position"`
	ExternalProfile  ExternalProfile `json:"external_profile"`
}

// Extattr 定义扩展属性结构体
type Extattr struct {
	Attrs []struct {
		Type int    `json:"type"`
		Name string `json:"name"`
		Text struct {
			Value string `json:"value"`
		} `json:"text"`
		Web struct {
			URL   string `json:"url"`
			Title string `json:"title"`
		} `json:"web"`
	} `json:"attrs"`
}

// ExternalProfile 定义对外属性结构体
type ExternalProfile struct {
	ExternalCorpName string `json:"external_corp_name"`
	WeChatChannels   struct {
		Nickname string `json:"nickname"`
		Status   int    `json:"run"`
	} `json:"wechat_channels"`
	ExternalAttr []struct {
		Type int    `json:"type"`
		Name string `json:"name"`
		Text struct {
			Value string `json:"value"`
		} `json:"text"`
		Web struct {
			URL   string `json:"url"`
			Title string `json:"title"`
		} `json:"web"`
		MiniProgram struct {
			AppID    string `json:"appid"`
			PagePath string `json:"pagepath"`
			Title    string `json:"title"`
		} `json:"miniprogram"`
	} `json:"external_attr"`
}

// getUserList 发送请求并获取部门成员详情
func getUserList(accessToken string, departmentID int) (*UserListResponse, error) {
	// 创建HTTP客户端，并设置超时
	client := &http.Client{
		Timeout: timeout,
	}

	// 定义重试的请求函数
	retryFunc := func() (*http.Response, error) {
		url := fmt.Sprintf("%s?access_token=%s&department_id=%d", userListUrl, accessToken, departmentID)
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
	var userListResp UserListResponse
	err = json.Unmarshal(body, &userListResp)
	if err != nil {
		return nil, err
	}

	// 校验状态码
	if userListResp.Errcode != 0 {
		return nil, fmt.Errorf("API error: %s", userListResp.Errmsg)
	}

	// 成功获取部门成员详情
	return &userListResp, nil
}

func List(departmentID int, accessToken string) []Details {
	// 获取部门成员详情
	userListResp, err := getUserList(accessToken, departmentID)
	if err != nil {
		log.Fatalf("Get user list failed: %v\n", err)
	}

	return userListResp.UserList
}
