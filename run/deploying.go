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

func Deploying() {
	for {
		WaitDeployList, err := mongodb.MyApprovalManager.GetApprovals(bson.M{"deploy_status": DeployStatusWaitDeploy})
		if err != nil {
			log.Printf("get commit approval err: %v\n", err)
			continue
		}

		appList, err := mongodb.MyAppStackManager.ListAppStacks()
		if err != nil {
			log.Printf("get app list err: %v", err)
			continue
		}

		now := time.Now()
		for _, approval := range WaitDeployList {
			deployTime := time.Unix(approval.DeployTime, 0)

			// 检查当前时间是否大于等于部署时间
			if now.After(deployTime) || now.Equal(deployTime) {
				// 查找对应的应用堆栈信息
				for _, app := range appList {
					if approval.AppName == app.AppName {
						err := postToWebhook(app.DeployWebhook, "")
						if err != nil {
							log.Printf("post to webhook %v err: %v", app.DeployWebhook, err)
							break
						}
						_, err = mongodb.MyApprovalManager.UpdateApprovalMany(bson.M{"deploy_status": DeployStatusWaitDeploy}, bson.M{"deploy_status": DeployStatusDeploying})
						if err != nil {
							log.Printf("update status %v approval err: %v", DeployStatusWaitDeploy, err)
							break
						}
					}
				}
			}
		}

		deployingApplist, err := mongodb.MyApprovalManager.GetApprovals(bson.M{"deploy_status": DeployStatusDeploying})
		if err != nil {
			log.Printf("get status Deploying approval err: %v\n", err)
			continue
		}

		spNoApprovalMap := make(map[string][]mongodb.Approval)
		for _, approval := range deployingApplist {
			spNoApprovalMap[approval.SpNo] = append(spNoApprovalMap[approval.SpNo], approval)
		}

		for spNo, approvals := range spNoApprovalMap {
			spNoFilter := bson.M{"sp_no": spNo}
			spNoApprovalList, err := mongodb.MyApprovalManager.GetApprovals(spNoFilter)
			if err != nil {
				log.Printf("get spNo %v approval err: %v", spNo, err)
				break
			}
			if len(spNoApprovalList) == len(approvals) {
				var applist []string
				for _, approval := range approvals {
					applist = append(applist, approval.AppName)
				}

				msg := "【已到发布时间，开始部署应用】" +
					"\n审批单号：" + spNo +
					"\n申请时间：" + p.TimestampInt64ToTime(approvals[0].ApplyTime) +
					"\n发布时间：" + p.TimestampInt64ToTime(approvals[0].DeployTime) +
					"\n发布内容：" + approvals[0].DeployContent +
					"\n发布应用：" + strings.Join(applist, ",") +
					"\n应用地址：https://devops.aliyun.com/appstack/apps" +
					"\n应用正在部署中!!!" +
					"\n应用正在部署中!!!" +
					"\n应用正在部署中!!!" +
					"\n申请人：" + approvals[0].Applyer
				data := fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s","mentioned_list":["%s"]}}`, msg, "@all")
				err := p.Send(data)
				if err != nil {
					log.Println("post msg err: " + err.Error())
				}
			}
		}
		time.Sleep(time.Second * 30)
	}
}
