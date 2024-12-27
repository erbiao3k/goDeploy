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

var endingMap = map[string]string{
	DeployStatusTestingCanceled: DeployStatusFinalCanceled,
	DeployStatusTestingFailed:   DeployStatusFinalFailed,
	DeployStatusTestingSuccess:  DeployStatusFinalSuccess,
}

func Tested() {
	for {
		for startStatus, endStatus := range endingMap {
			err := StatusNotice(startStatus, endStatus)
			if err != nil {
				log.Println(err)
				continue
			}
		}
		time.Sleep(3 * time.Second)
	}
}

func StatusNotice(startDeployStatus, endDeployStatus string) error {
	filter := bson.M{"deploy_status": startDeployStatus}
	approvalList, err := mongodb.MyApprovalManager.GetApprovals(filter)
	if err != nil {
		return errors.New(fmt.Sprintf("get %v approval err: %v", startDeployStatus, err))
	}
	if len(approvalList) == 0 {
		return nil
	}
	spNoApprovalsMap := make(map[string][]mongodb.Approval)
	for _, approval := range approvalList {
		spNoApprovalsMap[approval.SpNo] = append(spNoApprovalsMap[approval.SpNo], approval)
	}

	for spNo, approvals := range spNoApprovalsMap {
		var applist []string
		for _, approval := range approvals {
			applist = append(applist, approval.AppName)
		}
		msg := "【发布结束自动化测试】" +
			"\n审批单号：" + spNo +
			"\n发布时间：" + p.TimestampInt64ToTime(approvals[0].DeployTime) +
			"\n应用列表：" + strings.Join(applist, ",") +
			"\n发布内容：" + approvals[0].DeployContent +
			"\n测试结果：" + startDeployStatus +
			"\n自动化测试：https://workspace.apipost.net/30f939f79464000/testing" +
			"\n" + startDeployStatus + "，发布结束！！！" +
			"\n申请人：" + approvals[0].Applyer
		data := fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s","mentioned_list":["%s"]}}`, msg, "@all")
		err = p.Send(data)
		if err != nil {
			log.Println("post msg err: " + err.Error())
		}
	}
	_, err = mongodb.MyApprovalManager.UpdateApprovalMany(bson.M{"deploy_status": startDeployStatus}, bson.M{"deploy_status": endDeployStatus})
	if err != nil {
		return errors.New(fmt.Sprintf("update %v approval err: %v", startDeployStatus, err))
	}
	return nil
}
