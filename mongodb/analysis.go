package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"goDeploy/config"
	"time"
)

const analysisCollection = "analysis"

// Analysis 统计表结构体
type Analysis struct {
	AppName            string `bson:"app_name"`
	DeployCount        int    `bson:"deploy_count"`
	DeployFailedCount  int    `bson:"deploy_failed_count"`
	DeploySuccessCount int    `bson:"deploy_success_count"`
	DeployLastTime     int64  `bson:"deploy_last_time"`
}

// AnalysisManager 统计表管理器
type AnalysisManager struct {
	BaseDBManager
}

// NewAnalysisManager 初始化统计表管理器
func NewAnalysisManager(uri string) *AnalysisManager {
	dbConfig := &DBConfig{
		URI:                    uri,
		ConnectTimeout:         10 * time.Second,
		ServerSelectionTimeout: 5 * time.Second,
		HeartbeatInterval:      10 * time.Second,
		MaxPoolSize:            100,
		MinPoolSize:            5,
		RetryWrites:            true,
		DatabaseName:           config.ApprovalDb,
	}

	db, err := NewDBManager(dbConfig)
	if err != nil {
		panic(err)
	}

	return &AnalysisManager{
		BaseDBManager: BaseDBManager{Client: db.Client, DB: db.DB},
	}
}

// AddAnalysis 添加统计数据
func (m *AnalysisManager) AddAnalysis(analysis Analysis) (primitive.ObjectID, error) {
	coll := m.DB.Collection(analysisCollection)
	result, err := coll.InsertOne(context.TODO(), analysis)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// DeleteAnalysis 删除统计数据
func (m *AnalysisManager) DeleteAnalysis(appName string) (int64, error) {
	coll := m.DB.Collection(analysisCollection)
	result, err := coll.DeleteOne(context.TODO(), bson.M{"app_name": appName})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

// GetAnalysis 查询统计数据
func (m *AnalysisManager) GetAnalysis(appName string) (*Analysis, error) {
	coll := m.DB.Collection(analysisCollection)
	var analysis Analysis
	err := coll.FindOne(context.TODO(), bson.M{"app_name": appName}).Decode(&analysis)
	if err != nil {
		return nil, err
	}
	return &analysis, nil
}

// UpdateAnalysis 修改统计数据
func (m *AnalysisManager) UpdateAnalysis(appName string, update bson.M) (int64, error) {
	coll := m.DB.Collection(analysisCollection)
	result, err := coll.UpdateOne(context.TODO(), bson.M{"app_name": appName}, bson.M{"$set": update})
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

// IncrementDeployCount 增加发布次数
func (m *AnalysisManager) IncrementDeployCount(appName string) error {
	coll := m.DB.Collection(analysisCollection)
	_, err := coll.UpdateOne(context.TODO(), bson.M{"app_name": appName}, bson.M{"$inc": bson.M{"deploy_count": 1}})
	return err
}

// IncrementDeploySuccessCount 增加成功发布次数
func (m *AnalysisManager) IncrementDeploySuccessCount(appName string) error {
	coll := m.DB.Collection(analysisCollection)
	_, err := coll.UpdateOne(context.TODO(), bson.M{"app_name": appName}, bson.M{"$inc": bson.M{"deploy_success_count": 1}})
	return err
}

// IncrementDeployFailedCount 增加失败发布次数
func (m *AnalysisManager) IncrementDeployFailedCount(appName string) error {
	coll := m.DB.Collection(analysisCollection)
	_, err := coll.UpdateOne(context.TODO(), bson.M{"app_name": appName}, bson.M{"$inc": bson.M{"deploy_failed_count": 1}})
	return err
}

// UpdateDeployLastTime 更新最近一次发布时间
func (m *AnalysisManager) UpdateDeployLastTime(appName string) error {
	coll := m.DB.Collection(analysisCollection)
	currentTime := time.Now().Unix()
	_, err := coll.UpdateOne(context.TODO(), bson.M{"app_name": appName}, bson.M{"$set": bson.M{"deploy_last_time": currentTime}})
	return err
}

// Close 关闭MongoDB连接
func (m *AnalysisManager) Close() error {
	return m.Client.Disconnect(context.TODO())
}
