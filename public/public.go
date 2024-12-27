package public

import (
	"errors"
	"fmt"
	"goDeploy/config"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// WithRetry 执行重试逻辑
func WithRetry(retryFunc func() (*http.Response, error), maxAttempts int, delay time.Duration) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i < maxAttempts; i++ {
		resp, err = retryFunc()
		if err != nil {
			fmt.Printf("Attempt %d failed: %v\n", i+1, err)
			if i < maxAttempts-1 {
				time.Sleep(delay) // 简单的指数退避策略
				continue
			}
			return nil, err
		}
		break
	}
	return resp, nil
}

func TimestampInt64ToTime(timestamp int64) string {
	// 将 int64 时间戳转换为 time.Time 类型
	t := time.Unix(timestamp, 0)
	// 格式化为可读的字符串
	// 这里使用 "2006-01-02 15:04:05" 作为格式化模板
	return t.Format("2006-01-02 15:04:05")
}

func TimestampStringTOTime(timestamp string) string {
	// 将字符串时间戳转换为 int64
	timestampInt64 := TimestampStringToInt64(timestamp)
	return TimestampInt64ToTime(timestampInt64)
}

func TimestampStringToInt64(timestamp string) int64 {
	timestampInt64, _ := strconv.ParseInt(timestamp, 10, 64)
	return timestampInt64
}

func Send(postData string) error {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())
	// 产生一个 [0, len(wecomRobotAddr)-1) 的随机整数
	randNum := rand.Intn(len(config.WecomRobotAddr) - 1)

	if _, err := http.Post(config.WecomRobotAddr[randNum], config.ContentType, strings.NewReader(postData)); err != nil {
		return errors.New("post wecom robot failed: " + err.Error())
	}
	return nil
}

// InStringSlice 检查字符串是否在字符串切片中
func InStringSlice(str string, slice []string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
