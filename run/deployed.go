package run

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"goDeploy/mongodb"
	p "goDeploy/public"
	"log"
	"strings"
	"time"
)

func Deployed() {
	for {
		canceledApprovalList, err := getUpdate(DeployStatusDeployCanceled, DeployStatusFinalCanceled)
		if err != nil {
			log.Println(err)
			continue
		}

		failedApprovalList, err := getUpdate(DeployStatusDeployFailed, DeployStatusFinalFailed)
		if err != nil {
			log.Println(err)
			continue
		}

		successApprovalList, err := getUpdate(DeployStatusDeploySuccess, DeployStatusDeploySuccessFinish)
		if err != nil {
			log.Println(err)
			continue
		}

		spNoApprovalMap := make(map[string][]mongodb.Approval)
		for _, approval := range append(canceledApprovalList, append(failedApprovalList, successApprovalList...)...) {
			spNoApprovalMap[approval.SpNo] = append(spNoApprovalMap[approval.SpNo], approval)
		}

		for spNo, approvals := range spNoApprovalMap {
			spNoFilter := bson.M{"sp_no": spNo}
			spNoApprovalList, err := mongodb.MyApprovalManager.GetApprovals(spNoFilter)
			if err != nil {
				log.Printf("get spNo %v approval err: %v", spNo, err)
				break
			}

			var appStatuslist []string
			for _, approval := range approvals {
				appStatuslist = append(appStatuslist, approval.AppName+fmt.Sprintf("(%v)", approval.DeployStatus))
			}
			if len(approvals) == len(spNoApprovalList) {
				msg := "【发布已完成应用部署】" +
					"\n审批单号：" + spNo +
					"\n申请时间：" + p.TimestampInt64ToTime(approvals[0].ApplyTime) +
					"\n发布时间：" + p.TimestampInt64ToTime(approvals[0].DeployTime) +
					"\n发布内容：" + approvals[0].DeployContent +
					"\n发布应用：" + strings.Join(appStatuslist, ",") +
					"\n应用地址：https://devops.aliyun.com/appstack/apps" +
					"\n所有应用部署成功时，才进行自动化测试!!!" +
					"\n所有应用部署成功时，才进行自动化测试!!!" +
					"\n所有应用部署成功时，才进行自动化测试!!!" +
					"\n申请人：" + approvals[0].Applyer
				data := fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s","mentioned_list":["%s"]}}`, msg, "@all")
				err = p.Send(data)
				if err != nil {
					log.Println("post msg err: " + err.Error())
				}
			}
		}
		time.Sleep(3 * time.Second)
	}
}

func getUpdate(initStatus, endStatus string) ([]mongodb.Approval, error) {
	filter := bson.M{"deploy_status": initStatus}
	canceledApprovalList, err := mongodb.MyApprovalManager.GetApprovals(filter)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("get %v approval err: %v", initStatus, err))
	}
	updateStatus := bson.M{"deploy_status": endStatus}
	_, err = mongodb.MyApprovalManager.UpdateApprovalMany(filter, updateStatus)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("update approval status %v to %v err: %v", initStatus, endStatus, err))
	}
	return canceledApprovalList, nil
}
