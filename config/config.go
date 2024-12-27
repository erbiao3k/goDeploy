package config

import (
	"os"
	"strings"
)

const (
	ContentType = "application/json"

	ApprovalDb = "goDeploy"
)

var (
	WecomRobotAddr   = strings.Split(os.Getenv("WECOM_ROBOT_ADDR"), ",")
	MongodbUri       = os.Getenv("MONGODB_URI")
	ProviderSecret   = os.Getenv("PROVIDER_SECRET")
	DeployTemplateId = os.Getenv("DEPLOY_TEMPLATE_ID")
	Corpid           = os.Getenv("CORP_ID")
	AutoNodeList     = os.Getenv("AUTO_NODE_LIST")
	AutoNodePeople   = strings.Split(AutoNodeList, ",")
)
