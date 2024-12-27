package run

import (
	"goDeploy/mongodb"
	"goDeploy/workwx/user"
	"log"
	"strings"
	"time"
)

func User() {
	for {
		ak, err := mongodb.MyTokenManager.GetToken()
		if err != nil {
			log.Println("get token err:", err)
			continue
		}

		departmentIds := user.Ids(ak.AK)
		for _, id := range departmentIds {
			userList := user.List(id, ak.AK)
			for _, details := range userList {
				userDetails := mongodb.UserDetails{
					UserID:       details.UserID,
					Name:         details.Name,
					Mobile:       details.Mobile,
					Email:        details.Email,
					BizMail:      details.BizMail,
					DirectLeader: strings.Join(details.DirectLeader, ","),
				}
				err := mongodb.MyUserManager.AddOrUpdateUser(userDetails)
				if err != nil {
					log.Println("update userDetails err:", err)
					break
				}
			}
		}
		time.Sleep(12 * time.Hour)
	}
}
