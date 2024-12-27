package run

import (
	"bytes"
	"errors"
	"goDeploy/config"
	"net/http"
)

const (
	// 定义发布状态的常量
	DeployStatusCommit     = "commit"     // 审批阶段，已在企业微信提交发布申请
	DeployStatusWaitDeploy = "waitDeploy" // 审批阶段，已通过申请待触发发布地址

	DeployStatusDeploying      = "正在部署中" // 部署运行阶段，已触发发布地址，部署中
	DeployStatusDeployCanceled = "部署取消"  // 部署运行阶段，已触发发布地址，已取消部署
	DeployStatusDeploySuccess  = "部署成功"  // 部署运行阶段，部署完成成功待自动化测试
	DeployStatusDeployFailed   = "部署失败"  // 部署运行阶段，部署完成失败

	DeployStatusDeploySuccessFinish  = "部署成功且归档" // 部署归档阶段，成功完成时归档
	DeployStatusDeployFailedFinish   = "部署失败且归档" // 部署归档阶段，失败完成时归档
	DeployStatusDeployCanceledFinish = "部署取消且归档" // 部署归档阶段，取消完成时归档

	DeployStatusTesting         = "正在测试中" // 测试运行阶段，已触发自动化测试，测试中
	DeployStatusTestingSuccess  = "测试通过"  // 测试运行阶段，已触发自动化测试待归档
	DeployStatusTestingCanceled = "测试取消"  // 测试运行阶段，已触发自动化测试，取消
	DeployStatusTestingFailed   = "测试不通过" // 测试运行阶段，测试失败

	DeployStatusFinalSuccess  = "发布成功" // 测试归档阶段，成功完成时归档
	DeployStatusFinalFailed   = "发布失败" // 测试归档阶段，失败完成时归档
	DeployStatusFinalCanceled = "发布取消" // 测试归档阶段，取消完成时归档

	// 申请单状态常量
	SPStatusPending   = 1  // 审批中
	SPStatusApproved  = 2  // 已通过
	SPStatusRejected  = 3  // 已驳回
	SPStatusWithdrawn = 4  // 已撤销
	SPStatusRevoked   = 6  // 通过后撤销
	SPStatusDeleted   = 7  // 已删除
	SPStatusPaid      = 10 // 已支付

	// 节点类型常量
	NodeTypeApprover = 1 // 审批人
	NodeTypeCc       = 2 // 抄送人
	NodeTypeHandler  = 3 // 办理人

)

// postToWebhook 发送POST请求到指定的webhook URL并检查响应是否成功
func postToWebhook(url, data string) error {
	if len(data) == 0 {
		data = "{}"
	}
	// 设置请求的正文
	requestBody := bytes.NewBufferString(data)

	// 创建HTTP客户端
	client := &http.Client{}

	// 构建请求
	req, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return err
	}

	// 设置请求头部
	req.Header.Set("Content-Type", config.ContentType)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查HTTP响应状态码
	if resp.StatusCode < 400 {
		// 响应成功
		return nil
	} else {
		return errors.New(resp.Status)
	}
}
