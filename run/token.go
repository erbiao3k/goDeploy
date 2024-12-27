package run

import (
	"fmt"
	"goDeploy/mongodb"
	"goDeploy/workwx/provider"
	"log"
	"time"
)

func Token() {
	for {
		tk, err := provider.Token()
		if err != nil {
			log.Fatal("get provider token err:", err)
		}
		err = mongodb.MyTokenManager.UpsertToken(&mongodb.Token{AK: *tk})
		if err != nil {
			log.Fatal(fmt.Errorf("upsert token err:%v", err))
		}
		time.Sleep(30 * time.Minute)
	}
}
