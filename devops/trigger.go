package devops

import (
	"bytes"
	"log"
	"net/http"
)

// Trigger 发送POST请求到指定的webhook URL并检查响应是否成功
func Trigger(url string) bool {
	// 设置请求的正文
	requestBody := bytes.NewBufferString("{}")

	// 创建HTTP客户端
	client := &http.Client{}

	// 构建请求
	req, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return false
	}

	// 设置请求头部
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return false
	}
	defer resp.Body.Close()

	// 检查HTTP响应状态码
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// 响应成功
		return true
	}

	// 响应失败
	return false
}
