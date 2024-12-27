package main

import (
	"goDeploy/mongodb"
	"goDeploy/router"
	"goDeploy/run"
	"log"
	"time"
)

var err error

func main() {
	defer mongodb.MyAnalysisManager.Close()
	defer mongodb.MyAppStackManager.Close()
	defer mongodb.MyApprovalManager.Close()
	defer mongodb.MyUserManager.Close()

	go run.Token()
	time.Sleep(5 * time.Second)
	go run.User()
	time.Sleep(5 * time.Second)

	err = mongodb.MyApprovalManager.EnsureIndexes()
	if err != nil {
		log.Fatalln("ensure approval index err: ", err)
	}

	err = mongodb.MyAppStackManager.EnsureIndexes()
	if err != nil {
		log.Fatalln("ensure appstack index err: ", err)
	}

	err = mongodb.MyUserManager.EnsureIndexes()
	if err != nil {
		log.Fatalln("ensure user index err: ", err)
	}

	go run.CommitApproval()
	go run.WaitDeploy()
	go run.Deploying()
	go run.Deployed()
	go run.Testing()
	go run.Tested()
	router.Run()
}
