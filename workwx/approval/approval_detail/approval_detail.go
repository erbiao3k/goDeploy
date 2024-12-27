package approval_detail

import (
	"bytes"
	"encoding/json"
	"fmt"
	p "goDeploy/public"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	maxRetries = 3
	timeout    = 5 * time.Second
	retryDelay = 2 * time.Second
	apiURL     = "https://qyapi.weixin.qq.com/cgi-bin/oa/getapprovaldetail"
)

// RequestBody 定义请求体结构
type RequestBody struct {
	SpNo string `json:"sp_no"`
}

// ResponseInfo 定义响应体中的info结构
type ResponseInfo struct {
	SpNo        string      `json:"sp_no"`
	SpName      string      `json:"sp_name"`
	SpStatus    int         `json:"sp_status"`
	TemplateID  string      `json:"template_id"`
	ApplyTime   int64       `json:"apply_time"`
	Applyer     ApplyerInfo `json:"applyer"`
	SpRecord    []SpRecord  `json:"sp_record"`
	Notifyer    []Notifier  `json:"notifyer"`
	ApplyData   *ApplyData  `json:"apply_data"` // 使用指针类型以处理空值
	Comments    []Comment   `json:"comments"`
	ProcessList ProcessList `json:"process_list"`
}

// Response 定义响应体结构
type Response struct {
	Errcode int          `json:"errcode"`
	Errmsg  string       `json:"errmsg"`
	Info    ResponseInfo `json:"info"`
}

// ApplyerInfo 申请人信息
type ApplyerInfo struct {
	UserID  string `json:"userid"`
	PartyID string `json:"partyid"`
}

// SpRecord 审批流程信息
type SpRecord struct {
	SpStatus     int      `json:"sp_status"`
	ApproverAttr int      `json:"approverattr"`
	Details      []Detail `json:"details"`
}

// Detail 审批节点详情
type Detail struct {
	Approver Approver `json:"approver"`
	Speech   string   `json:"speech"`
	SpStatus int      `json:"sp_status"`
	Sptime   int      `json:"sptime"`
	MediaID  []string `json:"media_id"`
}

// Approver 分支审批人
type Approver struct {
	UserID string `json:"userid"`
}

// Notifier 抄送信息
type Notifier struct {
	UserID string `json:"userid"`
}

// ApplyData 审批申请数据
type ApplyData struct {
	Contents []Content `json:"contents"`
}

// Content 审批申请详情
type Content struct {
	Control string  `json:"control"`
	ID      string  `json:"id"`
	Title   []Title `json:"title"`
	Value   *Value  `json:"value"` // 使用指针类型以处理空值
}

// Title 控件名称
type Title struct {
	Text string `json:"text"`
	Lang string `json:"lang"`
}

// Value 控件值
type Value struct {
	Text            string    `json:"text"`
	Tips            []string  `json:"tips"`
	Date            *Date     `json:"date"`  // 使用指针类型以处理空值
	Files           []File    `json:"files"` // 修改为 File 类型的切片
	Children        []string  `json:"children"`
	StatField       []string  `json:"stat_field"`
	SumField        []string  `json:"sum_field"`
	RelatedApproval []string  `json:"related_approval"`
	Students        []string  `json:"students"`
	Classes         []string  `json:"classes"`
	Docs            []string  `json:"docs"`
	WedriveFiles    []string  `json:"wedrive_files"`
	Selector        *Selector `json:"selector"` // 新增Selector字段
}

// File 文件结构
type File struct {
	FileID string `json:"file_id"`
}

// Date 包含日期时间的字段
type Date struct {
	Type       string `json:"type"`
	STimestamp string `json:"s_timestamp"`
}

// Comment 审批申请备注信息
type Comment struct {
	CommentUserInfo CommentUserInfo `json:"commentUserInfo"`
	CommentTime     int64           `json:"commenttime"`
	CommentContent  string          `json:"commentcontent"`
	CommentID       string          `json:"commentid"`
	MediaID         []string        `json:"media_id"`
}

// CommentUserInfo 备注人信息
type CommentUserInfo struct {
	UserID string `json:"userid"`
}

// ProcessList 审批流程列表
type ProcessList struct {
	NodeList []Node `json:"node_list"`
}

// Node 流程节点
type Node struct {
	NodeType    int       `json:"node_type"`
	SpStatus    int       `json:"sp_status"`
	ApvRel      int       `json:"apv_rel"`
	SubNodeList []SubNode `json:"sub_node_list"`
}

// SubNode 子节点列表
type SubNode struct {
	UserID   string   `json:"userid"`
	Speech   string   `json:"speech"`
	SpYj     int      `json:"sp_yj"`
	Sptime   int      `json:"sptime"`
	MediaIDs []string `json:"media_ids"`
}

// Selector 选择器结构
type Selector struct {
	Type    string   `json:"type"`
	Options []Option `json:"options"`
}

// Option 选项结构
type Option struct {
	Key   string  `json:"key"`
	Value []Title `json:"value"`
}

// approvalDetail 发送请求并获取审批申请详情
func approvalDetail(token string, reqBody *RequestBody) (*Response, error) {
	var resp Response

	// JSON编码请求体
	requestBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	// 创建HTTP客户端，并设置超时
	client := &http.Client{
		Timeout: timeout,
	}

	// 定义重试的请求函数
	retryFunc := func() (*http.Response, error) {
		url := fmt.Sprintf("%s?access_token=%s", apiURL, token)
		return client.Post(url, "application/json", bytes.NewBuffer(requestBody))
	}

	// 执行重试逻辑
	respBody, err := p.WithRetry(retryFunc, maxRetries, retryDelay)
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
		return nil, fmt.Errorf("API error: %s", resp.Errmsg)
	}

	// 成功获取审批申请详情
	return &resp, nil
}

func Get(spNo, accessToken string) (*ResponseInfo, error) {
	reqBody := &RequestBody{
		SpNo: spNo,
	}

	// 获取审批申请详情
	resp, err := approvalDetail(accessToken, reqBody)
	if err != nil {
		return nil, err
	}

	// 打印审批申请详情
	return &resp.Info, nil
}
