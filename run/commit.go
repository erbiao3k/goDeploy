package run

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"goDeploy/config"
	"goDeploy/mongodb"
	p "goDeploy/public"
	"goDeploy/workwx/approval/approval_detail"
	"goDeploy/workwx/approval/approval_info"
	"goDeploy/workwx/media"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

func CommitApproval() {
	for {
		now := time.Now()

		ak, err := mongodb.MyTokenManager.GetToken()
		if err != nil {
			log.Println("get token err:", err)
			continue
		}
		startTimestamp := strconv.FormatInt(now.AddDate(0, 0, -7).Unix(), 10)

		endTimestamp := strconv.FormatInt(now.Unix(), 10)
		spNoList, err := approval_info.Get(startTimestamp, endTimestamp, config.DeployTemplateId, ak.AK)
		if err != nil {
			log.Println("get approval info err:", err)
			continue
		}

		for _, spNo := range spNoList {
			detail, err := approval_detail.Get(spNo, ak.AK)
			if err != nil {
				log.Println("get approval detail err: ", err)
				break
			}

			deployTimeSet := detail.ApplyData.Contents[0]
			deployAppSet := detail.ApplyData.Contents[2]
			deployContentSet := detail.ApplyData.Contents[3]
			deploySqlSet := detail.ApplyData.Contents[4]

			var deployAppList []string
			for _, option := range deployAppSet.Value.Selector.Options {
				deployAppList = append(deployAppList, option.Value[1].Text)
			}

			userInfo, err := mongodb.MyUserManager.GetUserByID(detail.Applyer.UserID)
			if err != nil {
				log.Println(detail)
				log.Println("get user info err:", err)
				break
			}

			var sqlsByte [][]byte
			for index, file := range deploySqlSet.Value.Files {
				fileName := detail.Applyer.UserID + "-" + strconv.FormatInt(detail.ApplyTime, 10) + "-" + strconv.FormatInt(int64(index), 10) + "-file.sql"
				err = media.Get(file.FileID, fileName, ak.AK)
				if err != nil {
					log.Println(detail)
					log.Println("get file err:", err)
					break
				}
				// 读取文件内容
				sqlContent, err := ioutil.ReadFile(fileName)
				if err != nil {
					log.Printf("Error reading file: %v\n", err)
					break
				}
				sqlsByte = append(sqlsByte, sqlContent)
			}

			var noticeBool bool
			for _, app := range deployAppList {
				noticeBool = true
				_, err = mongodb.MyApprovalManager.AddApproval(mongodb.Approval{
					SpName:        detail.SpName,
					SpNo:          detail.SpNo,
					ApplyTime:     detail.ApplyTime,
					Applyer:       userInfo.Name,
					UserId:        detail.Applyer.UserID,
					Tel:           userInfo.Mobile,
					DeployTime:    p.TimestampStringToInt64(deployTimeSet.Value.Date.STimestamp),
					AppName:       app,
					DeployContent: deployContentSet.Value.Text,
					SqlContent:    sqlsByte,
					DeployStatus:  DeployStatusCommit,
				})
				if err != nil {
					if mongo.IsDuplicateKeyError(err) {
						//log.Println("add AddApproval duplicate data：", err)
						noticeBool = false
					} else {
						noticeBool = false
						log.Printf("add approval err: %v\n", err)
						break
					}
				}
			}
			if noticeBool {
				msg := "【收到新发布申请】" +
					"\n审批名称：" + detail.SpName +
					"\n审批单号：" + detail.SpNo +
					"\n申请时间：" + p.TimestampInt64ToTime(detail.ApplyTime) +
					"\n发布时间：" + p.TimestampStringTOTime(deployTimeSet.Value.Date.STimestamp) +
					"\n发布应用：" + strings.Join(deployAppList, ",") +
					"\n发布内容：" + deployContentSet.Value.Text +
					"\nSQL文件：" + strconv.Itoa(len(deploySqlSet.Value.Files)) + "个" +
					"\n注意：【" + config.AutoNodeList + "】审批通过后将执行自动部署!!!" +
					"\n注意：【" + config.AutoNodeList + "】审批通过后将执行自动部署!!!" +
					"\n注意：【" + config.AutoNodeList + "】审批通过后将执行自动部署!!!"
				data := fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s","mentioned_list":["%s"]}}`, msg, detail.Applyer.UserID)
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
