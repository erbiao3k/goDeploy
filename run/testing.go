package run

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"goDeploy/mongodb"
	p "goDeploy/public"
	"log"
	"strings"
	"time"
)

func Testing() {
	for {
		autotestAppinfo, err := mongodb.MyAppStackManager.GetAppStackByName("xod-autotest")
		if err != nil {
			log.Printf("get appstack xod-autotest failed: %v", err)
			continue
		}

		DeploySuccessFinishFilter := bson.M{"deploy_status": DeployStatusDeploySuccessFinish}
		DeploySuccessFinishApprovalList, err := mongodb.MyApprovalManager.GetApprovals(DeploySuccessFinishFilter)
		if err != nil {
			log.Printf("get %v approval err: %v", DeployStatusDeploySuccessFinish, err)
		}
		if len(DeploySuccessFinishApprovalList) == 0 {
			continue
		}

		spNoApprovalMap := make(map[string][]mongodb.Approval)
		for _, approval := range DeploySuccessFinishApprovalList {
			spNoApprovalMap[approval.SpNo] = append(spNoApprovalMap[approval.SpNo], approval)
		}

		// 判断每个spNo部署单下DeployStatusDeploySuccessFinish状态下部署应用数与提交的总数是否相等
		// 相等时则说明SpNo所有应用的部署都成功完成
		// 此时可以进行自动化测试
		for spNo, approvals := range spNoApprovalMap {
			spNoFilter := bson.M{"sp_no": spNo}
			spNoApprovalList, err := mongodb.MyApprovalManager.GetApprovals(spNoFilter)
			if err != nil {
				log.Printf("get spNo %v approval err: %v", spNo, err)
				break
			}
			if len(spNoApprovalList) == len(approvals) {
				err := postToWebhook(autotestAppinfo.DeployWebhook, "")
				if err != nil {
					log.Printf("post to webhook %v err: %v", autotestAppinfo.DeployWebhook, err)
					break
				}
				_, err = mongodb.MyApprovalManager.UpdateApprovalMany(bson.M{"sp_no": spNo}, bson.M{"deploy_status": DeployStatusTesting})
				if err != nil {
					log.Printf("update spNo %v status %v err: %v", spNo, DeployStatusTesting, err)
					break
				}

				var applist []string
				for _, approval := range approvals {
					applist = append(applist, approval.AppName)
				}

				msg := "【发布开始自动化测试】" +
					"\n审批单号：" + spNo +
					"\n发布时间：" + p.TimestampInt64ToTime(approvals[0].DeployTime) +
					"\n应用列表：" + strings.Join(applist, ",") +
					"\n发布内容：" + approvals[0].DeployContent +
					"\n自动化测试：https://workspace.apipost.net/30f939f79464000/testing" +
					"\n已开始自动化测试!!!" +
					"\n已开始自动化测试!!!" +
					"\n已开始自动化测试!!!" +
					"\n申请人：" + approvals[0].Applyer

				data := fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s","mentioned_list":["%s"]}}`, msg, "@all")
				err = p.Send(data)
				if err != nil {
					log.Println("post msg err: " + err.Error())
				}
			}
		}
		time.Sleep(time.Second * 10)
	}
}
