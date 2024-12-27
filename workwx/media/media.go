package media

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	maxRetries = 3
	retryDelay = 2 * time.Second
	apiUrl     = "https://qyapi.weixin.qq.com/cgi-bin/media/get"
)

// downloadFile 支持断点下载的函数
func downloadFile(mediaID string, accessToken string, destFile string) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?access_token=%s&media_id=%s", apiUrl, accessToken, mediaID), nil)
	if err != nil {
		return err
	}

	// 尝试打开或创建文件，如果文件存在则截断文件内容
	file, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 发起请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected run code: %d", resp.StatusCode)
	}

	// 将响应体写入文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func Get(mediaId, destFile, accessToken string) error {
	err := downloadFile(mediaId, accessToken, destFile)
	if err != nil {
		return errors.New("Download error：" + err.Error())
	}
	return nil
}
