package run

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"goDeploy/config"
	"goDeploy/mongodb"
	p "goDeploy/public"
	"goDeploy/workwx/approval/approval_detail"
	"log"
	"time"
)

func WaitDeploy() {

	var userIds []string

	for _, name := range config.AutoNodePeople {
		userInfo, err := mongodb.MyUserManager.GetUserByName(name)
		if err != nil {
			log.Fatal("get auto node user info err: ", err)
		}
		userIds = append(userIds, userInfo.UserID)
	}

	for {
		ak, err := mongodb.MyTokenManager.GetToken()
		if err != nil {
			log.Panicln("get token err:", err)
		}

		commitedApproval, err := mongodb.MyApprovalManager.GetApprovals(bson.M{"deploy_status": DeployStatusCommit})
		if err != nil {
			log.Printf("get commit approval err: %v\n", err)
			continue
		}
		for index, approval := range commitedApproval {
			detail, err := approval_detail.Get(approval.SpNo, ak.AK)
			if err != nil {
				log.Println(detail)
				log.Printf("get approval detail err: %v\n", err)
				break
			}

			var changeStaus bool
			for _, node := range detail.ProcessList.NodeList {
				if node.NodeType == NodeTypeApprover && node.SpStatus == SPStatusApproved {
					for _, subNode := range node.SubNodeList {
						if p.InStringSlice(subNode.UserID, userIds) {
							changeStaus = true
						}
					}
				}
			}
			if changeStaus {
				_, err = mongodb.MyApprovalManager.UpdateApprovalMany(bson.M{"sp_no": approval.SpNo}, bson.M{"deploy_status": DeployStatusWaitDeploy})
				if err != nil {
					log.Printf("update many  status DeployStatusWaitDeploy approval err: %v\n", err)
					break
				}
			}

			if index+1 == len(commitedApproval) && changeStaus {
				msg := "【发布申请已通过】" +
					"\n审批名称：" + approval.SpName +
					"\n审批单号：" + approval.SpNo +
					"\n申请时间：" + p.TimestampInt64ToTime(approval.ApplyTime) +
					"\n发布时间：" + p.TimestampInt64ToTime(approval.DeployTime) +
					"\n发布申请已通过，部署将在指定时间自动执行!!!" +
					"\n发布申请已通过，部署将在指定时间自动执行!!!" +
					"\n发布申请已通过，部署将在指定时间自动执行!!!" +
					"\n申请人：" + approval.Applyer
				data := fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s","mentioned_list":["%s"]}}`, msg, approval.UserId)
				err := p.Send(data)
				if err != nil {
					log.Println(err)
					break
				}
			}
		}
		time.Sleep(time.Second * 10)
	}
}
