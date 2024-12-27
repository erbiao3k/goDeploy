package mongodb

import (
	"goDeploy/config"
)

var (
	MyAnalysisManager = NewAnalysisManager(config.MongodbUri)

	MyApprovalManager = NewApprovalManager(config.MongodbUri)

	MyAppStackManager = NewAppStackManager(config.MongodbUri)

	MyTokenManager = NewTokenManager(config.MongodbUri)

	MyUserManager = NewUserManager(config.MongodbUri)
)
